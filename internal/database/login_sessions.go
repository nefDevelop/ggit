package database

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// LoginSession represents a user's active login session.
type LoginSession struct {
	ID          int64  `gorm:"primaryKey"`
	UserID      int64  `gorm:"index;not null"`
	SID         string `gorm:"uniqueIndex;column:sid;not null"`
	UserAgent   string
	RemoteAddr  string
	CreatedUnix int64 `gorm:"not null"`
	UpdatedUnix int64 `gorm:"not null"`
}

// BeforeCreate implements the GORM BeforeCreate hook.
func (s *LoginSession) BeforeCreate(tx *gorm.DB) error {
	if s.CreatedUnix == 0 {
		s.CreatedUnix = time.Now().Unix()
		s.UpdatedUnix = s.CreatedUnix
	}
	return nil
}

// LoginSessionsStore is the storage layer for login sessions.
type LoginSessionsStore struct {
	db *gorm.DB
}

func newLoginSessionsStore(db *gorm.DB) *LoginSessionsStore {
	return &LoginSessionsStore{db: db}
}

// Create creates a new login session.
func (s *LoginSessionsStore) Create(ctx context.Context, userID int64, sid, ua, ip string) (*LoginSession, error) {
	sess := &LoginSession{
		UserID:     userID,
		SID:        sid,
		UserAgent:  ua,
		RemoteAddr: ip,
	}
	return sess, s.db.WithContext(ctx).Create(sess).Error
}

// GetBySID returns the login session with the given SID.
func (s *LoginSessionsStore) GetBySID(ctx context.Context, sid string) (*LoginSession, error) {
	sess := new(LoginSession)
	err := s.db.WithContext(ctx).Where("sid = ?", sid).First(sess).Error
	if err != nil {
		return nil, err
	}
	return sess, nil
}

// ListByUserID returns all active login sessions for the given user.
func (s *LoginSessionsStore) ListByUserID(ctx context.Context, userID int64) ([]*LoginSession, error) {
	var sessions []*LoginSession
	return sessions, s.db.WithContext(ctx).Where("user_id = ?", userID).Order("updated_unix DESC").Find(&sessions).Error
}

// DeleteBySID deletes the login session with the given SID.
func (s *LoginSessionsStore) DeleteBySID(ctx context.Context, sid string) error {
	return s.db.WithContext(ctx).Where("sid = ?", sid).Delete(&LoginSession{}).Error
}

// DeleteByIDAndUserID deletes the login session with the given ID and UserID.
func (s *LoginSessionsStore) DeleteByIDAndUserID(ctx context.Context, id, userID int64) error {
	return s.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&LoginSession{}).Error
}

// DeleteOthersByUserID deletes all login sessions for the given user except the one with the given SID.
func (s *LoginSessionsStore) DeleteOthersByUserID(ctx context.Context, userID int64, sid string) error {
	return s.db.WithContext(ctx).Where("user_id = ? AND sid <> ?", userID, sid).Delete(&LoginSession{}).Error
}

// TouchBySID updates the updated time of the login session with the given SID.
func (s *LoginSessionsStore) TouchBySID(ctx context.Context, sid string) error {
	return s.db.WithContext(ctx).Model(&LoginSession{}).Where("sid = ?", sid).Update("updated_unix", time.Now().Unix()).Error
}
