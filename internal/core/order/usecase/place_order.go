package usecase

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/bag"
	"github.com/Jeff1507/go-to-save-api/internal/core/order"
	"github.com/Jeff1507/go-to-save-api/internal/core/pickup"
	"github.com/Jeff1507/go-to-save-api/internal/infra/transaction"
	"github.com/Jeff1507/go-to-save-api/internal/pkg/uuid"
)

type ItemInput struct {
	BagID    string
	Quantity int32
}

type PlaceOrderInput struct {
	CustomerID string
	Items      []ItemInput
}

type PlaceOrderOutput struct {
	ID string
}

type PlaceOrder struct {
	transactions transaction.Manager
	orderRepo    order.Repository
	bagRepo      bag.Repository
	clock        clock.Clock
}

func NewPlaceOrder(
	transactions transaction.Manager,
	orderRepo order.Repository,
	bagRepo bag.Repository,
	clock clock.Clock,
) *PlaceOrder {
	return &PlaceOrder{
		transactions: transactions,
		orderRepo:    orderRepo,
		bagRepo:      bagRepo,
		clock:        clock,
	}
}

// TODO(prioridade baixa): Revisar retorno dos erros sentinelas
func (uc *PlaceOrder) Execute(
	ctx context.Context,
	input PlaceOrderInput,
) (PlaceOrderOutput, error) {
	now := uc.clock.Now()
	var output PlaceOrderOutput

	if len(input.Items) == 0 {
		return output, order.ErrCannotPlaceOrder
	}

	err := uc.transactions.Do(ctx, func(ctx context.Context) error {
		hasPending, err := uc.orderRepo.HasPendingByCustomer(ctx, input.CustomerID)
		if err != nil {
			return fmt.Errorf("check pending order: %w", err)
		}

		if hasPending {
			return order.ErrCustomerHasPendingOrder
		}

		items := make([]ItemInput, len(input.Items))
		copy(items, input.Items)

		seen := make(map[string]bool, len(items))
		for _, item := range items {
			if seen[item.BagID] {
				return order.ErrCannotPlaceOrder
			}

			seen[item.BagID] = true
		}

		sort.Slice(items, func(i, j int) bool {
			return items[i].BagID < items[j].BagID
		})

		var (
			referenceStoreID string
			referenceWindow  pickup.Window
			orderItems       = make([]order.Item, 0, len(items))
		)

		for i, item := range items {
			var unitPrice int64

			if err := uc.bagRepo.FindAndUpdate(ctx, item.BagID, func(b *bag.Bag) error {
				if i == 0 {
					referenceStoreID = b.StoreID()
					referenceWindow = b.PickupWindow()
				} else {
					if b.StoreID() != referenceStoreID {
						return order.ErrCannotPlaceOrder
					}

					if !b.PickupWindow().Equal(referenceWindow) {
						return order.ErrCannotPlaceOrder
					}
				}

				if err := b.Reserve(item.Quantity, now); err != nil {
					return fmt.Errorf("reserve bag: %w: %v", order.ErrCannotPlaceOrder, err)
				}

				unitPrice = b.FinalPrice()

				return nil
			}); err != nil {
				if errors.Is(err, bag.ErrNotFound) {
					return order.ErrItemNotFound
				}

				return fmt.Errorf("find and update bag(%s): %w", item.BagID, err)
			}

			orderItem, err := order.NewItem(item.BagID, unitPrice, item.Quantity)
			if err != nil {
				return fmt.Errorf("new item: %w", err)
			}

			orderItems = append(orderItems, orderItem)
		}

		// Adiciona +30 minutos caso a janela de retirada se expire em 15 minutos ou menos
		if referenceWindow.IsNearExpiry(now, time.Minute*15) {
			extendedWindow, err := referenceWindow.Extended(time.Minute * 30)
			if err != nil {
				return fmt.Errorf("extend pickup window: %w", err)
			}

			referenceWindow = extendedWindow
		}

		ord, err := order.New(
			uuid.Generate(),
			input.CustomerID,
			referenceStoreID,
			referenceWindow,
			orderItems,
			now,
		)
		if err != nil {
			return fmt.Errorf("new order: %w", err)
		}

		if err := uc.orderRepo.Save(ctx, ord); err != nil {
			return fmt.Errorf("save order: %w", err)
		}

		output = PlaceOrderOutput{
			ID: ord.ID(),
		}

		return nil
	})

	return output, err
}
