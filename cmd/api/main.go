package main

import (
	"fmt"
	"log"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	pickupHandler "github.com/Jeff1507/go-to-save-api/internal/core/pickup/handler"
	pickupRepository "github.com/Jeff1507/go-to-save-api/internal/core/pickup/repository"
	pickupUseCase "github.com/Jeff1507/go-to-save-api/internal/core/pickup/usecase"
	productHandler "github.com/Jeff1507/go-to-save-api/internal/core/product/handler"
	productRepository "github.com/Jeff1507/go-to-save-api/internal/core/product/repository"
	productUseCase "github.com/Jeff1507/go-to-save-api/internal/core/product/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/routes"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/server"
	"github.com/Jeff1507/go-to-save-api/internal/infra/storage"
)

func main() {
	if err := start(); err != nil {
		log.Fatalf("start api: %v", err)
	}
}

func start() error {
	sysClock := clock.SystemClock{}

	localStorage, err := storage.NewLocalStorage("http://localhost:8080", "./uploads/")
	if err != nil {
		return fmt.Errorf("new local storage: %v", err)
	}

	pickupRepo := pickupRepository.NewInMemRepo()
	productRepo := productRepository.NewInMemRepo()

	createPickupTemplate := pickupUseCase.NewCreatePickupTemplate(pickupRepo)
	updatePickupTemplate := pickupUseCase.NewUpdatePickupTemplate(pickupRepo)
	deletPickupTemplate := pickupUseCase.NewDeletePickupTemplate(pickupRepo)

	createProduct := productUseCase.NewCreateProduct(sysClock, localStorage, productRepo)
	updateProduct := productUseCase.NewUpdateProduct(sysClock, localStorage, productRepo)
	deleteProduct := productUseCase.NewDeleteProduct(localStorage, productRepo)

	createPickupTemplateHandler := pickupHandler.NewCreatePickupTemplate(createPickupTemplate)
	updatePickupTemplateHandler := pickupHandler.NewUpdatePickupTemplate(updatePickupTemplate)
	deletePickupTemplateHandler := pickupHandler.NewDeletePickupTemplate(deletPickupTemplate)

	createProductHandler := productHandler.NewCreateProduct(createProduct)
	updateProductHandler := productHandler.NewUpdateProduct(updateProduct)
	deleteProductHandler := productHandler.NewDeleteProduct(deleteProduct)

	pickupTemplateRoutes := routes.NewPickupTemplateRoutes(
		createPickupTemplateHandler,
		updatePickupTemplateHandler,
		deletePickupTemplateHandler,
	)

	productRoutes := routes.NewProductRoutes(
		createProductHandler,
		updateProductHandler,
		deleteProductHandler,
	)

	srv := server.New()

	srv.AddRoutes(
		pickupTemplateRoutes,
		productRoutes,
	)

	if err := srv.Run("8080"); err != nil {
		return fmt.Errorf("run server: %v", err)
	}

	return nil
}
