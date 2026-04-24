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
        if not '/api/' in self.base_url:
            self.base_url = f"{self.base_url}/api/v2"
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
        return self._request("get", "user/check")

    def create_repo(self, repo_name: str, is_private: bool = False) -> bool:
        payload = {
            "name": repo_name,
            "description": fake.sentence(nb_words=10),
            "private": is_private,
            "auto_init": True,
            "wiki_folder_name": "docs" # Aseguramos que use docs
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
            "title": title, "content": body, "head_branch": head_branch, "base_branch": base_branch
        }
        resp = self._request("post", f"{self.username}/{repo_name}/pulls", json=payload)
        return resp.json().get("index")

    def create_issue(self, repo_name: str, title: str, content: str, labels: Optional[List[str]] = None) -> Optional[int]:
        payload = {"title": title, "content": content}
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
            "draft": False, "prerelease": prerelease
        }
        self._request("post", f"{self.username}/{repo_name}/releases", json=payload)

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
    for _ in range(random.randint(1, 2)):
        branch_name = f"feature/{random_string(6)}"
        client.create_branch(repo_name, branch_name)
        branches.append(branch_name)
        logging.info(f"Created branch: {branch_name}")
    
    # 3. Create rich Wiki/Docs structure
    docs_files = [
        {
            "path": "docs/Home.md", 
            "content": f"# Welcome to {repo_name}\n\n{fake.paragraph(nb_sentences=5)}\n\nCheck out the [API Guide](api/v2/endpoints).",
            "message": "docs: initialize home page"
        },
        {
            "path": "docs/getting-started/installation.mdx",
            "content": f"""# Installation Guide

<Note>
  Make sure you have Go 1.25+ installed before proceeding.
</Note>

## Quick Start
```bash
curl -sSL https://ggit.io/install | sh
ggit init {repo_name}
```

## Manual Setup
{fake.paragraph()}""",
            "message": "docs: add installation guide"
        },
        {
            "path": "docs/api/v2/endpoints.mdx",
            "content": f"""# API Endpoints

<Warning>
  All requests must include a valid Bearer Token in the headers.
</Warning>

### Users
`GET /api/v2/users`

```json
{{
  "users": ["alice", "bob"],
  "count": 2
}}
```

### Repositories
{fake.paragraph()}""",
            "message": "docs: document api endpoints"
        },
        {
            "path": "docs/dev/database-schema.md",
            "content": f"## Database Schema\n\n{fake.paragraph()}\n\n![Schema](https://via.placeholder.com/600x400)",
            "message": "docs: add database schema"
        }
    ]
    client.create_commit(repo_name, docs_files)
    logging.info("Created Doc-as-Code wiki structure in docs/")
    
    # 4. Create source code commits
    for _ in range(random.randint(2, 4)):
        component = random.choice(COMPONENTS)
        message = random.choice(COMMIT_MESSAGES).format(
            component=component, 
            feature=random.choice(FEATURES),
            doc_section=random.choice(DOC_SECTIONS),
            module=random.choice(MODULES),
            file=fake.file_name()
        )
        
        files = []
        for _ in range(random.randint(1, 2)):
            file_path = f"src/{component}/{random_string(8)}.{random.choice(['go', 'py', 'js'])}"
            files.append({
                "path": file_path,
                "content": f"// {message}\npackage {component}\n\n{fake.paragraph()}",
                "message": message,
                "is_new": True
            })
        client.create_commit(repo_name, files)

    # 5. Create Issues
    issue_ids = []
    for _ in range(random.randint(2, 3)):
        title = random.choice(ISSUE_TITLES).format(
            component=random.choice(COMPONENTS),
            action=random.choice(ACTIONS),
            feature=random.choice(FEATURES),
            topic=random.choice(TOPICS),
            module=random.choice(MODULES)
        )
        issue_id = client.create_issue(repo_name, title, f"Check the docs at [[API Reference]]\n\n{fake.paragraph()}")
        if issue_id:
            issue_ids.append(issue_id)
            # Add comments
            client.add_issue_comment(repo_name, issue_id, f"I found a similar issue in repo #1. {fake.sentence()}")

    # 6. Create PRs
    for branch in branches:
        title = f"feat: implement {random.choice(FEATURES)}"
        pr_id = client.create_pull_request(repo_name, title, f"Closes #{random.choice(issue_ids) if issue_ids else 1}\n\n{fake.paragraph()}", branch)
        if pr_id:
            logging.info(f"Created PR #{pr_id}")

    # 7. Create a Release
    client.create_release(repo_name, "v1.0.0", "Release One", f"### Changes\n- {fake.sentence()}\n- Fixed #{issue_ids[0] if issue_ids else 1}")
    logging.info("Finished seeding repository")

def main():
    parser = argparse.ArgumentParser(description="Ggit Test Data Generator")
    parser.add_argument("--url", default="http://localhost:5555/api/v2", help="API URL")
    parser.add_argument("--username", required=True, help="Your username")
    parser.add_argument("--count", type=int, default=1, help="Number of repositories")
    args = parser.parse_args()

    password = getpass.getpass(f"Password for {args.username}: ")
    
    client = GgitClient(args.url, args.username, password)
    
    try:
        client.check_auth()
    except Exception:
        print("[!] Authentication failed.")
        return

    print(f"[+] Authenticated as {args.username}\n")
    
    for i in range(args.count):
        repo_name = random_repo_name()
        print(f"\nRepository {i+1}/{args.count}: {repo_name}")
        seed_repository(client, repo_name)
        if i < args.count - 1:
            time.sleep(1)

if __name__ == "__main__":
    main()
