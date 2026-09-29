package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/bag/query"
	"github.com/gin-gonic/gin"
)

type GetReservableBagParams struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type ReservableBagResponse struct {
	ID            string               `json:"id"`
	Title         string               `json:"title"`
	Description   *string              `json:"description"`
	Stock         int32                `json:"stock"`
	Discount      int32                `json:"discount"`
	OriginalPrice int64                `json:"originalPrice"`
	FinalPrice    int64                `json:"finalPrice"`
	PickupWindow  PickupWindowResponse `json:"pickupWindow"`
	Items         []ItemResponse       `json:"items"`
}

type GetReservableBagHandler struct {
	getReservableBag *query.GetReservableBag
}

func NewGetReservableBagHandler(
	getReservableBag *query.GetReservableBag,
) *GetReservableBagHandler {
	return &GetReservableBagHandler{
		getReservableBag: getReservableBag,
	}
}

func (h *GetReservableBagHandler) Handle(c *gin.Context) {
	var params GetReservableBagParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		return
	}

	b, err := h.getReservableBag.Execute(c.Request.Context(), params.ID)
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
		"data": ReservableBagResponse{
			ID:            b.ID,
			Title:         b.Title,
			Description:   b.Description,
			Stock:         b.Stock,
			Discount:      b.Discount,
			OriginalPrice: b.OriginalPrice,
			FinalPrice:    b.FinalPrice,
			PickupWindow: PickupWindowResponse{
				Start: b.PickupStart,
				End:   b.PickupEnd,
			},
			Items: items,
		},
	})
}
