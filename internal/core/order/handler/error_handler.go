package handler

import (
	"errors"
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/order"
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

	case errors.Is(err, order.ErrInvalidState):
		return http.StatusConflict, "Este pedido já foi processado"

	case errors.Is(err, order.ErrNotFound):
		return http.StatusNotFound, "Pedido não encontrado"

	case errors.Is(err, order.ErrCustomerHasPendingOrder):
		return http.StatusBadRequest, "Você já possui um pedido pendente"

	case errors.Is(err, order.ErrCannotPlaceOrder):
		return http.StatusBadRequest, "Não foi possível realizar o pedido"

	case errors.Is(err, order.ErrInvalidCode):
		return http.StatusUnprocessableEntity, "Código de confirmação inválido"

	case errors.Is(err, order.ErrDuplicateItems):
		return http.StatusBadRequest, "Itens duplicados"

	case errors.Is(err, order.ErrItemNotFound):
		return http.StatusNotFound, "Item não encontrado"

	default:
		return http.StatusInternalServerError, "Erro interno no servidor"
	}
}
