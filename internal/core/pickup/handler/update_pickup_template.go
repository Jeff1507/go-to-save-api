package handler

import (
	"errors"
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/pickup"
	"github.com/Jeff1507/go-to-save-api/internal/core/pickup/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/render"
	"github.com/Jeff1507/go-to-save-api/internal/infra/validation"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type UpdatePickupTemplateRequestBody struct {
	Name  string `json:"name" validate:"required,notblank,min=3,max=50"`
	Start string `json:"start" validate:"required,time"`
	End   string `json:"end" validate:"required,time"`
}

type UpdatePickupTemplate struct {
	v  *validator.Validate
	uc *usecase.UpdatePickupTemplate
}

func NewUpdatePickupTemplate(
	v *validator.Validate,
	uc *usecase.UpdatePickupTemplate,
) *UpdatePickupTemplate {
	return &UpdatePickupTemplate{
		v:  v,
		uc: uc,
	}
}

func (h *UpdatePickupTemplate) Handle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id := chi.URLParam(r, "id")

	var body UpdatePickupTemplateRequestBody
	if err := render.Bind(r, &body); err != nil {
		render.InvalidRequestBody(w)
		return
	}

	if errs := h.validateRequest(id, body); errs.HasErrors() {
		render.ValidationError(w, errs)
		return
	}

	startTime, err := parseTime(body.Start)
	if err != nil {
		render.InternalServerError(w)
		return
	}

	endTime, err := parseTime(body.End)
	if err != nil {
		render.InternalServerError(w)
		return
	}

	if err := h.uc.Execute(r.Context(), usecase.UpdatePickupTemplateInput{
		ID:   id,
		Name: body.Name,
		Start: usecase.UpdatePickupTemplateTimeInput{
			Hour:   startTime.Hour(),
			Minute: startTime.Minute(),
		},
		End: usecase.UpdatePickupTemplateTimeInput{
			Hour:   endTime.Hour(),
			Minute: endTime.Minute(),
		},
	}); err != nil {
		if errors.Is(err, pickup.ErrTemplateNotFound) {
			render.Error(w, http.StatusNotFound, "Template de retirada não encontrado")
			return
		}

		if errors.Is(err, pickup.ErrTemplateAlreadyExists) {
			render.Error(w, http.StatusConflict, "Template de retirada já cadastrado")
			return
		}

		if errors.Is(err, pickup.ErrInvalidInterval) {
			render.Error(w, http.StatusUnprocessableEntity, "Intervalo de retirada inválido")
			return
		}

		render.InternalServerError(w)
		return
	}

	render.NoContent(w)
}

func (h *UpdatePickupTemplate) validateRequest(
	id string,
	body UpdatePickupTemplateRequestBody,
) validation.Errors {
	errs := validation.NewErrors()

	validation.UUID(h.v, errs, id, "id")

	validation.Struct(h.v, errs, body)

	return errs
}
