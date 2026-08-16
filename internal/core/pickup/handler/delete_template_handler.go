package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/pickup/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/gin-gonic/gin"
)

type DeleteTemplateParams struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type DeleteTemplateHandler struct {
	deleteTemplate *usecase.DeleteTemplate
}

func NewDeleteTemplateHandler(deleteTemplate *usecase.DeleteTemplate) *DeleteTemplateHandler {
	return &DeleteTemplateHandler{
		deleteTemplate: deleteTemplate,
	}
}

func (h *DeleteTemplateHandler) Handle(c *gin.Context) {
	var params DeleteTemplateParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		return
	}

	storeID, ok := request.StoreID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	if err := h.deleteTemplate.Execute(c.Request.Context(), usecase.DeleteTemplateInput{
		ID:      params.ID,
		StoreID: storeID,
	}); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
