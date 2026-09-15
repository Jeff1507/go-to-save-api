package handler

import (
	"net/http"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/core/order/query"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/gin-gonic/gin"
)

type OrderSummaryResponse struct {
	ID           string               `json:"id"`
	Status       string               `json:"status"`
	TotalPrice   int64                `json:"totalPrice"`
	CreatedAt    time.Time            `json:"createdAt"`
	CompletedAt  *time.Time           `json:"completedAt"`
	PickupWindow PickupWindowResponse `json:"pickupWindow"`
	Store        StoreResponse        `json:"store"`
}

type ListOrdersHandler struct {
	listOrders *query.ListOrders
}

func NewListOrdersHandler(listOrders *query.ListOrders) *ListOrdersHandler {
	return &ListOrdersHandler{
		listOrders: listOrders,
	}
}

func (h *ListOrdersHandler) Handle(c *gin.Context) {
	customerID, ok := request.CustomerID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	orders, err := h.listOrders.Execute(c.Request.Context(), customerID)
	if err != nil {
		c.Error(err)
		return
	}

	response := make([]OrderSummaryResponse, 0, len(orders))
	for _, ord := range orders {
		response = append(response, OrderSummaryResponse{
			ID:          ord.ID,
			Status:      ord.Status,
			TotalPrice:  ord.TotalPrice,
			CreatedAt:   ord.CreatedAt,
			CompletedAt: ord.CompletedAt,
			PickupWindow: PickupWindowResponse{
				Start: ord.PickupStart,
				End:   ord.PickupEnd,
			},
			Store: StoreResponse{
				ID:        ord.Store.ID,
				Name:      ord.Store.Name,
				AvatarURL: ord.Store.AvatarURL,
				Rating:    RatingResponse(ord.Store.Rating),
				Address:   AddressResponse(ord.Store.Address),
			},
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}
