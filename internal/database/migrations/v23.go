package migrations

import (
	"gorm.io/gorm"
)

func addDisableNewsFeedToUser(db *gorm.DB) error {
	type user struct {
		DisableNewsFeed bool
	}
	if db.Migrator().HasColumn(&user{}, "DisableNewsFeed") {
		return errMigrationSkipped
	}
	return db.Migrator().AddColumn(&user{}, "DisableNewsFeed")
}
