package routes

import (
	"github.com/Jeff1507/go-to-save-api/internal/core/bag/handler"
	"github.com/go-chi/chi/v5"
)

type BagRoutes struct {
	createBag *handler.CreateBag
	updateBag *handler.UpdateBag
	closeBag  *handler.CloseBag
}

func NewBagRoutes(
	createBag *handler.CreateBag,
	updateBag *handler.UpdateBag,
	closeBag *handler.CloseBag,
) *BagRoutes {
	return &BagRoutes{
		createBag: createBag,
		updateBag: updateBag,
		closeBag:  closeBag,
	}
}

func (r *BagRoutes) Register(router *chi.Mux) {
	router.Route("/bags", func(bags chi.Router) {
		bags.Post("/", r.createBag.Handle)
		bags.Put("/{id}", r.updateBag.Handle)
		bags.Patch("/{id}/close", r.closeBag.Handle)
	})
}
