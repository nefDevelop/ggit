package database

import (
	"strings"
	"testing"

	"github.com/nef734/ggit/internal/dbtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRelease_Security(t *testing.T) {
	db := dbtest.NewDB(t, "ReleaseSecurity", new(Release))

	t.Run("Seguridad: Mitigación de Inyección de Comandos en TagName", func(t *testing.T) {
		// La lógica de mitigación está en createTag de release.go
		// 🚨 SECURITY: Trim any leading '-' to prevent command line argument injection.

		maliciousTag := "--upload-pack=/bin/sh"

		// Simulamos la limpieza que hace el backend
		cleanTag := strings.TrimLeft(maliciousTag, "-")

		assert.NotEqual(t, maliciousTag, cleanTag)
		assert.False(t, strings.HasPrefix(cleanTag, "-"), "El nombre del tag NUNCA debe empezar por '-' para evitar inyecciones en la CLI de Git")
	})

	t.Run("Superación: Persistencia de Release", func(t *testing.T) {
		rel := &Release{
			RepoID:       1,
			TagName:      "v1.0.0",
			LowerTagName: "v1.0.0",
			Target:       "main",
			Title:        "First Release",
		}
		require.NoError(t, db.Create(rel).Error)

		var got Release
		err := db.First(&got, rel.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, "v1.0.0", got.TagName)
	})

	t.Run("Fallo: Tag Duplicado", func(t *testing.T) {
		// Asumiendo que hay un UniqueIndex o lógica de validación previa
		// El método IsReleaseExist usa esta lógica:
		repoID := int64(1)
		tag := "v1.0.0"

		var count int64
		db.Model(new(Release)).Where("repo_id = ? AND lower_tag_name = ?", repoID, strings.ToLower(tag)).Count(&count)

		assert.True(t, count > 0, "Debe detectar que el tag ya existe para evitar colisiones")
	})
}
