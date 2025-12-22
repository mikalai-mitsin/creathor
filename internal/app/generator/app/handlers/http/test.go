package http

import (
	"path"

	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
	"github.com/mikalai-mitsin/creathor/internal/pkg/tmpl"
)

type TestGenerator struct {
	entityConfig configs.EntityConfig
}

func NewTestGenerator(entityConfig configs.EntityConfig) *TestGenerator {
	return &TestGenerator{entityConfig: entityConfig}
}

func (g *TestGenerator) Sync() error {
	test := tmpl.Template{
		SourcePath: "templates/internal/domain/handlers/http/crud_test.go.tmpl",
		DestinationPath: path.Join(
			destinationPath,
			"internal",
			"app",
			g.entityConfig.AppConfig.AppName(),
			"handlers",
			"http",
			g.entityConfig.DirName(),
			g.entityConfig.TestFileName(),
		),
		Name: "test http handler",
	}
	if err := test.RenderToFile(&g.entityConfig); err != nil {
		return err
	}
	return nil
}
