package dbx

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nef734/ggit/internal/conf"
)

func TestNewDSN(t *testing.T) {
	t.Run("bad dialect", func(t *testing.T) {
		_, err := NewDSN(conf.DatabaseOpts{
			Type: "bad_dialect",
		})
		assert.Equal(t, "unrecognized dialect: bad_dialect", fmt.Sprintf("%v", err))
	})

	tests := []struct {
		name    string
		opts    conf.DatabaseOpts
		wantDSN string
	}{
		{
			name: "sqlite3",
			opts: conf.DatabaseOpts{
				Type: "sqlite3",
				Path: "/tmp/gogs.db",
			},
			wantDSN: "file:///tmp/gogs.db?cache=shared&mode=rwc&_journal_mode=WAL&_synchronous=NORMAL",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path, _ := filepath.Abs(test.opts.Path)
			test.wantDSN = "file://" + path + "?cache=shared&mode=rwc&_journal_mode=WAL&_synchronous=NORMAL"

			dsn, err := NewDSN(test.opts)
			require.NoError(t, err)
			assert.Equal(t, test.wantDSN, dsn)
		})
	}
}

func TestQuote(t *testing.T) {
	assert.Equal(t, "`table`", Quote("table"))
}
