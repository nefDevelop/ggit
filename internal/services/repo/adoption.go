package repo

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogs/git-module"
	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/database"
	log "github.com/nef734/ggit/internal/log"
	"github.com/nef734/ggit/internal/repox"
	"gorm.io/gorm"
)

// UnmanagedRepository represents a repository folder on disk that is not in the database.
type UnmanagedRepository struct {
	OwnerName string `json:"owner_name"`
	RepoName  string `json:"repo_name"`
	Path      string `json:"path"`
}

// ListUnmanagedRepositories scans the repository root and returns repositories not in the database.
func (s *RepositoryService) ListUnmanagedRepositories(ctx context.Context) ([]*UnmanagedRepository, error) {
	root := conf.Repository.Root
	var unmanaged []*UnmanagedRepository

	owners, err := os.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("read repo root"+": %w", err)
	}

	for _, ownerDir := range owners {
		if !ownerDir.IsDir() {
			continue
		}
		ownerName := ownerDir.Name()

		owner, _ := database.Handle.Users().GetByUsername(ctx, ownerName)

		repos, err := os.ReadDir(filepath.Join(root, ownerName))
		if err != nil {
			continue
		}

		for _, repoDir := range repos {
			if !repoDir.IsDir() || !strings.HasSuffix(repoDir.Name(), ".git") {
				continue
			}
			repoName := strings.TrimSuffix(repoDir.Name(), ".git")

			if owner != nil {
				_, err = database.Handle.Repositories().GetByName(ctx, owner.ID, repoName)
				if err == nil {
					continue // Already in DB
				}
			}

			unmanaged = append(unmanaged, &UnmanagedRepository{
				OwnerName: ownerName,
				RepoName:  repoName,
				Path:      filepath.Join(root, ownerName, repoDir.Name()),
			})
		}
	}

	return unmanaged, nil
}

// AdoptRepository links an existing folder on disk to the database.
func (s *RepositoryService) AdoptRepository(ctx context.Context, ownerName, repoName string) (*database.Repository, error) {
	// 🚨 SECURITY: Prevent Path Traversal & System Hijacking
	// Strict validation: Only alphanumeric, dash, underscore, and dot allowed.
	if !database.IsValidRepoName(repoName) {
		return nil, fmt.Errorf("invalid or malicious repository name: %q", repoName)
	}

	owner, err := database.Handle.Users().GetByUsername(ctx, ownerName)
	if err != nil {
		return nil, fmt.Errorf("get owner"+": %w", err)
	}

	var repo *database.Repository
	err = database.Handle.DB().Transaction(func(tx *gorm.DB) error {
		domainRepo, err := database.NewRepositoriesStore(tx).Create(ctx, owner.ID, database.CreateRepoOptions{
			Name:          repoName,
			DefaultBranch: conf.Repository.DefaultBranch,
			EnableWiki:    true,
			EnableIssues:  true,
			EnablePulls:   true,
		})
		if err != nil {
			return fmt.Errorf("create db record"+": %w", err)
		}
		repo = database.FromDomainRepository(domainRepo)

		repoPath := repo.RepoPath()
		if !strings.HasSuffix(repoPath, ".git") {
			repoPath += ".git"
		}

		if _, err := os.Stat(repoPath); os.IsNotExist(err) {
			return fmt.Errorf("repository folder does not exist: %s", repoPath)
		}

		// Fix hooks
		if err := repox.CreateDelegateHooks(repoPath); err != nil {
			return fmt.Errorf("create delegate hooks"+": %w", err)
		}

		// Detect git info
		gitRepo, err := git.Open(repoPath)
		if err != nil {
			return fmt.Errorf("open git repo: %w", err)
		}

		branches, _ := gitRepo.Branches()
		if len(branches) > 0 {
			repo.IsBare = false
			refspec, err := gitRepo.SymbolicRef()
			if err == nil {
				repo.DefaultBranch = git.RefShortName(refspec)
			} else {
				repo.DefaultBranch = branches[0]
			}
		} else {
			repo.IsBare = true
		}

		// Save these initial fields
		if err := tx.Model(repo).Updates(map[string]any{
			"is_bare":        repo.IsBare,
			"default_branch": repo.DefaultBranch,
		}).Error; err != nil {
			return fmt.Errorf("update repo info: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Update stats (commits, branches, etc.)
	if !repo.IsBare {
		if err := repo.UpdateSize(); err != nil {
			log.Error("UpdateSize: %v", err)
		}
		if err := repo.UpdateCounters(); err != nil {
			log.Error("UpdateCounters: %v", err)
		}
		if err := UpdateRepositoryLanguages(repo); err != nil {
			log.Error("UpdateRepositoryLanguages: %v", err)
		}
	}

	return repo, nil
}

