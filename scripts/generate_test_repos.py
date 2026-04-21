import requests
import sys
import getpass
import random
import string
import time
from datetime import datetime, timedelta
from requests.auth import HTTPBasicAuth
from faker import Faker  # Necesitas: pip install faker

BASE_URL = "http://localhost:3005/api/v2"
fake = Faker()

# Listas para generar contenido realista
COMMIT_MESSAGES = [
    "fix: resolve issue with {component}",
    "feat: add new {feature} functionality",
    "docs: update {doc_section} documentation",
    "refactor: clean up {module} code structure",
    "test: add unit tests for {component}",
    "chore: update dependencies",
    "perf: optimize {component} performance",
    "style: format {file} according to standards"
]

COMPONENTS = ["auth", "database", "api", "frontend", "backend", "config"]
FEATURES = ["user authentication", "data export", "search", "notifications", "logging"]
DOC_SECTIONS = ["README", "API docs", "installation guide", "contributing"]
MODULES = ["core", "utils", "models", "controllers", "services"]

ISSUE_TITLES = [
    "Bug: {component} crashes when {action}",
    "Feature request: add {feature} support",
    "Documentation missing for {topic}",
    "Performance issue in {component}",
    "Security vulnerability in {module}"
]

ACTIONS = ["loading large files", "saving data", "deleting records", "updating settings"]
TOPICS = ["deployment", "configuration", "API endpoints", "error handling"]

PR_TITLES = [
    "Implement {feature}",
    "Fix {bug_type} in {component}",
    "Update documentation for {topic}",
    "Refactor {module} for better performance",
    "Add tests for {component}"
]

BUG_TYPES = ["null pointer", "off-by-one error", "race condition", "memory leak"]

def random_string(length=8):
    return ''.join(random.choices(string.ascii_lowercase + string.digits, k=length))

def random_repo_name():
    prefixes = ["awesome", "test", "demo", "sample", "example", "my", "cool", "project"]
    suffixes = ["app", "service", "lib", "tool", "suite", "platform"]
    return f"{random.choice(prefixes)}-{random.choice(suffixes)}-{random_string(4)}"

def random_date(start_date=None, end_date=None):
    if not start_date:
        start_date = datetime.now() - timedelta(days=30)
    if not end_date:
        end_date = datetime.now()
    return fake.date_time_between(start_date=start_date, end_date=end_date)

def create_repo(auth, username, repo_name, is_private=False):
    """Create a new repository"""
    payload = {
        "name": repo_name,
        "description": fake.sentence(nb_words=10),
        "private": is_private,
        "auto_init": True,
        "created_at": random_date().isoformat()
    }
    resp = requests.post(f"{BASE_URL}/user/repos", json=payload, auth=auth)
    return resp.status_code == 201 or resp.status_code == 422

def create_branch(auth, username, repo_name, branch_name, base_branch="main"):
    """Create a new branch"""
    payload = {
        "old_ref_name": base_branch,
        "new_ref_name": branch_name
    }
    resp = requests.post(f"{BASE_URL}/{username}/{repo_name}/branches", json=payload, auth=auth)
    return resp.status_code == 201

def create_commit(auth, username, repo_name, files, branch="main"):
    """Create commits with multiple files"""
    for file_info in files:
        payload = {
            "tree_path": file_info["path"],
            "content": file_info["content"],
            "commit_summary": file_info["message"],
            "commit_description": fake.paragraph(),
            "branch": branch,
            "is_new_file": file_info.get("is_new", True)
        }
        requests.post(f"{BASE_URL}/{username}/{repo_name}/editor/edit", json=payload, auth=auth)
        time.sleep(0.1)  # Pequeña pausa para evitar rate limiting

def create_pull_request(auth, username, repo_name, title, body, head_branch, base_branch="main"):
    """Create a pull request"""
    payload = {
        "title": title,
        "body": body,
        "head": head_branch,
        "base": base_branch,
        "created_at": random_date().isoformat()
    }
    resp = requests.post(f"{BASE_URL}/{username}/{repo_name}/pulls", json=payload, auth=auth)
    if resp.status_code == 201:
        return resp.json().get("index")
    return None

def create_issue(auth, username, repo_name, title, content, labels=None):
    """Create an issue with optional labels"""
    payload = {
        "title": title,
        "content": content,
        "created_at": random_date().isoformat()
    }
    if labels:
        payload["labels"] = labels
    
    resp = requests.post(f"{BASE_URL}/{username}/{repo_name}/issues", json=payload, auth=auth)
    if resp.status_code == 201:
        return resp.json().get("index")
    return None

def add_issue_comment(auth, username, repo_name, issue_index, comment):
    """Add comment to an issue"""
    payload = {"body": comment}
    requests.post(f"{BASE_URL}/{username}/{repo_name}/issues/{issue_index}/comments", 
                 json=payload, auth=auth)

def create_wiki_page(auth, username, repo_name, title, content):
    """Create a wiki page"""
    payload = {"title": title, "content": content}
    requests.post(f"{BASE_URL}/{username}/{repo_name}/wiki", json=payload, auth=auth)

def create_release(auth, username, repo_name, tag, name, body, prerelease=False):
    """Create a release"""
    payload = {
        "tag_name": tag,
        "name": name,
        "body": body,
        "draft": "false",
        "prerelease": prerelease
    }
    requests.post(f"{BASE_URL}/{username}/{repo_name}/releases", json=payload, auth=auth)

def add_collaborator(auth, username, repo_name, collaborator_name, permission="write"):
    """Add a collaborator to the repository"""
    payload = {"permission": permission}
    requests.post(f"{BASE_URL}/{username}/{repo_name}/collaborators/{collaborator_name}", 
                 json=payload, auth=auth)

def seed_repository(auth, username, repo_name):
    print(f"\n--- Seeding Repository: {repo_name} ---")
    
    # 1. Create repository
    is_private = random.choice([True, False])
    if create_repo(auth, username, repo_name, is_private):
        print(f"[+] Repository created (private: {is_private})")
    else:
        print(f"[!] Error creating repo")
        return
    
    # 2. Create multiple branches
    branches = []
    for _ in range(random.randint(1, 3)):
        branch_name = f"feature/{random_string(6)}"
        if create_branch(auth, username, repo_name, branch_name):
            branches.append(branch_name)
            print(f"[+] Created branch: {branch_name}")
    
    # 3. Create commits on main branch
    for _ in range(random.randint(2, 5)):
        component = random.choice(COMPONENTS)
        message = random.choice(COMMIT_MESSAGES).format(
            component=component, 
            feature=random.choice(FEATURES),
            doc_section=random.choice(DOC_SECTIONS),
            module=random.choice(MODULES),
            file=fake.file_name()
        )
        
        files = []
        for _ in range(random.randint(1, 3)):
            file_path = f"src/{random_string(8)}.{random.choice(['py', 'js', 'go', 'rs', 'java'])}"
            files.append({
                "path": file_path,
                "content": fake.paragraph(nb_sentences=random.randint(5, 20)),
                "message": message,
                "is_new": True
            })
        
        create_commit(auth, username, repo_name, files)
        print(f"[+] Commit: {message[:50]}...")
    
    # 4. Create commits on feature branches
    for branch in branches:
        for _ in range(random.randint(1, 3)):
            message = f"feat: implement {random.choice(FEATURES)}"
            files = [{
                "path": f"features/{random_string(10)}.md",
                "content": fake.paragraph(nb_sentences=10),
                "message": message,
                "is_new": True
            }]
            create_commit(auth, username, repo_name, files, branch=branch)
            print(f"[+] Commit on {branch}: {message}")
    
    # 5. Create Pull Requests
    for branch in branches:
        if random.choice([True, False]):  # 50% chance to create PR for each branch
            title = random.choice(PR_TITLES).format(
                feature=random.choice(FEATURES),
                bug_type=random.choice(BUG_TYPES),
                component=random.choice(COMPONENTS),
                topic=random.choice(TOPICS),
                module=random.choice(MODULES)
            )
            body = f"""## Description
{fake.paragraph()}

## Changes made
- {fake.sentence()}
- {fake.sentence()}
- {fake.sentence()}

## Testing
{fake.paragraph()}"""
            
            pr_index = create_pull_request(auth, username, repo_name, title, body, branch)
            if pr_index:
                print(f"[+] Created PR #{pr_index}: {title}")
                
                # Add comments to PR
                if random.choice([True, False]):
                    time.sleep(0.5)
                    add_issue_comment(auth, username, repo_name, pr_index, 
                                    f"Thanks for the contribution! I'll review this soon. {fake.sentence()}")
                    print(f"   Added comment to PR #{pr_index}")
    
    # 6. Create Issues
    issues = []
    for _ in range(random.randint(2, 4)):
        title = random.choice(ISSUE_TITLES).format(
            component=random.choice(COMPONENTS),
            action=random.choice(ACTIONS),
            feature=random.choice(FEATURES),
            topic=random.choice(TOPICS),
            module=random.choice(MODULES)
        )
        content = f"""## Description
{fake.paragraph()}

## Steps to reproduce
1. {fake.sentence()}
2. {fake.sentence()}
3. {fake.sentence()}

## Expected behavior
{fake.sentence()}

## Actual behavior
{fake.sentence()}

## Environment
- OS: {random.choice(['Ubuntu 22.04', 'macOS 14', 'Windows 11'])}
- Version: {random.choice(['1.0.0', '1.1.0', '2.0.0'])}"""
        
        issue_idx = create_issue(auth, username, repo_name, title, content)
        if issue_idx:
            issues.append(issue_idx)
            print(f"[+] Created issue #{issue_idx}: {title[:50]}...")
            
            # Add comments to issues
            for _ in range(random.randint(1, 3)):
                time.sleep(0.3)
                comment = random.choice([
                    f"I can reproduce this. {fake.sentence()}",
                    f"Looking into it. {fake.sentence()}",
                    f"Thanks for reporting! {fake.sentence()}",
                    f"Fixed in latest commit. {fake.sentence()}"
                ])
                add_issue_comment(auth, username, repo_name, issue_idx, comment)
                print(f"   Added comment to issue #{issue_idx}")
    
    # 7. Create Wiki pages
    wiki_pages = [
        {"title": "Home", "content": f"# Welcome to {repo_name}\n\n{fake.paragraph(nb_sentences=5)}"},
        {"title": "Getting Started", "content": f"## Installation\n```bash\npip install {repo_name}\n```\n\n## Usage\n{fake.paragraph()}"},
        {"title": "API Reference", "content": f"## Endpoints\n\n### GET /api/v1/users\n{fake.paragraph()}\n\n### POST /api/v1/users\n{fake.paragraph()}"},
        {"title": "Contributing", "content": f"We welcome contributions! Please read our guidelines.\n\n{fake.paragraph(nb_sentences=3)}"},
        {"title": "FAQ", "content": f"### Question 1\n{fake.sentence()}\n\n### Answer\n{fake.paragraph()}\n\n### Question 2\n{fake.sentence()}\n\n### Answer\n{fake.paragraph()}"}
    ]
    
    selected_pages = random.sample(wiki_pages, random.randint(2, len(wiki_pages)))
    for page in selected_pages:
        create_wiki_page(auth, username, repo_name, page["title"], page["content"])
        print(f"[+] Created wiki page: {page['title']}")
    
    # 8. Create Releases
    releases = [
        ("v1.0.0", "Initial Release", False),
        ("v1.1.0", "Feature Update", False),
        ("v2.0.0-beta", "Beta Release", True)
    ]
    
    selected_releases = random.sample(releases, random.randint(1, 2))
    for tag, name, prerelease in selected_releases:
        body = f"""## What's Changed
- {fake.sentence()}
- {fake.sentence()}
- {fake.sentence()}

## Installation
```bash
git clone http://localhost:3005/{username}/{repo_name}.git
cd {repo_name}
git checkout {tag}
```

## Contributors
- @{username}
{fake.name() if random.choice([True, False]) else ''}

**Full Changelog**: `http://localhost:3005/{username}/{repo_name}/compare/{tag}...v1.0.0`"""
        
        create_release(auth, username, repo_name, tag, name, body, prerelease)
        print(f"[+] Created release: {name} ({tag}){' [prerelease]' if prerelease else ''}")
    
    print(f"[✓] Finished seeding {repo_name}")

def main():
    print("=== Ggit Test Data Generator ===")
    print("Generates realistic test data for version control server testing\n")
    
    username = input("Username: ")
    password = getpass.getpass("Password: ")
    
    try:
        count = int(input("How many repositories to create? (e.g. 5): "))
    except ValueError:
        print("Invalid number")
        return
    
    auth = HTTPBasicAuth(username, password)
    
    # Verify credentials
    test_resp = requests.get(f"{BASE_URL}/user/check", auth=auth)
    if test_resp.status_code != 200:
        print(f"[!] Authentication failed: {test_resp.text}")
        return
    print(f"[+] Authenticated as {username}\n")
    
    # Generate repositories
    for i in range(1, count + 1):
        repo_name = random_repo_name()
        print(f"\n{'='*60}")
        print(f"Repository {i}/{count}: {repo_name}")
        print(f"{'='*60}")
        
        seed_repository(auth, username, repo_name)
        
        # Wait between repositories to avoid overwhelming the server
        if i < count:
            wait_time = random.uniform(1, 3)
            print(f"\nWaiting {wait_time:.1f} seconds before next repository...")
            time.sleep(wait_time)
    
    print("\n" + "="*60)
    print("[✔] Done! All repositories have been seeded successfully.")
    print(f"Generated {count} repositories with realistic data (branches, commits, PRs, issues, wiki, releases)")

if __name__ == "__main__":
    main()
