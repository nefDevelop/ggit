package repo

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nef734/ggit/internal/database"
)

// auto_init with no README named used to write nothing at all: "git add --all"
// staged nothing, "git commit" aborted with "nothing to commit" and the API
// answered 500. The web form always sent a README name, so only API clients
// hit it.
func TestAutoInitWithoutReadmeStillInitialises(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	owner := t.TempDir()

	for _, tc := range []struct {
		name       string
		opts       CreateOptions
		wantReadme bool
	}{
		{name: "auto_init without readme", opts: CreateOptions{AutoInit: true}, wantReadme: true},
		{name: "auto_init with readme", opts: CreateOptions{AutoInit: true, Readme: "Default"}, wantReadme: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repoName := "autoinit-" + strings.ReplaceAll(tc.name, " ", "-")
			repoPath := filepath.Join(owner, repoName+".git")
			// prepareRepoCommit clones, so the fixture has to be a real bare
			// repository rather than an empty directory.
			if out, err := exec.Command("git", "init", "--bare", "-q", repoPath).CombinedOutput(); err != nil {
				t.Fatalf("git init --bare: %v\n%s", err, out)
			}

			repo := &database.Repository{
				Name:  repoName,
				Owner: &database.User{ID: 1, Name: "owner", LowerName: "owner"},
			}

			tmp := filepath.Join(t.TempDir(), "work")
			if err := prepareRepoCommit(repo, tmp, repoPath, tc.opts); err != nil {
				t.Fatalf("prepareRepoCommit with %+v: %v", tc.opts, err)
			}

			if _, err := os.Stat(filepath.Join(tmp, "README.md")); err != nil {
				t.Fatalf("README.md was not written, so there is nothing to commit: %v", err)
			}
		})
	}
}
