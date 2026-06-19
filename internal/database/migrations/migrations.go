package migrations

import (
	"errors"
	"fmt"

	log "github.com/nef734/ggit/internal/log"
	"gorm.io/gorm"
)

const minDBVersion = 19

type Migration interface {
	Description() string
	Migrate(*gorm.DB) error
}

type migration struct {
	description string
	migrate     func(*gorm.DB) error
}

func NewMigration(desc string, fn func(*gorm.DB) error) Migration {
	return &migration{desc, fn}
}

func (m *migration) Description() string {
	return m.description
}

func (m *migration) Migrate(db *gorm.DB) error {
	return m.migrate(db)
}

// Version represents the version table. It should have only one row with `id == 1`.
type Version struct {
	ID      int64
	Version int64
}

// This is a sequence of migrations. Add new migrations to the bottom of the list.
// If you want to "retire" a migration, remove it from the top of the list and
// update _MIN_VER_DB accordingly
var migrations = []Migration{
	// v0 -> v4 : before 0.6.0 -> last support 0.7.33
	// v4 -> v10: before 0.7.0 -> last support 0.9.141
	// v10 -> v19: before 0.11.55 -> last support 0.12.0

	// Add new migration here, example:
	// v18 -> v19:v0.11.55
	// NewMigration("clean unlinked webhook and hook_tasks", cleanUnlinkedWebhookAndHookTasks),

	// v19 -> v20:v0.13.0
	NewMigration("migrate access tokens to store SHA56", migrateAccessTokenToSHA256),
	// v20 -> v21:v0.13.0
	NewMigration("add index to action.user_id", addIndexToActionUserID),
	// v21 -> v22:v0.13.0
	//
	// NOTE: There was a bug in calculating the value of the `version.version`
	// column after a migration is done, thus some instances are on v21 but some are
	// on v22. Let's make a noop v22 to make sure every instance will not miss a
	// real future migration.
	NewMigration("noop", func(*gorm.DB) error { return nil }),
	// v22 -> v23:v0.13.0
	NewMigration("add disable_news_feed to user table", addDisableNewsFeedToUser),
	// v23 -> v24:v0.13.0
	NewMigration("add is_archived to repository table", addIsArchivedToRepository),
	// v24 -> v25:v0.13.0
	NewMigration("add num_commits and num_branches to repository table", addNumCommitsAndBranchesToRepository),
	// v25 -> v26:v0.13.0
	NewMigration("create login_session table", createLoginSessionTable),
	// v26 -> v27:v0.13.0
	NewMigration("add branch protection and reviews", addBranchProtectionAndReviews),
	// v27 -> v28:v0.13.0
	NewMigration("create user_block table", createUserBlockTable),
	// v28 -> v29:v0.13.0
	NewMigration("add activity timestamps and sentinel to repository table", addActivityTimestampsAndSentinelToRepository),
	// v29 -> v30:v0.13.0
	NewMigration("add optimized indexes to action and commit_status tables", addOptimizedIndexes),
	// v30 -> v31:v0.13.0
	NewMigration("add hill_state to milestone table", v31),
	// v31 -> v32:v0.13.1
	NewMigration("add adr_folder_name to repository table", v32),
	// v32 -> v33:v0.13.2
	NewMigration("add size to attachment table", v33),
}

var errMigrationSkipped = errors.New("the migration has been skipped")

// Migrate migrates the database schema and/or data to the current version.
func Migrate(db *gorm.DB) error {
	// NOTE: GORM has problem migrating tables that happen to have columns with the
	// same name, see https://github.com/gogs/gogs/issues/7056.
	if !db.Migrator().HasTable(new(Version)) {
		err := db.AutoMigrate(new(Version))
		if err != nil {
			return fmt.Errorf(`auto migrate "version" table`+": %v", err)
		}
	}

	var current Version
	err := db.Where("id = ?", 1).First(&current).Error
	if err == gorm.ErrRecordNotFound {
		err = db.Create(
			&Version{
				ID:      1,
				Version: int64(minDBVersion + len(migrations)),
			},
		).Error
		if err != nil {
			return fmt.Errorf("create the version record"+": %v", err)
		}
		return nil

	} else if err != nil {
		return fmt.Errorf("get the version record"+": %v", err)
	}

	if minDBVersion > current.Version {
		log.Fatal(`
Hi there, thank you for using Ggit for so long!
However, Ggit has stopped supporting auto-migration from your previously installed version.
But the good news is, it's very easy to fix this problem!
You can migrate your older database using a previous release, then you can upgrade to the newest version.

Please save following instructions to somewhere and start working:

- If you were using below 0.6.0 (e.g. 0.5.x), download last supported archive from following link:
	https://github.com/gogs/gogs/releases/tag/v0.7.33
- If you were using below 0.7.0 (e.g. 0.6.x), download last supported archive from following link:
	https://github.com/gogs/gogs/releases/tag/v0.9.141
- If you were using below 0.11.55 (e.g. 0.9.141), download last supported archive from following link:
	https://github.com/gogs/gogs/releases/tag/v0.12.0

Once finished downloading:

1. Extract the archive and to upgrade steps as usual.
2. Run it once. To verify, you should see some migration traces.
3. Once it starts web server successfully, stop it.
4. Now it's time to put back the release archive you originally intent to upgrade.
5. Enjoy!

In case you're stilling getting this notice, go through instructions again until it disappears.`)
		return nil
	}

	if int(current.Version-minDBVersion) > len(migrations) {
		// User downgraded Gogs.
		current.Version = int64(len(migrations) + minDBVersion)
		return db.Where("id = ?", current.ID).Updates(current).Error
	}

	for _, m := range migrations[current.Version-minDBVersion:] {
		log.Info("Migration: %s", m.Description())
		if err = m.Migrate(db); err != nil {
			if err != errMigrationSkipped {
				return fmt.Errorf("do migrate"+": %v", err)
			}
			log.Trace("The migration %q has been skipped", m.Description())
		}

		current.Version++
		err = db.Where("id = ?", current.ID).Updates(current).Error
		if err != nil {
			return fmt.Errorf("update the version record"+": %v", err)
		}
	}
	return nil
}
