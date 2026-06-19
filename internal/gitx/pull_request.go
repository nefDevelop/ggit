package gitx

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gogs/git-module"
	log "github.com/nef734/ggit/internal/log"
)

// PullRequestMeta contains metadata for a pull request.
type PullRequestMeta struct {
	// The merge base of the pull request.
	MergeBase string
	// The commits that are requested to be merged.
	Commits []*git.Commit
	// The number of files changed.
	NumFiles int
}

func (module) PullRequestMeta(headPath, basePath, headBranch, baseBranch string) (*PullRequestMeta, error) {
	tmpRemoteBranch := baseBranch

	// We need to create a temporary remote when the pull request is sent from a forked repository.
	if headPath != basePath {
		tmpRemote := fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Int63())
		err := Module.RemoteAdd(headPath, tmpRemote, basePath, git.RemoteAddOptions{Fetch: true})
		if err != nil {
			return nil, fmt.Errorf("add remote: %v", err)
		}
		defer func() {
			err := Module.RemoteRemove(headPath, tmpRemote)
			if err != nil {
				log.Error("Failed to remove remote %q [path: %s]: %v", tmpRemote, headPath, err)
				return
			}
		}()

		tmpRemoteBranch = "remotes/" + tmpRemote + "/" + baseBranch
	}

	mergeBase, err := Module.MergeBase(headPath, tmpRemoteBranch, headBranch)
	if err != nil {
		return nil, fmt.Errorf("get merge base"+": %w", err)
	}

	commits, err := Module.Log(headPath, mergeBase+"..."+headBranch)
	if err != nil {
		return nil, fmt.Errorf("get commits"+": %w", err)
	}

	// Count number of changed files
	names, err := Module.DiffNameOnly(headPath, tmpRemoteBranch, headBranch, git.DiffNameOnlyOptions{NeedsMergeBase: true})
	if err != nil {
		return nil, fmt.Errorf("get changed files"+": %w", err)
	}

	return &PullRequestMeta{
		MergeBase: mergeBase,
		Commits:   commits,
		NumFiles:  len(names),
	}, nil
}
