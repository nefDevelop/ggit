package database

import (
	"context"

	"github.com/nef734/ggit/internal/domain"
	"gorm.io/gorm"
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
	if err := NewIssue(mRepo, mIssue, labelIDs, attachments); err != nil {
		return err
	}

	// NewIssue assigns the index on mIssue, which is a throwaway model built here
	// rather than the object the caller holds. Without this the caller keeps
	// reading zero, which is what the API answered as "number": 0.
	issue.Index = mIssue.Index
	issue.ID = mIssue.ID
	issue.Created = mIssue.Created
	issue.Updated = mIssue.Updated

	return nil
}

func (s *IssuesStore) GetByID(ctx context.Context, id int64) (*domain.Issue, error) {
	u, err := GetIssueByID(id)
	if err != nil {
		return nil, err
	}
	return u.ToDomain(), nil
}

func (s *IssuesStore) GetByIndex(ctx context.Context, repoID, index int64) (*domain.Issue, error) {
	u, err := GetIssueByIndex(repoID, index)
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
