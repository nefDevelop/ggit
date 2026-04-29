package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommitStatuses_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}

	db := newTestDB(t, "CommitStatuses_Create")
	ctx := context.Background()

	repo := &Repository{
		OwnerID: 1,
		Name:    "test-repo",
	}
	err := db.Create(repo).Error
	require.NoError(t, err)

	handle := &DB{db: db}
	status := &CommitStatus{
		RepoID:      repo.ID,
		SHA:         "abcdef1234567890abcdef1234567890abcdef12",
		State:       CommitStatusSuccess,
		TargetURL:   "https://ci.example.com/build/1",
		Description: "Build passed",
		Context:     "ci/jenkins",
		CreatorID:   1,
	}

	t.Run("Success", func(t *testing.T) {
		err := handle.CommitStatuses().Create(ctx, status)
		assert.NoError(t, err)
		assert.NotZero(t, status.ID)
		assert.NotZero(t, status.CreatedAt)
	})

	t.Run("Duplicate context for same SHA", func(t *testing.T) {
		// Should allow creating another status for same context (history)
		status2 := &CommitStatus{
			RepoID:      repo.ID,
			SHA:         status.SHA,
			State:       CommitStatusFailure,
			Context:     status.Context,
			CreatorID:   1,
		}
		err := handle.CommitStatuses().Create(ctx, status2)
		assert.NoError(t, err)
	})

	t.Run("Different context", func(t *testing.T) {
		status3 := &CommitStatus{
			RepoID:      repo.ID,
			SHA:         status.SHA,
			State:       CommitStatusPending,
			Context:     "security/snyk",
			CreatorID:   1,
		}
		err := handle.CommitStatuses().Create(ctx, status3)
		assert.NoError(t, err)
	})
}

func TestCommitStatuses_GetLatestBySHA(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}

	db := newTestDB(t, t.Name())
	handle := &DB{db: db}
	ctx := context.Background()

	repoID := int64(1)
	sha := "abcdef1234567890abcdef1234567890abcdef12"

	// Add some statuses
	statuses := []*CommitStatus{
		{RepoID: repoID, SHA: sha, State: CommitStatusSuccess, Context: "ci", CreatedAt: 100},
		{RepoID: repoID, SHA: sha, State: CommitStatusFailure, Context: "ci", CreatedAt: 200}, // Latest for "ci"
		{RepoID: repoID, SHA: sha, State: CommitStatusPending, Context: "security", CreatedAt: 150}, // Latest for "security"
		{RepoID: repoID, SHA: "other", State: CommitStatusSuccess, Context: "ci"},
	}
	for _, s := range statuses {
		err := db.Create(s).Error
		require.NoError(t, err)
	}

	t.Run("Get latest", func(t *testing.T) {
		latest, err := handle.CommitStatuses().GetLatestBySHA(ctx, repoID, sha)
		assert.NoError(t, err)
		assert.Len(t, latest, 2)

		foundCi := false
		foundSec := false
		for _, s := range latest {
			if s.Context == "ci" {
				assert.Equal(t, CommitStatusFailure, s.State)
				foundCi = true
			}
			if s.Context == "security" {
				assert.Equal(t, CommitStatusPending, s.State)
				foundSec = true
			}
		}
		assert.True(t, foundCi)
		assert.True(t, foundSec)
	})
}

func TestCommitStatuses_GetCombinedStatusBySHA(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}

	db := newTestDB(t, t.Name())
	handle := &DB{db: db}
	ctx := context.Background()

	repoID := int64(1)
	
	tests := []struct {
		name     string
		statuses []*CommitStatus
		expected CommitStatusState
	}{
		{
			name: "All success",
			statuses: []*CommitStatus{
				{State: CommitStatusSuccess, Context: "ci"},
				{State: CommitStatusSuccess, Context: "sec"},
			},
			expected: CommitStatusSuccess,
		},
		{
			name: "One failure",
			statuses: []*CommitStatus{
				{State: CommitStatusSuccess, Context: "ci"},
				{State: CommitStatusFailure, Context: "sec"},
			},
			expected: CommitStatusFailure,
		},
		{
			name: "One error",
			statuses: []*CommitStatus{
				{State: CommitStatusSuccess, Context: "ci"},
				{State: CommitStatusError, Context: "sec"},
			},
			expected: CommitStatusFailure, // Error is treated as failure in combined
		},
		{
			name: "Pending has precedence over success",
			statuses: []*CommitStatus{
				{State: CommitStatusSuccess, Context: "ci"},
				{State: CommitStatusPending, Context: "sec"},
			},
			expected: CommitStatusPending,
		},
		{
			name: "Failure has precedence over pending",
			statuses: []*CommitStatus{
				{State: CommitStatusFailure, Context: "ci"},
				{State: CommitStatusPending, Context: "sec"},
			},
			expected: CommitStatusFailure,
		},
		{
			name: "Warning",
			statuses: []*CommitStatus{
				{State: CommitStatusSuccess, Context: "ci"},
				{State: CommitStatusWarning, Context: "sec"},
			},
			expected: CommitStatusWarning,
		},
		{
			name: "Empty",
			statuses: []*CommitStatus{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sha := tt.name
			for _, s := range tt.statuses {
				s.RepoID = repoID
				s.SHA = sha
				err := db.Create(s).Error
				require.NoError(t, err)
			}

			combined := handle.CommitStatuses().GetCombinedStatusBySHA(ctx, repoID, sha)
			assert.Equal(t, tt.expected, combined)
		})
	}
}
