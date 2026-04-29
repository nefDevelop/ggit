package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDanglingRepositoriesCleanup(t *testing.T) {
	// Simulamos el comportamiento del limpiador expuesto a fallos de infraestructura
	simulateCleanupDecision := func(repoPath string, diskMounted bool) bool {
		// Mock de osx.IsDir (si el disco no está montado, simula que la carpeta no existe)
		isDir := func(path string) bool {
			if !diskMounted {
				return false
			}
			// En un disco montado, verificamos si existe
			_, err := os.Stat(path)
			return err == nil
		}

		shouldDelete := !isDir(repoPath)
		return shouldDelete
	}

	t.Run("Superación: Repositorio activo (Existe en disco)", func(t *testing.T) {
		// Creamos un directorio temporal real
		tmpDir, _ := os.MkdirTemp("", "active-repo")
		defer os.RemoveAll(tmpDir)

		shouldDelete := simulateCleanupDecision(tmpDir, true)
		assert.False(t, shouldDelete, "Un repositorio que existe físicamente JAMÁS debe ser marcado para borrado")
	})

	t.Run("Superación: Repositorio huérfano (No existe en disco sano)", func(t *testing.T) {
		missingPath := "/tmp/this-repo-was-deleted-manually-by-admin.git"
		
		shouldDelete := simulateCleanupDecision(missingPath, true)
		assert.True(t, shouldDelete, "Si el disco está sano y la carpeta no está, el limpiador debe borrar el registro huérfano")
	})

	t.Run("Error Catastrófico: Falso Positivo por fallo de disco de red (NFS/EFS)", func(t *testing.T) {
		// El repositorio es válido y crítico para la empresa
		criticalRepoPath := "/mnt/nfs_storage/ggit-repositories/admin/critical-app.git"
		
		// ERROR: El disco de red se desconecta temporalmente (diskMounted = false)
		shouldDelete := simulateCleanupDecision(criticalRepoPath, false)
		
		// El test demuestra que la lógica actual CAUSARÍA PÉRDIDA DE DATOS.
		// En una implementación segura, el limpiador debe comprobar si la raíz del almacenamiento está disponible
		// antes de procesar el borrado de repositorios individuales.
		assert.True(t, shouldDelete, "VULNERABILIDAD: El sistema actual borrará el repositorio de la DB si el disco de red falla temporalmente, causando pérdida de metadatos (issues, PRs, wiki)")
	})
}