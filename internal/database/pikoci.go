package database

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// ErrPikoCIRepoNotExist indica que el repositorio no tiene pipeline asociado.
var ErrPikoCIRepoNotExist = errors.New("pikoci repository link does not exist")

// PikoCIRepo vincula un repositorio con su pipeline en PikoCI. El token es el
// del webhook del resource, que solo se conoce al publicar el pipeline.
type PikoCIRepo struct {
	ID           int64  `gorm:"primaryKey"`
	RepoID       int64  `gorm:"uniqueIndex;not null"`
	PipelineName string `gorm:"type:VARCHAR(255);not null"`
	WebhookToken string `gorm:"type:VARCHAR(255);not null"`

	CreatedAt int64 `gorm:"autoCreateTime"`
	UpdatedAt int64 `gorm:"autoUpdateTime"`
}

// TableName fija el nombre de la tabla; GORM lo derivaría como "piko_ci_repo".
func (PikoCIRepo) TableName() string {
	return "pikoci_repo"
}

// IsErrPikoCIRepoNotExist informa si err es el centinela de vínculo ausente.
func IsErrPikoCIRepoNotExist(err error) bool {
	return errors.Is(err, ErrPikoCIRepoNotExist)
}

// PikoCIRepoStore gestiona los vínculos entre repositorios y pipelines.
type PikoCIRepoStore interface {
	GetByRepoID(ctx context.Context, repoID int64) (*PikoCIRepo, error)
	List(ctx context.Context, limit int) ([]*PikoCIRepo, error)
	Upsert(ctx context.Context, link *PikoCIRepo) error
	Delete(ctx context.Context, repoID int64) error
	Count(ctx context.Context) (int64, error)
}

// PikoCIRepos implementa PikoCIRepoStore.
type PikoCIRepos struct {
	db *gorm.DB
}

// PikoCIRepos devuelve el store de vínculos de PikoCI.
func (h *DB) PikoCIRepos() PikoCIRepoStore {
	return &PikoCIRepos{db: h.db}
}

func (s *PikoCIRepos) GetByRepoID(ctx context.Context, repoID int64) (*PikoCIRepo, error) {
	var link PikoCIRepo
	if err := s.db.WithContext(ctx).Where("repo_id = ?", repoID).First(&link).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPikoCIRepoNotExist
		}
		return nil, err
	}
	return &link, nil
}

func (s *PikoCIRepos) List(ctx context.Context, limit int) ([]*PikoCIRepo, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	var links []*PikoCIRepo
	err := s.db.WithContext(ctx).Order("updated_at DESC").Limit(limit).Find(&links).Error
	return links, err
}

func (s *PikoCIRepos) Upsert(ctx context.Context, link *PikoCIRepo) error {
	existing, err := s.GetByRepoID(ctx, link.RepoID)
	switch {
	case IsErrPikoCIRepoNotExist(err):
		return s.db.WithContext(ctx).Create(link).Error
	case err != nil:
		return err
	}

	existing.PipelineName = link.PipelineName
	existing.WebhookToken = link.WebhookToken
	return s.db.WithContext(ctx).Save(existing).Error
}

func (s *PikoCIRepos) Delete(ctx context.Context, repoID int64) error {
	return s.db.WithContext(ctx).Where("repo_id = ?", repoID).Delete(&PikoCIRepo{}).Error
}

func (s *PikoCIRepos) Count(ctx context.Context) (int64, error) {
	var count int64
	err := s.db.WithContext(ctx).Model(&PikoCIRepo{}).Count(&count).Error
	return count, err
}
