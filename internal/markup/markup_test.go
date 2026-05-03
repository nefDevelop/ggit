package markup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHighlight(t *testing.T) {
	t.Run("Superación: Resaltar Go", func(t *testing.T) {
		code := `package main
func main() { println("hello") }`
		html, err := Highlight(code, ".go", "github")
		assert.NoError(t, err)
		assert.Contains(t, html, "color:")
		assert.Contains(t, html, "package")
		assert.NotContains(t, html, "<pre", "No debe incluir la etiqueta <pre> para permitir procesamiento en el frontend")
	})

	t.Run("Fallo/Fallback: Lenguaje desconocido", func(t *testing.T) {
		code := `some random text`
		html, err := Highlight(code, ".unknown_ext_999", "github")
		assert.NoError(t, err)
		assert.NotEmpty(t, html, "Debe devolver contenido incluso si no reconoce el lenguaje")
	})

	t.Run("Seguridad: Mitigación de XSS", func(t *testing.T) {
		// Chroma HTML formatter debe escapar caracteres peligrosos
		code := `<script>alert("xss")</script>`
		html, err := Highlight(code, ".html", "github")
		assert.NoError(t, err)
		assert.NotContains(t, html, "<script>", "Debe escapar los tags HTML para evitar XSS")
		// Verificar que el contenido está escapado (buscamos la entidad &lt; dentro de los spans)
		assert.Contains(t, html, "&lt;", "El tag script debe estar escapado como entidad HTML")
	})
}
