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
