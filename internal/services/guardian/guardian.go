package guardian

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/gogs/git-module"
	log "unknwon.dev/clog/v2"

	"github.com/nef734/ggit/internal/database"
)

// Violation represents a found secret in a push.
type Violation struct {
	FileName string
	LineNum  int
	Pattern  string
	Content  string // The line content with the secret
}

// ScanPush scans the changes between two commits for secrets.
func ScanPush(ctx context.Context, repoID int64, oldCommitID, newCommitID string) ([]Violation, error) {
	repo, err := database.GetRepositoryByID(repoID)
	if err != nil {
		return nil, fmt.Errorf("get repository: %v", err)
	}

	if !repo.EnableSecretGuardian {
		return nil, nil
	}

	patterns, err := database.Handle.SecretPatterns().ListByRepo(ctx, repoID)
	if err != nil {
		return nil, fmt.Errorf("list patterns: %v", err)
	}

	if len(patterns) == 0 {
		return nil, nil
	}

	// Compile regexes
	type compiled struct {
		name string
		re   *regexp.Regexp
	}
	var res []compiled
	for _, p := range patterns {
		re, err := regexp.Compile(p.Regex)
		if err != nil {
			log.Warn("Guardian: Invalid regex pattern %q: %v", p.Name, err)
			continue
		}
		res = append(res, compiled{name: p.Name, re: re})
	}

	gitRepo, err := git.Open(repo.RepoPath())
	if err != nil {
		return nil, fmt.Errorf("open git repo: %v", err)
	}

	// If oldCommitID is empty, it's a new branch push.
	// We scan the entire tree of newCommitID or just the new commits?
	// Standard practice is to scan the diff between old and new.
	if oldCommitID == git.EmptyID {
		// For new branches, we could scan the whole commit content, 
		// but let's stick to the new commit's changes compared to its parent for performance.
		// Or if it's the very first push, compare against empty tree.
		oldCommitID = "4b825dc642cb6eb9a060e54bf8d69288fbee4904" // Empty tree SHA
	}

	diff, err := gitRepo.Diff(newCommitID, 500, 1000, 1000, git.DiffOptions{
		Base: oldCommitID,
	})
	if err != nil {
		return nil, fmt.Errorf("get diff: %v", err)
	}

	var violations []Violation
	for _, file := range diff.Files {
		// 1. Check filename
		for _, cp := range res {
			if cp.re.MatchString(file.Name) {
				violations = append(violations, Violation{
					FileName: file.Name,
					Pattern:  cp.name,
					Content:  "[FILENAME MATCH]",
				})
			}
		}

		// 2. Check content (added lines only)
		for _, section := range file.Sections {
			for _, line := range section.Lines {
				if line.Type != git.DiffLineAdd {
					continue
				}

				content := strings.TrimPrefix(line.Content, "+")
				for _, cp := range res {
					if cp.re.MatchString(content) {
						violations = append(violations, Violation{
							FileName: file.Name,
							LineNum:  line.RightLine,
							Pattern:  cp.name,
							Content:  content,
						})
					}
				}
			}
		}
	}

	return violations, nil
}
