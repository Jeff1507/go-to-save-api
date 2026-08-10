package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/customer/usecase"
	"github.com/gin-gonic/gin"
)

type SignupRequest struct {
	FirstName string `json:"firstName" binding:"required,notblank,min=2,max=50"`
	LastName  string `json:"lastName" binding:"required,notblank,min=2,max=50"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8,max=72"`
}

type SignupResponse struct {
	UserID    string `json:"userId"`
	SessionID string `json:"sessionId"`
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
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
		Password:  body.Password,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": SignupResponse{
			UserID:    output.UserID,
			SessionID: output.SessionID,
		},
	})
}
