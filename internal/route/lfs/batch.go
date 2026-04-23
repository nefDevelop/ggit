package lfs

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "unknwon.dev/clog/v2"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/lfsx"
	"github.com/nef734/ggit/internal/strx"
)

// POST /{owner}/{repo}.git/info/lfs/object/batch
func serveBatch(store Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		owner := c.MustGet("owner").(*database.User)
		repo := c.MustGet("repo").(*database.Repository)
		var request batchRequest
		defer func() { _ = c.Request.Body.Close() }()
		err := json.NewDecoder(c.Request.Body).Decode(&request)
		if err != nil {
			responseJSON(c.Writer, http.StatusBadRequest, responseError{
				Message: strx.ToUpperFirst(err.Error()),
			})
			return
		}

		// NOTE: We only support basic transfer as of now.
		transfer := transferBasic
		// Example: https://try.github.com/nef734/ggit/gogs.git/info/lfs/object/basic
		baseHref := fmt.Sprintf("%s%s/%s.git/info/lfs/objects/basic", conf.Server.ExternalURL, owner.Name, repo.Name)

		objects := make([]batchObject, 0, len(request.Objects))
		switch request.Operation {
		case basicOperationUpload:
			for _, obj := range request.Objects {
				var actions batchActions
				if lfsx.ValidOID(obj.Oid) {
					actions = batchActions{
						Upload: &batchAction{
							Href: fmt.Sprintf("%s/%s", baseHref, obj.Oid),
							Header: map[string]string{
								// NOTE: git-lfs v2.5.0 sets the Content-Type based on the uploaded file.
								// This ensures that the client always uses the designated value for the header.
								"Content-Type": "application/octet-stream",
							},
						},
						Verify: &batchAction{
							Href: fmt.Sprintf("%s/verify", baseHref),
						},
					}
				} else {
					actions = batchActions{
						Error: &batchError{
							Code:    http.StatusUnprocessableEntity,
							Message: "Object has invalid oid",
						},
					}
				}

				objects = append(objects, batchObject{
					Oid:     obj.Oid,
					Size:    obj.Size,
					Actions: actions,
				})
			}

		case basicOperationDownload:
			oids := make([]lfsx.OID, 0, len(request.Objects))
			for _, obj := range request.Objects {
				oids = append(oids, obj.Oid)
			}
			stored, err := store.GetLFSObjectsByOIDs(c.Request.Context(), repo.ID, oids...)
			if err != nil {
				internalServerError(c.Writer)
				log.Error("Failed to get objects [repo_id: %d, oids: %v]: %v", repo.ID, oids, err)
				return
			}
			storedSet := make(map[lfsx.OID]*database.LFSObject, len(stored))
			for _, obj := range stored {
				storedSet[obj.OID] = obj
			}

			for _, obj := range request.Objects {
				var actions batchActions
				if stored := storedSet[obj.Oid]; stored != nil {
					if stored.Size != obj.Size {
						actions.Error = &batchError{
							Code:    http.StatusUnprocessableEntity,
							Message: "Object size mismatch",
						}
					} else {
						actions.Download = &batchAction{
							Href: fmt.Sprintf("%s/%s", baseHref, obj.Oid),
						}
					}
				} else {
					actions.Error = &batchError{
						Code:    http.StatusNotFound,
						Message: "Object does not exist",
					}
				}

				objects = append(objects, batchObject{
					Oid:     obj.Oid,
					Size:    obj.Size,
					Actions: actions,
				})
			}

		default:
			responseJSON(c.Writer, http.StatusBadRequest, responseError{
				Message: "Operation not recognized",
			})
			return
		}

		responseJSON(c.Writer, http.StatusOK, batchResponse{
			Transfer: transfer,
			Objects:  objects,
		})
	}
}

// batchRequest defines the request payload for the batch endpoint.
type batchRequest struct {
	Operation string `json:"operation"`
	Objects   []struct {
		Oid  lfsx.OID `json:"oid"`
		Size int64    `json:"size"`
	} `json:"objects"`
}

type batchError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type batchAction struct {
	Href   string            `json:"href"`
	Header map[string]string `json:"header,omitempty"`
}

type batchActions struct {
	Download *batchAction `json:"download,omitempty"`
	Upload   *batchAction `json:"upload,omitempty"`
	Verify   *batchAction `json:"verify,omitempty"`
	Error    *batchError  `json:"error,omitempty"`
}

type batchObject struct {
	Oid     lfsx.OID     `json:"oid"`
	Size    int64        `json:"size"`
	Actions batchActions `json:"actions"`
}

// batchResponse defines the response payload for the batch endpoint.
type batchResponse struct {
	Transfer string        `json:"transfer"`
	Objects  []batchObject `json:"objects"`
}

type responseError struct {
	Message string `json:"message"`
}

const contentType = "application/vnd.git-lfs+json"

func responseJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		log.Error("Failed to encode JSON: %v", err)
		return
	}
}
