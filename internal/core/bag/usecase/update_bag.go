package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/bag"
	"github.com/Jeff1507/go-to-save-api/internal/core/product"
)

type UpdateBagPickupWindowInput struct {
	Start time.Time
	End   time.Time
}

type UpdateBagItemInput struct {
	ID       string
	Quantity int32
}

type UpdateBagInput struct {
	ID            string
	Title         string
	Description   string
	Stock         int32
	OriginalPrice int64
	Discount      int
	PickupWindow  UpdateBagPickupWindowInput
	Items         []UpdateBagItemInput
}

type UpdateBag struct {
	clock       clock.Clock
	bagRepo     bag.Repository
	productRepo product.Repository
}

func NewUpdateBag(
	clock clock.Clock,
	bagRepo bag.Repository,
	productRepo product.Repository,
) *UpdateBag {
	return &UpdateBag{
		clock:       clock,
		bagRepo:     bagRepo,
		productRepo: productRepo,
	}
}

func (uc *UpdateBag) Execute(ctx context.Context, input UpdateBagInput) error {
	now := uc.clock.Now()

	b, err := uc.bagRepo.Find(ctx, input.ID)
	if err != nil {
		if errors.Is(err, bag.ErrNotFound) {
			return fmt.Errorf("find bag by id: %w", err)
		}

		return fmt.Errorf("find bag by id: %v", err)
	}

	if b.IsClosed() {
		return bag.ErrClosed
	}

	if b.IsExpired(now) {
		return bag.ErrExpired
	}

	// TODO: Verificar se a bag não pertence a nenhum pedido.
	// Se pertencer a um pedido, a edição é bloqueada.

	if err := b.ChangeTitle(input.Title, now); err != nil {
		return fmt.Errorf("change bag title: %v", err)
	}

	if err := b.ChangeDescription(input.Description, now); err != nil {
		return fmt.Errorf("change bag description: %v", err)
	}

	if err := b.ChangeStock(input.Stock, now); err != nil {
		return fmt.Errorf("change bag stock: %v", err)
	}

	if err := b.ChangePrice(input.OriginalPrice, input.Discount, now); err != nil {
		return fmt.Errorf("change bag price: %v", err)
	}

	pickupWin, err := bag.NewPickupWindow(input.PickupWindow.Start, input.PickupWindow.End)
	if err != nil {
		if errors.Is(err, bag.ErrInvalidPickupWindowInterval) {
			return fmt.Errorf("new pickup window: %w", err)
		}

		return fmt.Errorf("new pickup window: %v", err)
	}

	if err := b.ChangePickupWindow(pickupWin, now); err != nil {
		return fmt.Errorf("change bag pickup window: %v", err)
	}

	productIDs := make([]string, 0, len(input.Items))
	for _, item := range input.Items {
		productIDs = append(productIDs, item.ID)
	}

	products, err := uc.productRepo.List(ctx, productIDs)
	if err != nil {
		return fmt.Errorf("list products by id: %v", err)
	}

	if len(products) != len(productIDs) {
		return product.ErrNotFound
	}

	for _, prod := range products {
		if prod.StoreID() != b.StoreID() {
			return product.ErrNotFound
		}
	}

	itemsByID := make(map[string]UpdateBagItemInput, len(input.Items))
	for _, item := range input.Items {
		itemsByID[item.ID] = item
	}

	bagItems := make([]bag.Item, 0, len(products))
	for _, prod := range products {
		item := itemsByID[prod.ID()]

		bagItem, err := bag.NewItem(prod.ID(), item.Quantity)
		if err != nil {
			return fmt.Errorf("new item: %v", err)
		}

		bagItems = append(bagItems, bagItem)
	}

	if err := b.ReplaceItems(bagItems, now); err != nil {
		return fmt.Errorf("replace bag items: %v", err)
	}

	if err := uc.bagRepo.Update(ctx, b); err != nil {
		return fmt.Errorf("update bag: %v", err)
	}

	return nil
}
