package containers

import (
	"context"

	"github.com/018bf/example/internal/configs"
	"github.com/018bf/example/internal/interceptors"
	gatewayInterface "github.com/018bf/example/internal/interfaces/gateway"
	grpcInterface "github.com/018bf/example/internal/interfaces/grpc"
	postgresInterface "github.com/018bf/example/internal/interfaces/postgres"
	restInterface "github.com/018bf/example/internal/interfaces/rest"
	"github.com/018bf/example/internal/interfaces/uptrace"
	uptraceInterface "github.com/018bf/example/internal/interfaces/uptrace"
	jwtRepositories "github.com/018bf/example/internal/repositories/jwt"
	postgresRepositories "github.com/018bf/example/internal/repositories/postgres"
	"github.com/018bf/example/internal/usecases"
	"github.com/018bf/example/pkg/clock"
	"github.com/018bf/example/pkg/log"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

var FXModule = fx.Options(fx.WithLogger(func(logger log.Logger) fxevent.Logger {
	return logger
}), fx.Provide(func(config *configs.Config) (log.Logger, error) {
	return log.NewLog(config.LogLevel)
}, context.Background, configs.ParseConfig, clock.NewRealClock, postgresInterface.NewDatabase, postgresInterface.NewMigrateManager, grpcInterface.NewServer, grpcInterface.NewRequestIDMiddleware, grpcInterface.NewAuthMiddleware, grpcInterface.NewAuthServiceServer, grpcInterface.NewUserServiceServer, uptraceInterface.NewProvider, restInterface.NewServer, restInterface.NewAuthMiddleware, restInterface.NewAuthHandler, restInterface.NewUserHandler, gatewayInterface.NewServer, interceptors.NewAuthInterceptor, usecases.NewAuthUseCase, jwtRepositories.NewAuthRepository, postgresRepositories.NewPermissionRepository, interceptors.NewUserInterceptor, usecases.NewUserUseCase, postgresRepositories.NewPostgresUserRepository, grpcInterface.NewSessionServiceServer, restInterface.NewSessionHandler, interceptors.NewSessionInterceptor, usecases.NewSessionUseCase, postgresRepositories.NewSessionRepository, grpcInterface.NewEquipmentServiceServer, restInterface.NewEquipmentHandler, interceptors.NewEquipmentInterceptor, usecases.NewEquipmentUseCase, postgresRepositories.NewEquipmentRepository, grpcInterface.NewPlanServiceServer, restInterface.NewPlanHandler, interceptors.NewPlanInterceptor, usecases.NewPlanUseCase, postgresRepositories.NewPlanRepository, grpcInterface.NewDayServiceServer, restInterface.NewDayHandler, interceptors.NewDayInterceptor, usecases.NewDayUseCase, postgresRepositories.NewDayRepository, grpcInterface.NewArchServiceServer, restInterface.NewArchHandler, interceptors.NewArchInterceptor, usecases.NewArchUseCase, postgresRepositories.NewArchRepository), fx.Invoke(func(lifecycle fx.Lifecycle, server *uptrace.Provider, config *configs.Config) {
	lifecycle.Append(fx.Hook{OnStart: server.Start, OnStop: server.Stop})
}))

func NewMigrateContainer(config string) *fx.App {
	app := fx.New(fx.Provide(func() string {
		return config
	}), FXModule, fx.Invoke(func(lifecycle fx.Lifecycle, logger log.Logger, manager *postgresInterface.MigrateManager, shutdowner fx.Shutdowner) {
		lifecycle.Append(fx.Hook{OnStart: func(ctx context.Context) error {
			go func() {
				err := manager.Up(ctx)
				if err != nil {
					logger.Error("shutdown", log.Any("error", err))
					_ = shutdowner.Shutdown()
				}
			}()
			return nil
		}})
	}))
	return app
}
func NewGRPCContainer(config string) *fx.App {
	app := fx.New(fx.Provide(func() string {
		return config
	}), FXModule, fx.Invoke(func(lifecycle fx.Lifecycle, logger log.Logger, server *grpcInterface.Server, shutdowner fx.Shutdowner) {
		lifecycle.Append(fx.Hook{OnStart: func(ctx context.Context) error {
			go func() {
				err := server.Start(ctx)
				if err != nil {
					logger.Error("shutdown", log.Any("error", err))
					_ = shutdowner.Shutdown()
				}
			}()
			return nil
		}, OnStop: server.Stop})
	}))
	return app
}
func NewGatewayContainer(config string) *fx.App {
	app := fx.New(fx.Provide(func() string {
		return config
	}), FXModule, fx.Invoke(func(lifecycle fx.Lifecycle, logger log.Logger, server *gatewayInterface.Server, shutdowner fx.Shutdowner) {
		lifecycle.Append(fx.Hook{OnStart: func(ctx context.Context) error {
			go func() {
				err := server.Start(ctx)
				if err != nil {
					logger.Error("shutdown", log.Any("error", err))
					_ = shutdowner.Shutdown()
				}
			}()
			return nil
		}})
	}))
	return app
}
func NewRESTContainer(config string) *fx.App {
	app := fx.New(fx.Provide(func() string {
		return config
	}), FXModule, fx.Invoke(func(lifecycle fx.Lifecycle, logger log.Logger, server *restInterface.Server, shutdowner fx.Shutdowner) {
		lifecycle.Append(fx.Hook{OnStart: func(ctx context.Context) error {
			go func() {
				err := server.Start(ctx)
				if err != nil {
					logger.Error("shutdown", log.Any("error", err))
					_ = shutdowner.Shutdown()
				}
			}()
			return nil
		}, OnStop: server.Stop})
	}))
	return app
}
