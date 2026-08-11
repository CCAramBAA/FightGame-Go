// Pure Node.js self-signed certificate generator (zero dependencies)
'use strict';
const crypto = require('crypto');
const fs = require('fs');
const path = require('path');

const SSL_DIR = path.join(__dirname, '..', 'docker', 'nginx', 'ssl');

// ============================================================
// ASN.1 / DER helpers (minimal, just enough for X.509 certs)
// ============================================================

function derLength(len) {
  if (len < 0x80) return Buffer.from([len]);
  const octets = [];
  let l = len;
  while (l > 0) { octets.unshift(l & 0xff); l >>>= 8; }
  return Buffer.concat([Buffer.from([0x80 | octets.length]), Buffer.from(octets)]);
}

function derInteger(buf) {
  // Remove leading zeros (but keep one if MSB is set to avoid negative)
  let i = 0;
  while (i < buf.length - 1 && buf[i] === 0 && !(buf[i + 1] & 0x80)) i++;
  const val = buf.subarray(i);
  return Buffer.concat([Buffer.from([0x02]), derLength(val.length), val]);
}

function derBitString(buf) {
  return Buffer.concat([Buffer.from([0x03]), derLength(buf.length + 1), Buffer.from([0x00]), buf]);
}

function derOctetString(buf) {
  return Buffer.concat([Buffer.from([0x04]), derLength(buf.length), buf]);
}

function derOID(oid) {
  const parts = oid.split('.').map(Number);
  const vals = [40 * parts[0] + parts[1]];
  for (let i = 2; i < parts.length; i++) {
    let v = parts[i];
    if (v < 128) { vals.push(v); continue; }
    const bytes = [];
    while (v > 0) { bytes.unshift(v & 0x7f); v >>>= 7; }
    for (let j = 0; j < bytes.length - 1; j++) bytes[j] |= 0x80;
    vals.push(...bytes);
  }
  return Buffer.concat([Buffer.from([0x06]), derLength(vals.length), Buffer.from(vals)]);
}

function derNull() { return Buffer.from([0x05, 0x00]); }

function derUTCTime(date) {
  const s = date.getUTCFullYear().toString().slice(2).padStart(2, '0') +
    (date.getUTCMonth() + 1).toString().padStart(2, '0') +
    date.getUTCDate().toString().padStart(2, '0') +
    date.getUTCHours().toString().padStart(2, '0') +
    date.getUTCMinutes().toString().padStart(2, '0') +
    date.getUTCSeconds().toString().padStart(2, '0') + 'Z';
  return Buffer.concat([Buffer.from([0x17]), derLength(s.length), Buffer.from(s)]);
}

function derSequence(items) {
  const body = Buffer.concat(items);
  return Buffer.concat([Buffer.from([0x30]), derLength(body.length), body]);
}

function derSet(items) {
  const body = Buffer.concat(items);
  return Buffer.concat([Buffer.from([0x31]), derLength(body.length), body]);
}

function derName(entries) {
  // Each entry: { oid, value }
  const rdns = entries.map(e => {
    const attr = Buffer.concat([
      derOID(e.oid),
      Buffer.from([0x0c]), derLength(Buffer.byteLength(e.value)), Buffer.from(e.value) // UTF8String
    ]);
    return derSet([derSequence([attr])]);
  });
  return derSequence(rdns);
}

function derSubjectAltNames(names) {
  const entries = names.map(n => {
    if (n.startsWith('DNS:')) {
      const dns = n.slice(4);
      return Buffer.concat([Buffer.from([0x82]), derLength(dns.length), Buffer.from(dns)]);
    }
    if (n.startsWith('IP:')) {
      const ip = n.slice(3);
      const parts = ip.split('.').map(Number);
      return Buffer.concat([Buffer.from([0x87, 0x04]), Buffer.from(parts)]);
    }
    throw new Error('Unknown SAN type: ' + n);
  });
  return derOctetString(Buffer.concat(entries));
}

function derExtensions(extensions) {
  // extensions: [{ oid, critical, value }]
  const items = extensions.map(e => {
    const ext = Buffer.concat([
      derOID(e.oid),
      ...(e.critical ? [Buffer.from([0x01, 0x01, 0xff])] : []), // BOOLEAN TRUE = FF
      derOctetString(e.value)
    ]);
    return derSequence([ext]);
  });
  return derSequence(items);
}

// ============================================================
// PEM encode
// ============================================================
function toPEM(label, derBuffer) {
  const b64 = derBuffer.toString('base64');
  const lines = [];
  for (let i = 0; i < b64.length; i += 64) lines.push(b64.slice(i, i + 64));
  return `-----BEGIN ${label}-----\n${lines.join('\n')}\n-----END ${label}-----\n`;
}

// ============================================================
// Main: generate key + self-signed cert
// ============================================================

function generate() {
  console.log('Generating RSA 2048-bit key pair...');
  
  // Generate key pair
  const { privateKey: privPEM, publicKey: pubPEM } = crypto.generateKeyPairSync('rsa', {
    modulusLength: 2048,
    publicKeyEncoding: { type: 'spki', format: 'der' },
    privateKeyEncoding: { type: 'pkcs8', format: 'pem' },
  });

  // Write private key
  fs.writeFileSync(path.join(SSL_DIR, 'server.key'), privPEM);
  console.log('  -> server.key');

  // Parse public key DER to extract modulus + exponent
  // SPKI: SEQUENCE { AlgorithmIdentifier, BIT STRING { RSAPublicKey } }
  let pos = 2; // skip SEQUENCE + length
  // skip AlgorithmIdentifier
  const algLen = pubPEM[pos + 1];
  pos += 2 + algLen;
  // BIT STRING: tag=03, len, unusedBits, SEQUENCE
  pos += 3; // skip BIT STRING tag + len + unusedBits byte
  // SEQUENCE { INTEGER modulus, INTEGER exponent }
  pos += 2; // skip SEQUENCE tag + len
  // INTEGER modulus
  const modLen = pubPEM[pos + 1];
  const modulus = pubPEM.subarray(pos + 2, pos + 2 + modLen);
  pos += 2 + modLen;
  // INTEGER exponent
  const expLen = pubPEM[pos + 1];
  const exponent = pubPEM.subarray(pos + 2, pos + 2 + expLen);

  // Build RSAPublicKey SEQUENCE
  const rsaPubKey = derSequence([derInteger(modulus), derInteger(exponent)]);
  
  // AlgorithmIdentifier: OID 1.2.840.113549.1.1.1 (rsaEncryption) + NULL
  const algId = derSequence([derOID('1.2.840.113549.1.1.1'), derNull()]);
  
  // SubjectPublicKeyInfo
  const spki = derSequence([algId, derBitString(rsaPubKey)]);

  // Serial number (random 16 bytes)
  const serial = crypto.randomBytes(16);
  // Ensure MSB is 0 to keep it positive in DER
  serial[0] &= 0x7f;
  
  // Validity
  const notBefore = new Date();
  const notAfter = new Date();
  notAfter.setFullYear(notBefore.getFullYear() + 1);
  
  // Subject & Issuer (same for self-signed)
  const name = derName([
    { oid: '2.5.4.6',  value: 'CN' },       // countryName
    { oid: '2.5.4.8',  value: 'Beijing' },   // stateOrProvinceName
    { oid: '2.5.4.7',  value: 'Beijing' },   // localityName
    { oid: '2.5.4.10', value: 'FightGame' }, // organizationName
    { oid: '2.5.4.3',  value: 'localhost' }, // commonName
  ]);
  
  // Subject Alt Names extension
  const sanExtValue = derSubjectAltNames(['DNS:localhost', 'IP:127.0.0.1']);
  
  // Basic Constraints: CA:TRUE
  const basicConstraints = derSequence([Buffer.from([0x01, 0x01, 0xff])]); // BOOLEAN TRUE
  
  const extensions = derExtensions([
    { oid: '2.5.29.19', critical: true,  value: basicConstraints }, // basicConstraints
    { oid: '2.5.29.17', critical: false, value: sanExtValue },      // subjectAltName
  ]);

  // TBSCertificate
  const version = Buffer.from([0xa0, 0x03, 0x02, 0x01, 0x02]); // [0] EXPLICIT INTEGER 2 (v3)
  const tbsCert = derSequence([
    version,
    derInteger(serial),
    algId,          // signature algorithm
    name,           // issuer
    derSequence([derUTCTime(notBefore), derUTCTime(notAfter)]), // validity
    name,           // subject
    spki,           // subjectPublicKeyInfo
    Buffer.from([0xa3]), derLength(extensions.length), extensions, // [3] EXPLICIT extensions
  ]);

  // Sign with private key
  const privKeyObj = crypto.createPrivateKey(privPEM);
  const signature = crypto.createSign('RSA-SHA256').update(tbsCert).sign(privKeyObj);

  // SignatureAlgorithm same as algorithm in TBSCert
  const sigAlg = derSequence([derOID('1.2.840.113549.1.1.11'), derNull()]); // sha256WithRSAEncryption
  
  // Complete certificate
  const cert = derSequence([
    tbsCert,
    sigAlg,
    derBitString(signature),
  ]);

  fs.writeFileSync(path.join(SSL_DIR, 'server.crt'), toPEM('CERTIFICATE', cert));
  console.log('  -> server.crt');
  console.log('✅ SSL certificate generated successfully!');
}

try {
  generate();
} catch (err) {
  console.error('❌ Failed:', err.message);
  process.exit(1);
}
