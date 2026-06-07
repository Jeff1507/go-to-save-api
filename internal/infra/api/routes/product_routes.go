package routes

import (
	"github.com/Jeff1507/go-to-save-api/internal/core/product/handler"
	"github.com/go-chi/chi/v5"
)

type ProductRoutes struct {
	create      *handler.CreateProduct
	update      *handler.UpdateProduct
	delete      *handler.DeleteProduct
	listByStore *handler.ListProductsByStore
}

func NewProductRoutes(
	create *handler.CreateProduct,
	update *handler.UpdateProduct,
	delete *handler.DeleteProduct,
	listByStore *handler.ListProductsByStore,
) *ProductRoutes {
	return &ProductRoutes{
		create:      create,
		update:      update,
		delete:      delete,
		listByStore: listByStore,
	}
}

func (r *ProductRoutes) Register(router *chi.Mux) {
	router.Route("/products", func(products chi.Router) {
		products.Post("/", r.create.Handle)
		products.Put("/{id}", r.update.Handle)
		products.Delete("/{id}", r.delete.Handle)
		products.Get("/", r.listByStore.Handle)
	})
}
