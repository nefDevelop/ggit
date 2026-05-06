package v2

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gogs/git-module"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/osx"
	"github.com/nef734/ggit/internal/tool"
	log "unknwon.dev/clog/v2"
)

// GetRepoArchive serves the repository as a zip or tar.gz archive.
func GetRepoArchive(c *context.APIContext) {
	archiveName := c.Params("name")
	if archiveName == "" {
		archiveName = c.Params("*")
	}
	var (
		ext           string
		archivePath   string
		archiveFormat git.ArchiveFormat
	)

	switch {
	case strings.HasSuffix(archiveName, ".zip"):
		ext = ".zip"
		archivePath = filepath.Join(c.Repo.GitRepo.Path(), "archives", "zip")
		archiveFormat = git.ArchiveZip
	case strings.HasSuffix(archiveName, ".tar.gz"):
		ext = ".tar.gz"
		archivePath = filepath.Join(c.Repo.GitRepo.Path(), "archives", "targz")
		archiveFormat = git.ArchiveTarGz
	default:
		c.ErrorStatus(http.StatusBadRequest, fmt.Errorf("Unsupported archive format: %s", archiveName))
		return
	}

	refName := strings.TrimSuffix(archiveName, ext)
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

	log.Trace("[DEBUG] Resolving archive commit for ref: %s (Repo ID: %d)", refName, c.Repo.Repository.ID)

	// Try branch first
	commit, err = gitRepo.BranchCommit(refName)
	if err == nil {
		log.Trace("[DEBUG] Resolved as branch: %s", refName)
	} else {
		// Try tag
		commit, err = gitRepo.TagCommit(refName)
		if err == nil {
			log.Trace("[DEBUG] Resolved as tag: %s", refName)
		}
	}
	
	if err != nil {
		log.Trace("[DEBUG] Not found in Git branches/tags, trying release fallback for: %s", refName)
		// Fallback: Check if there is a release with this tag name and use its stored SHA.
		rel, err2 := database.GetRelease(c.Repo.Repository.ID, refName)
		if err2 == nil && rel.Sha1 != "" {
			log.Trace("[DEBUG] Found release in DB with SHA: %s", rel.Sha1)
			commit, err = gitRepo.CatFileCommit(rel.Sha1)
			if err == nil {
				log.Trace("[DEBUG] Resolved via release SHA: %s", rel.Sha1)
			}
		} else if err2 != nil {
			log.Trace("[DEBUG] Release not found in DB: %v", err2)
		} else {
			log.Trace("[DEBUG] Release found but has no SHA")
		}
	}

	if err != nil && (len(refName) >= 7 && len(refName) <= 40) {
		log.Trace("[DEBUG] Trying as direct SHA: %s", refName)
		commit, err = gitRepo.CatFileCommit(refName)
	}

	if err != nil {
		log.Trace("[DEBUG] Failed to resolve commit for %s: %v", refName, err)
		c.Error(err, "get commit")
		return
	}

	archivePath = path.Join(archivePath, tool.ShortSHA1(commit.ID.String())+ext)
	if !osx.IsFile(archivePath) {
		if err := commit.CreateArchive(archiveFormat, archivePath); err != nil {
			c.Error(err, "create archive")
			return
		}
	}

	c.ServeFile(archivePath, c.Repo.Repository.Name+"-"+refName+ext)
}

// GetRepoRawFile serves the raw content of a file.
func GetRepoRawFile(c *context.APIContext) {
	serveRawFile(c, false)
}

// GetRepoMediaFile serves the raw content of a file as a download.
func GetRepoMediaFile(c *context.APIContext) {
	serveRawFile(c, true)
}

func serveRawFile(c *context.APIContext, isMedia bool) {
	refName := c.Params("ref")
	treePath := strings.TrimPrefix(c.Params("path"), "/")
	if treePath == "" {
		treePath = strings.TrimPrefix(c.Params("*"), "/")
	}

	gitRepo := c.Repo.GitRepo
	var (
		commit *git.Commit
		err    error
	)
	
	if gitRepo.HasBranch(refName) {
		commit, err = gitRepo.BranchCommit(refName)
	} else if gitRepo.HasTag(refName) {
		commit, err = gitRepo.TagCommit(refName)
	} else {
		commit, err = gitRepo.CatFileCommit(refName)
	}

	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	blob, err := commit.Blob(treePath)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	data, err := blob.Bytes()
	if err != nil {
		c.Error(err, "get blob bytes")
		return
	}

	// Set headers
	if isMedia {
		c.Resp.Header().Set("Content-Disposition", "attachment; filename=\""+path.Base(treePath)+"\"")
		c.Resp.Header().Set("Content-Transfer-Encoding", "binary")
	} else {
		if !tool.IsTextFile(data) {
			if !tool.IsImageFile(data) {
				c.Resp.Header().Set("Content-Disposition", "attachment; filename=\""+path.Base(treePath)+"\"")
				c.Resp.Header().Set("Content-Transfer-Encoding", "binary")
			}
		} else if !conf.Repository.EnableRawFileRenderMode || c.Query("render") != "true" {
			c.Resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
		}
	}

	c.Resp.Header().Set("Last-Modified", commit.Committer.When.Format(http.TimeFormat))
	if _, err := c.Resp.Write(data); err != nil {
		c.Error(err, "write response")
	}
}
