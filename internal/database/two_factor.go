package database

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"

	"gogs.io/gogs/internal/conf"
	"gogs.io/gogs/internal/cryptox"
)

// TwoFactor is a 2FA token of a user.
type TwoFactor struct {
	ID          int64 `gorm:"primaryKey"`
	UserID      int64 `gorm:"unique"`
	Secret      string
	Created     time.Time `gorm:"-" json:"-"`
	CreatedUnix int64
}

// ValidateTOTP returns true if given passcode is valid for two-factor authentication token.
// It also returns possible validation error.
func (t *TwoFactor) ValidateTOTP(passcode string) (bool, error) {
	secret, err := base64.StdEncoding.DecodeString(t.Secret)
	if err != nil {
		return false, errors.Newf("DecodeString: %v", err)
	}
	decryptSecret, err := cryptox.AESGCMDecrypt(cryptox.MD5Bytes(conf.Security.SecretKey), secret)
	if err != nil {
		return false, errors.Newf("AESGCMDecrypt: %v", err)
	}
	return totp.Validate(passcode, string(decryptSecret)), nil
}

// DeleteTwoFactor removes two-factor authentication token and recovery codes of given user.
func DeleteTwoFactor(userID int64) (err error) {
	return Handle.db.Transaction(func(tx *gorm.DB) error {
		if err = tx.Where("user_id = ?", userID).Delete(new(TwoFactor)).Error; err != nil {
			return errors.Wrap(err, "delete two-factor")
		} else if err = deleteRecoveryCodesByUserID(tx, userID); err != nil {
			return errors.Wrap(err, "deleteRecoveryCodesByUserID")
		}
		return nil
	})
}

// TwoFactorRecoveryCode represents a two-factor authentication recovery code.
type TwoFactorRecoveryCode struct {
	ID     int64  `gorm:"primaryKey"`
	UserID int64  `gorm:"index"`
	Code   string `gorm:"type:VARCHAR(11)"`
	IsUsed bool
}

func (s *TwoFactorsStore) GetRecoveryCodesByUserID(ctx context.Context, userID int64) ([]*TwoFactorRecoveryCode, error) {
	recoveryCodes := make([]*TwoFactorRecoveryCode, 0, 10)
	return recoveryCodes, s.db.WithContext(ctx).Where("user_id = ?", userID).Find(&recoveryCodes).Error
}

// GetRecoveryCodesByUserID returns recovery codes of given user.
func GetRecoveryCodesByUserID(userID int64) ([]*TwoFactorRecoveryCode, error) {
	return Handle.TwoFactors().GetRecoveryCodesByUserID(context.Background(), userID)
}

func deleteRecoveryCodesByUserID(tx *gorm.DB, userID int64) error {
	return tx.Where("user_id = ?", userID).Delete(new(TwoFactorRecoveryCode)).Error
}

// RegenerateRecoveryCodes regenerates new set of recovery codes for given user.
func RegenerateRecoveryCodes(userID int64) error {
	recoveryCodes, err := generateRecoveryCodes(userID, 10)
	if err != nil {
		return errors.Wrap(err, "generateRecoveryCodes")
	}

	return Handle.db.Transaction(func(tx *gorm.DB) error {
		if err = deleteRecoveryCodesByUserID(tx, userID); err != nil {
			return errors.Wrap(err, "deleteRecoveryCodesByUserID")
		} else if err = tx.Create(recoveryCodes).Error; err != nil {
			return errors.Wrap(err, "insert new recovery codes")
		}
		return nil
	})
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
