package dbx

import (
	"fmt"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"github.com/nef734/ggit/internal/conf"
	"gorm.io/gorm"
)

// NewDSN takes given database options and returns parsed DSN.
func NewDSN(opts conf.DatabaseOpts) (dsn string, err error) {
	if opts.Type != "sqlite3" {
		return "", fmt.Errorf("unrecognized dialect: %s", opts.Type)
	}

	path, _ := filepath.Abs(opts.Path)

	// Two decisions here are load bearing under concurrency:
	//
	// cache=shared is deliberately not set. In shared cache mode SQLite reports
	// a conflicting table lock as SQLITE_LOCKED (6) instead of SQLITE_BUSY (5),
	// and the busy handler is not consulted for it, so busy_timeout silently
	// stops helping and concurrent writers fail outright with "database table is
	// locked: database is deadlocked". WAL plus ordinary file locking already
	// lets readers run beside a writer, which is all a single-process server
	// needs.
	//
	// _txlock is deliberately left at the deferred default. Immediate mode would
	// avoid the read-then-upgrade deadlock, but it also takes the write lock at
	// BEGIN for every transaction, including read only ones, so all of them
	// serialise. That measured worse than it fixed: busy_timeout could not absorb
	// the contention. Shrinking the transactions that span slow work belongs in
	// the callers instead.
	return "file://" + path + "?mode=rwc&_journal_mode=WAL&_synchronous=NORMAL", nil
}

// OpenDB opens a new database connection encapsulated as gorm.DB using given
// database options and GORM config.
func OpenDB(opts conf.DatabaseOpts, cfg *gorm.Config) (*gorm.DB, error) {
	dsn, err := NewDSN(opts)
	if err != nil {
		return nil, fmt.Errorf("parse DSN"+": %w", err)
	}

	return gorm.Open(sqlite.Open(dsn), cfg)
}
