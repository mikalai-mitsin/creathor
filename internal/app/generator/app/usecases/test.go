package usecases

import (
	"path/filepath"

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
	//return nil
	test := tmpl.Template{
		SourcePath: "templates/internal/domain/usecases/crud_test.go.tmpl",
		DestinationPath: filepath.Join(
			destinationPath,
			"internal",
			"app",
			g.domain.DirName(),
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
