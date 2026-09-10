package usecase

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/bag"
)

type UpdateStockInput struct {
	ID      string
	StoreID string
	Delta   int32 // positivo = adicionar, negativo = remover
}

type UpdateStock struct {
	bagRepo bag.Repository
	clock   clock.Clock
}

func NewUpdateStock(bagRepo bag.Repository, clock clock.Clock) *UpdateStock {
	return &UpdateStock{
		bagRepo: bagRepo,
		clock:   clock,
	}
}

func (uc *UpdateStock) Execute(ctx context.Context, input UpdateStockInput) error {
	now := uc.clock.Now()

	if input.Delta == 0 {
		return bag.ErrZeroDelta
	}

	return uc.bagRepo.FindAndUpdate(ctx, input.ID, func(b *bag.Bag) error {
		if b.StoreID() != input.StoreID {
			return bag.ErrNotFound
		}

		if err := b.AdjustStock(input.Delta, now); err != nil {
			return fmt.Errorf("adjust stock: %w", err)
		}

		return nil
	})
}
