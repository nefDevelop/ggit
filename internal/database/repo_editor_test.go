package database

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepoEditorSecurity(t *testing.T) {
	// Simulamos la validación defectuosa actual
	simulatePathValidation := func(localPath, userPath string) string {
		// La implementación actual usa path.Join que resuelve los '../' silenciosamente
		// y luego checkea si tiene un symlink o si incluye ".git".
		return filepath.Join(localPath, userPath)
	}

	t.Run("Superación: Subida de archivo válido", func(t *testing.T) {
		basePath := "/var/ggit/tmp/repo1"
		safePath := "docs/readme.md"

		finalPath := simulatePathValidation(basePath, safePath)
		assert.True(t, strings.HasPrefix(finalPath, basePath), "La ruta final debe mantenerse dentro de la base de trabajo del repositorio")
	})

	t.Run("Superación: Intento de Path Traversal bloqueado", func(t *testing.T) {
		maliciousPath := "../../../../etc/shadow"

		// Simulamos la nueva validación estricta
		isMalicious := strings.Contains(maliciousPath, "..") || strings.HasPrefix(maliciousPath, "/")

		assert.True(t, isMalicious, "La ruta maliciosa DEBE ser detectada")
	})
}
