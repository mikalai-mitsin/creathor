package postgres

import (
	"path/filepath"

	"github.com/mikalai-mitsin/creathor/internal/pkg/app"

	"github.com/mikalai-mitsin/creathor/internal/pkg/tmpl"
)

type TestGenerator struct {
	domain *app.BaseEntity
}

func NewTestGenerator(domain *app.BaseEntity) *TestGenerator {
	return &TestGenerator{domain: domain}
}

func (g *TestGenerator) Sync() error {
	//return nil
	test := tmpl.Template{
		SourcePath: "templates/internal/domain/repositories/postgres/crud_test.go.tmpl",
		DestinationPath: filepath.Join(
			destinationPath,
			"internal",
			"app",
			g.domain.AppName(),
			"repositories",
			"postgres",
			g.domain.TestFileName(),
		),
		Name: "repository test",
	}
	if err := test.RenderToFile(g.domain); err != nil {
		return err
	}
	return nil
}
