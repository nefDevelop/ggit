package database

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/domain"
)

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
	ID        int64  `gorm:"primaryKey"`
	GistID    int64  `gorm:"index"`
	Filename  string `gorm:"size:255"`
	Language  string `gorm:"size:50"`
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
			return errors.Wrap(err, "create gist directory")
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
				return errors.Wrap(err, "write gist file content")
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
		for _, file := range gist.Files {
			// Find latest version
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
			// Rolling history logic
			if newVer > conf.Gist.MaxVersions {
				// Shift versions? No, simpler: just delete v1, rename v2->v1, etc.
				// OR just use a circular index.
				// For now, let's keep it simple: if > Max, delete v1 and shift all others down.
				// Actually, a simpler way is to just delete the oldest one and use newVer = Max.
				// But that changes the order.
				// Let's do shifting for now.
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
