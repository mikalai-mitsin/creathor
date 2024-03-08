package uuid

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
			SourcePath:      "templates/internal/pkg/uuid/uuid.go.tmpl",
			DestinationPath: path.Join(destinationPath, "internal", "pkg", "uuid", "uuid.go"),
			Name:            "utils pointer",
		},
	}
	for _, file := range files {
		if err := file.RenderToFile(c.project); err != nil {
			return err
		}
	}
	return nil
}
