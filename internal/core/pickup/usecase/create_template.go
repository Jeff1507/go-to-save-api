package usecase

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/core/pickup"
	"github.com/Jeff1507/go-to-save-api/internal/pkg/uuid"
)

type CreateTemplateInput struct {
	StoreID string
	Name    string
	Start   TimeInput
	End     TimeInput
}

type CreateTemplateOutput struct {
	ID string
}

type CreateTemplate struct {
	pickupRepo pickup.Repository
}

func NewCreateTemplate(pickupRepo pickup.Repository) *CreateTemplate {
	return &CreateTemplate{
		pickupRepo: pickupRepo,
	}
}

func (uc *CreateTemplate) Execute(
	ctx context.Context,
	input CreateTemplateInput,
) (CreateTemplateOutput, error) {
	start, err := pickup.NewTime(input.Start.Hour, input.Start.Minute)
	if err != nil {
		return CreateTemplateOutput{}, fmt.Errorf("new start time: %w", err)
	}

	end, err := pickup.NewTime(input.End.Hour, input.End.Minute)
	if err != nil {
		return CreateTemplateOutput{}, fmt.Errorf("new end time: %w", err)
	}

	templ, err := pickup.NewTemplate(uuid.Generate(), input.StoreID, input.Name, start, end)
	if err != nil {
		return CreateTemplateOutput{}, fmt.Errorf("new pickup template: %w", err)
	}

	if err := uc.pickupRepo.Save(ctx, templ); err != nil {
		return CreateTemplateOutput{}, fmt.Errorf("save pickup template: %w", err)
	}

	return CreateTemplateOutput{
		ID: templ.ID(),
	}, nil
}
