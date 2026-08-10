package middleware

import (
	"net/http"
	"strings"

	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/Jeff1507/go-to-save-api/internal/infra/auth/session"
	"github.com/gin-gonic/gin"
)

// TODO: Separar as sessões de Customer e Store com role
type Authentication struct {
	sessionManager session.Manager
}

func NewAuthentication(sessionManager session.Manager) *Authentication {
	return &Authentication{
		sessionManager: sessionManager,
	}
}

func (a Authentication) extractTokenFromHeader(c *gin.Context) string {
	return strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
}

func (a *Authentication) VerifyCustomer(c *gin.Context) {
	sessionID := a.extractTokenFromHeader(c)
	if sessionID == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	customerID, err := a.sessionManager.Validate(c.Request.Context(), sessionID)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx := request.WithCustomerID(c.Request.Context(), customerID)
	ctx = request.WithSessionID(ctx, sessionID)

	c.Request = c.Request.WithContext(ctx)

	c.Next()
}
