package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/product/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/render"
	"github.com/Jeff1507/go-to-save-api/internal/infra/validation"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type DeleteProduct struct {
	v  *validator.Validate
	uc *usecase.DeleteProduct
}

func NewDeleteProduct(v *validator.Validate, uc *usecase.DeleteProduct) *DeleteProduct {
	return &DeleteProduct{
		v:  v,
		uc: uc,
	}
}

func (h *DeleteProduct) Handle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if errs := h.validateRequest(id); errs.HasErrors() {
		render.ValidationError(w, errs)
		return
	}

	if err := h.uc.Execute(r.Context(), id); err != nil {
		render.InternalServerError(w)
		return
	}

	render.NoContent(w)
}

func (h *DeleteProduct) validateRequest(id string) validation.Errors {
	errs := validation.NewErrors()

	validation.UUID(h.v, errs, id, "id")

	return errs
}
