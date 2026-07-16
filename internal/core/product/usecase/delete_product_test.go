package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/Jeff1507/go-to-save-api/internal/core/product"
	"github.com/Jeff1507/go-to-save-api/internal/core/product/repository"
	"github.com/Jeff1507/go-to-save-api/internal/core/product/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type deleteProductDeps struct {
	repo product.Repository
	uc   *usecase.DeleteProduct
}

func newDeleteProductDeps() *deleteProductDeps {
	repo := repository.NewInMemRepo()
	uc := usecase.NewDeleteProduct(repo)

	return &deleteProductDeps{
		repo: repo,
		uc:   uc,
	}
}

func TestDeleteProduct_Success(t *testing.T) {
	deps := newDeleteProductDeps()
	date := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)

	prod, err := product.New("id", "store-id", "Coxinha", "coxinha.png", date)
	require.NoError(t, err)

	err = deps.repo.Save(context.Background(), prod)
	require.NoError(t, err)

	err = deps.uc.Execute(context.Background(), prod.ID())
	require.NoError(t, err)

	_, err = deps.repo.Find(context.Background(), prod.ID())
	require.Error(t, err)

	assert.ErrorIs(t, err, product.ErrNotFound)
}

func TestDeleteProduct_NotFound_Success(t *testing.T) {
	deps := newDeleteProductDeps()

	err := deps.uc.Execute(context.Background(), "not-found-id")
	require.NoError(t, err)
}
