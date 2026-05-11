package repox

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nef734/ggit/internal/conf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHooks_Integration(t *testing.T) {
	// Setup temporary repo environment
	tmpDir, err := os.MkdirTemp("", "ggit-hooks-integration-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	repoPath := filepath.Join(tmpDir, "test-repo.git")
	err = os.MkdirAll(filepath.Join(repoPath, "hooks"), 0755)
	require.NoError(t, err)

	// Mock configuration for hooks
	conf.App.RunUser = "git"
	conf.Repository.ScriptType = "bash"
	conf.CustomConf = "/etc/ggit/app.ini"

	t.Run("SelfHealHooks creates correct hook files", func(t *testing.T) {
		err := SelfHealHooks(repoPath)
		require.NoError(t, err)

		hooks := []string{"pre-receive", "update", "post-receive"}
		for _, name := range hooks {
			path := filepath.Join(repoPath, "hooks", name)
			assert.FileExists(t, path)

			content, err := os.ReadFile(path)
			require.NoError(t, err)

			// Verificar que el contenido delega al binario de ggit con la config correcta
			assert.Contains(t, string(content), "#!/usr/bin/env bash")
			assert.Contains(t, string(content), "hook --config='/etc/ggit/app.ini' "+name)
		}
	})

	t.Run("SelfHealHooks repairs corrupted hooks", func(t *testing.T) {
		hookPath := filepath.Join(repoPath, "hooks", "update")
		
		// Corromper el archivo
		err := os.WriteFile(hookPath, []byte("contenido basura"), 0644)
		require.NoError(t, err)

		// Ejecutar sanación
		err = SelfHealHooks(repoPath)
		require.NoError(t, err)

		content, _ := os.ReadFile(hookPath)
		assert.Contains(t, string(content), "update $1 $2 $3", "El hook 'update' debería haber sido reparado")
	})
}
