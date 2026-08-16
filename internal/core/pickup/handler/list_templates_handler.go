package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/pickup/query"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/gin-gonic/gin"
)

type TemplateResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Start string `json:"start"`
	End   string `json:"end"`
}

type ListTemplatesHandler struct {
	listTemplates *query.ListTemplates
}

func NewListTemplatesHandler(listTemplates *query.ListTemplates) *ListTemplatesHandler {
	return &ListTemplatesHandler{
		listTemplates: listTemplates,
	}
}

func (h *ListTemplatesHandler) Handle(c *gin.Context) {
	storeID, ok := request.StoreID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	templates, err := h.listTemplates.Execute(c.Request.Context(), storeID)
	if err != nil {
		c.Error(err)
		return
	}

	response := make([]TemplateResponse, 0, len(templates))
	for _, templ := range templates {
		response = append(response, TemplateResponse{
			ID:    templ.ID,
			Name:  templ.Name,
			Start: templ.Start,
			End:   templ.End,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}
