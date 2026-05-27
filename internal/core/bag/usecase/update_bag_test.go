package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/bag"
	"github.com/Jeff1507/go-to-save-api/internal/core/bag/repository"
	"github.com/Jeff1507/go-to-save-api/internal/core/bag/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/core/product"
	productRepository "github.com/Jeff1507/go-to-save-api/internal/core/product/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type updateBagDeps struct {
	clk         clock.Clock
	bagRepo     bag.Repository
	productRepo product.Repository
	uc          *usecase.UpdateBag
}

func newUpdateBagDeps() *updateBagDeps {
	clk := clock.NewFakeClock()
	bagRepo := repository.NewInMemRepo()
	productRepo := productRepository.NewInMemRepo()
	uc := usecase.NewUpdateBag(clk, bagRepo, productRepo)

	return &updateBagDeps{
		clk:         clk,
		bagRepo:     bagRepo,
		productRepo: productRepo,
		uc:          uc,
	}
}

func createBagToUpdate(
	t *testing.T,
	productRepo product.Repository,
	bagRepo bag.Repository,
	pickupStart time.Time,
	pickupEnd time.Time,
) *bag.Bag {
	t.Helper()

	date := time.Date(2024, time.April, 1, 0, 0, 0, 0, time.UTC)

	prod, err := product.New("prod-id", "store-id", "Coxinha", "coxinha.png", date)
	require.NoError(t, err)

	err = productRepo.Save(context.Background(), prod)
	require.NoError(t, err)

	win, err := bag.NewPickupWindow(pickupStart, pickupEnd)
	require.NoError(t, err)

	item, err := bag.NewItem(prod.ID(), 2)
	require.NoError(t, err)

	b, err := bag.New("id", "store-id", "XDD", "", 1, 500, 50, win, []bag.Item{item}, pickupStart)
	require.NoError(t, err)

	err = bagRepo.Save(context.Background(), b)
	require.NoError(t, err)

	return b
}

func makeBagToUpdate(
	t *testing.T,
	status bag.Status,
	bagRepo bag.Repository,
	pickupStart time.Time,
	pickupEnd time.Time,
) *bag.Bag {
	t.Helper()

	win, err := bag.NewPickupWindow(pickupStart, pickupEnd)
	require.NoError(t, err)

	b := bag.From(
		"id",
		"store-id",
		"XDD",
		"",
		1,
		500,
		250,
		status,
		win,
		[]bag.Item{},
		pickupStart,
		pickupStart,
	)

	err = bagRepo.Save(context.Background(), b)
	require.NoError(t, err)

	return b
}

func TestUpdateBag_Success(t *testing.T) {
	deps := newUpdateBagDeps()

	now := deps.clk.Now()

	pickupStart := now.Add(time.Hour * 24)
	pickupEnd := pickupStart.Add(time.Hour)

	b := createBagToUpdate(t, deps.productRepo, deps.bagRepo, pickupStart, pickupEnd)

	input := usecase.UpdateBagInput{
		ID:            b.ID(),
		Title:         "Nova Sacola",
		Description:   "Nova descrição",
		Stock:         2,
		OriginalPrice: 1000,
		Discount:      50,
		PickupWindow: usecase.UpdateBagPickupWindowInput{
			Start: pickupStart,
			End:   pickupEnd,
		},
		Items: []usecase.UpdateBagItemInput{
			{ID: b.Items()[0].ID(), Quantity: 1},
		},
	}

	err := deps.uc.Execute(context.Background(), input)
	require.NoError(t, err)

	updated, err := deps.bagRepo.Find(context.Background(), b.ID())
	require.NoError(t, err)

	assert.Equal(t, input.Title, updated.Title())
	assert.Equal(t, input.Description, updated.Description())
	assert.Equal(t, input.Stock, updated.Stock())
	assert.Equal(t, input.OriginalPrice, updated.OriginalPrice())
	assert.Equal(t, int64(500), updated.FinalPrice())
	assert.Equal(t, input.PickupWindow.Start, updated.PickupWindow().Start())
	assert.Equal(t, input.PickupWindow.End, updated.PickupWindow().End())
	assert.Equal(t, now, updated.UpdatedAt())

	bagItems := make(map[string]bag.Item, len(updated.Items()))
	for _, bagItem := range updated.Items() {
		bagItems[bagItem.ID()] = bagItem
	}

	assert.Len(t, updated.Items(), len(input.Items))
	for _, inputItem := range input.Items {
		bagItem, ok := bagItems[inputItem.ID]

		assert.True(t, ok, "item %s should be in bag", inputItem.ID)
		assert.Equal(t, inputItem.Quantity, bagItem.Quantity())
	}
}

func TestUpdateBag_NotFoundError(t *testing.T) {
	deps := newUpdateBagDeps()

	now := deps.clk.Now()

	input := usecase.UpdateBagInput{
		ID:            "not-found-id",
		Title:         "Sacola não encontrada",
		Stock:         1,
		OriginalPrice: 100,
		Discount:      30,
		PickupWindow: usecase.UpdateBagPickupWindowInput{
			Start: now,
			End:   now.Add(time.Hour),
		},
		Items: []usecase.UpdateBagItemInput{
			{ID: "prod-id", Quantity: 1},
		},
	}

	err := deps.uc.Execute(context.Background(), input)
	require.Error(t, err)

	assert.ErrorIs(t, err, bag.ErrNotFound)
}

func TestUpdateBag_IsClosedError(t *testing.T) {
	deps := newUpdateBagDeps()

	now := deps.clk.Now()

	pickupStart := now.Add(time.Hour * 24)
	pickupEnd := pickupStart.Add(time.Hour * 6)

	b := makeBagToUpdate(t, bag.Closed, deps.bagRepo, pickupStart, pickupEnd)

	input := usecase.UpdateBagInput{
		ID:            b.ID(),
		Title:         "Sacola encerrada",
		Stock:         1,
		OriginalPrice: 100,
		Discount:      30,
		PickupWindow: usecase.UpdateBagPickupWindowInput{
			Start: pickupStart,
			End:   pickupEnd,
		},
		Items: []usecase.UpdateBagItemInput{
			{ID: "prod-id", Quantity: 1},
		},
	}

	err := deps.uc.Execute(context.Background(), input)
	require.Error(t, err)

	assert.ErrorIs(t, err, bag.ErrClosed)

	updated, err := deps.bagRepo.Find(context.Background(), b.ID())
	require.NoError(t, err)

	assert.Equal(t, b.Title(), updated.Title())
	assert.Equal(t, b.Description(), updated.Description())
	assert.Equal(t, b.Stock(), updated.Stock())
	assert.Equal(t, b.OriginalPrice(), updated.OriginalPrice())
	assert.Equal(t, b.FinalPrice(), updated.FinalPrice())
	assert.Equal(t, b.PickupWindow().Start(), updated.PickupWindow().Start())
	assert.Equal(t, b.PickupWindow().End(), updated.PickupWindow().End())
	assert.Equal(t, b.UpdatedAt(), updated.UpdatedAt())
	assert.Equal(t, b.Items(), updated.Items())
}

func TestUpdateBag_IsExpiredError(t *testing.T) {
	deps := newUpdateBagDeps()

	now := deps.clk.Now()

	pickupStart := now.Add(-time.Hour * 24)
	pickupEnd := pickupStart.Add(time.Hour * 6)

	b := makeBagToUpdate(t, bag.Active, deps.bagRepo, pickupStart, pickupEnd)

	input := usecase.UpdateBagInput{
		ID:            b.ID(),
		Title:         "Sacola expirada",
		Stock:         1,
		OriginalPrice: 100,
		Discount:      30,
		PickupWindow: usecase.UpdateBagPickupWindowInput{
			Start: pickupStart,
			End:   pickupEnd,
		},
		Items: []usecase.UpdateBagItemInput{
			{ID: "prod-id", Quantity: 1},
		},
	}

	err := deps.uc.Execute(context.Background(), input)
	require.Error(t, err)

	assert.ErrorIs(t, err, bag.ErrExpired)

	updated, err := deps.bagRepo.Find(context.Background(), b.ID())
	require.NoError(t, err)

	assert.Equal(t, b.Title(), updated.Title())
	assert.Equal(t, b.Description(), updated.Description())
	assert.Equal(t, b.Stock(), updated.Stock())
	assert.Equal(t, b.OriginalPrice(), updated.OriginalPrice())
	assert.Equal(t, b.FinalPrice(), updated.FinalPrice())
	assert.Equal(t, b.PickupWindow().Start(), updated.PickupWindow().Start())
	assert.Equal(t, b.PickupWindow().End(), updated.PickupWindow().End())
	assert.Equal(t, b.UpdatedAt(), updated.UpdatedAt())
	assert.Equal(t, b.Items(), updated.Items())
}

func TestUpdateBag_ProductNotFoundError(t *testing.T) {
	deps := newUpdateBagDeps()

	now := deps.clk.Now()

	pickupStart := now.Add(time.Hour * 24)
	pickupEnd := pickupStart.Add(time.Hour)

	b := createBagToUpdate(t, deps.productRepo, deps.bagRepo, pickupStart, pickupEnd)

	input := usecase.UpdateBagInput{
		ID:            b.ID(),
		Title:         "Sacola com produto não encontrado",
		Stock:         1,
		OriginalPrice: 1000,
		Discount:      50,
		PickupWindow: usecase.UpdateBagPickupWindowInput{
			Start: pickupStart,
			End:   pickupEnd,
		},
		Items: []usecase.UpdateBagItemInput{
			{ID: b.Items()[0].ID(), Quantity: 1},
			{ID: "not-found-id", Quantity: 1},
		},
	}

	err := deps.uc.Execute(context.Background(), input)
	require.Error(t, err)

	assert.ErrorIs(t, err, product.ErrNotFound)

	updated, err := deps.bagRepo.Find(context.Background(), b.ID())
	require.NoError(t, err)

	assert.Equal(t, b.Title(), updated.Title())
	assert.Equal(t, b.Description(), updated.Description())
	assert.Equal(t, b.Stock(), updated.Stock())
	assert.Equal(t, b.OriginalPrice(), updated.OriginalPrice())
	assert.Equal(t, b.FinalPrice(), updated.FinalPrice())
	assert.Equal(t, b.PickupWindow().Start(), updated.PickupWindow().Start())
	assert.Equal(t, b.PickupWindow().End(), updated.PickupWindow().End())
	assert.Equal(t, b.UpdatedAt(), updated.UpdatedAt())
	assert.Equal(t, b.Items(), updated.Items())
}

func TestUpdateBag_DifferentStoreProductsError(t *testing.T) {
	deps := newUpdateBagDeps()

	now := deps.clk.Now()

	pickupStart := now.Add(time.Hour * 24)
	pickupEnd := pickupStart.Add(time.Hour)

	otherProd, err := product.New(
		"other-prod-id",
		"other-store-id",
		"Coxinha",
		"",
		pickupStart,
	)
	require.NoError(t, err)

	err = deps.productRepo.Save(context.Background(), otherProd)
	require.NoError(t, err)

	b := createBagToUpdate(t, deps.productRepo, deps.bagRepo, pickupStart, pickupEnd)

	input := usecase.UpdateBagInput{
		ID:            b.ID(),
		Title:         "Sacola com produtos de outra loja",
		Stock:         1,
		OriginalPrice: 1000,
		Discount:      50,
		PickupWindow: usecase.UpdateBagPickupWindowInput{
			Start: pickupStart,
			End:   pickupEnd,
		},
		Items: []usecase.UpdateBagItemInput{
			{ID: b.Items()[0].ID(), Quantity: 1},
			{ID: otherProd.ID(), Quantity: 1},
		},
	}

	err = deps.uc.Execute(context.Background(), input)
	require.Error(t, err)

	assert.ErrorIs(t, err, product.ErrNotFound)

	updated, err := deps.bagRepo.Find(context.Background(), b.ID())
	require.NoError(t, err)

	assert.Equal(t, b.Title(), updated.Title())
	assert.Equal(t, b.Description(), updated.Description())
	assert.Equal(t, b.Stock(), updated.Stock())
	assert.Equal(t, b.OriginalPrice(), updated.OriginalPrice())
	assert.Equal(t, b.FinalPrice(), updated.FinalPrice())
	assert.Equal(t, b.PickupWindow().Start(), updated.PickupWindow().Start())
	assert.Equal(t, b.PickupWindow().End(), updated.PickupWindow().End())
	assert.Equal(t, b.UpdatedAt(), updated.UpdatedAt())
	assert.Equal(t, b.Items(), updated.Items())
}

func TestUpdateBag_InvalidPickupIntervalError(t *testing.T) {
	t.Run("pickup window interval is less than 30 minutes", func(t *testing.T) {
		deps := newUpdateBagDeps()

		now := deps.clk.Now()

		pickupStart := now.Add(time.Hour * 24)
		pickupEnd := pickupStart.Add(time.Hour)

		b := createBagToUpdate(t, deps.productRepo, deps.bagRepo, pickupStart, pickupEnd)

		input := usecase.UpdateBagInput{
			ID:            b.ID(),
			Title:         "Sacola com intervalo muito curto",
			Stock:         1,
			OriginalPrice: 1000,
			Discount:      50,
			PickupWindow: usecase.UpdateBagPickupWindowInput{
				Start: pickupStart,
				End:   pickupStart.Add(time.Minute * 29),
			},
			Items: []usecase.UpdateBagItemInput{
				{ID: b.Items()[0].ID(), Quantity: 1},
			},
		}

		err := deps.uc.Execute(context.Background(), input)
		require.Error(t, err)

		assert.ErrorIs(t, err, bag.ErrInvalidPickupWindowInterval)

		updated, err := deps.bagRepo.Find(context.Background(), b.ID())
		require.NoError(t, err)

		assert.Equal(t, b.Title(), updated.Title())
		assert.Equal(t, b.Description(), updated.Description())
		assert.Equal(t, b.Stock(), updated.Stock())
		assert.Equal(t, b.OriginalPrice(), updated.OriginalPrice())
		assert.Equal(t, b.FinalPrice(), updated.FinalPrice())
		assert.Equal(t, b.PickupWindow().Start(), updated.PickupWindow().Start())
		assert.Equal(t, b.PickupWindow().End(), updated.PickupWindow().End())
		assert.Equal(t, b.UpdatedAt(), updated.UpdatedAt())
		assert.Equal(t, b.Items(), updated.Items())
	})

	t.Run("pickup window interval is greater than 6 hours", func(t *testing.T) {
		deps := newUpdateBagDeps()

		now := deps.clk.Now()

		pickupStart := now.Add(time.Hour * 24)
		pickupEnd := pickupStart.Add(time.Hour)

		b := createBagToUpdate(t, deps.productRepo, deps.bagRepo, pickupStart, pickupEnd)

		input := usecase.UpdateBagInput{
			ID:            b.ID(),
			Title:         "Sacola com intervalo muito longo",
			Stock:         1,
			OriginalPrice: 1000,
			Discount:      50,
			PickupWindow: usecase.UpdateBagPickupWindowInput{
				Start: pickupStart,
				End:   pickupStart.Add(time.Hour*6 + time.Minute),
			},
			Items: []usecase.UpdateBagItemInput{
				{ID: b.Items()[0].ID(), Quantity: 1},
			},
		}

		err := deps.uc.Execute(context.Background(), input)
		require.Error(t, err)

		assert.ErrorIs(t, err, bag.ErrInvalidPickupWindowInterval)

		updated, err := deps.bagRepo.Find(context.Background(), b.ID())
		require.NoError(t, err)

		assert.Equal(t, b.Title(), updated.Title())
		assert.Equal(t, b.Description(), updated.Description())
		assert.Equal(t, b.Stock(), updated.Stock())
		assert.Equal(t, b.OriginalPrice(), updated.OriginalPrice())
		assert.Equal(t, b.FinalPrice(), updated.FinalPrice())
		assert.Equal(t, b.PickupWindow().Start(), updated.PickupWindow().Start())
		assert.Equal(t, b.PickupWindow().End(), updated.PickupWindow().End())
		assert.Equal(t, b.UpdatedAt(), updated.UpdatedAt())
		assert.Equal(t, b.Items(), updated.Items())
	})
}
