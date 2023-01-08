package containers

import (
	"context"
	"github.com/018bf/example/internal/interfaces/grpc"

	"github.com/018bf/example/pkg/log"
	"go.uber.org/fx/fxevent"

	"github.com/018bf/example/internal/interceptors"
	"github.com/018bf/example/internal/interfaces/postgres"
	"github.com/018bf/example/internal/repositories"
	"github.com/018bf/example/internal/usecases"

	"github.com/018bf/example/pkg/clock"

	"github.com/018bf/example/internal/configs"
	"go.uber.org/fx"
)

var appModule = fx.Options(
	fx.Provide(
		context.Background,
		clock.NewRealClock,
		func(config *configs.Config) (log.Logger, error) {
			return log.NewLog(config.LogLevel)
		},
	),
	configs.FXModule,
	repositories.FXModule,
	usecases.FXModule,
	interceptors.FXModule,
	postgres.FXModule,
)

func NewExample(config string) *fx.App {
	app := fx.New(
		fx.Provide(func() string { return config }),
		appModule,
		fx.WithLogger(
			func(logger log.Logger) fxevent.Logger {
				return logger
			},
		),
		grpc.FXModule,
	)
	return app
}
