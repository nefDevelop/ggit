package database

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebhookSignatureLogic(t *testing.T) {
	tests := []struct {
		name     string
		secret   string
		payload  string
		expected string // HMAC-SHA256 hex
	}{
		{
			name:    "Simple payload",
			secret:  "test-secret",
			payload: `{"test":"data"}`,
			// Actual SHA256 HMAC for this secret and payload
			expected: "e914119c41a12d9a06ed4306f9f3c1ce32ee6618cbb06f54f00a280ba2dc7985",
		},
		{
			name:     "Empty secret",
			secret:   "",
			payload:  `{"test":"data"}`,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var signature string
			if len(tt.secret) > 0 {
				sig := hmac.New(sha256.New, []byte(tt.secret))
				_, _ = sig.Write([]byte(tt.payload))
				signature = hex.EncodeToString(sig.Sum(nil))
			}

			if tt.secret == "" {
				assert.Empty(t, signature)
			} else {
				assert.Equal(t, tt.expected, signature)
			}
		})
	}
}

func TestWebhookSignatureVerification(t *testing.T) {
	secret := "my-secure-webhook-secret"
	validPayload := `{"action": "push", "commits": []}`

	// Generar firma válida real para el test
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(validPayload))
	validSignature := hex.EncodeToString(mac.Sum(nil))

	verifySignature := func(secret, payload, signature string) bool {
		expectedMac := hmac.New(sha256.New, []byte(secret))
		expectedMac.Write([]byte(payload))
		expectedSignature := hex.EncodeToString(expectedMac.Sum(nil))

		// Usamos ConstantTimeCompare para evitar ataques de timing (Timing Attacks)
		return hmac.Equal([]byte(signature), []byte(expectedSignature))
	}

	t.Run("Superación: Firma Válida", func(t *testing.T) {
		assert.True(t, verifySignature(secret, validPayload, validSignature), "La firma válida debe ser aceptada")
	})

	t.Run("Malintención: Payload alterado (Man-in-the-Middle)", func(t *testing.T) {
		// El atacante intercepta la petición y cambia el payload pero mantiene la firma original
		alteredPayload := `{"action": "delete_repository", "commits": []}`
		assert.False(t, verifySignature(secret, alteredPayload, validSignature), "Debe rechazar payloads alterados que no coinciden con la firma original")
	})

	t.Run("Malintención: Firma falsificada", func(t *testing.T) {
		// El atacante intenta adivinar o forzar una firma falsa
		fakeSignature := "f914119c41a12d9a06ed4306f9f3c1ce32ee6618cbb06f54f00a280ba2dc7985"
		assert.False(t, verifySignature(secret, validPayload, fakeSignature), "Debe rechazar firmas falsificadas")
	})

	t.Run("Error: Secreto vacío o nulo", func(t *testing.T) {
		assert.False(t, verifySignature("", validPayload, validSignature), "No debe validar si no hay secreto configurado en el backend")
	})
}
