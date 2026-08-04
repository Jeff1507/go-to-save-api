package handler

import (
	"errors"
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/customer"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/render"
	"github.com/Jeff1507/go-to-save-api/internal/infra/auth/otp"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func ErrorHandler(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			handleError(w, err)
		}
	}
}

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, customer.ErrEmailAlreadyInUse):
		render.Error(w, http.StatusConflict, "Conta já cadastrada")
	case errors.Is(err, customer.ErrInvalidCredentials):
		render.Error(w, http.StatusUnauthorized, "Credenciais inválidas")
	case errors.Is(err, customer.ErrNotFound):
		render.Error(w, http.StatusNotFound, "Usuário não encontrado")
	case errors.Is(err, customer.ErrInvalidImageFormat):
		render.Error(w, http.StatusUnsupportedMediaType, "Formato de imagem não suportado")
	case errors.Is(err, otp.ErrInvalidOTP):
		render.Error(w, http.StatusUnprocessableEntity, "Código inválido ou expirado")
	default:
		render.InternalServerError(w)
	}
}
