package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/store/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/cookie"
	"github.com/Jeff1507/go-to-save-api/internal/infra/auth/session"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type LoginHandler struct {
	login *usecase.Login
}

func NewLoginHandler(login *usecase.Login) *LoginHandler {
	return &LoginHandler{
		login: login,
	}
}

func (h *LoginHandler) Handle(c *gin.Context) {
	var body LoginRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.Error(err)
		return
	}

	output, err := h.login.Execute(c.Request.Context(), usecase.LoginInput{
		Email:    body.Email,
		Password: body.Password,
	})
	if err != nil {
		c.Error(err)
		return
	}

	cookie.SetSessionCookie(c, output.SessionID, session.TTL)

	c.Status(http.StatusNoContent)
}
