package routes

import (
	"github.com/Jeff1507/go-to-save-api/internal/core/review/handler"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/middleware"
	"github.com/gin-gonic/gin"
)

type ReviewRoutes struct {
	auth *middleware.Authentication
	add  *handler.AddReviewHandler
}

func NewReviewRoutes(
	auth *middleware.Authentication,
	add *handler.AddReviewHandler,
) *ReviewRoutes {
	return &ReviewRoutes{
		auth: auth,
		add:  add,
	}
}

func (r *ReviewRoutes) Register(router *gin.Engine) {
	reviews := router.Group("/reviews")
	reviews.Use(handler.ErrorHandler)
	reviews.Use(r.auth.AuthenticateCustomer)
	reviews.POST("/:orderId", r.add.Handle)
}
