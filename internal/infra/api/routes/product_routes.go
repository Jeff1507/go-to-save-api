package routes

import (
	"github.com/Jeff1507/go-to-save-api/internal/core/product/handler"
	"github.com/go-chi/chi/v5"
)

type ProductRoutes struct {
	createProduct *handler.CreateProduct
	updateProduct *handler.UpdateProduct
	deleteProduct *handler.DeleteProduct
}

func NewProductRoutes(
	createProduct *handler.CreateProduct,
	updateProduct *handler.UpdateProduct,
	deleteProduct *handler.DeleteProduct,
) *ProductRoutes {
	return &ProductRoutes{
		createProduct: createProduct,
		updateProduct: updateProduct,
		deleteProduct: deleteProduct,
	}
}

func (r *ProductRoutes) Register(router *chi.Mux) {
	router.Route("/products", func(products chi.Router) {
		products.Post("/", r.createProduct.Handle)
		products.Put("/{id}", r.updateProduct.Handle)
		products.Delete("/{id}", r.deleteProduct.Handle)
	})
}
