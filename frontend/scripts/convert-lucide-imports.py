#!/usr/bin/env python3
"""Convert barrel imports from lucide-svelte to tree-shakeable per-icon imports."""

import re
import sys
from pathlib import Path

SRC_DIR = Path(__file__).resolve().parent.parent / "src"

def pascal_to_kebab(name: str) -> str:
    """Convert PascalCase to kebab-case: GitPullRequest -> git-pull-request, Trash2 -> trash-2"""
    name = re.sub(r'(?<=[a-z])(?=[A-Z])', '-', name)
    name = re.sub(r'(?<=[A-Z])(?=[A-Z][a-z])', '-', name)
    name = re.sub(r'(?<=[a-z])(?=[0-9])', '-', name)
    return name.lower()

def process_file(filepath: Path) -> bool:
    content = filepath.read_text()
    
    # Match: import { Icon1, Icon2, ... } from "lucide-svelte"
    pattern = re.compile(
        r'import\s*\{([^}]+)\}\s*from\s*["\']lucide-svelte["\']\s*;?'
    )
    
    if not pattern.search(content):
        return False
    
    def replace_import(match):
        icon_names = [name.strip() for name in match.group(1).split(",")]
        new_imports = []
        for name in icon_names:
            kebab = pascal_to_kebab(name)
            new_imports.append(f'import {name} from "lucide-svelte/icons/{kebab}";')
        return "\n" + "\n".join(new_imports) + "\n"
    
    new_content = pattern.sub(replace_import, content)
    filepath.write_text(new_content)
    return True

def main():
    modified = 0
    for f in sorted(SRC_DIR.rglob("*.svelte")):
        if process_file(f):
            print(f"  {f.relative_to(SRC_DIR.parent)}")
            modified += 1
    
    for f in sorted(SRC_DIR.rglob("*.ts")):
        if process_file(f):
            print(f"  {f.relative_to(SRC_DIR.parent)}")
            modified += 1
    
    print(f"\nModified {modified} files")

if __name__ == "__main__":
    main()
