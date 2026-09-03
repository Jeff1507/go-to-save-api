package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/preset/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/gin-gonic/gin"
)

type UpdatePresetParams struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type UpdatePresetRequest struct {
	Name  string `json:"name" binding:"required,notblank,min=3,max=50"`
	Start string `json:"start" binding:"required,time"`
	End   string `json:"end" binding:"required,time"`
}

type UpdatePresetHandler struct {
	updatePreset *usecase.UpdatePreset
}

func NewUpdatePresetHandler(updatePreset *usecase.UpdatePreset) *UpdatePresetHandler {
	return &UpdatePresetHandler{
		updatePreset: updatePreset,
	}
}

func (h *UpdatePresetHandler) Handle(c *gin.Context) {
	var params UpdatePresetParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		return
	}

	var body UpdatePresetRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.Error(err)
		return
	}

	storeID, ok := request.StoreID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	if err := h.updatePreset.Execute(c.Request.Context(), usecase.UpdatePresetInput{
		ID:      params.ID,
		StoreID: storeID,
		Name:    body.Name,
		Start:   body.Start,
		End:     body.End,
	}); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
