package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/bag/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/gin-gonic/gin"
)

type UpdateStockParams struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type UpdateStockRequest struct {
	Delta int32 `json:"delta" binding:"required"`
}

type UpdateStockHandler struct {
	updateStock *usecase.UpdateStock
}

func NewUpdateStockHandler(updateStock *usecase.UpdateStock) *UpdateStockHandler {
	return &UpdateStockHandler{
		updateStock: updateStock,
	}
}

func (h *UpdateStockHandler) Handle(c *gin.Context) {
	var params UpdateStockParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		return
	}

	var body UpdateStockRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.Error(err)
		return
	}

	storeID, ok := request.StoreID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	if err := h.updateStock.Execute(c.Request.Context(), usecase.UpdateStockInput{
		ID:      params.ID,
		StoreID: storeID,
		Delta:   body.Delta,
	}); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
