package database

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/cryptox"
	"github.com/nef734/ggit/internal/domain"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"
)

// TwoFactor is a 2FA token of a user.
type TwoFactor struct {
	ID          int64 `gorm:"primaryKey"`
	UserID      int64 `gorm:"unique"`
	Secret      string
	Created     time.Time `gorm:"-" json:"-"`
	CreatedUnix int64
}

func (t *TwoFactor) ToDomain() *domain.TwoFactor {
	if t == nil {
		return nil
	}
	return &domain.TwoFactor{
		ID:      t.ID,
		UserID:  t.UserID,
		Secret:  t.Secret,
		Created: t.Created,
	}
}

func FromDomainTwoFactor(t *domain.TwoFactor) *TwoFactor {
	if t == nil {
		return nil
	}
	return &TwoFactor{
		ID:          t.ID,
		UserID:      t.UserID,
		Secret:      t.Secret,
		CreatedUnix: t.Created.Unix(),
	}
}

// ValidateTOTP returns true if given passcode is valid for two-factor authentication token.
// It also returns possible validation error.
func (t *TwoFactor) ValidateTOTP(passcode string) (bool, error) {
	secret, err := base64.StdEncoding.DecodeString(t.Secret)
	if err != nil {
		return false, fmt.Errorf("DecodeString: %v", err)
	}
	decryptSecret, err := cryptox.AESGCMDecrypt(cryptox.MD5Bytes(conf.Security.SecretKey), secret)
	if err != nil {
		return false, fmt.Errorf("AESGCMDecrypt: %v", err)
	}
	return totp.Validate(passcode, string(decryptSecret)), nil
}

// DeleteTwoFactor removes two-factor authentication token and recovery codes of given user.
//
// Deprecated: Use TwoFactorsStore.Delete instead.
func DeleteTwoFactor(userID int64) (err error) {
	return Handle.TwoFactors().Delete(context.Background(), userID)
}

// TwoFactorRecoveryCode represents a two-factor authentication recovery code.
type TwoFactorRecoveryCode struct {
	ID     int64  `gorm:"primaryKey"`
	UserID int64  `gorm:"index"`
	Code   string `gorm:"type:VARCHAR(11)"`
	IsUsed bool
}

func (c *TwoFactorRecoveryCode) ToDomain() *domain.TwoFactorRecoveryCode {
	if c == nil {
		return nil
	}
	return &domain.TwoFactorRecoveryCode{
		ID:     c.ID,
		UserID: c.UserID,
		Code:   c.Code,
		IsUsed: c.IsUsed,
	}
}

func FromDomainTwoFactorRecoveryCode(c *domain.TwoFactorRecoveryCode) *TwoFactorRecoveryCode {
	if c == nil {
		return nil
	}
	return &TwoFactorRecoveryCode{
		ID:     c.ID,
		UserID: c.UserID,
		Code:   c.Code,
		IsUsed: c.IsUsed,
	}
}

// GetRecoveryCodesByUserID returns recovery codes of given user.
func GetRecoveryCodesByUserID(userID int64) ([]*TwoFactorRecoveryCode, error) {
	codes, err := Handle.TwoFactors().GetRecoveryCodesByUserID(context.Background(), userID)
	if err != nil {
		return nil, err
	}
	res := make([]*TwoFactorRecoveryCode, len(codes))
	for i, c := range codes {
		res[i] = FromDomainTwoFactorRecoveryCode(c)
	}
	return res, nil
}

func deleteRecoveryCodesByUserID(tx *gorm.DB, userID int64) error {
	return tx.Where("user_id = ?", userID).Delete(new(TwoFactorRecoveryCode)).Error
}

// RegenerateRecoveryCodes regenerates new set of recovery codes for given user.
//
// Deprecated: Use TwoFactorsStore.RegenerateRecoveryCodes instead.
func RegenerateRecoveryCodes(userID int64) error {
	return Handle.TwoFactors().RegenerateRecoveryCodes(context.Background(), userID)
}

type ErrTwoFactorRecoveryCodeNotFound struct {
	Code string
}

func IsTwoFactorRecoveryCodeNotFound(err error) bool {
	_, ok := err.(ErrTwoFactorRecoveryCodeNotFound)
	return ok
}

func (err ErrTwoFactorRecoveryCodeNotFound) Error() string {
	return fmt.Sprintf("two-factor recovery code does not found [code: %s]", err.Code)
}
