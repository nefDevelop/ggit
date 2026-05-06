package migrations

import (
	"gorm.io/gorm"
)

func v33(db *gorm.DB) error {
	type Attachment struct {
		Size int64 `gorm:"not null;default:0"`
	}
	return db.AutoMigrate(&Attachment{})
}
