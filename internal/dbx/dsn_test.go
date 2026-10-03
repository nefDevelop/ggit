package dbx

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nef734/ggit/internal/conf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDSN(t *testing.T) {
	t.Run("bad dialect", func(t *testing.T) {
		_, err := NewDSN(conf.DatabaseOpts{
			Type: "bad_dialect",
		})
		assert.Equal(t, "unrecognized dialect: bad_dialect", fmt.Sprintf("%v", err))
	})

	t.Run("sqlite3", func(t *testing.T) {
		dsn, err := NewDSN(conf.DatabaseOpts{
			Type: "sqlite3",
			Path: "/tmp/ggit.db",
		})
		require.NoError(t, err)

		path, _ := filepath.Abs("/tmp/ggit.db")
		assert.True(t, strings.HasPrefix(dsn, "file://"+path+"?"))
		assert.Contains(t, dsn, "_journal_mode=WAL")

		// Shared cache turns a conflicting table lock into SQLITE_LOCKED, which the
		// busy handler does not retry, so concurrent writers fail outright.
		assert.NotContains(t, dsn, "cache=shared",
			"cache=shared provoca SQLITE_LOCKED y anula busy_timeout")

		// Immediate mode was measured and rejected: it takes the write lock at
		// BEGIN for every transaction, so all of them serialise.
		assert.NotContains(t, dsn, "_txlock=immediate")
	})
}
