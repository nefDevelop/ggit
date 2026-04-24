package database

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nef734/ggit/internal/conf"
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

	// 3. Initialize Test DB and User
	db := newTestDB(t, "migrate")
	Handle = &DB{db: db}
	owner := &User{
		Name:            "testowner",
		MaxRepoCreation: -1,
	}
	err = db.Create(owner).Error
	require.NoError(t, err)

	// 4. Run Migration Tests
	t.Run("Migrate local repository with data", func(t *testing.T) {
		opts := MigrateRepoOptions{
			Name:        "migrated-repo",
			RemoteAddr:  "file://" + sourcePath,
			Description: "A migrated repository",
		}

		repo, err := MigrateRepository(owner, owner, opts)
		require.NoError(t, err)
		assert.NotNil(t, repo)
		assert.Equal(t, "main", repo.DefaultBranch)
		assert.False(t, repo.IsBare)

		dbRepo := new(Repository)
		err = db.Where("id = ?", repo.ID).First(dbRepo).Error
		require.NoError(t, err)
		assert.False(t, dbRepo.IsBare)
	})

	t.Run("Migrate empty repository", func(t *testing.T) {
		emptySource := filepath.Join(baseDir, "empty-source")
		err = os.MkdirAll(emptySource, 0755)
		require.NoError(t, err)
		runGit(emptySource, "init", "-b", "master")

		opts := MigrateRepoOptions{
			Name:       "empty-migrated",
			RemoteAddr: "file://" + emptySource,
		}

		repo, err := MigrateRepository(owner, owner, opts)
		require.NoError(t, err)
		assert.True(t, repo.IsBare, "Empty repo should be marked as Bare")
	})

	t.Run("Invalid remote address", func(t *testing.T) {
		opts := MigrateRepoOptions{
			Name:       "invalid-addr",
			RemoteAddr: "https://non-existent-domain-734.com/repo.git",
		}

		_, err := MigrateRepository(owner, owner, opts)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "clone")
	})

	t.Run("Malicious command injection attempt", func(t *testing.T) {
		opts := MigrateRepoOptions{
			Name:       "injection-attempt",
			RemoteAddr: "https://github.com/test/repo.git; touch " + filepath.Join(baseDir, "hacked.txt"),
		}

		_, err := MigrateRepository(owner, owner, opts)
		assert.Error(t, err)
		// Should not create the file hacked.txt because it's not executed by shell
		assert.NoFileExists(t, filepath.Join(baseDir, "hacked.txt"))
	})
}
