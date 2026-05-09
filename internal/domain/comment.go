package domain

import (
	"context"
	"time"
)

// Comment represents a comment in commit and issue page.
type Comment struct {
	ID        int64
	Type      int
	PosterID  int64
	Poster    *User
	IssueID   int64
	Issue     *Issue
	CommitID  int64
	Line      int64
	Content   string
	Created   time.Time
	Updated   time.Time
	CommitSHA string
}

// CommentsStore defines the interface for managing comments.
type CommentsStore interface {
	GetByID(ctx context.Context, id int64) (*Comment, error)
	ListByIssueID(ctx context.Context, issueID int64) ([]*Comment, error)
}
