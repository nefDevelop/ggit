// @vitest-environment node
import { describe, it, expect } from 'vitest';
import { readFileSync, readdirSync, statSync } from 'node:fs';
import { join, extname, basename } from 'node:path';

const SRC = join(process.cwd(), 'src');

/**
 * Components allowed to hold raw colour literals, with the reason.
 *
 * These are local palettes rather than theme colours. Tokenising them would
 * break the thing they exist to provide, so they are listed instead of being
 * rewritten:
 *
 *  - ThemeSelector / ThemeSwitch define the themes. Assigning
 *    var(--color-warning) to --color-warning would be circular, and the per
 *    theme colours are the input the picker edits.
 *  - Identicon derives a deterministic colour per user. That colour is data.
 *  - MarkdownRenderer carries syntax highlighting colours.
 *  - CodeEditor / CodeViewer carry the dark editor palette.
 */
const allowed = new Set([
  'ThemeSelector.svelte',
  'ThemeSwitch.svelte',
  'Identicon.svelte',
  'MarkdownRenderer.svelte',
  'CodeEditor.svelte',
  'CodeViewer.svelte',
]);

/**
 * Literals that are allowed even outside the palette files, with the reason.
 *
 * #fff and #000 are the achromatic extremes. A token for them adds indirection
 * without adding meaning, and the palette already resolves the cases that
 * matter through --md-sys-color-on-* where white sits on a tinted surface.
 */
const allowedLiterals = new Set(['#fff', '#fff8', '#ffffff', '#000', '#0008', '#000000']);

/** Colour literals that are values rather than styling, not theme colours. */
const valuePatterns: RegExp[] = [
  // A label or column colour the user picks, stored as data.
  /let\s+\w*[cC]olor\w*\s*=\s*\$state\(['"]#[0-9a-fA-F]{3,8}['"]\)/,
  /\b\w*[cC]olor\w*\s*=\s*['"]#[0-9a-fA-F]{3,8}['"]/,
  // A colour returned as a value, such as the Kanban column palette.
  /return\s+['"]#[0-9a-fA-F]{3,8}['"]/,
  // A placeholder telling the user what to type.
  /placeholder=['"]#[0-9a-fA-F]{3,8}['"]/,
  // Writing a theme token at runtime, which is how a theme is applied.
  /setProperty\(\s*['"]--[\w-]+['"]\s*,\s*['"]#[0-9a-fA-F]{3,8}['"]/,
];

function walk(dir: string, out: string[] = []): string[] {
  for (const entry of readdirSync(dir)) {
    const full = join(dir, entry);
    if (statSync(full).isDirectory()) walk(full, out);
    else if (extname(full) === '.svelte') out.push(full);
  }
  return out;
}

describe('colour literals', () => {
  const files = walk(SRC).filter((f) => !allowed.has(basename(f)));

  it('finds components to check', () => {
    expect(files.length).toBeGreaterThan(50);
  });

  it('uses tokens instead of raw colours', () => {
    const offenders: string[] = [];

    for (const file of files) {
      const lines = readFileSync(file, 'utf8').split('\n');
      lines.forEach((line, i) => {
        for (const m of line.matchAll(/#[0-9a-fA-F]{3,8}\b/g)) {
          if (allowedLiterals.has(m[0].toLowerCase())) continue;
          if (valuePatterns.some((p) => p.test(line))) continue;
          offenders.push(`${basename(file)}:${i + 1}  ${m[0]}  |  ${line.trim().slice(0, 70)}`);
        }
      });
    }

    // The palette above used to be about 170 literals across sixty components.
    // Status colours went to --status-* and the dark output panels to
    // --output-*. A new literal means a colour that will not follow the theme.
    expect(
      offenders,
      `These colours will not follow the theme. Add a token, or add the file to the\n` +
        `allowlist if it is a local palette that has to keep its own values:\n  ${offenders.join('\n  ')}`,
    ).toEqual([]);
  });

  it('defines every token the components reference', () => {
    const css = ['tokens.css', 'base.css', 'components.css', 'layout.css', 'utils.css', 'tables.css']
      .map((f) => readFileSync(join(SRC, 'styles', f), 'utf8'))
      .join('\n');

    // Tokens declared inside a component count too. PageHeader, for instance,
    // declares its own icon sizes locally, and failing those would push authors
    // to hoist perfectly reasonable local values into the global sheet.
    const local = new Set<string>();
    for (const file of walk(SRC)) {
      const src = readFileSync(file, 'utf8');
      for (const m of src.matchAll(/(^|[;{\s])(--[\w-]+)\s*:/gm)) local.add(m[2]);
    }

    const declared = new Set([
      ...[...css.matchAll(/(--[\w-]+)\s*:/g)].map((m) => m[1]),
      ...local,
    ]);

    const missing = new Map<string, string[]>();
    for (const file of files) {
      const src = readFileSync(file, 'utf8');
      for (const m of src.matchAll(/var\((--[\w-]+)\)/g)) {
        if (declared.has(m[1]) || m[1].startsWith('--md-sys-color')) continue;
        if (!missing.has(m[1])) missing.set(m[1], []);
        missing.get(m[1])!.push(basename(file));
      }
    }

    const report = [...missing.entries()].map(
      ([token, where]) => `${token} (usada en ${[...new Set(where)].join(', ')})`,
    );

    // A var() with no declaration resolves to nothing and silently loses the
    // colour, which is the failure this catches.
    expect(report, `Tokens referenciados sin declarar:\n  ${report.join('\n  ')}`).toEqual([]);
  });
});