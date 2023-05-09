package generators

import (
	interceptorInterfaces "github.com/018bf/creathor/internal/generators/domain/interceptors"
	"github.com/018bf/creathor/internal/generators/domain/models"
	repositoryInterfaces "github.com/018bf/creathor/internal/generators/domain/repositories"
	useCaseInterfaces "github.com/018bf/creathor/internal/generators/domain/usecases"
	"github.com/018bf/creathor/internal/generators/interceptors"
	"github.com/018bf/creathor/internal/generators/interfaces/grpc"
	"github.com/018bf/creathor/internal/generators/repositories/postgres"
	"github.com/018bf/creathor/internal/generators/usecases"
	"github.com/018bf/creathor/internal/mods"
)

type ModGenerator struct {
	mod *mods.Mod
}

func NewModGenerator(mod *mods.Mod) *ModGenerator {
	return &ModGenerator{mod: mod}
}

func (g *ModGenerator) Sync() error {
	generators := []Generator{
		interceptors.NewInterceptorCrud(g.mod),
		repositoryInterfaces.NewRepositoryInterfaceCrud(g.mod),
		useCaseInterfaces.NewUseCaseInterfaceCrud(g.mod),
		interceptorInterfaces.NewInterceptorInterfaceCrud(g.mod),
		usecases.NewUseCaseCrud(g.mod),
		postgres.NewRepositoryCrud(g.mod),
		grpc.NewHandler(g.mod),
	}
	for _, model := range g.mod.Models {
		generators = append(generators, models.NewModel(model, g.mod.Filename))
	}
	for _, generator := range generators {
		if err := generator.Sync(); err != nil {
			return err
		}
	}
	return nil
}
