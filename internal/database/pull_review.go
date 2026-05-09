package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	apiv2types "github.com/nef734/ggit/internal/route/api/v2/types"
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

	Comments []*PullReviewComment `json:"comments" gorm:"-"`

	CreatedUnix int64 `gorm:"not null"`
	UpdatedUnix int64 `gorm:"not null"`
}

func (r *PullReview) APIFormat() *apiv2types.PullReviewResponse {
	apiComments := make([]*apiv2types.PullReviewCommentResponse, len(r.Comments))
	for i, c := range r.Comments {
		apiComments[i] = c.APIFormat()
	}

	res := &apiv2types.PullReviewResponse{
		ID:        r.ID,
		Type:      int(r.Type),
		Content:   r.Content,
		Comments:  apiComments,
		CreatedAt: time.Unix(r.CreatedUnix, 0),
	}
	if r.Reviewer != nil {
		res.Reviewer = r.Reviewer.APIFormat()
	}
	return res
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
		return nil, fmt.Errorf("list reviews"+": %w", err)
	}
	return reviews, nil
}

// Create creates a new review.
func (s *PullReviewsStore) Create(ctx context.Context, r *PullReview) error {
	if r.IssueID <= 0 {
		return errors.New("invalid issue ID")
	}
	if r.Type < PullReviewTypePending || r.Type > PullReviewTypeReject {
		return errors.New("invalid review type")
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(r).Error; err != nil {
			return err
		}

		pr, err := getPullRequestByIssueID(tx, r.IssueID)
		if err != nil {
			return err
		}
		return pr.RefreshAttentionState(ctx, tx)
	})
}

type PullReviewStats struct {
	Approvals  int64
	Rejections int64
}

// GetReviewStatsByIssueID returns approval and rejection counts for a pull request.
func (s *PullReviewsStore) GetReviewStatsByIssueID(ctx context.Context, issueID int64) (*PullReviewStats, error) {
	var stats PullReviewStats
	err := s.db.WithContext(ctx).Model(&PullReview{}).
		Select("COUNT(CASE WHEN type = ? THEN 1 END) as approvals, COUNT(CASE WHEN type = ? THEN 1 END) as rejections", PullReviewTypeApprove, PullReviewTypeReject).
		Where("issue_id = ?", issueID).
		Scan(&stats).Error
	if err != nil {
		return nil, fmt.Errorf("get review stats"+": %w", err)
	}
	return &stats, nil
}
