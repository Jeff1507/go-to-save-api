package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/bag"
	"github.com/Jeff1507/go-to-save-api/internal/core/product"
	"github.com/Jeff1507/go-to-save-api/internal/pkg/uuid"
)

type CreateBagPickupWindowInput struct {
	Start time.Time
	End   time.Time
}

type CreateBagItemInput struct {
	ProductID string
	Quantity  int32
}

type CreateBagInput struct {
	StoreID       string
	Title         string
	Description   string
	Stock         int32
	OriginalPrice int64
	Discount      int
	PickupWindow  CreateBagPickupWindowInput
	Items         []CreateBagItemInput
}

type CreateBagOutput struct {
	ID string
}

type CreateBag struct {
	clock       clock.Clock
	bagRepo     bag.Repository
	productRepo product.Repository
}

func NewCreateBag(
	clock clock.Clock,
	bagRepo bag.Repository,
	productRepo product.Repository,
) *CreateBag {
	return &CreateBag{
		clock:       clock,
		bagRepo:     bagRepo,
		productRepo: productRepo,
	}
}

func (uc *CreateBag) Execute(ctx context.Context, input CreateBagInput) (CreateBagOutput, error) {
	now := uc.clock.Now()

	pickupWindow, err := bag.NewPickupWindow(input.PickupWindow.Start, input.PickupWindow.End)
	if err != nil {
		if errors.Is(err, bag.ErrInvalidPickupWindowInterval) {
			return CreateBagOutput{}, fmt.Errorf("new pickup window: %w", err)
		}

		return CreateBagOutput{}, fmt.Errorf("new pickup window: %v", err)
	}

	seen := make(map[string]bool)
	productIDs := make([]string, 0, len(input.Items))
	for _, item := range input.Items {
		if !seen[item.ProductID] {
			seen[item.ProductID] = true
			productIDs = append(productIDs, item.ProductID)
		}
	}

	products, err := uc.productRepo.List(ctx, productIDs)
	if err != nil {
		return CreateBagOutput{}, fmt.Errorf("list products by id: %v", err)
	}

	if len(products) != len(productIDs) {
		return CreateBagOutput{}, product.ErrNotFound
	}

	for _, prod := range products {
		if prod.StoreID() != input.StoreID {
			return CreateBagOutput{}, product.ErrNotFound
		}
	}

	itemsMap := make(map[string]CreateBagItemInput, len(input.Items))
	for _, item := range input.Items {
		itemsMap[item.ProductID] = item
	}

	bagItems := make([]bag.Item, 0, len(products))
	for _, prod := range products {
		item := itemsMap[prod.ID()]

		bagItem, err := bag.NewItem(uuid.Generate(), prod.Name(), prod.ImageKey(), item.Quantity)
		if err != nil {
			return CreateBagOutput{}, fmt.Errorf("new item: %v", err)
		}

		bagItems = append(bagItems, bagItem)
	}

	b, err := bag.New(
		uuid.Generate(),
		input.StoreID,
		input.Title,
		input.Description,
		input.Stock,
		input.OriginalPrice,
		input.Discount,
		pickupWindow,
		bagItems,
		now,
	)
	if err != nil {
		return CreateBagOutput{}, fmt.Errorf("new bag: %v", err)
	}

	if err := uc.bagRepo.Save(ctx, b); err != nil {
		return CreateBagOutput{}, fmt.Errorf("save bag: %v", err)
	}

	return CreateBagOutput{
		ID: b.ID(),
	}, nil
}
