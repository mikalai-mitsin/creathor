package pkg

import (
	"github.com/018bf/creathor/internal/configs"
	generators2 "github.com/018bf/creathor/internal/generators"
	"github.com/018bf/creathor/internal/generators/pkg/containers"
	"github.com/018bf/creathor/internal/generators/pkg/domain/repositories"
	"github.com/018bf/creathor/internal/generators/pkg/errs"
	"github.com/018bf/creathor/internal/generators/pkg/grpc"
	"github.com/018bf/creathor/internal/generators/pkg/uptrace"
)

type Generator struct {
	project *configs.Project
}

func NewGenerator(project *configs.Project) *Generator {
	return &Generator{project: project}
}

func (g *Generator) Sync() error {
	generators := []generators2.Generator{
		grpc.NewServer(g.project),
		errs.NewErrors(g.project),
		containers.NewFxContainer(g.project),
	}
	if g.project.UptraceEnabled {
		generators = append(generators, uptrace.NewProvider(g.project))
	}
	if g.project.KafkaEnabled {
		generators = append(
			generators,
			repositories.NewRepositoryInterfaceEvent(g.project),
		)
	}
	for _, generator := range generators {
		if err := generator.Sync(); err != nil {
			return err
		}
	}
	return nil
}
