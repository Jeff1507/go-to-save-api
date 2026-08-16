package usecase

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/core/pickup"
)

type DeleteTemplateInput struct {
	ID      string
	StoreID string
}

type DeleteTemplate struct {
	pickupRepo pickup.Repository
}

func NewDeleteTemplate(pickupRepo pickup.Repository) *DeleteTemplate {
	return &DeleteTemplate{
		pickupRepo: pickupRepo,
	}
}

func (uc *DeleteTemplate) Execute(ctx context.Context, input DeleteTemplateInput) error {
	templ, err := uc.pickupRepo.Find(ctx, input.ID)
	if err != nil {
		return fmt.Errorf("find pickup template by id: %w", err)
	}

	if templ.StoreID() != input.StoreID {
		return pickup.ErrTemplateNotFound
	}

	if err := uc.pickupRepo.Delete(ctx, input.ID); err != nil {
		return fmt.Errorf("delete pickup template: %w", err)
	}

	return nil
}
