package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/pickup/query"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/render"
)

type StorePickupTemplateResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Start string `json:"start"`
	End   string `json:"end"`
}

type ListStorePickupTemplates struct {
	qry *query.ListStorePickupTemplates
}

func NewListStorePickupTemplates(qry *query.ListStorePickupTemplates) *ListStorePickupTemplates {
	return &ListStorePickupTemplates{
		qry: qry,
	}
}

func (h *ListStorePickupTemplates) Handle(w http.ResponseWriter, r *http.Request) {
	// TODO: Pegar o storeID do ctx
	storeID := "b875eb25-fc40-48fc-8ed1-daff245deb94"

	templates, err := h.qry.Execute(r.Context(), storeID)
	if err != nil {
		render.InternalServerError(w)
		return
	}

	res := make([]StorePickupTemplateResponse, 0, len(templates))
	for _, templ := range templates {
		res = append(res, StorePickupTemplateResponse{
			ID:    templ.ID,
			Name:  templ.Name,
			Start: templ.Start,
			End:   templ.End,
		})
	}

	render.Success(w, http.StatusOK, res)
}
