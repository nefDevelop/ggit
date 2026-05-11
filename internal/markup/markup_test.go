package markup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHighlight(t *testing.T) {
	t.Run("Exito: Devolver código verbatim para el frontend", func(t *testing.T) {
		code := `package main
func main() { println("hello") }`
		html, err := Highlight(code, ".go", "github")
		assert.NoError(t, err)
		assert.Equal(t, code, html, "El backend debe devolver el código sin procesar para que el frontend realice el resaltado")
	})

	t.Run("Fallback: Lenguaje desconocido", func(t *testing.T) {
		code := `some random text`
		html, err := Highlight(code, ".unknown_ext_999", "github")
		assert.NoError(t, err)
		assert.Equal(t, code, html, "Debe devolver el contenido original incluso si no reconoce el lenguaje")
	})

	t.Run("Seguridad: Integridad del contenido", func(t *testing.T) {
		// En esta arquitectura, el backend no escapa porque el frontend necesita el código puro
		// La seguridad (XSS) es responsabilidad del componente de renderizado en el frontend
		code := `<script>alert("xss")</script>`
		html, err := Highlight(code, ".html", "github")
		assert.NoError(t, err)
		assert.Equal(t, code, html, "El backend no debe alterar el contenido, permitiendo que el frontend lo maneje de forma segura")
	})
}
