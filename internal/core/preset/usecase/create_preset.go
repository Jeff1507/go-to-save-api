package usecase

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/common/pickup"
	"github.com/Jeff1507/go-to-save-api/internal/core/preset"
	"github.com/Jeff1507/go-to-save-api/internal/pkg/uuid"
)

type CreatePresetInput struct {
	StoreID string
	Name    string
	Start   string
	End     string
}

type CreatePresetOutput struct {
	ID string
}

type CreatePreset struct {
	presetRepo preset.Repository
}

func NewCreatePreset(presetRepo preset.Repository) *CreatePreset {
	return &CreatePreset{
		presetRepo: presetRepo,
	}
}

func (uc *CreatePreset) Execute(
	ctx context.Context,
	input CreatePresetInput,
) (CreatePresetOutput, error) {
	start, err := pickup.TimeFromString(input.Start)
	if err != nil {
		return CreatePresetOutput{}, fmt.Errorf("parse start time: %w", err)
	}

	end, err := pickup.TimeFromString(input.End)
	if err != nil {
		return CreatePresetOutput{}, fmt.Errorf("parse end time: %w", err)
	}

	pset, err := preset.New(uuid.Generate(), input.StoreID, input.Name, start, end)
	if err != nil {
		return CreatePresetOutput{}, fmt.Errorf("new preset: %w", err)
	}

	if err := uc.presetRepo.Save(ctx, pset); err != nil {
		return CreatePresetOutput{}, fmt.Errorf("save preset: %w", err)
	}

	return CreatePresetOutput{
		ID: pset.ID(),
	}, nil
}
