package usecase

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/store"
)

type UpdateAddressInput struct {
	ID      string
	Address AddressInput
}

type UpdateAddress struct {
	storeRepo store.Repository
	clock     clock.Clock
}

func NewUpdateAddress(storeRepo store.Repository, clock clock.Clock) *UpdateAddress {
	return &UpdateAddress{
		storeRepo: storeRepo,
		clock:     clock,
	}
}

func (uc *UpdateAddress) Execute(ctx context.Context, input UpdateAddressInput) error {
	now := uc.clock.Now()

	s, err := uc.storeRepo.Find(ctx, input.ID)
	if err != nil {
		return fmt.Errorf("find store by id: %w", err)
	}

	address, err := store.NewAddress(
		input.Address.PostalCode,
		input.Address.Street,
		input.Address.Number,
		input.Address.Neighborhood,
		input.Address.Complement,
		input.Address.City,
		input.Address.State,
		input.Address.Latitude,
		input.Address.Longitude,
	)
	if err != nil {
		return fmt.Errorf("new address: %w", err)
	}

	if err := s.ChangeAddress(address, now); err != nil {
		return fmt.Errorf("change address: %w", err)
	}

	if err := uc.storeRepo.Update(ctx, s); err != nil {
		return fmt.Errorf("update store: %w", err)
	}

	return nil
}
