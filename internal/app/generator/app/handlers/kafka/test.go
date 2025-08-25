package kafka

import (
	"path"

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
	test := tmpl.Template{
		SourcePath: "templates/internal/domain/handlers/kafka/crud_test.go.tmpl",
		DestinationPath: path.Join(
			"internal",
			"app",
			g.domain.AppConfig.AppName(),
			"handlers",
			"kafka",
			g.domain.DirName(),
			g.domain.TestFileName(),
		),
		Name: "test http handler",
	}
	if err := test.RenderToFile(g.domain); err != nil {
		return err
	}
	return nil
}
