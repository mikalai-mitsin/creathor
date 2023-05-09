package generators

import (
	"github.com/018bf/creathor/internal/configs"
	"github.com/018bf/creathor/internal/generators/domain/errs"
	interceptorInterfaces "github.com/018bf/creathor/internal/generators/domain/interceptors"
	"github.com/018bf/creathor/internal/generators/domain/models"
	repositoryInterfaces "github.com/018bf/creathor/internal/generators/domain/repositories"
	useCaseInterfaces "github.com/018bf/creathor/internal/generators/domain/usecases"
	"github.com/018bf/creathor/internal/generators/interceptors"
	"github.com/018bf/creathor/internal/generators/interfaces/grpc"
	"github.com/018bf/creathor/internal/generators/interfaces/uptrace"
	"github.com/018bf/creathor/internal/generators/usecases"
)

type LayoutGenerator struct {
	project *configs.Project
}

func NewLayoutGenerator(project *configs.Project) *LayoutGenerator {
	return &LayoutGenerator{project: project}
}

func (g *LayoutGenerator) Sync() error {
	generators := []Generator{
		grpc.NewServer(g.project),
		errs.NewErrors(g.project),
		models.NewModelTypes(g.project),
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
