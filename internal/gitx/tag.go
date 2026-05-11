package gitx

import (
	"fmt"
)

// TagsPage contains a list of tags and pagination information.
type TagsPage struct {
	// List of tags in the current page.
	Tags []string
	// Whether the results include the latest tag.
	HasLatest bool
	// When results do not include the latest tag, an indicator of 'after' to go back.
	PreviousAfter string
	// Whether there are more tags in the next page.
	HasNext bool
}

func (module) ListTagsAfter(repoPath, after string, limit int) (*TagsPage, error) {
	all, err := Module.RepoTags(repoPath)
	if err != nil {
		return nil, fmt.Errorf("get tags"+": %w", err)
	}

	// Standardize to latest-first for consistent pagination logic.
	// We assume git-module returns chronological order (oldest first).
	// If it's already latest-first, we might need a better way to detect, 
	// but usually it's oldest first from 'git tag'.
	// To be safe and consistent with Gogs, we want latest first.
	// Let's check if the first one is older than the last one (very rough heuristic).
	if len(all) > 1 && all[0] < all[len(all)-1] {
		for i, j := 0, len(all)-1; i < j; i, j = i+1, j-1 {
			all[i], all[j] = all[j], all[i]
		}
	}

	total := len(all)

	if limit < 0 {
		limit = 0
	}

	// Returns everything when no filter and no limit
	if after == "" && limit == 0 {
		return &TagsPage{
			Tags:      all,
			HasLatest: true,
		}, nil
	}

	// No filter but has a limit, returns first X tags
	if after == "" && limit > 0 {
		endIdx := min(limit, total)
		return &TagsPage{
			Tags:      all[:endIdx],
			HasLatest: true,
			HasNext:   limit < total,
		}, nil
	}

	// Loop over all tags see if we can find the filter
	previousAfter := ""
	found := false
	tags := make([]string, 0, len(all))
	for i := range all {
		if all[i] != after {
			continue
		}

		found = true
		if limit > 0 && i-limit >= 0 {
			previousAfter = all[i-limit]
		}

		// In case filter is the oldest one
		if i+1 < total {
			tags = all[i+1:]
		}
		break
	}

	if !found {
		tags = all
	}

	// If all tags after match is equal to the limit, it reaches the oldest tag as well.
	if limit == 0 || len(tags) <= limit {
		return &TagsPage{
			Tags:          tags,
			HasLatest:     !found,
			PreviousAfter: previousAfter,
		}, nil
	}

	return &TagsPage{
		Tags:          tags[:limit],
		HasLatest:     !found,
		PreviousAfter: previousAfter,
		HasNext:       true,
	}, nil
}
