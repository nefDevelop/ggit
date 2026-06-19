package repo

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gogs/git-module"
	"github.com/nef734/ggit/internal/context"
	log "github.com/nef734/ggit/internal/log"
	"github.com/nef734/ggit/internal/osx"
	"github.com/nef734/ggit/internal/tool"
)

var archiveLocks sync.Map

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

	// Prevent concurrent writes to the same archive path (TOCTOU race)
	v, _ := archiveLocks.LoadOrStore(archivePath, &sync.Mutex{})
	mu := v.(*sync.Mutex)
	mu.Lock()
	if !osx.IsFile(archivePath) {
		if err := commit.CreateArchive(archiveFormat, archivePath); err != nil {
			mu.Unlock()
			c.Error(err, "creates archive")
			return
		}
	}
	mu.Unlock()

	c.ServeFile(archivePath, c.Repo.Repository.Name+"-"+refName+ext)
}
