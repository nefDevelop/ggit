package lfsx

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLFS_Integration(t *testing.T) {
	// Setup temporary storage
	tmpDir, err := os.MkdirTemp("", "ggit-lfs-integration-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	storage := &LocalStorage{
		Root:    filepath.Join(tmpDir, "lfs-objects"),
		TempDir: filepath.Join(tmpDir, "tmp"),
	}

	content := []byte("esto es un archivo grande simulado para LFS")
	oid := OID("9d51e4e7230efc556d63747b03dd1e288bd30c78396c4e0671b1dc20e3aefcf9") // SHA256 manual

	t.Run("Upload and Verify Storage Path", func(t *testing.T) {
		rc := io.NopCloser(bytes.NewReader(content))
		written, err := storage.Upload(oid, rc)
		require.NoError(t, err)
		assert.Equal(t, int64(len(content)), written)

		// Verificar que el archivo está en la ruta esperada (sharding de 2 niveles)
		// 9/d/9d51e4e7...
		expectedPath := filepath.Join(storage.Root, string(oid[0]), string(oid[1]), string(oid))
		assert.FileExists(t, expectedPath)
	})

	t.Run("Download and Verify Content", func(t *testing.T) {
		var buf bytes.Buffer
		err := storage.Download(oid, &buf)
		require.NoError(t, err)
		assert.Equal(t, content, buf.Bytes())
	})

	t.Run("Error on OID Mismatch", func(t *testing.T) {
		wrongOID := OID("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
		rc := io.NopCloser(bytes.NewReader(content))
		_, err := storage.Upload(wrongOID, rc)
		assert.ErrorIs(t, err, ErrOIDMismatch)
	})

	t.Run("Error on Download Non-existent", func(t *testing.T) {
		missingOID := OID("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
		err := storage.Download(missingOID, io.Discard)
		assert.ErrorIs(t, err, ErrObjectNotExist)
	})
}
