package handler

import (
	"net/http"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/core/product/query"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/gin-gonic/gin"
)

type ProductResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	ImageURL  *string   `json:"imageUrl"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ListProductsHandler struct {
	listProducts *query.ListProducts
}

func NewListProductsHandler(listProducts *query.ListProducts) *ListProductsHandler {
	return &ListProductsHandler{
		listProducts: listProducts,
	}
}

func (h *ListProductsHandler) Handle(c *gin.Context) {
	storeID, ok := request.StoreID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	products, err := h.listProducts.Execute(c.Request.Context(), storeID)
	if err != nil {
		c.Error(err)
		return
	}

	response := make([]ProductResponse, 0, len(products))
	for _, prod := range products {
		response = append(response, ProductResponse{
			ID:        prod.ID,
			Name:      prod.Name,
			ImageURL:  prod.ImageURL,
			CreatedAt: prod.CreatedAt,
			UpdatedAt: prod.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}
