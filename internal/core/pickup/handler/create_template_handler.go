package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/pickup/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/gin-gonic/gin"
)

type CreateTemplateRequest struct {
	Name  string `json:"name" binding:"required,notblank,min=3,max=50"`
	Start string `json:"start" binding:"required,time"`
	End   string `json:"end" binding:"required,time"`
}

type CreateTemplateResponse struct {
	ID string `json:"id"`
}

type CreateTemplateHandler struct {
	createTemplate *usecase.CreateTemplate
}

func NewCreateTemplateHandler(createTemplate *usecase.CreateTemplate) *CreateTemplateHandler {
	return &CreateTemplateHandler{
		createTemplate: createTemplate,
	}
}

func (h *CreateTemplateHandler) Handle(c *gin.Context) {
	var body CreateTemplateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.Error(err)
		return
	}

	storeID, ok := request.StoreID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	start := parseTime(body.Start)
	end := parseTime(body.End)

	output, err := h.createTemplate.Execute(c.Request.Context(), usecase.CreateTemplateInput{
		StoreID: storeID,
		Name:    body.Name,
		Start: usecase.CreateTemplateTimeInput{
			Hour:   start.Hour(),
			Minute: start.Minute(),
		},
		End: usecase.CreateTemplateTimeInput{
			Hour:   end.Hour(),
			Minute: end.Minute(),
		},
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": CreateTemplateResponse{
			ID: output.ID,
		},
	})
}
