package v2

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/dbtest"
)

func NewTestAPIContext(req *http.Request) (*context.APIContext, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	
	gctx := &context.Context{
		Context: c,
		Data:    make(map[string]any),
		Req:     req,
	}
	
	return &context.APIContext{Context: gctx}, w
}

func TestUpdateUserAvatar(t *testing.T) {
	db := dbtest.NewDB(t, "UpdateUserAvatar", new(database.User))
	database.SetHandle(db)
	
	user := &database.User{ID: 1, Name: "alice", LowerName: "alice", Email: "alice@example.com"}
	require.NoError(t, db.Create(user).Error)

	t.Run("Superación: Subir Avatar Mágico (Simulado)", func(t *testing.T) {
		// Crear una imagen PNG en memoria
		img := image.NewRGBA(image.Rect(0, 0, 10, 10))
		var buf bytes.Buffer
		err := png.Encode(&buf, img)
		require.NoError(t, err)

		// Preparar request multipart
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("avatar", "avatar.png")
		require.NoError(t, err)
		_, err = io.Copy(part, &buf)
		require.NoError(t, err)
		require.NoError(t, writer.Close())

		req, _ := http.NewRequest(http.MethodPost, "/api/v2/user/settings/avatar", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		c, w := NewTestAPIContext(req)
		c.User = user
		
		UpdateUserAvatar(c)
		assert.Equal(t, http.StatusOK, w.Code)
		
		// Verificar que el usuario ahora tiene avatar personalizado
		updatedUser, _ := database.Handle.Users().GetByID(c.Req.Context(), user.ID)
		assert.True(t, updatedUser.UseCustomAvatar)
	})

	t.Run("Fallo: Archivo vacío", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/api/v2/user/settings/avatar", nil)
		c, w := NewTestAPIContext(req)
		c.User = user
		
		UpdateUserAvatar(c)
		assert.Equal(t, http.StatusOK, w.Code, "Debe devolver OK incluso si no hay archivo (pero no cambia nada)")
	})
}
