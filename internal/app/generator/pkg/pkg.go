package pkg

import (
	"github.com/mikalai-mitsin/creathor/internal/app/generator"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/pkg/clock"
	cfg "github.com/mikalai-mitsin/creathor/internal/app/generator/pkg/configs"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/pkg/containers"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/pkg/dtx"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/pkg/errs"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/pkg/grpc"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/pkg/http"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/pkg/kafka"
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
	generators := []generator.Generator{
		clock.NewGenerator(g.project),
		cfg.NewGenerator(g.project),
		containers.NewGenerator(g.project),
		errs.NewGenerator(g.project),
		log.NewGenerator(g.project),
		pointer.NewGenerator(g.project),
		postgres.NewGenerator(g.project),
		uuid.NewGenerator(g.project),
		dtx.NewGenerator(g.project),
	}
	if g.project.KafkaEnabled {
		generators = append(
			generators,
			kafka.NewConfigGenerator(g.project),
			kafka.NewConsumerGenerator(g.project),
			kafka.NewProducerGenerator(g.project),
		)
	}
	if g.project.HTTPEnabled {
		generators = append(generators, http.NewConfig(g.project), http.NewServer(g.project))
	}
	if g.project.GRPCEnabled {
		generators = append(
			generators,
			grpc.NewConfig(g.project),
			grpc.NewMiddlewares(g.project),
			grpc.NewServer(g.project),
		)
	}
	if g.project.UptraceEnabled {
		generators = append(generators, uptrace.NewProvider(g.project))
	}
	for _, gen := range generators {
		if err := gen.Sync(); err != nil {
			return err
		}
	}
	return nil
}
