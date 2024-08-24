package auth

import (
	authGrpcHandlers "github.com/mikalai-mitsin/creathor/internal/app/generator/auth/handlers/grpc"
	authInterceptors "github.com/mikalai-mitsin/creathor/internal/app/generator/auth/interceptors"
	authModel "github.com/mikalai-mitsin/creathor/internal/app/generator/auth/models"
	authRepositoriesJwt "github.com/mikalai-mitsin/creathor/internal/app/generator/auth/repositories/jwt"
	authRepositoriesPosgres "github.com/mikalai-mitsin/creathor/internal/app/generator/auth/repositories/postgres"
	authUseCases "github.com/mikalai-mitsin/creathor/internal/app/generator/auth/usecases"
	"path"

	"github.com/mikalai-mitsin/creathor/internal/app/generator"
	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
	"github.com/mikalai-mitsin/creathor/internal/pkg/tmpl"
)

type Generator struct {
	project *configs.Project
}

func NewGenerator(project *configs.Project) *Generator {
	return &Generator{project: project}
}

var destinationPath = "."

func (g *Generator) Sync() error {
	var authGenerators []generator.Generator
	authGenerators = append(
		authGenerators,
		authModel.NewModelAuth(g.project),
		authRepositoriesJwt.NewRepository(g.project),
		authRepositoriesPosgres.NewRepository(g.project),
		//Use case and interfaces
		authUseCases.NewUseCaseAuth(g.project),
		authUseCases.NewRepositoryInterfaceAuth(g.project),
		//Interceptor and interfaces
		authInterceptors.NewInterceptorAuth(g.project),
		authInterceptors.NewUseCaseInterfaceAuth(g.project),
		//Handlers and interfaces
		authGrpcHandlers.NewInterceptorInterfaceAuth(g.project),
		authGrpcHandlers.NewHandler(g.project),
		authGrpcHandlers.NewProto(g.project),

		authModel.NewModelPermission(g.project),

		//App constructor
		NewAppAuth(g.project),
	)
	for _, authGenerator := range authGenerators {
		if err := authGenerator.Sync(); err != nil {
			return err
		}
	}
	tests := []*tmpl.Template{
		{
			SourcePath: "templates/internal/auth/usecases/auth_test.go.tmpl",
			DestinationPath: path.Join(
				destinationPath,
				"internal",
				"app",
				"auth",
				"usecases",
				"auth_test.go",
			),
			Name: "test auth usecase implementation",
		},
		{
			SourcePath: "templates/internal/auth/interceptors/auth_test.go.tmpl",
			DestinationPath: path.Join(
				destinationPath,
				"internal",
				"app",
				"auth",
				"interceptors",
				"auth_test.go",
			),
			Name: "test auth interceptor implementation",
		},
		{
			SourcePath: "templates/internal/auth/handlers/grpc/auth_test.go.tmpl",
			DestinationPath: path.Join(
				destinationPath,
				"internal",
				"app",
				"auth",
				"handlers",
				"grpc",
				"auth_test.go",
			),
			Name: "grpc auth test",
		},
	}
	for _, test := range tests {
		if err := test.RenderToFile(g.project); err != nil {
			return err
		}
	}
	return nil
}
