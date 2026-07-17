package handler

import (
	"errors"
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/pickup"
	"github.com/Jeff1507/go-to-save-api/internal/core/pickup/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/render"
	"github.com/Jeff1507/go-to-save-api/internal/infra/validation"
	"github.com/go-playground/validator/v10"
)

type CreatePickupTemplateRequest struct {
	Name  string `json:"name" validate:"required,notblank,min=3,max=50"`
	Start string `json:"start" validate:"required,time"`
	End   string `json:"end" validate:"required,time"`
}

type CreatePickupTemplateResponse struct {
	ID string `json:"id"`
}

type CreatePickupTemplate struct {
	v  *validator.Validate
	uc *usecase.CreatePickupTemplate
}

func NewCreatePickupTemplate(
	v *validator.Validate,
	uc *usecase.CreatePickupTemplate,
) *CreatePickupTemplate {
	return &CreatePickupTemplate{
		v:  v,
		uc: uc,
	}
}

func (h *CreatePickupTemplate) Handle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var body CreatePickupTemplateRequest
	if err := render.Bind(r, &body); err != nil {
		render.InvalidRequestBody(w)
		return
	}

	if errs := h.validateRequest(body); errs.HasErrors() {
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

	// TODO: Pegar o storeID do ctx
	storeID := "b875eb25-fc40-48fc-8ed1-daff245deb94"

	output, err := h.uc.Execute(r.Context(), usecase.CreatePickupTemplateInput{
		StoreID: storeID,
		Name:    body.Name,
		Start: usecase.CreatePickupTemplateTimeInput{
			Hour:   startTime.Hour(),
			Minute: startTime.Minute(),
		},
		End: usecase.CreatePickupTemplateTimeInput{
			Hour:   endTime.Hour(),
			Minute: endTime.Minute(),
		},
	})
	if err != nil {
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

	render.Success(w, http.StatusCreated, CreatePickupTemplateResponse{
		ID: output.ID,
	})
}

func (h *CreatePickupTemplate) validateRequest(
	body CreatePickupTemplateRequest,
) validation.Errors {
	errs := validation.NewErrors()

	validation.Struct(h.v, errs, body)

	return errs
}
