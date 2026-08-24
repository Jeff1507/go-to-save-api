package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/store/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/gin-gonic/gin"
)

type RemoveAvatarHandler struct {
	removeAvatar *usecase.RemoveAvatar
}

func NewRemoveAvatarHandler(removeAvatar *usecase.RemoveAvatar) *RemoveAvatarHandler {
	return &RemoveAvatarHandler{
		removeAvatar: removeAvatar,
	}
}

func (h *RemoveAvatarHandler) Handle(c *gin.Context) {
	id, ok := request.StoreID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	if err := h.removeAvatar.Execute(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
