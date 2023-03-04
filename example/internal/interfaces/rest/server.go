package rest

import (
	"context"
	"net/http"

	"github.com/018bf/example/internal/configs"

	"github.com/018bf/example/pkg/log"
	"github.com/gin-gonic/gin"
)

type ctxKey string

const (
	UserContextKey      ctxKey = "user"
	RequestIDContextKey ctxKey = "request_id"
)

type Server struct {
	router *gin.Engine
	config *configs.Config
	logger log.Logger
}

// NewServer        godoc
// @title           example
// @version         0.1.0
// @description     TBD
// @host      127.0.0.1:8000
// @BasePath  /api/v1
func NewServer(
	logger log.Logger,
	config *configs.Config,
	authMiddleware *AuthMiddleware,
	authHandler *AuthHandler,
	userHandler *UserHandler,
) *Server {
	router := gin.Default()
	router.Use(authMiddleware.Middleware())
	router.Use(RequestMiddleware)
	router.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})
	apiV1 := router.Group("api").Group("v1")
	authHandler.Register(apiV1)
	userHandler.Register(apiV1)
	return &Server{
		router: router,
		config: config,
		logger: logger,
	}
}

func (s *Server) Start(_ context.Context) error {
	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *Server) Stop(_ context.Context) error {
	return nil
}
