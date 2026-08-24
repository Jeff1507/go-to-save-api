package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/store/usecase"
	"github.com/gin-gonic/gin"
)

type VerifyEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,numeric,len=6"`
}

type VerifyEmailHandler struct {
	verifyEmail *usecase.VerifyEmail
}

func NewVerifyEmailHandler(verifyEmail *usecase.VerifyEmail) *VerifyEmailHandler {
	return &VerifyEmailHandler{
		verifyEmail: verifyEmail,
	}
}

func (h *VerifyEmailHandler) Handle(c *gin.Context) {
	var body VerifyEmailRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.Error(err)
		return
	}

	if err := h.verifyEmail.Execute(c.Request.Context(), usecase.VerifyEmailInput{
		Email: body.Email,
		Code:  body.Code,
	}); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
