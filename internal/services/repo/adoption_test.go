package repo

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/nef734/ggit/internal/database"
	"github.com/stretchr/testify/assert"
)

func TestAdoptionSecurity(t *testing.T) {
	// Simulación de la validación al adoptar una carpeta usando la lógica real
	simulateAdoptionCheck := func(root, owner, repoName string) (string, bool) {
		if !database.IsValidRepoName(repoName) {
			return "", false
		}
		return filepath.Join(root, owner, repoName), true
	}

	t.Run("Superación: Adopción legítima", func(t *testing.T) {
		root := "/ggit-repositories"
		owner := "admin"
		repo := "legacy-app"

		finalPath, ok := simulateAdoptionCheck(root, owner, repo)
		assert.True(t, ok, "Un repo legítimo debe ser aceptado")
		assert.True(t, strings.HasPrefix(finalPath, filepath.Join(root, owner)), "Un repo legítimo se mantiene en la carpeta del dueño")
	})

	t.Run("Malintención: Adopción de carpeta del sistema mediante inyección", func(t *testing.T) {
		root := "/ggit-repositories"
		owner := "admin"

		// Un administrador malintencionado intenta adoptar una carpeta sensible
		maliciousRepos := []string{
			"../../../../../var/log",
			"/etc/passwd",
			".ssh",
			"../config",
		}

		for _, mal := range maliciousRepos {
			_, ok := simulateAdoptionCheck(root, owner, mal)
			assert.False(t, ok, "VULNERABILIDAD CRÍTICA: Se aceptó un nombre de repo malicioso: %s", mal)
		}
	})
}
