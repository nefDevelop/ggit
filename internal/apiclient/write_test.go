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

func TestCreateIssue(t *testing.T) {
	var (
		gotMethod, gotPath, gotAuth string
		gotBody                     map[string]any
	)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotAuth = r.Header.Get("Authorization")
		require.NoError(t, json.NewDecoder(r.Body).Decode(&gotBody))

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":31,"number":12,"title":"Falta el inbox","state":"open"}`))
	}))
	defer server.Close()

	issue, err := New(server.URL, "t").CreateIssue(context.Background(), "pi", "ggit", CreateIssueOptions{
		Title:       "Falta el inbox",
		Content:     "El inbox no pagina.",
		MilestoneID: 3,
		AssigneeID:  7,
		LabelIDs:    []int64{1, 2},
	})

	require.NoError(t, err)
	assert.Equal(t, int64(12), issue.Index)
	assert.Equal(t, "open", string(issue.State))
	assert.Equal(t, http.MethodPost, gotMethod)
	assert.Equal(t, "/api/v2/pi/ggit/issues", gotPath)
	assert.Equal(t, "token t", gotAuth)
	assert.Equal(t, "Falta el inbox", gotBody["title"])
	assert.Equal(t, "El inbox no pagina.", gotBody["content"])
	assert.Equal(t, float64(3), gotBody["milestone_id"])
	assert.Equal(t, float64(7), gotBody["assignee_id"])
	assert.Equal(t, []any{float64(1), float64(2)}, gotBody["label_ids"])
}

func TestCreateIssueOmitsEmptyFields(t *testing.T) {
	var gotBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.NoError(t, json.NewDecoder(r.Body).Decode(&gotBody))
		_, _ = w.Write([]byte(`{"number":1,"title":"solo título"}`))
	}))
	defer server.Close()

	_, err := New(server.URL, "t").CreateIssue(context.Background(), "pi", "ggit", CreateIssueOptions{
		Title: "solo título",
	})

	require.NoError(t, err)
	assert.NotContains(t, gotBody, "label_ids")
	assert.NotContains(t, gotBody, "milestone_id")
	assert.NotContains(t, gotBody, "content")
}

func TestCreateIssueComment(t *testing.T) {
	var (
		gotPath string
		gotBody map[string]string
	)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		require.NoError(t, json.NewDecoder(r.Body).Decode(&gotBody))
		_, _ = w.Write([]byte(`{"id":9,"body":"hola","html_url":"https://ggit.example.com/pi/ggit/issues/12#issuecomment-9"}`))
	}))
	defer server.Close()

	comment, err := New(server.URL, "t").CreateIssueComment(context.Background(), "pi", "ggit", 12, "hola")

	require.NoError(t, err)
	assert.Equal(t, int64(9), comment.ID)
	assert.Contains(t, comment.HTMLURL, "#issuecomment-9")
	assert.Equal(t, "/api/v2/pi/ggit/issues/12/comments", gotPath)
	assert.Equal(t, map[string]string{"body": "hola"}, gotBody)
}

func TestPostSurfacesValidationError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"message":"Key: 'CreateIssueOptions.Title' Error:Field validation for 'Title' failed on the 'required' tag"}`))
	}))
	defer server.Close()

	_, err := New(server.URL, "t").CreateIssue(context.Background(), "pi", "ggit", CreateIssueOptions{})

	var apiErr *Error
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusBadRequest, apiErr.Status)
	assert.Contains(t, apiErr.Message, "Title")
}
