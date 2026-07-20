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
	"github.com/Jeff1507/go-to-save-api/internal/pkg/uuid"
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

	prod, err := product.New("prod-id-abc", "store-id", "Bolo de morango", "bolo.png", date)
	require.NoError(t, err)

	err = productRepo.Save(context.Background(), prod)
	require.NoError(t, err)

	win, err := bag.NewPickupWindow(pickupStart, pickupEnd)
	require.NoError(t, err)

	item, err := bag.NewItem("item-id", prod.Name(), prod.ImageKey(), 2)
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

func createProductsToUpdateBag(t *testing.T, productRepo product.Repository) []*product.Product {
	t.Helper()

	date := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)

	prod1, err := product.New("prod-1", "store-id", "Coxinha", "coxinha.png", date)
	require.NoError(t, err)

	err = productRepo.Save(context.Background(), prod1)
	require.NoError(t, err)

	prod2, err := product.New("prod-2", "store-id", "Croissant", "croissant.png", date)
	require.NoError(t, err)

	err = productRepo.Save(context.Background(), prod2)
	require.NoError(t, err)

	prod3, err := product.New("prod-3", "store-id", "Pão de batata", "", date)
	require.NoError(t, err)

	err = productRepo.Save(context.Background(), prod3)
	require.NoError(t, err)

	return []*product.Product{prod1, prod2, prod3}
}

func TestUpdateBag_Success(t *testing.T) {
	deps := newUpdateBagDeps()

	now := deps.clk.Now()

	pickupStart := now.Add(time.Hour * 24)
	pickupEnd := pickupStart.Add(time.Hour)

	b := createBagToUpdate(t, deps.productRepo, deps.bagRepo, pickupStart, pickupEnd)
	prods := createProductsToUpdateBag(t, deps.productRepo)

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
			{ProductID: prods[0].ID(), Quantity: 1},
			{ProductID: prods[1].ID(), Quantity: 2},
			{ProductID: prods[2].ID(), Quantity: 1},
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

	bagItemsByName := make(map[string]bag.Item, len(updated.Items()))
	for _, item := range updated.Items() {
		bagItemsByName[item.Name()] = item
	}

	assert.Len(t, updated.Items(), len(input.Items))

	for _, prod := range prods {
		item, ok := bagItemsByName[prod.Name()]

		assert.True(t, ok, "item %s should be in bag", prod.Name())
		assert.True(t, uuid.IsValid(item.ID()))
		assert.Equal(t, prod.Name(), item.Name())
		assert.Equal(t, prod.ImageKey(), item.ImageKey())
	}

	inputByProductID := make(map[string]usecase.UpdateBagItemInput)
	for _, item := range input.Items {
		inputByProductID[item.ProductID] = item
	}

	for _, prod := range prods {
		item := bagItemsByName[prod.Name()]
		inputItem := inputByProductID[prod.ID()]
		assert.Equal(t, inputItem.Quantity, item.Quantity())
	}
}

func TestUpdateBag_NotFoundError(t *testing.T) {
	deps := newUpdateBagDeps()

	now := deps.clk.Now()

	prods := createProductsToUpdateBag(t, deps.productRepo)

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
			{ProductID: prods[0].ID(), Quantity: 1},
			{ProductID: prods[1].ID(), Quantity: 2},
			{ProductID: prods[2].ID(), Quantity: 1},
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
	prods := createProductsToUpdateBag(t, deps.productRepo)

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
			{ProductID: prods[0].ID(), Quantity: 1},
			{ProductID: prods[1].ID(), Quantity: 2},
			{ProductID: prods[2].ID(), Quantity: 1},
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
	prods := createProductsToUpdateBag(t, deps.productRepo)

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
			{ProductID: prods[0].ID(), Quantity: 1},
			{ProductID: prods[1].ID(), Quantity: 2},
			{ProductID: prods[2].ID(), Quantity: 1},
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
	prods := createProductsToUpdateBag(t, deps.productRepo)

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
			{ProductID: prods[0].ID(), Quantity: 1},
			{ProductID: prods[1].ID(), Quantity: 2},
			{ProductID: prods[2].ID(), Quantity: 3},
			{ProductID: "prod-id-4", Quantity: 1},
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
	prods := createProductsToUpdateBag(t, deps.productRepo)

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
			{ProductID: prods[0].ID(), Quantity: 1},
			{ProductID: prods[1].ID(), Quantity: 2},
			{ProductID: prods[2].ID(), Quantity: 3},
			{ProductID: otherProd.ID(), Quantity: 1},
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
		prods := createProductsToUpdateBag(t, deps.productRepo)

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
				{ProductID: prods[0].ID(), Quantity: 1},
				{ProductID: prods[1].ID(), Quantity: 2},
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
		prods := createProductsToUpdateBag(t, deps.productRepo)

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
				{ProductID: prods[0].ID(), Quantity: 1},
				{ProductID: prods[1].ID(), Quantity: 2},
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
