package migrations

import (
	"gorm.io/gorm"
)

func createUserBlockTable(db *gorm.DB) error {
	type UserBlock struct {
		ID          int64 `gorm:"primaryKey"`
		UserID      int64 `gorm:"uniqueIndex:user_block_unique;not null"`
		BlockID     int64 `gorm:"uniqueIndex:user_block_unique;not null"`
		CreatedUnix int64 `gorm:"not null"`
	}
	return db.AutoMigrate(&UserBlock{})
}
