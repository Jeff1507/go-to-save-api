package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/product"
	"github.com/Jeff1507/go-to-save-api/internal/infra/storage"
)

type UpdateProductImageInput struct {
	Content     io.Reader
	ContentType string
}

type UpdateProductInput struct {
	ID    string
	Name  string
	Image *UpdateProductImageInput
}

type UpdateProduct struct {
	clock       clock.Clock
	fileStorage storage.Storage
	productRepo product.Repository
}

func NewUpdateProduct(
	clock clock.Clock,
	fileStorage storage.Storage,
	productRepo product.Repository,
) *UpdateProduct {
	return &UpdateProduct{
		clock:       clock,
		fileStorage: fileStorage,
		productRepo: productRepo,
	}
}

func (uc *UpdateProduct) Execute(ctx context.Context, input UpdateProductInput) error {
	now := uc.clock.Now()

	prod, err := uc.productRepo.Find(ctx, input.ID)
	if err != nil {
		if errors.Is(err, product.ErrNotFound) {
			return fmt.Errorf("find product by id: %w", err)
		}

		return fmt.Errorf("find product by id: %v", err)
	}

	nameChanged := prod.Name() != input.Name
	imageChanged := input.Image != nil

	if !nameChanged && !imageChanged {
		return nil
	}

	if nameChanged {
		if err := prod.ChangeName(input.Name, now); err != nil {
			return fmt.Errorf("change product name: %v", err)
		}
	}

	oldImgKey := prod.ImageKey()

	var imgKey string
	if imageChanged {
		imgKey, err = uc.fileStorage.Upload(
			input.Image.Content,
			input.Image.ContentType,
			storage.ProductsPath,
		)
		if err != nil {
			if errors.Is(err, storage.ErrContentTypeNotAllowed) {
				return product.ErrInvalidImageFormat
			}

			return fmt.Errorf("upload image: %v", err)
		}

		if err := prod.ChangeImage(imgKey, now); err != nil {
			uc.fileStorage.Delete(imgKey)
			return fmt.Errorf("change product image: %v", err)
		}
	}

	if err := uc.productRepo.Update(ctx, prod); err != nil {
		uc.fileStorage.Delete(imgKey)

		if errors.Is(err, product.ErrAlreadyExists) {
			return fmt.Errorf("update product: %w", err)
		}

		return fmt.Errorf("update product: %v", err)
	}

	if imageChanged {
		uc.fileStorage.Delete(oldImgKey)
	}

	return nil
}
