package usecases

import (
	"path/filepath"

	"github.com/mikalai-mitsin/creathor/internal/pkg/app"

	"github.com/mikalai-mitsin/creathor/internal/pkg/tmpl"
)

type TestGenerator struct {
	domain *app.BaseEntity
}

func NewTestGenerator(domain *app.BaseEntity) *TestGenerator {
	return &TestGenerator{domain: domain}
}

func (g *TestGenerator) Sync() error {
	//return nil
	test := tmpl.Template{
		SourcePath: "templates/internal/domain/usecases/crud_test.go.tmpl",
		DestinationPath: filepath.Join(
			destinationPath,
			"internal",
			"app",
			g.domain.AppName(),
			"usecases",
			g.domain.TestFileName(),
		),
		Name: "usecase test",
	}
	if err := test.RenderToFile(g.domain); err != nil {
		return err
	}
	return nil
}
