package usecase

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/core/preset"
)

type DeletePresetInput struct {
	ID      string
	StoreID string
}

type DeletePreset struct {
	presetRepo preset.Repository
}

func NewDeletePreset(presetRepo preset.Repository) *DeletePreset {
	return &DeletePreset{
		presetRepo: presetRepo,
	}
}

func (uc *DeletePreset) Execute(ctx context.Context, input DeletePresetInput) error {
	pset, err := uc.presetRepo.Find(ctx, input.ID)
	if err != nil {
		return fmt.Errorf("find preset by id: %w", err)
	}

	if pset.StoreID() != input.StoreID {
		return preset.ErrNotFound
	}

	if err := uc.presetRepo.Delete(ctx, input.ID); err != nil {
		return fmt.Errorf("delete preset: %w", err)
	}

	return nil
}
