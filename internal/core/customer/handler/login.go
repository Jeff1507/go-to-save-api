package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/customer/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/render"
	"github.com/Jeff1507/go-to-save-api/internal/infra/validation"
	"github.com/go-playground/validator/v10"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type LoginResponse struct {
	SessionID string `json:"sessionId"`
}

type Login struct {
	v  *validator.Validate
	uc *usecase.Login
}

func NewLogin(v *validator.Validate, uc *usecase.Login) *Login {
	return &Login{
		v:  v,
		uc: uc,
	}
}

func (h *Login) Handle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var body LoginRequest
	if err := render.Bind(r, &body); err != nil {
		render.InvalidRequestBody(w)
		return
	}

	if errs := h.validateRequest(body); errs.HasErrors() {
		render.ValidationError(w, errs)
		return
	}

	output, err := h.uc.Execute(r.Context(), usecase.LoginInput{
		Email:    body.Email,
		Password: body.Password,
	})
	if err != nil {
		handleError(w, err)
		return
	}

	render.Success(w, http.StatusOK, LoginResponse{
		SessionID: output.SessionID,
	})
}

func (h *Login) validateRequest(body LoginRequest) validation.Errors {
	errs := validation.NewErrors()

	validation.Struct(h.v, errs, body)

	return errs
}
