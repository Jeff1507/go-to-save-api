package handler

import (
	"errors"
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/product"
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

	code, msg := domainError(err)
	c.JSON(code, gin.H{
		"error": msg,
	})
}

func domainError(err error) (int, string) {
	switch {
	case errors.Is(err, product.ErrAlreadyExists):
		return http.StatusConflict, "Produto já cadastrado"

	case errors.Is(err, product.ErrNotFound):
		return http.StatusNotFound, "Produto não encontrado"

	case errors.Is(err, product.ErrInvalidImageFormat):
		return http.StatusUnsupportedMediaType, "Formato de imagem não suportado"

	default:
		return http.StatusInternalServerError, "Erro interno no servidor"
	}
}
