package main

import (
	"fmt"
	"log"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	bagHandler "github.com/Jeff1507/go-to-save-api/internal/core/bag/handler"
	bagRepository "github.com/Jeff1507/go-to-save-api/internal/core/bag/repository"
	bagUseCase "github.com/Jeff1507/go-to-save-api/internal/core/bag/usecase"
	pickupHandler "github.com/Jeff1507/go-to-save-api/internal/core/pickup/handler"
	pickupRepository "github.com/Jeff1507/go-to-save-api/internal/core/pickup/repository"
	pickupUseCase "github.com/Jeff1507/go-to-save-api/internal/core/pickup/usecase"
	productHandler "github.com/Jeff1507/go-to-save-api/internal/core/product/handler"
	productRepository "github.com/Jeff1507/go-to-save-api/internal/core/product/repository"
	productUseCase "github.com/Jeff1507/go-to-save-api/internal/core/product/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/routes"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/server"
	"github.com/Jeff1507/go-to-save-api/internal/infra/storage"
	"github.com/Jeff1507/go-to-save-api/internal/infra/validation"
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

	// Repositories
	pickupRepo := pickupRepository.NewInMemRepo()
	productRepo := productRepository.NewInMemRepo()
	bagRepo := bagRepository.NewInMemRepo()

	// Use Cases
	createPickupTemplate := pickupUseCase.NewCreatePickupTemplate(pickupRepo)
	updatePickupTemplate := pickupUseCase.NewUpdatePickupTemplate(pickupRepo)
	deletPickupTemplate := pickupUseCase.NewDeletePickupTemplate(pickupRepo)

	createProduct := productUseCase.NewCreateProduct(sysClock, localStorage, productRepo)
	updateProduct := productUseCase.NewUpdateProduct(sysClock, localStorage, productRepo)
	deleteProduct := productUseCase.NewDeleteProduct(localStorage, productRepo)

	createBag := bagUseCase.NewCreateBag(sysClock, bagRepo, productRepo)
	updateBag := bagUseCase.NewUpdateBag(sysClock, bagRepo, productRepo)
	closeBag := bagUseCase.NewCloseBag(sysClock, bagRepo)

	validator := validation.NewValidator()

	// Handlers
	createPickupTemplateHandler := pickupHandler.NewCreatePickupTemplate(
		validator,
		createPickupTemplate,
	)
	updatePickupTemplateHandler := pickupHandler.NewUpdatePickupTemplate(
		validator,
		updatePickupTemplate,
	)
	deletePickupTemplateHandler := pickupHandler.NewDeletePickupTemplate(
		validator,
		deletPickupTemplate,
	)

	createProductHandler := productHandler.NewCreateProduct(validator, createProduct)
	updateProductHandler := productHandler.NewUpdateProduct(validator, updateProduct)
	deleteProductHandler := productHandler.NewDeleteProduct(validator, deleteProduct)

	createBagHandler := bagHandler.NewCreateBag(validator, createBag)
	updateBagHandler := bagHandler.NewUpdateBag(validator, updateBag)
	closeBagHandler := bagHandler.NewCloseBag(validator, closeBag)

	// Routes
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

	bagRoutes := routes.NewBagRoutes(
		createBagHandler,
		updateBagHandler,
		closeBagHandler,
	)

	// Server
	srv := server.New([]string{"http://localhost:5173"})

	srv.AddRoutes(
		pickupTemplateRoutes,
		productRoutes,
		bagRoutes,
	)

	if err := srv.Run("8080"); err != nil {
		return fmt.Errorf("run server: %v", err)
	}

	return nil
}
