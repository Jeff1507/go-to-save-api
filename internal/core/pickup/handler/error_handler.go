package handler

import (
	"errors"
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/pickup"
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

	code, msg := domainError(err)
	c.JSON(code, gin.H{"error": msg})
}

func domainError(err error) (int, string) {
	switch {
	case errors.Is(err, pickup.ErrInvalidInterval):
		return http.StatusUnprocessableEntity, "Intervalo de retirada inválido"

	case errors.Is(err, pickup.ErrTemplateAlreadyExists):
		return http.StatusConflict, "Template de retirada já cadastrado"

	case errors.Is(err, pickup.ErrTemplateNotFound):
		return http.StatusNotFound, "Template de retirada não encontrado"

	default:
		return http.StatusInternalServerError, "Erro interno no servidor"
	}
}
