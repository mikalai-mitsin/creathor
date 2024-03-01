package layout

import (
	"path"

	"github.com/018bf/creathor/internal/pkg/configs"
	"github.com/018bf/creathor/internal/pkg/tmpl"
)

type CmdGenerator struct {
	project *configs.Project
}

func NewCmdGenerator(project *configs.Project) *CmdGenerator {
	return &CmdGenerator{project: project}
}

func (c *CmdGenerator) Sync() error {
	files := []*tmpl.Template{
		{
			SourcePath:      "templates/cmd/service/main.go.tmpl",
			DestinationPath: path.Join(destinationPath, "cmd", c.project.Name, "main.go"),
			Name:            "service main",
		},
		{
			SourcePath:      "templates/go.mod.tmpl",
			DestinationPath: path.Join(destinationPath, "go.mod"),
			Name:            "go.mod",
		},
		{
			SourcePath:      "templates/version.go.tmpl",
			DestinationPath: path.Join(destinationPath, "version.go"),
			Name:            "version",
		},
	}
	for _, file := range files {
		if err := file.RenderToFile(c.project); err != nil {
			return err
		}
	}
	return nil
}
