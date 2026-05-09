package repo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLanguageDetection(t *testing.T) {
	t.Run("Superación: Identificar Go", func(t *testing.T) {
		filename := "main.go"
		lang := getLanguageByExtension(filename)
		assert.Equal(t, "Go", lang)
	})

	t.Run("Superación: Identificar Svelte", func(t *testing.T) {
		filename := "App.svelte"
		lang := getLanguageByExtension(filename)
		assert.Equal(t, "Svelte", lang)
	})

	t.Run("Fallo: Extensión desconocida", func(t *testing.T) {
		filename := "file.unknown_ext_999"
		lang := getLanguageByExtension(filename)
		assert.Empty(t, lang, "No debe identificar lenguajes inexistentes")
	})

	t.Run("Seguridad/Filtro: Archivos de configuración y documentación", func(t *testing.T) {
		// Estos archivos NO deben sumarse a las estadísticas de lenguaje según walkTree
		assert.True(t, isConfiguration("package-lock.json"))
		assert.True(t, isDocumentation("README.md"))
		assert.True(t, isVendor("node_modules/svelte/index.js"))
	})
}
