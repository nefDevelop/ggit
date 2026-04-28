package backup

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/nef734/ggit/internal/conf"
)

// CreateMetadataBackup creates a ZIP archive with DB and assets.
func CreateMetadataBackup() (string, error) {
	backupDir := filepath.Join(conf.Server.AppDataPath, "backups")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", errors.Wrap(err, "create backup directory")
	}

	backupName := fmt.Sprintf("ggit-metadata-%s.zip", time.Now().Format("20060102-150405"))
	backupPath := filepath.Join(backupDir, backupName)
	
	zipFile, err := os.Create(backupPath)
	if err != nil {
		return "", err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	// 1. Database
	// NOTE: For SQLite, copying the file while active might be slightly risky but usually works for simple backups.
	// In the future, we could use "VACUUM INTO" for a clean copy.
	if err := addFileToZip(archive, conf.Database.Path, "ggit.db"); err != nil {
		return "", errors.Wrap(err, "add database")
	}

	// 2. Attachments
	_ = addDirToZip(archive, conf.Attachment.Path, "attachments")

	// 3. User Avatars
	_ = addDirToZip(archive, conf.Picture.AvatarUploadPath, "avatars/users")

	// 4. Repo Avatars
	_ = addDirToZip(archive, conf.Picture.RepositoryAvatarUploadPath, "avatars/repos")

	return backupPath, nil
}

func addFileToZip(archive *zip.Writer, path, name string) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	w, err := archive.Create(name)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, file)
	return err
}

func addDirToZip(archive *zip.Writer, root, prefix string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, _ := filepath.Rel(root, path)
		return addFileToZip(archive, path, filepath.Join(prefix, relPath))
	})
}

// CleanupOldBackups removes backups older than 7 days.
func CleanupOldBackups() error {
	backupDir := filepath.Join(conf.Server.AppDataPath, "backups")
	files, err := os.ReadDir(backupDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	now := time.Now()
	for _, f := range files {
		info, err := f.Info()
		if err != nil {
			continue
		}
		if now.Sub(info.ModTime()) > 7*24*time.Hour {
			_ = os.Remove(filepath.Join(backupDir, f.Name()))
		}
	}
	return nil
}
