package lfs

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	log "unknwon.dev/clog/v2"

	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/lfsx"
	"github.com/nef734/ggit/internal/strx"
)

const transferBasic = "basic"
const (
	basicOperationUpload   = "upload"
	basicOperationDownload = "download"
)

type basicHandler struct {
	store Store
	// The default storage backend for uploading new objects.
	defaultStorage lfsx.Storage
	// The list of available storage backends to access objects.
	storagers map[lfsx.Storage]lfsx.Storager
}

// DefaultStorager returns the default storage backend.
func (h *basicHandler) DefaultStorager() lfsx.Storager {
	return h.storagers[h.defaultStorage]
}

// Storager returns the given storage backend.
func (h *basicHandler) Storager(storage lfsx.Storage) lfsx.Storager {
	return h.storagers[storage]
}

// GET /{owner}/{repo}.git/info/lfs/object/basic/{oid}
func (h *basicHandler) serveDownload(c *gin.Context) {
	repo := c.MustGet("repo").(*database.Repository)
	oid := c.MustGet("oid").(lfsx.OID)
	object, err := h.store.GetLFSObjectByOID(c.Request.Context(), repo.ID, oid)
	if err != nil {
		if database.IsErrLFSObjectNotExist(err) {
			responseJSON(c.Writer, http.StatusNotFound, responseError{
				Message: "Object does not exist",
			})
		} else {
			internalServerError(c.Writer)
			log.Error("Failed to get object [repo_id: %d, oid: %s]: %v", repo.ID, oid, err)
		}
		return
	}

	s := h.Storager(object.Storage)
	if s == nil {
		internalServerError(c.Writer)
		log.Error("Failed to locate the object [repo_id: %d, oid: %s]: storage %q not found", object.RepoID, object.OID, object.Storage)
		return
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", strconv.FormatInt(object.Size, 10))
	c.Status(http.StatusOK)

	err = s.Download(object.OID, c.Writer)
	if err != nil {
		log.Error("Failed to download object [oid: %s]: %v", object.OID, err)
		return
	}
}

// PUT /{owner}/{repo}.git/info/lfs/object/basic/{oid}
func (h *basicHandler) serveUpload(c *gin.Context) {
	repo := c.MustGet("repo").(*database.Repository)
	oid := c.MustGet("oid").(lfsx.OID)

	// NOTE: LFS client will retry upload the same object if there was a partial failure,
	// therefore we would like to skip ones that already exist.
	_, err := h.store.GetLFSObjectByOID(c.Request.Context(), repo.ID, oid)
	if err == nil {
		// Object exists, drain the request body and we're good.
		_, _ = io.Copy(io.Discard, c.Request.Body)
		c.Request.Body.Close()
		c.Status(http.StatusOK)
		return
	} else if !database.IsErrLFSObjectNotExist(err) {
		internalServerError(c.Writer)
		log.Error("Failed to get object [repo_id: %d, oid: %s]: %v", repo.ID, oid, err)
		return
	}

	s := h.DefaultStorager()
	written, tmpPath, err := s.UploadTemp(oid, c.Request.Body)
	if err != nil {
		if err == lfsx.ErrInvalidOID || err == lfsx.ErrOIDMismatch {
			responseJSON(c.Writer, http.StatusBadRequest, responseError{
				Message: err.Error(),
			})
		} else {
			internalServerError(c.Writer)
			log.Error("Failed to upload object [storage: %s, oid: %s]: %v", s.Storage(), oid, err)
		}
		return
	}

	if tmpPath != "" {
		defer os.Remove(tmpPath)
		finalPath := s.(*lfsx.LocalStorage).StoragePath(oid)
		err = database.Handle.LFS().CreateObjectAtomic(c.Request.Context(), repo.ID, oid, written, s.Storage(), tmpPath, finalPath)
	} else {
		err = h.store.CreateLFSObject(c.Request.Context(), repo.ID, oid, written, s.Storage())
	}

	if err != nil {
		internalServerError(c.Writer)
		log.Error("Failed to create object [repo_id: %d, oid: %s]: %v", repo.ID, oid, err)
		return
	}
	c.Status(http.StatusOK)

	log.Trace("[LFS] Object created %q", oid)
}

// POST /{owner}/{repo}.git/info/lfs/object/basic/verify
func (h *basicHandler) serveVerify(c *gin.Context) {
	repo := c.MustGet("repo").(*database.Repository)
	var request basicVerifyRequest
	defer func() { _ = c.Request.Body.Close() }()

	err := json.NewDecoder(c.Request.Body).Decode(&request)
	if err != nil {
		responseJSON(c.Writer, http.StatusBadRequest, responseError{
			Message: strx.ToUpperFirst(err.Error()),
		})
		return
	}

	if !lfsx.ValidOID(request.Oid) {
		responseJSON(c.Writer, http.StatusBadRequest, responseError{
			Message: "Invalid oid",
		})
		return
	}

	object, err := h.store.GetLFSObjectByOID(c.Request.Context(), repo.ID, request.Oid)
	if err != nil {
		if database.IsErrLFSObjectNotExist(err) {
			responseJSON(c.Writer, http.StatusNotFound, responseError{
				Message: "Object does not exist",
			})
		} else {
			internalServerError(c.Writer)
			log.Error("Failed to get object [repo_id: %d, oid: %s]: %v", repo.ID, request.Oid, err)
		}
		return
	}

	if object.Size != request.Size {
		responseJSON(c.Writer, http.StatusBadRequest, responseError{
			Message: "Object size mismatch",
		})
		return
	}
	c.Status(http.StatusOK)
}

type basicVerifyRequest struct {
	Oid  lfsx.OID `json:"oid"`
	Size int64    `json:"size"`
}
