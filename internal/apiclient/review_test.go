package apiclient

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetIssueClosed(t *testing.T) {
	tests := []struct {
		name       string
		closed     bool
		wantClosed bool
	}{
		{name: "cerrar", closed: true, wantClosed: true},
		{name: "reabrir", closed: false, wantClosed: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				gotMethod, gotPath string
				gotBody            map[string]any
			)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotMethod = r.Method
				gotPath = r.URL.Path
				require.NoError(t, json.NewDecoder(r.Body).Decode(&gotBody))
				_, _ = w.Write([]byte(`{"id":31,"index":12,"is_closed":true}`))
			}))
			defer server.Close()

			raw, err := New(server.URL, "t").SetIssueClosed(context.Background(), "pi", "ggit", 12, tt.closed)

			require.NoError(t, err)
			assert.JSONEq(t, `{"id":31,"index":12,"is_closed":true}`, string(raw))
			assert.Equal(t, http.MethodPut, gotMethod)
			assert.Equal(t, "/api/v2/pi/ggit/issues/12", gotPath)
			assert.Equal(t, tt.wantClosed, gotBody["is_closed"])
		})
	}
}

func TestCreatePullReview(t *testing.T) {
	var (
		gotPath string
		gotBody map[string]any
	)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		require.NoError(t, json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":5,"type":1,"content":"LGTM"}`))
	}))
	defer server.Close()

	raw, err := New(server.URL, "t").CreatePullReview(context.Background(), "pi", "ggit", 4, CreatePullReviewOptions{
		Type:    ReviewApprove,
		Content: "LGTM",
	})

	require.NoError(t, err)
	assert.JSONEq(t, `{"id":5,"type":1,"content":"LGTM"}`, string(raw))
	assert.Equal(t, "/api/v2/pi/ggit/pulls/4/reviews", gotPath)
	assert.Equal(t, float64(1), gotBody["type"])
	assert.Equal(t, "LGTM", gotBody["content"])
}

func TestCreatePullReviewOmitsEmptyContent(t *testing.T) {
	var gotBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.NoError(t, json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":6,"type":1}`))
	}))
	defer server.Close()

	_, err := New(server.URL, "t").CreatePullReview(context.Background(), "pi", "ggit", 4, CreatePullReviewOptions{
		Type: ReviewApprove,
	})

	require.NoError(t, err)
	assert.NotContains(t, gotBody, "content")
	assert.Equal(t, float64(1), gotBody["type"])
}

func TestMergePullRequest(t *testing.T) {
	var (
		gotMethod, gotPath string
		gotBody            string
	)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		raw, _ := io.ReadAll(r.Body)
		gotBody = string(raw)
		_, _ = w.Write([]byte(`{"status":"merged"}`))
	}))
	defer server.Close()

	status, err := New(server.URL, "t").MergePullRequest(context.Background(), "pi", "ggit", 4)

	require.NoError(t, err)
	assert.Equal(t, "merged", status)
	assert.Equal(t, http.MethodPost, gotMethod)
	assert.Equal(t, "/api/v2/pi/ggit/pulls/4/merge", gotPath)
	assert.Empty(t, gotBody)
}

func TestMergePullRequestReportsServerRefusal(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"message":"Pull request is not mergeable."}`))
	}))
	defer server.Close()

	_, err := New(server.URL, "t").MergePullRequest(context.Background(), "pi", "ggit", 4)

	var apiErr *Error
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusUnprocessableEntity, apiErr.Status)
	assert.Equal(t, "Pull request is not mergeable.", apiErr.Message)
}
