package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/order/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/gin-gonic/gin"
)

type CompleteOrderParams struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type CompleteOrderRequest struct {
	ConfirmationCode string `json:"confirmationCode" binding:"required,numeric,len=4"`
}

type CompleteOrderHandler struct {
	completeOrder *usecase.CompleteOrder
}

func NewCompleteOrderHandler(completeOrder *usecase.CompleteOrder) *CompleteOrderHandler {
	return &CompleteOrderHandler{
		completeOrder: completeOrder,
	}
}

func (h *CompleteOrderHandler) Handle(c *gin.Context) {
	var params CompleteOrderParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		return
	}

	var body CompleteOrderRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.Error(err)
		return
	}

	storeID, ok := request.StoreID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	if err := h.completeOrder.Execute(c.Request.Context(), usecase.CompleteOrderInput{
		ID:               params.ID,
		StoreID:          storeID,
		ConfirmationCode: body.ConfirmationCode,
	}); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
