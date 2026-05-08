package database

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/cockroachdb/errors"
	"gorm.io/gorm"

	"github.com/nef734/ggit/internal/auth"
	"github.com/nef734/ggit/internal/auth/ldap"
	"github.com/nef734/ggit/internal/auth/oidc"
	"github.com/nef734/ggit/internal/domain"
	"github.com/nef734/ggit/internal/errx"
)

// LoginSource represents an external way for authorizing users.
type LoginSource struct {
	ID        int64 `gorm:"primaryKey"`
	Type      auth.Type
	Name      string        `gorm:"unique"`
	IsActived bool          `gorm:"not null"`
	IsDefault bool          ``
	Provider  auth.Provider `gorm:"-"`
	Config    string        `gorm:"column:cfg;type:TEXT" json:"RawConfig"`

	Created     time.Time `gorm:"-" json:"-"`
	CreatedUnix int64
	Updated     time.Time `gorm:"-" json:"-"`
	UpdatedUnix int64

	File loginSourceFileStore `gorm:"-" json:"-"`
}

// BeforeSave implements the GORM save hook.
func (s *LoginSource) BeforeSave(_ *gorm.DB) (err error) {
	if s.Provider == nil {
		return nil
	}
	data, err := json.Marshal(s.Provider.Config())
	s.Config = string(data)
	return err
}

// BeforeCreate implements the GORM create hook.
func (s *LoginSource) BeforeCreate(tx *gorm.DB) error {
	if s.CreatedUnix == 0 {
		s.CreatedUnix = tx.NowFunc().Unix()
		s.UpdatedUnix = s.CreatedUnix
	}
	return nil
}

// BeforeUpdate implements the GORM update hook.
func (s *LoginSource) BeforeUpdate(tx *gorm.DB) error {
	s.UpdatedUnix = tx.NowFunc().Unix()
	return nil
}

type mockProviderConfig struct {
	ExternalAccount *auth.ExternalAccount
}

// AfterFind implements the GORM query hook.
func (s *LoginSource) AfterFind(_ *gorm.DB) error {
	s.Created = time.Unix(s.CreatedUnix, 0).Local()
	s.Updated = time.Unix(s.UpdatedUnix, 0).Local()

	switch s.Type {
	case auth.LDAP:
		var cfg ldap.Config
		err := json.Unmarshal([]byte(s.Config), &cfg)
		if err != nil {
			return err
		}
		s.Provider = ldap.NewProvider(false, &cfg)

	case auth.DLDAP:
		var cfg ldap.Config
		err := json.Unmarshal([]byte(s.Config), &cfg)
		if err != nil {
			return err
		}
		s.Provider = ldap.NewProvider(true, &cfg)

	case auth.Mock:
		var cfg mockProviderConfig
		err := json.Unmarshal([]byte(s.Config), &cfg)
		if err != nil {
			return err
		}
		mockProvider := NewMockProvider()
		mockProvider.AuthenticateFunc.SetDefaultReturn(cfg.ExternalAccount, nil)
		s.Provider = mockProvider

	case auth.OIDC:
		var cfg oidc.Config
		err := json.Unmarshal([]byte(s.Config), &cfg)
		if err != nil {
			return err
		}
		s.Provider = oidc.NewProvider(&cfg)

	default:
		return errors.Newf("unrecognized login source type: %v", s.Type)
	}
	return nil
}

func (s *LoginSource) TypeName() string {
	return auth.Name(s.Type)
}

func (s *LoginSource) IsLDAP() bool {
	return s.Type == auth.LDAP
}

func (s *LoginSource) IsDLDAP() bool {
	return s.Type == auth.DLDAP
}

func (s *LoginSource) IsOIDC() bool {
	return s.Type == auth.OIDC
}

func (s *LoginSource) LDAP() *ldap.Config {
	return s.Provider.Config().(*ldap.Config)
}

func (s *LoginSource) OIDC() *oidc.Config {
	return s.Provider.Config().(*oidc.Config)
}

func (s *LoginSource) ToDomain() *domain.LoginSource {
	if s == nil {
		return nil
	}
	return &domain.LoginSource{
		ID:        s.ID,
		Type:      int(s.Type),
		Name:      s.Name,
		IsActived: s.IsActived,
		IsDefault: s.IsDefault,
		Config:    s.Config,
		Created:   time.Unix(s.CreatedUnix, 0),
		Updated:   time.Unix(s.UpdatedUnix, 0),
	}
}

func FromDomainLoginSource(s *domain.LoginSource) *LoginSource {
	if s == nil {
		return nil
	}
	return &LoginSource{
		ID:          s.ID,
		Type:        auth.Type(s.Type),
		Name:        s.Name,
		IsActived:   s.IsActived,
		IsDefault:   s.IsDefault,
		Config:      s.Config,
		CreatedUnix: s.Created.Unix(),
		UpdatedUnix: s.Updated.Unix(),
	}
}

// LoginSourcesStore is the storage layer for login sources.
type LoginSourcesStore struct {
	db    *gorm.DB
	files loginSourceFilesStore
}

func newLoginSourcesStore(db *gorm.DB, files loginSourceFilesStore) *LoginSourcesStore {
	return &LoginSourcesStore{
		db:    db,
		files: files,
	}
}

type ErrLoginSourceAlreadyExist struct {
	args errx.Args
}

func IsErrLoginSourceAlreadyExist(err error) bool {
	return errors.As(err, &ErrLoginSourceAlreadyExist{})
}

func (err ErrLoginSourceAlreadyExist) Error() string {
	return fmt.Sprintf("login source already exists: %v", err.args)
}

// Create creates a new login source and persists it to the database. It returns
// ErrLoginSourceAlreadyExist when a login source with same name already exists.
func (s *LoginSourcesStore) Create(ctx context.Context, opts domain.CreateLoginSourceOptions) (*domain.LoginSource, error) {
	err := s.db.WithContext(ctx).Where("name = ?", opts.Name).First(new(LoginSource)).Error
	if err == nil {
		return nil, ErrLoginSourceAlreadyExist{args: errx.Args{"name": opts.Name}}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	source := &LoginSource{
		Type:      auth.Type(opts.Type),
		Name:      opts.Name,
		IsActived: opts.Activated,
		IsDefault: opts.Default,
	}
	data, err := json.Marshal(opts.Config)
	source.Config = string(data)
	if err != nil {
		return nil, err
	}
	err = s.db.WithContext(ctx).Create(source).Error
	if err != nil {
		return nil, err
	}
	return source.ToDomain(), nil
}

// Update updates fields for the given login source.
func (s *LoginSourcesStore) Update(ctx context.Context, id int64, opts domain.UpdateLoginSourceOptions) error {
	updates := map[string]any{
		"updated_unix": s.db.NowFunc().Unix(),
	}

	if opts.Name != nil {
		updates["name"] = *opts.Name
	}
	if opts.Activated != nil {
		updates["is_actived"] = *opts.Activated
	}
	if opts.Default != nil {
		updates["is_default"] = *opts.Default
	}
	if opts.Config != nil {
		data, err := json.Marshal(opts.Config)
		if err != nil {
			return errors.Wrap(err, "marshal config")
		}
		updates["cfg"] = string(data)
	}

	err := s.db.WithContext(ctx).Model(new(LoginSource)).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return err
	}

	if opts.Default != nil && *opts.Default {
		return s.ResetNonDefault(ctx, &domain.LoginSource{ID: id})
	}
	return nil
}

// Count returns the total number of login sources.
func (s *LoginSourcesStore) Count(ctx context.Context) int64 {
	var count int64
	s.db.WithContext(ctx).Model(new(LoginSource)).Count(&count)
	return count + int64(s.files.Len())
}

type ErrLoginSourceInUse struct {
	args errx.Args
}

func IsErrLoginSourceInUse(err error) bool {
	return errors.As(err, &ErrLoginSourceInUse{})
}

func (err ErrLoginSourceInUse) Error() string {
	return fmt.Sprintf("login source is still used by some users: %v", err.args)
}

// DeleteByID deletes a login source by given ID. It returns ErrLoginSourceInUse
// if at least one user is associated with the login source.
func (s *LoginSourcesStore) DeleteByID(ctx context.Context, id int64) error {
	var count int64
	err := s.db.WithContext(ctx).Model(new(User)).Where("login_source = ?", id).Count(&count).Error
	if err != nil {
		return err
	} else if count > 0 {
		return ErrLoginSourceInUse{args: errx.Args{"id": id}}
	}

	return s.db.WithContext(ctx).Where("id = ?", id).Delete(new(LoginSource)).Error
}

func (s *LoginSourcesStore) getByID(ctx context.Context, id int64) (*LoginSource, error) {
	source := new(LoginSource)
	err := s.db.WithContext(ctx).Where("id = ?", id).First(source).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.files.GetByID(id)
		}
		return nil, err
	}
	return source, nil
}

// GetByID returns the login source with given ID. It returns
// ErrLoginSourceNotExist when not found.
func (s *LoginSourcesStore) GetByID(ctx context.Context, id int64) (*domain.LoginSource, error) {
	source, err := s.getByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return source.ToDomain(), nil
}

// List returns a list of login sources filtered by options.
func (s *LoginSourcesStore) List(ctx context.Context, opts domain.ListLoginSourceOptions) ([]*domain.LoginSource, error) {
	var sources []*LoginSource
	query := s.db.WithContext(ctx).Order("id ASC")
	if opts.OnlyActivated {
		query = query.Where("is_actived = ?", true)
	}
	err := query.Find(&sources).Error
	if err != nil {
		return nil, err
	}

	res := make([]*domain.LoginSource, 0, len(sources)+s.files.Len())
	for _, src := range sources {
		res = append(res, src.ToDomain())
	}
	for _, src := range s.files.List(opts) {
		res = append(res, src.ToDomain())
	}
	return res, nil
}

// ResetNonDefault clears default flag for all the other login sources.
func (s *LoginSourcesStore) ResetNonDefault(ctx context.Context, dflt *domain.LoginSource) error {
	err := s.db.WithContext(ctx).
		Model(new(LoginSource)).
		Where("id != ?", dflt.ID).
		Updates(map[string]any{"is_default": false}).
		Error
	if err != nil {
		return err
	}

	for _, source := range s.files.List(domain.ListLoginSourceOptions{}) {
		if source.File != nil && source.ID != dflt.ID {
			source.File.SetGeneral("is_default", "false")
			if err = source.File.Save(); err != nil {
				return errors.Wrap(err, "save file")
			}
		}
	}

	s.files.Update(FromDomainLoginSource(dflt))
	return nil
}

// Save persists all values of given login source to database or local file. The
// Updated field is set to current time automatically.
func (s *LoginSourcesStore) Save(ctx context.Context, source *domain.LoginSource) error {
	m, err := s.files.GetByID(source.ID)
	if err == nil && m != nil {
		// It's a file-backed source
		m.Name = source.Name
		m.IsActived = source.IsActived
		m.IsDefault = source.IsDefault
		m.Config = source.Config

		m.File.SetGeneral("name", m.Name)
		m.File.SetGeneral("is_activated", strconv.FormatBool(m.IsActived))
		m.File.SetGeneral("is_default", strconv.FormatBool(m.IsDefault))

		// If it's a file-backed source, AfterFind was already called during loading,
		// and Provider should be set.
		if m.Provider != nil {
			if err := m.File.SetConfig(m.Provider.Config()); err != nil {
				return errors.Wrap(err, "set config")
			}
		}
		if err := m.File.Save(); err != nil {
			return errors.Wrap(err, "save file")
		}
		return nil
	} else if err != nil && !IsErrLoginSourceNotExist(err) {
		return err
	}

	m = FromDomainLoginSource(source)
	return s.db.WithContext(ctx).Save(m).Error
}
