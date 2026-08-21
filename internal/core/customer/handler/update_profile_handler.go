package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/customer/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/gin-gonic/gin"
)

type UpdateProfileRequest struct {
	FirstName string `json:"firstName" binding:"required,notblank,min=2,max=50"`
	LastName  string `json:"lastName" binding:"required,notblank,min=2,max=50"`
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

	id, ok := request.CustomerID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	if err := h.updateProfile.Execute(c.Request.Context(), usecase.UpdateProfileInput{
		ID:        id,
		FirstName: body.FirstName,
		LastName:  body.LastName,
	}); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
