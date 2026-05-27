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

type createBagDeps struct {
	clk         clock.Clock
	bagRepo     bag.Repository
	productRepo product.Repository
	uc          *usecase.CreateBag
}

func newCreateBagDeps() *createBagDeps {
	clk := clock.NewFakeClock()
	productRepo := productRepository.NewInMemRepo()
	bagRepo := repository.NewInMemRepo()
	uc := usecase.NewCreateBag(clk, bagRepo, productRepo)

	return &createBagDeps{
		clk:         clk,
		productRepo: productRepo,
		bagRepo:     bagRepo,
		uc:          uc,
	}
}

func createProductsToCreateBag(t *testing.T, productRepo product.Repository) []*product.Product {
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

func TestCreateBag_Success(t *testing.T) {
	deps := newCreateBagDeps()

	now := deps.clk.Now()

	prods := createProductsToCreateBag(t, deps.productRepo)

	pickupStart := time.Date(2025, time.June, 1, 14, 0, 0, 0, time.UTC)
	pickupEnd := pickupStart.Add(time.Hour)

	input := usecase.CreateBagInput{
		StoreID:       "store-id",
		Title:         "Sacola salgada",
		Description:   "Sacola com salgados para o café da tarde",
		Stock:         2,
		OriginalPrice: 1500,
		Discount:      30,
		PickupWindow: usecase.CreateBagPickupWindowInput{
			Start: pickupStart,
			End:   pickupEnd,
		},
		Items: []usecase.CreateBagItemInput{
			{ID: prods[0].ID(), Quantity: 3},
			{ID: prods[1].ID(), Quantity: 1},
			{ID: prods[2].ID(), Quantity: 2},
		},
	}

	output, err := deps.uc.Execute(context.Background(), input)
	require.NoError(t, err)

	b, err := deps.bagRepo.Find(context.Background(), output.ID)
	require.NoError(t, err)

	assert.True(t, uuid.IsValid(output.ID), "bag id should be an uuid")
	assert.Equal(t, input.StoreID, b.StoreID())
	assert.Equal(t, input.Title, b.Title())
	assert.Equal(t, input.Description, b.Description())
	assert.Equal(t, input.Stock, b.Stock())
	assert.Equal(t, input.OriginalPrice, b.OriginalPrice())
	assert.Equal(t, int64(1050), b.FinalPrice())
	assert.Equal(t, input.PickupWindow.Start, b.PickupWindow().Start())
	assert.Equal(t, input.PickupWindow.End, b.PickupWindow().End())
	assert.Equal(t, now, b.CreatedAt())
	assert.Equal(t, now, b.UpdatedAt())

	bagItems := make(map[string]bag.Item, len(b.Items()))
	for _, bagItem := range b.Items() {
		bagItems[bagItem.ID()] = bagItem
	}

	assert.Len(t, b.Items(), len(input.Items))
	for _, inputItem := range input.Items {
		bagItem, ok := bagItems[inputItem.ID]

		assert.True(t, ok, "item %s should be in bag", inputItem.ID)
		assert.Equal(t, inputItem.Quantity, bagItem.Quantity())
	}
}

func TestCreateBag_ProductNotFoundError(t *testing.T) {
	deps := newCreateBagDeps()

	prods := createProductsToCreateBag(t, deps.productRepo)

	pickupStart := time.Date(2025, time.June, 1, 14, 0, 0, 0, time.UTC)
	pickupEnd := pickupStart.Add(time.Hour)

	input := usecase.CreateBagInput{
		StoreID:       "store-id",
		Title:         "Sacola mista",
		Description:   "Essa é uma sacola",
		Stock:         1,
		OriginalPrice: 1000,
		Discount:      50,
		PickupWindow: usecase.CreateBagPickupWindowInput{
			Start: pickupStart,
			End:   pickupEnd,
		},
		Items: []usecase.CreateBagItemInput{
			{ID: prods[0].ID(), Quantity: 2},
			{ID: prods[1].ID(), Quantity: 1},
			{ID: prods[2].ID(), Quantity: 3},
			{ID: "prod-4", Quantity: 2},
		},
	}

	_, err := deps.uc.Execute(context.Background(), input)
	require.Error(t, err)

	assert.ErrorIs(t, err, product.ErrNotFound)
}

func TestCreateBag_DifferentStoreProductsError(t *testing.T) {
	deps := newCreateBagDeps()

	prods := createProductsToCreateBag(t, deps.productRepo)

	pickupStart := time.Date(2025, time.June, 1, 14, 0, 0, 0, time.UTC)
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

	input := usecase.CreateBagInput{
		StoreID:       "store-id",
		Title:         "Sacola com produtos de outra loja",
		Stock:         1,
		OriginalPrice: 1000,
		Discount:      50,
		PickupWindow: usecase.CreateBagPickupWindowInput{
			Start: pickupStart,
			End:   pickupEnd,
		},
		Items: []usecase.CreateBagItemInput{
			{ID: prods[0].ID(), Quantity: 2},
			{ID: prods[1].ID(), Quantity: 1},
			{ID: otherProd.ID(), Quantity: 3},
		},
	}

	_, err = deps.uc.Execute(context.Background(), input)
	require.Error(t, err)

	assert.ErrorIs(t, err, product.ErrNotFound)
}

func TestCreateBag_InvalidPickupWindowIntervalError(t *testing.T) {
	t.Run("pickup window interval is less than 30 minutes", func(t *testing.T) {
		deps := newCreateBagDeps()

		prods := createProductsToCreateBag(t, deps.productRepo)

		pickupStart := time.Date(2025, time.June, 1, 14, 0, 0, 0, time.UTC)
		pickupEnd := pickupStart.Add(time.Minute * 29)

		input := usecase.CreateBagInput{
			StoreID:       "store-id",
			Title:         "Sacola com intervalo muito curto",
			Stock:         1,
			OriginalPrice: 1000,
			Discount:      50,
			PickupWindow: usecase.CreateBagPickupWindowInput{
				Start: pickupStart,
				End:   pickupEnd,
			},
			Items: []usecase.CreateBagItemInput{
				{ID: prods[0].ID(), Quantity: 2},
				{ID: prods[1].ID(), Quantity: 1},
			},
		}

		_, err := deps.uc.Execute(context.Background(), input)
		require.Error(t, err)

		assert.ErrorIs(t, err, bag.ErrInvalidPickupWindowInterval)
	})

	t.Run("pickup window interval is greater than 6 hours", func(t *testing.T) {
		deps := newCreateBagDeps()

		prods := createProductsToCreateBag(t, deps.productRepo)

		pickupStart := time.Date(2025, time.June, 1, 14, 0, 0, 0, time.UTC)
		pickupEnd := pickupStart.Add(time.Hour*6 + time.Minute)

		input := usecase.CreateBagInput{
			StoreID:       "store-id",
			Title:         "Sacola com intervalo muito longo",
			Stock:         1,
			OriginalPrice: 1000,
			Discount:      50,
			PickupWindow: usecase.CreateBagPickupWindowInput{
				Start: pickupStart,
				End:   pickupEnd,
			},
			Items: []usecase.CreateBagItemInput{
				{ID: prods[0].ID(), Quantity: 2},
				{ID: prods[1].ID(), Quantity: 1},
			},
		}

		_, err := deps.uc.Execute(context.Background(), input)
		require.Error(t, err)

		assert.ErrorIs(t, err, bag.ErrInvalidPickupWindowInterval)
	})
}
