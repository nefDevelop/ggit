package database

import (
	"context"
	"regexp"

	"gorm.io/gorm"

	"github.com/nef734/ggit/internal/domain"
)

// MatchPattern is a helper to check if a string matches a regex pattern.
func MatchPattern(pattern, content string) bool {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return re.MatchString(content)
}

// SecretPattern represents a regex pattern to detect secrets in code.
type SecretPattern struct {
	ID          int64  `gorm:"primaryKey"`
	Name        string `gorm:"type:VARCHAR(255);not null"`
	Regex       string `gorm:"type:TEXT;not null"`
	Description string `gorm:"type:TEXT"`
	IsGlobal    bool   `gorm:"not null;default:FALSE"`
	RepoID      int64  `gorm:"index"` // 0 for global patterns
}

func (p *SecretPattern) ToDomain() *domain.SecretPattern {
	return &domain.SecretPattern{
		ID:          p.ID,
		Name:        p.Name,
		Regex:       p.Regex,
		Description: p.Description,
		IsGlobal:    p.IsGlobal,
		RepoID:      p.RepoID,
	}
}

func fromDomain(p *domain.SecretPattern) *SecretPattern {
	return &SecretPattern{
		ID:          p.ID,
		Name:        p.Name,
		Regex:       p.Regex,
		Description: p.Description,
		IsGlobal:    p.IsGlobal,
		RepoID:      p.RepoID,
	}
}

type SecretPatternsStore struct {
	db *gorm.DB
}

func newSecretPatternsStore(db *gorm.DB) *SecretPatternsStore {
	return &SecretPatternsStore{db: db}
}

func (s *SecretPatternsStore) Create(ctx context.Context, p *domain.SecretPattern) error {
	m := fromDomain(p)
	err := s.db.WithContext(ctx).Create(m).Error
	if err == nil {
		p.ID = m.ID
	}
	return err
}

func (s *SecretPatternsStore) ListByRepo(ctx context.Context, repoID int64) ([]*domain.SecretPattern, error) {
	var models []*SecretPattern
	err := s.db.WithContext(ctx).Where("repo_id = ? OR is_global = ?", repoID, true).Find(&models).Error
	if err != nil {
		return nil, err
	}

	patterns := make([]*domain.SecretPattern, len(models))
	for i, m := range models {
		patterns[i] = m.ToDomain()
	}
	return patterns, err
}

func (s *SecretPatternsStore) ListGlobal(ctx context.Context) ([]*domain.SecretPattern, error) {
	var models []*SecretPattern
	err := s.db.WithContext(ctx).Where("is_global = ?", true).Find(&models).Error
	if err != nil {
		return nil, err
	}

	patterns := make([]*domain.SecretPattern, len(models))
	for i, m := range models {
		patterns[i] = m.ToDomain()
	}
	return patterns, err
}

func (s *SecretPatternsStore) DeleteByID(ctx context.Context, id int64) error {
	return s.db.WithContext(ctx).Delete(&SecretPattern{}, id).Error
}

func (s *SecretPatternsStore) InitDefaultPatterns(ctx context.Context) error {
	defaults := []domain.SecretPattern{
		{
			Name:        "AWS Access Key",
			Regex:       `(A3T[A-Z0-9]|AKIA|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16}`,
			Description: "Amazon Web Services Access Key ID",
			IsGlobal:    true,
		},
		{
			Name:        "Private Key",
			Regex:       `-----BEGIN (RSA|OPENSSH|DSA|EC|PGP) PRIVATE KEY-----`,
			Description: "Archivos de clave privada SSH/PGP",
			IsGlobal:    true,
		},
		{
			Name:        "GitHub Token",
			Regex:       `gh[pousr]_[a-zA-Z0-9]{36}`,
			Description: "GitHub Personal Access Token",
			IsGlobal:    true,
		},
		{
			Name:        "Google Cloud API Key",
			Regex:       `AIza[0-9A-Za-z\\-_]{35}`,
			Description: "Google Cloud Platform API Key",
			IsGlobal:    true,
		},
		{
			Name:        "Firebase Secret",
			Regex:       `firebase-adminsdk-[a-z0-9]+-[a-z0-9]+\.json`,
			Description: "Firebase Admin SDK credentials file",
			IsGlobal:    true,
		},
		{
			Name:        ".env File",
			Regex:       `\.env$`,
			Description: "Archivos de configuración de entorno (secrets)",
			IsGlobal:    true,
		},
		{
			Name:        "Generic Secret",
			Regex:       `(?i)(api_key|apikey|secret|password|auth_token|token)(.{0,20})?['\"][0-9a-zA-Z]{16,128}['\"]`,
			Description: "Patrón genérico de secretos en código",
			IsGlobal:    true,
		},
	}

	for _, p := range defaults {
		var count int64
		s.db.WithContext(ctx).Model(&SecretPattern{}).Where("name = ? AND is_global = ?", p.Name, true).Count(&count)
		if count == 0 {
			if err := s.Create(ctx, &p); err != nil {
				return err
			}
		}
	}
	return nil
}

func (db *DB) SecretPatterns() domain.SecretPatternsStore {
	return newSecretPatternsStore(db.db)
}

