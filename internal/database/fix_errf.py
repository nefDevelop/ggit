import re
import glob

for fp in glob.glob('**/*.go', recursive=True):
    if fp.endswith('_test.go'):
        continue
    with open(fp) as f:
        content = f.read()
    
    # Replace only fmt.Errorf("...: %v", err) patterns
    # where %v is the last format verb before the closing quote
    # and err is the last argument
    new_content = re.sub(
        r'fmt\.Errorf\(("[^"]*):\s*%v(", err\))',
        r'fmt.Errorf(\1: %w\2',
        content
    )
    
    if new_content != content:
        with open(fp, 'w') as f:
            f.write(new_content)
        print(f'  {fp}')

print('Done')
