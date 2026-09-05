package usecase

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/core/bag"
	"github.com/Jeff1507/go-to-save-api/internal/core/product"
	"github.com/Jeff1507/go-to-save-api/internal/pkg/uuid"
)

type ResolveBagItems struct {
	productRepo product.Repository
}

func NewResolveBagItems(productRepo product.Repository) *ResolveBagItems {
	return &ResolveBagItems{
		productRepo: productRepo,
	}
}

func (uc *ResolveBagItems) Execute(
	ctx context.Context,
	storeID string,
	items []ItemInput,
) ([]bag.Item, error) {
	itemsMap := make(map[string]ItemInput, len(items))
	productIDs := make([]string, 0, len(items))
	for _, item := range items {
		if _, exists := itemsMap[item.ProductID]; exists {
			return nil, bag.ErrDuplicateItems
		}

		itemsMap[item.ProductID] = item
		productIDs = append(productIDs, item.ProductID)
	}

	products, err := uc.productRepo.List(ctx, productIDs)
	if err != nil {
		return nil, fmt.Errorf("list products by id: %w", err)
	}

	if len(products) != len(productIDs) {
		return nil, bag.ErrProductNotFound
	}

	bagItems := make([]bag.Item, 0, len(products))
	for _, prod := range products {
		if prod.StoreID() != storeID {
			return nil, bag.ErrProductNotFound
		}

		item := itemsMap[prod.ID()]

		bagItem, err := bag.NewItem(uuid.Generate(), prod.Name(), prod.ImageKey(), item.Quantity)
		if err != nil {
			return nil, fmt.Errorf("new item: %w", err)
		}

		bagItems = append(bagItems, bagItem)
	}

	return bagItems, nil
}
