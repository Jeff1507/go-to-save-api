package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Route interface {
	Register(router *gin.Engine)
}

type Server struct {
	router *gin.Engine
}

func New(allowedOrigins []string, filesDir string) *Server {
	setupValidator()

	router := gin.Default()

	router.SetTrustedProxies(nil)

	router.MaxMultipartMemory = 10 << 20 // 10 MB

	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300 * time.Second,
	}))

	router.Static("/uploads", filesDir)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "TUDO OK 👍")
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
	if err := srv.router.Run(addr); err != nil {
		return fmt.Errorf("run gin router: %w", err)
	}

	return nil
}
