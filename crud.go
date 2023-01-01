package main

import (
	"path/filepath"
)

func CreateCRUD(data Model) error {

	files := []*Template{
		{
			SourcePath:      "templates/domain/model.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "domain", "models", data.FileName()),
			Name:            "model",
		},
		{
			SourcePath:      "templates/domain/model_mock.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "domain", "models", "mock", data.FileName()),
			Name:            "model_mock",
		},
		{
			SourcePath:      "templates/domain/repository.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "domain", "repositories", data.FileName()),
			Name:            "repository",
		},
		{
			SourcePath:      "templates/domain/usecase.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "domain", "usecases", data.FileName()),
			Name:            "usecase",
		},
		{
			SourcePath:      "templates/domain/interceptor.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "domain", "interceptors", data.FileName()),
			Name:            "interceptor",
		},
		{
			SourcePath:      "templates/implementations/usecase.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "usecases", data.FileName()),
			Name:            "usecase",
		},
		{
			SourcePath:      "templates/implementations/usecase_test.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "usecases", data.TestFileName()),
			Name:            "usecase test",
		},
		{
			SourcePath:      "templates/implementations/interceptor.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "interceptors", data.FileName()),
			Name:            "interceptor",
		},
		{
			SourcePath:      "templates/implementations/interceptor_test.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "interceptors", data.TestFileName()),
			Name:            "interceptor test",
		},
		{
			SourcePath:      "templates/implementations/repository.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "repositories", data.FileName()),
			Name:            "repository",
		},
		{
			SourcePath:      "templates/implementations/repository_test.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "repositories", data.TestFileName()),
			Name:            "repository test",
		},
	}
	for _, tmpl := range files {
		if err := tmpl.renderToFile(data); err != nil {
			return err
		}
	}
	return nil
}
