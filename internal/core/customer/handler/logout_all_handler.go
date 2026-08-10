package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/customer/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/gin-gonic/gin"
)

type LogoutAllHandler struct {
	logoutAll *usecase.LogoutAll
}

func NewLogoutAllHandler(logoutAll *usecase.LogoutAll) *LogoutAllHandler {
	return &LogoutAllHandler{
		logoutAll: logoutAll,
	}
}

func (h *LogoutAllHandler) Handle(c *gin.Context) {
	id, ok := request.CustomerID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	if err := h.logoutAll.Execute(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
