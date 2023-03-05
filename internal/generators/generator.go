package generators

import (
	"github.com/018bf/creathor/internal/configs"
	"github.com/018bf/creathor/internal/generators/domain"
	"github.com/018bf/creathor/internal/generators/implementations"
	"github.com/018bf/creathor/internal/generators/interfaces/grpc"
	"github.com/018bf/creathor/internal/generators/interfaces/uptrace"
)

type Generator interface {
	Sync() error
}

type CrudGenerator struct {
	project *configs.Project
}

func NewCrudGenerator(project *configs.Project) *CrudGenerator {
	return &CrudGenerator{project: project}
}

func (g CrudGenerator) Sync() error {
	generators := []Generator{
		grpc.NewServer(g.project),
		uptrace.NewProvider(g.project),
		domain.NewErrors(g.project),
	}
	for _, m := range g.project.Models {
		generators = append(
			generators,
			domain.NewMainModel(m),
			domain.NewCreateModel(m),
			domain.NewUpdateModel(m),
			domain.NewFilterModel(m),
			domain.NewRepositoryInterface(m),
			domain.NewUseCaseInterface(m),
			domain.NewInterceptorInterface(m),

			implementations.NewInterceptor(m),
			implementations.NewUseCase(m),
			implementations.NewPostgresRepository(m),

			grpc.NewHandler(m),
		)
	}
	for _, generator := range generators {
		if err := generator.Sync(); err != nil {
			return err
		}
	}
	return nil
}
