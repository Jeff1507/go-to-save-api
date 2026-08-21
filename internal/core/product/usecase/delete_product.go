package usecase

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/core/product"
)

type DeleteProductInput struct {
	ID      string
	StoreID string
}

type DeleteProduct struct {
	productRepo product.Repository
}

func NewDeleteProduct(productRepo product.Repository) *DeleteProduct {
	return &DeleteProduct{
		productRepo: productRepo,
	}
}

func (uc *DeleteProduct) Execute(ctx context.Context, input DeleteProductInput) error {
	prod, err := uc.productRepo.Find(ctx, input.ID)
	if err != nil {
		return fmt.Errorf("find product by id: %w", err)
	}

	if prod.StoreID() != input.StoreID {
		return product.ErrNotFound
	}

	if err := uc.productRepo.Delete(ctx, input.ID); err != nil {
		return fmt.Errorf("delete product: %w", err)
	}

	// não deleta a imagem antiga — pode estar em uso por bags existentes
	// (snapshot em bag_items referencia a image_key anterior)
	// TODO: job de limpeza periódico para imagens órfãs

	return nil
}
