package kafka

import (
	"path/filepath"

	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"

	"github.com/mikalai-mitsin/creathor/internal/pkg/tmpl"
)

type TestGenerator struct {
	domain *configs.BaseEntity
}

func NewProducerTestGenerator(domain *configs.BaseEntity) *TestGenerator {
	return &TestGenerator{domain: domain}
}

func (g *TestGenerator) Sync() error {
	//return nil
	test := tmpl.Template{
		SourcePath: "templates/internal/domain/repositories/kafka/event_test.go.tmpl",
		DestinationPath: filepath.Join(
			".",
			"internal",
			"app",
			g.domain.AppName(),
			"repositories",
			"kafka",
			g.domain.DirName(),
			g.domain.TestFileName(),
		),
		Name: "producer test",
	}
	if err := test.RenderToFile(g.domain); err != nil {
		return err
	}
	return nil
}
