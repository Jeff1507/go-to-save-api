package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/preset/query"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/gin-gonic/gin"
)

type PresetResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Start string `json:"start"`
	End   string `json:"end"`
}

type ListPresetsHandler struct {
	listPresets *query.ListPresets
}

func NewListPresetsHandler(listPresets *query.ListPresets) *ListPresetsHandler {
	return &ListPresetsHandler{
		listPresets: listPresets,
	}
}

func (h *ListPresetsHandler) Handle(c *gin.Context) {
	storeID, ok := request.StoreID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	presets, err := h.listPresets.Execute(c.Request.Context(), storeID)
	if err != nil {
		c.Error(err)
		return
	}

	response := make([]PresetResponse, 0, len(presets))
	for _, pset := range presets {
		response = append(response, PresetResponse{
			ID:    pset.ID,
			Name:  pset.Name,
			Start: pset.Start,
			End:   pset.End,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}
