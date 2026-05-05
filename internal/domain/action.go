package domain

import (
	"context"
	"time"
)

// Action represents a user operation that can be logged and displayed in the news feed.
type Action struct {
	ID           int64
	UserID       int64
	OpType       int
	ActUserID    int64
	ActUserName  string
	RepoID       int64
	RepoUserName string
	RepoName     string
	RefName      string
	IsPrivate    bool
	Content      string
	Created      time.Time

	// UI convenience
	ActAvatar string `json:"act_avatar"`
}

// ActionsStore defines the interface for managing actions.
type ActionsStore interface {
	ListByUser(ctx context.Context, userID, actorID, afterID int64, isProfile bool) ([]*Action, error)
	CommitRepo(ctx context.Context, opts any) error
	TransferRepo(ctx context.Context, doer, oldOwner, newOwner *User, repo *Repository) error
	RenameRepo(ctx context.Context, doer, owner *User, oldRepoName string, repo *Repository) error
	MergePullRequest(ctx context.Context, doer, owner *User, repo *Repository, pull *Issue) error
	PushTag(ctx context.Context, opts any) error
	UpdateConfig(ctx context.Context, doer *User, section, content string) error
	ListAuditLog(ctx context.Context, page int) ([]*Action, error)
}
