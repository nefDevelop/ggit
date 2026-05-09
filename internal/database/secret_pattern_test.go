package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecretPatterns(t *testing.T) {
	// This requires a real DB connection, which is handled in a more complex way in Ggit tests.
	// For now, I'll test the MatchPattern utility which is standalone.
}

func TestMatchPattern(t *testing.T) {
	assert.True(t, MatchPattern(`\.env$`, ".env"))
	assert.True(t, MatchPattern(`\.env$`, "configs/.env"))
	assert.False(t, MatchPattern(`\.env$`, ".env.example"))

	assert.True(t, MatchPattern(`(?i)secret`, "MY_SECRET=123"))
	assert.True(t, MatchPattern(`(?i)secret`, "my-Secret-key"))
}
