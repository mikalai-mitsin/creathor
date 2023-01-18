package rest

import (
	"net/http"

	"github.com/018bf/example/pkg/log"
	"github.com/gin-gonic/gin"
)

type ctxKey string

const (
	RequestIDContextKey ctxKey = "request_id"
)

func NewRouter(
	logger log.Logger, userHandler *UserHandler, equipmentHandler *EquipmentHandler, sessionHandler *SessionHandler, approachHandler *ApproachHandler,
) *gin.Engine {
	router := gin.Default()
	router.Use(RequestMiddleware)
	router.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})
	userHandler.Register(router)
	equipmentHandler.Register(router)
	sessionHandler.Register(router)
	approachHandler.Register(router)
	return router
}
