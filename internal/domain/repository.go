package domain

import (
	"context"
	"time"
)

// Repository represents a git repository.
type Repository struct {
	ID              int64
	OwnerID         int64
	Owner           *User
	LowerName       string
	Name            string
	Description     string
	Website         string
	DefaultBranch   string
	Size            int64
	ProjectFolder   string
	UseCustomAvatar bool

	// Counters
	NumWatches          int
	NumStars            int
	NumForks            int
	NumIssues           int
	NumClosedIssues     int
	NumPulls            int
	NumClosedPulls      int
	NumMilestones       int
	NumClosedMilestones int
	NumCommits          int64
	NumBranches         int

	IsPrivate     bool
	IsFork        bool
	IsUnlisted    bool
	IsBare        bool
	IsArchived    bool
	IsMaintenance bool

	// Advanced settings
	EnableWiki                          bool
	AllowPublicWiki                     bool
	EnableExternalWiki                  bool
	ExternalWikiURL                     string
	WikiFolderName                      string
	ADRFolderName                       string
	EnableIssues                        bool
	AllowPublicIssues                   bool
	EnableExternalTracker               bool
	ExternalTrackerURL                  string
	ExternalTrackerFormat               string
	ExternalTrackerStyle                string
	EnablePulls                         bool
	PullsIgnoreWhitespace               bool
	PullsAllowRebase                    bool
	PullsAllowManualMerge               bool
	BlockMergeOnPendingReviews          bool
	BlockMergeOnUnresolvedConversations bool
	RequiredApprovals                   int

	PrimaryLanguage string
	Languages       string

	// Secret Guardian
	EnableSecretGuardian bool
	SecretGuardianStrict bool

	// Janitor & Sentinel
	JanitorLastCheckUnix int64
	LastHealthCheckUnix  int64
	CommitsSinceCheck    int

	// Activity Timestamps
	LastCommitUnix int64
	LastIssueUnix  int64
	LastPullUnix   int64
	LastBranchUnix int64

	ForkID   int64
	BaseRepo *Repository
	Created  time.Time
	Updated  time.Time
	Topics   []string
}

// Topic represents a repository topic.
type Topic struct {
	ID       int64
	Name     string
	NumRepos int
}

// RepositoriesStore defines the interface for managing repositories.
type RepositoriesStore interface {
	GetByID(ctx context.Context, id int64) (*Repository, error)
	GetByName(ctx context.Context, ownerID int64, name string) (*Repository, error)
	List(ctx context.Context, opts any) ([]*Repository, int64, error)
	Create(ctx context.Context, ownerID int64, opts any) (*Repository, error)
	Update(ctx context.Context, repoID int64, opts any) error
	Star(ctx context.Context, userID, repoID int64) error
	Watch(ctx context.Context, userID, repoID int64) error
	Touch(ctx context.Context, id int64) error
	ListForks(ctx context.Context, repoID int64, page, pageSize int) ([]*Repository, error)
	ListStarred(ctx context.Context, userID int64, page, pageSize int) ([]*Repository, error)
	ListStars(ctx context.Context, repoID int64, page, pageSize int) ([]*User, error)
	GetTopicsByRepoID(ctx context.Context, repoID int64) ([]string, error)
	AddTopics(ctx context.Context, repoID int64, topics []string) error
	HasForkedBy(ctx context.Context, repoID, userID int64) bool
	GetByCollaboratorID(ctx context.Context, collaboratorID int64, limit int, orderBy string) ([]*Repository, error)
	ListPopularTopics(ctx context.Context, limit int) ([]*Topic, error)
	ListWatches(ctx context.Context, repoID int64) ([]*User, error)
}
