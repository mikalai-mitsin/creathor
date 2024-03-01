package postgres

import (
	"path"

	"github.com/018bf/creathor/internal/pkg/configs"
	"github.com/018bf/creathor/internal/pkg/tmpl"
)

var destinationPath = "."

type Generator struct {
	project *configs.Project
}

func NewGenerator(project *configs.Project) *Generator {
	return &Generator{project: project}
}

func (c *Generator) Sync() error {
	files := []*tmpl.Template{
		{
			SourcePath: "templates/internal/pkg/postgres/postgres.go.tmpl",
			DestinationPath: path.Join(
				destinationPath,
				"internal",
				"pkg",
				"postgres",
				"postgres.go",
			),
			Name: "postgres",
		},
		{
			SourcePath:      "templates/internal/pkg/postgres/search.go.tmpl",
			DestinationPath: path.Join(destinationPath, "internal", "pkg", "postgres", "search.go"),
			Name:            "search",
		},
		{
			SourcePath: "templates/internal/pkg/postgres/testing.go.tmpl",
			DestinationPath: path.Join(
				destinationPath,
				"internal",
				"pkg",
				"postgres",
				"testing.go",
			),
			Name: "postgres testing",
		},
		{
			SourcePath: "templates/internal/pkg/postgres/migrations/init.sql.tmpl",
			DestinationPath: path.Join(
				destinationPath,
				"internal",
				"pkg",
				"postgres",
				"migrations",
				"000001_init.up.sql",
			),
			Name: "postgres init migration",
		},
	}
	for _, file := range files {
		if err := file.RenderToFile(c.project); err != nil {
			return err
		}
	}
	return nil
}
