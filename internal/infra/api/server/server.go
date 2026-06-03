package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Route interface {
	Register(router *chi.Mux)
}

type Server struct {
	router *chi.Mux
}

func New(allowedOrigins []string) *Server {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("TUDO OK 👍"))
	})

	return &Server{
		router: router,
	}
}

func (srv *Server) AddRoutes(routes ...Route) {
	for _, route := range routes {
		route.Register(srv.router)
	}
}

func (srv *Server) Run(port string) error {
	addr := ":" + port
	if err := http.ListenAndServe(addr, srv.router); err != nil {
		return fmt.Errorf("listen and serve: %v", err)
	}

	return nil
}
