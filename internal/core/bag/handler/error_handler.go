package handler

import (
	"errors"
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/bag"
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
	c.JSON(code, gin.H{
		"error": msg,
	})
}

func domainError(err error) (int, string) {
	switch {
	case errors.Is(err, pickup.ErrInvalidWindowInterval):
		return http.StatusUnprocessableEntity, "Intervalo da janela de retirada inválido"

	case errors.Is(err, pickup.ErrWindowStartInPast):
		return http.StatusBadRequest, "A janela de retirada deve iniciar no futuro"

	case errors.Is(err, bag.ErrNotFound):
		return http.StatusNotFound, "Sacola não encontrada"

	case errors.Is(err, bag.ErrClosed):
		return http.StatusUnprocessableEntity, "Sacola já encerrada"

	case errors.Is(err, bag.ErrExpired):
		return http.StatusUnprocessableEntity, "Sacola expirada"

	case errors.Is(err, bag.ErrInsufficientStock):
		return http.StatusUnprocessableEntity, "Estoque da sacola insuficiente"

	case errors.Is(err, bag.ErrDuplicateItems):
		return http.StatusBadRequest, "Itens duplicados"

	case errors.Is(err, bag.ErrItemNotFound):
		return http.StatusNotFound, "Item não encontrado"

	case errors.Is(err, bag.ErrHasPendingOrder):
		return http.StatusNotFound, "A sacola não pode ser editada com um pedido pendente"

	case errors.Is(err, bag.ErrZeroDelta):
		return http.StatusBadRequest, "O ajuste de estoque não pode ser zero"

	default:
		return http.StatusInternalServerError, "Erro interno no servidor"
	}
}
