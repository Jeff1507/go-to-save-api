package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/core/customer/query"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/render"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/reqctx"
)

type CustomerProfileResponse struct {
	ID         string     `json:"id"`
	Email      string     `json:"email"`
	FirstName  string     `json:"firstName"`
	LastName   string     `json:"lastName"`
	AvatarURL  *string    `json:"avatarUrl"`
	VerifiedAt *time.Time `json:"verifiedAt"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

type GetCustomer struct {
	qry *query.GetCustomer
}

func NewGetCustomer(qry *query.GetCustomer) *GetCustomer {
	return &GetCustomer{
		qry: qry,
	}
}

func (h *GetCustomer) Handle(w http.ResponseWriter, r *http.Request) error {
	customerID, ok := reqctx.CustomerID(r.Context())
	if !ok {
		return nil
	}

	output, err := h.qry.Execute(r.Context(), customerID)
	if err != nil {
		return fmt.Errorf("get customer: %w", err)
	}

	render.Success(w, http.StatusOK, CustomerProfileResponse{
		ID:         output.ID,
		Email:      output.Email,
		FirstName:  output.FirstName,
		LastName:   output.LastName,
		AvatarURL:  output.PictureURL,
		VerifiedAt: output.VerifiedAt,
		CreatedAt:  output.CreatedAt,
		UpdatedAt:  output.UpdatedAt,
	})

	return nil
}
