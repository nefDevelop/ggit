package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	log "unknwon.dev/clog/v2"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/dbx"
)

func newLogWriter() (logger.Writer, error) {
	sec := conf.File.Section("log.gorm")
	w, err := log.NewFileWriter(
		filepath.Join(conf.Log.RootPath, "gorm.log"),
		log.FileRotationConfig{
			Rotate:  sec.Key("ROTATE").MustBool(true),
			Daily:   sec.Key("ROTATE_DAILY").MustBool(true),
			MaxSize: sec.Key("MAX_SIZE").MustInt64(100) * 1024 * 1024,
			MaxDays: sec.Key("MAX_DAYS").MustInt64(3),
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, `create "gorm.log"`)
	}
	return &dbx.Logger{Writer: w}, nil
}

// Tables is the list of struct-to-table mappings.
//
// NOTE: Lines are sorted in alphabetical order, each letter in its own line.
//
// ⚠️ WARNING: This list is meant to be read-only.
var Tables = []any{
	new(Access), new(AccessToken), new(Action), new(Attachment),
	new(Collaboration), new(Comment),
	new(DeployKey),
	new(EmailAddress),
	new(HookTask),
	new(Issue), new(IssueLabel), new(IssueUser),
	new(Label), new(LFSObject), new(LoginSession), new(LoginSource),
	new(Milestone),
	new(Notice),
	new(OrgUser),
	new(ProtectBranch), new(ProtectBranchWhitelist), new(PublicKey), new(PullRequest),
	new(Release), new(Repository), new(RepoTopic),
	new(Star),
	new(Team), new(TeamRepo), new(TeamUser), new(Topic), new(TwoFactor), new(TwoFactorRecoveryCode),
	new(Upload), new(User),
	new(Version),
	new(Watch), new(Webhook),
}

// NewConnection returns a new database connection with the given logger.
func NewConnection(w logger.Writer) (*gorm.DB, error) {
	level := logger.Info
	if conf.IsProdMode() {
		level = logger.Warn
	}

	// NOTE: AutoMigrate does not respect logger passed in gorm.Config.
	logger.Default = logger.New(w, logger.Config{
		SlowThreshold: 100 * time.Millisecond,
		LogLevel:      level,
	})

	// Ensure database directory exists
	if conf.Database.Type == "sqlite3" {
		if err := os.MkdirAll(filepath.Dir(conf.Database.Path), 0755); err != nil {
			return nil, errors.Wrap(err, "create database directory")
		}
	}

	db, err := dbx.OpenDB(
		conf.Database,
		&gorm.Config{
			SkipDefaultTransaction: true,
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			NowFunc: func() time.Time {
				return time.Now().UTC().Truncate(time.Microsecond)
			},
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "open database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.Wrap(err, "get underlying *sql.DB")
	}

	maxOpenConns := conf.Database.MaxOpenConns
	if conf.Database.Type == "sqlite3" {
		maxOpenConns = 2
	}
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(conf.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Minute)

	if conf.Database.Type == "sqlite3" {
		conf.UseSQLite3 = true
	} else {
		panic("unreachable")
	}

	// NOTE: GORM has problem detecting existing columns, see
	// https://github.com/gogs/gogs/issues/6091. Therefore, only use it to create new
	// tables, and do customize migration with future changes.
	for _, table := range Tables {
		if db.Migrator().HasTable(table) {
			continue
		}

		name := strings.TrimPrefix(fmt.Sprintf("%T", table), "*database.")
		err = db.Migrator().AutoMigrate(table)
		if err != nil {
			return nil, errors.Wrapf(err, "auto migrate %q", name)
		}
		log.Trace("Auto migrated %q", name)
	}

	loadedLoginSourceFilesStore, err = loadLoginSourceFiles(filepath.Join(conf.CustomDir(), "conf", "auth.d"), db.NowFunc)
	if err != nil {
		return nil, errors.Wrap(err, "load login source files")
	}

	// Initialize the database handle.
	Handle = &DB{db: db}
	return db, nil
}

// Engine is a type alias for *gorm.DB to maintain compatibility with legacy
// function signatures that were using the XORM Engine interface.
type Engine = *gorm.DB

// DB is the database handler for the storage layer.
type DB struct {
	db *gorm.DB
}

// Handle is the global database handle. It could be `nil` during the
// installation mode.
//
// NOTE: Because we need to register all the routes even during the installation
// mode (which initially has no database configuration), we have to use a global
// variable since we can't pass a database handler around before it's available.
//
// NOTE: It is not guarded by a mutex because it is only written once either
// during the service start or during the installation process (which is a
// single-thread process).
var Handle *DB

// WithTransaction returns a new DB handler with the given transaction.
func (db *DB) WithTransaction(tx *gorm.DB) *DB {
	return &DB{db: tx}
}

// WithContext returns a new DB handler with the given context.
func (db *DB) WithContext(ctx context.Context) *DB {
	return &DB{db: db.db.WithContext(ctx)}
}

// Transaction executes the given function within a database transaction.
func (db *DB) Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	return db.db.Transaction(fc, opts...)
}

func (db *DB) AccessTokens() *AccessTokensStore {
	return newAccessTokensStore(db.db)
}

func (db *DB) Actions() *ActionsStore {
	return newActionsStore(db.db)
}

func (db *DB) LFS() *LFSStore {
	return newLFSStore(db.db)
}

// NOTE: It is not guarded by a mutex because it only gets written during the
// service start.
var loadedLoginSourceFilesStore loginSourceFilesStore

func (db *DB) LoginSources() *LoginSourcesStore {
	return newLoginSourcesStore(db.db, loadedLoginSourceFilesStore)
}

func (db *DB) LoginSessions() *LoginSessionsStore {
	return newLoginSessionsStore(db.db)
}

func (db *DB) Notices() *NoticesStore {
	return newNoticesStore(db.db)
}

func (db *DB) Organizations() *OrganizationsStore {
	return newOrganizationsStoreStore(db.db)
}

func (db *DB) Permissions() *PermissionsStore {
	return newPermissionsStore(db.db)
}

func (db *DB) PublicKey() *PublicKeysStore {
	return newPublicKeysStore(db.db)
}

func (db *DB) Repositories() *RepositoriesStore {
	return newReposStore(db.db)
}

func (db *DB) Issues() *IssuesStore {
	return newIssuesStore(db.db)
}

func (db *DB) TwoFactors() *TwoFactorsStore {
	return newTwoFactorsStore(db.db)
}

func (db *DB) Users() *UsersStore {
	return newUsersStore(db.db)
}
