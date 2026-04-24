package context

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/editorconfig/editorconfig-core-go/v2"
	"github.com/gin-gonic/gin"

	"github.com/gogs/git-module"
	log "unknwon.dev/clog/v2"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/repox"
)

type PullRequest struct {
	BaseRepo *database.Repository
	Allowed  bool
	SameRepo bool
	HeadInfo string // [<user>:]<branch>
}

type Repository struct {
	AccessMode   database.AccessMode
	IsWatching   bool
	IsViewBranch bool
	IsViewTag    bool
	IsViewCommit bool
	Repository   *database.Repository
	Owner        *database.User
	Commit       *git.Commit
	Tag          *git.Tag
	GitRepo      *git.Repository
	BranchName   string
	TagName      string
	TreePath     string
	CommitID     string
	RepoLink     string
	CloneLink    repox.CloneLink
	CommitsCount int64
	PullRequest  *PullRequest
}

// ... helper methods (IsOwner, IsAdmin, etc.) same as before ...
func (r *Repository) IsOwner() bool   { return r.AccessMode >= database.AccessModeOwner }
func (r *Repository) IsAdmin() bool   { return r.AccessMode >= database.AccessModeAdmin }
func (r *Repository) IsWriter() bool  { return r.AccessMode >= database.AccessModeWrite }
func (r *Repository) HasAccess() bool { return r.AccessMode >= database.AccessModeRead }

func (r *Repository) CanEnableEditor() bool {
	return r.Repository.CanEnableEditor() && r.IsViewBranch && r.IsWriter() && !r.Repository.IsBranchRequirePullRequest(r.BranchName)
}

func (r *Repository) Editorconfig() (*editorconfig.Editorconfig, error) {
	commit, err := r.GitRepo.BranchCommit(r.Repository.DefaultBranch)
	if err != nil {
		return nil, errors.Wrapf(err, "get commit of branch %q ", r.Repository.DefaultBranch)
	}

	entry, err := commit.TreeEntry(".editorconfig")
	if err != nil {
		return nil, errors.Wrap(err, "get .editorconfig")
	}

	p, err := entry.Blob().Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "read .editorconfig")
	}
	return editorconfig.Parse(bytes.NewReader(p))
}

func (r *Repository) MakeURL(location any) string {
	switch location := location.(type) {
	case string:
		tempURL := url.URL{
			Path: r.RepoLink + "/" + location,
		}
		return tempURL.String()
	case url.URL:
		location.Path = r.RepoLink + "/" + location.Path
		return location.String()
	default:
		panic("location type must be either string or url.URL")
	}
}

func (r *Repository) PullRequestURL(baseBranch, headBranch string) string {
	repoLink := r.RepoLink
	if r.PullRequest.BaseRepo != nil {
		repoLink = r.PullRequest.BaseRepo.Link()
	}
	return fmt.Sprintf("%s/compare/%s...%s:%s", repoLink, baseBranch, r.Owner.Name, headBranch)
}

func RepoAssignment(pages ...bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		gctx := c.MustGet("context").(*Context)
		var (
			owner        *database.User
			err          error
			isIssuesPage bool
		)

		if len(pages) > 0 {
			isIssuesPage = pages[0]
		}

		ownerName := c.Param("username")
		repoName := strings.TrimSuffix(c.Param("reponame"), ".git")

		// Check if Ggit is installed/locked
		if !database.HasEngine || database.Handle == nil {
			return
		}

		// Check if the user is the same as the repository owner
		if gctx.IsLogged && gctx.User.LowerName == strings.ToLower(ownerName) {
			owner = gctx.User
		} else {
			owner, err = database.Handle.Users().GetByUsername(c.Request.Context(), ownerName)
			if err != nil {
				gctx.NotFoundOrError(err, "get user by name")
				return
			}
		}
		gctx.Repo.Owner = owner
		gctx.Data["Username"] = gctx.Repo.Owner.Name

		repo, err := database.GetRepositoryByName(owner.ID, repoName)
		if err != nil {
			gctx.NotFoundOrError(err, "get repository by name")
			return
		}

		gctx.Repo.Repository = repo
		gctx.Data["RepoName"] = gctx.Repo.Repository.Name
		gctx.Data["IsBareRepo"] = gctx.Repo.Repository.IsBare
		gctx.Repo.RepoLink = repo.Link()
		gctx.Data["RepoLink"] = gctx.Repo.RepoLink
		gctx.Data["RepoRelPath"] = gctx.Repo.Owner.Name + "/" + gctx.Repo.Repository.Name

		// Admin has super access
		if gctx.IsLogged && gctx.User.IsAdmin {
			gctx.Repo.AccessMode = database.AccessModeOwner
		} else {
			gctx.Repo.AccessMode = database.Handle.Permissions().AccessMode(c.Request.Context(), gctx.UserID(), repo.ID,
				database.AccessModeOptions{
					OwnerID: repo.OwnerID,
					Private: repo.IsPrivate,
				},
			)
		}

		if gctx.Repo.AccessMode == database.AccessModeNone && repo.BaseRepo != nil {
			mode := min(
				database.Handle.Permissions().AccessMode(c.Request.Context(), gctx.UserID(), repo.BaseRepo.ID,
					database.AccessModeOptions{
						OwnerID: repo.BaseRepo.OwnerID,
						Private: repo.BaseRepo.IsPrivate,
					},
				),
				database.AccessModeWrite,
			)
			gctx.Repo.AccessMode = mode
		}

		if gctx.Repo.AccessMode == database.AccessModeNone {
			if repo.IsPartialPublic() &&
				(!isIssuesPage || !repo.CanGuestViewIssues()) {
				if repo.CanGuestViewIssues() {
					gctx.Redirect(repo.Link() + "/issues")
				} else {
					gctx.NotFound()
				}
				return
			}

			if !repo.IsPartialPublic() ||
				(isIssuesPage && !repo.CanGuestViewIssues()) {
				gctx.NotFound()
				return
			}

			gctx.Repo.Repository.EnableIssues = repo.CanGuestViewIssues()
		}

		gitRepo, err := git.Open(database.RepoPath(ownerName, repoName))
		if err != nil {
			if gctx.IsLogged && (gctx.User.IsAdmin || gctx.User.ID == repo.OwnerID) {
				log.Warn("Repository directory missing for %s/%s: %v", ownerName, repoName, err)
				gctx.Data["Title"] = owner.Name + "/" + repo.Name
				gctx.Data["Repository"] = repo
				gctx.Data["Owner"] = gctx.Repo.Repository.Owner
				gctx.Data["IsRepositoryOwner"] = gctx.Repo.IsOwner()
				gctx.Data["IsRepositoryAdmin"] = gctx.Repo.IsAdmin()
				gctx.Data["IsRepositoryWriter"] = gctx.Repo.IsWriter()
				gctx.Data["IsBrokenRepo"] = true
				return
			}
			gctx.Error(err, "open repository")
			return
		}
		gctx.Repo.GitRepo = gitRepo

		tags, err := gctx.Repo.GitRepo.Tags()
		if err != nil {
			gctx.Error(err, "get tags")
			return
		}
		gctx.Data["Tags"] = tags
		gctx.Repo.Repository.NumTags = len(tags)

		gctx.Data["Title"] = owner.Name + "/" + repo.Name
		gctx.Data["Repository"] = repo
		gctx.Data["Owner"] = gctx.Repo.Repository.Owner
		gctx.Data["IsRepositoryOwner"] = gctx.Repo.IsOwner()
		gctx.Data["IsRepositoryAdmin"] = gctx.Repo.IsAdmin()
		gctx.Data["IsRepositoryWriter"] = gctx.Repo.IsWriter()

		gctx.Data["DisableSSH"] = conf.SSH.Disabled
		gctx.Data["DisableHTTP"] = conf.Repository.DisableHTTPGit
		gctx.Data["CloneLink"] = repo.CloneLink()

		if gctx.IsLogged {
			gctx.Data["IsWatchingRepo"] = database.IsWatching(gctx.User.ID, repo.ID)
			gctx.Data["IsStaringRepo"] = database.IsStaring(gctx.User.ID, repo.ID)
		}

		if gctx.Repo.Repository.IsBare {
			return
		}

		gctx.Data["TagName"] = gctx.Repo.TagName
		branches, err := gctx.Repo.GitRepo.Branches()
		if err != nil {
			gctx.Error(err, "get branches")
			return
		}
		gctx.Data["Branches"] = branches
		gctx.Data["BranchCount"] = len(branches)

		if gctx.Repo.BranchName == "" {
			if len(gctx.Repo.Repository.DefaultBranch) > 0 && gitRepo.HasBranch(gctx.Repo.Repository.DefaultBranch) {
				gctx.Repo.BranchName = gctx.Repo.Repository.DefaultBranch
			} else if len(branches) > 0 {
				gctx.Repo.BranchName = branches[0]
			}
		}
		gctx.Data["BranchName"] = gctx.Repo.BranchName
		gctx.Data["CommitID"] = gctx.Repo.CommitID

		gctx.Data["IsGuest"] = !gctx.Repo.HasAccess()
	}
}

func RepoRef() gin.HandlerFunc {
	return func(c *gin.Context) {
		gctx := c.MustGet("context").(*Context)
		if gctx.Repo.Repository == nil {
			return
		}
		if gctx.Repo.Repository.IsBare {
			return
		}

		var (
			refName string
			err     error
		)

		if gctx.Repo.GitRepo == nil {
			startOpen := time.Now()
			repoPath := database.RepoPath(gctx.Repo.Owner.Name, gctx.Repo.Repository.Name)
			gctx.Repo.GitRepo, err = git.Open(repoPath)
			if err != nil {
				gctx.Error(err, "open repository")
				return
			}
			log.Trace("git.Open took %v for %s", time.Since(startOpen), repoPath)
		}

		// In Gin, the "*" parameter is named differently depending on the route.
		// We'll assume the param name matches the wildcard in the route.
		path := c.Param("path")
		if path == "" {
			path = c.Param("*")
		}
		path = strings.Trim(path, "/")

		if path == "" {
			refName = gctx.Repo.Repository.DefaultBranch
			if refName == "" || !gctx.Repo.GitRepo.HasBranch(refName) {
				branches, err := gctx.Repo.GitRepo.Branches()
				if err != nil {
					gctx.Error(err, "get branches")
					return
				}

				if len(branches) > 0 {
					refName = branches[0]
					log.Trace("Fallback to first branch: %s", refName)
					// Update default branch if it was empty
					if gctx.Repo.Repository.DefaultBranch == "" {
						gctx.Repo.Repository.DefaultBranch = refName
						_ = gctx.Repo.Repository.Update()
					}
				} else {
					log.Warn("Repository %s/%s has no branches", gctx.Repo.Owner.Name, gctx.Repo.Repository.Name)
					gctx.Error(fmt.Errorf("no branches found"), "get branches")
					return
				}
			}
			log.Trace("Resolved ref: %s, path: %s", refName, gctx.Repo.TreePath)
			gctx.Repo.Commit, err = gctx.Repo.GitRepo.BranchCommit(refName)
			if err != nil {
				gctx.Error(err, "get branch commit")
				return
			}
			gctx.Repo.CommitID = gctx.Repo.Commit.ID.String()
			gctx.Repo.IsViewBranch = true

		} else {
			hasMatched := false
			parts := strings.Split(path, "/")
			for i, part := range parts {
				refName = strings.TrimPrefix(refName+"/"+part, "/")

				if gctx.Repo.GitRepo.HasBranch(refName) ||
					gctx.Repo.GitRepo.HasTag(refName) {
					if i < len(parts)-1 {
						gctx.Repo.TreePath = strings.Join(parts[i+1:], "/")
					}
					hasMatched = true
					break
				}
			}
			if !hasMatched {
				if len(parts[0]) == 40 {
					refName = parts[0]
					gctx.Repo.TreePath = strings.Join(parts[1:], "/")
				} else {
					refName = gctx.Repo.Repository.DefaultBranch
					gctx.Repo.TreePath = path
				}
				log.Trace("Path resolution matched: %v, ref: %s, path: %s", hasMatched, refName, gctx.Repo.TreePath)
			}

			if gctx.Repo.GitRepo.HasBranch(refName) {
				gctx.Repo.IsViewBranch = true
				gctx.Repo.Commit, err = gctx.Repo.GitRepo.BranchCommit(refName)
				if err != nil {
					gctx.Error(err, "get branch commit")
					return
				}
				gctx.Repo.CommitID = gctx.Repo.Commit.ID.String()

			} else if gctx.Repo.GitRepo.HasTag(refName) {
				gctx.Repo.IsViewTag = true
				gctx.Repo.Commit, err = gctx.Repo.GitRepo.TagCommit(refName)
				if err != nil {
					gctx.Error(err, "get tag commit")
					return
				}
				gctx.Repo.CommitID = gctx.Repo.Commit.ID.String()
			} else if len(refName) == 40 {
				gctx.Repo.IsViewCommit = true
				gctx.Repo.CommitID = refName
				gctx.Repo.Commit, err = gctx.Repo.GitRepo.CatFileCommit(refName)
				if err != nil {
					gctx.NotFound()
					return
				}
			} else {
				gctx.NotFound()
				return
			}
		}

		gctx.Repo.BranchName = refName
		gctx.Data["BranchName"] = gctx.Repo.BranchName
		gctx.Data["CommitID"] = gctx.Repo.CommitID
		gctx.Data["TreePath"] = gctx.Repo.TreePath
		gctx.Data["IsViewBranch"] = gctx.Repo.IsViewBranch
		gctx.Data["IsViewTag"] = gctx.Repo.IsViewTag
		gctx.Data["IsViewCommit"] = gctx.Repo.IsViewCommit

		if gctx.Repo.IsWriter() || (gctx.IsLogged && database.Handle.Repositories().HasForkedBy(c.Request.Context(), gctx.Repo.Repository.ID, gctx.User.ID)) {
			if gctx.Repo.Repository.BaseRepo != nil {
				if gctx.Repo.Repository.BaseRepo.AllowsPulls() {
					gctx.Repo.PullRequest.Allowed = true
					if gctx.Repo.IsWriter() {
						gctx.Data["BaseRepo"] = gctx.Repo.Repository.BaseRepo
						gctx.Repo.PullRequest.BaseRepo = gctx.Repo.Repository.BaseRepo
						gctx.Repo.PullRequest.HeadInfo = gctx.Repo.Owner.Name + ":" + gctx.Repo.BranchName
					} else {
						gctx.Data["BaseRepo"] = gctx.Repo.Repository
						gctx.Repo.PullRequest.BaseRepo = gctx.Repo.Repository
						gctx.Repo.PullRequest.HeadInfo = gctx.User.Name + ":" + gctx.Repo.BranchName
					}
				}
			} else {
				if gctx.Repo.Repository.AllowsPulls() {
					gctx.Data["BaseRepo"] = gctx.Repo.Repository
					gctx.Repo.PullRequest.BaseRepo = gctx.Repo.Repository
					gctx.Repo.PullRequest.Allowed = true
					gctx.Repo.PullRequest.SameRepo = true
					gctx.Repo.PullRequest.HeadInfo = gctx.Repo.BranchName
				}
			}
		}
		gctx.Data["PullRequestCtx"] = gctx.Repo.PullRequest
	}
}

func RequireRepoAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		gctx := c.MustGet("context").(*Context)
		if !gctx.IsLogged || (!gctx.Repo.IsAdmin() && !gctx.User.IsAdmin) {
			gctx.NotFound()
			return
		}
	}
}

func RequireRepoWriter() gin.HandlerFunc {
	return func(c *gin.Context) {
		gctx := c.MustGet("context").(*Context)
		if !gctx.IsLogged || (!gctx.Repo.IsWriter() && !gctx.User.IsAdmin) {
			gctx.NotFound()
			return
		}
	}
}

func GitHookService() gin.HandlerFunc {
	return func(c *gin.Context) {
		gctx := c.MustGet("context").(*Context)
		if !gctx.User.CanEditGitHook() {
			gctx.NotFound()
			return
		}
	}
}

func min(a, b database.AccessMode) database.AccessMode {
	if a < b {
		return a
	}
	return b
}
