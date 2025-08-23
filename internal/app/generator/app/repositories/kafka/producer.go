package kafka

import (
	"os"
	"path"
	"path/filepath"

	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"

	"github.com/mikalai-mitsin/creathor/internal/pkg/tmpl"
)

type ProducerGenerator struct {
	domain *configs.EntityConfig
}

func NewProducerGenerator(domain *configs.EntityConfig) *ProducerGenerator {
	return &ProducerGenerator{domain: domain}
}

func (r *ProducerGenerator) Sync() error {
	err := os.MkdirAll(path.Dir(r.filename()), 0777)
	if err != nil {
		return err
	}
	test := tmpl.Template{
		SourcePath:      "templates/internal/domain/repositories/kafka/event.go.tmpl",
		DestinationPath: r.filename(),
		Name:            "producer",
	}
	if err := test.RenderToFile(r.domain); err != nil {
		return err
	}
	return nil
}

func (r *ProducerGenerator) filename() string {
	return filepath.Join(
		".",
		"internal",
		"app",
		r.domain.AppConfig.AppName(),
		"repositories",
		"kafka",
		r.domain.DirName(),
		r.domain.FileName(),
	)
}
