import os
import re

def fix_imports(file_path):
    with open(file_path, 'r', encoding='utf-8') as f:
        lines = f.readlines()

    content = "".join(lines)
    
    # Check usages
    uses_log = bool(re.search(r'\blog\.(Trace|Info|Warn|Error|Fatal|NewConsole|ErrorDepth|Success|LevelTrace|LevelInfo|LevelWarn|LevelError|LevelFatal|Handle)\b', content))
    uses_errors = bool(re.search(r'\berrors\.(New|Is|As|Unwrap)\b', content))
    uses_fmt = bool(re.search(r'\bfmt\.\b', content))

    # Better approach: parse imports, keep track of them, and regenerate.
    
    new_lines = []
    in_import_block = False
    imports = []
    
    # We'll use a simple state machine to find and extract imports
    other_lines_before = []
    other_lines_after = []
    found_imports = False
    finished_imports = False
    
    i = 0
    while i < len(lines):
        line = lines[i]
        if not found_imports:
            if line.strip().startswith('import ('):
                in_import_block = True
                found_imports = True
                i += 1
                continue
            elif line.strip().startswith('import "'):
                found_imports = True
                finished_imports = True
                imports.append(line.strip()[len('import '):])
                i += 1
                continue
            elif line.strip().startswith('import'):
                 # handle import log "..."
                 match = re.match(r'import\s+(\w+)\s+"(.*?)"', line.strip())
                 if match:
                     found_imports = True
                     finished_imports = True
                     imports.append(f'{match.group(1)} "{match.group(2)}"')
                     i += 1
                     continue
            other_lines_before.append(line)
        elif in_import_block:
            if line.strip() == ')':
                in_import_block = False
                finished_imports = True
                i += 1
                continue
            if line.strip():
                # Split multiple imports on same line if any (though rare in Go)
                imports.append(line.strip())
        else:
            other_lines_after.append(line)
        i += 1

    # Clean up existing imports
    new_imports = []
    for imp in imports:
        # Normalize
        imp = imp.strip()
        if not imp: continue
        
        # Remove cockroachdb/errors
        if "github.com/cockroachdb/errors" in imp:
            continue
        
        # We will re-add needed ones later
        if imp == '"errors"' or imp == 'errors "errors"': continue
        if imp == '"fmt"' or imp == 'fmt "fmt"': continue
        if '"github.com/nef734/ggit/internal/log"' in imp: continue
        if imp == '"log"' or imp == 'log "log"': continue
        
        if imp not in new_imports:
            new_imports.append(imp)
    
    # Re-add needed imports
    if uses_log:
        new_imports.append('log "github.com/nef734/ggit/internal/log"')
    if uses_errors:
        new_imports.append('"errors"')
    if uses_fmt:
        if '"fmt"' not in new_imports:
            new_imports.append('"fmt"')
    
    # Sort imports roughly: stdlib first, then third party
    stdlib = []
    third_party = []
    for imp in new_imports:
        path = imp.split()[-1].strip('"')
        # Naive stdlib check
        if "." not in path.split('/')[0] and "github.com/nef734/ggit" not in imp:
             stdlib.append(imp)
        else:
            third_party.append(imp)
            
    stdlib.sort()
    third_party.sort()
    
    final_imports = stdlib + ([''] if stdlib and third_party else []) + third_party
    
    # Reconstruct file
    new_content = "".join(other_lines_before)
    if final_imports:
        new_content += "import (\n"
        for imp in final_imports:
            if imp == '':
                new_content += "\n"
            else:
                new_content += f"\t{imp}\n"
        new_content += ")\n"
    new_content += "".join(other_lines_after)
    
    if new_content != content:
        with open(file_path, 'w', encoding='utf-8') as f:
            f.write(new_content)
        return True
    return False

# Main loop
changed_files = 0
for root, dirs, files in os.walk('.'):
    if '.git' in dirs: dirs.remove('.git')
    if 'vendor' in dirs: dirs.remove('vendor')
    if 'node_modules' in dirs: dirs.remove('node_modules')
    
    for file in files:
        if file.endswith('.go'):
            file_path = os.path.join(root, file)
            try:
                if fix_imports(file_path):
                    changed_files += 1
                    print(f"Fixed imports in {file_path}")
            except Exception as e:
                print(f"Error processing {file_path}: {e}")

print(f"Total files changed: {changed_files}")
