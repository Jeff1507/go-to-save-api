package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/store/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/gin-gonic/gin"
)

type UpdateAddressRequest struct {
	Address AddressRequest `json:"address" binding:"required"`
}

type UpdateAddressHandler struct {
	updateAddress *usecase.UpdateAddress
}

func NewUpdateAddressHandler(updateAddress *usecase.UpdateAddress) *UpdateAddressHandler {
	return &UpdateAddressHandler{
		updateAddress: updateAddress,
	}
}

func (h *UpdateAddressHandler) Handle(c *gin.Context) {
	var body UpdateAddressRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.Error(err)
		return
	}

	id, ok := request.StoreID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	if err := h.updateAddress.Execute(c.Request.Context(), usecase.UpdateAddressInput{
		ID: id,
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
	}); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
