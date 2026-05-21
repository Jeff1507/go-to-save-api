package usecase

import (
	"context"
	"fmt"
	"io"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/product"
	"github.com/Jeff1507/go-to-save-api/internal/infra/storage"
	"github.com/Jeff1507/go-to-save-api/internal/pkg/uuid"
)

type CreateProductImageInput struct {
	Content     io.Reader
	ContentType string
}

type CreateProductInput struct {
	StoreID string
	Name    string
	Image   *CreateProductImageInput
}

type CreateProductOutput struct {
	ID string
}

type CreateProduct struct {
	clock       clock.Clock
	fileStorage storage.Storage
	productRepo product.Repository
}

func NewCreateProduct(
	clock clock.Clock,
	fileStorage storage.Storage,
	productRepo product.Repository,
) *CreateProduct {
	return &CreateProduct{
		clock:       clock,
		fileStorage: fileStorage,
		productRepo: productRepo,
	}
}

func (uc *CreateProduct) Execute(
	ctx context.Context,
	input CreateProductInput,
) (CreateProductOutput, error) {
	now := uc.clock.Now()

	exists, err := uc.productRepo.Exists(ctx, input.StoreID, input.Name)
	if err != nil {
		return CreateProductOutput{}, fmt.Errorf("check if product already exists: %v", err)
	}

	if exists {
		return CreateProductOutput{}, product.ErrAlreadyExists
	}

	var imgKey string
	if input.Image != nil {
		imgKey, err = uc.fileStorage.Upload(
			input.Image.Content,
			input.Image.ContentType,
			storage.ProductsPath,
		)
		if err != nil {
			return CreateProductOutput{}, fmt.Errorf("upload image: %v", err)
		}
	}

	prod, err := product.New(uuid.Generate(), input.StoreID, input.Name, imgKey, now)
	if err != nil {
		return CreateProductOutput{}, fmt.Errorf("new product: %v", err)
	}

	if err := uc.productRepo.Save(ctx, prod); err != nil {
		uc.fileStorage.Delete(imgKey)
		return CreateProductOutput{}, fmt.Errorf("save product: %v", err)
	}

	return CreateProductOutput{
		ID: prod.ID(),
	}, nil
}
