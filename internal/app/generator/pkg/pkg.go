package pkg

import (
	generators2 "github.com/mikalai-mitsin/creathor/internal/app/generator"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/pkg/auth"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/pkg/clock"
	cg "github.com/mikalai-mitsin/creathor/internal/app/generator/pkg/configs"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/pkg/containers"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/pkg/domain/repositories"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/pkg/errs"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/pkg/grpc"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/pkg/log"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/pkg/pointer"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/pkg/postgres"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/pkg/uptrace"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/pkg/uuid"
	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
)

type Generator struct {
	project *configs.Project
}

func NewGenerator(project *configs.Project) *Generator {
	return &Generator{project: project}
}

func (g *Generator) Sync() error {
	generators := []generators2.Generator{
		clock.NewGenerator(g.project),
		cg.NewConfigGenerator(g.project),
		containers.NewFxContainer(g.project),
		errs.NewErrors(g.project),
		grpc.NewMiddlewares(g.project),
		grpc.NewServer(g.project),
		log.NewGenerator(g.project),
		pointer.NewGenerator(g.project),
		postgres.NewGenerator(g.project),
		uuid.NewGenerator(g.project),
	}
	if g.project.Auth {
		generators = append(generators, auth.NewGenerator(g.project))
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
