package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/core/pickup"
	"github.com/Jeff1507/go-to-save-api/internal/pkg/uuid"
)

type CreatePickupTemplateTimeInput struct {
	Hour   int
	Minute int
}

type CreatePickupTemplateInput struct {
	StoreID string
	Name    string
	Start   CreatePickupTemplateTimeInput
	End     CreatePickupTemplateTimeInput
}

type CreatePickupTemplateOutput struct {
	ID string
}

type CreatePickupTemplate struct {
	pickupRepo pickup.Repository
}

func NewCreatePickupTemplate(pickupRepo pickup.Repository) *CreatePickupTemplate {
	return &CreatePickupTemplate{
		pickupRepo: pickupRepo,
	}
}

func (uc *CreatePickupTemplate) Execute(
	ctx context.Context,
	input CreatePickupTemplateInput,
) (CreatePickupTemplateOutput, error) {
	start, err := pickup.NewTime(input.Start.Hour, input.Start.Minute)
	if err != nil {
		return CreatePickupTemplateOutput{}, fmt.Errorf("new start time: %v", err)
	}

	end, err := pickup.NewTime(input.End.Hour, input.End.Minute)
	if err != nil {
		return CreatePickupTemplateOutput{}, fmt.Errorf("new end time: %v", err)
	}

	templ, err := pickup.NewTemplate(uuid.Generate(), input.StoreID, input.Name, start, end)
	if err != nil {
		if errors.Is(err, pickup.ErrInvalidInterval) {
			return CreatePickupTemplateOutput{}, fmt.Errorf("new pickup template: %w", err)
		}

		return CreatePickupTemplateOutput{}, fmt.Errorf("new pickup template: %v", err)
	}

	if err := uc.pickupRepo.Save(ctx, templ); err != nil {
		if errors.Is(err, pickup.ErrTemplateAlreadyExists) {
			return CreatePickupTemplateOutput{}, fmt.Errorf("save pickup template: %w", err)
		}

		return CreatePickupTemplateOutput{}, fmt.Errorf("save pickup template: %v", err)
	}

	return CreatePickupTemplateOutput{
		ID: templ.ID(),
	}, nil
}
