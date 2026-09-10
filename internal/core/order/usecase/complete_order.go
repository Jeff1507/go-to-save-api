package usecase

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/order"
)

type CompleteOrderInput struct {
	ID               string
	StoreID          string
	ConfirmationCode string
}

type CompleteOrder struct {
	orderRepo order.Repository
	clock     clock.Clock
}

func NewCompleteOrder(orderRepo order.Repository, clock clock.Clock) *CompleteOrder {
	return &CompleteOrder{
		orderRepo: orderRepo,
		clock:     clock,
	}
}

// TODO: Job para marcar pedidos expirados em 1 hora ou mais como NO_SHOW
func (uc *CompleteOrder) Execute(ctx context.Context, input CompleteOrderInput) error {
	now := uc.clock.Now()

	ord, err := uc.orderRepo.Find(ctx, input.ID)
	if err != nil {
		return fmt.Errorf("find order by id: %w", err)
	}

	if ord.StoreID() != input.StoreID {
		return order.ErrNotFound
	}

	if err := ord.Complete(input.ConfirmationCode, now); err != nil {
		return fmt.Errorf("complete order: %w", err)
	}

	if err := uc.orderRepo.Update(ctx, ord); err != nil {
		return fmt.Errorf("update order: %w", err)
	}

	return nil
}
