package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nef734/ggit/internal/cryptox"
	"github.com/nef734/ggit/internal/domain"
	"github.com/nef734/ggit/internal/errx"
	"gorm.io/gorm"
)

// AccessToken is a personal access token.
type AccessToken struct {
	ID     int64 `gorm:"primarykey"`
	UserID int64 `gorm:"column:uid;index"`
	Name   string
	Sha1   string `gorm:"type:VARCHAR(40);unique"`
	SHA256 string `gorm:"type:VARCHAR(64);unique;not null"`

	Created           time.Time `gorm:"-" json:"-"`
	CreatedUnix       int64
	Updated           time.Time `gorm:"-" json:"-"`
	UpdatedUnix       int64
	HasRecentActivity bool `gorm:"-" json:"-"`
	HasUsed           bool `gorm:"-" json:"-"`
}

// BeforeCreate implements the GORM create hook.
func (t *AccessToken) BeforeCreate(tx *gorm.DB) error {
	if t.CreatedUnix == 0 {
		t.CreatedUnix = tx.NowFunc().Unix()
	}
	return nil
}

// AfterFind implements the GORM query hook.
func (t *AccessToken) AfterFind(tx *gorm.DB) error {
	t.Created = time.Unix(t.CreatedUnix, 0).Local()
	if t.UpdatedUnix > 0 {
		t.Updated = time.Unix(t.UpdatedUnix, 0).Local()
		t.HasUsed = t.Updated.After(t.Created)
		t.HasRecentActivity = t.Updated.Add(7 * 24 * time.Hour).After(tx.NowFunc())
	}
	return nil
}

func (t *AccessToken) ToDomain() *domain.AccessToken {
	if t == nil {
		return nil
	}
	return &domain.AccessToken{
		ID:                t.ID,
		UserID:            t.UserID,
		Name:              t.Name,
		Sha1:              t.Sha1,
		SHA256:            t.SHA256,
		Created:           t.Created,
		Updated:           t.Updated,
		HasRecentActivity: t.HasRecentActivity,
		HasUsed:           t.HasUsed,
	}
}

func FromDomainAccessToken(t *domain.AccessToken) *AccessToken {
	if t == nil {
		return nil
	}
	return &AccessToken{
		ID:          t.ID,
		UserID:      t.UserID,
		Name:        t.Name,
		Sha1:        t.Sha1,
		SHA256:      t.SHA256,
		CreatedUnix: t.Created.Unix(),
		UpdatedUnix: t.Updated.Unix(),
	}
}

// AccessTokensStore is the storage layer for access tokens.
type AccessTokensStore struct {
	db *gorm.DB
}

func newAccessTokensStore(db *gorm.DB) *AccessTokensStore {
	return &AccessTokensStore{db: db}
}

type ErrAccessTokenAlreadyExist struct {
	args errx.Args
}

func IsErrAccessTokenAlreadyExist(err error) bool {
	var target ErrAccessTokenAlreadyExist
	return errors.As(err, &target)
}

func (err ErrAccessTokenAlreadyExist) Error() string {
	return fmt.Sprintf("access token already exists: %v", err.args)
}

// Create creates a new access token and persist to database. It returns
// ErrAccessTokenAlreadyExist when an access token with same name already exists
// for the user.
func (s *AccessTokensStore) Create(ctx context.Context, userID int64, name string) (*domain.AccessToken, string, error) {
	err := s.db.WithContext(ctx).Where("uid = ? AND name = ?", userID, name).First(new(AccessToken)).Error
	if err == nil {
		return nil, "", ErrAccessTokenAlreadyExist{args: errx.Args{"userID": userID, "name": name}}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, "", err
	}

	token := cryptox.SHA1(uuid.New().String())
	sha256 := cryptox.SHA256(token)

	accessToken := &AccessToken{
		UserID: userID,
		Name:   name,
		Sha1:   sha256[:40], // To pass the column unique constraint, keep the length of SHA1.
		SHA256: sha256,
	}
	if err = s.db.WithContext(ctx).Create(accessToken).Error; err != nil {
		return nil, "", err
	}

	// Set back the raw access token value, for the sake of the caller.
	accessToken.Sha1 = token
	return accessToken.ToDomain(), token, nil
}

// DeleteByID deletes the access token by given ID.
//
// 🚨 SECURITY: The "userID" is required to prevent attacker deletes arbitrary
// access token that belongs to another user.
func (s *AccessTokensStore) DeleteByID(ctx context.Context, userID, id int64) error {
	return s.db.WithContext(ctx).Where("id = ? AND uid = ?", id, userID).Delete(new(AccessToken)).Error
}

var _ errx.NotFound = (*ErrAccessTokenNotExist)(nil)

type ErrAccessTokenNotExist struct {
	args errx.Args
}

// IsErrAccessTokenNotExist returns true if the underlying error has the type
// ErrAccessTokenNotExist.
func IsErrAccessTokenNotExist(err error) bool {
	var target ErrAccessTokenNotExist
	return errors.As(err, &target)
}

func (err ErrAccessTokenNotExist) Error() string {
	return fmt.Sprintf("access token does not exist: %v", err.args)
}

func (ErrAccessTokenNotExist) NotFound() bool {
	return true
}

// GetBySHA1 returns the access token with given SHA1. It returns
// ErrAccessTokenNotExist when not found.
func (s *AccessTokensStore) GetBySHA1(ctx context.Context, sha1 string) (*domain.AccessToken, error) {
	// No need to waste a query for an empty SHA1.
	if sha1 == "" {
		return nil, ErrAccessTokenNotExist{args: errx.Args{"sha": sha1}}
	}

	sha256 := cryptox.SHA256(sha1)
	token := new(AccessToken)
	err := s.db.WithContext(ctx).Where("sha256 = ?", sha256).First(token).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrAccessTokenNotExist{args: errx.Args{"sha": sha1}}
	} else if err != nil {
		return nil, err
	}
	return token.ToDomain(), nil
}

// GetBySHA256 returns the access token with given SHA256. It returns
// ErrAccessTokenNotExist when not found.
func (s *AccessTokensStore) GetBySHA256(ctx context.Context, sha256 string) (*domain.AccessToken, error) {
	if sha256 == "" {
		return nil, ErrAccessTokenNotExist{args: errx.Args{"sha256": sha256}}
	}

	token := new(AccessToken)
	err := s.db.WithContext(ctx).Where("sha256 = ?", sha256).First(token).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrAccessTokenNotExist{args: errx.Args{"sha256": sha256}}
	} else if err != nil {
		return nil, err
	}
	return token.ToDomain(), nil
}

// List returns all access tokens belongs to given user.
func (s *AccessTokensStore) List(ctx context.Context, userID int64) ([]*domain.AccessToken, error) {
	var tokens []*AccessToken
	err := s.db.WithContext(ctx).Where("uid = ?", userID).Order("id ASC").Find(&tokens).Error
	if err != nil {
		return nil, err
	}

	res := make([]*domain.AccessToken, len(tokens))
	for i, t := range tokens {
		res[i] = t.ToDomain()
	}
	return res, nil
}

// TouchByID updates the updated time of the given access token to the current time.
func (s *AccessTokensStore) TouchByID(ctx context.Context, id int64) error {
	return s.db.WithContext(ctx).
		Model(new(AccessToken)).
		Where("id = ?", id).
		UpdateColumn("updated_unix", s.db.NowFunc().Unix()).
		Error
}
