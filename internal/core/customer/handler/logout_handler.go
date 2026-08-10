package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/customer/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/gin-gonic/gin"
)

type LogoutHandler struct {
	logout *usecase.Logout
}

func NewLogoutHandler(logout *usecase.Logout) *LogoutHandler {
	return &LogoutHandler{
		logout: logout,
	}
}

func (h *LogoutHandler) Handle(c *gin.Context) {
	sessionID, ok := request.SessionID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	if err := h.logout.Execute(c.Request.Context(), sessionID); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
