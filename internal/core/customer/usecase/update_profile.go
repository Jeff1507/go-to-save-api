package usecase

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/customer"
)

type UpdateProfileInput struct {
	ID        string
	FirstName string
	LastName  string
}

type UpdateProfile struct {
	customerRepo customer.Repository
	clock        clock.Clock
}

func NewUpdateProfile(customerRepo customer.Repository, clock clock.Clock) *UpdateProfile {
	return &UpdateProfile{
		customerRepo: customerRepo,
		clock:        clock,
	}
}

func (uc *UpdateProfile) Execute(ctx context.Context, input UpdateProfileInput) error {
	now := uc.clock.Now()

	cust, err := uc.customerRepo.Find(ctx, input.ID)
	if err != nil {
		return fmt.Errorf("find customer by id: %w", err)
	}

	if err := cust.ChangeName(input.FirstName, input.LastName, now); err != nil {
		return fmt.Errorf("change first and last name: %w", err)
	}

	if err := uc.customerRepo.Update(ctx, cust); err != nil {
		return fmt.Errorf("update customer: %w", err)
	}

	return nil
}
