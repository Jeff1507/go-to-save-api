package handler

import (
	"net/http"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/core/order/query"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/gin-gonic/gin"
)

type GetSaleParams struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type SaleResponse struct {
	ID           string               `json:"id"`
	Status       string               `json:"status"`
	TotalPrice   int64                `json:"totalPrice"`
	CreatedAt    time.Time            `json:"createdAt"`
	CompletedAt  *time.Time           `json:"completedAt"`
	PickupWindow PickupWindowResponse `json:"pickupWindow"`
	Customer     CustomerResponse     `json:"customer"`
	Bags         []BagResponse        `json:"bags"`
}

type GetSaleHandler struct {
	getSale *query.GetSale
}

func NewGetSaleHandler(getSale *query.GetSale) *GetSaleHandler {
	return &GetSaleHandler{
		getSale: getSale,
	}
}

func (h *GetSaleHandler) Handle(c *gin.Context) {
	var params GetSaleParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		return
	}

	storeID, ok := request.StoreID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	sale, err := h.getSale.Execute(c.Request.Context(), query.GetSaleInput{
		ID:      params.ID,
		StoreID: storeID,
	})
	if err != nil {
		c.Error(err)
		return
	}

	bags := make([]BagResponse, len(sale.Bags))
	for i, b := range sale.Bags {
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

	c.JSON(http.StatusOK, gin.H{
		"data": SaleResponse{
			ID:          sale.ID,
			Status:      sale.Status,
			TotalPrice:  sale.TotalPrice,
			CreatedAt:   sale.CreatedAt,
			CompletedAt: sale.CompletedAt,
			PickupWindow: PickupWindowResponse{
				Start: sale.PickupStart,
				End:   sale.PickupEnd,
			},
			Customer: CustomerResponse(sale.Customer),
			Bags:     bags,
		},
	})
}
