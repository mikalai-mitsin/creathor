package layout

import (
	"github.com/018bf/creathor/internal/configs"
	generators2 "github.com/018bf/creathor/internal/generators"
	authGrpcHandlers "github.com/018bf/creathor/internal/generators/auth/handlers/grpc"
	authInterceptors "github.com/018bf/creathor/internal/generators/auth/interceptors"
	authModel "github.com/018bf/creathor/internal/generators/auth/models"
	authUseCases "github.com/018bf/creathor/internal/generators/auth/usecases"
	"github.com/018bf/creathor/internal/generators/layout/domain/repositories"
	"github.com/018bf/creathor/internal/generators/layout/errs"
	"github.com/018bf/creathor/internal/generators/layout/interfaces/grpc"
	"github.com/018bf/creathor/internal/generators/layout/interfaces/uptrace"
	userModel "github.com/018bf/creathor/internal/generators/user/models"
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
	}
	if g.project.UptraceEnabled {
		generators = append(generators, uptrace.NewProvider(g.project))
	}
	if g.project.Auth {
		generators = append(
			generators,
			authModel.NewModelAuth(g.project),
			//TODO: auth repository
			//Use case and interfaces
			authUseCases.NewUseCaseAuth(g.project),
			authUseCases.NewRepositoryInterfaceAuth(g.project),
			//Interceptor and interfaces
			authInterceptors.NewInterceptorAuth(g.project),
			authInterceptors.NewUseCaseInterfaceAuth(g.project),
			//Handlers and interfaces
			authGrpcHandlers.NewInterceptorInterfaceAuth(g.project),
			//
			//userModel.NewModelUser(g.project),
			userModel.NewModelPermission(g.project),
			//
			//Use case and interfaces
			//userUseCases.NewUseCaseUser(g.project),
			//userUseCases.NewRepositoryInterfaceUser(g.project),
			//Interceptor and interfaces
			//userInterceptors.NewInterceptorUser(g.project),
			//userInterceptors.NewUseCaseInterfaceUser(g.project),
			//Handlers and interfaces
			//userGrpcHandlers.OldNewInterceptorInterfaceUser(g.project),
		)
	}
	if g.project.KafkaEnabled {
		generators = append(
			generators,
			//models2.NewModelEvent(g.project),
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
