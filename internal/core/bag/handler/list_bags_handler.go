package handler

import (
	"net/http"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/core/bag/query"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/gin-gonic/gin"
)

type BagSummaryResponse struct {
	ID            string               `json:"id"`
	Title         string               `json:"storeId"`
	Description   *string              `json:"description"`
	Stock         int32                `json:"stock"`
	Discount      int32                `json:"discount"`
	OriginalPrice int64                `json:"originalPrice"`
	FinalPrice    int64                `json:"finalPrice"`
	Status        string               `json:"status"`
	PickupWindow  PickupWindowResponse `json:"pickupWindow"`
	CreatedAt     time.Time            `json:"createdAt"`
	UpdatedAt     time.Time            `json:"updatedAt"`
}

type ListBagsHandler struct {
	listBags *query.ListBags
}

func NewListBagsHandler(listBags *query.ListBags) *ListBagsHandler {
	return &ListBagsHandler{
		listBags: listBags,
	}
}

func (h *ListBagsHandler) Handle(c *gin.Context) {
	storeID, ok := request.StoreID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	bags, err := h.listBags.Execute(c.Request.Context(), storeID)
	if err != nil {
		c.Error(err)
		return
	}

	response := make([]BagSummaryResponse, 0, len(bags))
	for _, b := range bags {
		response = append(response, BagSummaryResponse{
			ID:            b.ID,
			Title:         b.Title,
			Description:   b.Description,
			Stock:         b.Stock,
			Discount:      b.Discount,
			OriginalPrice: b.OriginalPrice,
			FinalPrice:    b.FinalPrice,
			Status:        b.Status,
			PickupWindow: PickupWindowResponse{
				Start: b.PickupStart,
				End:   b.PickupEnd,
			},
			CreatedAt: b.CreatedAt,
			UpdatedAt: b.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}
