package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/preset/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/gin-gonic/gin"
)

type DeletePresetParams struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type DeletePresetHandler struct {
	deletePreset *usecase.DeletePreset
}

func NewDeletePresetHandler(deletePreset *usecase.DeletePreset) *DeletePresetHandler {
	return &DeletePresetHandler{
		deletePreset: deletePreset,
	}
}

func (h *DeletePresetHandler) Handle(c *gin.Context) {
	var params DeletePresetParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		return
	}

	storeID, ok := request.StoreID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	if err := h.deletePreset.Execute(c.Request.Context(), usecase.DeletePresetInput{
		ID:      params.ID,
		StoreID: storeID,
	}); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
