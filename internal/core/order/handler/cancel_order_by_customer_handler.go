package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/order/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/gin-gonic/gin"
)

type CancelOrderByCustomerParams struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type CancelOrderByCustomerHandler struct {
	cancelOrderByCustomer *usecase.CancelOrderByCustomer
}

func NewCancelOrderByCustomerHandler(
	cancelOrderByCustomer *usecase.CancelOrderByCustomer,
) *CancelOrderByCustomerHandler {
	return &CancelOrderByCustomerHandler{
		cancelOrderByCustomer: cancelOrderByCustomer,
	}
}

func (h *CancelOrderByCustomerHandler) Handle(c *gin.Context) {
	var params CancelOrderByCustomerParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		return
	}

	customerID, ok := request.CustomerID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	if err := h.cancelOrderByCustomer.Execute(c.Request.Context(), usecase.CancelOrderByCustomerInput{
		ID:         params.ID,
		CustomerID: customerID,
	}); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
