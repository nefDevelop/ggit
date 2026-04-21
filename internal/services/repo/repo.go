package repo

import (
	"context"
	"os"

	"github.com/cockroachdb/errors"
	"gorm.io/gorm"
	log "unknwon.dev/clog/v2"

	"gogs.io/gogs/internal/database"
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
	var repo *database.Repository
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Crear el registro en la base de datos
		var err error
		repo, err = s.db.Repositories().Create(ctx, owner.ID, database.CreateRepoOptions{
			Name:         opts.Name,
			Description:  opts.Description,
			Private:      opts.IsPrivate,
			EnableWiki:   true,
			EnableIssues: true,
			EnablePulls:  true,
			CreatedAt:    opts.CreatedAt,
		})
		if err != nil {
			return err
		}

		repo.Owner = owner
		repoPath := repo.RepoPath()

		// 2. Inicializar el repositorio Git
		err = initRepository(repoPath, doer, repo, opts)
		if err != nil {
			// Intentar limpiar el directorio si falló la init
			_ = os.RemoveAll(repoPath)
			return errors.Wrap(err, "init repository")
		}

		log.Trace("Repository created: %s/%s", owner.Name, repo.Name)

		// 3. Detectar lenguajes si se auto-inicializó
		if opts.AutoInit {
			if err := UpdateRepositoryLanguages(repo); err != nil {
				log.Error("UpdateRepositoryLanguages: %v", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return repo, nil
}
