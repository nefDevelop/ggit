package database

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/nef734/ggit/internal/dbtest"
	apiv2types "github.com/nef734/ggit/internal/route/api/v2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestReleaseListMatchesPerIssueLoading is the guard that the batch loader did
// not change any result: the same releases go through both paths and the fields
// APIFormat serialises are compared.
//
// It also covers the mistake worth covering, that release attachments hang off
// release_id and not issue_id. Loading them with the wrong column silently
// returned an empty list instead of failing.
func TestReleaseListMatchesPerIssueLoading(t *testing.T) {
	// suite keeps the database files apart; suffix is part of the seeded names and
	// must match across both runs, otherwise the comparison would be between two
	// different datasets rather than two loaders.
	build := func(t *testing.T, suite, suffix string) []*Release {
		t.Helper()
		db := dbtest.NewDB(t, "ReleaseListEquivalence"+suite,
			new(User), new(Repository), new(Release), new(Attachment), new(Watch), new(IssueLabel))
		SetHandle(db)
		ctx := context.Background()

		// Timestamps are pinned. Left to the clock, two runs straddling a second
		// boundary differ in CreatedAt and the comparison fails for a reason that
		// has nothing to do with the loader.
		seeded := time.Unix(1700000000, 0)

		publisher, err := Handle.Users().Create(ctx, "rel-publisher-"+suffix, "rel-"+suffix+"@example.com", CreateUserOptions{})
		require.NoError(t, err)
		require.NoError(t, db.Model(FromDomainUser(publisher)).Updates(map[string]any{
			"created_unix": seeded.Unix(),
			"updated_unix": seeded.Unix(),
		}).Error)

		repo, err := Handle.Repositories().Create(ctx, publisher.ID, CreateRepoOptions{Name: "rel-repo-" + suffix})
		require.NoError(t, err)
		domainRepo := FromDomainRepository(repo)

		releases := make([]*Release, 0, 3)
		for i := 1; i <= 3; i++ {
			release := &Release{
				RepoID:      domainRepo.ID,
				PublisherID: publisher.ID,
				TagName:     "v1." + strconv.Itoa(i),
				Title:       "release " + strconv.Itoa(i),
				Note:        "notas",
			}
			require.NoError(t, db.Create(release).Error)
			require.NoError(t, db.Model(release).Update("created_unix", seeded.Unix()).Error)

			// Every release gets one attachment, plus a second on the first, so
			// both the single and the multi attachment paths are covered.
			require.NoError(t, db.Create(&Attachment{
				ReleaseID: release.ID, UUID: "uuid-" + strconv.Itoa(i), Name: "a.txt",
			}).Error)
			if i == 1 {
				require.NoError(t, db.Create(&Attachment{
					ReleaseID: release.ID, UUID: "uuid-extra", Name: "b.txt",
				}).Error)
			}
			releases = append(releases, release)
		}
		return releases
	}

	perRelease := build(t, "PerRelease", "d")
	for i := range perRelease {
		perRelease[i].Attachments = nil
		perRelease[i].Publisher = nil
		require.NoError(t, perRelease[i].LoadAttributes())
	}
	perFormatted := make([]*apiv2types.RepositoryReleaseResponse, len(perRelease))
	for i := range perRelease {
		perFormatted[i] = perRelease[i].APIFormat()
	}

	batched := build(t, "Batched", "d")
	for i := range batched {
		batched[i].Attachments = nil
		batched[i].Publisher = nil
	}
	require.NoError(t, ReleaseList(batched).LoadAttributes())
	batchFormatted := make([]*apiv2types.RepositoryReleaseResponse, len(batched))
	for i := range batched {
		batchFormatted[i] = batched[i].APIFormat()
	}

	require.Len(t, batchFormatted, len(perFormatted))
	for i := range perFormatted {
		assert.Equal(t, perFormatted[i], batchFormatted[i], "release %d", i)
	}
	assert.Len(t, batchFormatted[0].Attachments, 2, "el primero tiene dos adjuntos")
	assert.Len(t, batchFormatted[1].Attachments, 1)
}
