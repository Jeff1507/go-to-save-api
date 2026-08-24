package handler

import (
	"errors"
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/store"
	"github.com/Jeff1507/go-to-save-api/internal/infra/auth/otp"
	"github.com/Jeff1507/go-to-save-api/internal/infra/validation"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ErrorHandler(c *gin.Context) {
	c.Next()

	if len(c.Errors) == 0 {
		return
	}

	err := c.Errors.Last().Err

	if vErrs, ok := errors.AsType[validator.ValidationErrors](err); ok {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": validation.TranslateValidationErrors(vErrs),
		})
		return
	}

	if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{
			"error": "Arquivo muito grande",
		})
		return
	}

	if errors.Is(err, http.ErrMissingFile) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Arquivo obrigatório",
		})
		return
	}

	code, msg := domainError(err)
	c.JSON(code, gin.H{
		"error": msg,
	})
}

func domainError(err error) (int, string) {
	switch {
	case errors.Is(err, store.ErrEmailAlreadyInUse):
		return http.StatusConflict, "Este email já está em uso"

	case errors.Is(err, store.ErrDuplicateCNPJ):
		return http.StatusConflict, "Este cnpj já está em uso"

	case errors.Is(err, store.ErrInvalidCredentials):
		return http.StatusUnauthorized, "Credenciais inválidas"

	case errors.Is(err, store.ErrNotFound):
		return http.StatusNotFound, "Usuário não encontrado"

	case errors.Is(err, store.ErrInvalidImageFormat):
		return http.StatusUnsupportedMediaType, "Formato de imagem não suportado"

	case errors.Is(err, store.ErrIncorrectPassword):
		return http.StatusUnprocessableEntity, "Senha atual incorreta"

	case errors.Is(err, otp.ErrInvalidOTP):
		return http.StatusUnprocessableEntity, "Código inválido ou expirado"

	default:
		return http.StatusInternalServerError, "Erro interno no servidor"
	}
}
