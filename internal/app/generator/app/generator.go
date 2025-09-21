package app

import (
	"github.com/mikalai-mitsin/creathor/internal/app/generator"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/app/api/proto"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/app/entities"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/app/handlers/grpc"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/app/handlers/http"
	handlersKafka "github.com/mikalai-mitsin/creathor/internal/app/generator/app/handlers/kafka"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/app/repositories/kafka"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/app/repositories/postgres"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/app/services"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/app/usecases"
	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
)

type Generator struct {
	appConfig *configs.AppConfig
}

func NewGenerator(d *configs.AppConfig) *Generator {
	return &Generator{appConfig: d}
}

func (g *Generator) Sync() error {
	appGenerators := []generator.Generator{NewApp(g.appConfig)}
	for _, entity := range g.appConfig.Entities {
		appGenerators = append(appGenerators,
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
		if g.appConfig.KafkaEnabled {
			appGenerators = append(
				appGenerators,
				proto.NewProtoGenerator(entity),
				kafka.NewProducerGenerator(entity),
				kafka.NewInterfacesGenerator(entity),
				kafka.NewProducerTestGenerator(entity),
				kafka.NewProtoDecoder(entity),
				services.NewEventService(entity),
				handlersKafka.NewHandlerGenerator(entity),
				handlersKafka.NewInterfacesGenerator(entity),
			)
		}
		if g.appConfig.HTTPEnabled {
			appGenerators = append(
				appGenerators,
				http.NewDTOGenerator(entity),
				http.NewHandlerGenerator(entity),
				http.NewInterfacesGenerator(entity),
			)
		}
		if g.appConfig.GRPCEnabled {
			appGenerators = append(
				appGenerators,
				proto.NewProtoGenerator(entity),
				grpc.NewInterfacesGenerator(entity),
				grpc.NewHandlerGenerator(entity),
				grpc.NewTestGenerator(entity),
				grpc.NewProtoEncoder(entity),
				grpc.NewProtoDecoder(entity),
			)
		}
		for _, baseEntity := range entity.Entities {
			appGenerators = append(appGenerators, entities.NewModel(baseEntity, entity))
		}
	}
	for _, appGenerator := range appGenerators {
		if err := appGenerator.Sync(); err != nil {
			return err
		}
	}
	return nil
}
