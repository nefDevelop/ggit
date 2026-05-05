package domain

import "context"

// SecretPattern represents a regex pattern to detect secrets in code.
type SecretPattern struct {
	ID          int64
	Name        string
	Regex       string
	Description string
	IsGlobal    bool
	RepoID      int64 // 0 for global patterns
}

// SecretPatternsStore defines the interface for managing secret patterns.
type SecretPatternsStore interface {
	Create(ctx context.Context, p *SecretPattern) error
	ListByRepo(ctx context.Context, repoID int64) ([]*SecretPattern, error)
	ListGlobal(ctx context.Context) ([]*SecretPattern, error)
	DeleteByID(ctx context.Context, id int64) error
	InitDefaultPatterns(ctx context.Context) error
}
