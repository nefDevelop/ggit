package repo

import (
	stdctx "context"

	"github.com/nef734/ggit/internal/database"
)

// Store is the data layer carrier for context middleware. This interface is
// meant to abstract away and limit the exposure of the underlying data layer to
// the handler through a thin-wrapper.
type Store interface {
	// GetAccessTokenBySHA1 returns the access token with given SHA1. It returns
	// database.ErrAccessTokenNotExist when not found.
	GetAccessTokenBySHA1(ctx stdctx.Context, sha1 string) (*database.AccessToken, error)
	// TouchAccessTokenByID updates the updated time of the given access token to
	// the current time.
	TouchAccessTokenByID(ctx stdctx.Context, id int64) error

	// GetRepositoryByName returns the repository with given owner and name. It
	// returns database.ErrRepoNotExist when not found.
	GetRepositoryByName(ctx stdctx.Context, ownerID int64, name string) (*database.Repository, error)

	// IsTwoFactorEnabled returns true if the user has enabled 2FA.
	IsTwoFactorEnabled(ctx stdctx.Context, userID int64) bool

	// GetUserByID returns the user with given ID. It returns
	// database.ErrUserNotExist when not found.
	GetUserByID(ctx stdctx.Context, id int64) (*database.User, error)
	// GetUserByUsername returns the user with given username. It returns
	// database.ErrUserNotExist when not found.
	GetUserByUsername(ctx stdctx.Context, username string) (*database.User, error)
	// CreateUser creates a new user and persists to database. It returns
	// database.ErrNameNotAllowed if the given name or pattern of the name is not
	// allowed as a username, or database.ErrUserAlreadyExist when a user with same
	// name already exists, or database.ErrEmailAlreadyUsed if the email has been
	// verified by another user.
	CreateUser(ctx stdctx.Context, username, email string, opts database.CreateUserOptions) (*database.User, error)
	// AuthenticateUser validates username and password via given login source ID.
	// It returns database.ErrUserNotExist when the user was not found.
	//
	// When the "loginSourceID" is negative, it aborts the process and returns
	// database.ErrUserNotExist if the user was not found in the database.
	//
	// When the "loginSourceID" is non-negative, it returns
	// database.ErrLoginSourceMismatch if the user has different login source ID
	// than the "loginSourceID".
	//
	// When the "loginSourceID" is positive, it tries to authenticate via given
	// login source and creates a new user when not yet exists in the database.
	AuthenticateUser(ctx stdctx.Context, login, password string, loginSourceID int64) (*database.User, error)
}

type store struct{}

// NewStore returns a new Store using the global database handle.
func NewStore() Store {
	return &store{}
}

func (*store) GetAccessTokenBySHA1(ctx stdctx.Context, sha1 string) (*database.AccessToken, error) {
	t, err := database.Handle.AccessTokens().GetBySHA1(ctx, sha1)
	if err != nil {
		return nil, err
	}
	return database.FromDomainAccessToken(t), nil
}

func (*store) TouchAccessTokenByID(ctx stdctx.Context, id int64) error {
	return database.Handle.AccessTokens().TouchByID(ctx, id)
}

func (*store) GetRepositoryByName(ctx stdctx.Context, ownerID int64, name string) (*database.Repository, error) {
	r, err := database.Handle.Repositories().GetByName(ctx, ownerID, name)
	if err != nil {
		return nil, err
	}
	return database.FromDomainRepository(r), nil
}

func (*store) IsTwoFactorEnabled(ctx stdctx.Context, userID int64) bool {
	return database.Handle.TwoFactors().IsEnabled(ctx, userID)
}

func (*store) GetUserByID(ctx stdctx.Context, id int64) (*database.User, error) {
	u, err := database.Handle.Users().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return database.FromDomainUser(u), nil
}

func (*store) GetUserByUsername(ctx stdctx.Context, username string) (*database.User, error) {
	u, err := database.Handle.Users().GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return database.FromDomainUser(u), nil
}

func (*store) CreateUser(ctx stdctx.Context, username, email string, opts database.CreateUserOptions) (*database.User, error) {
	u, err := database.Handle.Users().Create(ctx, username, email, opts)
	if err != nil {
		return nil, err
	}
	return database.FromDomainUser(u), nil
}

func (*store) AuthenticateUser(ctx stdctx.Context, login, password string, loginSourceID int64) (*database.User, error) {
	u, err := database.Handle.Users().Authenticate(ctx, login, password, loginSourceID)
	if err != nil {
		return nil, err
	}
	return database.FromDomainUser(u), nil
}
