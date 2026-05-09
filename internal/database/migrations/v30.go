package migrations

import (
	"gorm.io/gorm"
)

func addOptimizedIndexes(db *gorm.DB) error {
	type action struct {
		UserID int64 `gorm:"index:idx_action_user_id_id"`
		ID     int64 `gorm:"index:idx_action_user_id_id"`
		RepoID int64 `gorm:"index:idx_action_repo_id_id"`
		// ID included in RepoID index too for pagination
		ID2 int64 `gorm:"column:id;index:idx_action_repo_id_id"`
	}

	type commit_status struct {
		RepoID  int64  `gorm:"index:idx_commit_status_repo_sha_context"`
		SHA     string `gorm:"index:idx_commit_status_repo_sha_context"`
		Context string `gorm:"index:idx_commit_status_repo_sha_context"`
	}

	// We use Migrator().CreateIndex instead of AutoMigrate for composite indexes
	// to have better control over names and existing indexes.

	// Action indexes
	if !db.Migrator().HasIndex(&action{}, "idx_action_user_id_id") {
		err := db.Exec("CREATE INDEX idx_action_user_id_id ON action(user_id, id DESC)").Error
		if err != nil {
			return err
		}
	}
	if !db.Migrator().HasIndex(&action{}, "idx_action_repo_id_id") {
		err := db.Exec("CREATE INDEX idx_action_repo_id_id ON action(repo_id, id DESC)").Error
		if err != nil {
			return err
		}
	}

	// Commit Status index
	if !db.Migrator().HasIndex(&commit_status{}, "idx_commit_status_repo_sha_context") {
		err := db.Exec("CREATE INDEX idx_commit_status_repo_sha_context ON commit_status(repo_id, sha, context)").Error
		if err != nil {
			return err
		}
	}

	return nil
}
