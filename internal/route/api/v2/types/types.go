package types

import (
	"encoding/json"
	"time"
)

// USER MODELS

// UserResponse represents a user in the API v2.
type UserResponse struct {
	ID              int64  `json:"id"`
	Username        string `json:"username"`
	FullName        string `json:"full_name"`
	Email           string `json:"email"`
	AvatarURL       string    `json:"avatar_url"`
	IsAdmin         bool      `json:"is_admin"`
	DisableNewsFeed bool      `json:"disable_news_feed"`
	ProhibitLogin    bool      `json:"prohibit_login"`
	AllowGitHook     bool      `json:"allow_git_hook"`
	AllowImportLocal bool      `json:"allow_import_local"`
	MaxRepoCreation  int       `json:"max_repo_creation"`
	CreatedAt       time.Time `json:"created_at"`
	}


// UserEmailResponse represents an email address of a user.
type UserEmailResponse struct {
	Email    string `json:"email"`
	Verified bool   `json:"verified"`
	Primary  bool   `json:"primary"`
}

// UserAccessTokenResponse represents an access token of a user.
type UserAccessTokenResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Sha1 string `json:"sha1"`
}

// LoginSessionResponse represents a login session of a user.
type LoginSessionResponse struct {
	ID          int64     `json:"id"`
	SID         string    `json:"sid"` // Masked or truncated
	UserAgent   string    `json:"user_agent"`
	RemoteAddr  string    `json:"remote_addr"`
	Current     bool      `json:"current"`
	CreatedUnix int64     `json:"created_unix"`
	UpdatedUnix int64     `json:"updated_unix"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// REPOSITORY MODELS

// RepositoryResponse represents a repository in the API v2.
type RepositoryResponse struct {
	ID          int64         `json:"id"`
	Name          string        `json:"name"`
	ProjectFolder string        `json:"project_folder"`
	FullName    string        `json:"full_name"`
	Description string        `json:"description"`
	Owner       *UserResponse `json:"owner"`
	AvatarURL   string        `json:"avatar_url"`
	NumStars      int                   `json:"num_stars"`
	NumForks      int                   `json:"num_forks"`
	NumWatches    int                   `json:"num_watches"`
	OpenIssues    int                   `json:"open_issues_count"`
	OpenPulls     int                   `json:"open_pulls_count"`
	NumCommits    int64                 `json:"num_commits"`
	NumBranches   int                   `json:"num_branches"`
	NumReleases   int                   `json:"num_releases_count"`
	IsPrivate     bool                  `json:"is_private"`
	HTMLURL       string                `json:"html_url"`
	CloneURL      string                `json:"clone_url"`
	SSHURL        string                `json:"ssh_url"`
	DefaultBranch  string                `json:"default_branch"`
	WikiFolderName string                `json:"wiki_folder_name"`
	ADRFolderName  string                `json:"adr_folder_name"`
	EnableWiki            bool   `json:"enable_wiki"`
	EnableExternalWiki    bool   `json:"enable_external_wiki"`
	ExternalWikiURL       string `json:"external_wiki_url"`
	EnableIssues          bool   `json:"enable_issues"`
	EnableExternalTracker bool   `json:"enable_external_tracker"`
	ExternalTrackerURL    string `json:"external_tracker_url"`
	ExternalTrackerFormat string `json:"external_tracker_format"`
	ExternalTrackerStyle  string `json:"external_tracker_style"`
	EnablePulls           bool   `json:"enable_pulls"`
	PullsIgnoreWhitespace bool   `json:"pulls_ignore_whitespace"`
	PullsAllowRebase      bool   `json:"pulls_allow_rebase"`
	CreatedAt      time.Time             `json:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at"`
	Permissions   *RepositoryPermission `json:"permissions,omitempty"`
	Parent        *RepositoryResponse   `json:"parent,omitempty"`
	Topics        []string              `json:"topics"`
	
	// Activity Timestamps
	LastCommitAt   time.Time `json:"last_commit_at"`
	LastIssueAt    time.Time `json:"last_issue_at"`
	LastPullAt     time.Time `json:"last_pull_at"`
	LastBranchAt   time.Time `json:"last_branch_at"`

	// Sentinel Info
	LastHealthCheckAt time.Time `json:"last_health_check_at"`
}

// RepositoryBranchResponse represents a branch in a repository.
type RepositoryBranchResponse struct {
	Name   string                `json:"name"`
	Commit *WebhookPayloadCommit `json:"commit"`
}

// AttachmentResponse represents an attachment in the API v2.
type AttachmentResponse struct {
	ID          int64  `json:"id"`
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	DownloadURL string `json:"download_url"`
}

// RepositoryReleaseResponse represents a release of a repository.
type RepositoryReleaseResponse struct {
	ID              int64                 `json:"id"`
	TagName         string                `json:"tag_name"`
	TargetCommitish string                `json:"target_commitish"`
	Name            string                `json:"name"`
	Body            string                `json:"body"`
	Draft           bool                  `json:"draft"`
	Prerelease      bool                  `json:"prerelease"`
	Author          *UserResponse         `json:"author"`
	Attachments     []*AttachmentResponse `json:"attachments"`
	CreatedAt       int64                 `json:"created_at"`
}

// RepositoryHookResponse represents a webhook for a repository.
type RepositoryHookResponse struct {
	ID        int64             `json:"id"`
	Type      string            `json:"type"`
	Config    map[string]string `json:"config"`
	Events    []string          `json:"events"`
	Active    bool              `json:"active"`
	UpdatedAt time.Time         `json:"updated_at"`
	CreatedAt time.Time         `json:"created_at"`
}

// RepositoryDeployKeyResponse represents an SSH deploy key for a repository.
type RepositoryDeployKeyResponse struct {
	ID        int64     `json:"id"`
	Key       string    `json:"key"`
	URL       string    `json:"url"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	ReadOnly  bool      `json:"read_only"`
}

// RepositoryPermission represents the permissions for a repository.
type RepositoryPermission struct {
	Admin bool `json:"admin"`
	Push  bool `json:"push"`
	Pull  bool `json:"pull"`
}

// RepositoryCollaboratorResponse represents a collaborator in a repository.
type RepositoryCollaboratorResponse struct {
	*UserResponse
	Permissions RepositoryPermission `json:"permissions"`
}

// ISSUE MODELS

// IssueStateType represents the state of an issue or pull request.
type IssueStateType string

const (
	IssueStateOpen   IssueStateType = "open"
	IssueStateClosed IssueStateType = "closed"
)

// IssueLabelResponse represents a label for an issue.
type IssueLabelResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
	URL   string `json:"url"`
}

// IssueMilestoneResponse represents a milestone for a repository.
type IssueMilestoneResponse struct {
	ID           int64          `json:"id"`
	Title        string         `json:"title"`
	Description  string         `json:"description"`
	State        IssueStateType `json:"state"`
	OpenIssues   int            `json:"open_issues"`
	ClosedIssues int            `json:"closed_issues"`
	Completeness int            `json:"completeness"`
	HillState    int            `json:"hill_state"`
	Deadline     *time.Time     `json:"due_on"`
}

// CreateMilestoneOption represents options for creating a milestone.
type CreateMilestoneOption struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Deadline    *time.Time `json:"due_on"`
	HillState   int        `json:"hill_state"`
}

// EditMilestoneOption represents options for editing a milestone.
type EditMilestoneOption struct {
	Title       *string     `json:"title"`
	Description *string     `json:"description"`
	State       *string     `json:"state"`
	Deadline    *time.Time  `json:"due_on"`
	HillState   *int        `json:"hill_state"`
}

// PullRequestMetaResponse represents metadata for a pull request within an issue.
type PullRequestMetaResponse struct {
	HasMerged bool       `json:"merged"`
	MergedAt  *time.Time `json:"merged_at"`
}

// IssueResponse represents an issue in the API v2.
type IssueResponse struct {
	ID          int64                    `json:"id"`
	Index       int64                    `json:"number"`
	Poster      *UserResponse            `json:"user"`
	Title       string                   `json:"title"`
	Body        string                   `json:"body"`
	BodyHTML    string                   `json:"body_html"`
	Labels      []*IssueLabelResponse    `json:"labels"`
	Milestone   *IssueMilestoneResponse  `json:"milestone"`
	Assignee    *UserResponse            `json:"assignee"`
	State       IssueStateType           `json:"state"`
	Comments    int                      `json:"comments"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
	PullRequest *PullRequestMetaResponse `json:"pull_request"`
}

// IssueCommentResponse represents a comment on an issue.
type IssueCommentResponse struct {
	ID        int64         `json:"id"`
	HTMLURL   string        `json:"html_url"`
	Poster    *UserResponse `json:"user"`
	Body      string        `json:"body"`
	BodyHTML  string        `json:"body_html"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// PULL REQUEST MODELS

// PullRequestResponse represents a pull request in the API v2.
type PullRequestResponse struct {
	ID             int64               `json:"id"`
	Index          int64               `json:"number"`
	Poster         *UserResponse       `json:"user"`
	Title          string              `json:"title"`
	Body           string              `json:"body"`
	BodyHTML       string              `json:"body_html"`
	Labels         []*IssueLabelResponse `json:"labels"`
	Milestone      *IssueMilestoneResponse `json:"milestone"`
	Assignee       *UserResponse       `json:"assignee"`
	State          IssueStateType      `json:"state"`
	Comments       int                 `json:"comments"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
	HeadBranch     string              `json:"head_branch"`
	HeadRepo       *RepositoryResponse `json:"head_repo"`
	BaseBranch     string              `json:"base_branch"`
	BaseRepo       *RepositoryResponse `json:"base_repo"`
	HTMLURL        string              `json:"html_url"`
	Mergeable      *bool               `json:"mergeable"`
	HasMerged      bool                `json:"merged"`
	MergedAt       *time.Time          `json:"merged_at"`
	MergedCommitID *string             `json:"merged_commit_sha"`
	MergedBy       *UserResponse       `json:"merged_by"`
	AttentionState int                 `json:"attention_state"`
	CombinedStatus string              `json:"combined_status"`
}

// PullReviewCommentResponse represents an inline comment on a pull request review.
type PullReviewCommentResponse struct {
	ID         int64     `json:"id"`
	ReviewID   int64     `json:"review_id"`
	TreePath   string    `json:"tree_path"`
	Line       int64     `json:"line"`
	CommitID   string    `json:"commit_id"`
	Content    string    `json:"content"`
	IsResolved bool      `json:"is_resolved"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// PullReviewResponse represents a review on a pull request.
type PullReviewResponse struct {
	ID         int64                        `json:"id"`
	Type       int                          `json:"type"`
	Reviewer   *UserResponse                `json:"reviewer"`
	Content    string                       `json:"content"`
	Comments   []*PullReviewCommentResponse `json:"comments"`
	CreatedAt  time.Time                    `json:"created_at"`
}

// ORGANIZATION MODELS

// OrganizationResponse represents an organization in the API v2.
type OrganizationResponse struct {
	ID          int64         `json:"id"`
	Username    string        `json:"username"`
	FullName    string        `json:"full_name"`
	Description string        `json:"description"`
	Owner       *UserResponse `json:"owner"`
	NumTeams    int           `json:"num_teams"`
	NumMembers  int           `json:"num_members"`
}

// OTHER MODELS

// AuthSourceResponse represents a login source in the API v2.
type AuthSourceResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	TypeName    string    `json:"type_name"`
	IsActivated bool      `json:"is_activated"`
	CreatedAt   time.Time `json:"created_at"`
	RawConfig   string    `json:"raw_config,omitempty"`
}

// WebhookResponse represents a webhook in the API v2.
type WebhookResponse struct {
	ID          int64    `json:"id"`
	URL         string   `json:"url"`
	ContentType string   `json:"content_type"`
	Secret      string   `json:"secret"`
	Events      []string `json:"events"`
	IsSSL       bool     `json:"is_ssl"`
	IsActive    bool     `json:"active"`
	LastStatus  int      `json:"last_status"`
	Updated     int64    `json:"updated"`
}

// CreateWebhookRequest represents a request to create a new webhook.
type CreateWebhookRequest struct {
	URL         string   `json:"url" binding:"required"`
	ContentType string   `json:"content_type" binding:"required"`
	Secret      string   `json:"secret"`
	Active      bool     `json:"active"`
	Events      []string `json:"events"`
}

// UpdateWebhookRequest represents a request to update an existing webhook.
type UpdateWebhookRequest struct {
	URL         string   `json:"url"`
	ContentType string   `json:"content_type"`
	Secret      string   `json:"secret"`
	Active      *bool    `json:"active"`
	Events      []string `json:"events"`
}

// WebhookHistoryResponse represents a delivery task of a webhook.

type WebhookHistoryResponse struct {
	ID              int64     `json:"id"`
	UUID            string    `json:"uuid"`
	EventType       string    `json:"event_type"`
	IsSucceed       bool      `json:"is_succeed"`
	Delivered       time.Time `json:"delivered_at"`
	RequestHeaders  map[string]string `json:"request_headers"`
	RequestBody     string            `json:"request_body"`
	ResponseHeaders map[string]string `json:"response_headers"`
	ResponseBody    string            `json:"response_body"`
	ResponseStatus  int               `json:"response_status"`
}

// NoticeResponse represents a system notice in the API v2.
type NoticeResponse struct {
	ID          int64     `json:"id"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
// AuditLogResponse represents a system action/audit log entry in the API v2.
type AuditLogResponse struct {
	ID        int64     `json:"id"`
	OpType    int       `json:"op_type"`
	OpName    string    `json:"op_name"`
	ActorID   int64     `json:"actor_id"`
	ActorName string    `json:"actor_name"`
	RepoID    int64     `json:"repo_id"`
	RepoPath  string    `json:"repo_path"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// CommitStatusResponse represents a commit status in the API v2.
type CommitStatusResponse struct {
	ID          int64     `json:"id"`
	State       string    `json:"state"`
	TargetURL   string    `json:"target_url"`
	Description string    `json:"description"`
	Context     string    `json:"context"`
	Creator     *UserResponse `json:"creator"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CombinedCommitStatusResponse represents combined commit statuses for a SHA.
type CombinedCommitStatusResponse struct {
	State    string                  `json:"state"`
	Statuses []*CommitStatusResponse `json:"statuses"`
	Repo     *RepositoryResponse     `json:"repository"`
	SHA      string                  `json:"sha"`
}

// SystemStatusResponse represents the system health statistics in the API v2.
type SystemStatusResponse struct {
	Uptime       string `json:"uptime"`
	NumGoroutine int    `json:"num_goroutine"`
	MemAllocated string `json:"mem_allocated"`
	MemTotal     string `json:"mem_total"`
	MemSys       string `json:"mem_sys"`
	Lookups      uint64 `json:"lookups"`
	MemMallocs   uint64 `json:"mem_mallocs"`
	MemFrees     uint64 `json:"mem_frees"`
	HeapAlloc    string `json:"heap_alloc"`
	HeapSys      string `json:"heap_sys"`
	HeapIdle     string `json:"heap_idle"`
	HeapInuse    string `json:"heap_inuse"`
	HeapReleased string `json:"heap_released"`
	HeapObjects  uint64 `json:"heap_objects"`
	StackInuse   string `json:"stack_inuse"`
	StackSys     string `json:"stack_sys"`
	MSpanInuse   string `json:"mspan_inuse"`
	MSpanSys     string `json:"mspan_sys"`
	MCacheInuse  string `json:"mcache_inuse"`
	MCacheSys    string `json:"mcache_sys"`
	BuckHashSys  string `json:"buck_hash_sys"`
	GCSys        string `json:"gc_sys"`
	OtherSys     string `json:"other_sys"`
	NextGC       string `json:"next_gc"`
	LastGC       string `json:"last_gc"`
	PauseTotalNs string `json:"pause_total_ns"`
	PauseNs      string `json:"pause_ns"`
	NumGC        uint32 `json:"num_gc"`
}

// WEBHOOK PAYLOADS

type WebhookPusherType string

const (
	WebhookPusherTypeUser WebhookPusherType = "user"
	WebhookPusherTypeOrg  WebhookPusherType = "organization"
)

type WebhookPayloadUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	UserName string `json:"username"`
}

type WebhookPayloadCommit struct {
	ID        string              `json:"id"`
	Message   string              `json:"message"`
	URL       string              `json:"url"`
	Author    *WebhookPayloadUser `json:"author"`
	Committer *WebhookPayloadUser `json:"committer"`
	Added     []string            `json:"added"`
	Removed   []string            `json:"removed"`
	Modified  []string            `json:"modified"`
	Timestamp time.Time           `json:"timestamp"`
}

type WebhookIssueAction string

const (
	WebhookIssueOpened       WebhookIssueAction = "opened"
	WebhookIssueClosed       WebhookIssueAction = "closed"
	WebhookIssueReopened     WebhookIssueAction = "reopened"
	WebhookIssueEdited       WebhookIssueAction = "edited"
	WebhookIssueAssigned     WebhookIssueAction = "assigned"
	WebhookIssueUnassigned   WebhookIssueAction = "unassigned"
	WebhookIssueLabelUpdated WebhookIssueAction = "label_updated"
	WebhookIssueLabelCleared WebhookIssueAction = "label_cleared"
	WebhookIssueMilestoned   WebhookIssueAction = "milestoned"
	WebhookIssueDemilestoned WebhookIssueAction = "demilestoned"
	WebhookIssueSynchronized WebhookIssueAction = "synchronized"
)

type WebhookIssueCommentAction string

const (
	WebhookIssueCommentCreated WebhookIssueCommentAction = "created"
	WebhookIssueCommentEdited  WebhookIssueCommentAction = "edited"
	WebhookIssueCommentDeleted WebhookIssueCommentAction = "deleted"
)

type WebhookReleaseAction string

const WebhookReleasePublished WebhookReleaseAction = "published"

type WebhookChangesFromPayload struct {
	From string `json:"from"`
}

type WebhookChangesPayload struct {
	Title *WebhookChangesFromPayload `json:"title,omitempty"`
	Body  *WebhookChangesFromPayload `json:"body,omitempty"`
}

type WebhookPayloader interface {
	JSONPayload() ([]byte, error)
}

type WebhookCreatePayload struct {
	Ref           string              `json:"ref"`
	RefType       string              `json:"ref_type"`
	Sha           string              `json:"sha"`
	DefaultBranch string              `json:"default_branch"`
	Repo          *RepositoryResponse `json:"repository"`
	Sender        *UserResponse       `json:"sender"`
}

func (p *WebhookCreatePayload) JSONPayload() ([]byte, error) {
	return json.Marshal(p)
}

type WebhookDeletePayload struct {
	Ref        string              `json:"ref"`
	RefType    string              `json:"ref_type"`
	PusherType WebhookPusherType   `json:"pusher_type"`
	Repo       *RepositoryResponse `json:"repository"`
	Sender     *UserResponse       `json:"sender"`
}

func (p *WebhookDeletePayload) JSONPayload() ([]byte, error) {
	return json.Marshal(p)
}

type WebhookForkPayload struct {
	Forkee *RepositoryResponse `json:"forkee"`
	Repo   *RepositoryResponse `json:"repository"`
	Sender *UserResponse       `json:"sender"`
}

func (p *WebhookForkPayload) JSONPayload() ([]byte, error) {
	return json.Marshal(p)
}

type WebhookPushPayload struct {
	Ref        string                  `json:"ref"`
	Before     string                  `json:"before"`
	After      string                  `json:"after"`
	CompareURL string                  `json:"compare_url"`
	Commits    []*WebhookPayloadCommit `json:"commits"`
	Repo       *RepositoryResponse     `json:"repository"`
	Pusher     *UserResponse           `json:"pusher"`
	Sender     *UserResponse           `json:"sender"`
}

func (p *WebhookPushPayload) JSONPayload() ([]byte, error) {
	return json.Marshal(p)
}

type WebhookIssuesPayload struct {
	Action     WebhookIssueAction     `json:"action"`
	Index      int64                  `json:"number"`
	Issue      *IssueResponse         `json:"issue"`
	Changes    *WebhookChangesPayload `json:"changes,omitempty"`
	Repository *RepositoryResponse    `json:"repository"`
	Sender     *UserResponse          `json:"sender"`
}

func (p *WebhookIssuesPayload) JSONPayload() ([]byte, error) {
	return json.Marshal(p)
}

type WebhookIssueCommentPayload struct {
	Action     WebhookIssueCommentAction `json:"action"`
	Issue      *IssueResponse            `json:"issue"`
	Comment    *IssueCommentResponse     `json:"comment"`
	Changes    *WebhookChangesPayload    `json:"changes,omitempty"`
	Repository *RepositoryResponse       `json:"repository"`
	Sender     *UserResponse             `json:"sender"`
}

func (p *WebhookIssueCommentPayload) JSONPayload() ([]byte, error) {
	return json.Marshal(p)
}

type WebhookPullRequestPayload struct {
	Action      WebhookIssueAction     `json:"action"`
	Index       int64                  `json:"number"`
	PullRequest *PullRequestResponse   `json:"pull_request"`
	Changes     *WebhookChangesPayload `json:"changes,omitempty"`
	Repository  *RepositoryResponse    `json:"repository"`
	Sender      *UserResponse          `json:"sender"`
}

func (p *WebhookPullRequestPayload) JSONPayload() ([]byte, error) {
	return json.Marshal(p)
}

type WebhookReleasePayload struct {
	Action     WebhookReleaseAction       `json:"action"`
	Release    *RepositoryReleaseResponse `json:"release"`
	Repository *RepositoryResponse        `json:"repository"`
	Sender     *UserResponse              `json:"sender"`
}

func (p *WebhookReleasePayload) JSONPayload() ([]byte, error) {
	return json.Marshal(p)
}
