package handler

import (
	"net/http"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/core/store/query"
	"github.com/gin-gonic/gin"
)

type ListNearbyStoresQuery struct {
	Latitude  float64 `form:"latitude" binding:"required,latitude"`
	Longitude float64 `form:"longitude" binding:"required,longitude"`
	Radius    int32   `form:"radius" binding:"required,gt=0"`
}

type StoreSummaryResponse struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	CNPJ         string          `json:"cnpj"`
	AvatarURL    *string         `json:"avatarUrl"`
	Description  *string         `json:"description"`
	BagsCount    int32           `json:"bagsCount"`
	BagsMinPrice int64           `json:"bagsMinPrice"`
	DistanceKM   float64         `json:"distanceKm"`
	CreatedAt    time.Time       `json:"createdAt"`
	Rating       RatingResponse  `json:"rating"`
	Address      AddressResponse `json:"address"`
}

type ListNearbyStoresHandler struct {
	listNearbyStores *query.ListNearbyStores
}

func NewListNearbyStoresHandler(
	listNearbyStores *query.ListNearbyStores,
) *ListNearbyStoresHandler {
	return &ListNearbyStoresHandler{
		listNearbyStores: listNearbyStores,
	}
}

func (h *ListNearbyStoresHandler) Handle(c *gin.Context) {
	var qry ListNearbyStoresQuery
	if err := c.ShouldBindQuery(&qry); err != nil {
		c.Error(err)
		return
	}

	stores, err := h.listNearbyStores.Execute(c.Request.Context(), query.ListNearbyStoresInput{
		Latitude:  qry.Latitude,
		Longitude: qry.Longitude,
		Radius:    qry.Radius,
	})
	if err != nil {
		c.Error(err)
		return
	}

	response := make([]StoreSummaryResponse, 0, len(stores))
	for _, s := range stores {
		response = append(response, StoreSummaryResponse{
			ID:           s.ID,
			Name:         s.Name,
			CNPJ:         s.CNPJ,
			AvatarURL:    s.AvatarURL,
			Description:  s.Description,
			BagsCount:    s.BagsCount,
			BagsMinPrice: s.BagsMinPrice,
			DistanceKM:   s.DistanceKM,
			CreatedAt:    s.CreatedAt,
			Rating:       RatingResponse(s.Rating),
			Address:      AddressResponse(s.Address),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}
