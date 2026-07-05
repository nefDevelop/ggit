package repo

import (
	"context"
	"fmt"

	"github.com/nef734/ggit/internal/database"
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
	tempRepoPath := finalRepoPath + ".tmp"
	fsTx := fsx.NewTransaction()
	fsTx.StageCreate(tempRepoPath, finalRepoPath)
	defer fsTx.Rollback()

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		log.Trace("Creating repository record in DB for %s/%s", owner.Name, opts.Name)
		// 1. Crear el registro en la base de datos
		domainRepo, err := s.db.WithTransaction(tx).Repositories().Create(ctx, owner.ID, database.CreateRepoOptions{
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
		repo = database.FromDomainRepository(domainRepo)

		repo.Owner = owner

		// 2. Inicializar el repositorio Git en ruta temporal
		log.Trace("Initializing Git repository at temp path: %s", tempRepoPath)
		err = initRepository(tempRepoPath, doer, repo, opts)
		if err != nil {
			log.Error("Failed to initialize Git repository: %v", err)
			return fmt.Errorf("init repository"+": %w", err)
		}

		// 3. Commit del FS dentro de la transacción de DB para atomicidad total
		log.Trace("Committing filesystem transaction for %s/%s", owner.Name, repo.Name)
		if err := fsTx.Commit(); err != nil {
			log.Error("Failed to commit filesystem transaction: %v", err)
			return fmt.Errorf("commit filesystem transaction"+": %w", err)
		}

		log.Trace("Repository created successfully: %s/%s", owner.Name, repo.Name)
		return nil
	})
	if err != nil {
		return nil, err
	}

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
