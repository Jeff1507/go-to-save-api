package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/customer/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/Jeff1507/go-to-save-api/internal/infra/storage"
	"github.com/gin-gonic/gin"
)

type UpdateAvatarHandler struct {
	updateAvatar *usecase.UpdateAvatar
}

func NewUpdateAvatarHandler(updateAvatar *usecase.UpdateAvatar) *UpdateAvatarHandler {
	return &UpdateAvatarHandler{
		updateAvatar: updateAvatar,
	}
}

func (h *UpdateAvatarHandler) Handle(c *gin.Context) {
	file, _, err := c.Request.FormFile("avatar")
	if err != nil {
		c.Error(err)
		return
	}
	defer file.Close()

	if c.Request.MultipartForm != nil {
		defer c.Request.MultipartForm.RemoveAll()
	}

	contentType, err := storage.GetContentType(file)
	if err != nil {
		c.Error(err)
		return
	}

	id, ok := request.CustomerID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	if err := h.updateAvatar.Execute(c.Request.Context(), usecase.UpdateAvatarInput{
		ID: id,
		Avatar: usecase.AvatarInput{
			Content:     file,
			ContentType: contentType,
		},
	}); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
