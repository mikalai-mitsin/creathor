package main

import (
	"fmt"
	"path/filepath"
)

func CreateCRUD(data Model) error {
	files := []*Template{
		{
			SourcePath:      "templates/internal/domain/models/crud.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "domain", "models", data.FileName()),
			Name:            "model",
		},
		{
			SourcePath:      "templates/internal/domain/models/crud_mock.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "domain", "models", "mock", data.FileName()),
			Name:            "model_mock",
		},
		{
			SourcePath:      "templates/internal/domain/repositories/crud.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "domain", "repositories", data.FileName()),
			Name:            "repository",
		},
		{
			SourcePath:      "templates/internal/domain/usecases/crud.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "domain", "usecases", data.FileName()),
			Name:            "usecase",
		},
		{
			SourcePath:      "templates/internal/domain/interceptors/crud.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "domain", "interceptors", data.FileName()),
			Name:            "interceptor",
		},
		{
			SourcePath:      "templates/internal/usecases/crud.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "usecases", data.FileName()),
			Name:            "usecase",
		},
		{
			SourcePath:      "templates/internal/usecases/crud_test.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "usecases", data.TestFileName()),
			Name:            "usecase test",
		},
		{
			SourcePath:      "templates/internal/interceptors/crud.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "interceptors", data.FileName()),
			Name:            "interceptor",
		},
		{
			SourcePath:      "templates/internal/interceptors/crud_test.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "interceptors", data.TestFileName()),
			Name:            "interceptor test",
		},
		{
			SourcePath:      "templates/internal/repositories/crud.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "repositories", data.FileName()),
			Name:            "repository",
		},
		{
			SourcePath:      "templates/internal/repositories/crud_test.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "repositories", data.TestFileName()),
			Name:            "repository test",
		},
	}
	for _, tmpl := range files {
		if err := tmpl.renderToFile(data); err != nil {
			return err
		}
	}
	if err := addToDI("usecases", fmt.Sprintf("New%s", data.UseCaseTypeName())); err != nil {
		return err
	}
	if err := addToDI("interceptors", fmt.Sprintf("New%s", data.InterceptorTypeName())); err != nil {
		return err
	}
	if err := addToDI("repositories", fmt.Sprintf("New%s", data.RepositoryTypeName())); err != nil {
		return err
	}
	return nil
}
