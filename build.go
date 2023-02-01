package main

import (
	"github.com/018bf/creathor/models"
	"path/filepath"
)

func CreateBuild(data *models.Project) error {
	files := []*Template{
		{
			SourcePath:      "templates/build/Dockerfile.tmpl",
			DestinationPath: filepath.Join(destinationPath, "build", "Dockerfile"),
			Name:            "Dockerfile",
		},
	}
	if data.MakeEnabled {
		files = append(
			files,
			&Template{
				SourcePath:      "templates/build/Makefile.tmpl",
				DestinationPath: filepath.Join(destinationPath, "Makefile"),
				Name:            "Makefile",
			},
		)
	}
	if data.TaskEnabled {
		files = append(
			files,
			&Template{
				SourcePath:      "templates/build/Taskfile.yaml.tmpl",
				DestinationPath: filepath.Join(destinationPath, "Taskfile.yaml"),
				Name:            "Taskfile.yaml",
			},
		)
	}
	for _, tmpl := range files {
		if err := tmpl.renderToFile(data); err != nil {
			return err
		}
	}
	return nil
}
