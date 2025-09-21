package kafka

import (
	"os"
	"path"
	"path/filepath"

	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"

	"github.com/mikalai-mitsin/creathor/internal/pkg/tmpl"
)

type ProducerGenerator struct {
	entityConfig configs.EntityConfig
}

func NewProducerGenerator(entityConfig configs.EntityConfig) *ProducerGenerator {
	return &ProducerGenerator{entityConfig: entityConfig}
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
	if err := test.RenderToFile(&r.entityConfig); err != nil {
		return err
	}
	return nil
}

func (r *ProducerGenerator) filename() string {
	return filepath.Join(
		".",
		"internal",
		"app",
		r.entityConfig.AppConfig.AppName(),
		"repositories",
		"kafka",
		r.entityConfig.DirName(),
		r.entityConfig.FileName(),
	)
}
