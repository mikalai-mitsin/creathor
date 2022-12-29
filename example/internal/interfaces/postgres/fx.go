package postgres

import (
	"context"
	"github.com/018bf/example/internal/configs"

	"go.uber.org/fx"
)

var FXModule = fx.Options(
	fx.Provide(
		NewDatabase,
	),
)

var FXApp = fx.Options(
	fx.Provide(
		NewMigrateManager,
	),
	fx.Invoke(func(lifecycle fx.Lifecycle, manager *MigrateManager, config *configs.Config) {
		lifecycle.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				if err := manager.Up(config); err != nil {
					return err
				}
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return nil
			},
		})
	}),
)
