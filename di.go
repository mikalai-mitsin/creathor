package main

import (
	"os"
	"path"
)

func CreateDI(data *Project) error {
	directories := []string{
		path.Join(destinationPath, "internal", "containers"),
	}
	for _, directory := range directories {
		if err := os.MkdirAll(directory, 0777); err != nil {
			return NewUnexpectedBehaviorError(err.Error())
		}
	}
	files := []*Template{
		{
			SourcePath:      "templates/containers/fx.go.tmpl",
			DestinationPath: path.Join(destinationPath, "internal", "containers", "fx.go"),
			Name:            "Uber FX DI container",
		},
		{
			SourcePath:      "templates/containers/configs.go.tmpl",
			DestinationPath: path.Join(destinationPath, "internal", "configs", "fx.go"),
			Name:            "Configs FX module",
		},
		{
			SourcePath:      "templates/containers/repositories.go.tmpl",
			DestinationPath: path.Join(destinationPath, "internal", "repositories", "fx.go"),
			Name:            "Repositories FX module",
		},
		{
			SourcePath:      "templates/containers/usecases.go.tmpl",
			DestinationPath: path.Join(destinationPath, "internal", "usecases", "fx.go"),
			Name:            "Use Cases FX module",
		},
		{
			SourcePath:      "templates/containers/interceptors.go.tmpl",
			DestinationPath: path.Join(destinationPath, "internal", "interceptors", "fx.go"),
			Name:            "Interceptors FX module",
		},
	}

	for _, tmpl := range files {
		if err := tmpl.renderToFile(data); err != nil {
			return err
		}
	}
	return nil
}
