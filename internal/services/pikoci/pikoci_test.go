package pikoci

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthOK(t *testing.T) {
	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"status":"ok","db":"connected","version":"0.8.0"}`))
	}))
	defer server.Close()

	health, err := New(server.URL).Health(context.Background())

	require.NoError(t, err)
	assert.Equal(t, "/health", gotPath)
	assert.Equal(t, "ok", health.Status)
	assert.Equal(t, "connected", health.DB)
	assert.Equal(t, "0.8.0", health.Version)
}

// Con la base de datos caída PikoCI responde 503, pero el cuerpo sigue siendo
// válido y el panel debe poder mostrarlo.
func TestHealthDatabaseDown(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"status":"error","db":"disconnected","error":"no such table","version":"0.8.0"}`))
	}))
	defer server.Close()

	health, err := New(server.URL).Health(context.Background())

	require.NoError(t, err)
	assert.Equal(t, "error", health.Status)
	assert.Equal(t, "disconnected", health.DB)
	assert.Contains(t, health.Error, "no such table")
}

func TestHealthUnexpectedStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	_, err := New(server.URL).Health(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "respondió 500")
}

func TestHealthInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`no es json`))
	}))
	defer server.Close()

	_, err := New(server.URL).Health(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "respuesta inválida")
}

func TestReachable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"status":"ok","db":"connected"}`))
	}))
	defer server.Close()

	down := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"status":"error","db":"disconnected"}`))
	}))
	defer down.Close()

	assert.True(t, New(server.URL).Reachable(context.Background()))
	assert.False(t, New(down.URL).Reachable(context.Background()))
	assert.False(t, New("http://127.0.0.1:1").Reachable(context.Background()))
}

func TestBaseURLTrimsTrailingSlash(t *testing.T) {
	assert.Equal(t, "http://127.0.0.1:18080", New("http://127.0.0.1:18080/").BaseURL())
}

func TestTriggerWebhook(t *testing.T) {
	var gotMethod, gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.EscapedPath()
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	err := New(server.URL).TriggerWebhook(context.Background(), "tok/en")

	require.NoError(t, err)
	assert.Equal(t, http.MethodPost, gotMethod)
	// El token viaja escapado: un "/" no debe abrir un segmento de ruta nuevo.
	assert.Equal(t, "/webhooks/tok%2Fen", gotPath)
}

func TestTriggerWebhookWithoutTokenDoesNothing(t *testing.T) {
	var calls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
	}))
	defer server.Close()

	require.NoError(t, New(server.URL).TriggerWebhook(context.Background(), ""))
	assert.Equal(t, 0, calls)
}

func TestTriggerWebhookSurfacesError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	err := New(server.URL).TriggerWebhook(context.Background(), "token")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "respondió 404")
}

func TestTriggerWebhookUnreachable(t *testing.T) {
	err := New("http://127.0.0.1:1").TriggerWebhook(context.Background(), "token")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no se pudo avisar")
}
