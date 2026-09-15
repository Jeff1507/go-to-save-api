package handler

import (
	"errors"
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/review"
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
	c.JSON(code, gin.H{
		"error": msg,
	})
}

func domainError(err error) (int, string) {
	switch {
	case errors.Is(err, review.ErrAlreadyReviewed):
		return http.StatusConflict, "Pedido já avaliado"

	case errors.Is(err, review.ErrCannotReview):
		return http.StatusBadRequest, "Não foi possível realizar a avaliação"

	default:
		return http.StatusInternalServerError, "Erro interno no servidor"
	}
}
