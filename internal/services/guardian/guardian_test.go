package guardian

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/database"
)

func TestScanPush(t *testing.T) {
	// Simulación de la lógica de intercepción del Guardian durante un pre-receive hook
	simulateGuardianCheck := func(content string) error {
		// Evaluamos contra los patrones críticos
		patterns := []string{
			`(A3T[A-Z0-9]|AKIA|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16}`,
			`-----BEGIN (RSA|OPENSSH|DSA|EC|PGP) PRIVATE KEY-----`,
			`gh[pousr]_[a-zA-Z0-9]{36}`,
		}
		for _, p := range patterns {
			if database.MatchPattern(p, content) {
				return os.ErrPermission // Simula el rechazo del push
			}
		}
		return nil
	}

	t.Run("Superación: Push de código limpio permitido", func(t *testing.T) {
		cleanCode := `
			package main
			import "fmt"
			func main() { fmt.Println("Hello Ggit!") }
		`
		err := simulateGuardianCheck(cleanCode)
		assert.NoError(t, err, "El código sin secretos debería ser aceptado")
	})

	t.Run("Malintención: Intento de subir clave privada SSH (Bloqueo)", func(t *testing.T) {
		maliciousCode := `
			// Backup of my key
			-----BEGIN OPENSSH PRIVATE KEY-----
			b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
			-----END OPENSSH PRIVATE KEY-----
		`
		err := simulateGuardianCheck(maliciousCode)
		assert.ErrorIs(t, err, os.ErrPermission, "El push debe ser rechazado por contener una clave SSH")
	})

	t.Run("Malintención: Intento de filtrar token de AWS", func(t *testing.T) {
		maliciousCode := `const awsAccess = "AKIAIOSFODNN7EXAMPLE";`
		err := simulateGuardianCheck(maliciousCode)
		assert.ErrorIs(t, err, os.ErrPermission, "El push debe ser rechazado por contener un token AWS")
	})
}

func TestRegexPatterns(t *testing.T) {
	tests := []struct {
		name    string
		regex   string
		content string
		want    bool
	}{
		{
			name:    "AWS Access Key",
			regex:   `(A3T[A-Z0-9]|AKIA|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16}`,
			content: "AKIAIOSFODNN7EXAMPLE",
			want:    true,
		},
		{
			name:    "SSH Private Key",
			regex:   `-----BEGIN (RSA|OPENSSH|DSA|EC|PGP) PRIVATE KEY-----`,
			content: "-----BEGIN RSA PRIVATE KEY-----",
			want:    true,
		},
		{
			name:    "GitHub Token",
			regex:   `gh[pousr]_[a-zA-Z0-9]{36}`,
			content: "ghp_123456789012345678901234567890123456",
			want:    true,
		},
		{
			name:    ".env File",
			regex:   `\.env$`,
			content: "path/to/.env",
			want:    true,
		},
		{
			name:    "Generic API Key",
			regex:   `(?i)(api_key|apikey|secret|password|auth_token|token)(.{0,20})?['\"][0-9a-zA-Z]{16,128}['\"]`,
			content: `api_key: "abcdef1234567890abcdef1234567890"`,
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := database.MatchPattern(tt.regex, tt.content)
			assert.Equal(t, tt.want, got, "Pattern: %s, Content: %s", tt.regex, tt.content)
		})
	}
}
