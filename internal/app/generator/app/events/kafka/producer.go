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
		r.domain.AppName(),
		"events",
		r.domain.DirName(),
		r.domain.FileName(),
	)
}
