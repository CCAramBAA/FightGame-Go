const { execSync } = require('child_process');
const path = require('path');
const frontendDir = 'd:/Development/ALL-PROJECTS/FightGame/frontend';
// Use node to run the tsc CLI script
const tscJs = path.join(frontendDir, 'node_modules', 'typescript', 'bin', 'tsc');
try {
  const r = execSync(`node "${tscJs}" --noEmit --pretty`, { cwd: frontendDir, encoding: 'utf8', stdio: 'pipe', maxBuffer: 2*1024*1024 });
  console.log('ALL CLEAN - no type errors');
} catch(e) {
  const err = e.stderr || e.stdout || e.message;
  console.log(err.substring(0, 5000));
}
