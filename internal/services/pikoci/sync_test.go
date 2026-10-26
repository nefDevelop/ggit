package pikoci

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeLinks struct {
	link *database.PikoCIRepo
}

func (f *fakeLinks) GetByRepoID(_ context.Context, repoID int64) (*database.PikoCIRepo, error) {
	if f.link == nil || f.link.RepoID != repoID {
		return nil, database.ErrPikoCIRepoNotExist
	}
	return f.link, nil
}

func (f *fakeLinks) List(context.Context, int) ([]*database.PikoCIRepo, error) {
	if f.link == nil {
		return nil, nil
	}
	return []*database.PikoCIRepo{f.link}, nil
}

func (f *fakeLinks) Upsert(_ context.Context, link *database.PikoCIRepo) error {
	f.link = link
	return nil
}

func (f *fakeLinks) Delete(context.Context, int64) error {
	f.link = nil
	return nil
}

func (f *fakeLinks) Count(context.Context) (int64, error) {
	if f.link == nil {
		return 0, nil
	}
	return 1, nil
}

func gitRun(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "git %v: %s", args, out)
	return strings.TrimSpace(string(out))
}

func initTestRepo(t *testing.T, root, owner, name, pipeline string) {
	t.Helper()
	dir := filepath.Join(root, owner, name+".git")
	require.NoError(t, os.MkdirAll(filepath.Join(dir, ".ggit"), 0o755))
	gitRun(t, dir, "init", "-b", "main")
	gitRun(t, dir, "config", "user.email", "ci@example.com")
	gitRun(t, dir, "config", "user.name", "CI")
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".ggit", "ci.hcl"), []byte(pipeline), 0o644))
	gitRun(t, dir, "add", ".")
	gitRun(t, dir, "commit", "-m", "pipeline")
}

func TestSyncRepoPipeline(t *testing.T) {
	root := t.TempDir()
	oldRoot, oldURL, oldFile := conf.Repository.Root, conf.Server.ExternalURL, conf.PikoCI.PipelineFile
	oldBot, oldPublic := conf.PikoCI.BotToken, conf.PikoCI.PublicURL
	defer func() {
		conf.Repository.Root, conf.Server.ExternalURL, conf.PikoCI.PipelineFile = oldRoot, oldURL, oldFile
		conf.PikoCI.BotToken, conf.PikoCI.PublicURL = oldBot, oldPublic
	}()

	conf.Repository.Root = root
	conf.Server.ExternalURL = "https://ggit.example.com"
	conf.PikoCI.PipelineFile = ".ggit/ci.hcl"
	conf.PikoCI.BotToken = "bot-token"
	conf.PikoCI.PublicURL = "https://ci.example.com"

	const pipeline = "job \"test\" {}\n"
	initTestRepo(t, root, "acme", "web", pipeline)

	var (
		gotPath string
		gotBody map[string]any
	)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPut:
			gotPath = r.URL.Path
			require.NoError(t, json.NewDecoder(r.Body).Decode(&gotBody))
			_, _ = w.Write([]byte(`{}`))
		case strings.HasSuffix(r.URL.Path, "/webhook_token"):
			_, _ = w.Write([]byte(`{"token":"hook-token"}`))
		}
	}))
	defer server.Close()

	service := NewWithAuth(server.URL, "main", "pko_api")
	links := &fakeLinks{}
	service.links = links

	link, err := service.SyncRepoPipeline(context.Background(), 7, "acme", "web", "main")

	require.NoError(t, err)
	assert.Equal(t, "/teams/main/pipelines/acme-web", gotPath)
	assert.Equal(t, "acme-web", link.PipelineName)
	assert.Equal(t, "hook-token", link.WebhookToken)
	require.NotNil(t, links.link)
	assert.Equal(t, int64(7), links.link.RepoID)

	decoded, err := base64.StdEncoding.DecodeString(gotBody["config"].(string))
	require.NoError(t, err)
	assert.Equal(t, pipeline, string(decoded))

	vars := gotBody["vars"].(map[string]any)
	assert.Equal(t, "acme/web", vars["repo"])
	assert.Equal(t, "main", vars["branch"])
	assert.Equal(t, "bot-token", vars["ggit_token"])
	assert.Equal(t, "https://ggit.example.com", vars["ggit_url"])
	assert.Equal(t, "https://ggit.example.com/acme/web.git", vars["git_url"])
	assert.Equal(t, "https://ci.example.com", vars["pikoci_url"])
	assert.Equal(t, "acme-web", vars["serial_group"])
}

func TestSyncRepoPipelineWithoutFile(t *testing.T) {
	root := t.TempDir()
	oldRoot, oldFile := conf.Repository.Root, conf.PikoCI.PipelineFile
	defer func() {
		conf.Repository.Root, conf.PikoCI.PipelineFile = oldRoot, oldFile
	}()

	conf.Repository.Root = root
	conf.PikoCI.PipelineFile = ".ggit/ci.hcl"

	initTestRepo(t, root, "acme", "empty", "# sin pipeline\n")
	require.NoError(t, os.Remove(filepath.Join(root, "acme", "empty.git", ".ggit", "ci.hcl")))
	gitRun(t, filepath.Join(root, "acme", "empty.git"), "commit", "-am", "remove")

	service := NewWithAuth("http://127.0.0.1:1", "main", "pko")
	service.links = &fakeLinks{}

	_, err := service.SyncRepoPipeline(context.Background(), 8, "acme", "empty", "main")

	require.Error(t, err)
	assert.ErrorIs(t, err, ErrPipelineFileNotFound)
}

func TestPipelineFileChanged(t *testing.T) {
	root := t.TempDir()
	oldRoot := conf.Repository.Root
	defer func() { conf.Repository.Root = oldRoot }()
	conf.Repository.Root = root

	initTestRepo(t, root, "acme", "web", "# v1\n")
	dir := filepath.Join(root, "acme", "web.git")
	first := gitRun(t, dir, "rev-parse", "HEAD")

	require.NoError(t, os.WriteFile(filepath.Join(dir, ".ggit", "ci.hcl"), []byte("# v2\n"), 0o644))
	gitRun(t, dir, "commit", "-am", "update")
	second := gitRun(t, dir, "rev-parse", "HEAD")

	changed, err := PipelineFileChanged("acme", "web", first, second, ".ggit/ci.hcl")
	require.NoError(t, err)
	assert.True(t, changed)

	changed, err = PipelineFileChanged("acme", "web", second, second, ".ggit/ci.hcl")
	require.NoError(t, err)
	assert.False(t, changed)
}
