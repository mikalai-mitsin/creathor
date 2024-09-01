package services

import (
	"github.com/mikalai-mitsin/creathor/internal/pkg/domain"
	"path/filepath"

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
