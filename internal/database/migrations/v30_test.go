package migrations

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nef734/ggit/internal/dbtest"
)

func Test_addOptimizedIndexes(t *testing.T) {
	db := dbtest.NewDB(t, "addOptimizedIndexes")

	// Create tables first
	type action struct {
		ID     int64 `gorm:"primaryKey"`
		UserID int64
		RepoID int64
	}
	type commit_status struct {
		ID      int64 `gorm:"primaryKey"`
		RepoID  int64
		SHA     string
		Context string
	}
	err := db.AutoMigrate(&action{}, &commit_status{})
	require.NoError(t, err)

	// Run migration
	err = addOptimizedIndexes(db)
	assert.NoError(t, err)

	// Verify indexes
	assert.True(t, db.Migrator().HasIndex(&action{}, "idx_action_user_id_id"))
	assert.True(t, db.Migrator().HasIndex(&action{}, "idx_action_repo_id_id"))
	assert.True(t, db.Migrator().HasIndex(&commit_status{}, "idx_commit_status_repo_sha_context"))
}
