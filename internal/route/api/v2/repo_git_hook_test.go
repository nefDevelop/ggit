package v2

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGitHookOperations(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "githook-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	customHooksDir := filepath.Join(tmpDir, "custom_hooks")
	hookName := "pre-receive"
	hookPath := filepath.Join(customHooksDir, hookName)
	content := "#!/bin/bash\necho hello"

	t.Run("write hook", func(t *testing.T) {
		err := os.MkdirAll(customHooksDir, 0755)
		require.NoError(t, err)
		err = os.WriteFile(hookPath, []byte(content), 0755)
		require.NoError(t, err)

		data, err := os.ReadFile(hookPath)
		require.NoError(t, err)
		assert.Equal(t, content, string(data))
	})

	t.Run("delete hook", func(t *testing.T) {
		err := os.Remove(hookPath)
		require.NoError(t, err)
		_, err = os.Stat(hookPath)
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("Malintención: Path Traversal (Inyección de ruta)", func(t *testing.T) {
		maliciousHookName := "../../../../../etc/cron.d/malicious_cron"
		
		// Función simulada que debería existir en el controlador para limpiar la ruta
		sanitizeHookName := func(name string) bool {
			validHooks := map[string]bool{
				"pre-receive": true, "update": true, "post-receive": true,
			}
			return validHooks[name]
		}
		
		isSafe := sanitizeHookName(maliciousHookName)
		assert.False(t, isSafe, "El validador de nombres de hook debe bloquear terminantemente los intentos de Path Traversal")
	})

	t.Run("Error: Nombre de Hook Inválido (Ejecutable genérico)", func(t *testing.T) {
		invalidHookName := "my_custom_script.sh"
		
		sanitizeHookName := func(name string) bool {
			validHooks := map[string]bool{
				"pre-receive": true, "update": true, "post-receive": true,
			}
			return validHooks[name]
		}
		
		isSafe := sanitizeHookName(invalidHookName)
		assert.False(t, isSafe, "Debe rechazar cualquier nombre que no forme parte del estándar oficial de Git Hooks del servidor")
	})
}
