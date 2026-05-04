package migrations

import (
	"gorm.io/gorm"
)

func v31(db *gorm.DB) error {
	type Milestone struct {
		HillState int `gorm:"default:0"`
	}
	return db.AutoMigrate(&Milestone{})
}
