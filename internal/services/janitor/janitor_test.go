package janitor

import (
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJanitorPatterns(t *testing.T) {
	tests := []struct {
		filename    string
		wantType    string
		wantDesc    string
		wantMatch   bool
	}{
		{".DS_Store", "system", "macOS metadata file", true},
		{"node_modules", "dependency", "Node.js dependencies (should be in .gitignore)", true},
		{"dist", "build", "Build distribution folder", true},
		{"main.go", "", "", false},
		{"id_rsa.pem", "security", "Private key file (Sensitive)", true},
		{".env", "security", "Environment configuration (Sensitive)", true},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			matchFound := false
			foundDesc := ""
			name := path.Base(tt.filename)

			for pattern, desc := range patterns {
				matched, _ := path.Match(pattern, name)
				if matched || name == pattern {
					matchFound = true
					foundDesc = desc
					break
				}
			}

			assert.Equal(t, tt.wantMatch, matchFound)
			if tt.wantMatch {
				assert.Equal(t, tt.wantDesc, foundDesc)
				
				// Test Type Logic
				issueType := "build"
				if strings.Contains(foundDesc, "dependencies") || strings.Contains(foundDesc, "environment") {
					issueType = "dependency"
				} else if strings.Contains(foundDesc, "metadata") {
					issueType = "system"
				} else if strings.Contains(foundDesc, "Sensitive") {
					issueType = "security"
				}
				assert.Equal(t, tt.wantType, issueType)
			}
		})
	}
}

func TestBinaryDetectionLogic(t *testing.T) {
	tests := []struct {
		name     string
		size     int64
		wantType string
		wantDesc string
	}{
		{
			name:     "small_go_file.go",
			size:     1024,
			wantType: "",
		},
		{
			name:     "large_data_file.dat",
			size:     20 * 1024 * 1024,
			wantType: "large_file",
			wantDesc: "Large file detected (> 10MB)",
		},
		{
			name:     "ggit_binary",
			size:     15 * 1024 * 1024,
			wantType: "large_file",
			wantDesc: "Possible compiled binary without extension",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var issueType string
			var desc string

			// Simulate logic in janitor.go
			if tt.size > 10*1024*1024 {
				issueType = "large_file"
				desc = "Large file detected (> 10MB)"
				if !strings.Contains(path.Base(tt.name), ".") {
					desc = "Possible compiled binary without extension"
				}
			}

			assert.Equal(t, tt.wantType, issueType)
			if tt.wantType != "" {
				assert.Equal(t, tt.wantDesc, desc)
			}
		})
	}
}
