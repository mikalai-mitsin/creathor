package rest

import (
	"context"
	"errors"
	"net/http"

	"github.com/018bf/example/internal/configs"
	"github.com/018bf/example/internal/domain/errs"
	"github.com/gin-contrib/cors"

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
// @schemes https http
//
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
//
// @security ApiKeyAuth
func NewServer(
	logger log.Logger,
	config *configs.Config,
	authMiddleware *AuthMiddleware,
	authHandler *AuthHandler,
	userHandler *UserHandler,
	sessionHandler *SessionHandler,
	equipmentHandler *EquipmentHandler,
	planHandler *PlanHandler,
	dayHandler *DayHandler,
	archHandler *ArchHandler,
) *Server {
	router := gin.Default()
	router.Use(authMiddleware.Middleware())
	router.Use(cors.Default())
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
	sessionHandler.Register(apiV1)
	equipmentHandler.Register(apiV1)
	planHandler.Register(apiV1)
	dayHandler.Register(apiV1)
	archHandler.Register(apiV1)
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

func decodeError(ctx *gin.Context, err error) {
	var domainError *errs.Error
	if errors.As(err, &domainError) {
		switch domainError.Code {
		case errs.ErrorCodeOK:
			ctx.JSON(http.StatusOK, err)
		case errs.ErrorCodeCanceled:
			ctx.JSON(http.StatusInternalServerError, err)
		case errs.ErrorCodeUnknown:
			ctx.JSON(http.StatusInternalServerError, err)
		case errs.ErrorCodeInvalidArgument:
			ctx.JSON(http.StatusBadRequest, err)
		case errs.ErrorCodeDeadlineExceeded:
			ctx.JSON(http.StatusInternalServerError, err)
		case errs.ErrorCodeNotFound:
			ctx.JSON(http.StatusNotFound, err)
		case errs.ErrorCodeAlreadyExists:
			ctx.JSON(http.StatusBadRequest, err)
		case errs.ErrorCodePermissionDenied:
			ctx.JSON(http.StatusForbidden, err)
		case errs.ErrorCodeResourceExhausted:
			ctx.JSON(http.StatusInternalServerError, err)
		case errs.ErrorCodeFailedPrecondition:
			ctx.JSON(http.StatusBadRequest, err)
		case errs.ErrorCodeAborted:
			ctx.JSON(http.StatusInternalServerError, err)
		case errs.ErrorCodeOutOfRange:
			ctx.JSON(http.StatusInternalServerError, err)
		case errs.ErrorCodeUnimplemented:
			ctx.JSON(http.StatusMethodNotAllowed, err)
		case errs.ErrorCodeInternal:
			ctx.JSON(http.StatusInternalServerError, err)
		case errs.ErrorCodeUnavailable:
			ctx.JSON(http.StatusServiceUnavailable, err)
		case errs.ErrorCodeDataLoss:
			ctx.JSON(http.StatusInternalServerError, err)
		case errs.ErrorCodeUnauthenticated:
			ctx.JSON(http.StatusUnauthorized, err)
		default:
			ctx.JSON(http.StatusInternalServerError, err)
		}
		return
	}
	ctx.String(http.StatusInternalServerError, err.Error())
}
