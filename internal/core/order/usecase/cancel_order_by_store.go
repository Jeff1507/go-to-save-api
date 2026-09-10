package usecase

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/bag"
	"github.com/Jeff1507/go-to-save-api/internal/core/order"
	"github.com/Jeff1507/go-to-save-api/internal/infra/transaction"
)

type CancelOrderByStoreInput struct {
	ID      string
	StoreID string
}

type CancelOrderByStore struct {
	transactions transaction.Manager
	orderRepo    order.Repository
	bagRepo      bag.Repository
	clock        clock.Clock
}

func NewCancelOrderByStore(
	transactions transaction.Manager,
	orderRepo order.Repository,
	bagRepo bag.Repository,
	clock clock.Clock,
) *CancelOrderByStore {
	return &CancelOrderByStore{
		transactions: transactions,
		orderRepo:    orderRepo,
		bagRepo:      bagRepo,
		clock:        clock,
	}
}

// TODO: Job para marcar pedidos expirados em 1 hora ou mais como NO_SHOW
func (uc *CancelOrderByStore) Execute(ctx context.Context, input CancelOrderByStoreInput) error {
	now := uc.clock.Now()

	return uc.transactions.Do(ctx, func(ctx context.Context) error {
		ord, err := uc.orderRepo.Find(ctx, input.ID)
		if err != nil {
			return fmt.Errorf("find order by id: %w", err)
		}

		if ord.StoreID() != input.StoreID {
			return order.ErrNotFound
		}

		if err := ord.Cancel(now); err != nil {
			return fmt.Errorf("cancel order: %w", err)
		}

		for _, item := range ord.Items() {
			if err := uc.bagRepo.FindAndUpdate(ctx, item.BagID(), func(b *bag.Bag) error {
				return b.Restock(item.Quantity(), now)
			}); err != nil {
				return fmt.Errorf("find and update bag(%s): %w", item.BagID(), err)
			}
		}

		if err := uc.orderRepo.Update(ctx, ord); err != nil {
			return fmt.Errorf("update order: %w", err)
		}

		return nil
	})
}
