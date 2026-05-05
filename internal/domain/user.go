package domain

import (
	"context"
	"time"
)

// UserType indicates the type of the user account.
type UserType int

const (
	UserTypeIndividual UserType = iota
	UserTypeOrganization
	UserTypeOrganizationPrivate
)

// User represents the object of an individual or an organization.
type User struct {
	ID                 int64
	LowerName          string
	Name               string
	FullName           string
	Email              string
	Password           string
	LoginSource        int64
	LoginName          string
	Type               UserType
	Location           string
	Website            string
	Rands              string
	Salt               string
	Created            time.Time
	Updated            time.Time
	LastRepoVisibility bool
	MaxRepoCreation    int
	IsActive           bool
	IsAdmin            bool
	AllowGitHook       bool
	AllowImportLocal   bool
	ProhibitLogin      bool
	Avatar             string
	AvatarEmail        string
	UseCustomAvatar    bool
	NumStars           int
	NumRepos           int
	Description        string
	NumTeams           int
	NumMembers         int
	DisableNewsFeed    bool
}

// UsersStore defines the interface for managing users.
type UsersStore interface {
	Authenticate(ctx context.Context, login, password string, loginSourceID int64) (*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	List(ctx context.Context, page, pageSize int) ([]*User, error)
	SearchByName(ctx context.Context, keyword string, page, pageSize int, orderBy string) ([]*User, int64, error)
	Create(ctx context.Context, username, email string, opts any) (*User, error)
	Update(ctx context.Context, userID int64, opts any) error
	DeleteByID(ctx context.Context, userID int64, skipRewriteAuthorizedKeys bool) error
	UseCustomAvatar(ctx context.Context, userID int64, avatar []byte) error
	Count(ctx context.Context) int64
}
