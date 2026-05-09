package migrations

import (
	"gorm.io/gorm"
)

func addNumCommitsAndBranchesToRepository(db *gorm.DB) error {
	type repository struct {
		NumCommits  int64
		NumBranches int
	}

	if err := db.Migrator().AutoMigrate(&repository{}); err != nil {
		return err
	}
	return nil
}
