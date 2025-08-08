package kafka

import (
	"os"
	"path"
	"path/filepath"

	"github.com/mikalai-mitsin/creathor/internal/pkg/app"

	"github.com/mikalai-mitsin/creathor/internal/pkg/tmpl"
)

type ProducerGenerator struct {
	domain *app.BaseEntity
}

func NewProducerGenerator(domain *app.BaseEntity) *ProducerGenerator {
	return &ProducerGenerator{domain: domain}
}

func (g *ProducerGenerator) Sync() error {
	err := os.MkdirAll(path.Dir(g.filename()), 0777)
	if err != nil {
		return err
	}
	test := tmpl.Template{
		SourcePath:      "templates/internal/domain/repositories/kafka/event.go.tmpl",
		DestinationPath: g.filename(),
		Name:            "producer",
	}
	if err := test.RenderToFile(g.domain); err != nil {
		return err
	}
	return nil
}

func (g *ProducerGenerator) filename() string {
	return filepath.Join(
		".",
		"internal",
		"app",
		g.domain.AppName(),
		"events",
		g.domain.DirName(),
		"producer.go",
	)
}
