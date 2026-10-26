package pikoci

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

func TestPublishPipeline(t *testing.T) {
	var (
		gotMethod, gotPath, gotAuth string
		gotBody                     map[string]any
	)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotAuth = r.Header.Get("Authorization")
		require.NoError(t, json.NewDecoder(r.Body).Decode(&gotBody))
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	hcl := []byte("job \"test\" {}\n")
	err := NewWithAuth(server.URL, "ci", "pko_api").PublishPipeline(context.Background(), "acme-web", hcl, map[string]string{
		"repo": "acme/web",
	})

	require.NoError(t, err)
	assert.Equal(t, http.MethodPut, gotMethod)
	assert.Equal(t, "/teams/ci/pipelines/acme-web", gotPath)
	assert.Equal(t, "Bearer pko_api", gotAuth)

	// La API espera el HCL en base64 dentro de "config".
	decoded, err := base64.StdEncoding.DecodeString(gotBody["config"].(string))
	require.NoError(t, err)
	assert.Equal(t, string(hcl), string(decoded))
	assert.Equal(t, map[string]any{"repo": "acme/web"}, gotBody["vars"])
}

func TestPublishPipelineRequiresAPIToken(t *testing.T) {
	err := New("http://127.0.0.1:1").PublishPipeline(context.Background(), "p", nil, nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "API_TOKEN")
}

func TestResourceWebhookToken(t *testing.T) {
	var gotMethod, gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"token":"hook-token"}`))
	}))
	defer server.Close()

	token, err := NewWithAuth(server.URL, "main", "pko").ResourceWebhookToken(context.Background(), "acme-web", "repo")

	require.NoError(t, err)
	assert.Equal(t, "hook-token", token)
	assert.Equal(t, http.MethodPost, gotMethod)
	assert.Equal(t, "/teams/main/pipelines/acme-web/resources/repo/webhook_token", gotPath)
}

func TestResourceWebhookTokenRequiresToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	_, err := NewWithAuth(server.URL, "main", "pko").ResourceWebhookToken(context.Background(), "p", "repo")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no devolvió el token")
}

func TestDeletePipelineIgnoresNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	assert.NoError(t, NewWithAuth(server.URL, "main", "pko").DeletePipeline(context.Background(), "p"))
}

func TestAdminErrorSurfacesMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"message":"invalid pipeline config"}`))
	}))
	defer server.Close()

	err := NewWithAuth(server.URL, "main", "pko").PublishPipeline(context.Background(), "p", []byte("x"), nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "respondió 422")
	assert.Contains(t, err.Error(), "invalid pipeline config")
}

func TestPipelineName(t *testing.T) {
	assert.Equal(t, "acme-web", PipelineName("Acme", "Web"))
	assert.Equal(t, "acme-web-app", PipelineName("acme", "web app"))
	assert.Equal(t, "acme-web-app", PipelineName("acme", "web_app"))
	assert.Equal(t, "acme-web", PipelineName("acme", "web."))
}
