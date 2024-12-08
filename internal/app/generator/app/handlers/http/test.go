package http

import (
	"path"

	"github.com/mikalai-mitsin/creathor/internal/pkg/domain"
	"github.com/mikalai-mitsin/creathor/internal/pkg/tmpl"
)

type TestGenerator struct {
	domain *domain.Domain
}

func NewTestGenerator(domain *domain.Domain) *TestGenerator {
	return &TestGenerator{domain: domain}
}

func (g *TestGenerator) Sync() error {
	test := tmpl.Template{
		SourcePath: "templates/internal/domain/handlers/http/crud_test.go.tmpl",
		DestinationPath: path.Join(
			destinationPath,
			"internal",
			"app",
			g.domain.DirName(),
			"handlers",
			"http",
			g.domain.TestFileName(),
		),
		Name: "test http handler",
	}
	if err := test.RenderToFile(g.domain); err != nil {
		return err
	}
	return nil
}
