const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

const srcDir = path.join(__dirname, '..', 'src');

// Find all individual imports from lucide-svelte/icons/
const result = execSync(`grep -rn 'lucide-svelte/icons/' ${srcDir} --include='*.svelte' --include='*.ts' -l`, {encoding: 'utf8'});
const files = result.trim().split('\n').filter(Boolean);

console.log(`Found ${files.length} files with individual icon imports`);

// For each file, find all lines with lucide-svelte/icons/ imports,
// collect them, and replace with a single barrel import
let modified = 0;
for (const filepath of files) {
  let content = fs.readFileSync(filepath, 'utf8');
  
  // Match all import lines: `import X from "lucide-svelte/icons/x";`
  const iconLines = content.match(/import \w+ from "lucide-svelte\/icons\/[^"]+";/g);
  if (!iconLines || iconLines.length <= 1) continue;
  
  // Extract icon names
  const names = iconLines.map(line => line.match(/import (\w+) from/)[1]);
  
  // Create barrel import
  const barrelImport = `import { ${names.join(', ')} } from "lucide-svelte";`;
  
  // Calculate the exact range to replace
  // Find the first and last icon import line
  const firstIndex = content.indexOf(iconLines[0]);
  const lastLine = iconLines[iconLines.length - 1];
  const lastIndex = content.indexOf(lastLine) + lastLine.length;
  
  // Get the text between first and last (including newlines)
  const before = content.substring(0, firstIndex);
  const after = content.substring(lastIndex);
  
  content = before + barrelImport + after;
  fs.writeFileSync(filepath, content);
  modified++;
}

console.log(`Reverted ${modified} files to barrel imports`);
