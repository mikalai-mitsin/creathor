package module

import (
	generators2 "github.com/018bf/creathor/internal/generators"
	"github.com/018bf/creathor/internal/generators/models"
	"github.com/018bf/creathor/internal/generators/module/domain/interceptors"
	"github.com/018bf/creathor/internal/generators/module/domain/repositories"
	"github.com/018bf/creathor/internal/generators/module/domain/usecases"
	interceptors2 "github.com/018bf/creathor/internal/generators/module/interceptors"
	"github.com/018bf/creathor/internal/generators/module/interfaces/grpc"
	"github.com/018bf/creathor/internal/generators/module/repositories/postgres"
	usecases2 "github.com/018bf/creathor/internal/generators/module/usecases"
	mods "github.com/018bf/creathor/internal/module"
)

type Generator struct {
	mod *mods.Mod
}

func NewGenerator(mod *mods.Mod) *Generator {
	return &Generator{mod: mod}
}

func (g *Generator) Sync() error {
	generators := []generators2.Generator{
		interceptors2.NewInterceptorCrud(g.mod),
		repositories.NewRepositoryInterfaceCrud(g.mod),
		usecases.NewUseCaseInterfaceCrud(g.mod),
		interceptors.NewInterceptorInterfaceCrud(g.mod),
		usecases2.NewUseCaseCrud(g.mod),
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
