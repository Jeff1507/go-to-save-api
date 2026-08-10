package handler

import (
	"net/http"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/core/customer/query"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/gin-gonic/gin"
)

type ProfileResponse struct {
	ID         string     `json:"id"`
	Email      string     `json:"email"`
	FirstName  string     `json:"firstName"`
	LastName   string     `json:"lastName"`
	AvatarURL  *string    `json:"avatarUrl"`
	VerifiedAt *time.Time `json:"verifiedAt"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

type GetProfileHandler struct {
	getProfile *query.GetProfile
}

func NewGetProfileHandler(getProfile *query.GetProfile) *GetProfileHandler {
	return &GetProfileHandler{
		getProfile: getProfile,
	}
}

func (h *GetProfileHandler) Handle(c *gin.Context) {
	id, ok := request.CustomerID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	output, err := h.getProfile.Execute(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": ProfileResponse{
			ID:         output.ID,
			Email:      output.Email,
			FirstName:  output.FirstName,
			LastName:   output.LastName,
			AvatarURL:  output.AvatarURL,
			VerifiedAt: output.VerifiedAt,
			CreatedAt:  output.CreatedAt,
			UpdatedAt:  output.UpdatedAt,
		},
	})
}
