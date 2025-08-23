package app

import (
	"github.com/mikalai-mitsin/creathor/internal/app/generator"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/app/entities"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/app/handlers/grpc"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/app/handlers/http"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/app/repositories/kafka"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/app/repositories/postgres"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/app/services"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/app/usecases"
	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
)

type Generator struct {
	domain *configs.GeneratorConfig
}

func NewGenerator(d *configs.GeneratorConfig) *Generator {
	return &Generator{domain: d}
}

func (g *Generator) Sync() error {
	domainGenerators := []generator.Generator{NewApp(g.domain)}
	for _, entity := range g.domain.Entities {
		domainGenerators = append(domainGenerators,
			usecases.NewInterfacesGenerator(entity),
			usecases.NewUseCaseGenerator(entity),
			usecases.NewTestGenerator(entity),

			services.NewInterfacesGenerator(entity),
			services.NewServiceGenerator(entity),
			services.NewTestGenerator(entity),

			postgres.NewInterfacesGenerator(entity),
			postgres.NewRepositoryGenerator(entity),
			postgres.NewTestGenerator(entity),
		)
		if g.domain.Config.KafkaEnabled {
			domainGenerators = append(
				domainGenerators,
				kafka.NewProducerGenerator(entity),
				kafka.NewInterfacesGenerator(entity),
				kafka.NewProducerTestGenerator(entity),
			)
		}
		if g.domain.Config.HTTPEnabled {
			domainGenerators = append(
				domainGenerators,
				http.NewDTOGenerator(entity),
				http.NewHandlerGenerator(entity),
				http.NewInterfacesGenerator(entity),
			)
		}
		if g.domain.Config.GRPCEnabled {
			domainGenerators = append(
				domainGenerators,
				grpc.NewProtoGenerator(entity),
				grpc.NewInterfacesGenerator(entity),
				grpc.NewHandlerGenerator(entity),
				grpc.NewTestGenerator(entity),
			)
		}
		for _, baseEntity := range entity.Entities {
			domainGenerators = append(domainGenerators, entities.NewModel(baseEntity, entity))
		}
	}
	for _, domainGenerator := range domainGenerators {
		if err := domainGenerator.Sync(); err != nil {
			return err
		}
	}
	return nil
}
