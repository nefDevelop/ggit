package apiclient

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveConfigRoundTrip(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv(EnvHost, "")
	t.Setenv(EnvToken, "")

	require.NoError(t, SaveConfig(Config{Host: "https://ggit.example.com", Token: "secreto"}))

	path := filepath.Join(dir, "ggit-cli", "config.json")
	info, err := os.Stat(path)
	require.NoError(t, err)
	// El archivo contiene el token: nunca debe quedar legible por otros.
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())

	cfg, err := LoadConfig()
	require.NoError(t, err)
	assert.Equal(t, "https://ggit.example.com", cfg.Host)
	assert.Equal(t, "secreto", cfg.Token)
}

func TestLoadConfigEnvOverridesFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	require.NoError(t, SaveConfig(Config{Host: "https://archivo.example.com", Token: "del-archivo"}))

	t.Setenv(EnvHost, "https://entorno.example.com")
	t.Setenv(EnvToken, "del-entorno")

	cfg, err := LoadConfig()
	require.NoError(t, err)
	assert.Equal(t, "https://entorno.example.com", cfg.Host)
	assert.Equal(t, "del-entorno", cfg.Token)
}

func TestRemoveConfig(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv(EnvHost, "")
	t.Setenv(EnvToken, "")

	require.NoError(t, SaveConfig(Config{Host: "https://ggit.example.com", Token: "secreto"}))
	require.NoError(t, RemoveConfig())

	_, err := os.Stat(filepath.Join(dir, "ggit-cli", "config.json"))
	require.ErrorIs(t, err, os.ErrNotExist)

	// Borrar dos veces no es un error.
	require.NoError(t, RemoveConfig())
}

func TestLoadConfigWithoutFile(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv(EnvHost, "https://solo-entorno.example.com")
	t.Setenv(EnvToken, "")

	cfg, err := LoadConfig()
	require.NoError(t, err)
	assert.Equal(t, "https://solo-entorno.example.com", cfg.Host)
	assert.Empty(t, cfg.Token)
}
