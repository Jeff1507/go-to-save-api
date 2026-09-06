package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/bag/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/gin-gonic/gin"
)

type CreateBagRequest struct {
	Title         string              `json:"title" binding:"required,notblank,min=3,max=50"`
	Description   string              `json:"description" binding:"omitempty,notblank,max=255"`
	Stock         int32               `json:"stock" binding:"required,gt=0"`
	OriginalPrice int64               `json:"originalPrice" binding:"required,gt=0"`
	Discount      int                 `json:"discount" binding:"required,gte=30,lte=90"`
	PickupWindow  PickupWindowRequest `json:"pickupWindow" binding:"required"`
	Items         []ItemRequest       `json:"items" binding:"required,min=1,dive"`
}

type CreateBagResponse struct {
	ID string `json:"id"`
}

type CreateBagHandler struct {
	createBag *usecase.CreateBag
}

func NewCreateBagHandler(createBag *usecase.CreateBag) *CreateBagHandler {
	return &CreateBagHandler{
		createBag: createBag,
	}
}

func (h *CreateBagHandler) Handle(c *gin.Context) {
	var body CreateBagRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.Error(err)
		return
	}

	storeID, ok := request.StoreID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	items := make([]usecase.ItemInput, len(body.Items))
	for i, item := range body.Items {
		items[i] = usecase.ItemInput{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	output, err := h.createBag.Execute(c.Request.Context(), usecase.CreateBagInput{
		StoreID:       storeID,
		Title:         body.Title,
		Description:   body.Description,
		Stock:         body.Stock,
		OriginalPrice: body.OriginalPrice,
		Discount:      body.Discount,
		PickupStart:   body.PickupWindow.Start,
		PickupEnd:     body.PickupWindow.End,
		Items:         items,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": CreateBagResponse{
			ID: output.ID,
		},
	})
}
