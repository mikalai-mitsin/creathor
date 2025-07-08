package services

import (
	"path/filepath"

	"github.com/mikalai-mitsin/creathor/internal/pkg/domain"

	"github.com/mikalai-mitsin/creathor/internal/pkg/tmpl"
)

type TestGenerator struct {
	domain *domain.App
}

func NewTestGenerator(domain *domain.App) *TestGenerator {
	return &TestGenerator{domain: domain}
}

func (g *TestGenerator) Sync() error {
	//return nil
	test := tmpl.Template{
		SourcePath: "templates/internal/domain/services/crud_test.go.tmpl",
		DestinationPath: filepath.Join(
			destinationPath,
			"internal",
			"app",
			g.domain.DirName(),
			"services",
			g.domain.TestFileName(),
		),
		Name: "service test",
	}
	if err := test.RenderToFile(g.domain); err != nil {
		return err
	}
	return nil
}
