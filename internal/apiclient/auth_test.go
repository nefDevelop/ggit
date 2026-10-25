package apiclient

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithBasicAuthSendsBasicHeader(t *testing.T) {
	var gotAuth, gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"id":1,"name":"ggit-cli","sha1":"nuevo-token"}`))
	}))
	defer server.Close()

	token, err := New(server.URL, "").WithBasicAuth("pi", "secreto").CreateAccessToken(context.Background(), "ggit-cli")

	require.NoError(t, err)
	assert.Equal(t, "nuevo-token", token)
	assert.Equal(t, "/api/v2/user/settings/tokens", gotPath)
	want := "Basic " + base64.StdEncoding.EncodeToString([]byte("pi:secreto"))
	assert.Equal(t, want, gotAuth)
}

func TestTokenAuthWinsOverBasicAuth(t *testing.T) {
	var gotAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		_, _ = w.Write([]byte(`{"status":"ok","user":{"username":"pi"}}`))
	}))
	defer server.Close()

	_, err := New(server.URL, "existente").WithBasicAuth("pi", "secreto").AuthStatus(context.Background())

	require.NoError(t, err)
	assert.Equal(t, "token existente", gotAuth)
}

func TestCreateAccessTokenSendsName(t *testing.T) {
	var gotBody map[string]string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.NoError(t, json.NewDecoder(r.Body).Decode(&gotBody))
		_, _ = w.Write([]byte(`{"id":2,"name":"cli-termux","sha1":"abc123"}`))
	}))
	defer server.Close()

	token, err := New(server.URL, "").WithBasicAuth("pi", "x").CreateAccessToken(context.Background(), "cli-termux")

	require.NoError(t, err)
	assert.Equal(t, "abc123", token)
	assert.Equal(t, map[string]string{"name": "cli-termux"}, gotBody)
}

func TestCreateAccessTokenRequiresSha(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"id":2,"name":"cli","sha1":""}`))
	}))
	defer server.Close()

	_, err := New(server.URL, "").WithBasicAuth("pi", "x").CreateAccessToken(context.Background(), "cli")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no devolvió el token")
}

func TestCreateAccessTokenSurfacesDuplicateName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"message":"Access token name has already been used."}`))
	}))
	defer server.Close()

	_, err := New(server.URL, "").WithBasicAuth("pi", "x").CreateAccessToken(context.Background(), "ggit-cli")

	var apiErr *Error
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusUnprocessableEntity, apiErr.Status)
	assert.Contains(t, apiErr.Message, "already been used")
}
