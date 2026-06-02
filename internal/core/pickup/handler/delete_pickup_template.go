package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/pickup/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/render"
	"github.com/Jeff1507/go-to-save-api/internal/infra/validation"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type DeletePickupTemplate struct {
	v  *validator.Validate
	uc *usecase.DeletePickupTemplate
}

func NewDeletePickupTemplate(
	v *validator.Validate,
	uc *usecase.DeletePickupTemplate,
) *DeletePickupTemplate {
	return &DeletePickupTemplate{
		v:  v,
		uc: uc,
	}
}

func (h *DeletePickupTemplate) Handle(w http.ResponseWriter, r *http.Request) {
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

func (h *DeletePickupTemplate) validateRequest(id string) validation.Errors {
	errs := validation.NewErrors()

	validation.UUID(h.v, errs, id, "id")

	return errs
}
