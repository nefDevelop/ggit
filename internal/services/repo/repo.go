package repo

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/domain"
	"github.com/nef734/ggit/internal/fsx"
	log "github.com/nef734/ggit/internal/log"
	"github.com/nef734/ggit/internal/repox"
	"gorm.io/gorm"
)

// RepositoryService define las operaciones de negocio sobre los repositorios.
type RepositoryService struct {
	db *database.DB
}

// Repos es el singleton del servicio de repositorios.
var Repos *RepositoryService

// Init inicializa el singleton del servicio.
func Init(db *database.DB) {
	Repos = &RepositoryService{db: db}
}

// Create crea un nuevo repositorio tanto en base de datos como en disco.
func (s *RepositoryService) Create(ctx context.Context, doer, owner *database.User, opts CreateOptions) (*database.Repository, error) {
	if opts.DefaultBranch == "" {
		opts.DefaultBranch = "main"
	}
	var repo *database.Repository

	finalRepoPath := repox.RepositoryPath(owner.Name, opts.Name)

	// Every attempt needs its own staging directory. A shared path derived from
	// the repository name, as "<name>.git.tmp", makes two concurrent creations of
	// the same name fight over one directory: the first rename succeeds and the
	// rest fail with "file exists", so none of them can report success.
	if err := os.MkdirAll(filepath.Dir(finalRepoPath), os.ModePerm); err != nil {
		return nil, fmt.Errorf("create repository parent directory: %w", err)
	}
	stagingDir, err := os.MkdirTemp(filepath.Dir(finalRepoPath), ".create-"+strings.ToLower(opts.Name)+"-")
	if err != nil {
		return nil, fmt.Errorf("create staging directory: %w", err)
	}
	// Safe in both directions: after a successful commit tempRepoPath has been
	// renamed away and stagingDir is empty, and a rollback leaves it to be swept.
	defer func() {
		_ = os.RemoveAll(stagingDir)
	}()

	tempRepoPath := filepath.Join(stagingDir, "repo.git")
	fsTx := fsx.NewTransaction()
	fsTx.StageCreate(tempRepoPath, finalRepoPath)
	defer fsTx.Rollback()

	// The repository is built on disk first and the row is inserted afterwards,
	// so no write transaction is ever held across git work.
	//
	// Holding one was what broke concurrent creation of the same name: the lock
	// stayed taken for the whole of git init, the delegate hooks, the initial
	// clone and the initial commit. On SQLite a transaction that read first and
	// writes later cannot be retried by the busy handler when it fails to promote,
	// so every one of the racers died with SQLITE_BUSY rather than letting the
	// unique index reject all but one.
	//
	// There is no atomicity across a row and a directory either way: the rename is
	// not part of the database transaction. The order below keeps the duplicate
	// check in the database, where the unique index enforces it, and pays for it
	// with an explicit cleanup when the insert loses.
	placeholder := &database.Repository{
		OwnerID:        owner.ID,
		LowerName:      strings.ToLower(opts.Name),
		Name:           opts.Name,
		Description:    opts.Description,
		Website:        opts.Website,
		DefaultBranch:  opts.DefaultBranch,
		IsPrivate:      opts.IsPrivate,
		IsUnlisted:     opts.IsUnlisted,
		EnableWiki:     true,
		EnableIssues:   true,
		EnablePulls:    true,
		WikiFolderName: "docs",
		Owner:          owner,
	}

	log.Trace("Initializing Git repository at temp path: %s", tempRepoPath)
	if err = initRepository(tempRepoPath, doer, placeholder, opts); err != nil {
		log.Error("Failed to initialize Git repository: %v", err)
		return nil, fmt.Errorf("init repository"+": %w", err)
	}

	log.Trace("Committing filesystem transaction for %s/%s", owner.Name, opts.Name)
	if err = fsTx.Commit(); err != nil {
		log.Error("Failed to commit filesystem transaction: %v", err)
		return nil, fmt.Errorf("commit filesystem transaction"+": %w", err)
	}

	// From here the directory is in place, so a failed insert has to remove it
	// again or the forge is left with an unregistered repository on disk.
	var domainRepo *domain.Repository
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		log.Trace("Creating repository record in DB for %s/%s", owner.Name, opts.Name)
		created, err := s.db.WithTransaction(tx).Repositories().Create(ctx, owner.ID, database.CreateRepoOptions{
			Name:           opts.Name,
			Description:    opts.Description,
			Website:        opts.Website,
			DefaultBranch:  opts.DefaultBranch,
			Private:        opts.IsPrivate,
			IsUnlisted:     opts.IsUnlisted,
			AutoInit:       opts.AutoInit,
			Gitignores:     opts.Gitignores,
			License:        opts.License,
			Readme:         opts.Readme,
			IsArchived:     opts.IsArchived,
			ProjectFolder:  opts.ProjectFolder,
			EnableWiki:     true,
			EnableIssues:   true,
			EnablePulls:    true,
			CreatedAt:      opts.CreatedAt,
			WikiFolderName: "docs", // Default
		})
		if err != nil {
			log.Error("Failed to create repository record: %v", err)
			return err
		}
		domainRepo = created
		return nil
	})
	if err != nil {
		if removeErr := os.RemoveAll(finalRepoPath); removeErr != nil {
			log.Error("Failed to remove %s after a failed insert: %v", finalRepoPath, removeErr)
		}
		return nil, err
	}

	repo = database.FromDomainRepository(domainRepo)
	repo.Owner = owner

	log.Trace("Repository created successfully: %s/%s", owner.Name, repo.Name)

	if opts.AutoInit {
		_ = s.db.Repositories().Touch(ctx, repo.ID)
	}

	// 4. Detectar lenguajes si se auto-inicializó (ya en la ruta final)
	if opts.AutoInit {
		if err := UpdateRepositoryLanguages(repo); err != nil {
			log.Error("UpdateRepositoryLanguages: %v", err)
		}
	}

	return repo, nil
}
