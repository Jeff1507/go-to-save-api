package routes

import (
	"github.com/Jeff1507/go-to-save-api/internal/core/pickup/handler"
	"github.com/go-chi/chi/v5"
)

type PickupTemplateRoutes struct {
	createPickupTemplate *handler.CreatePickupTemplate
	updatePickupTemplate *handler.UpdatePickupTemplate
	deletePickupTemplate *handler.DeletePickupTemplate
}

func NewPickupTemplateRoutes(
	createPickupTemplate *handler.CreatePickupTemplate,
	updatePickupTemplate *handler.UpdatePickupTemplate,
	deletePickupTemplate *handler.DeletePickupTemplate,
) *PickupTemplateRoutes {
	return &PickupTemplateRoutes{
		createPickupTemplate: createPickupTemplate,
		updatePickupTemplate: updatePickupTemplate,
		deletePickupTemplate: deletePickupTemplate,
	}
}

func (r *PickupTemplateRoutes) Register(router *chi.Mux) {
	router.Route("/pickup-templates", func(pickupTemplates chi.Router) {
		pickupTemplates.Post("/", r.createPickupTemplate.Handle)
		pickupTemplates.Put("/{id}", r.updatePickupTemplate.Handle)
		pickupTemplates.Delete("/{id}", r.deletePickupTemplate.Handle)
	})
}
