package handler

import (
	"net/http"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/core/order/query"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/gin-gonic/gin"
)

type GetOrderParams struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type ReviewResponse struct {
	ID         string    `json:"id"`
	Rating     int16     `json:"rating"`
	Comment    *string   `json:"comment"`
	ReviewedAt time.Time `json:"reviewedAt"`
}

type OrderResponse struct {
	ID               string               `json:"id"`
	Status           string               `json:"status"`
	TotalPrice       int64                `json:"totalPrice"`
	ConfirmationCode string               `json:"confirmationCode"`
	CreatedAt        time.Time            `json:"createdAt"`
	CompletedAt      *time.Time           `json:"completedAt"`
	PickupWindow     PickupWindowResponse `json:"pickupWindow"`
	Store            StoreResponse        `json:"store"`
	Bags             []BagResponse        `json:"bags"`
	Review           *ReviewResponse      `json:"review"`
}

type GetOrderHandler struct {
	getOrder *query.GetOrder
}

func NewGetOrderHandler(getOrder *query.GetOrder) *GetOrderHandler {
	return &GetOrderHandler{
		getOrder: getOrder,
	}
}

func (h *GetOrderHandler) Handle(c *gin.Context) {
	var params GetOrderParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		return
	}

	customerID, ok := request.CustomerID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	ord, err := h.getOrder.Execute(c.Request.Context(), query.GetOrderInput{
		ID:         params.ID,
		CustomerID: customerID,
	})
	if err != nil {
		c.Error(err)
		return
	}

	bags := make([]BagResponse, len(ord.Bags))
	for i, b := range ord.Bags {
		bags[i] = BagResponse{
			ID:        b.ID,
			Title:     b.Title,
			UnitPrice: b.UnitPrice,
			Quantity:  b.Quantity,
		}

		items := make([]BagItemResponse, len(b.Items))
		for i, item := range b.Items {
			items[i] = BagItemResponse{
				ID:       item.ID,
				Name:     item.Name,
				ImageURL: item.ImageURL,
				Quantity: item.Quantity,
			}
		}

		bags[i].Items = items
	}

	var review *ReviewResponse
	if ord.Review != nil {
		review = &ReviewResponse{
			ID:         ord.Review.ID,
			Rating:     ord.Review.Rating,
			Comment:    ord.Review.Comment,
			ReviewedAt: ord.Review.ReviewedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": OrderResponse{
			ID:               ord.ID,
			Status:           ord.Status,
			TotalPrice:       ord.TotalPrice,
			ConfirmationCode: ord.ConfirmationCode,
			CreatedAt:        ord.CreatedAt,
			CompletedAt:      ord.CompletedAt,
			PickupWindow: PickupWindowResponse{
				Start: ord.PickupStart,
				End:   ord.PickupEnd,
			},
			Store: StoreResponse{
				ID:        ord.Store.ID,
				Name:      ord.Store.Name,
				AvatarURL: ord.Store.AvatarURL,
				Address:   AddressResponse(ord.Store.Address),
			},
			Bags:   bags,
			Review: review,
		},
	})
}
