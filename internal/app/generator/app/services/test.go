package services

import (
	"path/filepath"

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
	//return nil
	test := tmpl.Template{
		SourcePath: "templates/internal/domain/services/crud_test.go.tmpl",
		DestinationPath: filepath.Join(
			destinationPath,
			"internal",
			"app",
			g.entityConfig.AppConfig.AppName(),
			"services",
			g.entityConfig.DirName(),
			g.entityConfig.TestFileName(),
		),
		Name: "service test",
	}
	if err := test.RenderToFile(&g.entityConfig); err != nil {
		return err
	}
	return nil
}
