package repox

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gogs/git-module"
	"github.com/nef734/ggit/internal/conf"
)

// CloneLink represents different types of clone URLs of repository.
type CloneLink struct {
	SSH   string
	HTTPS string
}

// NewCloneLink returns clone URLs using given owner and repository name.
func NewCloneLink(owner, repo string) *CloneLink {
	cl := new(CloneLink)
	if conf.SSH.Port != 22 {
		cl.SSH = fmt.Sprintf("ssh://%s@%s:%d/%s/%s.git", conf.App.RunUser, conf.SSH.Domain, conf.SSH.Port, owner, repo)
	} else {
		cl.SSH = fmt.Sprintf("%s@%s:%s/%s.git", conf.App.RunUser, conf.SSH.Domain, owner, repo)
	}
	cl.HTTPS = HTTPSCloneURL(owner, repo)
	return cl
}

// HTTPSCloneURL returns HTTPS clone URL using given owner and repository name.
func HTTPSCloneURL(owner, repo string) string {
	return fmt.Sprintf("%s%s/%s.git", conf.Server.ExternalURL, owner, repo)
}

// HTMLURL returns HTML URL using given owner and repository name.
func HTMLURL(owner, repo string) string {
	return conf.Server.ExternalURL + owner + "/" + repo
}

// CompareCommitsPath returns the comparison path using given owner, repository,
// and commit IDs.
func CompareCommitsPath(owner, repo, oldCommitID, newCommitID string) string {
	return fmt.Sprintf("%s/%s/compare/%s...%s", owner, repo, oldCommitID, newCommitID)
}

// UserPath returns the absolute path for storing user repositories.
func UserPath(user string) string {
	return filepath.Join(conf.Repository.Root, strings.ToLower(user))
}

// RepositoryPath returns the absolute path using given user and repository
// name.
func RepositoryPath(owner, repo string) string {
	return filepath.Join(UserPath(owner), strings.ToLower(repo)+".git")
}

// RepositoryLocalPath returns the absolute path of the repository local copy
// with the given ID.
func RepositoryLocalPath(repoID int64) string {
	return filepath.Join(conf.Server.AppDataPath, "tmp", "local-repo", strconv.FormatInt(repoID, 10))
}

var hooksTpls = map[git.HookName]string{
	"pre-receive":  "#!/usr/bin/env %s\n\"%s\" hook --config='%s' pre-receive\n",
	"update":       "#!/usr/bin/env %s\n\"%s\" hook --config='%s' update $1 $2 $3\n",
	"post-receive": "#!/usr/bin/env %s\n\"%s\" hook --config='%s' post-receive\n",
}

// CreateDelegateHooks creates delegate hooks for the given repository path.
func CreateDelegateHooks(repoPath string) error {
	return SelfHealHooks(repoPath)
}

// SelfHealHooks checks and repairs hooks for the given repository path.
// It ensures hooks exist, have the correct content (binary path, config path), and correct permissions.
func SelfHealHooks(repoPath string) error {
	for _, name := range git.ServerSideHooks {
		hookPath := filepath.Join(repoPath, "hooks", string(name))
		expectedContent := fmt.Sprintf(hooksTpls[name], conf.Repository.ScriptType, conf.AppPath(), conf.CustomConf)

		shouldWrite := false
		info, err := os.Stat(hookPath)
		if err != nil {
			if os.IsNotExist(err) {
				shouldWrite = true
			} else {
				return fmt.Errorf("stat hook '%s': %v", hookPath, err)
			}
		} else {
			// Check content and permissions
			content, err := os.ReadFile(hookPath)
			if err != nil || string(content) != expectedContent || info.Mode() != os.ModePerm {
				shouldWrite = true
			}
		}

		if shouldWrite {
			if err := os.MkdirAll(filepath.Dir(hookPath), os.ModePerm); err != nil {
				return fmt.Errorf("create hooks directory: %v", err)
			}
			if err := os.WriteFile(hookPath, []byte(expectedContent), os.ModePerm); err != nil {
				return fmt.Errorf("repair delegate hook '%s': %v", hookPath, err)
			}
		}
	}
	return nil
}
