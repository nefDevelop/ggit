package issue

import (
	"time"
)

// CreateIssueOptions contiene los datos necesarios para crear un issue.
type CreateIssueOptions struct {
	RepoID      int64
	PosterID    int64
	Title       string
	Content     string
	MilestoneID int64
	AssigneeID  int64
	LabelIDs    []int64
	Attachments []string // UUIDs de adjuntos
}

// CreatePullRequestOptions contiene los datos necesarios para crear un pull request.
type CreatePullRequestOptions struct {
	CreateIssueOptions
	HeadRepoID int64
	BaseRepoID int64
	HeadBranch string
	BaseBranch string
	MergeBase  string
	Patch      []byte
}

// CreateCommentOptions contiene los datos para crear un comentario.
type CreateCommentOptions struct {
	IssueID  int64
	PosterID int64
	Content  string
	Type     int // CommentType (Action, Comment, etc)
}

// UpdateIssueOptions contiene los datos para actualizar un issue.
type UpdateIssueOptions struct {
	Title       *string
	Content     *string
	MilestoneID *int64
	AssigneeID  *int64
	IsClosed    *bool
	Deadline    *time.Time
}
