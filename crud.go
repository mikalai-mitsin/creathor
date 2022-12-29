package main

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"path/filepath"
	"strings"
)

type Model struct {
	Model  string
	Module string
}

func CreateCRUD(name, module string) error {
	name = cases.Title(language.English).String(name)
	filename := fmt.Sprintf("%s.go", strings.ToLower(name))
	testFilename := fmt.Sprintf("%s_test.go", strings.ToLower(name))
	files := []*Template{
		{
			SourcePath:      "templates/domain/model.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "domain", "models", filename),
			Name:            "model",
		},
		{
			SourcePath:      "templates/domain/model_mock.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "domain", "models", "mock", filename),
			Name:            "model_mock",
		},
		{
			SourcePath:      "templates/domain/repository.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "domain", "repositories", filename),
			Name:            "repository",
		},
		{
			SourcePath:      "templates/domain/usecase.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "domain", "usecases", filename),
			Name:            "usecase",
		},
		{
			SourcePath:      "templates/domain/interceptor.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "domain", "interceptors", filename),
			Name:            "interceptor",
		},
		{
			SourcePath:      "templates/implementations/usecase.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "usecases", filename),
			Name:            "usecase",
		},
		{
			SourcePath:      "templates/implementations/usecase_test.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "usecases", testFilename),
			Name:            "usecase test",
		},
		{
			SourcePath:      "templates/implementations/interceptor.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "interceptors", filename),
			Name:            "interceptor",
		},
		{
			SourcePath:      "templates/implementations/repository.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "repositories", filename),
			Name:            "repository",
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
