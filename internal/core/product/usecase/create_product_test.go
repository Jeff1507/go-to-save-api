package usecase_test

import (
	"context"
	"strings"
	"testing"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/product"
	"github.com/Jeff1507/go-to-save-api/internal/core/product/repository"
	"github.com/Jeff1507/go-to-save-api/internal/core/product/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/storage"
	"github.com/Jeff1507/go-to-save-api/internal/pkg/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type createProductDeps struct {
	clk  clock.Clock
	st   *storage.FakeStorage
	repo product.Repository
	uc   *usecase.CreateProduct
}

func newCreateProductDeps() *createProductDeps {
	clk := clock.NewFakeClock()
	st := storage.NewFakeStorage()
	repo := repository.NewInMemRepo()
	uc := usecase.NewCreateProduct(clk, st, repo)

	return &createProductDeps{
		clk:  clk,
		st:   st,
		repo: repo,
		uc:   uc,
	}
}

func TestCreateProduct_WithImage_Success(t *testing.T) {
	deps := newCreateProductDeps()

	input := usecase.CreateProductInput{
		StoreID: "store-id",
		Name:    "Coxinha",
		Image: &usecase.CreateProductImageInput{
			Content:     strings.NewReader("fake-img-content"),
			ContentType: ".jpg",
		},
	}

	output, err := deps.uc.Execute(context.Background(), input)
	require.NoError(t, err)

	prod, err := deps.repo.Find(context.Background(), output.ID)
	require.NoError(t, err)

	assert.True(t, uuid.IsValid(prod.ID()), "product id should be an uuid")
	assert.Equal(t, input.StoreID, prod.StoreID())
	assert.Equal(t, input.Name, prod.Name())
	assert.Equal(t, deps.clk.Now(), prod.CreatedAt())
	assert.Equal(t, deps.clk.Now(), prod.UpdatedAt())
	assert.True(t, deps.st.Uploaded, "upload should be called")
	assert.Equal(t, deps.st.UploadedKey, prod.ImageKey())
}

func TestCreateProduct_WithoutImage_Success(t *testing.T) {
	deps := newCreateProductDeps()

	input := usecase.CreateProductInput{
		StoreID: "store-id",
		Name:    "Coxinha",
		Image:   nil,
	}

	output, err := deps.uc.Execute(context.Background(), input)
	require.NoError(t, err)

	prod, err := deps.repo.Find(context.Background(), output.ID)
	require.NoError(t, err)

	assert.True(t, uuid.IsValid(prod.ID()), "product id should be an uuid")
	assert.Equal(t, input.StoreID, prod.StoreID())
	assert.Equal(t, input.Name, prod.Name())
	assert.Equal(t, deps.clk.Now(), prod.CreatedAt())
	assert.Equal(t, deps.clk.Now(), prod.UpdatedAt())
	assert.False(t, deps.st.Uploaded, "upload should not be called")
	assert.Empty(t, prod.ImageKey(), "product image key should be empty")
}

func TestCreateProduct_AlreadyExistsError(t *testing.T) {
	deps := newCreateProductDeps()

	existed, err := product.New("id", "store-id", "Coxinha", "", deps.clk.Now())
	require.NoError(t, err)

	err = deps.repo.Save(context.Background(), existed)
	require.NoError(t, err)

	input := usecase.CreateProductInput{
		StoreID: existed.StoreID(),
		Name:    existed.Name(),
		Image:   nil,
	}

	output, err := deps.uc.Execute(context.Background(), input)
	require.Error(t, err)

	assert.ErrorIs(t, err, product.ErrAlreadyExists)

	_, err = deps.repo.Find(context.Background(), output.ID)
	require.Error(t, err)

	assert.ErrorIs(t, err, product.ErrNotFound)
}
