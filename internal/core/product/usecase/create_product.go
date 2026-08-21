package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/product"
	"github.com/Jeff1507/go-to-save-api/internal/infra/storage"
	"github.com/Jeff1507/go-to-save-api/internal/pkg/uuid"
)

type CreateProductInput struct {
	StoreID string
	Name    string
	Image   ImageInput
}

type CreateProductOutput struct {
	ID string
}

type CreateProduct struct {
	productRepo product.Repository
	storage     storage.Storage
	clock       clock.Clock
}

func NewCreateProduct(
	productRepo product.Repository,
	storage storage.Storage,
	clock clock.Clock,
) *CreateProduct {
	return &CreateProduct{
		productRepo: productRepo,
		storage:     storage,
		clock:       clock,
	}
}

func (uc *CreateProduct) Execute(
	ctx context.Context,
	input CreateProductInput,
) (CreateProductOutput, error) {
	now := uc.clock.Now()

	var imageKey string
	if input.Image.Content != nil && input.Image.ContentType != "" {
		key, err := uc.storage.Upload(
			input.Image.Content,
			input.Image.ContentType,
			storage.ProductsPath,
		)
		if err != nil {
			if errors.Is(err, storage.ErrContentTypeNotAllowed) {
				return CreateProductOutput{}, product.ErrInvalidImageFormat
			}

			return CreateProductOutput{}, fmt.Errorf("upload image: %w", err)
		}

		imageKey = key
	}

	prod, err := product.New(uuid.Generate(), input.StoreID, input.Name, imageKey, now)
	if err != nil {
		return CreateProductOutput{}, fmt.Errorf("new product: %w", err)
	}

	if err := uc.productRepo.Save(ctx, prod); err != nil {
		uc.storage.Delete(imageKey)
		return CreateProductOutput{}, fmt.Errorf("save product: %w", err)
	}

	return CreateProductOutput{
		ID: prod.ID(),
	}, nil
}
