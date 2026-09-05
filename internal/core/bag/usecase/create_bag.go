package usecase

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/bag"
	"github.com/Jeff1507/go-to-save-api/internal/core/pickup"
	"github.com/Jeff1507/go-to-save-api/internal/core/product"
	"github.com/Jeff1507/go-to-save-api/internal/pkg/uuid"
)

type CreateBagInput struct {
	StoreID       string
	Title         string
	Description   string
	Stock         int32
	OriginalPrice int64
	Discount      int
	PickupStart   string
	PickupEnd     string
	Items         []ItemInput
}

type CreateBagOutput struct {
	ID string
}

type CreateBag struct {
	bagRepo     bag.Repository
	productRepo product.Repository
	clock       clock.Clock
}

func NewCreateBag(
	bagRepo bag.Repository,
	productRepo product.Repository,
	clock clock.Clock,
) *CreateBag {
	return &CreateBag{
		bagRepo:     bagRepo,
		productRepo: productRepo,
		clock:       clock,
	}
}

func (uc *CreateBag) Execute(ctx context.Context, input CreateBagInput) (CreateBagOutput, error) {
	now := uc.clock.Now()

	start, err := pickup.TimeFromString(input.PickupStart)
	if err != nil {
		return CreateBagOutput{}, fmt.Errorf("parse start time: %w", err)
	}

	end, err := pickup.TimeFromString(input.PickupEnd)
	if err != nil {
		return CreateBagOutput{}, fmt.Errorf("parse end time: %w", err)
	}

	pickupWindow, err := pickup.ResolveWindow(start, end, now)
	if err != nil {
		return CreateBagOutput{}, fmt.Errorf("resolve pickup window: %w", err)
	}

	resolveBagItems := NewResolveBagItems(uc.productRepo)
	bagItems, err := resolveBagItems.Execute(ctx, input.StoreID, input.Items)
	if err != nil {
		return CreateBagOutput{}, fmt.Errorf("resolve bag items: %w", err)
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
		return CreateBagOutput{}, fmt.Errorf("new bag: %w", err)
	}

	if err := uc.bagRepo.Save(ctx, b); err != nil {
		return CreateBagOutput{}, fmt.Errorf("save bag: %w", err)
	}

	return CreateBagOutput{
		ID: b.ID(),
	}, nil
}
