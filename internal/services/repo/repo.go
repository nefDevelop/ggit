package repo

import (
	"context"
	"fmt"
	"os"

	"github.com/cockroachdb/errors"
	"gorm.io/gorm"
	log "unknwon.dev/clog/v2"

	"gogs.io/gogs/internal/conf"
	"gogs.io/gogs/internal/database"
	"gogs.io/gogs/internal/osx"
	"gogs.io/gogs/internal/process"
	"gogs.io/gogs/internal/repox"
)

// RepositoryService es la interfaz para las operaciones de negocio con repositorios.
type RepositoryService interface {
	Create(ctx context.Context, doer, owner *database.User, opts CreateOptions) (*database.Repository, error)
}

type repositoryService struct {
	db *database.DB
}

// Repos es el singleton global para el servicio de repositorios.
var Repos RepositoryService

// Init inicializa el servicio de repositorios.
func Init(db *database.DB) {
	Repos = &repositoryService{
		db: db,
	}
}

func (s *repositoryService) Create(ctx context.Context, doer, owner *database.User, opts CreateOptions) (*database.Repository, error) {
	repoPath := repox.RepositoryPath(owner.Name, opts.Name)
	if osx.Exist(repoPath) {
		return nil, errors.Errorf("repository directory already exists: %s", repoPath)
	}

	if !owner.CanCreateRepo() {
		return nil, database.ErrReachLimitOfRepo{Limit: owner.MaxNumRepos()}
	}

	var repo *database.Repository
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error
		// 1. Crear el registro en la base de datos (GORM)
		repo, err = s.db.WithTransaction(tx).Repositories().Create(ctx, owner.ID, database.CreateRepoOptions{
			Name:          opts.Name,
			Description:   opts.Description,
			DefaultBranch: conf.Repository.DefaultBranch,
			Private:       opts.IsPrivate,
			EnableWiki:    true,
			EnableIssues:  true,
			EnablePulls:   true,
		})
		if err != nil {
			return err
		}

		repo.Owner = owner

		// 2. Inicializar el repositorio Git
		err = initRepository(repoPath, doer, repo, opts)
		if err != nil {
			// Intentar limpiar el directorio si falló la init
			_ = os.RemoveAll(repoPath)
			return err
		}

		// 3. Actualizar la información de Git (update-server-info)
		_, stderr, err := process.ExecDir(-1,
			repoPath, fmt.Sprintf("CreateRepository 'git update-server-info': %s", repoPath),
			"git", "update-server-info")
		if err != nil {
			return errors.Newf("CreateRepository 'git update-server-info': %s", stderr)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 4. Actualizar preferencia de visibilidad del usuario
	err = s.db.Users().Update(ctx, owner.ID, database.UpdateUserOptions{
		LastRepoVisibility: &repo.IsPrivate,
	})
	if err != nil {
		log.Error("Failed to update user last repo visibility: %v", err)
	}

	return repo, nil
}
