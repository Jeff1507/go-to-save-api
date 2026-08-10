package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/customer/usecase"
	"github.com/gin-gonic/gin"
)

type RequestEmailVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type RequestEmailVerificationHandler struct {
	requestEmailVerification *usecase.RequestEmailVerification
}

func NewRequestEmailVerificationHandler(
	requestEmailVerification *usecase.RequestEmailVerification,
) *RequestEmailVerificationHandler {
	return &RequestEmailVerificationHandler{
		requestEmailVerification: requestEmailVerification,
	}
}

func (h *RequestEmailVerificationHandler) Handle(c *gin.Context) {
	var body RequestEmailVerificationRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.Error(err)
		return
	}

	if err := h.requestEmailVerification.Execute(c.Request.Context(), body.Email); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusAccepted)
}
