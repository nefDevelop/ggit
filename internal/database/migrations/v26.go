package migrations

import (
	"gorm.io/gorm"
)

func createLoginSessionTable(db *gorm.DB) error {
	type LoginSession struct {
		ID          int64  `gorm:"primaryKey"`
		UserID      int64  `gorm:"index;not null"`
		SID         string `gorm:"uniqueIndex;column:sid;not null"`
		UserAgent   string
		RemoteAddr  string
		CreatedUnix int64 `gorm:"not null"`
		UpdatedUnix int64 `gorm:"not null"`
	}
	return db.AutoMigrate(&LoginSession{})
}
