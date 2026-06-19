package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nef734/ggit/internal/domain"
	"gorm.io/gorm"
)

// UserBlock represents a user blocking another user.
type UserBlock struct {
	ID          int64 `gorm:"primaryKey"`
	UserID      int64 `gorm:"uniqueIndex:user_block_unique;not null"`
	BlockID     int64 `gorm:"uniqueIndex:user_block_unique;not null"`
	CreatedUnix int64 `gorm:"not null"`
}

func (b *UserBlock) BeforeCreate(tx *gorm.DB) error {
	if b.CreatedUnix == 0 {
		b.CreatedUnix = time.Now().Unix()
	}
	return nil
}

// UserBlocksStore is the storage layer for user blocks.
type UserBlocksStore struct {
	db *gorm.DB
}

func newUserBlocksStore(db *gorm.DB) *UserBlocksStore {
	return &UserBlocksStore{db: db}
}

// ListByUserID returns all users blocked by the given user.
func (s *UserBlocksStore) ListByUserID(ctx context.Context, userID int64) ([]*domain.User, error) {
	users := make([]*User, 0, 5)
	err := s.db.WithContext(ctx).
		Table("user").
		Joins("JOIN user_block ON user_block.block_id = user.id").
		Where("user_block.user_id = ?", userID).
		Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("list blocked users"+": %v", err)
	}

	res := make([]*domain.User, len(users))
	for i, user := range users {
		res[i] = user.ToDomain()
	}
	return res, nil
}

// Block blocks a user.
func (s *UserBlocksStore) Block(ctx context.Context, userID, blockID int64) error {
	if userID == blockID {
		return errors.New("cannot block yourself")
	}
	return s.db.WithContext(ctx).FirstOrCreate(&UserBlock{UserID: userID, BlockID: blockID}).Error
}

// Unblock unblocks a user.
func (s *UserBlocksStore) Unblock(ctx context.Context, userID, blockID int64) error {
	return s.db.WithContext(ctx).Where("user_id = ? AND block_id = ?", userID, blockID).Delete(&UserBlock{}).Error
}

// IsBlocked returns true if blockID is blocked by userID.
func (s *UserBlocksStore) IsBlocked(ctx context.Context, userID, blockID int64) bool {
	var count int64
	s.db.WithContext(ctx).Model(&UserBlock{}).Where("user_id = ? AND block_id = ?", userID, blockID).Count(&count)
	return count > 0
}
