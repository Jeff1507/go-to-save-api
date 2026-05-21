package usecase_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/product"
	"github.com/Jeff1507/go-to-save-api/internal/core/product/repository"
	"github.com/Jeff1507/go-to-save-api/internal/core/product/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type updateProductDeps struct {
	clk  clock.Clock
	st   *storage.FakeStorage
	repo product.Repository
	uc   *usecase.UpdateProduct
}

func newUpdateProductDeps() *updateProductDeps {
	clk := clock.NewFakeClock()
	st := storage.NewFakeStorage()
	repo := repository.NewInMemRepo()
	uc := usecase.NewUpdateProduct(clk, st, repo)

	return &updateProductDeps{
		clk:  clk,
		st:   st,
		repo: repo,
		uc:   uc,
	}
}

func createProductToUpdate(t *testing.T, repo product.Repository) *product.Product {
	t.Helper()

	date := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)

	prod, err := product.New("id", "store-id", "Coxinha", "coxinha.png", date)
	require.NoError(t, err)

	err = repo.Save(context.Background(), prod)
	require.NoError(t, err)

	return prod
}

func TestUpdateProduct_WithImage_Success(t *testing.T) {
	deps := newUpdateProductDeps()

	prod := createProductToUpdate(t, deps.repo)
	oldImgKey := prod.ImageKey()

	input := usecase.UpdateProductInput{
		ID:   prod.ID(),
		Name: "Coxinha de frango",
		Image: &usecase.UpdateProductImageInput{
			Content:     strings.NewReader("fake-image-content"),
			ContentType: ".jpg",
		},
	}

	err := deps.uc.Execute(context.Background(), input)
	require.NoError(t, err)

	updated, err := deps.repo.Find(context.Background(), prod.ID())
	require.NoError(t, err)

	assert.Equal(t, input.Name, updated.Name())
	assert.Equal(t, deps.clk.Now(), updated.UpdatedAt())
	assert.True(t, deps.st.Uploaded, "upload should be called")
	assert.True(t, deps.st.Deleted, "delete should be called")
	assert.Equal(t, deps.st.UploadedKey, updated.ImageKey())
	assert.Equal(t, deps.st.DeletedKey, oldImgKey)
}

func TestUpdateProduct_WithoutImage_Success(t *testing.T) {
	deps := newUpdateProductDeps()

	prod := createProductToUpdate(t, deps.repo)

	input := usecase.UpdateProductInput{
		ID:    prod.ID(),
		Name:  "Coxinha de frango",
		Image: nil,
	}

	err := deps.uc.Execute(context.Background(), input)
	require.NoError(t, err)

	updated, err := deps.repo.Find(context.Background(), prod.ID())
	require.NoError(t, err)

	assert.Equal(t, input.Name, updated.Name())
	assert.Equal(t, deps.clk.Now(), updated.UpdatedAt())
	assert.Equal(t, prod.ImageKey(), updated.ImageKey())
	assert.False(t, deps.st.Uploaded, "upload should not be called")
	assert.False(t, deps.st.Deleted, "delete should not be called")
}

func TestUpdateProduct_SameName_Success(t *testing.T) {
	deps := newUpdateProductDeps()

	prod := createProductToUpdate(t, deps.repo)
	oldImgKey := prod.ImageKey()

	input := usecase.UpdateProductInput{
		ID:   prod.ID(),
		Name: prod.Name(),
		Image: &usecase.UpdateProductImageInput{
			Content:     strings.NewReader("fake-image-content"),
			ContentType: ".jpg",
		},
	}

	err := deps.uc.Execute(context.Background(), input)
	require.NoError(t, err)

	updated, err := deps.repo.Find(context.Background(), prod.ID())
	require.NoError(t, err)

	assert.Equal(t, prod.Name(), updated.Name())
	assert.Equal(t, deps.clk.Now(), updated.UpdatedAt())
	assert.True(t, deps.st.Uploaded, "upload should be called")
	assert.True(t, deps.st.Deleted, "delete should be called")
	assert.Equal(t, deps.st.UploadedKey, updated.ImageKey())
	assert.Equal(t, deps.st.DeletedKey, oldImgKey)
}

func TestUpdateProduct_NoFoundError(t *testing.T) {
	deps := newUpdateProductDeps()

	input := usecase.UpdateProductInput{
		ID:    "not-found-id",
		Name:  "Bolo de morango",
		Image: nil,
	}

	err := deps.uc.Execute(context.Background(), input)
	require.Error(t, err)

	assert.ErrorIs(t, err, product.ErrNotFound)
}

func TestUpdateProduct_AlreadyExistsError(t *testing.T) {
	deps := newUpdateProductDeps()

	existed, err := product.New("existed-id", "store-id", "Pão de batata", "", deps.clk.Now())
	require.NoError(t, err)

	err = deps.repo.Save(context.Background(), existed)
	require.NoError(t, err)

	prod := createProductToUpdate(t, deps.repo)

	input := usecase.UpdateProductInput{
		ID:   prod.ID(),
		Name: existed.Name(),
		Image: &usecase.UpdateProductImageInput{
			Content:     strings.NewReader("fake-image-content"),
			ContentType: ".png",
		},
	}

	err = deps.uc.Execute(context.Background(), input)
	require.Error(t, err)

	assert.ErrorIs(t, err, product.ErrAlreadyExists)

	notUpdated, err := deps.repo.Find(context.Background(), prod.ID())

	assert.Equal(t, prod.Name(), notUpdated.Name())
	assert.Equal(t, prod.ImageKey(), notUpdated.ImageKey())
	assert.Equal(t, prod.UpdatedAt(), notUpdated.UpdatedAt())
	assert.False(t, deps.st.Uploaded, "upload should not be called")
	assert.False(t, deps.st.Deleted, "delete should not be called")
}
