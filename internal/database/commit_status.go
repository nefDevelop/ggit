package database

import (
	"context"

	"gorm.io/gorm"
)

type CommitStatusState string

const (
	CommitStatusPending CommitStatusState = "pending"
	CommitStatusSuccess CommitStatusState = "success"
	CommitStatusError   CommitStatusState = "error"
	CommitStatusFailure CommitStatusState = "failure"
	CommitStatusWarning CommitStatusState = "warning"
)

// CommitStatus represents a status of a commit from an external CI/CD service.
type CommitStatus struct {
	ID          int64
	RepoID      int64  `gorm:"index"`
	SHA         string `gorm:"type:VARCHAR(40);index"`
	State       CommitStatusState
	TargetURL   string `gorm:"type:TEXT"`
	Description string `gorm:"type:TEXT"`
	Context     string `gorm:"type:VARCHAR(255);index"` // e.g. "ci/jenkins", "security/snyk"
	CreatorID   int64

	CreatedAt int64 `gorm:"autoCreateTime"`
	UpdatedAt int64 `gorm:"autoUpdateTime"`
}

type CommitStatusStore interface {
	Create(ctx context.Context, status *CommitStatus) error
	GetLatestBySHA(ctx context.Context, repoID int64, sha string) ([]*CommitStatus, error)
	GetCombinedStatusBySHA(ctx context.Context, repoID int64, sha string) CommitStatusState
}

func (s *CommitStatuses) Create(ctx context.Context, status *CommitStatus) error {
	return s.db.WithContext(ctx).Create(status).Error
}

func (s *CommitStatuses) GetLatestBySHA(ctx context.Context, repoID int64, sha string) ([]*CommitStatus, error) {
	// Get only the latest status for each context
	var statuses []*CommitStatus
	subQuery := s.db.Select("MAX(id)").Table("commit_status").
		Where("repo_id = ? AND sha = ?", repoID, sha).
		Group("context")
	
	err := s.db.WithContext(ctx).Where("id IN (?)", subQuery).Find(&statuses).Error
	return statuses, err
}

func (s *CommitStatuses) GetCombinedStatusBySHA(ctx context.Context, repoID int64, sha string) CommitStatusState {
	statuses, err := s.GetLatestBySHA(ctx, repoID, sha)
	if err != nil || len(statuses) == 0 {
		return ""
	}

	states := make(map[CommitStatusState]int)
	for _, st := range statuses {
		states[st.State]++
	}

	if states[CommitStatusError] > 0 || states[CommitStatusFailure] > 0 {
		return CommitStatusFailure
	}
	if states[CommitStatusPending] > 0 {
		return CommitStatusPending
	}
	if states[CommitStatusWarning] > 0 {
		return CommitStatusWarning
	}
	return CommitStatusSuccess
}

type CommitStatuses struct {
	db *gorm.DB
}

func (h *DB) CommitStatuses() CommitStatusStore {
	return &CommitStatuses{db: h.db}
}
