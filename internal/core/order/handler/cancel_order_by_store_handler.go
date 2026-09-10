package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/order/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/gin-gonic/gin"
)

type CancelOrderByStoreParams struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type CancelOrderByStoreHandler struct {
	cancelOrderByStore *usecase.CancelOrderByStore
}

func NewCancelOrderByStoreHandler(
	cancelOrderByStore *usecase.CancelOrderByStore,
) *CancelOrderByStoreHandler {
	return &CancelOrderByStoreHandler{
		cancelOrderByStore: cancelOrderByStore,
	}
}

func (h *CancelOrderByStoreHandler) Handle(c *gin.Context) {
	var params CancelOrderByStoreParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		return
	}

	storeID, ok := request.StoreID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	if err := h.cancelOrderByStore.Execute(c.Request.Context(), usecase.CancelOrderByStoreInput{
		ID:      params.ID,
		StoreID: storeID,
	}); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
