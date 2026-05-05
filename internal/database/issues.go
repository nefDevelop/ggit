package database

import (
	"context"

	"gorm.io/gorm"

	"github.com/nef734/ggit/internal/domain"
)

// IssuesStore es el almacén de datos para issues utilizando GORM.
type IssuesStore struct {
	db *gorm.DB
}

func newIssuesStore(db *gorm.DB) *IssuesStore {
	return &IssuesStore{db: db}
}

// Create crea un nuevo issue en la base de datos de forma transaccional.
func (s *IssuesStore) Create(ctx context.Context, repo *domain.Repository, issue *domain.Issue, labelIDs []int64, attachments []string) error {
	mRepo := FromDomainRepository(repo)
	mIssue := FromDomainIssue(issue)
	return NewIssue(mRepo, mIssue, labelIDs, attachments)
}

func (s *IssuesStore) GetByID(ctx context.Context, id int64) (*domain.Issue, error) {
	u, err := GetIssueByID(id)
	if err != nil {
		return nil, err
	}
	return u.ToDomain(), nil
}

func (s *IssuesStore) ListByRepoID(ctx context.Context, repoID int64, page, pageSize int) ([]*domain.Issue, error) {
	// TODO: Implement using pure GORM
	return nil, nil
}

// database.go necesita conocer este nuevo store.
