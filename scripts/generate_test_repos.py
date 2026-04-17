import requests
import sys
import getpass
from requests.auth import HTTPBasicAuth

BASE_URL = "http://localhost:3000/api/v2"

def seed_repository(auth, username, repo_name):
    print(f"\n--- Seeding Repository: {repo_name} ---")
    
    # 1. Create Repo
    payload = {
        "name": repo_name,
        "description": f"Automated test repository for {repo_name}",
        "private": False,
        "auto_init": True
    }
    resp = requests.post(f"{BASE_URL}/user/repos", json=payload, auth=auth)
    if resp.status_code == 201:
        print("[+] Repository created")
    elif resp.status_code == 422:
        print("[!] Repository already exists, adding content to it...")
    else:
        print(f"[!] Error creating repo: {resp.text}")
        return

    # 2. Commits (Editor API)
    files = [
        {"path": "CONTRIBUTING.md", "content": "# How to contribute\nJust send a PR!"},
        {"path": f"src/{repo_name}.txt", "content": f"Sample content for {repo_name}"},
        {"path": "TODO.md", "content": "- [ ] Task 1\n- [x] Task 2 (Done)"}
    ]
    for f in files:
        payload = {
            "tree_path": f["path"],
            "content": f["content"],
            "commit_summary": f"Initial commit for {f['path']}",
            "is_new_file": True
        }
        requests.post(f"{BASE_URL}/{username}/{repo_name}/editor/edit", json=payload, auth=auth)
    print(f"[+] Committed {len(files)} files")

    # 3. Issues and Comments
    for i in range(1, 3):
        issue_payload = {
            "title": f"Bug found in {repo_name} #{i}",
            "content": "Description of a very serious bug that needs fixing."
        }
        iresp = requests.post(f"{BASE_URL}/{username}/{repo_name}/issues", json=issue_payload, auth=auth)
        if iresp.status_code == 201:
            issue_data = iresp.json()
            idx = issue_data.get('index')
            if idx:
                requests.post(f"{BASE_URL}/{username}/{repo_name}/issues/{idx}/comments", 
                            json={"body": "I'm looking into this right now."}, auth=auth)
            else:
                print(f" [!] Issue created but 'index' not found in response: {issue_data}")
    print("[+] Created 2 issues with comments")

    # 4. Wiki
    wiki_pages = [
        {"title": "Home", "content": f"# Welcome to {repo_name} Wiki"},
        {"title": "Setup", "content": "## Requirements\n- Python 3\n- Requests library"}
    ]
    for page in wiki_pages:
        requests.post(f"{BASE_URL}/{username}/{repo_name}/wiki", json=page, auth=auth)
    print("[+] Created Wiki pages")

    # 5. Release
    release_payload = {
        "tag_name": "v1.0.0",
        "name": "First Release",
        "body": "Automatically generated release for testing.",
        "draft": False,
        "prerelease": False
    }
    requests.post(f"{BASE_URL}/{username}/{repo_name}/releases", json=release_payload, auth=auth)
    print("[+] Created Release v1.0.0")

def main():
    print("=== Ggit Test Data Generator ===")
    username = input("Username: ")
    password = getpass.getpass("Password: ")
    
    try:
        count = int(input("How many repositories to create? (e.g. 3): "))
    except ValueError:
        print("Invalid number")
        return

    auth = HTTPBasicAuth(username, password)

    # Verify credentials first
    test_resp = requests.get(f"{BASE_URL}/user/check", auth=auth)
    if test_resp.status_code != 200:
        print(f"[!] Authentication failed: {test_resp.text}")
        return
    print(f"[+] Authenticated as {username}")

    for i in range(1, count + 1):
        repo_name = f"test-project-{i}"
        seed_repository(auth, username, repo_name)

    print("\n[✔] Done! All repositories have been seeded.")

if __name__ == "__main__":
    main()
