import re

def fix_file(fp):
    with open(fp) as f:
        lines = f.readlines()
    
    changed = False
    for i, line in enumerate(lines):
        stripped = line.strip()
        
        # Match log.Error|Fatal|Warn|Trace("%v" followed by args
        m = re.match(r'(log\.(?:Error|Fatal|Warn|Trace))\("%v"(.*)', stripped)
        if m:
            logfn = m.group(1)
            rest = m.group(2)
            
            # Look at surrounding code for context keywords
            context = ""
            for j in range(max(0, i-5), min(len(lines), i+3)):
                if j != i:
                    ctx_line = lines[j].strip()
                    for kw in ['save', 'create', 'update', 'delete', 'load', 'fetch', 'prepare',
                               'merge', 'push', 'sync', 'clean', 'init', 'migrate', 'import',
                               'export', 'validate', 'check', 'verify', 'send', 'receive',
                               'attachment', 'commit', 'branch', 'repository', 'issue', 'pull',
                               'webhook', 'release', 'milestone', 'action', 'comment', 'notice',
                               'access', 'token', 'user', 'key', 'org', 'team', 'hook',
                               'git', 'fsck', 'gc', 'archive', 'blob', 'tree', 'tag']:
                        if kw in ctx_line.lower():
                            context = kw
                            break
            
            if not context:
                context = 'operation'
            
            indent = line[:len(line) - len(line.lstrip())]
            lines[i] = f'{indent}{logfn}("failed to {context}: %v"{rest}\n'
            changed = True
    
    if changed:
        with open(fp, 'w') as f:
            f.writelines(lines)
        return True
    return False

import glob
files = set()
for fp in glob.glob('/home/nef734/Documentos/Proyectos/ggit/internal/**/*.go', recursive=True):
    if '_test.go' in fp:
        continue
    if fix_file(fp):
        files.add(fp)

for f in sorted(files):
    print(f'  {f.replace("/home/nef734/Documentos/Proyectos/ggit/", "")}')
print(f'\nFixed {len(files)} files')
