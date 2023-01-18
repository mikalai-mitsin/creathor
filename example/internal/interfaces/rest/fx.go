package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/018bf/example/internal/configs"
	"go.uber.org/fx"
)

var FXModule = fx.Options(
	fx.Provide(NewSessionHandler, NewEquipmentHandler),
	fx.Provide(NewAuthMiddleware, NewAuthHandler, NewUserHandler),
	fx.Provide(NewRouter),
	fx.Invoke(func(lifecycle fx.Lifecycle, router *gin.Engine, config *configs.Config) {
		lifecycle.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				return http.ListenAndServe(config.BindAddr, router)
			},
			OnStop: func(ctx context.Context) error {
				return nil
			},
		})
	}),
)
