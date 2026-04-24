package issue

import (
	"time"
)

// CreateIssueOptions contiene los datos necesarios para crear un issue.
type CreateIssueOptions struct {
	RepoID      int64     `json:"repo_id"`
	PosterID    int64     `json:"poster_id"`
	Title       string    `json:"title" binding:"required"`
	Content     string    `json:"content"`
	MilestoneID int64     `json:"milestone_id"`
	AssigneeID  int64     `json:"assignee_id"`
	LabelIDs    []int64   `json:"label_ids"`
	Attachments []string  `json:"attachments"` // UUIDs de adjuntos
	CreatedAt   time.Time `json:"created_at"`
}

// CreatePullRequestOptions contiene los datos necesarios para crear un pull request.
type CreatePullRequestOptions struct {
	CreateIssueOptions
	HeadRepoID int64  `json:"head_repo_id"`
	BaseRepoID int64  `json:"base_repo_id"`
	HeadBranch string `json:"head_branch"`
	BaseBranch string `json:"base_branch"`
	MergeBase  string `json:"merge_base"`
	Patch      []byte `json:"patch"`
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
