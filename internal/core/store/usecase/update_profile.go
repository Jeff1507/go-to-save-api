package usecase

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/store"
)

type UpdateProfileInput struct {
	ID          string
	Name        string
	Description string
}

type UpdateProfile struct {
	storeRepo store.Repository
	clock     clock.Clock
}

func NewUpdateProfile(storeRepo store.Repository, clock clock.Clock) *UpdateProfile {
	return &UpdateProfile{
		storeRepo: storeRepo,
		clock:     clock,
	}
}

func (uc *UpdateProfile) Execute(ctx context.Context, input UpdateProfileInput) error {
	now := uc.clock.Now()

	s, err := uc.storeRepo.Find(ctx, input.ID)
	if err != nil {
		return fmt.Errorf("find store by id: %w", err)
	}

	if err := s.ChangeName(input.Name, now); err != nil {
		return fmt.Errorf("change name: %w", err)
	}

	if err := s.ChangeDescription(input.Description, now); err != nil {
		return fmt.Errorf("change description: %w", err)
	}

	if err := uc.storeRepo.Update(ctx, s); err != nil {
		return fmt.Errorf("update store: %w", err)
	}

	return nil
}
