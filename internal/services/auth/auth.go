package auth

import (
	gocontext "context"

	"github.com/gin-contrib/sessions"
	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/email"
	log "github.com/nef734/ggit/internal/log"
	"github.com/nef734/ggit/internal/strx"
)

// AuthService define las operaciones de negocio para la gestión de usuarios y acceso.
type AuthService interface {
	Authenticate(ctx gocontext.Context, opts LoginOptions) (*database.User, error)
	Register(ctx gocontext.Context, opts SignupOptions) (*database.User, error)
	CreateSession(ctx gocontext.Context, sess sessions.Session, u *database.User) error
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

func (s *authService) Authenticate(ctx gocontext.Context, opts LoginOptions) (*database.User, error) {
	u, err := s.db.Users().Authenticate(ctx, opts.Username, opts.Password, opts.LoginSource)
	if err != nil {
		return nil, err
	}
	return database.FromDomainUser(u), nil
}

func (s *authService) Register(ctx gocontext.Context, opts SignupOptions) (*database.User, error) {
	du, err := s.db.Users().Create(
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
	user := database.FromDomainUser(du)

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

func (s *authService) CreateSession(ctx gocontext.Context, sess sessions.Session, u *database.User) error {
	sess.Set("uid", u.ID)
	sess.Set("uname", u.Name)

	sid, _ := strx.RandomChars(32)
	sess.Set("sid", sid)
	_, err := s.db.LoginSessions().Create(ctx, u.ID, sid, "", "") // UA and IP are optional here or handled by caller if needed
	if err != nil {
		log.Error("Failed to create login session: %v", err)
	}

	return sess.Save()
}
