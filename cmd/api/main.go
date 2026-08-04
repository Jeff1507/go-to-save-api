package main

import (
	"fmt"
	"log"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	bagHandler "github.com/Jeff1507/go-to-save-api/internal/core/bag/handler"
	bagQuery "github.com/Jeff1507/go-to-save-api/internal/core/bag/query"
	bagRepository "github.com/Jeff1507/go-to-save-api/internal/core/bag/repository"
	bagUseCase "github.com/Jeff1507/go-to-save-api/internal/core/bag/usecase"
	customerHandler "github.com/Jeff1507/go-to-save-api/internal/core/customer/handler"
	customerQuery "github.com/Jeff1507/go-to-save-api/internal/core/customer/query"
	customerRepository "github.com/Jeff1507/go-to-save-api/internal/core/customer/repository"
	customerUseCase "github.com/Jeff1507/go-to-save-api/internal/core/customer/usecase"
	pickupHandler "github.com/Jeff1507/go-to-save-api/internal/core/pickup/handler"
	pickupQuery "github.com/Jeff1507/go-to-save-api/internal/core/pickup/query"
	pickupRepository "github.com/Jeff1507/go-to-save-api/internal/core/pickup/repository"
	pickupUseCase "github.com/Jeff1507/go-to-save-api/internal/core/pickup/usecase"
	productHandler "github.com/Jeff1507/go-to-save-api/internal/core/product/handler"
	productQuery "github.com/Jeff1507/go-to-save-api/internal/core/product/query"
	productRepository "github.com/Jeff1507/go-to-save-api/internal/core/product/repository"
	productUseCase "github.com/Jeff1507/go-to-save-api/internal/core/product/usecase"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/middleware"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/routes"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/server"
	"github.com/Jeff1507/go-to-save-api/internal/infra/auth/session"
	"github.com/Jeff1507/go-to-save-api/internal/infra/db"
	"github.com/Jeff1507/go-to-save-api/internal/infra/otp"
	"github.com/Jeff1507/go-to-save-api/internal/infra/storage"
	"github.com/Jeff1507/go-to-save-api/internal/infra/validation"
	"github.com/Jeff1507/go-to-save-api/internal/pkg/config"
)

func main() {
	if err := start(); err != nil {
		log.Fatalf("start api: %v", err)
	}
}

func start() error {
	config.LoadEnv()

	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("new config: %v", err)
	}

	sysClock := clock.SystemClock{}

	pg, err := db.NewPostgresPool(cfg.Postgres.URL)
	if err != nil {
		return fmt.Errorf("connect to postgres: %v", err)
	}
	defer pg.Close()

	rdb, err := db.NewRedisClient(
		cfg.Redis.Addr,
		cfg.Redis.Username,
		cfg.Redis.Password,
		cfg.Redis.DB,
	)
	if err != nil {
		return fmt.Errorf("connect to redis: %v", err)
	}
	defer rdb.Close()

	localStorage, err := storage.NewLocalStorage(cfg.Storage.BaseURL, cfg.Storage.BasePath)
	if err != nil {
		return fmt.Errorf("new local storage: %v", err)
	}

	validator := validation.NewValidator()

	// Auth
	sessionManager := session.NewRedisManager(rdb)
	otpManager := otp.NewRedisManager(rdb)

	// Repositories
	customerRepo := customerRepository.NewPostgresRepo(pg)
	pickupRepo := pickupRepository.NewPostgresRepo(pg)
	productRepo := productRepository.NewPostgresRepo(pg)
	bagRepo := bagRepository.NewPostgresRepo(pg)

	// Queries
	getCustomer := customerQuery.NewGetCustomer(pg, localStorage)

	listStoreProducts := productQuery.NewListStoreProducts(pg, localStorage)

	listStorePickupTemplates := pickupQuery.NewListStorePickupTemplates(pg)

	listStoreBags := bagQuery.NewListStoreBags(pg, localStorage)

	// Use Cases
	signupCustomer := customerUseCase.NewSignup(sysClock, customerRepo, sessionManager, otpManager)
	loginCustomer := customerUseCase.NewLogin(customerRepo, sessionManager)

	createPickupTemplate := pickupUseCase.NewCreatePickupTemplate(pickupRepo)
	updatePickupTemplate := pickupUseCase.NewUpdatePickupTemplate(pickupRepo)
	deletPickupTemplate := pickupUseCase.NewDeletePickupTemplate(pickupRepo)

	createProduct := productUseCase.NewCreateProduct(sysClock, localStorage, productRepo)
	updateProduct := productUseCase.NewUpdateProduct(sysClock, localStorage, productRepo)
	deleteProduct := productUseCase.NewDeleteProduct(productRepo)

	createBag := bagUseCase.NewCreateBag(sysClock, bagRepo, productRepo)
	updateBag := bagUseCase.NewUpdateBag(sysClock, bagRepo, productRepo)
	closeBag := bagUseCase.NewCloseBag(sysClock, bagRepo)

	// Handlers
	getCustomerHandler := customerHandler.NewGetCustomer(getCustomer)
	signupCustomerHandler := customerHandler.NewSignup(validator, signupCustomer)
	loginCustomerHandler := customerHandler.NewLogin(validator, loginCustomer)

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
	listStorePickupTemplatesHandler := pickupHandler.NewListStorePickupTemplates(
		listStorePickupTemplates,
	)

	createProductHandler := productHandler.NewCreateProduct(validator, createProduct)
	updateProductHandler := productHandler.NewUpdateProduct(validator, updateProduct)
	deleteProductHandler := productHandler.NewDeleteProduct(validator, deleteProduct)
	listStoreProductsHandler := productHandler.NewListStoreProducts(listStoreProducts)

	createBagHandler := bagHandler.NewCreateBag(validator, createBag)
	updateBagHandler := bagHandler.NewUpdateBag(validator, updateBag)
	closeBagHandler := bagHandler.NewCloseBag(validator, closeBag)
	listStoreBagsHandler := bagHandler.NewListStoreBags(listStoreBags)

	// Middlewares
	authentication := middleware.NewAuthentication(sessionManager)

	// Routes
	customerRoutes := routes.NewCustomerRoutes(
		authentication,
		getCustomerHandler,
		signupCustomerHandler,
		loginCustomerHandler,
	)

	pickupTemplateRoutes := routes.NewPickupTemplateRoutes(
		createPickupTemplateHandler,
		updatePickupTemplateHandler,
		deletePickupTemplateHandler,
		listStorePickupTemplatesHandler,
	)

	productRoutes := routes.NewProductRoutes(
		createProductHandler,
		updateProductHandler,
		deleteProductHandler,
		listStoreProductsHandler,
	)

	bagRoutes := routes.NewBagRoutes(
		createBagHandler,
		updateBagHandler,
		closeBagHandler,
		listStoreBagsHandler,
	)

	// Server
	srv := server.New([]string{cfg.App.WebOrigin}, cfg.Storage.BasePath)

	srv.AddRoutes(
		customerRoutes,
		pickupTemplateRoutes,
		productRoutes,
		bagRoutes,
	)

	if err := srv.Run(cfg.App.Port); err != nil {
		return fmt.Errorf("run server: %v", err)
	}

	return nil
}
