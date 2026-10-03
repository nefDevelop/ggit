package migrations

import (
	"gorm.io/gorm"
)

// addQueryIndexes adds the indexes that the queries in internal/database actually
// rely on. Every index below is justified by a concrete statement:
//
//   - access(repo_id, mode): recalculateTeamAccesses filters on repo_id and mode
//     (repo.go), and three queries join access on access.repo_id. The only index
//     present was access(user_id, repo_id), which cannot serve a lookup keyed on
//     repo_id, so those scans walked the whole table.
//   - hook_task(is_delivered, repo_id): DeliverHooks reads every undelivered task
//     on each start and hook_task is never pruned, so the scan grew without bound.
//   - watch(repo_id) and star(repo_id): counters and watcher/star lists are keyed
//     by repo_id, including "SELECT COUNT(*) FROM watch WHERE repo_id = ?".
//     The existing unique index is on (user_id, repo_id), useless in reverse.
//   - repository(owner_id, updated_unix): the repository list of a user and the
//     explore listing both order by updated_unix within an owner.
//   - repository(is_private, is_unlisted, updated_unix): the public listing
//     filters on the two flags and then sorts by updated_unix.
//   - repository(fork_id): fork listings filter on it and nothing indexes it.
//   - notice(id DESC): the admin notice list sorts by id descending.
//   - action(act_user_id, created_unix): the activity feed filters by the acting
//     user and sorts by creation time; act_user_id has no index on its own.
func addQueryIndexes(db *gorm.DB) error {
	// Shadow structs with the indexed columns declared. GORM resolves the table
	// name from the value passed to Migrator().HasIndex, so an empty struct would
	// leave the check unable to name a table and it would always report the index
	// as missing, making the migration try to create it a second time.
	type Access struct {
		RepoID int64 `gorm:"column:repo_id"`
		Mode   int64 `gorm:"column:mode"`
	}
	type HookTask struct {
		IsDelivered bool  `gorm:"column:is_delivered"`
		RepoID      int64 `gorm:"column:repo_id"`
	}
	type Watch struct {
		RepoID int64 `gorm:"column:repo_id"`
	}
	type Star struct {
		RepoID int64 `gorm:"column:repo_id"`
	}
	type Repository struct {
		OwnerID     int64 `gorm:"column:owner_id"`
		UpdatedUnix int64 `gorm:"column:updated_unix"`
		IsPrivate   bool  `gorm:"column:is_private"`
		IsUnlisted  bool  `gorm:"column:is_unlisted"`
		ForkID      int64 `gorm:"column:fork_id"`
	}
	type Notice struct {
		ID int64 `gorm:"primaryKey;column:id"`
	}
	type Action struct {
		ActUserID   int64 `gorm:"column:act_user_id"`
		CreatedUnix int64 `gorm:"column:created_unix"`
	}

	indexes := []struct {
		owner any
		name  string
		stmt  string
	}{
		{&Access{}, "idx_access_repo_id_mode", "CREATE INDEX idx_access_repo_id_mode ON access(repo_id, mode)"},
		{&HookTask{}, "idx_hook_task_delivered", "CREATE INDEX idx_hook_task_delivered ON hook_task(is_delivered, repo_id)"},
		{&Watch{}, "idx_watch_repo_id", "CREATE INDEX idx_watch_repo_id ON watch(repo_id)"},
		{&Star{}, "idx_star_repo_id", "CREATE INDEX idx_star_repo_id ON star(repo_id)"},
		{&Repository{}, "idx_repo_owner_updated", "CREATE INDEX idx_repo_owner_updated ON repository(owner_id, updated_unix DESC)"},
		{&Repository{}, "idx_repo_public_updated", "CREATE INDEX idx_repo_public_updated ON repository(is_private, is_unlisted, updated_unix DESC)"},
		{&Repository{}, "idx_repo_fork_id", "CREATE INDEX idx_repo_fork_id ON repository(fork_id)"},
		{&Notice{}, "idx_notice_id_desc", "CREATE INDEX idx_notice_id_desc ON notice(id DESC)"},
		{&Action{}, "idx_action_act_user_created", "CREATE INDEX idx_action_act_user_created ON action(act_user_id, created_unix DESC)"},
	}

	for _, index := range indexes {
		if db.Migrator().HasIndex(index.owner, index.name) {
			continue
		}
		if err := db.Exec(index.stmt).Error; err != nil {
			return err
		}
	}

	return nil
}

func v34(db *gorm.DB) error {
	return addQueryIndexes(db)
}
