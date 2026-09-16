package main

import (
	"fmt"
	"log"

	customerHandler "github.com/Jeff1507/go-to-save-api/internal/core/customer/handler"
	customerQuery "github.com/Jeff1507/go-to-save-api/internal/core/customer/query"
	customerRepository "github.com/Jeff1507/go-to-save-api/internal/core/customer/repository"
	customerUseCase "github.com/Jeff1507/go-to-save-api/internal/core/customer/usecase"

	storeHandler "github.com/Jeff1507/go-to-save-api/internal/core/store/handler"
	storeQuery "github.com/Jeff1507/go-to-save-api/internal/core/store/query"
	storeRepository "github.com/Jeff1507/go-to-save-api/internal/core/store/repository"
	storeUseCase "github.com/Jeff1507/go-to-save-api/internal/core/store/usecase"

	productHandler "github.com/Jeff1507/go-to-save-api/internal/core/product/handler"
	productQuery "github.com/Jeff1507/go-to-save-api/internal/core/product/query"
	productRepository "github.com/Jeff1507/go-to-save-api/internal/core/product/repository"
	productUseCase "github.com/Jeff1507/go-to-save-api/internal/core/product/usecase"

	presetHandler "github.com/Jeff1507/go-to-save-api/internal/core/preset/handler"
	presetQuery "github.com/Jeff1507/go-to-save-api/internal/core/preset/query"
	presetRepository "github.com/Jeff1507/go-to-save-api/internal/core/preset/repository"
	presetUseCase "github.com/Jeff1507/go-to-save-api/internal/core/preset/usecase"

	bagHandler "github.com/Jeff1507/go-to-save-api/internal/core/bag/handler"
	bagQuery "github.com/Jeff1507/go-to-save-api/internal/core/bag/query"
	bagRepository "github.com/Jeff1507/go-to-save-api/internal/core/bag/repository"
	bagUseCase "github.com/Jeff1507/go-to-save-api/internal/core/bag/usecase"

	orderHandler "github.com/Jeff1507/go-to-save-api/internal/core/order/handler"
	orderQuery "github.com/Jeff1507/go-to-save-api/internal/core/order/query"
	orderRepository "github.com/Jeff1507/go-to-save-api/internal/core/order/repository"
	orderUseCase "github.com/Jeff1507/go-to-save-api/internal/core/order/usecase"

	reviewHandler "github.com/Jeff1507/go-to-save-api/internal/core/review/handler"
	reviewRepository "github.com/Jeff1507/go-to-save-api/internal/core/review/repository"
	reviewUseCase "github.com/Jeff1507/go-to-save-api/internal/core/review/usecase"

	"github.com/Jeff1507/go-to-save-api/internal/common/clock"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/middleware"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/routes"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/server"
	"github.com/Jeff1507/go-to-save-api/internal/infra/auth/otp"
	"github.com/Jeff1507/go-to-save-api/internal/infra/auth/session"
	"github.com/Jeff1507/go-to-save-api/internal/infra/db"
	"github.com/Jeff1507/go-to-save-api/internal/infra/storage"
	"github.com/Jeff1507/go-to-save-api/internal/infra/transaction"
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
		return fmt.Errorf("new config: %w", err)
	}

	sysClock := clock.SystemClock{}

	pg, err := db.NewPostgresPool(cfg.Postgres.URL)
	if err != nil {
		return fmt.Errorf("connect to postgres: %w", err)
	}
	defer pg.Close()

	rdb, err := db.NewRedisClient(
		cfg.Redis.Addr,
		cfg.Redis.Username,
		cfg.Redis.Password,
		cfg.Redis.DB,
	)
	if err != nil {
		return fmt.Errorf("connect to redis: %w", err)
	}
	defer rdb.Close()

	fileStorage, err := storage.NewLocalStorage(cfg.Storage.BasePath)
	if err != nil {
		return fmt.Errorf("new local storage: %w", err)
	}

	// Auth
	sessionManager := session.NewRedisManager(rdb)
	otpManager := otp.NewRedisManager(rdb)

	// Repositories
	txManager := transaction.NewPostgresManager(pg)

	customerRepo := customerRepository.NewPostgresRepo(pg)
	storeRepo := storeRepository.NewPostgresRepo(pg)
	productRepo := productRepository.NewPostgresRepo(pg)
	presetRepo := presetRepository.NewPostgresRepo(pg)
	bagRepo := bagRepository.NewPostgresRepo(pg)
	orderRepo := orderRepository.NewPostgresRepo(pg)
	reviewRepo := reviewRepository.NewPostgresRepo(pg)

	// Use Cases
	signupCustomer := customerUseCase.NewSignup(customerRepo, sessionManager, otpManager, sysClock)
	loginCustomer := customerUseCase.NewLogin(customerRepo, sessionManager)
	logoutCustomer := customerUseCase.NewLogout(sessionManager)
	logoutAllCustomer := customerUseCase.NewLogoutAll(sessionManager)
	updateCustomerProfile := customerUseCase.NewUpdateProfile(customerRepo, sysClock)
	requestCustomerEmailVerification := customerUseCase.NewRequestEmailVerification(customerRepo, otpManager)
	verifyCustomerEmail := customerUseCase.NewVerifyEmail(customerRepo, otpManager, sysClock)
	requestCustomerPasswordReset := customerUseCase.NewRequestPasswordReset(customerRepo, otpManager)
	resetCustomerPassword := customerUseCase.NewResetPassword(customerRepo, sessionManager, otpManager, sysClock)
	updateCustomerAvatar := customerUseCase.NewUpdateAvatar(customerRepo, fileStorage, sysClock)
	removeCustomerAvatar := customerUseCase.NewRemoveAvatar(customerRepo, fileStorage, sysClock)
	updateCustomerEmail := customerUseCase.NewUpdateEmail(customerRepo, otpManager, sysClock)
	updateCustomerPassword := customerUseCase.NewUpdatePassword(customerRepo, sessionManager, sysClock)

	signupStore := storeUseCase.NewSignup(storeRepo, sessionManager, otpManager, sysClock)
	loginStore := storeUseCase.NewLogin(storeRepo, sessionManager)
	logoutStore := storeUseCase.NewLogout(sessionManager)
	logoutAllStore := storeUseCase.NewLogoutAll(sessionManager)
	updateStoreProfile := storeUseCase.NewUpdateProfile(storeRepo, sysClock)
	requestStoreEmailVerification := storeUseCase.NewRequestEmailVerification(storeRepo, otpManager)
	verifyStoreEmail := storeUseCase.NewVerifyEmail(storeRepo, otpManager, sysClock)
	requestStorePasswordReset := storeUseCase.NewRequestPasswordReset(storeRepo, otpManager)
	resetStorePassword := storeUseCase.NewResetPassword(storeRepo, sessionManager, otpManager, sysClock)
	updateStoreAvatar := storeUseCase.NewUpdateAvatar(storeRepo, fileStorage, sysClock)
	removeStoreAvatar := storeUseCase.NewRemoveAvatar(storeRepo, fileStorage, sysClock)
	updateStoreEmail := storeUseCase.NewUpdateEmail(storeRepo, otpManager, sysClock)
	updateStorePassword := storeUseCase.NewUpdatePassword(storeRepo, sessionManager, sysClock)
	updateStoreAddress := storeUseCase.NewUpdateAddress(storeRepo, sysClock)

	createProduct := productUseCase.NewCreateProduct(productRepo, fileStorage, sysClock)
	updateProduct := productUseCase.NewUpdateProduct(productRepo, fileStorage, sysClock)
	deleteProduct := productUseCase.NewDeleteProduct(productRepo)

	createPreset := presetUseCase.NewCreatePreset(presetRepo)
	updatePreset := presetUseCase.NewUpdatePreset(presetRepo)
	deletePreset := presetUseCase.NewDeletePreset(presetRepo)

	createBag := bagUseCase.NewCreateBag(bagRepo, productRepo, sysClock)
	updateBag := bagUseCase.NewUpdateBag(bagRepo, productRepo, orderRepo, sysClock)
	updateBagStock := bagUseCase.NewUpdateStock(bagRepo, sysClock)
	closeBag := bagUseCase.NewCloseBag(bagRepo, sysClock)

	placeOrder := orderUseCase.NewPlaceOrder(txManager, orderRepo, bagRepo, sysClock)
	completeOrder := orderUseCase.NewCompleteOrder(orderRepo, sysClock)
	cancelOrderByStore := orderUseCase.NewCancelOrderByStore(txManager, orderRepo, bagRepo, sysClock)
	cancelOrderByCustomer := orderUseCase.NewCancelOrderByCustomer(txManager, orderRepo, bagRepo, sysClock)

	addReview := reviewUseCase.NewAddReview(txManager, reviewRepo, orderRepo, storeRepo, sysClock)

	// Queries
	getCustomerProfile := customerQuery.NewGetProfile(pg, fileStorage)

	getStoreProfile := storeQuery.NewGetProfile(pg, fileStorage)
	listNearbyStores := storeQuery.NewListNearbyStores(pg, fileStorage, sysClock)
	getStoreWithActiveBags := storeQuery.NewGetStoreWithActiveBags(pg, fileStorage, sysClock)

	listProducts := productQuery.NewListProducts(pg, fileStorage)

	listPresets := presetQuery.NewListPresets(pg)

	listBags := bagQuery.NewListBags(pg)
	getBag := bagQuery.NewGetBag(pg, fileStorage)
	getActiveBag := bagQuery.NewGetActiveBag(pg, fileStorage)

	listOrders := orderQuery.NewListOrders(pg, fileStorage)
	getOrder := orderQuery.NewGetOrder(pg, fileStorage)
	listSales := orderQuery.NewListSales(pg, fileStorage)
	getSale := orderQuery.NewGetSale(pg, fileStorage)

	// Middlewares
	authentication := middleware.NewAuthentication(sessionManager)

	// Routes
	customerRoutes := routes.NewCustomerRoutes(
		authentication,
		customerHandler.NewGetProfileHandler(getCustomerProfile),
		customerHandler.NewSignupHandler(signupCustomer),
		customerHandler.NewLoginHandler(loginCustomer),
		customerHandler.NewRequestEmailVerificationHandler(requestCustomerEmailVerification),
		customerHandler.NewVerifyEmailHandler(verifyCustomerEmail),
		customerHandler.NewRequestPasswordResetHandler(requestCustomerPasswordReset),
		customerHandler.NewResetPasswordHandler(resetCustomerPassword),
		customerHandler.NewLogoutHandler(logoutCustomer),
		customerHandler.NewLogoutAllHandler(logoutAllCustomer),
		customerHandler.NewUpdateProfileHandler(updateCustomerProfile),
		customerHandler.NewUpdateAvatarHandler(updateCustomerAvatar),
		customerHandler.NewRemoveAvatarHandler(removeCustomerAvatar),
		customerHandler.NewUpdateEmailHandler(updateCustomerEmail),
		customerHandler.NewUpdatePasswordHandler(updateCustomerPassword),
	)

	storeRoutes := routes.NewStoreRoutes(
		authentication,
		storeHandler.NewGetProfileHandler(getStoreProfile),
		storeHandler.NewSignupHandler(signupStore),
		storeHandler.NewLoginHandler(loginStore),
		storeHandler.NewRequestEmailVerificationHandler(requestStoreEmailVerification),
		storeHandler.NewVerifyEmailHandler(verifyStoreEmail),
		storeHandler.NewRequestPasswordResetHandler(requestStorePasswordReset),
		storeHandler.NewResetPasswordHandler(resetStorePassword),
		storeHandler.NewLogoutHandler(logoutStore),
		storeHandler.NewLogoutAllHandler(logoutAllStore),
		storeHandler.NewUpdateProfileHandler(updateStoreProfile),
		storeHandler.NewUpdateAvatarHandler(updateStoreAvatar),
		storeHandler.NewRemoveAvatarHandler(removeStoreAvatar),
		storeHandler.NewUpdateEmailHandler(updateStoreEmail),
		storeHandler.NewUpdatePasswordHandler(updateStorePassword),
		storeHandler.NewUpdateAddressHandler(updateStoreAddress),
		storeHandler.NewListNearbyStoresHandler(listNearbyStores),
		storeHandler.NewGetStoreWithActiveBagsHandler(getStoreWithActiveBags),
	)

	productRoutes := routes.NewProductRoutes(
		authentication,
		productHandler.NewListProductsHandler(listProducts),
		productHandler.NewCreateProductHandler(createProduct),
		productHandler.NewUpdateProductHandler(updateProduct),
		productHandler.NewDeleteProduct(deleteProduct),
	)

	presetRoutes := routes.NewPresetRoutes(
		authentication,
		presetHandler.NewListPresetsHandler(listPresets),
		presetHandler.NewCreatePresetHandler(createPreset),
		presetHandler.NewUpdatePresetHandler(updatePreset),
		presetHandler.NewDeletePresetHandler(deletePreset),
	)

	bagRoutes := routes.NewBagRoutes(
		authentication,
		bagHandler.NewListBagsHandler(listBags),
		bagHandler.NewGetBagHandler(getBag),
		bagHandler.NewCreateBagHandler(createBag),
		bagHandler.NewUpdateBagHandler(updateBag),
		bagHandler.NewUpdateStockHandler(updateBagStock),
		bagHandler.NewCloseBagHandler(closeBag),
		bagHandler.NewGetActiveBagHandler(getActiveBag),
	)

	orderRoutes := routes.NewOrderRoutes(
		authentication,
		orderHandler.NewListOrdersHandler(listOrders),
		orderHandler.NewGetOrderHandler(getOrder),
		orderHandler.NewListSalesHandler(listSales),
		orderHandler.NewGetSaleHandler(getSale),
		orderHandler.NewPlaceOrderHandler(placeOrder),
		orderHandler.NewCompleteOrderHandler(completeOrder),
		orderHandler.NewCancelOrderByStoreHandler(cancelOrderByStore),
		orderHandler.NewCancelOrderByCustomerHandler(cancelOrderByCustomer),
	)

	reviewRoutes := routes.NewReviewRoutes(
		authentication,
		reviewHandler.NewAddReviewHandler(addReview),
	)

	// Server
	srv := server.New([]string{cfg.App.WebOrigin}, cfg.Storage.BasePath)

	srv.AddRoutes(
		customerRoutes,
		storeRoutes,
		productRoutes,
		presetRoutes,
		bagRoutes,
		orderRoutes,
		reviewRoutes,
	)

	if err := srv.Run(cfg.App.Port); err != nil {
		return fmt.Errorf("run server: %w", err)
	}

	return nil
}
