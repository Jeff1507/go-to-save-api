package handler

import (
	"errors"
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/bag"
	"github.com/Jeff1507/go-to-save-api/internal/core/bag/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/render"
	"github.com/Jeff1507/go-to-save-api/internal/infra/validation"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type CloseBag struct {
	v  *validator.Validate
	uc *usecase.CloseBag
}

func NewCloseBag(v *validator.Validate, uc *usecase.CloseBag) *CloseBag {
	return &CloseBag{
		v:  v,
		uc: uc,
	}
}

func (h *CloseBag) Handle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if errs := h.validateRequest(id); errs.HasErrors() {
		render.ValidationError(w, errs)
		return
	}

	if err := h.uc.Execute(r.Context(), id); err != nil {
		if errors.Is(err, bag.ErrNotFound) {
			render.Error(w, http.StatusNotFound, "Sacola não encontrada")
			return
		}

		if errors.Is(err, bag.ErrInvalidOperation) {
			render.Error(w, http.StatusForbidden, "Não é possível realizar essa ação")
			return
		}

		render.InternalServerError(w)
		return
	}

	render.NoContent(w)
}

func (h *CloseBag) validateRequest(id string) validation.Errors {
	errs := validation.NewErrors()

	validation.UUID(h.v, errs, id, "id")

	return errs
}
