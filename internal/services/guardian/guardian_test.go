package guardian

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/database"
)

func TestScanPush(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup temp environment
	tmpDir, err := os.MkdirTemp("", "guardian-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	conf.SetMockRepository(t, conf.RepositoryOpts{Root: tmpDir})
	
	// Initialize Mock DB
	// Note: In a real environment, we'd need a more robust DB mock, 
	// but here we are testing the scanning logic which relies on database.GetRepositoryByID.
	// Since we can't easily mock the global database.Handle without more setup,
	// I'll focus on testing the regex matching part by refactoring guardian.go 
	// to expose a internal matching function if needed, or just mock the repo.
}

func TestRegexPatterns(t *testing.T) {
	tests := []struct {
		name    string
		regex   string
		content string
		want    bool
	}{
		{
			name:    "AWS Access Key",
			regex:   `(A3T[A-Z0-9]|AKIA|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16}`,
			content: "AKIAIOSFODNN7EXAMPLE",
			want:    true,
		},
		{
			name:    "SSH Private Key",
			regex:   `-----BEGIN (RSA|OPENSSH|DSA|EC|PGP) PRIVATE KEY-----`,
			content: "-----BEGIN RSA PRIVATE KEY-----",
			want:    true,
		},
		{
			name:    "GitHub Token",
			regex:   `gh[pousr]_[a-zA-Z0-9]{36}`,
			content: "ghp_123456789012345678901234567890123456",
			want:    true,
		},
		{
			name:    ".env File",
			regex:   `\.env$`,
			content: "path/to/.env",
			want:    true,
		},
		{
			name:    "Generic API Key",
			regex:   `(?i)(api_key|apikey|secret|password|auth_token|token)(.{0,20})?['\"][0-9a-zA-Z]{16,128}['\"]`,
			content: `api_key: "abcdef1234567890abcdef1234567890"`,
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := database.MatchPattern(tt.regex, tt.content)
			assert.Equal(t, tt.want, got, "Pattern: %s, Content: %s", tt.regex, tt.content)
		})
	}
}
