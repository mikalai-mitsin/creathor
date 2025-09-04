package dtx

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
			SourcePath:      "templates/internal/pkg/dtx/manager.go.tmpl",
			DestinationPath: path.Join(destinationPath, "internal", "pkg", "dtx", "manager.go"),
			Name:            "manager",
		},
		{
			SourcePath:      "templates/internal/pkg/dtx/tx.go.tmpl",
			DestinationPath: path.Join(destinationPath, "internal", "pkg", "dtx", "tx.go"),
			Name:            "tx",
		},
	}
	for _, file := range files {
		if err := file.RenderToFile(c.project); err != nil {
			return err
		}
	}
	return nil
}
