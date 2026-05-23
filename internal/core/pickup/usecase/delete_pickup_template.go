package usecase

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/core/pickup"
)

type DeletePickupTemplate struct {
	pickupRepo pickup.Repository
}

func NewDeletePickupTemplate(pickupRepo pickup.Repository) *DeletePickupTemplate {
	return &DeletePickupTemplate{
		pickupRepo: pickupRepo,
	}
}

func (uc *DeletePickupTemplate) Execute(ctx context.Context, id string) error {
	if err := uc.pickupRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete pickup template: %v", err)
	}

	return nil
}
