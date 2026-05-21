package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/core/product"
	"github.com/Jeff1507/go-to-save-api/internal/infra/storage"
)

type DeleteProduct struct {
	fileStorage storage.Storage
	productRepo product.Repository
}

func NewDeleteProduct(fileStorage storage.Storage, productRepo product.Repository) *DeleteProduct {
	return &DeleteProduct{
		fileStorage: fileStorage,
		productRepo: productRepo,
	}
}

func (uc *DeleteProduct) Execute(ctx context.Context, id string) error {
	prod, err := uc.productRepo.Find(ctx, id)
	if err != nil {
		if errors.Is(err, product.ErrNotFound) {
			return nil
		}

		return fmt.Errorf("find product by id: %v", err)
	}

	if err := uc.productRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete product by id: %v", err)
	}

	uc.fileStorage.Delete(prod.ImageKey())

	return nil
}
