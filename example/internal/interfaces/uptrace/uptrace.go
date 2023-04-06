package uptrace

import (
	"context"

	"github.com/018bf/example"
	"github.com/018bf/example/internal/configs"
	"github.com/uptrace/uptrace-go/uptrace"
)

type Provider struct {
	config *configs.Config
}

func (p Provider) Stop(ctx context.Context) error {
	return uptrace.Shutdown(ctx)
}
func (p Provider) Start(_ context.Context) error {
	if p.config.Otel.Enabled {
		uptrace.ConfigureOpentelemetry(
			uptrace.WithDSN(p.config.Otel.URL),
			uptrace.WithServiceName(example.Name),
			uptrace.WithServiceVersion(example.Version),
			uptrace.WithDeploymentEnvironment(p.config.Otel.Environment),
		)
	}
	return nil
}
func NewProvider(config *configs.Config) *Provider {
	return &Provider{config: config}
}
