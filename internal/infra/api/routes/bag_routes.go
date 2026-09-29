package routes

import (
	"github.com/Jeff1507/go-to-save-api/internal/core/bag/handler"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/middleware"
	"github.com/gin-gonic/gin"
)

type BagRoutes struct {
	auth          *middleware.Authentication
	list          *handler.ListBagsHandler
	get           *handler.GetBagHandler
	create        *handler.CreateBagHandler
	update        *handler.UpdateBagHandler
	updateStock   *handler.UpdateStockHandler
	close         *handler.CloseBagHandler
	getReservable *handler.GetReservableBagHandler
}

func NewBagRoutes(
	auth *middleware.Authentication,
	list *handler.ListBagsHandler,
	get *handler.GetBagHandler,
	create *handler.CreateBagHandler,
	update *handler.UpdateBagHandler,
	updateStock *handler.UpdateStockHandler,
	close *handler.CloseBagHandler,
	getReservable *handler.GetReservableBagHandler,
) *BagRoutes {
	return &BagRoutes{
		auth:          auth,
		list:          list,
		get:           get,
		create:        create,
		update:        update,
		updateStock:   updateStock,
		close:         close,
		getReservable: getReservable,
	}
}

func (r *BagRoutes) Register(router *gin.Engine) {
	bags := router.Group("/bags")
	bags.Use(handler.ErrorHandler)
	bags.Use(r.auth.AuthenticateCustomer)
	bags.GET("/:id", r.getReservable.Handle)

	store := router.Group("/store/bags")
	store.Use(handler.ErrorHandler)
	store.Use(r.auth.AuthenticateStore)
	store.POST("", r.create.Handle)
	store.GET("", r.list.Handle)
	store.GET("/:id", r.get.Handle)
	store.PUT("/:id", r.update.Handle)
	store.PATCH("/:id/stock", r.updateStock.Handle)
	store.PATCH("/:id/close", r.close.Handle)
}
