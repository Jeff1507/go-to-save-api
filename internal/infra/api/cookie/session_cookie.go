package cookie

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const cookieName = "session"

func SetSessionCookie(c *gin.Context, sessionID string, sessionTTL time.Duration) {
	c.SetCookieData(&http.Cookie{
		Name:     cookieName,
		Value:    sessionID,
		Path:     "/",
		MaxAge:   int(sessionTTL / time.Second),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func GetSessionCookie(c *gin.Context) (string, error) {
	val, err := c.Cookie(cookieName)
	if err != nil {
		return "", fmt.Errorf("get gin cookie: %w", err)
	}

	return val, nil
}

func ClearSessionCookie(c *gin.Context) {
	c.SetCookieData(&http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}
