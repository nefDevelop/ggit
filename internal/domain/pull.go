package domain

import (
	"context"
	"time"
)

// PullRequest represents a relation between pull request and repositories.
type PullRequest struct {
	ID             int64
	Type           int
	Status         int
	AttentionState int

	IssueID int64
	Issue   *Issue
	Index   int64

	HeadRepoID   int64
	HeadRepo     *Repository
	BaseRepoID   int64
	BaseRepo     *Repository
	HeadUserName string
	HeadBranch   string
	HeadCommitID string
	BaseBranch   string
	MergeBase    string

	HasMerged      bool
	MergedCommitID string
	MergerID       int64
	Merger         *User
	Merged         time.Time
}

// PullRequestsStore defines the interface for managing pull requests.
type PullRequestsStore interface {
	GetByID(ctx context.Context, id int64) (*PullRequest, error)
	GetByIssueID(ctx context.Context, issueID int64) (*PullRequest, error)
}
