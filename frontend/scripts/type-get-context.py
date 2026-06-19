import re
import glob

files = glob.glob('src/**/*.svelte', recursive=True) + glob.glob('src/**/*.ts', recursive=True)

for fp in sorted(files):
    with open(fp) as f:
        content = f.read()
    
    # Match: `const repo: any = getContext('repo')` or `const repo: any = getContext("repo")`
    # Also handles `const repoContext = getContext<any>('repo')`
    old = r"const\s+(\w+)\s*:\s*any\s*=\s*getContext(?:<any>)?\s*\(\s*['\"]repo['\"]\s*\)"
    new = r"const \1 = getContext<RepoContext>('repo')"
    
    new_content = re.sub(old, new, content)
    
    # Check if we need to add the import
    if new_content != content:
        # Add import if not present
        if "import type { RepoContext } from '$lib/types';" not in new_content and \
           'import type { RepoContext } from "$lib/types";' not in new_content:
            # Add after last existing import or script tag
            new_content = new_content.replace(
                '<script lang="ts">',
                '<script lang="ts">\n  import type { RepoContext } from "$lib/types";',
                1
            )
        
        with open(fp, 'w') as f:
            f.write(new_content)
        print(f'  {fp}')
