#!/usr/bin/env python3
"""Comprueba el contraste WCAG de los pares texto/fondo de cada tema.

Solo analiza valores hexadecimales; los alias y `color-mix` se omiten para no
inventar resultados. Sale con código 1 si algún par no llega a 4.5:1 (AA para
texto normal), que es el mínimo que usa la interfaz.
"""

import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
TOKENS = ROOT / "frontend" / "src" / "styles" / "tokens.css"

PAIRS = [
    ("--md-sys-color-on-background", "--md-sys-color-background"),
    ("--md-sys-color-on-surface", "--md-sys-color-surface"),
    ("--md-sys-color-on-surface", "--md-sys-color-surface-container"),
    ("--md-sys-color-on-surface", "--md-sys-color-surface-container-high"),
    ("--md-sys-color-on-surface-variant", "--md-sys-color-surface"),
    ("--md-sys-color-on-primary", "--md-sys-color-primary"),
    ("--md-sys-color-on-secondary", "--md-sys-color-secondary"),
    ("--md-sys-color-on-secondary-container", "--md-sys-color-secondary-container"),
    ("--md-sys-color-on-tertiary", "--md-sys-color-tertiary"),
    ("--md-sys-color-on-tertiary-container", "--md-sys-color-tertiary-container"),
    ("--md-sys-color-on-error", "--md-sys-color-error"),
    ("--md-sys-color-on-error-container", "--md-sys-color-error-container"),
]

BLOCKS = {
    "default/light": ':root, [data-theme="default"][data-color-scheme="light"]',
    "default/dark": '[data-color-scheme="dark"], [data-theme="default"][data-color-scheme="dark"]',
    "ayu/light": '[data-theme="ayu"][data-color-scheme="light"]',
    "ayu/dark": '[data-theme="ayu"][data-color-scheme="dark"]',
    "monokai/light": '[data-theme="monokai"][data-color-scheme="light"]',
    "monokai/dark": '[data-theme="monokai"][data-color-scheme="dark"]',
    "nord/light": '[data-theme="nord"][data-color-scheme="light"]',
    "nord/dark": '[data-theme="nord"][data-color-scheme="dark"]',
}

HEX = re.compile(r"^#(?:[0-9a-fA-F]{3}|[0-9a-fA-F]{6})$")


def parse_blocks(css):
    out, i = [], 0
    while i < len(css):
        m = re.search(r"([^{}]*)\{", css[i:])
        if not m:
            break
        selector = " ".join(re.sub(r"/\*.*?\*/", "", m.group(1), flags=re.S).split())
        start = i + m.end()
        depth, j = 1, start
        while j < len(css) and depth:
            if css[j] == "{":
                depth += 1
            elif css[j] == "}":
                depth -= 1
            j += 1
        body = css[start : j - 1]
        values = dict(re.findall(r"(--[a-z0-9-]+)\s*:\s*([^;]+);", body))
        out.append((selector, values))
        i = j
    return out


def rgb(value):
    v = value.strip()
    if not HEX.match(v):
        return None
    v = v[1:]
    if len(v) == 3:
        v = "".join(ch * 2 for ch in v)
    return tuple(int(v[k : k + 2], 16) for k in (0, 2, 4))


def luminance(color):
    def channel(c):
        c = c / 255
        return c / 12.92 if c <= 0.03928 else ((c + 0.055) / 1.055) ** 2.4

    r, g, b = (channel(c) for c in color)
    return 0.2126 * r + 0.7152 * g + 0.0722 * b


def contrast(a, b):
    la, lb = luminance(a), luminance(b)
    return (max(la, lb) + 0.05) / (min(la, lb) + 0.05)


def main():
    css = TOKENS.read_text()
    blocks = dict(parse_blocks(css))
    failures, skipped = [], 0

    for name, selector in BLOCKS.items():
        values = blocks.get(selector)
        if values is None:
            failures.append((name, "<bloque ausente>", 0))
            continue

        for fg, bg in PAIRS:
            fg_rgb, bg_rgb = rgb(values.get(fg, "")), rgb(values.get(bg, ""))
            if fg_rgb is None or bg_rgb is None:
                skipped += 1
                continue
            ratio = contrast(fg_rgb, bg_rgb)
            if ratio < 4.5:
                failures.append((name, f"{fg} sobre {bg}", ratio))

    if failures:
        print("Pares por debajo de 4.5:1:")
        for name, pair, ratio in failures:
            print(f"  {name:14} {pair:70} {ratio:.2f}:1")
        print(f"\n{len(failures)} fallos, {skipped} pares omitidos por no ser hex.")
        return 1

    print(f"Contraste OK en {len(BLOCKS)} bloques ({skipped} pares omitidos).")
    return 0


if __name__ == "__main__":
    sys.exit(main())
