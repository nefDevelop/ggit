import argparse
import getpass
import logging
import os
import requests
import random
import string
import time
from datetime import datetime, timedelta
from requests.auth import HTTPBasicAuth
from faker import Faker  # Necesitas: pip install faker
from typing import List, Dict, Any, Optional

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

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

class GgitClient:
    """Un cliente para interactuar coa API v2 de Ggit."""


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

class GgitClient:
    """Un cliente para interactuar coa API v2 de Ggit."""

    def __init__(self, base_url: str, username: str, password: str):
        self.base_url = base_url.rstrip('/')
        self.username = username
        self.session = requests.Session()
        self.session.auth = HTTPBasicAuth(username, password)

    def _request(self, method: str, endpoint: str, **kwargs) -> requests.Response:
        url = f"{self.base_url}/{endpoint}"
        try:
            resp = self.session.request(method, url, **kwargs)
            resp.raise_for_status()
            return resp
        except requests.exceptions.HTTPError as e:
            logging.error(f"HTTP Error for {method.upper()} {url}: {e.response.status_code} - {e.response.text}")
            raise
        except requests.exceptions.RequestException as e:
            logging.error(f"Request failed for {method.upper()} {url}: {e}")
            raise

    def check_auth(self):
        self._request("get", "user/check")

    def create_repo(self, repo_name: str, is_private: bool = False) -> bool:
        payload = {
            "name": repo_name,
            "description": fake.sentence(nb_words=10),
            "private": is_private,
            "auto_init": True,
            "created_at": random_date().isoformat()
        }
        try:
            self._request("post", "user/repos", json=payload)
            return True
        except requests.exceptions.HTTPError as e:
            if e.response.status_code == 422: # Repo already exists
                logging.warning(f"Repository '{repo_name}' already exists. Continuing...")
                return True
            return False

    def create_branch(self, repo_name: str, branch_name: str, base_branch: str = "main"):
        payload = {"old_ref_name": base_branch, "new_ref_name": branch_name}
        self._request("post", f"{self.username}/{repo_name}/branches", json=payload)

    def create_commit(self, repo_name: str, files: List[Dict[str, Any]], branch: str = "main"):
        for file_info in files:
            payload = {
                "tree_path": file_info["path"],
                "content": file_info["content"],
                "commit_summary": file_info["message"],
                "commit_description": fake.paragraph(),
                "branch": branch,
                "is_new_file": file_info.get("is_new", True)
            }
            self._request("post", f"{self.username}/{repo_name}/editor/edit", json=payload)
            time.sleep(0.1)

    def create_pull_request(self, repo_name: str, title: str, body: str, head_branch: str, base_branch: str = "main") -> Optional[int]:
        payload = {
            "title": title, "body": body, "head": head_branch, "base": base_branch,
            "created_at": random_date().isoformat()
        }
        resp = self._request("post", f"{self.username}/{repo_name}/pulls", json=payload)
        return resp.json().get("index")

    def create_issue(self, repo_name: str, title: str, content: str, labels: Optional[List[str]] = None) -> Optional[int]:
        payload = {"title": title, "content": content, "created_at": random_date().isoformat()}
        if labels:
            payload["labels"] = labels
        resp = self._request("post", f"{self.username}/{repo_name}/issues", json=payload)
        return resp.json().get("index")

    def add_issue_comment(self, repo_name: str, issue_index: int, comment: str):
        payload = {"body": comment}
        self._request("post", f"{self.username}/{repo_name}/issues/{issue_index}/comments", json=payload)

    def create_wiki_page(self, repo_name: str, title: str, content: str):
        payload = {"title": title, "content": content}
        self._request("post", f"{self.username}/{repo_name}/wiki", json=payload)

    def create_release(self, repo_name: str, tag: str, name: str, body: str, prerelease: bool = False):
        payload = {
            "tag_name": tag, "name": name, "body": body,
            "draft": "false", "prerelease": prerelease
        }
        self._request("post", f"{self.username}/{repo_name}/releases", json=payload)

    def add_collaborator(self, repo_name: str, collaborator_name: str, permission: str = "write"):
        payload = {"permission": permission}
        self._request("post", f"{self.username}/{repo_name}/collaborators/{collaborator_name}", json=payload)


def seed_repository(client: GgitClient, repo_name: str):
    logging.info(f"--- Seeding Repository: {repo_name} ---")
    
    # 1. Create repository
    is_private = random.choice([True, False])
    if client.create_repo(repo_name, is_private):
        logging.info(f"Repository created (private: {is_private})")
    else:
        logging.error(f"Failed to create repository '{repo_name}'. Skipping.")
        return
    
    # 2. Create multiple branches
    branches = []
    for _ in range(random.randint(1, 3)):
        branch_name = f"feature/{random_string(6)}"
        client.create_branch(repo_name, branch_name)
        branches.append(branch_name)
        logging.info(f"Created branch: {branch_name}")
    
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
        
        client.create_commit(repo_name, files)
        logging.info(f"Commit on main: {message[:50]}...")
    
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
            client.create_commit(repo_name, files, branch=branch)
            logging.info(f"Commit on {branch}: {message}")
    
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
            
            pr_index = client.create_pull_request(repo_name, title, body, branch)
            if pr_index:
                logging.info(f"Created PR #{pr_index}: {title}")
                
                # Add comments to PR
                if random.choice([True, False]):
                    time.sleep(0.5)
                    client.add_issue_comment(repo_name, pr_index, f"Thanks for the contribution! I'll review this soon. {fake.sentence()}")
                    logging.info(f"  -> Added comment to PR #{pr_index}")
    
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
        
        issue_idx = client.create_issue(repo_name, title, content)
        if issue_idx:
            issues.append(issue_idx)
            logging.info(f"Created issue #{issue_idx}: {title[:50]}...")
            
            # Add comments to issues
            for _ in range(random.randint(1, 3)):
                time.sleep(0.3)
                comment = random.choice([
                    f"I can reproduce this. {fake.sentence()}",
                    f"Looking into it. {fake.sentence()}",
                    f"Thanks for reporting! {fake.sentence()}",
                    f"Fixed in latest commit. {fake.sentence()}"
                ])
                client.add_issue_comment(repo_name, issue_idx, comment)
                logging.info(f"  -> Added comment to issue #{issue_idx}")
    
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
        client.create_wiki_page(repo_name, page["title"], page["content"])
        logging.info(f"Created wiki page: {page['title']}")
    
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
```bash\ngit clone {client.base_url.replace('/api/v2', '')}/{client.username}/{repo_name}.git
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
