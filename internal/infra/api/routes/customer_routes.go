package routes

import (
	"github.com/Jeff1507/go-to-save-api/internal/core/customer/handler"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/middleware"
	"github.com/go-chi/chi/v5"
)

type CustomerRoutes struct {
	auth   *middleware.Authentication
	get    *handler.GetCustomer
	signup *handler.Signup
	login  *handler.Login
}

func NewCustomerRoutes(
	auth *middleware.Authentication,
	get *handler.GetCustomer,
	signup *handler.Signup,
	login *handler.Login,
) *CustomerRoutes {
	return &CustomerRoutes{
		auth:   auth,
		get:    get,
		signup: signup,
		login:  login,
	}
}

func (r *CustomerRoutes) Register(router *chi.Mux) {
	router.Route("/customers", func(customers chi.Router) {

		// Rotas públicas
		customers.Group(func(public chi.Router) {
			public.Post("/signup", r.signup.Handle)
			public.Post("/login", r.login.Handle)
		})

		// Rotas privadas
		customers.Group(func(private chi.Router) {
			private.Use(r.auth.CustomerAuth)

			private.Get("/me", handler.ErrorHandler(r.get.Handle))
			private.Post("/logout", nil)
		})
	})
}
