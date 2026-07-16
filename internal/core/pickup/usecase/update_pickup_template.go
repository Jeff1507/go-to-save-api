package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/core/pickup"
)

type UpdatePickupTemplateTimeInput struct {
	Hour   int
	Minute int
}

type UpdatePickupTemplateInput struct {
	ID    string
	Name  string
	Start UpdatePickupTemplateTimeInput
	End   UpdatePickupTemplateTimeInput
}

type UpdatePickupTemplate struct {
	pickupRepo pickup.Repository
}

func NewUpdatePickupTemplate(pickupRepo pickup.Repository) *UpdatePickupTemplate {
	return &UpdatePickupTemplate{
		pickupRepo: pickupRepo,
	}
}

func (uc *UpdatePickupTemplate) Execute(
	ctx context.Context,
	input UpdatePickupTemplateInput,
) error {
	templ, err := uc.pickupRepo.Find(ctx, input.ID)
	if err != nil {
		if errors.Is(err, pickup.ErrTemplateNotFound) {
			return fmt.Errorf("find pickup template by id: %w", err)
		}

		return fmt.Errorf("find pickup template by id: %v", err)
	}

	if err := templ.ChangeName(input.Name); err != nil {
		return fmt.Errorf("change pickup template name: %v", err)
	}

	start, err := pickup.NewTime(input.Start.Hour, input.Start.Minute)
	if err != nil {
		return fmt.Errorf("new start time: %v", err)
	}

	end, err := pickup.NewTime(input.End.Hour, input.End.Minute)
	if err != nil {
		return fmt.Errorf("new end time: %v", err)
	}

	if err := templ.ChangeTime(start, end); err != nil {
		if errors.Is(err, pickup.ErrInvalidInterval) {
			return fmt.Errorf("change pickup template time: %w", err)
		}

		return fmt.Errorf("change pickup template time: %v", err)
	}

	if err := uc.pickupRepo.Update(ctx, templ); err != nil {
		if errors.Is(err, pickup.ErrTemplateAlreadyExists) {
			return fmt.Errorf("update pickup template: %w", err)
		}

		return fmt.Errorf("update pickup template: %v", err)
	}

	return nil
}
