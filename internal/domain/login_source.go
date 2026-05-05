package domain

import (
	"context"
)

// LoginSource represents an external way for authorizing users.
type LoginSource struct {
	ID        int64
	Type      int
	Name      string
	IsActived bool
	IsDefault bool
	Config    string
}

// CreateLoginSourceOptions contains options for creating a login source.
type CreateLoginSourceOptions struct {
	Type      int
	Name      string
	Activated bool
	Default   bool
	Config    any
}

// UpdateLoginSourceOptions contains options for updating a login source.
type UpdateLoginSourceOptions struct {
	Name      *string
	Activated *bool
	Default   *bool
	Config    any
}

// ListLoginSourceOptions contains options for listing login sources.
type ListLoginSourceOptions struct {
	OnlyActivated bool
}

// LoginSourcesStore defines the interface for managing login sources.
type LoginSourcesStore interface {
	Create(ctx context.Context, opts CreateLoginSourceOptions) (*LoginSource, error)
	Update(ctx context.Context, id int64, opts UpdateLoginSourceOptions) error
	Count(ctx context.Context) int64
	DeleteByID(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*LoginSource, error)
	List(ctx context.Context, opts ListLoginSourceOptions) ([]*LoginSource, error)
	ResetNonDefault(ctx context.Context, dflt *LoginSource) error
	Save(ctx context.Context, source *LoginSource) error
}
