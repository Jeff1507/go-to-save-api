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
	store := router.Group("/store/pickup-presets")
	store.Use(handler.ErrorHandler)
	store.Use(r.auth.AuthenticateStore)
	store.POST("", r.create.Handle)
	store.GET("", r.list.Handle)
	store.PUT("/:id", r.update.Handle)
	store.DELETE("/:id", r.delete.Handle)
}
