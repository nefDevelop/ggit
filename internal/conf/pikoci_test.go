package conf

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPikoCIDefaults(t *testing.T) {
	defer func() { MustInit("") }()
	require.NoError(t, Init(""))

	assert.False(t, PikoCI.Enabled)
	assert.Equal(t, "sidecar", PikoCI.Mode)
	assert.Equal(t, "http://127.0.0.1:18080", PikoCI.URL)
	assert.Equal(t, 2, PikoCI.Concurrency)
	assert.Empty(t, PikoCI.JWTSecret)
	assert.Empty(t, PikoCI.SecretKey)
}

func TestPikoCIMapping(t *testing.T) {
	defer func() { MustInit("") }()
	require.NoError(t, Init(filepath.Join("testdata", "pikoci.ini")))

	assert.True(t, PikoCI.Enabled)
	assert.Equal(t, "subprocess", PikoCI.Mode)
	assert.Equal(t, "http://127.0.0.1:19000", PikoCI.URL)
	assert.Equal(t, "https://ci.example.com", PikoCI.PublicURL)
	assert.Equal(t, "/usr/local/bin/pikoci", PikoCI.BinaryPath)
	assert.Equal(t, "/tmp/pikoci-data", PikoCI.DataDir)
	assert.Equal(t, 4, PikoCI.Concurrency)
	assert.Equal(t, "jwt-de-prueba", PikoCI.JWTSecret)
	assert.Equal(t, "secreto-de-prueba", PikoCI.SecretKey)
}
