package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/review/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/gin-gonic/gin"
)

type AddReviewParams struct {
	OrderID string `uri:"orderId" binding:"required,uuid"`
}

type AddReviewRequest struct {
	Rating  int16  `json:"rating" binding:"required,gte=1,lte=5"`
	Comment string `json:"comment" binding:"omitempty,notblank,max=255"`
}

type AddReviewResponse struct {
	ID string `json:"id"`
}

type AddReviewHandler struct {
	addReview *usecase.AddReview
}

func NewAddReviewHandler(addReview *usecase.AddReview) *AddReviewHandler {
	return &AddReviewHandler{
		addReview: addReview,
	}
}

func (h *AddReviewHandler) Handle(c *gin.Context) {
	var params AddReviewParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		return
	}

	var body AddReviewRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.Error(err)
		return
	}

	customerID, ok := request.CustomerID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	output, err := h.addReview.Execute(c.Request.Context(), usecase.AddReviewInput{
		OrderID:    params.OrderID,
		CustomerID: customerID,
		Rating:     body.Rating,
		Comment:    body.Comment,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": AddReviewResponse{
			ID: output.ID,
		},
	})
}
