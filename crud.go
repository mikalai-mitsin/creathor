package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Model struct {
	Model  string
	Module string
}

func CreateCRUD(name, module string) error {
	name = strings.Title(name)
	filename := fmt.Sprintf("%s.go", strings.ToLower(name))
	files := []*Template{
		{
			SourcePath:      "templates/models/model.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "domain", "models", filename),
			Name:            "model",
		},
		{
			SourcePath:      "templates/models/repository.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "domain", "repositories", filename),
			Name:            "repository",
		},
		{
			SourcePath:      "templates/models/usecase.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "domain", "usecases", filename),
			Name:            "usecase",
		},
		{
			SourcePath:      "templates/models/interceptor.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "domain", "interceptors", filename),
			Name:            "interceptor",
		},
	}
	data := Model{
		Model:  name,
		Module: module,
	}
	for _, tmpl := range files {
		if err := tmpl.renderToFile(data); err != nil {
			return err
		}
	}
	return nil
}
