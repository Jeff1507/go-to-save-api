package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/product/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/render"
	"github.com/Jeff1507/go-to-save-api/internal/infra/validation"
	"github.com/go-chi/chi/v5"
)

type DeleteProduct struct {
	uc *usecase.DeleteProduct
}

func NewDeleteProduct(uc *usecase.DeleteProduct) *DeleteProduct {
	return &DeleteProduct{
		uc: uc,
	}
}

func (h *DeleteProduct) Handle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if errs := h.validate(id); errs.HasErrors() {
		render.ValidationError(w, errs)
		return
	}

	if err := h.uc.Execute(r.Context(), id); err != nil {
		render.InternalServerError(w)
		return
	}

	render.NoContent(w)
}

func (h DeleteProduct) validate(id string) validation.Errors {
	errs := validation.NewErrors()

	validation.String(errs, "id", id).
		Required().
		UUID()

	return errs
}
