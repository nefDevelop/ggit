package database

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/nef734/ggit/internal/auth"
	"github.com/nef734/ggit/internal/auth/ldap"
	"github.com/nef734/ggit/internal/cryptox"
	"github.com/nef734/ggit/internal/dbtest"
	"github.com/nef734/ggit/internal/domain"
	"github.com/nef734/ggit/internal/lfsx"
	"github.com/nef734/ggit/internal/testx"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestDumpAndImport(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	t.Parallel()

	backupTables := []any{
		new(Access), new(AccessToken), new(Action), new(EmailAddress),
		new(LFSObject), new(LoginSource), new(Notice),
	}

	const wantTables = 7
	if len(backupTables) != wantTables {
		t.Fatalf("New table has added (want %d got %d), please add new tests for the table and update this check", wantTables, len(backupTables))
	}

	db := dbtest.NewDB(t, "dumpAndImport", backupTables...)
	setupDBToDump(t, db)
	dumpTables(t, backupTables, db)
	importTables(t, backupTables, db)

	// Dump and assert golden again to make sure data aren't changed.
	dumpTables(t, backupTables, db)
}

func setupDBToDump(t *testing.T, db *gorm.DB) {
	vals := []any{
		&Access{
			ID:     1,
			UserID: 1,
			RepoID: 11,
			Mode:   domain.AccessModeRead,
		},
		&Access{
			ID:     2,
			UserID: 2,
			RepoID: 22,
			Mode:   domain.AccessModeWrite,
		},

		&AccessToken{
			UserID:      1,
			Name:        "test1",
			Sha1:        cryptox.SHA1("2910d03d-c0b5-4f71-bad5-c4086e4efae3"),
			SHA256:      cryptox.SHA256(cryptox.SHA1("2910d03d-c0b5-4f71-bad5-c4086e4efae3")),
			CreatedUnix: 1588568886,
			UpdatedUnix: 1588572486, // 1 hour later
		},
		&AccessToken{
			UserID:      1,
			Name:        "test2",
			Sha1:        cryptox.SHA1("84117e17-7e67-4024-bd04-1c23e6e809d4"),
			SHA256:      cryptox.SHA256(cryptox.SHA1("84117e17-7e67-4024-bd04-1c23e6e809d4")),
			CreatedUnix: 1588568886,
		},
		&AccessToken{
			UserID:      2,
			Name:        "test1",
			Sha1:        cryptox.SHA1("da2775ce-73dd-47ba-b9d2-bbcc346585c4"),
			SHA256:      cryptox.SHA256(cryptox.SHA1("da2775ce-73dd-47ba-b9d2-bbcc346585c4")),
			CreatedUnix: 1588568886,
		},
		&AccessToken{
			UserID:      2,
			Name:        "test2",
			Sha1:        cryptox.SHA256(cryptox.SHA1("1b2dccd1-a262-470f-bb8c-7fc73192e9bb"))[:40],
			SHA256:      cryptox.SHA256(cryptox.SHA1("1b2dccd1-a262-470f-bb8c-7fc73192e9bb")),
			CreatedUnix: 1588568886,
		},

		&Action{
			ID:           1,
			UserID:       1,
			OpType:       ActionCreateBranch,
			ActUserID:    1,
			ActUserName:  "alice",
			RepoID:       1,
			RepoUserName: "alice",
			RepoName:     "example",
			RefName:      "main",
			IsPrivate:    false,
			Content:      `{"Len":1,"Commits":[],"CompareURL":""}`,
			CreatedUnix:  1588568886,
		},
		&Action{
			ID:           2,
			UserID:       1,
			OpType:       ActionCommitRepo,
			ActUserID:    1,
			ActUserName:  "alice",
			RepoID:       1,
			RepoUserName: "alice",
			RepoName:     "example",
			RefName:      "main",
			IsPrivate:    false,
			Content:      `{"Len":1,"Commits":[],"CompareURL":""}`,
			CreatedUnix:  1588568886,
		},
		&Action{
			ID:           3,
			UserID:       1,
			OpType:       ActionDeleteBranch,
			ActUserID:    1,
			ActUserName:  "alice",
			RepoID:       1,
			RepoUserName: "alice",
			RepoName:     "example",
			RefName:      "main",
			IsPrivate:    false,
			CreatedUnix:  1588568886,
		},

		&EmailAddress{
			ID:          1,
			UserID:      1,
			Email:       "alice@example.com",
			IsActivated: false,
		},
		&EmailAddress{
			ID:          2,
			UserID:      2,
			Email:       "bob@example.com",
			IsActivated: true,
		},

		&LFSObject{
			RepoID:      1,
			OID:         "ef797c8118f02dfb649607dd5d3f8c7623048c9c063d532cc95c5ed7a898a64f",
			Size:        100,
			Storage:     lfsx.StorageLocal,
			CreatedUnix: 1588568886,
		},
		&LFSObject{
			RepoID:      2,
			OID:         "ef797c8118f02dfb649607dd5d3f8c7623048c9c063d532cc95c5ed7a898a64f",
			Size:        100,
			Storage:     lfsx.StorageLocal,
			CreatedUnix: 1588568886,
		},

		&LoginSource{
			Type:      auth.LDAP,
			Name:      "My LDAP",
			IsActived: true,
			Provider: ldap.NewProvider(false, &ldap.Config{
				Host: "ldap.example.com",
			}),
			CreatedUnix: 1588568886,
			UpdatedUnix: 1588572486, // 1 hour later
		},
		&LoginSource{
			Type:      auth.DLDAP,
			Name:      "Direct LDAP",
			IsActived: true,
			Provider: ldap.NewProvider(true, &ldap.Config{
				Host: "ldap.example.com",
			}),
			CreatedUnix: 1588568886,
		},

		&Notice{
			ID:          1,
			Type:        NoticeTypeRepository,
			Description: "This is a notice",
			CreatedUnix: 1588568886,
		},
	}
	for _, val := range vals {
		err := db.Create(val).Error
		require.NoError(t, err)
	}
}

func dumpTables(t *testing.T, backupTables []any, db *gorm.DB) {
	ctx := context.Background()

	for _, table := range backupTables {
		tableName := getTableType(table)

		var buf bytes.Buffer
		err := dumpTable(ctx, db, table, &buf)
		if err != nil {
			t.Fatalf("%s: %v", tableName, err)
		}

		golden := filepath.Join("testdata", "backup", tableName+".golden.json")
		testx.AssertGolden(t, golden, testx.Update("TestDumpAndImport"), buf.String())
	}
}

func importTables(t *testing.T, backupTables []any, db *gorm.DB) {
	ctx := context.Background()

	for _, table := range backupTables {
		tableName := getTableType(table)

		err := func() error {
			golden := filepath.Join("testdata", "backup", tableName+".golden.json")
			f, err := os.Open(golden)
			if err != nil {
				return fmt.Errorf("open table file"+": %w", err)
			}
			defer func() { _ = f.Close() }()

			return importTable(ctx, db, table, f)
		}()
		if err != nil {
			t.Fatalf("%s: %v", tableName, err)
		}
	}
}
