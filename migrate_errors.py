import os
import re

files = [
    "cmd/ggit/admin.go",
    "cmd/ggit/backup.go",
    "cmd/ggit/import.go",
    "cmd/ggit/restore.go",
    "internal/auth/auth.go",
    "internal/auth/ldap/config.go",
    "internal/avatar/avatar.go",
    "internal/conf/conf.go",
    "internal/conf/log.go",
    "internal/conf/utils.go",
    "internal/context/auth.go",
    "internal/context/repo.go",
    "internal/cryptox/aes.go",
    "internal/database/access_tokens.go",
    "internal/database/actions.go",
    "internal/database/attachment.go",
    "internal/database/backup.go",
    "internal/database/backup_test.go",
    "internal/database/comment.go",
    "internal/database/database.go",
    "internal/database/gists.go",
    "internal/database/issue.go",
    "internal/database/issue_label.go",
    "internal/database/issue_mail.go",
    "internal/database/lfs.go",
    "internal/database/login_source_files.go",
    "internal/database/login_sources.go",
    "internal/database/migrations/migrations.go",
    "internal/database/migrations/v20.go",
    "internal/database/milestone.go",
    "internal/database/models.go",
    "internal/database/org.go",
    "internal/database/org_team.go",
    "internal/database/organizations.go",
    "internal/database/permissions.go",
    "internal/database/public_keys.go",
    "internal/database/pull.go",
    "internal/database/pull_review.go",
    "internal/database/pull_review_comment.go",
    "internal/database/release.go",
    "internal/database/repo.go",
    "internal/database/repo_branch.go",
    "internal/database/repo_collaboration.go",
    "internal/database/repo_editor.go",
    "internal/database/repo_tag.go",
    "internal/database/repositories.go",
    "internal/database/ssh_key.go",
    "internal/database/topics.go",
    "internal/database/two_factor.go",
    "internal/database/two_factors.go",
    "internal/database/update.go",
    "internal/database/user_block.go",
    "internal/database/users.go",
    "internal/database/webhook.go",
    "internal/dbx/dsn.go",
    "internal/email/email.go",
    "internal/email/message.go",
    "internal/errx/errx_test.go",
    "internal/fsx/transaction.go",
    "internal/gitx/diff.go",
    "internal/gitx/error.go",
    "internal/gitx/pull_request.go",
    "internal/gitx/pull_request_test.go",
    "internal/gitx/tag.go",
    "internal/iox/iox.go",
    "internal/lfsx/oid.go",
    "internal/lfsx/storage.go",
    "internal/process/manager.go",
    "internal/repox/repox.go",
    "internal/route/api/v2/repo.go",
    "internal/route/api/v2/repo_compare.go",
    "internal/route/api/v2/user_settings.go",
    "internal/route/init.go",
    "internal/route/repo/download.go",
    "internal/services/backup/backup.go",
    "internal/services/repo/adoption.go",
    "internal/services/repo/git.go",
    "internal/services/repo/repo.go",
    "internal/ssh/ssh.go",
    "internal/template/template.go",
    "internal/testx/exec.go",
    "internal/userx/userx.go"
]

def find_closing_paren(s, start):
    count = 1
    for i in range(start, len(s)):
        if s[i] == '(':
            count += 1
        elif s[i] == ')':
            count -= 1
            if count == 0:
                return i
    return -1

def migrate_file(filepath):
    with open(filepath, 'r') as f:
        content = f.read()

    # Reset common errors from previous run if any
    # (e.g. e.Name(, err)) - wait, it's better to just re-read or assume we start fresh
    # Actually, some files are already partially changed.

    # 1. Replace import
    content = content.replace('"github.com/cockroachdb/errors"', '"errors"')

    # Better Wrap/Wrapf replacement using a loop to find calls
    for func in ['errors.Wrap', 'errors.Wrapf']:
        start = 0
        while True:
            idx = content.find(func + '(', start)
            if idx == -1:
                break
            
            end_idx = find_closing_paren(content, idx + len(func) + 1)
            if end_idx == -1:
                start = idx + 1
                continue
            
            args_str = content[idx + len(func) + 1 : end_idx]
            # Split args by comma, but be careful with nested commas
            # For simplicity, we can use a small parser or just look for the first two args
            
            # Simple comma splitter that respects parens
            args = []
            current = ""
            p_count = 0
            for char in args_str:
                if char == '(' : p_count += 1
                elif char == ')' : p_count -= 1
                
                if char == ',' and p_count == 0:
                    args.append(current.strip())
                    current = ""
                else:
                    current += char
            args.append(current.strip())
            
            if len(args) >= 2:
                err = args[0]
                fmt = args[1]
                rest = args[2:]
                
                new_call = ""
                if func == 'errors.Wrap':
                    new_call = f'fmt.Errorf({fmt} + ": %w", {err})'
                else: # Wrapf
                    if rest:
                        new_call = f'fmt.Errorf({fmt[:-1]}: %w", {", ".join(rest)}, {err})'
                    else:
                        new_call = f'fmt.Errorf({fmt[:-1]}: %w", {err})'
                
                content = content[:idx] + new_call + content[end_idx+1:]
                start = idx + len(new_call)
            else:
                start = idx + 1

    # 4. errors.Errorf -> fmt.Errorf, errors.Newf -> fmt.Errorf
    content = re.sub(r'errors\.Errorf\(', 'fmt.Errorf(', content)
    content = re.sub(r'errors\.Newf\(', 'fmt.Errorf(', content)

    # 5. Handle Cause (simple cases)
    content = re.sub(r'errors\.Cause\(\s*([^)]+)\s*\)\s*==\s*([^ \n\t;]+)', r'errors.Is(\1, \2)', content)
    
    # 6. Ensure "fmt" is imported if fmt.Errorf is used
    if 'fmt.Errorf' in content and '"fmt"' not in content:
        if 'import (' in content:
            content = content.replace('import (', 'import (\n\t"fmt"', 1)
        else:
            content = re.sub(r'(package \w+)', r'\1\n\nimport "fmt"', content)

    with open(filepath, 'w') as f:
        f.write(content)

# Since I already ran it once, some files might have "fmt.Errorf(\"add %q: %w\", e.Name(, err))"
# I should probably revert or fix those first, or just make sure the new script handles them.
# Re-running on the current state might be messy.
# I'll try to revert the changes to these files first if possible, but I don't have git revert easy.
# Actually, I'll just run it and hope for the best, or use `git checkout` if I can.

for f in files:
    if os.path.exists(f):
        # Revert to git state to start fresh
        os.system(f"git checkout {f}")
        print(f"Migrating {f}...")
        migrate_file(f)
