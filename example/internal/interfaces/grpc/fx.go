package grpc

import (
	"context"
	"net"

	"github.com/018bf/example/internal/configs"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

var FXModule = fx.Options(
	fx.Provide(NewAuthMiddleware),
	fx.Provide(NewServer),
	fx.Invoke(func(lifecycle fx.Lifecycle, server *grpc.Server, config *configs.Config) {
		lifecycle.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				listener, err := net.Listen("tcp", config.BindAddr)
				if err != nil {
					return err
				}
				return server.Serve(listener)
			},
			OnStop: func(ctx context.Context) error {
				server.GracefulStop()
				return nil
			},
		})
	}),
)
