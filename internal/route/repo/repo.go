package repo

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	log "unknwon.dev/clog/v2"

	"github.com/gogs/git-module"

	"gogs.io/gogs/internal/conf"
	"gogs.io/gogs/internal/context"
	"gogs.io/gogs/internal/database"
	"gogs.io/gogs/internal/form"
	"gogs.io/gogs/internal/osx"
	"gogs.io/gogs/internal/services/repo"
	"gogs.io/gogs/internal/tool"
	"gogs.io/gogs/internal/urlx"
)

func Action(c *context.Context) {
	var err error
	switch c.Params(":action") {
	case "watch":
		err = database.WatchRepo(c.User.ID, c.Repo.Repository.ID, true)
	case "unwatch":
		if userID := c.QueryInt64("user_id"); userID != 0 {
			if c.User.IsAdmin {
				err = database.WatchRepo(userID, c.Repo.Repository.ID, false)
			}
		} else {
			err = database.WatchRepo(c.User.ID, c.Repo.Repository.ID, false)
		}
	case "star":
		err = database.StarRepo(c.User.ID, c.Repo.Repository.ID, true)
	case "unstar":
		err = database.StarRepo(c.User.ID, c.Repo.Repository.ID, false)
	case "desc": // FIXME: this is not used
		if !c.Repo.IsOwner() {
			c.NotFound()
			return
		}

		c.Repo.Repository.Description = c.Query("desc")
		c.Repo.Repository.Website = c.Query("site")
		err = database.UpdateRepository(c.Repo.Repository, false)
	}

	if err != nil {
		c.Error(err, fmt.Sprintf("action %q", c.Params(":action")))
		return
	}

	redirectTo := c.Query("redirect_to")
	if !urlx.IsSameSite(redirectTo) {
		redirectTo = c.Repo.RepoLink
	}
	c.Redirect(redirectTo)
}

func Download(c *context.Context) {
	var (
		uri           = c.Params("*")
		refName       string
		ext           string
		archivePath   string
		archiveFormat git.ArchiveFormat
	)

	switch {
	case strings.HasSuffix(uri, ".zip"):
		ext = ".zip"
		archivePath = filepath.Join(c.Repo.GitRepo.Path(), "archives", "zip")
		archiveFormat = git.ArchiveZip
	case strings.HasSuffix(uri, ".tar.gz"):
		ext = ".tar.gz"
		archivePath = filepath.Join(c.Repo.GitRepo.Path(), "archives", "targz")
		archiveFormat = git.ArchiveTarGz
	default:
		log.Trace("Unknown format: %s", uri)
		c.NotFound()
		return
	}
	refName = strings.TrimSuffix(uri, ext)

	if !osx.IsDir(archivePath) {
		if err := os.MkdirAll(archivePath, os.ModePerm); err != nil {
			c.Error(err, "create archive directory")
			return
		}
	}

	// Get corresponding commit.
	var (
		commit *git.Commit
		err    error
	)
	gitRepo := c.Repo.GitRepo
	if gitRepo.HasBranch(refName) {
		commit, err = gitRepo.BranchCommit(refName)
		if err != nil {
			c.Error(err, "get branch commit")
			return
		}
	} else if gitRepo.HasTag(refName) {
		commit, err = gitRepo.TagCommit(refName)
		if err != nil {
			c.Error(err, "get tag commit")
			return
		}
	} else if len(refName) >= 7 && len(refName) <= 40 {
		commit, err = gitRepo.CatFileCommit(refName)
		if err != nil {
			c.NotFound()
			return
		}
	} else {
		c.NotFound()
		return
	}

	archivePath = path.Join(archivePath, tool.ShortSHA1(commit.ID.String())+ext)
	if !osx.IsFile(archivePath) {
		if err := commit.CreateArchive(archiveFormat, archivePath); err != nil {
			c.Error(err, "creates archive")
			return
		}
	}

	c.ServeFile(archivePath, c.Repo.Repository.Name+"-"+refName+ext)
}
