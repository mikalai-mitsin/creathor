package layout

import (
	"github.com/018bf/creathor/internal/configs"
	generators2 "github.com/018bf/creathor/internal/generators"
	"github.com/018bf/creathor/internal/generators/layout/domain/errs"
	interceptors3 "github.com/018bf/creathor/internal/generators/layout/domain/interceptors"
	models2 "github.com/018bf/creathor/internal/generators/layout/domain/models"
	"github.com/018bf/creathor/internal/generators/layout/domain/repositories"
	usecases3 "github.com/018bf/creathor/internal/generators/layout/domain/usecases"
	interceptors2 "github.com/018bf/creathor/internal/generators/layout/interceptors"
	"github.com/018bf/creathor/internal/generators/layout/interfaces/grpc"
	"github.com/018bf/creathor/internal/generators/layout/interfaces/uptrace"
	usecases2 "github.com/018bf/creathor/internal/generators/layout/usecases"
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
		models2.NewModelTypes(g.project),
	}
	if g.project.UptraceEnabled {
		generators = append(generators, uptrace.NewProvider(g.project))
	}
	if g.project.Auth {
		generators = append(
			generators,
			models2.NewModelAuth(g.project),
			//models2.NewModelUser(g.project),
			models2.NewModelPermission(g.project),

			//repositories.NewRepositoryInterfaceUser(g.project),
			repositories.NewRepositoryInterfacePermission(g.project),
			repositories.NewRepositoryInterfaceAuth(g.project),
			usecases3.NewUseCaseInterfaceAuth(g.project),
			//usecases3.NewUseCaseInterfaceUser(g.project),
			interceptors3.NewInterceptorInterfaceAuth(g.project),
			//interceptors3.NewInterceptorInterfaceUser(g.project),

			//usecases2.NewUseCaseUser(g.project),
			usecases2.NewUseCaseAuth(g.project),
			//interceptors2.NewInterceptorUser(g.project),
			interceptors2.NewInterceptorAuth(g.project),
		)
	}
	if g.project.KafkaEnabled {
		generators = append(
			generators,
			models2.NewModelEvent(g.project),
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
