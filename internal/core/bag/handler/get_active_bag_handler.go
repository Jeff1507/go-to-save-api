package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/bag/query"
	"github.com/gin-gonic/gin"
)

type GetActiveBagParams struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type ActiveBagResponse struct {
	ID            string               `json:"id"`
	StoreID       string               `json:"storeId"`
	Title         string               `json:"title"`
	Description   *string              `json:"description"`
	Stock         int32                `json:"stock"`
	Discount      int32                `json:"discount"`
	OriginalPrice int64                `json:"originalPrice"`
	FinalPrice    int64                `json:"finalPrice"`
	PickupWindow  PickupWindowResponse `json:"pickupWindow"`
	Items         []ItemResponse       `json:"items"`
}

type GetActiveBagHandler struct {
	getActiveBag *query.GetActiveBag
}

func NewGetActiveBagHandler(getActiveBag *query.GetActiveBag) *GetActiveBagHandler {
	return &GetActiveBagHandler{
		getActiveBag: getActiveBag,
	}
}

func (h *GetActiveBagHandler) Handle(c *gin.Context) {
	var params GetActiveBagParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		return
	}

	active, err := h.getActiveBag.Execute(c.Request.Context(), params.ID)
	if err != nil {
		c.Error(err)
		return
	}

	items := make([]ItemResponse, len(active.Items))
	for i, item := range active.Items {
		items[i] = ItemResponse{
			ID:       item.ID,
			Name:     item.Name,
			ImageURL: item.ImageURL,
			Quantity: item.Quantity,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": ActiveBagResponse{
			ID:            active.ID,
			StoreID:       active.StoreID,
			Title:         active.Title,
			Description:   active.Description,
			Stock:         active.Stock,
			Discount:      active.Discount,
			OriginalPrice: active.OriginalPrice,
			FinalPrice:    active.FinalPrice,
			PickupWindow: PickupWindowResponse{
				Start: active.PickupStart,
				End:   active.PickupEnd,
			},
			Items: items,
		},
	})
}
