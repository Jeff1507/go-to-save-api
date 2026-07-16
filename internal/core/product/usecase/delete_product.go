package usecase

import (
	"context"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/core/product"
)

type DeleteProduct struct {
	productRepo product.Repository
}

func NewDeleteProduct(productRepo product.Repository) *DeleteProduct {
	return &DeleteProduct{
		productRepo: productRepo,
	}
}

func (uc *DeleteProduct) Execute(ctx context.Context, id string) error {
	if err := uc.productRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete product by id: %v", err)
	}

	// não deleta a imagem antiga — pode estar em uso por bags existentes
	// (snapshot em bag_items referencia a image_key anterior)
	// TODO: job de limpeza periódico para imagens órfãs

	return nil
}
