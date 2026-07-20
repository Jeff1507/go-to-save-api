package routes

import (
	"github.com/Jeff1507/go-to-save-api/internal/core/product/handler"
	"github.com/go-chi/chi/v5"
)

type ProductRoutes struct {
	create *handler.CreateProduct
	list   *handler.ListStoreProducts
	update *handler.UpdateProduct
	delete *handler.DeleteProduct
}

func NewProductRoutes(
	create *handler.CreateProduct,
	update *handler.UpdateProduct,
	delete *handler.DeleteProduct,
	list *handler.ListStoreProducts,
) *ProductRoutes {
	return &ProductRoutes{
		create: create,
		list:   list,
		update: update,
		delete: delete,
	}
}

func (r *ProductRoutes) Register(router *chi.Mux) {
	router.Route("/products", func(products chi.Router) {
		products.Post("/", r.create.Handle)
		products.Get("/", r.list.Handle)
		products.Put("/{id}", r.update.Handle)
		products.Delete("/{id}", r.delete.Handle)
	})
}
