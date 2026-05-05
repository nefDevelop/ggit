package domain

import (
	"context"
	"time"
)

// AccessToken represents a personal access token.
type AccessToken struct {
	ID                int64
	UserID            int64
	Name              string
	Sha1              string
	SHA256            string
	Created           time.Time
	Updated           time.Time
	HasRecentActivity bool
	HasUsed           bool
}

// AccessTokensStore defines the interface for managing access tokens.
type AccessTokensStore interface {
	Create(ctx context.Context, userID int64, name string) (*AccessToken, string, error)
	List(ctx context.Context, userID int64) ([]*AccessToken, error)
	DeleteByID(ctx context.Context, userID, id int64) error
	GetBySHA1(ctx context.Context, sha1 string) (*AccessToken, error)
	GetBySHA256(ctx context.Context, sha256 string) (*AccessToken, error)
	TouchByID(ctx context.Context, id int64) error
}
