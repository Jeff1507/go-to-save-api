package handler

import (
	"net/http"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/core/bag/query"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/gin-gonic/gin"
)

type GetBagParams struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type BagResponse struct {
	ID            string               `json:"id"`
	Title         string               `json:"title"`
	Description   *string              `json:"description"`
	Stock         int32                `json:"stock"`
	Discount      int32                `json:"discount"`
	OriginalPrice int64                `json:"originalPrice"`
	FinalPrice    int64                `json:"finalPrice"`
	Status        string               `json:"status"`
	CreatedAt     time.Time            `json:"createdAt"`
	UpdatedAt     time.Time            `json:"updatedAt"`
	PickupWindow  PickupWindowResponse `json:"pickupWindow"`
	Items         []ItemResponse       `json:"items"`
}

type GetBagHandler struct {
	getBag *query.GetBag
}

func NewGetBagHandler(getBag *query.GetBag) *GetBagHandler {
	return &GetBagHandler{
		getBag: getBag,
	}
}

func (h *GetBagHandler) Handle(c *gin.Context) {
	var params GetBagParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		return
	}

	storeID, ok := request.StoreID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	b, err := h.getBag.Execute(c.Request.Context(), query.GetBagInput{
		ID:      params.ID,
		StoreID: storeID,
	})
	if err != nil {
		c.Error(err)
		return
	}

	items := make([]ItemResponse, len(b.Items))
	for i, item := range b.Items {
		items[i] = ItemResponse{
			ID:       item.ID,
			Name:     item.Name,
			ImageURL: item.ImageURL,
			Quantity: item.Quantity,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": BagResponse{
			ID:            b.ID,
			Title:         b.Title,
			Description:   b.Description,
			Stock:         b.Stock,
			Discount:      b.Discount,
			OriginalPrice: b.OriginalPrice,
			FinalPrice:    b.FinalPrice,
			Status:        b.Status,
			CreatedAt:     b.CreatedAt,
			UpdatedAt:     b.UpdatedAt,
			PickupWindow: PickupWindowResponse{
				Start: b.PickupStart,
				End:   b.PickupEnd,
			},
			Items: items,
		},
	})
}
