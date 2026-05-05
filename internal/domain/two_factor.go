package domain

import (
	"context"
	"time"
)

// TwoFactor is a 2FA token of a user.
type TwoFactor struct {
	ID      int64
	UserID  int64
	Secret  string
	Created time.Time
}

// TwoFactorRecoveryCode represents a two-factor authentication recovery code.
type TwoFactorRecoveryCode struct {
	ID     int64
	UserID int64
	Code   string
	IsUsed bool
}

// TwoFactorsStore defines the interface for managing two-factor authentication.
type TwoFactorsStore interface {
	// Create creates a new 2FA token and recovery codes for given user.
	Create(ctx context.Context, userID int64, key, secret string) error
	// GetByUserID returns the 2FA token of given user.
	GetByUserID(ctx context.Context, userID int64) (*TwoFactor, error)
	// IsEnabled returns true if the user has enabled 2FA.
	IsEnabled(ctx context.Context, userID int64) bool
	// UseRecoveryCode validates a recovery code of given user and marks it as used if valid.
	UseRecoveryCode(ctx context.Context, userID int64, code string) error
	// GetRecoveryCodesByUserID returns recovery codes of given user.
	GetRecoveryCodesByUserID(ctx context.Context, userID int64) ([]*TwoFactorRecoveryCode, error)
	// RegenerateRecoveryCodes regenerates new set of recovery codes for given user.
	RegenerateRecoveryCodes(ctx context.Context, userID int64) error
	// Delete removes two-factor authentication token and recovery codes of given user.
	Delete(ctx context.Context, userID int64) error
}
