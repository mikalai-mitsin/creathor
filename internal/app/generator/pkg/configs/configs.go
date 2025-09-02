package configs

import (
	"path"

	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
	"github.com/mikalai-mitsin/creathor/internal/pkg/tmpl"
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
			SourcePath:      "templates/configs/config.toml.tmpl",
			DestinationPath: path.Join(destinationPath, "configs", "config.toml"),
			Name:            "main config",
		},
		{
			SourcePath:      "templates/configs/config.toml.tmpl",
			DestinationPath: path.Join(destinationPath, "configs", "ci.toml"),
			Name:            "ci config",
		},
		{
			SourcePath:      "templates/configs/config.toml.tmpl",
			DestinationPath: path.Join(destinationPath, "configs", "test.toml"),
			Name:            "test config",
		},
		{
			SourcePath:      "templates/internal/pkg/configs/config.go.tmpl",
			DestinationPath: path.Join(destinationPath, "internal", "pkg", "configs", "config.go"),
			Name:            "config struct",
		},
		{
			SourcePath: "templates/internal/pkg/configs/config_test.go.tmpl",
			DestinationPath: path.Join(
				destinationPath,
				"internal",
				"pkg",
				"configs",
				"config_test.go",
			),
			Name: "config tests",
		},
	}
	for _, file := range files {
		if err := file.RenderToFile(c.project); err != nil {
			return err
		}
	}
	return nil
}
