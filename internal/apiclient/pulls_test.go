package apiclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListPullRequestsQuery(t *testing.T) {
	var gotPath, gotQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotQuery = r.URL.RawQuery
		_, _ = w.Write([]byte(`[{"number":4,"title":"Code review","state":"open","head_branch":"feature/x","base_branch":"main"}]`))
	}))
	defer server.Close()

	pulls, err := New(server.URL, "t").ListPullRequests(context.Background(), "pi", "ggit", ListPullRequestsOptions{
		State: "open",
		Page:  2,
	})

	require.NoError(t, err)
	require.Len(t, pulls, 1)
	assert.Equal(t, int64(4), pulls[0].Index)
	assert.Equal(t, "feature/x", pulls[0].HeadBranch)
	assert.Equal(t, "/api/v2/pi/ggit/pulls", gotPath)
	assert.Contains(t, gotQuery, "state=open")
	assert.Contains(t, gotQuery, "page=2")
}

func TestListPullReviews(t *testing.T) {
	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`[{"id":1,"type":1,"reviewer":{"username":"nef"},"comments":[{"id":9,"tree_path":"main.go","line":3}]}]`))
	}))
	defer server.Close()

	reviews, err := New(server.URL, "t").ListPullReviews(context.Background(), "pi", "ggit", 4)

	require.NoError(t, err)
	require.Len(t, reviews, 1)
	assert.Equal(t, 1, reviews[0].Type)
	assert.Equal(t, "nef", reviews[0].Reviewer.Username)
	require.Len(t, reviews[0].Comments, 1)
	assert.Equal(t, "/api/v2/pi/ggit/pulls/4/reviews", gotPath)
}

// El servidor serializa internal/gitx.Diff sin tags JSON: el cliente debe
// entender los nombres de campo de Go tal cual.
func TestGetPullRequestDiffDecodesServerShape(t *testing.T) {
	const payload = `{"Files":[{"Name":"main.go","Type":2,"Index":"abc","OldIndex":"def","Sections":[{"Lines":[
		{"Type":4,"Content":"@@ -1,3 +1,4 @@","LeftLine":0,"RightLine":0},
		{"Type":1,"Content":" package main","LeftLine":1,"RightLine":1},
		{"Type":3,"Content":"-// viejo","LeftLine":2,"RightLine":0},
		{"Type":2,"Content":"+// nuevo","LeftLine":0,"RightLine":2}
	]}]}]}`

	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(payload))
	}))
	defer server.Close()

	diff, err := New(server.URL, "t").GetPullRequestDiff(context.Background(), "pi", "ggit", 4)

	require.NoError(t, err)
	require.Len(t, diff.Files, 1)
	file := diff.Files[0]
	assert.Equal(t, "main.go", file.Name)
	assert.Equal(t, "change", file.TypeName())

	additions, deletions := file.Stats()
	assert.Equal(t, 1, additions)
	assert.Equal(t, 1, deletions)
	assert.Equal(t, "/api/v2/pi/ggit/pulls/4/diff", gotPath)
}

func TestDiffFileTypeNames(t *testing.T) {
	tests := map[DiffFileType]string{
		DiffFileAdd:    "add",
		DiffFileChange: "change",
		DiffFileDelete: "delete",
		DiffFileRename: "rename",
	}
	for typ, want := range tests {
		assert.Equal(t, want, (&DiffFile{Type: typ}).TypeName())
	}
}
