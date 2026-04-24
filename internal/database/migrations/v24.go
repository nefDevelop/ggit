package migrations

import (
	"gorm.io/gorm"
)

func addIsArchivedToRepository(db *gorm.DB) error {
	type repository struct {
		IsArchived bool `gorm:"not null;default:FALSE"`
	}
	if db.Migrator().HasColumn(&repository{}, "is_archived") {
		return errMigrationSkipped
	}
	return db.Migrator().AddColumn(&repository{}, "is_archived")
}
