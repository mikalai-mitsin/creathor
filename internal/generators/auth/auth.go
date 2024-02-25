package auth

import (
	"github.com/018bf/creathor/internal/configs"
	generators2 "github.com/018bf/creathor/internal/generators"
	authGrpcHandlers "github.com/018bf/creathor/internal/generators/auth/handlers/grpc"
	authInterceptors "github.com/018bf/creathor/internal/generators/auth/interceptors"
	authModel "github.com/018bf/creathor/internal/generators/auth/models"
	authUseCases "github.com/018bf/creathor/internal/generators/auth/usecases"
)

type Generator struct {
	project *configs.Project
}

func NewGenerator(project *configs.Project) *Generator {
	return &Generator{project: project}
}

func (g *Generator) Sync() error {
	var generators []generators2.Generator
	if g.project.Auth {
		generators = append(
			generators,
			authModel.NewModelAuth(g.project),
			//Use case and interfaces
			authUseCases.NewUseCaseAuth(g.project),
			authUseCases.NewRepositoryInterfaceAuth(g.project),
			//Interceptor and interfaces
			authInterceptors.NewInterceptorAuth(g.project),
			authInterceptors.NewUseCaseInterfaceAuth(g.project),
			//Handlers and interfaces
			authGrpcHandlers.NewInterceptorInterfaceAuth(g.project),
			authModel.NewModelPermission(g.project),
		)
	}
	for _, generator := range generators {
		if err := generator.Sync(); err != nil {
			return err
		}
	}
	return nil
}
