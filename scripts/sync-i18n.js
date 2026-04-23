const fs = require('fs');
const path = require('path');

const localesDir = path.join(__dirname, '../frontend/src/locales');
const sourceFile = path.join(localesDir, 'en-US.json');

if (!fs.existsSync(sourceFile)) {
  console.error('Source file en-US.json not found!');
  process.exit(1);
}

const sourceData = JSON.parse(fs.readFileSync(sourceFile, 'utf8'));

function deepMerge(source, target) {
  for (const key in source) {
    if (source[key] instanceof Object && !Array.isArray(source[key])) {
      if (!target[key]) target[key] = {};
      deepMerge(source[key], target[key]);
    } else {
      if (target[key] === undefined) {
        target[key] = source[key];
      }
    }
  }
}

function syncLocales() {
  const files = fs.readdirSync(localesDir).filter(f => f.endsWith('.json') && f !== 'en-US.json');

  files.forEach(file => {
    const filePath = path.join(localesDir, file);
    try {
      const targetData = JSON.parse(fs.readFileSync(filePath, 'utf8'));
      
      console.log(`Syncing ${file}...`);
      deepMerge(sourceData, targetData);

      // Sort keys alphabetically for better git diffs
      const sortedData = {};
      Object.keys(targetData).sort().forEach(key => {
        sortedData[key] = targetData[key];
      });

      fs.writeFileSync(filePath, JSON.stringify(sortedData, null, 2) + '\n');
    } catch (e) {
      console.error(`Error syncing ${file}:`, e.message);
    }
  });

  console.log('Done! All locales are synced with en-US.json');
}

syncLocales();
