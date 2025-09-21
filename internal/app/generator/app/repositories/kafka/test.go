package kafka

import (
	"path/filepath"

	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"

	"github.com/mikalai-mitsin/creathor/internal/pkg/tmpl"
)

type TestGenerator struct {
	entityConfig configs.EntityConfig
}

func NewProducerTestGenerator(entityConfig configs.EntityConfig) *TestGenerator {
	return &TestGenerator{entityConfig: entityConfig}
}

func (g *TestGenerator) Sync() error {
	//return nil
	test := tmpl.Template{
		SourcePath: "templates/internal/domain/repositories/kafka/event_test.go.tmpl",
		DestinationPath: filepath.Join(
			".",
			"internal",
			"app",
			g.entityConfig.AppConfig.AppName(),
			"repositories",
			"kafka",
			g.entityConfig.DirName(),
			g.entityConfig.TestFileName(),
		),
		Name: "producer test",
	}
	if err := test.RenderToFile(&g.entityConfig); err != nil {
		return err
	}
	return nil
}
