package domain

import (
	"context"
)

// UserBlocksStore defines the interface for managing user blocks.
type UserBlocksStore interface {
	// ListByUserID returns all users blocked by the given user.
	ListByUserID(ctx context.Context, userID int64) ([]*User, error)
	// Block blocks a user.
	Block(ctx context.Context, userID, blockID int64) error
	// Unblock unblocks a user.
	Unblock(ctx context.Context, userID, blockID int64) error
	// IsBlocked returns true if blockID is blocked by userID.
	IsBlocked(ctx context.Context, userID, blockID int64) bool
}
