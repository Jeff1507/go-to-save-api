package handler

import (
	"net/http"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/core/store/query"
	"github.com/gin-gonic/gin"
)

type SearchStoresQuery struct {
	Name  string `form:"name" binding:"required,min=2"`
	Limit int32  `form:"limit"`
}

type StoreResultResponse struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	CNPJ        string          `json:"cnpj"`
	AvatarURL   *string         `json:"avatarUrl"`
	Description *string         `json:"description"`
	CreatedAt   time.Time       `json:"createdAt"`
	Rating      RatingResponse  `json:"rating"`
	Address     AddressResponse `json:"address"`
}

type SearchStoresHandler struct {
	searchStores *query.SearchStores
}

func NewSearchStoresHandler(searchStores *query.SearchStores) *SearchStoresHandler {
	return &SearchStoresHandler{
		searchStores: searchStores,
	}
}

func (h *SearchStoresHandler) Handle(c *gin.Context) {
	var qry SearchStoresQuery
	if err := c.ShouldBindQuery(&qry); err != nil {
		c.Error(err)
		return
	}

	stores, err := h.searchStores.Execute(c.Request.Context(), query.SearchStoresInput{
		Name:  qry.Name,
		Limit: qry.Limit,
	})
	if err != nil {
		c.Error(err)
		return
	}

	response := make([]StoreResultResponse, 0, len(stores))
	for _, s := range stores {
		response = append(response, StoreResultResponse{
			ID:          s.ID,
			Name:        s.Name,
			CNPJ:        s.CNPJ,
			AvatarURL:   s.AvatarURL,
			Description: s.Description,
			CreatedAt:   s.CreatedAt,
			Rating:      RatingResponse(s.Rating),
			Address:     AddressResponse(s.Address),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}
