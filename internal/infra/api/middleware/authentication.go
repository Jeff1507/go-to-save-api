package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Jeff1507/go-to-save-api/internal/infra/api/cookie"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/Jeff1507/go-to-save-api/internal/infra/auth/session"
	"github.com/gin-gonic/gin"
)

// TODO: Separar as sessões de Customer e Store com role
type Authentication struct {
	sessions session.Manager
}

func NewAuthentication(sessions session.Manager) *Authentication {
	return &Authentication{
		sessions: sessions,
	}
}

func (a Authentication) tokenFromHeader(c *gin.Context) string {
	return strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
}

func (a Authentication) tokenFromCookie(c *gin.Context) string {
	if token, err := cookie.GetSessionCookie(c); err == nil {
		return token
	}

	return ""
}

func (a *Authentication) AuthenticateCustomer(c *gin.Context) {
	sessionID := a.tokenFromHeader(c)
	if sessionID == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	customerID, err := a.sessions.Validate(c.Request.Context(), sessionID)
	if err != nil {
		if errors.Is(err, session.ErrInvalidSession) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Erro interno no servidor",
		})
		return
	}

	ctx := request.WithCustomerID(c.Request.Context(), customerID)
	ctx = request.WithSessionID(ctx, sessionID)

	c.Request = c.Request.WithContext(ctx)
	c.Next()
}

func (a *Authentication) AuthenticateStore(c *gin.Context) {
	sessionID := a.tokenFromCookie(c)
	if sessionID == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	storeID, err := a.sessions.Validate(c.Request.Context(), sessionID)
	if err != nil {
		if errors.Is(err, session.ErrInvalidSession) {
			cookie.ClearSessionCookie(c)

			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Erro interno no servidor",
		})
		return
	}

	ctx := request.WithStoreID(c.Request.Context(), storeID)

	c.Request = c.Request.WithContext(ctx)
	c.Next()
}
