package apiclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateRepo(t *testing.T) {
	var (
		gotMethod, gotPath string
		gotBody            map[string]any
	)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		require.NoError(t, json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"full_name":"pi/nuevo","is_private":true,"clone_url":"https://ggit.example.com/pi/nuevo.git","ssh_url":"git@ggit.example.com:pi/nuevo.git"}`))
	}))
	defer server.Close()

	repo, err := New(server.URL, "t").CreateRepo(context.Background(), CreateRepoOptions{
		Name:        "nuevo",
		Description: "un repo",
		Private:     true,
		AutoInit:    true,
	})

	require.NoError(t, err)
	assert.Equal(t, "pi/nuevo", repo.FullName)
	assert.True(t, repo.IsPrivate)
	assert.Equal(t, http.MethodPost, gotMethod)
	assert.Equal(t, "/api/v2/user/repos", gotPath)
	assert.Equal(t, "nuevo", gotBody["name"])
	assert.Equal(t, "un repo", gotBody["description"])
	assert.Equal(t, true, gotBody["private"])
	assert.Equal(t, true, gotBody["auto_init"])
	// Sin plantillas no se envían claves vacías.
	assert.NotContains(t, gotBody, "readme")
	assert.NotContains(t, gotBody, "gitignores")
	assert.NotContains(t, gotBody, "license")
	assert.NotContains(t, gotBody, "default_branch")
}

func TestCreateRepoInOrganization(t *testing.T) {
	var gotBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.NoError(t, json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"full_name":"acme/nuevo"}`))
	}))
	defer server.Close()

	_, err := New(server.URL, "t").CreateRepo(context.Background(), CreateRepoOptions{
		Name: "nuevo",
		UID:  9,
	})

	require.NoError(t, err)
	assert.Equal(t, float64(9), gotBody["uid"])
}

func TestCreateRepoSurfacesValidation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"message":"Key: 'createRepoPayload.Name' Error:Field validation for 'Name' failed on the 'required' tag"}`))
	}))
	defer server.Close()

	_, err := New(server.URL, "t").CreateRepo(context.Background(), CreateRepoOptions{})

	var apiErr *Error
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusBadRequest, apiErr.Status)
	assert.Contains(t, apiErr.Message, "Name")
}

func TestCreatePullRequest(t *testing.T) {
	var (
		gotPath string
		gotBody map[string]any
	)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		require.NoError(t, json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusCreated)
		// Forma del modelo interno de la incidencia.
		_, _ = w.Write([]byte(`{"ID":40,"Index":7,"Title":"Code review","IsPull":true}`))
	}))
	defer server.Close()

	raw, err := New(server.URL, "t").CreatePullRequest(context.Background(), "pi", "ggit", CreatePullRequestOptions{
		Title:      "Code review",
		Content:    "Añade aprobaciones.",
		HeadBranch: "feature/review",
		BaseBranch: "main",
	})

	require.NoError(t, err)
	assert.Equal(t, "/api/v2/pi/ggit/pulls", gotPath)
	assert.Equal(t, "Code review", gotBody["title"])
	assert.Equal(t, "Añade aprobaciones.", gotBody["content"])
	assert.Equal(t, "feature/review", gotBody["head_branch"])
	assert.Equal(t, "main", gotBody["base_branch"])

	var created struct {
		Index int64
	}
	require.NoError(t, json.Unmarshal(raw, &created))
	assert.Equal(t, int64(7), created.Index)
}

func TestCreatePullRequestOmitsEmptyContent(t *testing.T) {
	var gotBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.NoError(t, json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"Index":8}`))
	}))
	defer server.Close()

	_, err := New(server.URL, "t").CreatePullRequest(context.Background(), "pi", "ggit", CreatePullRequestOptions{
		Title:      "Sin cuerpo",
		HeadBranch: "feature/x",
		BaseBranch: "main",
	})

	require.NoError(t, err)
	assert.NotContains(t, gotBody, "content")
}

func TestCreatePullRequestSurfacesMissingBranch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"message":"Head branch does not exist."}`))
	}))
	defer server.Close()

	_, err := New(server.URL, "t").CreatePullRequest(context.Background(), "pi", "ggit", CreatePullRequestOptions{
		Title:      "x",
		HeadBranch: "no-existe",
		BaseBranch: "main",
	})

	var apiErr *Error
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusUnprocessableEntity, apiErr.Status)
	assert.Equal(t, "Head branch does not exist.", apiErr.Message)
}
