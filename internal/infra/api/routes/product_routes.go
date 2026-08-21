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
	products := router.Group("/products")
	products.Use(handler.ErrorHandler)
	products.Use(middleware.MaxRequestBodySize(5 << 20))

	{
		products.Use(r.auth.VerifyStore)

		products.POST("", r.create.Handle)
		products.GET("", r.list.Handle)
		products.PUT("/:id", r.update.Handle)
		products.DELETE("/:id", r.delete.Handle)
	}
}
