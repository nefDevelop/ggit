package backup

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nef734/ggit/internal/conf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBackup(t *testing.T) {
	// Setup temporary environment
	tmpDir, err := os.MkdirTemp("", "ggit-backup-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Mock configuration paths
	conf.Server.AppDataPath = tmpDir
	conf.Database.Path = filepath.Join(tmpDir, "ggit.db")
	conf.Attachment.Path = filepath.Join(tmpDir, "attachments")
	
	// Create dummy data
	err = os.WriteFile(conf.Database.Path, []byte("fake sqlite content"), 0644)
	require.NoError(t, err)
	
	err = os.MkdirAll(conf.Attachment.Path, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(conf.Attachment.Path, "test.txt"), []byte("hello"), 0644)
	require.NoError(t, err)

	t.Run("CreateMetadataBackup", func(t *testing.T) {
		path, err := CreateMetadataBackup()
		require.NoError(t, err)
		assert.FileExists(t, path)
		assert.Contains(t, path, "ggit-metadata-")
	})

	t.Run("CleanupOldBackups", func(t *testing.T) {
		backupDir := filepath.Join(tmpDir, "backups")
		err := os.MkdirAll(backupDir, 0755)
		require.NoError(t, err)

		// Create a new backup and an old one
		newBackup := filepath.Join(backupDir, "new.zip")
		oldBackup := filepath.Join(backupDir, "old.zip")
		
		err = os.WriteFile(newBackup, []byte("new"), 0644)
		require.NoError(t, err)
		err = os.WriteFile(oldBackup, []byte("old"), 0644)
		require.NoError(t, err)

		// Backdate the old backup
		oldTime := time.Now().Add(-10 * 24 * time.Hour)
		err = os.Chtimes(oldBackup, oldTime, oldTime)
		require.NoError(t, err)

		err = CleanupOldBackups()
		require.NoError(t, err)

		assert.FileExists(t, newBackup)
		assert.NoFileExists(t, oldBackup)
	})
}
