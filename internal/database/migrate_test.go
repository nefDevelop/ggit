package database

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gogs/git-module"
	"github.com/nef734/ggit/internal/conf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func Test_MigrateRepository(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	// 1. Setup Mock Environment
	baseDir := filepath.Join(os.TempDir(), fmt.Sprintf("ggit-migrate-test-%d", os.Getpid()))
	err := os.MkdirAll(baseDir, 0755)
	require.NoError(t, err)
	defer os.RemoveAll(baseDir)

	sourcePath := filepath.Join(baseDir, "source")
	destRoot := filepath.Join(baseDir, "dest")
	err = os.MkdirAll(sourcePath, 0755)
	require.NoError(t, err)

	conf.SetMockRepository(t, conf.RepositoryOpts{
		Root:             destRoot,
		MaxCreationLimit: -1,
		DefaultBranch:    "main",
	})
	conf.Git.Timeout.Migrate = 600

	// 2. Create a source repository with some commits
	runGit := func(dir string, args ...string) {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		err := cmd.Run()
		require.NoError(t, err, "git %v failed", args)
	}

	runGit(sourcePath, "init", "-b", "main")
	runGit(sourcePath, "config", "user.email", "test@example.com")
	runGit(sourcePath, "config", "user.name", "Test User")
	err = os.WriteFile(filepath.Join(sourcePath, "README.md"), []byte("# Test Repo"), 0644)
	require.NoError(t, err)
	runGit(sourcePath, "add", ".")
	runGit(sourcePath, "commit", "-m", "Initial commit")

	// 4. Run Migration Tests
	setup := func(subtest string) (*User, string, *gorm.DB) {
		db := newTestDB(t, "migrate-"+subtest)
		Handle = &DB{db: db}
		owner := &User{
			Name:            "owner-" + subtest,
			MaxRepoCreation: -1,
		}
		require.NoError(t, db.Create(owner).Error)
		return owner, "file://" + sourcePath, db
	}

	t.Run("Migrate local repository with data", func(t *testing.T) {
		owner, addr, db := setup("with-data")
		opts := MigrateRepoOptions{
			Name:        "migrated-repo",
			RemoteAddr:  addr,
			Description: "A migrated repository",
		}

		repo, err := MigrateRepository(owner, owner, opts)
		require.NoError(t, err)
		assert.False(t, repo.IsBare)

		// 5. TEST: LOAD CODE CONTENT
		gitRepo, err := git.Open(repo.RepoPath())
		require.NoError(t, err)

		commit, err := gitRepo.CatFileCommit("main")
		require.NoError(t, err)

		entry, err := commit.TreeEntry("README.md")
		require.NoError(t, err)

		blob := entry.Blob()
		require.NotNil(t, blob)

		data, err := blob.Bytes()
		require.NoError(t, err)
		assert.Equal(t, "# Test Repo", string(data), "File content should match the source")

		dbRepo := new(Repository)
		err = db.Where("id = ?", repo.ID).First(dbRepo).Error
		require.NoError(t, err)
		assert.False(t, dbRepo.IsBare)
	})

	t.Run("Migrate empty repository", func(t *testing.T) {
		owner, _, _ := setup("empty")
		emptySource := filepath.Join(baseDir, "empty-source")
		_ = os.MkdirAll(emptySource, 0755)
		runGit(emptySource, "init", "-b", "master")

		opts := MigrateRepoOptions{
			Name:       "empty-migrated",
			RemoteAddr: "file://" + emptySource,
		}

		repo, err := MigrateRepository(owner, owner, opts)
		require.NoError(t, err)
		assert.True(t, repo.IsBare)
	})

	t.Run("Invalid remote address", func(t *testing.T) {
		owner, _, _ := setup("invalid-addr")
		opts := MigrateRepoOptions{
			Name:       "invalid-addr",
			RemoteAddr: "https://non-existent-domain.com/repo.git",
		}

		_, err := MigrateRepository(owner, owner, opts)
		assert.Error(t, err)
	})

	t.Run("Malicious command injection attempt", func(t *testing.T) {
		owner, _, _ := setup("injection")
		opts := MigrateRepoOptions{
			Name:       "injection-attempt",
			RemoteAddr: "https://github.com/test/repo.git; touch " + filepath.Join(baseDir, "hacked-test.txt"),
		}

		_, err := MigrateRepository(owner, owner, opts)
		assert.Error(t, err)
		assert.NoFileExists(t, filepath.Join(baseDir, "hacked-test.txt"))
	})

	t.Run("Duplicate migration attempt", func(t *testing.T) {
		owner, addr, _ := setup("duplicate")
		opts := MigrateRepoOptions{
			Name:       "duplicate-test",
			RemoteAddr: addr,
		}

		_, err := MigrateRepository(owner, owner, opts)
		require.NoError(t, err)

		_, err = MigrateRepository(owner, owner, opts)
		require.Error(t, err)
		assert.True(t, IsErrRepoAlreadyExist(err), "expected ErrRepoAlreadyExist, got %v (%T)", err, err)
	})

	t.Run("Reserved repository name", func(t *testing.T) {
		owner, addr, _ := setup("reserved")
		opts := MigrateRepoOptions{
			Name:       "test.git", // Truly reserved pattern
			RemoteAddr: addr,
		}

		_, err := MigrateRepository(owner, owner, opts)
		require.Error(t, err) // This will fail the test elegantly if err is nil
		assert.Contains(t, err.Error(), "reserved")
	})

	t.Run("URL with special characters in credentials", func(t *testing.T) {
		owner, _, _ := setup("special-url")
		opts := MigrateRepoOptions{
			Name:       "special-url",
			RemoteAddr: "https://user:p@ssword!@github.com/test/repo.git",
		}

		_, err := MigrateRepository(owner, owner, opts)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "clone")
	})
}
