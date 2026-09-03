package routes

import (
	"github.com/Jeff1507/go-to-save-api/internal/core/preset/handler"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/middleware"
	"github.com/gin-gonic/gin"
)

type PresetRoutes struct {
	auth   *middleware.Authentication
	list   *handler.ListPresetsHandler
	create *handler.CreatePresetHandler
	update *handler.UpdatePresetHandler
	delete *handler.DeletePresetHandler
}

func NewPresetRoutes(
	auth *middleware.Authentication,
	list *handler.ListPresetsHandler,
	create *handler.CreatePresetHandler,
	update *handler.UpdatePresetHandler,
	delete *handler.DeletePresetHandler,
) *PresetRoutes {
	return &PresetRoutes{
		auth:   auth,
		list:   list,
		create: create,
		update: update,
		delete: delete,
	}
}

func (r *PresetRoutes) Register(router *gin.Engine) {
	presets := router.Group("/pickup-presets")
	presets.Use(handler.ErrorHandler)

	{
		presets.Use(r.auth.AuthenticateStore)

		presets.POST("", r.create.Handle)
		presets.GET("", r.list.Handle)
		presets.PUT("/:id", r.update.Handle)
		presets.DELETE("/:id", r.delete.Handle)
	}
}
