package database

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
	"gorm.io/gorm"
)

// PullReviewComment represents a comment on a specific line of a pull request diff.
type PullReviewComment struct {
	ID       int64  `gorm:"primaryKey"`
	ReviewID int64  `gorm:"index"` // Linked to PullReview
	IssueID  int64  `gorm:"index"` // Linked to PullRequest Issue
	
	TreePath string `gorm:"type:TEXT"`
	Line     int64  
	CommitID string `gorm:"type:VARCHAR(40)"` // SHA of the commit
	Content  string `gorm:"type:TEXT"`
	
	CreatedUnix int64 `gorm:"not null"`
	UpdatedUnix int64 `gorm:"not null"`
}

func (c *PullReviewComment) BeforeCreate(tx *gorm.DB) error {
	if c.CreatedUnix == 0 {
		c.CreatedUnix = time.Now().Unix()
		c.UpdatedUnix = c.CreatedUnix
	}
	return nil
}

func (c *PullReviewComment) BeforeUpdate(tx *gorm.DB) error {
	c.UpdatedUnix = time.Now().Unix()
	return nil
}

// PullReviewCommentsStore is the storage layer for pull review comments.
type PullReviewCommentsStore struct {
	db *gorm.DB
}

func newPullReviewCommentsStore(db *gorm.DB) *PullReviewCommentsStore {
	return &PullReviewCommentsStore{db: db}
}

// ListByIssueID returns all inline comments for a given pull request.
func (s *PullReviewCommentsStore) ListByIssueID(ctx context.Context, issueID int64) ([]*PullReviewComment, error) {
	comments := make([]*PullReviewComment, 0, 10)
	err := s.db.WithContext(ctx).Where("issue_id = ?", issueID).Order("created_unix ASC").Find(&comments).Error
	if err != nil {
		return nil, errors.Wrap(err, "list review comments")
	}
	return comments, nil
}

// ListByReviewID returns all inline comments for a given review.
func (s *PullReviewCommentsStore) ListByReviewID(ctx context.Context, reviewID int64) ([]*PullReviewComment, error) {
	comments := make([]*PullReviewComment, 0, 5)
	err := s.db.WithContext(ctx).Where("review_id = ?", reviewID).Order("created_unix ASC").Find(&comments).Error
	if err != nil {
		return nil, errors.Wrap(err, "list review comments by review ID")
	}
	return comments, nil
}

// Create creates a new review comment.
func (s *PullReviewCommentsStore) Create(ctx context.Context, c *PullReviewComment) error {
	return s.db.WithContext(ctx).Create(c).Error
}
