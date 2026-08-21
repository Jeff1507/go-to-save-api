package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/product"
	"github.com/Jeff1507/go-to-save-api/internal/infra/storage"
)

type UpdateProductInput struct {
	ID          string
	StoreID     string
	Name        string
	Image       ImageInput
	RemoveImage bool
}

type UpdateProduct struct {
	productRepo product.Repository
	storage     storage.Storage
	clock       clock.Clock
}

func NewUpdateProduct(
	productRepo product.Repository,
	storage storage.Storage,
	clock clock.Clock,
) *UpdateProduct {
	return &UpdateProduct{
		productRepo: productRepo,
		storage:     storage,
		clock:       clock,
	}
}

func (uc *UpdateProduct) Execute(ctx context.Context, input UpdateProductInput) error {
	now := uc.clock.Now()

	imageRemoved := input.RemoveImage
	imageChanged := input.Image.Content != nil && input.Image.ContentType != ""

	if imageRemoved && imageChanged {
		return errors.New("cannot change and remove image at the same time")
	}

	prod, err := uc.productRepo.Find(ctx, input.ID)
	if err != nil {
		return fmt.Errorf("find product by id: %w", err)
	}

	if prod.StoreID() != input.StoreID {
		return product.ErrNotFound
	}

	if err := prod.ChangeName(input.Name, now); err != nil {
		return fmt.Errorf("change name: %w", err)
	}

	if imageRemoved {
		if err := prod.RemoveImage(now); err != nil {
			return fmt.Errorf("remove image: %w", err)
		}
	}

	var newImageKey string
	if imageChanged {
		newImageKey, err = uc.storage.Upload(
			input.Image.Content,
			input.Image.ContentType,
			storage.ProductsPath,
		)
		if err != nil {
			if errors.Is(err, storage.ErrContentTypeNotAllowed) {
				return product.ErrInvalidImageFormat
			}

			return fmt.Errorf("upload image: %w", err)
		}

		if err := prod.ChangeImage(newImageKey, now); err != nil {
			uc.storage.Delete(newImageKey)
			return fmt.Errorf("change image: %w", err)
		}
	}

	if err := uc.productRepo.Update(ctx, prod); err != nil {
		uc.storage.Delete(newImageKey)
		return fmt.Errorf("update product: %w", err)
	}

	// não deleta a imagem antiga — pode estar em uso por bags existentes
	// (snapshot em bag_items referencia a image_key anterior)
	// TODO: job de limpeza periódico para imagens órfãs

	return nil
}
