package migrations

import (
	"gorm.io/gorm"
)

func addBranchProtectionAndReviews(db *gorm.DB) error {
	type ProtectBranch struct {
		RequiredApprovals     int
		EnableMergeWhitelist  bool
		MergeWhitelistUserIDs string `gorm:"type:TEXT"`
		MergeWhitelistTeamIDs string `gorm:"type:TEXT"`
		BlockOnOutdatedBranch bool
	}
	if err := db.AutoMigrate(&ProtectBranch{}); err != nil {
		return err
	}

	type PullReview struct {
		ID          int64  `gorm:"primaryKey"`
		Type        int    `gorm:"not null"`
		ReviewerID  int64  `gorm:"index"`
		IssueID     int64  `gorm:"index"`
		Content     string `gorm:"type:TEXT"`
		CreatedUnix int64  `gorm:"not null"`
		UpdatedUnix int64  `gorm:"not null"`
	}
	if err := db.AutoMigrate(&PullReview{}); err != nil {
		return err
	}

	type PullReviewComment struct {
		ID          int64  `gorm:"primaryKey"`
		ReviewID    int64  `gorm:"index"`
		IssueID     int64  `gorm:"index"`
		TreePath    string `gorm:"type:TEXT"`
		Line        int64
		CommitID    string `gorm:"type:VARCHAR(40)"`
		Content     string `gorm:"type:TEXT"`
		CreatedUnix int64  `gorm:"not null"`
		UpdatedUnix int64  `gorm:"not null"`
	}
	return db.AutoMigrate(&PullReviewComment{})
}
