package handler

import (
	"errors"
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/pickup"
	"github.com/Jeff1507/go-to-save-api/internal/core/pickup/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/render"
	"github.com/Jeff1507/go-to-save-api/internal/infra/validation"
)

type CreatePickupTemplateRequestBody struct {
	Name  string `json:"name"`
	Start string `json:"start"`
	End   string `json:"end"`
}

type CreatePickupTemplateResponseBody struct {
	ID string `json:"id"`
}

type CreatePickupTemplate struct {
	uc *usecase.CreatePickupTemplate
}

func NewCreatePickupTemplate(uc *usecase.CreatePickupTemplate) *CreatePickupTemplate {
	return &CreatePickupTemplate{
		uc: uc,
	}
}

func (h *CreatePickupTemplate) Handle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var body CreatePickupTemplateRequestBody
	if err := render.Bind(r, &body); err != nil {
		render.InvalidRequestBody(w)
		return
	}

	if errs := h.validate(body); errs.HasErrors() {
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
	storeID := "store-id"

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

	render.Success(w, http.StatusCreated, CreatePickupTemplateResponseBody{
		ID: output.ID,
	})
}

func (h CreatePickupTemplate) validate(body CreatePickupTemplateRequestBody) validation.Errors {
	errs := validation.NewErrors()

	validation.String(errs, "name", body.Name).
		Required().
		Min(3).
		Max(50)

	validation.String(errs, "start", body.Start).
		Required().
		Time()

	validation.String(errs, "end", body.End).
		Required().
		Time()

	return errs
}
