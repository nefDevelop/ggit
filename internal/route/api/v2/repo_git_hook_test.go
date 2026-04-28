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
}
