package domain

import (
	"context"
)

// AccessMode is the access mode of a user has to a repository.
type AccessMode int

const (
	AccessModeNone  AccessMode = iota // 0
	AccessModeRead                    // 1
	AccessModeWrite                   // 2
	AccessModeAdmin                   // 3
	AccessModeOwner                   // 4
)

func (mode AccessMode) String() string {
	switch mode {
	case AccessModeRead:
		return "read"
	case AccessModeWrite:
		return "write"
	case AccessModeAdmin:
		return "admin"
	case AccessModeOwner:
		return "owner"
	default:
		return "none"
	}
}

// AccessModeOptions contains options for calculating access mode.
type AccessModeOptions struct {
	OwnerID int64 // The ID of the repository owner.
	Private bool  // Whether the repository is private.
}

// PermissionsStore defines the interface for repository permissions.
type PermissionsStore interface {
	// AccessMode returns the access mode of given user has to the repository.
	AccessMode(ctx context.Context, userID, repoID int64, opts AccessModeOptions) AccessMode
	// Authorize returns true if the user has as good as desired access mode to the repository.
	Authorize(ctx context.Context, userID, repoID int64, desired AccessMode, opts AccessModeOptions) bool
	// SetRepoPerms does a full update to which users have which level of access to given repository.
	SetRepoPerms(ctx context.Context, repoID int64, accessMap map[int64]AccessMode) error
}
