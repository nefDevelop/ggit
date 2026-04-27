package database

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
	"gorm.io/gorm"
)

type PullReviewType int

const (
	PullReviewTypePending PullReviewType = iota
	PullReviewTypeApprove
	PullReviewTypeComment
	PullReviewTypeReject
)

// PullReview represents a review on a pull request.
type PullReview struct {
	ID         int64          `gorm:"primaryKey"`
	Type       PullReviewType `gorm:"not null"`
	ReviewerID int64          `gorm:"index"`
	Reviewer   *User          `json:"reviewer" gorm:"-"`
	IssueID    int64          `gorm:"index"` // Linked to the pull request's issue ID
	Content    string         `gorm:"type:TEXT"`
	
	Comments   []*PullReviewComment `json:"comments" gorm:"-"`

	CreatedUnix int64 `gorm:"not null"`
	UpdatedUnix int64 `gorm:"not null"`
}

func (r *PullReview) BeforeCreate(tx *gorm.DB) error {
	if r.CreatedUnix == 0 {
		r.CreatedUnix = time.Now().Unix()
		r.UpdatedUnix = r.CreatedUnix
	}
	return nil
}

func (r *PullReview) BeforeUpdate(tx *gorm.DB) error {
	r.UpdatedUnix = time.Now().Unix()
	return nil
}

// PullReviewsStore is the storage layer for pull reviews.
type PullReviewsStore struct {
	db *gorm.DB
}

func newPullReviewsStore(db *gorm.DB) *PullReviewsStore {
	return &PullReviewsStore{db: db}
}

// ListByIssueID returns all reviews for a given issue (pull request).
func (s *PullReviewsStore) ListByIssueID(ctx context.Context, issueID int64) ([]*PullReview, error) {
	reviews := make([]*PullReview, 0, 5)
	err := s.db.WithContext(ctx).Where("issue_id = ?", issueID).Order("created_unix ASC").Find(&reviews).Error
	if err != nil {
		return nil, errors.Wrap(err, "list reviews")
	}
	return reviews, nil
}

// Create creates a new review.
func (s *PullReviewsStore) Create(ctx context.Context, r *PullReview) error {
	return s.db.WithContext(ctx).Create(r).Error
}

type PullReviewStats struct {
	Approvals int
	Rejections int
}

// GetReviewStatsByIssueID returns approval and rejection counts for a pull request.
func (s *PullReviewsStore) GetReviewStatsByIssueID(ctx context.Context, issueID int64) (*PullReviewStats, error) {
	var stats PullReviewStats
	err := s.db.WithContext(ctx).Model(&PullReview{}).
		Select("COUNT(CASE WHEN type = ? THEN 1 END) as approvals, COUNT(CASE WHEN type = ? THEN 1 END) as rejections", PullReviewTypeApprove, PullReviewTypeReject).
		Where("issue_id = ?", issueID).
		Scan(&stats).Error
	if err != nil {
		return nil, errors.Wrap(err, "get review stats")
	}
	return &stats, nil
}
