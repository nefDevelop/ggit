#!/usr/bin/env python3
"""Comprueba que todos los temas definen el mismo contrato de tokens.

El tema por defecto es la referencia. Las sombras viven en el primer `:root`
global y no se exigen por tema. Exit code 1 si algún tema se queda corto, que es
lo que evita que las paletas se desincronicen al añadir un token nuevo.
"""

import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
TOKENS = ROOT / "frontend" / "src" / "styles" / "tokens.css"

THEMES = ("default", "ayu", "monokai", "nord")
GLOBAL_TOKENS = {
    "--shadow-sm",
    "--shadow-md",
    "--shadow-lg",
    "--shadow-xl",
    "--inner-shadow",
    "--island-shadow",
    "--scrim",
    "--sunken-1",
    "--sunken-2",
    "--sunken-3",
    "--shadow-popover",
    "--shadow-modal",
}

DEFAULT_LIGHT = ':root, [data-theme="default"][data-color-scheme="light"]'
DEFAULT_DARK = '[data-color-scheme="dark"], [data-theme="default"][data-color-scheme="dark"]'


def parse_blocks(css):
    blocks = []
    i = 0
    while i < len(css):
        match = re.search(r"([^{}]*)\{", css[i:])
        if not match:
            break

        selector = " ".join(re.sub(r"/\*.*?\*/", "", match.group(1), flags=re.S).split())
        start = i + match.end()
        depth = 1
        j = start
        while j < len(css) and depth:
            if css[j] == "{":
                depth += 1
            elif css[j] == "}":
                depth -= 1
            j += 1

        blocks.append((selector, css[start : j - 1]))
        i = j
    return blocks


def tokens_of(body):
    return set(re.findall(r"(--[a-z0-9-]+)\s*:", body))


def find_block(blocks, selector):
    for sel, body in blocks:
        if sel == selector:
            return body
    raise SystemExit(f"no se encontró el bloque {selector}")


def main():
    css = TOKENS.read_text()
    blocks = parse_blocks(css)

    required_light = tokens_of(find_block(blocks, DEFAULT_LIGHT)) - GLOBAL_TOKENS
    required_dark = tokens_of(find_block(blocks, DEFAULT_DARK)) - GLOBAL_TOKENS

    failures = []
    for theme in THEMES:
        for scheme, required in (("light", required_light), ("dark", required_dark)):
            selector = f'[data-theme="{theme}"][data-color-scheme="{scheme}"]'
            if theme == "default":
                selector = DEFAULT_LIGHT if scheme == "light" else DEFAULT_DARK

            try:
                body = find_block(blocks, selector)
            except SystemExit:
                failures.append((theme, scheme, ["<bloque ausente>"]))
                continue

            missing = sorted(required - tokens_of(body))
            if missing:
                failures.append((theme, scheme, missing))

    if failures:
        print("Contrato de tokens incompleto:")
        for theme, scheme, missing in failures:
            print(f"  {theme}/{scheme}: {', '.join(missing)}")
        return 1

    print(
        f"Contrato de tokens OK: {len(required_light)} tokens en claro y "
        f"{len(required_dark)} en oscuro, en {len(THEMES)} temas."
    )
    return 0


if __name__ == "__main__":
    sys.exit(main())
