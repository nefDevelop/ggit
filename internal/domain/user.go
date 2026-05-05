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

// EmailAddress represents an email address of a user.
type EmailAddress struct {
	ID          int64
	UserID      int64
	Email       string
	IsActivated bool
	IsPrimary   bool
}

// IsOrganization returns true if the user is an organization.
func (u *User) IsOrganization() bool {
	return u.Type == UserTypeOrganization || u.Type == UserTypeOrganizationPrivate
}

// AvatarURLPath returns the URL path of the user's avatar.
func (u *User) AvatarURLPath() string {
	if u.UseCustomAvatar {
		return "/avatars/" + u.Name
	}
	return "/avatars/default"
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
	IsUsernameUsed(ctx context.Context, username string, excludeUserID int64) bool
	GetMailableEmailsByUsernames(ctx context.Context, usernames []string) ([]string, error)
	GetEmail(ctx context.Context, userID int64, email string, needsActivated bool) (*EmailAddress, error)
	DeleteInactivated(ctx context.Context) error
	DeleteCustomAvatar(ctx context.Context, userID int64) error
	ChangeUsername(ctx context.Context, userID int64, newUsername string) error
	AddEmail(ctx context.Context, userID int64, email string, isActivated bool) error
	ListEmails(ctx context.Context, userID int64) ([]*EmailAddress, error)
	MarkEmailPrimary(ctx context.Context, userID int64, email string) error
	DeleteEmail(ctx context.Context, userID int64, email string) error
	MarkEmailActivated(ctx context.Context, userID int64, email string) error
	GetByKeyID(ctx context.Context, keyID int64) (*User, error)
}
