package middleware

import (
	"net/http"
	"strings"

	"github.com/Jeff1507/go-to-save-api/internal/infra/api/render"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/reqctx"
	"github.com/Jeff1507/go-to-save-api/internal/infra/auth/session"
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

// Mobile: Extrai o token do header
func (a Authentication) extractTokenFromHeader(r *http.Request) string {
	return strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
}

// Web: Extrai o token do cookie
func (a Authentication) extractTokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("session")
	if err != nil {
		return ""
	}

	return cookie.Value
}

func (a *Authentication) CustomerAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := a.extractTokenFromHeader(r)
		if token == "" {
			render.Unauthorized(w)
			return
		}

		customerID, err := a.sessionManager.Validate(r.Context(), token)
		if err != nil {
			render.Unauthorized(w)
			return
		}

		ctx := reqctx.WithCustomer(r.Context(), customerID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *Authentication) StoreAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := a.extractTokenFromCookie(r)
		if token == "" {
			render.Unauthorized(w)
			return
		}

		storeID, err := a.sessionManager.Validate(r.Context(), token)
		if err != nil {
			render.Unauthorized(w)
			return
		}

		ctx := reqctx.WithStore(r.Context(), storeID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
