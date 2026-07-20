package handler

import (
	"net/http"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/core/bag/query"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/render"
)

type StoreBagPickupWindowResponse struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type StoreBagItemResponse struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	ImageURL *string `json:"imageUrl"`
	Quantity int32   `json:"quantity"`
}

type StoreBagResponse struct {
	ID            string                       `json:"id"`
	Title         string                       `json:"title"`
	Description   *string                      `json:"description"`
	Stock         int32                        `json:"stock"`
	OriginalPrice int64                        `json:"originalPrice"`
	FinalPrice    int64                        `json:"finalPrice"`
	Status        string                       `json:"status"`
	PickupWindow  StoreBagPickupWindowResponse `json:"pickupWindow"`
	Items         []StoreBagItemResponse       `json:"items"`
	CreatedAt     time.Time                    `json:"createdAt"`
	UpdatedAt     time.Time                    `json:"updatedAt"`
}

type ListStoreBags struct {
	qry *query.ListStoreBags
}

func NewListStoreBags(qry *query.ListStoreBags) *ListStoreBags {
	return &ListStoreBags{
		qry: qry,
	}
}

func (h *ListStoreBags) Handle(w http.ResponseWriter, r *http.Request) {
	// TODO: Pegar o storeID do ctx
	storeID := "b875eb25-fc40-48fc-8ed1-daff245deb94"

	bags, err := h.qry.Execute(r.Context(), storeID)
	if err != nil {
		render.InternalServerError(w)
		return
	}

	res := make([]StoreBagResponse, 0, len(bags))
	for _, b := range bags {
		pickupWindow := StoreBagPickupWindowResponse{
			Start: b.PickupWindow.Start,
			End:   b.PickupWindow.End,
		}

		items := make([]StoreBagItemResponse, 0, len(b.Items))
		for _, item := range b.Items {
			items = append(items, StoreBagItemResponse{
				ID:       item.ID,
				Name:     item.Name,
				ImageURL: item.ImageURL,
				Quantity: item.Quantity,
			})
		}

		res = append(res, StoreBagResponse{
			ID:            b.ID,
			Title:         b.Title,
			Description:   b.Description,
			Stock:         b.Stock,
			OriginalPrice: b.OriginalPrice,
			FinalPrice:    b.FinalPrice,
			Status:        b.Status,
			PickupWindow:  pickupWindow,
			Items:         items,
			CreatedAt:     b.CreatedAt,
			UpdatedAt:     b.UpdatedAt,
		})
	}

	render.Success(w, http.StatusOK, res)
}
