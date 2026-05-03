package repo

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/go-enry/go-enry/v2"
)

func TestLanguageDetection(t *testing.T) {
	t.Run("Superación: Identificar Go", func(t *testing.T) {
		filename := "main.go"
		lang, _ := enry.GetLanguageByExtension(filename)
		assert.Equal(t, "Go", lang)
	})

	t.Run("Superación: Identificar Svelte", func(t *testing.T) {
		filename := "App.svelte"
		lang, _ := enry.GetLanguageByExtension(filename)
		assert.Equal(t, "Svelte", lang)
	})

	t.Run("Fallo: Extensión desconocida", func(t *testing.T) {
		filename := "file.unknown_ext_999"
		lang, _ := enry.GetLanguageByExtension(filename)
		assert.Empty(t, lang, "No debe identificar lenguajes inexistentes")
	})

	t.Run("Seguridad/Filtro: Archivos de configuración y documentación", func(t *testing.T) {
		// Estos archivos NO deben sumarse a las estadísticas de lenguaje según walkTree
		assert.True(t, enry.IsConfiguration("package-lock.json"))
		assert.True(t, enry.IsDocumentation("README.md"))
		assert.True(t, enry.IsVendor("node_modules/svelte/index.js"))
	})
}
