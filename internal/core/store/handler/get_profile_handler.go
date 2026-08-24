package handler

import (
	"net/http"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/core/store/query"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/gin-gonic/gin"
)

type ProfileResponse struct {
	ID          string          `json:"id"`
	Email       string          `json:"email"`
	Name        string          `json:"name"`
	CNPJ        string          `json:"cnpj"`
	AvatarURL   *string         `json:"avatarUrl"`
	Description *string         `json:"description"`
	VerifiedAt  *time.Time      `json:"verifiedAt"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
	Address     AddressResponse `json:"address"`
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
	id, ok := request.StoreID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	profile, err := h.getProfile.Execute(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": ProfileResponse{
			ID:          profile.ID,
			Email:       profile.Email,
			Name:        profile.Name,
			CNPJ:        profile.CNPJ,
			AvatarURL:   profile.AvatarURL,
			Description: profile.Description,
			VerifiedAt:  profile.VerifiedAt,
			CreatedAt:   profile.CreatedAt,
			UpdatedAt:   profile.UpdatedAt,
			Address:     AddressResponse(profile.Address),
		},
	})
}
