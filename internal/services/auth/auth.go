package auth

import (
	"context"

	log "unknwon.dev/clog/v2"

	"gogs.io/gogs/internal/conf"
	"gogs.io/gogs/internal/database"
	"gogs.io/gogs/internal/email"
)

// AuthService define las operaciones de negocio para la gestión de usuarios y acceso.
type AuthService interface {
	Authenticate(ctx context.Context, opts LoginOptions) (*database.User, error)
	Register(ctx context.Context, opts SignupOptions) (*database.User, error)
}

type authService struct {
	db *database.DB
}

// Auth es el singleton global para el servicio de autenticación.
var Auth AuthService

// Init inicializa el servicio de autenticación.
func Init(db *database.DB) {
	Auth = &authService{
		db: db,
	}
}

func (s *authService) Authenticate(ctx context.Context, opts LoginOptions) (*database.User, error) {
	return s.db.Users().Authenticate(ctx, opts.Username, opts.Password, opts.LoginSource)
}

func (s *authService) Register(ctx context.Context, opts SignupOptions) (*database.User, error) {
	user, err := s.db.Users().Create(
		ctx,
		opts.Username,
		opts.Email,
		database.CreateUserOptions{
			Password:  opts.Password,
			Activated: !conf.Auth.RequireEmailConfirmation,
		},
	)
	if err != nil {
		return nil, err
	}

	// Auto-set admin for the only user.
	if s.db.Users().Count(ctx) == 1 {
		v := true
		err := s.db.Users().Update(
			ctx,
			user.ID,
			database.UpdateUserOptions{
				IsActivated: &v,
				IsAdmin:     &v,
			},
		)
		if err != nil {
			return nil, err
		}
	}

	// Send confirmation email.
	if conf.Auth.RequireEmailConfirmation && user.ID > 1 {
		if err := email.SendActivateAccountMail(database.NewMailerUser(user), "Activate Account"); err != nil {
			log.Error("Failed to send activate account mail: %v", err)
		}
	}

	return user, nil
}
