package handler

import (
	"net/http"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/core/store/query"
	"github.com/gin-gonic/gin"
)

type GetStoreWithActiveBagsParams struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type PickupWindowResponse struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type BagSummaryResponse struct {
	ID            string               `json:"id"`
	Title         string               `json:"title"`
	Description   *string              `json:"description"`
	Stock         int32                `json:"stock"`
	Discount      int32                `json:"discount"`
	OriginalPrice int64                `json:"originalPrice"`
	FinalPrice    int64                `json:"finalPrice"`
	PickupWindow  PickupWindowResponse `json:"pickupWindow"`
}

type StoreResponse struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	CNPJ        string               `json:"cnpj"`
	AvatarURL   *string              `json:"avatarUrl"`
	Description *string              `json:"description"`
	CreatedAt   time.Time            `json:"createdAt"`
	Rating      RatingResponse       `json:"rating"`
	Address     AddressResponse      `json:"address"`
	Bags        []BagSummaryResponse `json:"bags"`
}

type GetStoreWithActiveBagsHandler struct {
	getStoreWithActiveBags *query.GetStoreWithActiveBags
}

func NewGetStoreWithActiveBagsHandler(
	getStoreWithActiveBags *query.GetStoreWithActiveBags,
) *GetStoreWithActiveBagsHandler {
	return &GetStoreWithActiveBagsHandler{
		getStoreWithActiveBags: getStoreWithActiveBags,
	}
}

func (h *GetStoreWithActiveBagsHandler) Handle(c *gin.Context) {
	var params GetStoreWithActiveBagsParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		return
	}

	s, err := h.getStoreWithActiveBags.Execute(c.Request.Context(), params.ID)
	if err != nil {
		c.Error(err)
		return
	}

	bags := make([]BagSummaryResponse, len(s.Bags))
	for i, b := range s.Bags {
		bags[i] = BagSummaryResponse{
			ID:            b.ID,
			Title:         b.Title,
			Description:   b.Description,
			Stock:         b.Stock,
			Discount:      b.Discount,
			OriginalPrice: b.OriginalPrice,
			FinalPrice:    b.FinalPrice,
			PickupWindow: PickupWindowResponse{
				Start: b.PickupStart,
				End:   b.PickupEnd,
			},
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": StoreResponse{
			ID:          s.ID,
			Name:        s.Name,
			CNPJ:        s.CNPJ,
			AvatarURL:   s.AvatarURL,
			Description: s.Description,
			CreatedAt:   s.CreatedAt,
			Rating:      RatingResponse(s.Rating),
			Address:     AddressResponse(s.Address),
			Bags:        bags,
		},
	})
}
