package database

import (
	"strings"
	"testing"
)

// TestIsRepoNameAllowed covers the whitelist that keeps a repository name from
// escaping conf.Repository.Root. Repository.RepoPath joins the name onto the
// root without further sanitising, so any separator reaching disk decides which
// directory git init and RemoveAll operate on.
func TestIsRepoNameAllowed(t *testing.T) {
	tests := []struct {
		name    string
		repo    string
		wantErr bool
	}{
		{name: "simple name", repo: "ggit", wantErr: false},
		{name: "name with dashes", repo: "mi-proyecto", wantErr: false},
		{name: "name with underscores", repo: "mi_proyecto", wantErr: false},
		{name: "name with dots", repo: "mi.proyecto", wantErr: false},
		{name: "mixed case", repo: "MiProyecto", wantErr: false},
		{name: "digits", repo: "proyecto2026", wantErr: false},

		// Path traversal, either on its own or hidden inside a valid-looking name.
		{name: "parent directory", repo: "..", wantErr: true},
		{name: "current directory", repo: ".", wantErr: true},
		{name: "leading parent traversal", repo: "../otro", wantErr: true},
		{name: "trailing parent traversal", repo: "otro/..", wantErr: true},
		{name: "embedded traversal", repo: "a/../../b", wantErr: true},
		{name: "absolute path", repo: "/etc/cron.d/x", wantErr: true},
		{name: "forward separator", repo: "foo/bar", wantErr: true},
		{name: "backslash separator", repo: `foo\bar`, wantErr: true},
		{name: "traversal to absolute", repo: "../../../../srv/pwn", wantErr: true},

		// Other rejections.
		{name: "leading dot", repo: ".oculto", wantErr: true},
		{name: "empty", repo: "", wantErr: true},
		{name: "whitespace only", repo: "   ", wantErr: true},
		{name: "git suffix", repo: "repositorio.git", wantErr: true},
		{name: "wiki suffix", repo: "repositorio.wiki", wantErr: true},
		{name: "null byte", repo: "proyecto\x00", wantErr: true},
		{name: "newline", repo: "proyecto\n", wantErr: true},
		{name: "space inside", repo: "mi proyecto", wantErr: true},
		{name: "non ascii", repo: "proyecto-ñ", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := isRepoNameAllowed(tt.repo)
			if tt.wantErr && err == nil {
				t.Fatalf("isRepoNameAllowed(%q) = nil, se esperaba error", tt.repo)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("isRepoNameAllowed(%q) = %v, se esperaba nil", tt.repo, err)
			}
		})
	}
}

// TestIsRepoNameAllowedRejectsTraversalSpellings guards the specific property the
// validator exists for: no accepted name may contain a path separator or a
// parent reference, whatever the combination.
func TestIsRepoNameAllowedRejectsTraversalSpellings(t *testing.T) {
	hostile := []string{
		"..", "../", "../../", "a/../b", "./a", "a/./b",
		"...", "a..b", "a/..", `..\..`, "a\\b", "/", "//",
	}

	for _, name := range hostile {
		t.Run(name, func(t *testing.T) {
			if err := isRepoNameAllowed(name); err == nil {
				t.Fatalf("isRepoNameAllowed(%q) acepto un nombre con recorrido de ruta", name)
			}
			// The whitelist is only meaningful if it is at least as strict as the
			// predicate used elsewhere in the codebase.
			if IsValidRepoName(name) && strings.ContainsAny(name, `/\`) {
				t.Fatalf("IsValidRepoName(%q) acepta separadores de ruta", name)
			}
		})
	}
}