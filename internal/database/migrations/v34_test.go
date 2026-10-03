package migrations

import (
	"testing"

	"github.com/nef734/ggit/internal/dbtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// expectedIndexes is the set addQueryIndexes is responsible for. The names are
// asserted rather than the statements so that a rename shows up as a failure.
var expectedIndexes = []string{
	"idx_access_repo_id_mode",
	"idx_hook_task_delivered",
	"idx_watch_repo_id",
	"idx_star_repo_id",
	"idx_repo_owner_updated",
	"idx_repo_public_updated",
	"idx_repo_fork_id",
	"idx_notice_id_desc",
	"idx_action_act_user_created",
}

func TestV34CreatesIndexes(t *testing.T) {
	db := dbtest.NewDB(t, "MigrationsV34")

	// The tables have to exist before the index statements can run. DDL is used
	// rather than AutoMigrate so the test does not depend on the model structs.
	for _, ddl := range []string{
		"CREATE TABLE access (id INTEGER PRIMARY KEY, repo_id INTEGER, user_id INTEGER, mode INTEGER)",
		"CREATE TABLE hook_task (id INTEGER PRIMARY KEY, is_delivered BOOLEAN, repo_id INTEGER)",
		"CREATE TABLE watch (id INTEGER PRIMARY KEY, repo_id INTEGER, user_id INTEGER)",
		"CREATE TABLE star (id INTEGER PRIMARY KEY, repo_id INTEGER, user_id INTEGER)",
		"CREATE TABLE repository (id INTEGER PRIMARY KEY, owner_id INTEGER, updated_unix INTEGER, is_private BOOLEAN, is_unlisted BOOLEAN, fork_id INTEGER)",
		"CREATE TABLE notice (id INTEGER PRIMARY KEY, created_unix INTEGER)",
		"CREATE TABLE action (id INTEGER PRIMARY KEY, act_user_id INTEGER, created_unix INTEGER)",
	} {
		require.NoError(t, db.Exec(ddl).Error, ddl)
	}

	require.NoError(t, v34(db))

	// Migrator().HasIndex derives the table name from the value it is given, so
	// the probes below are named after the real tables. An anonymous struct would
	// resolve to a table called "structs" and report every index as missing.
	type access struct{}
	type hookTask struct{}
	type watch struct{}
	type star struct{}
	type repository struct{}
	type notice struct{}
	type action struct{}

	owners := map[string]any{
		"idx_access_repo_id_mode":     &access{},
		"idx_hook_task_delivered":     &hookTask{},
		"idx_watch_repo_id":           &watch{},
		"idx_star_repo_id":            &star{},
		"idx_repo_owner_updated":      &repository{},
		"idx_repo_public_updated":     &repository{},
		"idx_repo_fork_id":            &repository{},
		"idx_notice_id_desc":          &notice{},
		"idx_action_act_user_created": &action{},
	}

	for _, name := range expectedIndexes {
		owner, ok := owners[name]
		require.True(t, ok, "indice sin propietario en el test: %s", name)
		assert.True(t, db.Migrator().HasIndex(owner, name), "falta el indice %s", name)
	}
}

// TestV34IsIdempotent guards the migration path itself: Migrate replays every
// migration against databases that already carry some of them, and a
// non-idempotent statement would abort startup.
func TestV34IsIdempotent(t *testing.T) {
	db := dbtest.NewDB(t, "MigrationsV34Idempotent")

	for _, ddl := range []string{
		"CREATE TABLE access (id INTEGER PRIMARY KEY, repo_id INTEGER, user_id INTEGER, mode INTEGER)",
		"CREATE TABLE hook_task (id INTEGER PRIMARY KEY, is_delivered BOOLEAN, repo_id INTEGER)",
		"CREATE TABLE watch (id INTEGER PRIMARY KEY, repo_id INTEGER, user_id INTEGER)",
		"CREATE TABLE star (id INTEGER PRIMARY KEY, repo_id INTEGER, user_id INTEGER)",
		"CREATE TABLE repository (id INTEGER PRIMARY KEY, owner_id INTEGER, updated_unix INTEGER, is_private BOOLEAN, is_unlisted BOOLEAN, fork_id INTEGER)",
		"CREATE TABLE notice (id INTEGER PRIMARY KEY, created_unix INTEGER)",
		"CREATE TABLE action (id INTEGER PRIMARY KEY, act_user_id INTEGER, created_unix INTEGER)",
	} {
		require.NoError(t, db.Exec(ddl).Error, ddl)
	}

	require.NoError(t, v34(db))
	require.NoError(t, v34(db), "la segunda pasada no debe fallar")

	var count int64
	require.NoError(t, db.Raw(
		"SELECT COUNT(*) FROM sqlite_master WHERE type = 'index' AND name LIKE 'idx_%'").Scan(&count).Error)
	assert.GreaterOrEqual(t, count, int64(len(expectedIndexes)))
}
