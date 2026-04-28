package janitor

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gogs/git-module"

	"github.com/nef734/ggit/internal/database"
)

// Issue represents a finding by the Janitor.
type Issue struct {
	Type        string `json:"type"` // "system", "dependency", "build", "large_file"
	Severity    string `json:"severity"` // "info", "warning", "critical"
	Path        string `json:"path"`
	Description string `json:"description"`
	Size        int64  `json:"size"`
}

// Report contains all findings for a repository.
type Report struct {
	Issues      []Issue `json:"issues"`
	TotalWaste  int64   `json:"total_waste"`
	GeneratedAt int64   `json:"generated_at"`
}

// Patterns for detection
var patterns = map[string]string{
	// System
	".DS_Store": "macOS metadata file",
	"Thumbs.db": "Windows thumbnail cache",
	"desktop.ini": "Windows folder settings",
	".Spotlight-V100": "macOS Spotlight index",
	".Trashes": "macOS Trash folder",
	
	// Dependencies
	"node_modules": "Node.js dependencies (should be in .gitignore)",
	"bower_components": "Bower components",
	"vendor": "Go/PHP dependencies (usually better to not commit)",
	"venv": "Python virtual environment",
	".venv": "Python virtual environment",
	"env": "Python virtual environment",
	"__pycache__": "Python bytecode cache",
	".bundle": "Ruby bundle folder",
	"jspm_packages": "JSPM packages",
	
	// Build / Artifacts
	"dist": "Build distribution folder",
	"build": "Build output folder",
	"bin": "Binary output folder",
	"obj": "Object files folder",
	"target": "Maven/Rust build target",
	"out": "IDE output folder",
	".next": "Next.js build folder",
	".nuxt": "Nuxt.js build folder",
	".sass-cache": "Sass cache folder",
	
	// IDEs
	".idea": "IntelliJ IDEA metadata",
	".vscode": "Visual Studio Code settings",
	"*.suo": "Visual Studio user options",
	"*.user": "Visual Studio user project settings",
	".project": "Eclipse project file",
	".settings": "Eclipse settings folder",
	
	// Logs & Coverage
	"*.log": "Log file",
	"npm-debug.log*": "npm debug log",
	"yarn-error.log*": "yarn error log",
	"coverage": "Code coverage report",
	".nyc_output": "Istanbul/nyc coverage data",
	"htmlcov": "Python coverage report",
}

// Scan analyzes the repository tree.
func Scan(ctx context.Context, repoID int64) (*Report, error) {
	repo, err := database.GetRepositoryByID(repoID)
	if err != nil {
		return nil, err
	}

	gitRepo, err := git.Open(repo.RepoPath())
	if err != nil {
		return nil, err
	}

	branch := repo.DefaultBranch
	if branch == "" {
		branch = "main"
	}
	
	if !gitRepo.HasBranch(branch) {
		if gitRepo.HasBranch("master") {
			branch = "master"
		} else {
			return nil, fmt.Errorf("default branch not found")
		}
	}

	tree, err := gitRepo.LsTree(branch)
	if err != nil {
		return nil, err
	}

	report := &Report{
		GeneratedAt: time.Now().Unix(),
		Issues:      []Issue{},
	}

	entries, err := tree.Entries()
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		name := path.Base(entry.Name())
		
		matchFound := false
		foundDesc := ""
		for pattern, desc := range patterns {
			matched, _ := path.Match(pattern, name)
			if matched || name == pattern {
				matchFound = true
				foundDesc = desc
				break
			}
		}

		if matchFound {
			severity := "warning"
			if entry.Type() == git.ObjectCommit {
				continue
			}
			
			issueType := "build"
			if strings.Contains(foundDesc, "dependencies") || strings.Contains(foundDesc, "environment") {
				issueType = "dependency"
			} else if strings.Contains(foundDesc, "metadata") {
				issueType = "system"
			} else if strings.Contains(foundDesc, "IDE") || strings.Contains(foundDesc, "settings") {
				issueType = "build" // or "ide"
			}

			report.Issues = append(report.Issues, Issue{
				Type:        issueType,
				Severity:    severity,
				Path:        entry.Name(),
				Description: foundDesc,
				Size:        0,
			})
		}
		
		if entry.Type() == git.ObjectBlob && entry.Size() > 50*1024*1024 {
			report.Issues = append(report.Issues, Issue{
				Type:        "large_file",
				Severity:    "warning",
				Path:        entry.Name(),
				Description: "Large file detected (> 50MB)",
				Size:        entry.Size(),
			})
			report.TotalWaste += entry.Size()
		}
	}

	data, _ := json.Marshal(report)
	err = database.Handle.DB().Model(repo).Updates(map[string]any{
		"janitor_report":          string(data),
		"janitor_last_check_unix": report.GeneratedAt,
	}).Error

	return report, err
}

// Cleanup removes the specified paths and updates .gitignore if needed.
func Cleanup(ctx context.Context, repoID int64, doer *database.User, paths []string) error {
	repo, err := database.GetRepositoryByID(repoID)
	if err != nil {
		return err
	}

	// Use internal editor logic to handle the commit
	// Since we need to delete multiple files, we'll implement a custom logic 
	// inspired by repo_editor.go
	
	// 1. Prepare local copy
	if err := repo.DiscardLocalRepoBranchChanges(repo.DefaultBranch); err != nil {
		return err
	} else if err = repo.UpdateLocalCopyBranch(repo.DefaultBranch); err != nil {
		return err
	}

	localPath := repo.LocalCopyPath()
	
	// 2. Remove files/folders
	removedAny := false
	for _, p := range paths {
		fullPath := filepath.Join(localPath, p)
		if info, err := os.Stat(fullPath); err == nil {
			if info.IsDir() {
				os.RemoveAll(fullPath)
			} else {
				os.Remove(fullPath)
			}
			removedAny = true
		}
	}

	if !removedAny {
		return fmt.Errorf("no files found to remove")
	}

	// 3. Update .gitignore (Suggestion logic)
	// We'll append the unique patterns to .gitignore
	gitignorePath := filepath.Join(localPath, ".gitignore")
	existingGitignore := ""
	if data, err := os.ReadFile(gitignorePath); err == nil {
		existingGitignore = string(data)
	}

	newPatterns := []string{}
	for _, p := range paths {
		pattern := p
		if info, err := os.Stat(filepath.Join(localPath, p)); err == nil && info.IsDir() {
			pattern = p + "/"
		}
		if !strings.Contains(existingGitignore, pattern) {
			newPatterns = append(newPatterns, pattern)
		}
	}

	if len(newPatterns) > 0 {
		f, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			if !strings.HasSuffix(existingGitignore, "\n") && existingGitignore != "" {
				f.WriteString("\n")
			}
			f.WriteString("\n# Added by Ggit Janitor\n")
			for _, np := range newPatterns {
				f.WriteString(np + "\n")
			}
			f.Close()
		}
	}

	// 4. Commit and Push
	if err := git.Add(localPath, git.AddOptions{All: true}); err != nil {
		return err
	}

	message := "chore: cleanup repository by Ggit Janitor\n\nRemoved unnecessary files and updated .gitignore."
	err = git.CreateCommit(
		localPath,
		&git.Signature{
			Name:  doer.DisplayName(),
			Email: doer.Email,
			When:  time.Now(),
		},
		message,
	)
	if err != nil {
		return err
	}

	return git.Push(localPath, "origin", repo.DefaultBranch,
		git.PushOptions{
			CommandOptions: git.CommandOptions{
				Envs: database.ComposeHookEnvs(database.ComposeHookEnvsOptions{
					AuthUser:  doer,
					OwnerName: repo.MustOwner().Name,
					OwnerSalt: repo.MustOwner().Salt,
					RepoID:    repo.ID,
					RepoName:  repo.Name,
					RepoPath:  repo.RepoPath(),
				}),
			},
		},
	)
}
