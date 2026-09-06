package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/bag/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/gin-gonic/gin"
)

type CloseBagParams struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type CloseBagHandler struct {
	closeBag *usecase.CloseBag
}

func NewCloseBagHandler(closeBag *usecase.CloseBag) *CloseBagHandler {
	return &CloseBagHandler{
		closeBag: closeBag,
	}
}

func (h *CloseBagHandler) Handle(c *gin.Context) {
	var params CloseBagParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		return
	}

	storeID, ok := request.StoreID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	if err := h.closeBag.Execute(c.Request.Context(), usecase.CloseBagInput{
		ID:      params.ID,
		StoreID: storeID,
	}); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
