package repo

import (
	"testing"
	"path/filepath"
	"strings"

	"github.com/stretchr/testify/assert"
)

func TestAdoptionSecurity(t *testing.T) {
	// Simulación de la validación al adoptar una carpeta
	simulateAdoptionCheck := func(root, owner, repoName string) string {
		// La lógica defectuosa simplemente lo une y confía
		return filepath.Join(root, owner, repoName)
	}

	t.Run("Superación: Adopción legítima", func(t *testing.T) {
		root := "/ggit-repositories"
		owner := "admin"
		repo := "legacy-app.git"

		finalPath := simulateAdoptionCheck(root, owner, repo)
		assert.True(t, strings.HasPrefix(finalPath, filepath.Join(root, owner)), "Un repo legítimo se mantiene en la carpeta del dueño")
	})

	t.Run("Malintención: Adopción de carpeta del sistema mediante inyección", func(t *testing.T) {
		root := "/ggit-repositories"
		owner := "admin"
		
		// Un administrador malintencionado adopta la carpeta root del SO para leer los logs
		maliciousRepo := "../../../../../var/log/"

		finalPath := simulateAdoptionCheck(root, owner, maliciousRepo)
		
		isConfined := strings.HasPrefix(finalPath, filepath.Join(root, owner))
		assert.False(t, !isConfined, "VULNERABILIDAD CRÍTICA: La adopción de repositorios permite apoderarse de carpetas del sistema. Se debe validar estrictamente el formato del nombre [a-zA-Z0-9-_.].")
	})
}
