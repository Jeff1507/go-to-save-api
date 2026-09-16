package routes

import (
	"github.com/Jeff1507/go-to-save-api/internal/core/product/handler"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/middleware"
	"github.com/gin-gonic/gin"
)

type ProductRoutes struct {
	auth   *middleware.Authentication
	list   *handler.ListProductsHandler
	create *handler.CreateProductHandler
	update *handler.UpdateProductHandler
	delete *handler.DeleteProductHandler
}

func NewProductRoutes(
	auth *middleware.Authentication,
	list *handler.ListProductsHandler,
	create *handler.CreateProductHandler,
	update *handler.UpdateProductHandler,
	delete *handler.DeleteProductHandler,
) *ProductRoutes {
	return &ProductRoutes{
		auth:   auth,
		list:   list,
		create: create,
		update: update,
		delete: delete,
	}
}

func (r *ProductRoutes) Register(router *gin.Engine) {
	store := router.Group("/store/products")
	store.Use(handler.ErrorHandler)
	store.Use(middleware.MaxRequestBodySize(5 << 20))
	store.Use(r.auth.AuthenticateStore)
	store.POST("", r.create.Handle)
	store.GET("", r.list.Handle)
	store.PUT("/:id", r.update.Handle)
	store.DELETE("/:id", r.delete.Handle)
}
