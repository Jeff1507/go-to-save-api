package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/bag/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/gin-gonic/gin"
)

type UpdateBagParams struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type UpdateBagRequest struct {
	Title         string              `json:"title" binding:"required,notblank,min=3,max=50"`
	Description   string              `json:"description" binding:"omitempty,notblank,max=255"`
	OriginalPrice int64               `json:"originalPrice" binding:"required,gt=0"`
	Discount      int                 `json:"discount" binding:"required,gte=30,lte=90"`
	PickupWindow  PickupWindowRequest `json:"pickupWindow" binding:"required"`
	Items         []ItemRequest       `json:"items" binding:"required,min=1,dive"`
}

type UpdateBagHandler struct {
	updateBag *usecase.UpdateBag
}

func NewUpdateBagHandler(updateBag *usecase.UpdateBag) *UpdateBagHandler {
	return &UpdateBagHandler{
		updateBag: updateBag,
	}
}

func (h *UpdateBagHandler) Handle(c *gin.Context) {
	var params UpdateBagParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		return
	}

	var body UpdateBagRequest
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

	if err := h.updateBag.Execute(c.Request.Context(), usecase.UpdateBagInput{
		ID:            params.ID,
		StoreID:       storeID,
		Title:         body.Title,
		Description:   body.Description,
		OriginalPrice: body.OriginalPrice,
		Discount:      body.Discount,
		PickupStart:   body.PickupWindow.Start,
		PickupEnd:     body.PickupWindow.End,
		Items:         items,
	}); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
