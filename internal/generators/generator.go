package generators

import (
	"github.com/018bf/creathor/internal/configs"
	"github.com/018bf/creathor/internal/generators/domain/errs"
	interceptorInterfaces "github.com/018bf/creathor/internal/generators/domain/interceptors"
	"github.com/018bf/creathor/internal/generators/domain/models"
	"github.com/018bf/creathor/internal/generators/domain/repositories"
	repositoryInterfaces "github.com/018bf/creathor/internal/generators/domain/repositories"
	useCaseInterfaces "github.com/018bf/creathor/internal/generators/domain/usecases"
	"github.com/018bf/creathor/internal/generators/interceptors"
	"github.com/018bf/creathor/internal/generators/interfaces/grpc"
	"github.com/018bf/creathor/internal/generators/interfaces/uptrace"
	"github.com/018bf/creathor/internal/generators/repositories/postgres"
	"github.com/018bf/creathor/internal/generators/usecases"
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
		errs.NewErrors(g.project),
		models.NewModelTypes(g.project),
	}
	for _, m := range g.project.Models {
		generators = append(
			generators,
			models.NewMainModel(m),
			models.NewCreateModel(m),
			models.NewUpdateModel(m),
			models.NewFilterModel(m),
			repositories.NewRepositoryInterfaceCrud(m),
			useCaseInterfaces.NewUseCaseInterfaceCrud(m),
			interceptorInterfaces.NewInterceptorInterfaceCrud(m),

			interceptors.NewInterceptorCrud(m),
			usecases.NewUseCaseCrud(m),
			postgres.NewRepositoryCrud(m),

			grpc.NewHandler(m),
		)
	}
	if g.project.UptraceEnabled {
		generators = append(generators, uptrace.NewProvider(g.project))
	}
	if g.project.Auth {
		generators = append(
			generators,
			models.NewModelAuth(g.project),
			models.NewModelUser(g.project),
			models.NewModelPermission(g.project),

			repositoryInterfaces.NewRepositoryInterfaceUser(g.project),
			repositoryInterfaces.NewRepositoryInterfacePermission(g.project),
			repositoryInterfaces.NewRepositoryInterfaceAuth(g.project),
			useCaseInterfaces.NewUseCaseInterfaceAuth(g.project),
			useCaseInterfaces.NewUseCaseInterfaceUser(g.project),
			interceptorInterfaces.NewInterceptorInterfaceAuth(g.project),
			interceptorInterfaces.NewInterceptorInterfaceUser(g.project),

			usecases.NewUseCaseUser(g.project),
			usecases.NewUseCaseAuth(g.project),
			interceptors.NewInterceptorUser(g.project),
			interceptors.NewInterceptorAuth(g.project),
		)
	}
	if g.project.KafkaEnabled {
		generators = append(
			generators,
			models.NewModelEvent(g.project),
			repositoryInterfaces.NewRepositoryInterfaceEvent(g.project),
		)
	}
	for _, generator := range generators {
		if err := generator.Sync(); err != nil {
			return err
		}
	}
	return nil
}
