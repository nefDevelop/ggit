package repo

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
	"gorm.io/gorm"
	log "unknwon.dev/clog/v2"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/repox"
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
		return nil, errors.Wrap(err, "read repo root")
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
		return nil, errors.Newf("invalid or malicious repository name: %q", repoName)
	}

	owner, err := database.Handle.Users().GetByUsername(ctx, ownerName)
	if err != nil {
		return nil, errors.Wrap(err, "get owner")
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
			return errors.Wrap(err, "create db record")
		}
		repo = database.FromDomainRepository(domainRepo)

		repoPath := repo.RepoPath()
		if !strings.HasSuffix(repoPath, ".git") {
			repoPath += ".git"
		}

		if _, err := os.Stat(repoPath); os.IsNotExist(err) {
			return errors.Newf("repository folder does not exist: %s", repoPath)
		}

		// Fix hooks
		if err := repox.CreateDelegateHooks(repoPath); err != nil {
			return errors.Wrap(err, "create delegate hooks")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Update stats (commits, branches, etc.)
	if err := UpdateRepositoryLanguages(repo); err != nil {
		log.Error("UpdateRepositoryLanguages: %v", err)
	}

	return repo, nil
}
