package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/customer/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/gin-gonic/gin"
)

type UpdateEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type UpdateEmailHandler struct {
	updateEmail *usecase.UpdateEmail
}

func NewUpdateEmailHandler(updateEmail *usecase.UpdateEmail) *UpdateEmailHandler {
	return &UpdateEmailHandler{
		updateEmail: updateEmail,
	}
}

func (h *UpdateEmailHandler) Handle(c *gin.Context) {
	var body UpdateEmailRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.Error(err)
		return
	}

	id, ok := request.CustomerID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	if err := h.updateEmail.Execute(c.Request.Context(), usecase.UpdateEmailInput{
		ID:    id,
		Email: body.Email,
	}); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
