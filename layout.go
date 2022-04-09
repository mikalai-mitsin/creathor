package main

import (
	"os"
	"os/exec"
	"path"
)

type Project struct {
	Name      string
	Module    string
	GoVersion string
}

func CreateLayout(data *Project) error {
	directories := []string{
		path.Join(destinationPath, "build"),
		path.Join(destinationPath, "cmd"),
		path.Join(destinationPath, "cmd", data.Name),
		path.Join(destinationPath, "configs"),
		path.Join(destinationPath, "deployments"),
		path.Join(destinationPath, "dist"),
		path.Join(destinationPath, "docs"),
		path.Join(destinationPath, "internal"),
		path.Join(destinationPath, "internal", "configs"),
		path.Join(destinationPath, "internal", "containers"),
		path.Join(destinationPath, "internal", "domain", "errs"),
		path.Join(destinationPath, "internal", "domain", "interceptors"),
		path.Join(destinationPath, "internal", "domain", "models"),
		path.Join(destinationPath, "internal", "domain", "repositories"),
		path.Join(destinationPath, "internal", "domain", "usecases"),
		path.Join(destinationPath, "internal", "interceptors"),
		path.Join(destinationPath, "internal", "interfaces"),
		path.Join(destinationPath, "internal", "usecases"),
		path.Join(destinationPath, "internal", "repositories"),
		path.Join(destinationPath, "pkg"),
		path.Join(destinationPath, "pkg", "clock"),
		path.Join(destinationPath, "pkg", "log"),
	}
	for _, directory := range directories {
		if err := os.MkdirAll(directory, 0777); err != nil {
			return NewUnexpectedBehaviorError(err.Error())
		}
	}
	files := []*Template{
		{
			SourcePath:      "templates/cmd/service/main.go.tmpl",
			DestinationPath: path.Join(destinationPath, "cmd", data.Name, "main.go"),
			Name:            "service main",
		},
		{
			SourcePath:      "templates/configs/config.toml.tmpl",
			DestinationPath: path.Join(destinationPath, "configs", "config.toml"),
			Name:            "main config",
		},
		{
			SourcePath:      "templates/configs/config.toml.tmpl",
			DestinationPath: path.Join(destinationPath, "configs", "ci.toml"),
			Name:            "ci config",
		},
		{
			SourcePath:      "templates/configs/config.toml.tmpl",
			DestinationPath: path.Join(destinationPath, "configs", "test.toml"),
			Name:            "test config",
		},
		{
			SourcePath:      "templates/configs/config.go.tmpl",
			DestinationPath: path.Join(destinationPath, "internal", "configs", "config.go"),
			Name:            "config struct",
		},
		{
			SourcePath:      "templates/configs/config_test.go.tmpl",
			DestinationPath: path.Join(destinationPath, "internal", "configs", "config_test.go"),
			Name:            "config tests",
		},
		{
			SourcePath:      "templates/errs/errors.go.tmpl",
			DestinationPath: path.Join(destinationPath, "internal", "domain", "errs", "errors.go"),
			Name:            "domain errors",
		},
		{
			SourcePath:      "templates/errs/errors_test.go.tmpl",
			DestinationPath: path.Join(destinationPath, "internal", "domain", "errs", "errors_test.go"),
			Name:            "domain errors tests",
		},
		{
			SourcePath:      "templates/clock/clock.go.tmpl",
			DestinationPath: path.Join(destinationPath, "pkg", "clock", "clock.go"),
			Name:            "clock",
		},
		{
			SourcePath:      "templates/log/logger.go.tmpl",
			DestinationPath: path.Join(destinationPath, "pkg", "log", "logger.go"),
			Name:            "logger interface",
		},
		{
			SourcePath:      "templates/go.mod.tmpl",
			DestinationPath: path.Join(destinationPath, "go.mod"),
			Name:            "go.mod",
		},
	}
	for _, file := range files {
		if err := file.renderToFile(data); err != nil {
			return err
		}
	}
	command := exec.Command("go", "mod", "tidy")
	if err := command.Run(); err != nil {
		return err
	}
	return nil
}
