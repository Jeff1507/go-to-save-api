package handler

import (
	"net/http"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/core/product/query"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/render"
)

type ProductResponseBody struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	ImageURL  *string   `json:"imageUrl"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ListProductsByStore struct {
	qry *query.ListProductsByStore
}

func NewListProductsByStore(qry *query.ListProductsByStore) *ListProductsByStore {
	return &ListProductsByStore{
		qry: qry,
	}
}

func (h *ListProductsByStore) Handle(w http.ResponseWriter, r *http.Request) {
	// TODO: Pegar o storeID do ctx
	storeID := "b875eb25-fc40-48fc-8ed1-daff245deb94"

	products, err := h.qry.Execute(r.Context(), storeID)
	if err != nil {
		render.InternalServerError(w)
		return
	}

	res := make([]ProductResponseBody, 0, len(products))
	for _, prod := range products {
		res = append(res, ProductResponseBody{
			ID:        prod.ID,
			Name:      prod.Name,
			ImageURL:  prod.ImageURL,
			CreatedAt: prod.CreatedAt,
			UpdatedAt: prod.UpdatedAt,
		})
	}

	render.Success(w, http.StatusOK, res)
}
