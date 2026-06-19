package database

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/domain"
	"gorm.io/gorm"
)

type GistUser struct {
	ID     int64
	Name   string
	Avatar string
}

// Gist is the database model for domain.Gist.
type Gist struct {
	ID          int64  `gorm:"primaryKey"`
	UUID        string `gorm:"uniqueIndex;size:36"`
	UserID      int64  `gorm:"index"`
	Description string `gorm:"size:255"`
	IsPublic    bool
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

// GistFile is the database model for domain.GistFile.
// We store filename and other metadata in DB, but content is usually in FS.
// For the hybrid model, we store "CurrentContent" in a separate table or just in FS.
type GistFile struct {
	ID        int64     `gorm:"primaryKey"`
	GistID    int64     `gorm:"index"`
	Filename  string    `gorm:"size:255"`
	Language  string    `gorm:"size:50"`
	CreatedAt time.Time `gorm:"not null"`
}

func (g *Gist) BeforeCreate(tx *gorm.DB) error {
	if g.UUID == "" {
		g.UUID = uuid.New().String()
	}
	return nil
}

type gistsStore struct {
	db *gorm.DB
}

func newGistsStore(db *gorm.DB) domain.GistStore {
	return &gistsStore{db: db}
}

func (s *gistsStore) gistDir(uuid string) string {
	return filepath.Join(conf.Server.AppDataPath, "gists", uuid)
}

func (s *gistsStore) Create(gist *domain.Gist) error {
	g := &Gist{
		UserID:      gist.UserID,
		Description: gist.Description,
		IsPublic:    gist.IsPublic,
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(g).Error; err != nil {
			return err
		}
		gist.ID = g.ID
		gist.UUID = g.UUID

		dir := s.gistDir(g.UUID)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create gist directory"+": %v", err)
		}

		for _, file := range gist.Files {
			gf := &GistFile{
				GistID:   g.ID,
				Filename: file.Filename,
				Language: file.Language,
			}
			if err := tx.Create(gf).Error; err != nil {
				return err
			}
			file.ID = gf.ID

			// Write content to v1.md
			filePath := filepath.Join(dir, fmt.Sprintf("%s.v1", gf.Filename))
			if err := os.WriteFile(filePath, []byte(file.Content), 0644); err != nil {
				return fmt.Errorf("write gist file content"+": %v", err)
			}
		}
		return nil
	})
}

func (s *gistsStore) GetByUUID(uuid string) (*domain.Gist, error) {
	var g Gist
	if err := s.db.Where("uuid = ?", uuid).First(&g).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrGistNotFound{UUID: uuid}
		}
		return nil, err
	}

	var user User
	if err := s.db.Where("id = ?", g.UserID).First(&user).Error; err != nil {
		return nil, err
	}

	var files []*GistFile
	if err := s.db.Where("gist_id = ?", g.ID).Find(&files).Error; err != nil {
		return nil, err
	}

	dir := s.gistDir(g.UUID)
	domainFiles := make([]*domain.GistFile, len(files))
	for i, f := range files {
		// Read current version (for now let's assume it's the latest one found in FS)
		// Better: store "CurrentVersion" in GistFile DB record.
		// For simplicity in this first draft, I'll just read the latest version file.
		content, err := s.readLatestVersion(dir, f.Filename)
		if err != nil {
			return nil, err
		}

		domainFiles[i] = &domain.GistFile{
			ID:       f.ID,
			GistID:   f.GistID,
			Filename: f.Filename,
			Language: f.Language,
			Content:  content,
		}
	}

	return &domain.Gist{
		ID:          g.ID,
		UUID:        g.UUID,
		UserID:      g.UserID,
		Description: g.Description,
		IsPublic:    g.IsPublic,
		CreatedAt:   g.CreatedAt,
		UpdatedAt:   g.UpdatedAt,
		Files:       domainFiles,
		Owner: &domain.User{
			ID:     user.ID,
			Name:   user.Name,
			Avatar: user.Avatar,
		},
	}, nil
}

func (s *gistsStore) readLatestVersion(dir, filename string) (string, error) {
	// Find latest .vX file
	var latestContent string
	var latestVer int
	for v := 1; v <= conf.Gist.MaxVersions; v++ {
		path := filepath.Join(dir, fmt.Sprintf("%s.v%d", filename, v))
		if _, err := os.Stat(path); err == nil {
			latestVer = v
		} else {
			break
		}
	}

	if latestVer == 0 {
		return "", errors.New("gist file not found on disk")
	}

	path := filepath.Join(dir, fmt.Sprintf("%s.v%d", filename, latestVer))
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	latestContent = string(content)
	return latestContent, nil
}

func (s *gistsStore) ListByUserID(userID int64, isOwner bool) ([]*domain.Gist, error) {
	var gists []*Gist
	query := s.db.Where("user_id = ?", userID)
	if !isOwner {
		query = query.Where("is_public = ?", true)
	}
	if err := query.Order("updated_at DESC").Find(&gists).Error; err != nil {
		return nil, err
	}

	domainGists := make([]*domain.Gist, len(gists))
	for i, g := range gists {
		domainGists[i] = &domain.Gist{
			ID:          g.ID,
			UUID:        g.UUID,
			UserID:      g.UserID,
			Description: g.Description,
			IsPublic:    g.IsPublic,
			CreatedAt:   g.CreatedAt,
			UpdatedAt:   g.UpdatedAt,
		}
	}
	return domainGists, nil
}

func (s *gistsStore) ListPublic() ([]*domain.Gist, error) {
	var gists []*Gist
	if err := s.db.Where("is_public = ?", true).Order("updated_at DESC").Limit(50).Find(&gists).Error; err != nil {
		return nil, err
	}

	domainGists := make([]*domain.Gist, len(gists))
	for i, g := range gists {
		domainGists[i] = &domain.Gist{
			ID:          g.ID,
			UUID:        g.UUID,
			UserID:      g.UserID,
			Description: g.Description,
			IsPublic:    g.IsPublic,
			CreatedAt:   g.CreatedAt,
			UpdatedAt:   g.UpdatedAt,
		}
	}
	return domainGists, nil
}

func (s *gistsStore) Update(gist *domain.Gist) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Gist{}).Where("id = ?", gist.ID).Updates(map[string]any{
			"description": gist.Description,
			"is_public":   gist.IsPublic,
		}).Error; err != nil {
			return err
		}

		dir := s.gistDir(gist.UUID)

		// Track IDs of files that are still present
		var currentFileIDs []int64
		for _, file := range gist.Files {
			var gf GistFile
			if file.ID > 0 {
				// Update existing file record
				if err := tx.Where("id = ? AND gist_id = ?", file.ID, gist.ID).First(&gf).Error; err == nil {
					// Check if filename changed to rename files on disk?
					// For now let's just update the DB record.
					if gf.Filename != file.Filename {
						// Optional: rename existing .vX files to new filename
						for v := 1; v <= conf.Gist.MaxVersions; v++ {
							oldPath := filepath.Join(dir, fmt.Sprintf("%s.v%d", gf.Filename, v))
							newPath := filepath.Join(dir, fmt.Sprintf("%s.v%d", file.Filename, v))
							_ = os.Rename(oldPath, newPath)
						}
					}

					gf.Filename = file.Filename
					gf.Language = file.Language
					if err := tx.Save(&gf).Error; err != nil {
						return err
					}
					currentFileIDs = append(currentFileIDs, gf.ID)
				}
			} else {
				// Create new file record
				gf = GistFile{
					GistID:   gist.ID,
					Filename: file.Filename,
					Language: file.Language,
				}
				if err := tx.Create(&gf).Error; err != nil {
					return err
				}
				file.ID = gf.ID
				currentFileIDs = append(currentFileIDs, gf.ID)
			}

			// Write new version to disk
			// Find latest version for THIS filename
			latestVer := 0
			for v := 1; v <= conf.Gist.MaxVersions; v++ {
				path := filepath.Join(dir, fmt.Sprintf("%s.v%d", file.Filename, v))
				if _, err := os.Stat(path); err == nil {
					latestVer = v
				} else {
					break
				}
			}

			newVer := latestVer + 1
			if newVer > conf.Gist.MaxVersions {
				for v := 1; v < conf.Gist.MaxVersions; v++ {
					oldPath := filepath.Join(dir, fmt.Sprintf("%s.v%d", file.Filename, v+1))
					newPath := filepath.Join(dir, fmt.Sprintf("%s.v%d", file.Filename, v))
					_ = os.Rename(oldPath, newPath)
				}
				newVer = conf.Gist.MaxVersions
			}

			path := filepath.Join(dir, fmt.Sprintf("%s.v%d", file.Filename, newVer))
			if err := os.WriteFile(path, []byte(file.Content), 0644); err != nil {
				return err
			}
		}

		// Delete files that were removed from the Gist
		if err := tx.Where("gist_id = ? AND id NOT IN ?", gist.ID, currentFileIDs).Delete(&GistFile{}).Error; err != nil {
			// If currentFileIDs is empty, NOT IN might fail in some DBs, but GORM handles it.
			// Actually, let's be safe.
			if len(currentFileIDs) == 0 {
				if err := tx.Where("gist_id = ?", gist.ID).Delete(&GistFile{}).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (s *gistsStore) Delete(uuid string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var g Gist
		if err := tx.Where("uuid = ?", uuid).First(&g).Error; err != nil {
			return err
		}

		if err := tx.Delete(&GistFile{}, "gist_id = ?", g.ID).Error; err != nil {
			return err
		}
		if err := tx.Delete(&g).Error; err != nil {
			return err
		}

		dir := s.gistDir(uuid)
		return os.RemoveAll(dir)
	})
}
