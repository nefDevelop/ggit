package repo

import (
	"context"

	"github.com/cockroachdb/errors"
	"gorm.io/gorm"
	log "unknwon.dev/clog/v2"

	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/fsx"
	"github.com/nef734/ggit/internal/repox"
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
			EnableWiki:     true,
			EnableIssues:   true,
			EnablePulls:    true,
			CreatedAt:      opts.CreatedAt,
			WikiFolderName: "docs", // Default
		})
		if err != nil {
			return err
		}
		repo = database.FromDomainRepository(domainRepo)

		repo.Owner = owner

		// 2. Inicializar el repositorio Git en ruta temporal
		err = initRepository(tempRepoPath, doer, repo, opts)
		if err != nil {
			return errors.Wrap(err, "init repository")
		}

		log.Trace("Repository DB record created and staged on FS: %s/%s", owner.Name, repo.Name)
		return nil
	})
	if err != nil {
		return nil, err
	}

	// 3. Commit del FS después de que la DB haya confirmado
	if err := fsTx.Commit(); err != nil {
		return nil, errors.Wrap(err, "commit filesystem transaction")
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

