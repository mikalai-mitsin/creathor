package layout

import (
	"path/filepath"

	"github.com/018bf/creathor/internal/pkg/tmpl"

	"github.com/018bf/creathor/internal/pkg/configs"
)

type BuilderGenerator struct {
	project *configs.Project
}

func NewBuilderGenerator(project *configs.Project) *BuilderGenerator {
	return &BuilderGenerator{project: project}
}

func (b *BuilderGenerator) Sync() error {
	files := []*tmpl.Template{
		{
			SourcePath:      "templates/build/Dockerfile.tmpl",
			DestinationPath: filepath.Join(destinationPath, "build", "Dockerfile"),
			Name:            "Dockerfile",
		},
	}
	if b.project.MakeEnabled {
		files = append(
			files,
			&tmpl.Template{
				SourcePath:      "templates/build/Makefile.tmpl",
				DestinationPath: filepath.Join(destinationPath, "Makefile"),
				Name:            "Makefile",
			},
		)
	}
	if b.project.TaskEnabled {
		files = append(
			files,
			&tmpl.Template{
				SourcePath:      "templates/build/Taskfile.yaml.tmpl",
				DestinationPath: filepath.Join(destinationPath, "Taskfile.yaml"),
				Name:            "Taskfile.yaml",
			},
		)
	}
	for _, file := range files {
		if err := file.RenderToFile(b.project); err != nil {
			return err
		}
	}
	return nil
}
