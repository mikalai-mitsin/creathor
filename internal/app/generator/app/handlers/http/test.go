package http

import (
	"path"

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
	test := tmpl.Template{
		SourcePath: "templates/internal/domain/handlers/http/crud_test.go.tmpl",
		DestinationPath: path.Join(
			destinationPath,
			"internal",
			"app",
			g.domain.AppName(),
			"handlers",
			"http",
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
