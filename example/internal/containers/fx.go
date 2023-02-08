package containers

import (
	"context"
	gatewayInterface "github.com/018bf/example/internal/interfaces/gateway"
	grpcInterface "github.com/018bf/example/internal/interfaces/grpc"
	postgresInterface "github.com/018bf/example/internal/interfaces/postgres"
	restInterface "github.com/018bf/example/internal/interfaces/rest"
	jwtRepositories "github.com/018bf/example/internal/repositories/jwt"
	postgresRepositories "github.com/018bf/example/internal/repositories/postgres"

	"github.com/018bf/example/pkg/log"
	"go.uber.org/fx/fxevent"

	"github.com/018bf/example/internal/interceptors"
	"github.com/018bf/example/internal/usecases"

	"github.com/018bf/example/pkg/clock"

	"github.com/018bf/example/internal/configs"
	"go.uber.org/fx"
)

var FXModule = fx.Options(
	fx.WithLogger(
		func(logger log.Logger) fxevent.Logger {
			return logger
		},
	),
	fx.Provide(
		context.Background,
		configs.ParseConfig,
		clock.NewRealClock,
		postgresInterface.NewDatabase,
		postgresInterface.NewMigrateManager,
		grpcInterface.NewServer,
		restInterface.NewServer,
		func(config *configs.Config) (log.Logger, error) {
			return log.NewLog(config.LogLevel)
		},
		usecases.NewAuthUseCase,
		interceptors.NewAuthInterceptor,
		jwtRepositories.NewAuthRepository,
		postgresRepositories.NewPermissionRepository,
		grpcInterface.NewAuthMiddleware,
		grpcInterface.NewAuthServiceServer,
		grpcInterface.NewUserServiceServer,
		restInterface.NewAuthHandler,
		restInterface.NewAuthMiddleware,
		restInterface.NewUserHandler,
		gatewayInterface.NewServer,
		usecases.NewUserUseCase,
		interceptors.NewUserInterceptor,
		postgresRepositories.NewPostgresUserRepository, usecases.NewSessionUseCase, interceptors.NewSessionInterceptor, postgresRepositories.NewSessionRepository, restInterface.NewSessionHandler, grpcInterface.NewSessionServiceServer, usecases.NewEquipmentUseCase, interceptors.NewEquipmentInterceptor, postgresRepositories.NewEquipmentRepository, restInterface.NewEquipmentHandler, grpcInterface.NewEquipmentServiceServer, usecases.NewPlanUseCase, interceptors.NewPlanInterceptor, postgresRepositories.NewPlanRepository, restInterface.NewPlanHandler, grpcInterface.NewPlanServiceServer, usecases.NewDayUseCase, interceptors.NewDayInterceptor, postgresRepositories.NewDayRepository, restInterface.NewDayHandler, grpcInterface.NewDayServiceServer, usecases.NewArchUseCase, interceptors.NewArchInterceptor, postgresRepositories.NewArchRepository, restInterface.NewArchHandler, grpcInterface.NewArchServiceServer,
	),
)

func NewGRPCExample(config string) *fx.App {
	app := fx.New(
		fx.Provide(func() string { return config }),
		FXModule,
		fx.Invoke(func(lifecycle fx.Lifecycle, server *grpcInterface.Server) {
			lifecycle.Append(fx.Hook{
				OnStart: server.Start,
				OnStop:  server.Stop,
			})
		}),
	)
	return app
}
func NewGatewayExample(config string) *fx.App {
	app := fx.New(
		fx.Provide(func() string { return config }),
		FXModule,
		fx.Invoke(func(lifecycle fx.Lifecycle, server *gatewayInterface.Server) {
			lifecycle.Append(fx.Hook{
				OnStart: server.Start,
			})
		}),
	)
	return app
}
func NewRESTExample(config string) *fx.App {
	app := fx.New(
		fx.Provide(func() string { return config }),
		FXModule,
		fx.Invoke(func(lifecycle fx.Lifecycle, server *restInterface.Server) {
			lifecycle.Append(fx.Hook{
				OnStart: server.Start,
				OnStop:  server.Stop,
			})
		}),
	)
	return app
}

func NewMigrate(config string) *fx.App {
	app := fx.New(
		fx.Provide(func() string { return config }),
		FXModule,
		fx.Invoke(func(lifecycle fx.Lifecycle, manager *postgresInterface.MigrateManager) {
			lifecycle.Append(fx.Hook{
				OnStart: manager.Up,
			})
		}),
	)
	return app
}
