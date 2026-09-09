package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/store/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/gin-gonic/gin"
)

type UpdateProfileRequest struct {
	Name        string `json:"name" binding:"required,notblank,min=2,max=100"`
	Description string `json:"description" binding:"omitempty,notblank,max=512"`
}

type UpdateProfileHandler struct {
	updateProfile *usecase.UpdateProfile
}

func NewUpdateProfileHandler(updateProfile *usecase.UpdateProfile) *UpdateProfileHandler {
	return &UpdateProfileHandler{
		updateProfile: updateProfile,
	}
}

func (h *UpdateProfileHandler) Handle(c *gin.Context) {
	var body UpdateProfileRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.Error(err)
		return
	}

	id, ok := request.StoreID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	if err := h.updateProfile.Execute(c.Request.Context(), usecase.UpdateProfileInput{
		ID:          id,
		Name:        body.Name,
		Description: body.Description,
	}); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
