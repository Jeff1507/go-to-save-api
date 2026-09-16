package routes

import (
	"github.com/Jeff1507/go-to-save-api/internal/core/order/handler"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/middleware"
	"github.com/gin-gonic/gin"
)

type OrderRoutes struct {
	auth             *middleware.Authentication
	list             *handler.ListOrdersHandler
	get              *handler.GetOrderHandler
	listSales        *handler.ListSalesHandler
	getSale          *handler.GetSaleHandler
	create           *handler.PlaceOrderHandler
	complete         *handler.CompleteOrderHandler
	cancelByStore    *handler.CancelOrderByStoreHandler
	cancelByCustomer *handler.CancelOrderByCustomerHandler
}

func NewOrderRoutes(
	auth *middleware.Authentication,
	list *handler.ListOrdersHandler,
	get *handler.GetOrderHandler,
	listSales *handler.ListSalesHandler,
	getSale *handler.GetSaleHandler,
	create *handler.PlaceOrderHandler,
	complete *handler.CompleteOrderHandler,
	cancelByStore *handler.CancelOrderByStoreHandler,
	cancelByCustomer *handler.CancelOrderByCustomerHandler,
) *OrderRoutes {
	return &OrderRoutes{
		auth:             auth,
		list:             list,
		get:              get,
		listSales:        listSales,
		getSale:          getSale,
		create:           create,
		complete:         complete,
		cancelByStore:    cancelByStore,
		cancelByCustomer: cancelByCustomer,
	}
}

func (r *OrderRoutes) Register(router *gin.Engine) {
	orders := router.Group("/orders")
	orders.Use(handler.ErrorHandler)
	orders.Use(r.auth.AuthenticateCustomer)
	orders.POST("", r.create.Handle)
	orders.GET("", r.list.Handle)
	orders.GET("/:id", r.get.Handle)
	orders.PATCH("/:id/cancel", r.cancelByCustomer.Handle)

	store := router.Group("/store/orders")
	store.Use(handler.ErrorHandler)
	store.Use(r.auth.AuthenticateStore)
	store.GET("", r.listSales.Handle)
	store.GET("/:id", r.getSale.Handle)
	store.PATCH("/:id/complete", r.complete.Handle)
	store.PATCH("/:id/cancel", r.cancelByStore.Handle)
}
