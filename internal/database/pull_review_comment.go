package database

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
	apiv2types "github.com/nef734/ggit/internal/route/api/v2/types"
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
	IsResolved bool `gorm:"not null;default:FALSE"`
	
	CreatedUnix int64 `gorm:"not null"`
	UpdatedUnix int64 `gorm:"not null"`
}

func (c *PullReviewComment) APIFormat() *apiv2types.PullReviewCommentResponse {
	return &apiv2types.PullReviewCommentResponse{
		ID:         c.ID,
		ReviewID:   c.ReviewID,
		TreePath:   c.TreePath,
		Line:       c.Line,
		CommitID:   c.CommitID,
		Content:    c.Content,
		IsResolved: c.IsResolved,
		CreatedAt:  time.Unix(c.CreatedUnix, 0),
		UpdatedAt:  time.Unix(c.UpdatedUnix, 0),
	}
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
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(c).Error; err != nil {
			return err
		}

		pr, err := getPullRequestByIssueID(tx, c.IssueID)
		if err != nil {
			return err
		}
		return pr.RefreshAttentionState(ctx, tx)
	})
}

// SetResolved sets the resolved status of a comment.
func (s *PullReviewCommentsStore) SetResolved(ctx context.Context, id int64, resolved bool) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		c := new(PullReviewComment)
		if err := tx.Where("id = ?", id).First(c).Error; err != nil {
			return err
		}
		c.IsResolved = resolved
		if err := tx.Model(c).Update("is_resolved", resolved).Error; err != nil {
			return err
		}

		pr, err := getPullRequestByIssueID(tx, c.IssueID)
		if err != nil {
			return err
		}
		return pr.RefreshAttentionState(ctx, tx)
	})
}
