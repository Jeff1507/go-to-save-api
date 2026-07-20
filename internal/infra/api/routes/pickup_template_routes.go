package routes

import (
	"github.com/Jeff1507/go-to-save-api/internal/core/pickup/handler"
	"github.com/go-chi/chi/v5"
)

type PickupTemplateRoutes struct {
	create *handler.CreatePickupTemplate
	list   *handler.ListStorePickupTemplates
	update *handler.UpdatePickupTemplate
	delete *handler.DeletePickupTemplate
}

func NewPickupTemplateRoutes(
	create *handler.CreatePickupTemplate,
	update *handler.UpdatePickupTemplate,
	delete *handler.DeletePickupTemplate,
	list *handler.ListStorePickupTemplates,
) *PickupTemplateRoutes {
	return &PickupTemplateRoutes{
		create: create,
		list:   list,
		update: update,
		delete: delete,
	}
}

func (r *PickupTemplateRoutes) Register(router *chi.Mux) {
	router.Route("/pickup-templates", func(pickupTemplates chi.Router) {
		pickupTemplates.Post("/", r.create.Handle)
		pickupTemplates.Get("/", r.list.Handle)
		pickupTemplates.Put("/{id}", r.update.Handle)
		pickupTemplates.Delete("/{id}", r.delete.Handle)
	})
}
