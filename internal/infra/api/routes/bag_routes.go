package routes

import (
	"github.com/Jeff1507/go-to-save-api/internal/core/bag/handler"
	"github.com/go-chi/chi/v5"
)

type BagRoutes struct {
	create *handler.CreateBag
	list   *handler.ListStoreBags
	update *handler.UpdateBag
	close  *handler.CloseBag
}

func NewBagRoutes(
	create *handler.CreateBag,
	update *handler.UpdateBag,
	close *handler.CloseBag,
	list *handler.ListStoreBags,
) *BagRoutes {
	return &BagRoutes{
		create: create,
		list:   list,
		update: update,
		close:  close,
	}
}

func (r *BagRoutes) Register(router *chi.Mux) {
	router.Route("/bags", func(bags chi.Router) {
		bags.Post("/", r.create.Handle)
		bags.Get("/", r.list.Handle)
		bags.Put("/{id}", r.update.Handle)
		bags.Patch("/{id}/close", r.close.Handle)
	})
}
