package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/pickup/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/gin-gonic/gin"
)

type UpdateTemplateParams struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type UpdateTemplateRequest struct {
	Name  string `json:"name" binding:"required,notblank,min=3,max=50"`
	Start string `json:"start" binding:"required,time"`
	End   string `json:"end" binding:"required,time"`
}

type UpdateTemplateHandler struct {
	updateTemplate *usecase.UpdateTemplate
}

func NewUpdateTemplateHandler(updateTemplate *usecase.UpdateTemplate) *UpdateTemplateHandler {
	return &UpdateTemplateHandler{
		updateTemplate: updateTemplate,
	}
}

func (h *UpdateTemplateHandler) Handle(c *gin.Context) {
	var params UpdateTemplateParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		return
	}

	var body UpdateTemplateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.Error(err)
		return
	}

	start := parseTime(body.Start)
	end := parseTime(body.End)

	storeID, ok := request.StoreID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	if err := h.updateTemplate.Execute(c.Request.Context(), usecase.UpdateTemplateInput{
		ID:      params.ID,
		StoreID: storeID,
		Name:    body.Name,
		Start: usecase.UpdateTemplateTimeInput{
			Hour:   start.Hour(),
			Minute: start.Minute(),
		},
		End: usecase.UpdateTemplateTimeInput{
			Hour:   end.Hour(),
			Minute: end.Minute(),
		},
	}); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
