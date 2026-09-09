package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/store/usecase"
	"github.com/gin-gonic/gin"
)

type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Code        string `json:"code" binding:"required,numeric,len=6"`
	NewPassword string `json:"newPassword" binding:"required,min=8,max=72"`
}

type ResetPasswordHandler struct {
	resetPassword *usecase.ResetPassword
}

func NewResetPasswordHandler(resetPassword *usecase.ResetPassword) *ResetPasswordHandler {
	return &ResetPasswordHandler{
		resetPassword: resetPassword,
	}
}

func (h *ResetPasswordHandler) Handle(c *gin.Context) {
	var body ResetPasswordRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.Error(err)
		return
	}

	if err := h.resetPassword.Execute(c.Request.Context(), usecase.ResetPasswordInput{
		Email:       body.Email,
		Code:        body.Code,
		NewPassword: body.NewPassword,
	}); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
