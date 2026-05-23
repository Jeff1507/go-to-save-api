package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Route interface {
	Register(router *chi.Mux)
}

type Server struct {
	router *chi.Mux
}

func New() *Server {
	router := chi.NewRouter()

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
	apiRouter := chi.NewRouter()

	for _, route := range routes {
		route.Register(apiRouter)
	}

	srv.router.Mount("/api", apiRouter)
}

func (srv *Server) Run(port string) error {
	addr := ":" + port
	if err := http.ListenAndServe(addr, srv.router); err != nil {
		return fmt.Errorf("listen and serve: %v", err)
	}

	return nil
}
