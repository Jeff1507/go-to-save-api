package handler

import (
	"net/http"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/core/order/query"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/gin-gonic/gin"
)

type ListSalesQuery struct {
	Status string `form:"status"`
}

type SaleSummaryResponse struct {
	ID           string               `json:"id"`
	Status       string               `json:"status"`
	TotalPrice   int64                `json:"totalPrice"`
	CreatedAt    time.Time            `json:"createdAt"`
	CompletedAt  *time.Time           `json:"completedAt"`
	PickupWindow PickupWindowResponse `json:"pickupWindow"`
	Customer     CustomerResponse     `json:"customer"`
}

type ListSalesHandler struct {
	listSales *query.ListSales
}

func NewListSalesHandler(listSales *query.ListSales) *ListSalesHandler {
	return &ListSalesHandler{
		listSales: listSales,
	}
}

func (h *ListSalesHandler) Handle(c *gin.Context) {
	var qry ListSalesQuery
	if err := c.ShouldBindQuery(&qry); err != nil {
		c.Error(err)
		return
	}

	storeID, ok := request.StoreID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	sales, err := h.listSales.Execute(c.Request.Context(), query.ListSalesInput{
		StoreID: storeID,
		Status:  qry.Status,
	})
	if err != nil {
		c.Error(err)
		return
	}

	response := make([]SaleSummaryResponse, 0, len(sales))
	for _, sale := range sales {
		response = append(response, SaleSummaryResponse{
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
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}
