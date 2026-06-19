const fs = require('fs');
const path = require('path');

const iconDir = path.join(__dirname, '..', 'node_modules', 'lucide-svelte', 'dist', 'icons');
const actualIcons = new Set(
  fs.readdirSync(iconDir)
    .filter(f => f.endsWith('.js'))
    .map(f => f.replace(/\.js$/, ''))
);

// Find all import paths in our code
const srcDir = path.join(__dirname, '..', 'src');
const { execSync } = require('child_process');
const result = execSync(`grep -r 'lucide-svelte/icons/' ${srcDir} --include='*.svelte' --include='*.ts' -h`);
const importPaths = result.toString().match(/lucide-svelte\/icons\/[^"']+/g) || [];
const usedIcons = [...new Set(importPaths.map(s => s.replace('lucide-svelte/icons/', '')))];

const missing = usedIcons.filter(i => !actualIcons.has(i));
if (missing.length > 0) {
  console.log('MISSING ICONS (' + missing.length + '):');
  missing.forEach(i => console.log('  ' + i));
} else {
  console.log('All ' + usedIcons.length + ' icon imports resolve correctly');
}
