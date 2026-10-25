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

func TestDeleteRepo(t *testing.T) {
	var (
		gotMethod, gotPath string
		gotBody            map[string]string
	)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		require.NoError(t, json.NewDecoder(r.Body).Decode(&gotBody))
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	status, err := New(server.URL, "t").DeleteRepo(context.Background(), "pi", "ggit")

	require.NoError(t, err)
	assert.Equal(t, "ok", status)
	assert.Equal(t, http.MethodDelete, gotMethod)
	assert.Equal(t, "/api/v2/pi/ggit/settings", gotPath)
	// El servidor exige el nombre como confirmación.
	assert.Equal(t, map[string]string{"name": "ggit"}, gotBody)
}

func TestDeleteRepoSurfacesMismatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"message":"Confirmation name does not match."}`))
	}))
	defer server.Close()

	_, err := New(server.URL, "t").DeleteRepo(context.Background(), "pi", "ggit")

	var apiErr *Error
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusUnprocessableEntity, apiErr.Status)
	assert.Equal(t, "Confirmation name does not match.", apiErr.Message)
}

func TestUpdateIssueSendsOnlyProvidedFields(t *testing.T) {
	var (
		gotPath string
		gotBody map[string]any
	)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		require.NoError(t, json.NewDecoder(r.Body).Decode(&gotBody))
		_, _ = w.Write([]byte(`{"Index":12,"Title":"nuevo"}`))
	}))
	defer server.Close()

	title := "nuevo"
	raw, err := New(server.URL, "t").UpdateIssue(context.Background(), "pi", "ggit", 12, UpdateIssueOptions{
		Title: &title,
	})

	require.NoError(t, err)
	assert.Equal(t, "/api/v2/pi/ggit/issues/12", gotPath)
	assert.Equal(t, map[string]any{"title": "nuevo"}, gotBody)
	assert.Contains(t, string(raw), `"Title":"nuevo"`)
}

func TestUpdateIssueSupportsZeroValues(t *testing.T) {
	var gotBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.NoError(t, json.NewDecoder(r.Body).Decode(&gotBody))
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	zero := int64(0)
	_, err := New(server.URL, "t").UpdateIssue(context.Background(), "pi", "ggit", 12, UpdateIssueOptions{
		AssigneeID: &zero,
		IsClosed:   ptrBool(true),
	})

	require.NoError(t, err)
	// un puntero a 0 sí viaja: sirve para desasignar.
	assert.Equal(t, float64(0), gotBody["assignee_id"])
	assert.Equal(t, true, gotBody["is_closed"])
}

func TestGetOrg(t *testing.T) {
	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"id":9,"username":"acme","full_name":"Acme","description":"org","private":false}`))
	}))
	defer server.Close()

	org, err := New(server.URL, "t").GetOrg(context.Background(), "acme")

	require.NoError(t, err)
	assert.Equal(t, int64(9), org.ID)
	assert.Equal(t, "acme", org.Username)
	assert.Equal(t, "/api/v2/orgs/acme", gotPath)
}

func TestAccessTokenListAndDelete(t *testing.T) {
	var (
		deletedPath string
		gotMethod   string
	)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		switch r.Method {
		case http.MethodGet:
			_, _ = w.Write([]byte(`[{"id":1,"name":"ggit-cli"},{"id":2,"name":"otro"}]`))
		case http.MethodDelete:
			deletedPath = r.URL.Path
			_, _ = w.Write([]byte(`{"status":"ok"}`))
		}
	}))
	defer server.Close()

	client := New(server.URL, "t")

	tokens, err := client.ListAccessTokens(context.Background())
	require.NoError(t, err)
	require.Len(t, tokens, 2)
	assert.Equal(t, "ggit-cli", tokens[0].Name)

	require.NoError(t, client.DeleteAccessToken(context.Background(), tokens[0].ID))
	assert.Equal(t, http.MethodDelete, gotMethod)
	assert.Equal(t, "/api/v2/user/settings/tokens/1", deletedPath)
}

func ptrBool(v bool) *bool { return &v }
