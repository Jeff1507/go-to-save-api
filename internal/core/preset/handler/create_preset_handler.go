package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/preset/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/gin-gonic/gin"
)

type CreatePresetRequest struct {
	Name  string `json:"name" binding:"required,notblank,min=3,max=50"`
	Start string `json:"start" binding:"required,time"`
	End   string `json:"end" binding:"required,time"`
}

type CreatePresetResponse struct {
	ID string `json:"id"`
}

type CreatePresetHandler struct {
	createPreset *usecase.CreatePreset
}

func NewCreatePresetHandler(createPreset *usecase.CreatePreset) *CreatePresetHandler {
	return &CreatePresetHandler{
		createPreset: createPreset,
	}
}

func (h *CreatePresetHandler) Handle(c *gin.Context) {
	var body CreatePresetRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.Error(err)
		return
	}

	storeID, ok := request.StoreID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	output, err := h.createPreset.Execute(c.Request.Context(), usecase.CreatePresetInput{
		StoreID: storeID,
		Name:    body.Name,
		Start:   body.Start,
		End:     body.End,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": CreatePresetResponse{
			ID: output.ID,
		},
	})
}
