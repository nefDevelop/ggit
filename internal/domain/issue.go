package domain

import (
	"context"
	"time"
)

// Issue represents an issue or pull request of repository.
type Issue struct {
	ID          int64
	RepoID      int64
	Repo        *Repository
	Index       int64
	PosterID    int64
	Poster      *User
	Title       string
	Content     string
	Labels      []*Label
	MilestoneID int64
	Milestone   *Milestone
	Priority    int
	AssigneeID  int64
	Assignee    *User
	IsClosed    bool
	IsPull      bool
	NumComments int
	Deadline    time.Time
	Created     time.Time
	Updated     time.Time
}

// Label represents a label for issues.
type Label struct {
	ID              int64
	RepoID          int64
	Name            string
	Color           string
	NumIssues       int
	NumClosedIssues int
}

// Milestone represents a milestone for issues.
type Milestone struct {
	ID              int64
	RepoID          int64
	Name            string
	Content         string
	IsClosed        bool
	NumIssues       int
	NumClosedIssues int
	Completeness    int
	Deadline        time.Time
	ClosedDate      time.Time
}

// IssuesStore defines the interface for managing issues.
type IssuesStore interface {
	Create(ctx context.Context, repo *Repository, issue *Issue, labelIDs []int64, attachments []string) error
	GetByID(ctx context.Context, id int64) (*Issue, error)
	GetByIndex(ctx context.Context, repoID, index int64) (*Issue, error)
	ListByRepoID(ctx context.Context, repoID int64, page, pageSize int) ([]*Issue, error)
}
