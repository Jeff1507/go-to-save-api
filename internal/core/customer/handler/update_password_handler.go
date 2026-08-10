package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/customer/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/gin-gonic/gin"
)

type UpdatePasswordRequest struct {
	Password    string `json:"password" binding:"required,min=8,max=72"`
	NewPassword string `json:"newPassword" binding:"required,min=8,max=72"`
}

type UpdatePasswordHandler struct {
	updatePassword *usecase.UpdatePassword
}

func NewUpdatePasswordHandler(updatePassword *usecase.UpdatePassword) *UpdatePasswordHandler {
	return &UpdatePasswordHandler{
		updatePassword: updatePassword,
	}
}

func (h *UpdatePasswordHandler) Handle(c *gin.Context) {
	var body UpdatePasswordRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.Error(err)
		return
	}

	id, ok := request.CustomerID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	if err := h.updatePassword.Execute(c.Request.Context(), usecase.UpdatePasswordInput{
		ID:          id,
		Password:    body.Password,
		NewPassword: body.NewPassword,
	}); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
