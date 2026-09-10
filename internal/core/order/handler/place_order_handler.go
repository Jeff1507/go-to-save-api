package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/order/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/gin-gonic/gin"
)

type ItemRequest struct {
	BagID    string `json:"bagId" binding:"required,uuid"`
	Quantity int32  `json:"quantity" binding:"required,gt=0"`
}

type PlaceOrderRequest struct {
	Items []ItemRequest `json:"items" binding:"required,min=1,dive"`
}

type PlaceOrderResponse struct {
	ID string `json:"id"`
}

type PlaceOrderHandler struct {
	placeOrder *usecase.PlaceOrder
}

func NewPlaceOrderHandler(placeOrder *usecase.PlaceOrder) *PlaceOrderHandler {
	return &PlaceOrderHandler{
		placeOrder: placeOrder,
	}
}

func (h *PlaceOrderHandler) Handle(c *gin.Context) {
	var body PlaceOrderRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.Error(err)
		return
	}

	customerID, ok := request.CustomerID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	items := make([]usecase.ItemInput, len(body.Items))
	for i, item := range body.Items {
		items[i] = usecase.ItemInput{
			BagID:    item.BagID,
			Quantity: item.Quantity,
		}
	}

	output, err := h.placeOrder.Execute(c.Request.Context(), usecase.PlaceOrderInput{
		CustomerID: customerID,
		Items:      items,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": PlaceOrderResponse{
			ID: output.ID,
		},
	})
}
