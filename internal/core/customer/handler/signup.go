package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/customer/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/render"
	"github.com/Jeff1507/go-to-save-api/internal/infra/validation"
	"github.com/go-playground/validator/v10"
)

type SignupRequest struct {
	FirstName string `json:"firstName" validate:"required,notblank,min=2,max=50"`
	LastName  string `json:"lastName" validate:"required,notblank,min=2,max=50"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,max=72"`
}

type SignupResponse struct {
	UserID    string `json:"userId"`
	SessionID string `json:"sessionId"`
}

type Signup struct {
	v  *validator.Validate
	uc *usecase.Signup
}

func NewSignup(v *validator.Validate, uc *usecase.Signup) *Signup {
	return &Signup{
		v:  v,
		uc: uc,
	}
}

func (h *Signup) Handle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var body SignupRequest
	if err := render.Bind(r, &body); err != nil {
		render.InvalidRequestBody(w)
		return
	}

	if errs := h.validateRequest(body); errs.HasErrors() {
		render.ValidationError(w, errs)
		return
	}

	output, err := h.uc.Execute(r.Context(), usecase.SignupInput{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
		Password:  body.Password,
	})
	if err != nil {
		handleError(w, err)
		return
	}

	render.Success(w, http.StatusCreated, SignupResponse{
		UserID:    output.UserID,
		SessionID: output.SessionID,
	})
}

func (h *Signup) validateRequest(body SignupRequest) validation.Errors {
	errs := validation.NewErrors()

	validation.Struct(h.v, errs, body)

	return errs
}
