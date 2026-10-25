package apiclient

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientSendsTokenHeader(t *testing.T) {
	var gotAuth, gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"id":1,"full_name":"pi/ggit"}]`))
	}))
	defer server.Close()

	client := New(server.URL, "secreto")
	repos, err := client.ListUserRepos(context.Background())

	require.NoError(t, err)
	require.Len(t, repos, 1)
	assert.Equal(t, "pi/ggit", repos[0].FullName)
	// El servidor espera el esquema "token", no "Bearer".
	assert.Equal(t, "token secreto", gotAuth)
	assert.Equal(t, "/api/v2/user/repos", gotPath)
}

func TestClientOmitsHeaderWithoutToken(t *testing.T) {
	var gotAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		_, _ = w.Write([]byte(`{"status":"ok","user":{"id":1,"username":"pi"}}`))
	}))
	defer server.Close()

	user, err := New(server.URL, "").AuthStatus(context.Background())

	require.NoError(t, err)
	assert.Equal(t, "pi", user.Username)
	assert.Empty(t, gotAuth)
}

func TestClientReturnsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"message":"Only authenticated user is allowed to call APIs."}`))
	}))
	defer server.Close()

	_, err := New(server.URL, "").AuthStatus(context.Background())

	var apiErr *Error
	require.True(t, errors.As(err, &apiErr))
	assert.Equal(t, http.StatusUnauthorized, apiErr.Status)
	assert.Equal(t, "Only authenticated user is allowed to call APIs.", apiErr.Message)
}

func TestClientSendsQueryAndDecodesSearchWrapper(t *testing.T) {
	var gotQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		_, _ = w.Write([]byte(`{"ok":true,"total":2,"repos":[{"full_name":"pi/ggit"},{"full_name":"pi/otro"}]}`))
	}))
	defer server.Close()

	repos, total, err := New(server.URL, "t").SearchRepos(context.Background(), "ggit", 3)

	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, repos, 2)
	assert.Contains(t, gotQuery, "q=ggit")
	assert.Contains(t, gotQuery, "page=3")
}

func TestGetRepoMatchesFullName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"ok":true,"total":2,"repos":[{"full_name":"pi/ggit-cli"},{"full_name":"pi/ggit"}]}`))
	}))
	defer server.Close()

	repo, err := New(server.URL, "t").GetRepo(context.Background(), "pi/ggit")

	require.NoError(t, err)
	assert.Equal(t, "pi/ggit", repo.FullName)

	_, err = New(server.URL, "t").GetRepo(context.Background(), "pi/inexistente")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no encontrado")
}

func TestSplitRepoName(t *testing.T) {
	tests := []struct {
		input   string
		owner   string
		name    string
		wantErr bool
	}{
		{input: "pi/ggit", owner: "pi", name: "ggit"},
		{input: "/pi/ggit/", owner: "pi", name: "ggit"},
		{input: "ggit", wantErr: true},
		{input: "pi/ggit/extra", wantErr: true},
		{input: "/", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			owner, name, err := SplitRepoName(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.owner, owner)
			assert.Equal(t, tt.name, name)
		})
	}
}

func TestNormalizeHost(t *testing.T) {
	assert.Equal(t, "https://ggit.example.com", NormalizeHost(" ggit.example.com "))
	assert.Equal(t, "http://localhost:3005", NormalizeHost("http://localhost:3005"))
	assert.Equal(t, "", NormalizeHost("  "))
}
