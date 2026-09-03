package usecase

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/core/pickup"
	"github.com/Jeff1507/go-to-save-api/internal/core/preset"
)

type UpdatePresetInput struct {
	ID      string
	StoreID string
	Name    string
	Start   string
	End     string
}

type UpdatePreset struct {
	presetRepo preset.Repository
}

func NewUpdatePreset(presetRepo preset.Repository) *UpdatePreset {
	return &UpdatePreset{
		presetRepo: presetRepo,
	}
}

func (uc *UpdatePreset) Execute(ctx context.Context, input UpdatePresetInput) error {
	pset, err := uc.presetRepo.Find(ctx, input.ID)
	if err != nil {
		return fmt.Errorf("find preset by id: %w", err)
	}

	if pset.StoreID() != input.StoreID {
		return preset.ErrNotFound
	}

	if err := pset.ChangeName(input.Name); err != nil {
		return fmt.Errorf("change name: %w", err)
	}

	start, err := pickup.TimeFromString(input.Start)
	if err != nil {
		return fmt.Errorf("parse start time: %w", err)
	}

	end, err := pickup.TimeFromString(input.End)
	if err != nil {
		return fmt.Errorf("parse end time: %w", err)
	}

	if err := pset.ChangePickupTime(start, end); err != nil {
		return fmt.Errorf("change pickup time: %w", err)
	}

	if err := uc.presetRepo.Update(ctx, pset); err != nil {
		return fmt.Errorf("update preset: %w", err)
	}

	return nil
}
