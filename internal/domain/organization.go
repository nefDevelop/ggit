package domain

import (
	"context"
)

// ListOrgsOptions contains options for listing organizations.
type ListOrgsOptions struct {
	// Filter by the membership with the given user ID.
	MemberID int64
	// Whether to include private memberships.
	IncludePrivateMembers bool
}

// OrganizationsStore defines the interface for managing organizations.
type OrganizationsStore interface {
	List(ctx context.Context, opts ListOrgsOptions) ([]*User, error)
	SearchByName(ctx context.Context, keyword string, page, pageSize int, orderBy string) ([]*User, int64, error)
	CountByUser(ctx context.Context, userID int64) (int64, error)
}
