package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/store/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/cookie"
	"github.com/Jeff1507/go-to-save-api/internal/infra/auth/session"
	"github.com/gin-gonic/gin"
)

type SignupRequest struct {
	Email    string         `json:"email" binding:"required,email"`
	Name     string         `json:"name" binding:"required,notblank,min=2,max=100"`
	CNPJ     string         `json:"cnpj" binding:"required,numeric,len=14"`
	Password string         `json:"password" binding:"required,min=8,max=72"`
	Address  AddressRequest `json:"address" binding:"required"`
}

type SignupResponse struct {
	UserID string `json:"userId"`
}

type SignupHandler struct {
	signup *usecase.Signup
}

func NewSignupHandler(signup *usecase.Signup) *SignupHandler {
	return &SignupHandler{
		signup: signup,
	}
}

func (h *SignupHandler) Handle(c *gin.Context) {
	var body SignupRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.Error(err)
		return
	}

	output, err := h.signup.Execute(c.Request.Context(), usecase.SignupInput{
		Email:    body.Email,
		Name:     body.Name,
		CNPJ:     body.CNPJ,
		Password: body.Password,
		Address: usecase.AddressInput{
			PostalCode:   body.Address.PostalCode,
			Street:       body.Address.Street,
			Number:       body.Address.Number,
			Neighborhood: body.Address.Neighborhood,
			Complement:   body.Address.Complement,
			City:         body.Address.City,
			State:        body.Address.State,
			Latitude:     *body.Address.Latitude,
			Longitude:    *body.Address.Longitude,
		},
	})
	if err != nil {
		c.Error(err)
		return
	}

	cookie.SetSessionCookie(c, output.SessionID, session.TTL)

	c.JSON(http.StatusCreated, gin.H{
		"data": SignupResponse{
			UserID: output.UserID,
		},
	})
}
