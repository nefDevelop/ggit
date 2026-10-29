// @vitest-environment node
import { describe, it, expect } from 'vitest';
import { readFileSync, readdirSync, statSync } from 'node:fs';
import { join, extname } from 'node:path';

const SRC = join(process.cwd(), 'src');
const LOCALES = join(SRC, 'locales');
const SOURCE_LOCALE = 'en-US.json';

/** Collects every file under a directory, recursively. */
function walk(dir: string, out: string[] = []): string[] {
  for (const entry of readdirSync(dir)) {
    const full = join(dir, entry);
    if (statSync(full).isDirectory()) {
      walk(full, out);
    } else if (['.svelte', '.ts'].includes(extname(full))) {
      out.push(full);
    }
  }
  return out;
}

/** Flattens a nested locale object into dotted keys. */
function flatten(value: unknown, prefix = ''): Set<string> {
  const keys = new Set<string>();
  if (value && typeof value === 'object' && !Array.isArray(value)) {
    for (const [key, child] of Object.entries(value as Record<string, unknown>)) {
      const path = prefix ? `${prefix}.${key}` : key;
      if (child && typeof child === 'object' && !Array.isArray(child)) {
        for (const nested of flatten(child, path)) keys.add(nested);
      } else {
        keys.add(path);
      }
    }
  }
  return keys;
}

/** Every key passed to i18n.t() in the source. */
function collectUsedKeys(): Map<string, Set<string>> {
  const used = new Map<string, Set<string>>();
  const pattern = /i18n\.t\(\s*["']([^"']+)["']/g;

  for (const file of walk(SRC)) {
    const content = readFileSync(file, 'utf8');
    for (const match of content.matchAll(pattern)) {
      const key = match[1];
      if (!used.has(key)) used.set(key, new Set());
      used.get(key)!.add(file.replace(SRC + '/', ''));
    }
  }
  return used;
}

const sourceKeys = flatten(JSON.parse(readFileSync(join(LOCALES, SOURCE_LOCALE), 'utf8')));
const used = collectUsedKeys();

describe('i18n key coverage', () => {
  it('resolves every key the interface asks for', () => {
    const missing = [...used.entries()]
      .filter(([key]) => !sourceKeys.has(key))
      .map(([key, files]) => `${key} (usada en ${[...files].join(', ')})`);

    // t() returns the key itself when the lookup fails, so a missing entry is not
    // a blank label but the raw identifier shown to the user, for example
    // "settings.saved_success" in the toast after saving a repository.
    expect(
      missing,
      `Estas claves se usan en el codigo pero no existen en ${SOURCE_LOCALE}:\n  ${missing.join('\n  ')}`,
    ).toEqual([]);
  });

  it('finds the locale files it needs', () => {
    const locales = readdirSync(LOCALES).filter((f) => f.endsWith('.json'));
    expect(locales).toContain(SOURCE_LOCALE);
  });
});