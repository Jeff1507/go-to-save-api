package routes

import (
	"github.com/Jeff1507/go-to-save-api/internal/core/pickup/handler"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/middleware"
	"github.com/gin-gonic/gin"
)

type PickupTemplateRoutes struct {
	auth   *middleware.Authentication
	list   *handler.ListTemplatesHandler
	create *handler.CreateTemplateHandler
	update *handler.UpdateTemplateHandler
	delete *handler.DeleteTemplateHandler
}

func NewPickupTemplateRoutes(
	auth *middleware.Authentication,
	list *handler.ListTemplatesHandler,
	create *handler.CreateTemplateHandler,
	update *handler.UpdateTemplateHandler,
	delete *handler.DeleteTemplateHandler,
) *PickupTemplateRoutes {
	return &PickupTemplateRoutes{
		auth:   auth,
		list:   list,
		create: create,
		update: update,
		delete: delete,
	}
}

func (r *PickupTemplateRoutes) Register(router *gin.Engine) {
	pickupTemplates := router.Group("/pickup-templates")
	pickupTemplates.Use(handler.ErrorHandler)

	{
		pickupTemplates.Use(r.auth.VerifyStore)

		pickupTemplates.POST("", r.create.Handle)
		pickupTemplates.GET("", r.list.Handle)
		pickupTemplates.PUT("/:id", r.update.Handle)
		pickupTemplates.DELETE("/:id", r.delete.Handle)
	}
}
