package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/store/usecase"
	"github.com/gin-gonic/gin"
)

type RequestPasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type RequestPasswordResetHandler struct {
	requestPasswordReset *usecase.RequestPasswordReset
}

func NewRequestPasswordResetHandler(
	requestPasswordReset *usecase.RequestPasswordReset,
) *RequestPasswordResetHandler {
	return &RequestPasswordResetHandler{
		requestPasswordReset: requestPasswordReset,
	}
}

func (h *RequestPasswordResetHandler) Handle(c *gin.Context) {
	var body RequestPasswordResetRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.Error(err)
		return
	}

	if err := h.requestPasswordReset.Execute(c.Request.Context(), body.Email); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusAccepted)
}
