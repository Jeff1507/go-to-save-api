package main

import (
	"fmt"
	"log"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/core/product/handler"
	"github.com/Jeff1507/go-to-save-api/internal/core/product/repository"
	"github.com/Jeff1507/go-to-save-api/internal/core/product/usecase"
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

	productRepo := repository.NewInMemRepo()

	createProduct := usecase.NewCreateProduct(sysClock, localStorage, productRepo)
	updateProduct := usecase.NewUpdateProduct(sysClock, localStorage, productRepo)
	deleteProduct := usecase.NewDeleteProduct(localStorage, productRepo)

	createProductHandler := handler.NewCreateProduct(createProduct)
	updateProductHandler := handler.NewUpdateProduct(updateProduct)
	deleteProductHandler := handler.NewDeleteProduct(deleteProduct)

	productRoutes := routes.NewProductRoutes(
		createProductHandler,
		updateProductHandler,
		deleteProductHandler,
	)

	srv := server.New()

	srv.AddRoutes(productRoutes)

	if err := srv.Run("8080"); err != nil {
		return fmt.Errorf("run server: %v", err)
	}

	return nil
}
