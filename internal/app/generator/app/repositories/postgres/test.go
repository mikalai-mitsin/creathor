package postgres

import (
	"path/filepath"

	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"

	"github.com/mikalai-mitsin/creathor/internal/pkg/tmpl"
)

type TestGenerator struct {
	domain *configs.EntityConfig
}

func NewTestGenerator(domain *configs.EntityConfig) *TestGenerator {
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
			g.domain.AppConfig.AppName(),
			"repositories",
			"postgres",
			g.domain.DirName(),
			g.domain.TestFileName(),
		),
		Name: "repository test",
	}
	if err := test.RenderToFile(g.domain); err != nil {
		return err
	}
	return nil
}
