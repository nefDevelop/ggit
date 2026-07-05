package domain

import (
	"context"
	"time"
)

// Action represents a user operation that can be logged and displayed in the news feed.
type Action struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	OpType       int       `json:"op_type"`
	ActUserID    int64     `json:"act_user_id"`
	ActUserName  string    `json:"act_user_name"`
	RepoID       int64     `json:"repo_id"`
	RepoUserName string    `json:"repo_user_name"`
	RepoName     string    `json:"repo_name"`
	RefName      string    `json:"ref_name"`
	IsPrivate    bool      `json:"is_private"`
	Content      string    `json:"content"`
	Created      time.Time `json:"created_at"`

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
