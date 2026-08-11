const { execSync } = require('child_process');
try {
  const r = execSync('wmic logicaldisk get DeviceID,FreeSpace,Size /format:csv', { encoding: 'utf8' });
  console.log(r);
} catch(e) {
  console.error(e.message);
}
