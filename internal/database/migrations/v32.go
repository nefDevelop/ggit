package migrations

import (
	"gorm.io/gorm"
)

func v32(db *gorm.DB) error {
	type Repository struct {
		ADRFolderName string `gorm:"type:VARCHAR(255);not null;default:'ADR'"`
	}
	return db.AutoMigrate(&Repository{})
}
