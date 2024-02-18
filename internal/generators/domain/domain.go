package domain

import (
	"github.com/018bf/creathor/internal/domain"
	"github.com/018bf/creathor/internal/generators"
	"github.com/018bf/creathor/internal/generators/domain/interceptors"
	"github.com/018bf/creathor/internal/generators/domain/interfaces/grpc"
	"github.com/018bf/creathor/internal/generators/domain/models"
	"github.com/018bf/creathor/internal/generators/domain/repositories/postgres"
	"github.com/018bf/creathor/internal/generators/domain/usecases"
)

type Generator struct {
	domain *domain.Domain
}

func NewGenerator(d *domain.Domain) *Generator {
	return &Generator{domain: d}
}

func (g *Generator) Sync() error {
	domainGenerators := []generators.Generator{
		interceptors.NewInterceptorCrud(g.domain),
		interceptors.NewUseCaseInterfaceCrud(g.domain),

		usecases.NewUseCaseCrud(g.domain),
		usecases.NewRepositoryInterfaceCrud(g.domain),

		postgres.NewRepositoryCrud(g.domain),

		grpc.NewHandler(g.domain),
		grpc.NewInterceptorInterfaceCrud(g.domain),
	}
	for _, model := range g.domain.Models {
		domainGenerators = append(domainGenerators, models.NewModel(model, g.domain.Filename, g.domain))
	}
	for _, generator := range domainGenerators {
		if err := generator.Sync(); err != nil {
			return err
		}
	}
	return nil
}
