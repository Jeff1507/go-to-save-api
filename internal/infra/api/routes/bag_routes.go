package routes

import (
	"github.com/Jeff1507/go-to-save-api/internal/core/bag/handler"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/middleware"
	"github.com/gin-gonic/gin"
)

type BagRoutes struct {
	auth      *middleware.Authentication
	list      *handler.ListBagsHandler
	get       *handler.GetBagHandler
	create    *handler.CreateBagHandler
	update    *handler.UpdateBagHandler
	close     *handler.CloseBagHandler
	getActive *handler.GetActiveBagHandler
}

func NewBagRoutes(
	auth *middleware.Authentication,
	list *handler.ListBagsHandler,
	get *handler.GetBagHandler,
	create *handler.CreateBagHandler,
	update *handler.UpdateBagHandler,
	close *handler.CloseBagHandler,
	getActive *handler.GetActiveBagHandler,
) *BagRoutes {
	return &BagRoutes{
		auth:      auth,
		list:      list,
		get:       get,
		create:    create,
		update:    update,
		close:     close,
		getActive: getActive,
	}
}

func (r *BagRoutes) Register(router *gin.Engine) {
	bags := router.Group("/bags")
	bags.Use(handler.ErrorHandler)

	{
		// Rotas para o app mobile
		reservation := bags.Group("/reservation")
		reservation.Use(r.auth.AuthenticateCustomer)

		reservation.GET("/:id", r.getActive.Handle)
	}

	{
		// Rotas para o dashboard web
		store := bags.Group("")
		store.Use(r.auth.AuthenticateStore)

		store.POST("", r.create.Handle)
		store.GET("", r.list.Handle)
		store.GET("/:id", r.get.Handle)
		store.PUT("/:id", r.update.Handle)
		store.PATCH("/:id/close", r.close.Handle)
	}
}
