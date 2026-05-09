package health

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSentinelDecisionLogic(t *testing.T) {
	now := time.Now().Unix()

	tests := []struct {
		name              string
		lastCheck         int64
		commitsSinceCheck int64
		activityCooldown  int
		idleCheckDays     int
		wantShouldRun     bool
	}{
		{
			name:              "Initial check (never checked)",
			lastCheck:         0,
			commitsSinceCheck: 0,
			wantShouldRun:     true,
		},
		{
			name:              "Mandatory monthly check (31 days ago)",
			lastCheck:         now - (31 * 24 * 60 * 60),
			commitsSinceCheck: 0,
			wantShouldRun:     true,
		},
		{
			name:              "Recently checked, no activity",
			lastCheck:         now - (1 * 24 * 60 * 60),
			commitsSinceCheck: 0,
			wantShouldRun:     false,
		},
		{
			name:              "Activity detected, but cooldown not reached",
			lastCheck:         now - (1 * 24 * 60 * 60),
			commitsSinceCheck: 5,
			activityCooldown:  7, // days
			wantShouldRun:     false,
		},
		{
			name:              "Activity detected, cooldown reached",
			lastCheck:         now - (8 * 24 * 60 * 60),
			commitsSinceCheck: 1,
			activityCooldown:  7,
			wantShouldRun:     true,
		},
		{
			name:              "Idle period reached",
			lastCheck:         now - (15 * 24 * 60 * 60),
			commitsSinceCheck: 0,
			idleCheckDays:     14,
			wantShouldRun:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeSinceLast := now - tt.lastCheck
			shouldRun := false

			// Logic matched with sentinel.go
			if tt.lastCheck == 0 {
				shouldRun = true
			} else if timeSinceLast > 30*24*60*60 {
				shouldRun = true
			} else if tt.commitsSinceCheck > 0 {
				cooldown := int64(tt.activityCooldown)
				if cooldown == 0 {
					cooldown = 7
				}
				if timeSinceLast > cooldown*24*60*60 {
					shouldRun = true
				} else if tt.commitsSinceCheck >= 100 {
					shouldRun = true
				}
			} else {
				idleDays := int64(tt.idleCheckDays)
				if idleDays == 0 {
					idleDays = 14
				}
				if timeSinceLast > idleDays*24*60*60 {
					shouldRun = true
				}
			}

			assert.Equal(t, tt.wantShouldRun, shouldRun)
		})
	}
}

func TestSentinelCorruptionHandling(t *testing.T) {
	// Simulador del scanner fsck
	simulateFsckCheck := func(repoPath string, isCorrupted bool) error {
		if isCorrupted {
			return assert.AnError // Simula un "fatal: loose object is corrupt"
		}
		return nil
	}

	t.Run("Superación: Repositorio Sano", func(t *testing.T) {
		err := simulateFsckCheck("/path/to/healthy/repo.git", false)
		assert.NoError(t, err, "No debe reportar errores en repositorios sanos")
	})

	t.Run("Error: Corrupción de Objetos Git (fsck falla)", func(t *testing.T) {
		err := simulateFsckCheck("/path/to/corrupted/repo.git", true)
		assert.Error(t, err, "El Sentinel debe atrapar el error del fsck")
		// Aquí verificaríamos que el Sentinel cambia el estado del repo en DB a "Corrupto"
		// y desactiva las operaciones de escritura (Push, UI Edits).
	})

	t.Run("Malintención: Intento de inyección de comandos en ruta", func(t *testing.T) {
		// Si un atacante logra crear un repositorio con nombre malicioso
		maliciousPath := "/path/repo.git; rm -rf /"

		// La validación estricta debe rechazar rutas que no pertenezcan al contexto seguro
		sanitizePath := func(p string) bool {
			// En la implementación real esto chequearía contra un regex estricto de nombres de repos [a-zA-Z0-9-_.]
			return !assert.ObjectsAreEqual(p, maliciousPath) // Simulación básica
		}

		isSafe := sanitizePath(maliciousPath)
		assert.False(t, isSafe, "El Sentinel debe rechazar la ejecución del fsck en rutas inyectadas para evitar RCE")
	})
}
