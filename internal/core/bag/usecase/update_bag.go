package usecase

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/bag"
	"github.com/Jeff1507/go-to-save-api/internal/core/order"
	"github.com/Jeff1507/go-to-save-api/internal/core/pickup"
	"github.com/Jeff1507/go-to-save-api/internal/core/product"
)

type UpdateBagInput struct {
	ID            string
	StoreID       string
	Title         string
	Description   string
	OriginalPrice int64
	Discount      int
	PickupStart   string
	PickupEnd     string
	Items         []ItemInput
}

type UpdateBag struct {
	bagRepo     bag.Repository
	productRepo product.Repository
	orderRepo   order.Repository
	clock       clock.Clock
}

func NewUpdateBag(
	bagRepo bag.Repository,
	productRepo product.Repository,
	orderRepo order.Repository,
	clock clock.Clock,
) *UpdateBag {
	return &UpdateBag{
		bagRepo:     bagRepo,
		productRepo: productRepo,
		orderRepo:   orderRepo,
		clock:       clock,
	}
}

func (uc *UpdateBag) Execute(ctx context.Context, input UpdateBagInput) error {
	now := uc.clock.Now()

	b, err := uc.bagRepo.Find(ctx, input.ID)
	if err != nil {
		return fmt.Errorf("find bag by id: %w", err)
	}

	if b.StoreID() != input.StoreID {
		return bag.ErrNotFound
	}

	if b.IsClosed() {
		return bag.ErrClosed
	}

	if b.IsExpired(now) {
		return bag.ErrExpired
	}

	// Verifica se a bag pertence a um pedido, se pertencer a edição é bloqueada.
	hasPendingOrder, err := uc.orderRepo.HasPendingByBag(ctx, b.ID())
	if err != nil {
		return fmt.Errorf("check pending order: %w", err)
	}

	if hasPendingOrder {
		return bag.ErrHasPendingOrder
	}

	if err := b.ChangeTitle(input.Title, now); err != nil {
		return fmt.Errorf("change title: %w", err)
	}

	if err := b.ChangeDescription(input.Description, now); err != nil {
		return fmt.Errorf("change description: %w", err)
	}

	if err := b.ChangePrice(input.OriginalPrice, input.Discount, now); err != nil {
		return fmt.Errorf("change price: %w", err)
	}

	start, err := pickup.TimeFromString(input.PickupStart)
	if err != nil {
		return fmt.Errorf("parse start time: %w", err)
	}

	end, err := pickup.TimeFromString(input.PickupEnd)
	if err != nil {
		return fmt.Errorf("parse end time: %w", err)
	}

	pickupWindow, err := pickup.ResolveWindow(start, end, now)
	if err != nil {
		return fmt.Errorf("new pickup window: %w", err)
	}

	if err := b.ChangePickupWindow(pickupWindow, now); err != nil {
		return fmt.Errorf("change pickup window: %w", err)
	}

	resolveBagItems := NewResolveBagItems(uc.productRepo)
	bagItems, err := resolveBagItems.Execute(ctx, input.StoreID, input.Items)
	if err != nil {
		return fmt.Errorf("resolve bag items: %w", err)
	}

	if err := b.ReplaceItems(bagItems, now); err != nil {
		return fmt.Errorf("replace items: %w", err)
	}

	if err := uc.bagRepo.Update(ctx, b); err != nil {
		return fmt.Errorf("update bag: %w", err)
	}

	return nil
}
