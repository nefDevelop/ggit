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
	return "file://" + path + "?cache=shared&mode=rwc&_journal_mode=WAL&_synchronous=NORMAL", nil
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
