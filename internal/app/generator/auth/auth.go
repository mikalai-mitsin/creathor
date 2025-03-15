package auth

import (
	"path"

	authModel "github.com/mikalai-mitsin/creathor/internal/app/generator/auth/entities"
	authGrpcHandlers "github.com/mikalai-mitsin/creathor/internal/app/generator/auth/handlers/grpc"
	authHttpHandlers "github.com/mikalai-mitsin/creathor/internal/app/generator/auth/handlers/http"
	authRepositoriesJwt "github.com/mikalai-mitsin/creathor/internal/app/generator/auth/repositories/jwt"
	authRepositoriesPosgres "github.com/mikalai-mitsin/creathor/internal/app/generator/auth/repositories/postgres"
	authServices "github.com/mikalai-mitsin/creathor/internal/app/generator/auth/services"
	authUseCases "github.com/mikalai-mitsin/creathor/internal/app/generator/auth/usecases"

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
		authServices.NewServiceAuth(g.project),
		authServices.NewRepositoryInterfaceAuth(g.project),
		//UseCase and interfaces
		authUseCases.NewUseCaseAuth(g.project),
		authUseCases.NewServiceInterfaceAuth(g.project),
		//gRPC Handlers and interfaces
		authGrpcHandlers.NewInterfaces(g.project),
		authGrpcHandlers.NewHandler(g.project),
		authGrpcHandlers.NewMiddlewares(g.project),
		authGrpcHandlers.NewProto(g.project),
		// HTTP Handlers and interfaces
		authHttpHandlers.NewInterfaces(g.project),
		authHttpHandlers.NewHandler(g.project),
		authHttpHandlers.NewDTOGenerator(g.project),
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
			SourcePath: "templates/internal/auth/services/auth_test.go.tmpl",
			DestinationPath: path.Join(
				destinationPath,
				"internal",
				"app",
				"auth",
				"services",
				"auth_test.go",
			),
			Name: "test auth service implementation",
		},
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
