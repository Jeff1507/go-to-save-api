package usecase

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/core/pickup"
)

type UpdateTemplateTimeInput struct {
	Hour   int
	Minute int
}

type UpdateTemplateInput struct {
	ID      string
	StoreID string
	Name    string
	Start   UpdateTemplateTimeInput
	End     UpdateTemplateTimeInput
}

type UpdateTemplate struct {
	pickupRepo pickup.Repository
}

func NewUpdateTemplate(pickupRepo pickup.Repository) *UpdateTemplate {
	return &UpdateTemplate{
		pickupRepo: pickupRepo,
	}
}

func (uc *UpdateTemplate) Execute(ctx context.Context, input UpdateTemplateInput) error {
	templ, err := uc.pickupRepo.Find(ctx, input.ID)
	if err != nil {
		return fmt.Errorf("find pickup template by id: %w", err)
	}

	if templ.StoreID() != input.StoreID {
		return pickup.ErrTemplateNotFound
	}

	if err := templ.ChangeName(input.Name); err != nil {
		return fmt.Errorf("change name: %w", err)
	}

	start, err := pickup.NewTime(input.Start.Hour, input.Start.Minute)
	if err != nil {
		return fmt.Errorf("new start time: %w", err)
	}

	end, err := pickup.NewTime(input.End.Hour, input.End.Minute)
	if err != nil {
		return fmt.Errorf("new end time: %w", err)
	}

	if err := templ.ChangeTime(start, end); err != nil {
		return fmt.Errorf("change time: %w", err)
	}

	if err := uc.pickupRepo.Update(ctx, templ); err != nil {
		return fmt.Errorf("update pickup template: %w", err)
	}

	return nil
}
