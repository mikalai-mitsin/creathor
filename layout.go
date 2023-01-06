package main

import (
	"os"
	"path"
)

type Project struct {
	Name      string
	Module    string
	GoVersion string
	Auth      bool
}

func CreateLayout(data *Project) error {
	directories := []string{
		path.Join(destinationPath, "build"),
		path.Join(destinationPath, "cmd"),
		path.Join(destinationPath, "cmd", data.Name),
		path.Join(destinationPath, "configs"),
		path.Join(destinationPath, "dist"),
		path.Join(destinationPath, "docs"),
		path.Join(destinationPath, "docs", ".chglog"),
		path.Join(destinationPath, "internal"),
		path.Join(destinationPath, "internal", "configs"),
		path.Join(destinationPath, "internal", "domain", "errs"),
		path.Join(destinationPath, "internal", "domain", "interceptors"),
		path.Join(destinationPath, "internal", "domain", "models"),
		path.Join(destinationPath, "internal", "domain", "models", "mock"),
		path.Join(destinationPath, "internal", "domain", "repositories"),
		path.Join(destinationPath, "internal", "domain", "usecases"),
		path.Join(destinationPath, "internal", "interceptors"),
		path.Join(destinationPath, "internal", "interfaces"),
		path.Join(destinationPath, "internal", "interfaces", "postgres"),
		path.Join(destinationPath, "internal", "interfaces", "postgres", "migrations"),
		path.Join(destinationPath, "internal", "usecases"),
		path.Join(destinationPath, "internal", "repositories"),
		path.Join(destinationPath, "pkg"),
		path.Join(destinationPath, "pkg", "clock"),
		path.Join(destinationPath, "pkg", "log"),
		path.Join(destinationPath, "pkg", "utils"),
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
			SourcePath:      "templates/utils/pointer.go.tmpl",
			DestinationPath: path.Join(destinationPath, "pkg", "utils", "pointer.go"),
			Name:            "utils pointer",
		},
		{
			SourcePath:      "templates/go.mod.tmpl",
			DestinationPath: path.Join(destinationPath, "go.mod"),
			Name:            "go.mod",
		},
		{
			SourcePath:      "templates/docs/README.md.tmpl",
			DestinationPath: path.Join(destinationPath, "README.md"),
			Name:            "README.md",
		},
		{
			SourcePath:      "templates/docs/chglog/CHANGELOG.tpl.md.tmpl",
			DestinationPath: path.Join(destinationPath, "docs", ".chglog", "CHANGELOG.tpl.md"),
			Name:            ".chglog/CHANGELOG.tpl.md",
		},
		{
			SourcePath:      "templates/docs/chglog/config.yml.tmpl",
			DestinationPath: path.Join(destinationPath, "docs", ".chglog", "config.yml"),
			Name:            ".chglog/config.yml",
		},
		{
			SourcePath:      "templates/docs/CHANGELOG.md.tmpl",
			DestinationPath: path.Join(destinationPath, "docs", "CHANGELOG.md"),
			Name:            "CHANGELOG.md",
		},
		{
			SourcePath:      "templates/interfaces/postgres/fx.go.tmpl",
			DestinationPath: path.Join(destinationPath, "internal", "interfaces", "postgres", "fx.go"),
			Name:            "postgres fx",
		},
		{
			SourcePath:      "templates/interfaces/postgres/postgres.go.tmpl",
			DestinationPath: path.Join(destinationPath, "internal", "interfaces", "postgres", "postgres.go"),
			Name:            "postgres",
		},
		{
			SourcePath:      "templates/interfaces/postgres/testing.go.tmpl",
			DestinationPath: path.Join(destinationPath, "internal", "interfaces", "postgres", "testing.go"),
			Name:            "postgres testing",
		},
		{
			SourcePath:      "templates/interfaces/postgres/migrations/init.sql.tmpl",
			DestinationPath: path.Join(destinationPath, "internal", "interfaces", "postgres", "migrations", "init.sql"),
			Name:            "postgres init migration",
		},
	}
	if authEnabled {
		files = append(
			files,
			&Template{
				SourcePath:      "templates/domain/usecase_auth.go.tmpl",
				DestinationPath: path.Join(destinationPath, "internal", "domain", "usecases", "auth.go"),
				Name:            "auth usecase",
			},
			&Template{
				SourcePath:      "templates/domain/model_auth.go.tmpl",
				DestinationPath: path.Join(destinationPath, "internal", "domain", "models", "auth.go"),
				Name:            "auth models",
			},
			&Template{
				SourcePath:      "templates/domain/model_permission.go.tmpl",
				DestinationPath: path.Join(destinationPath, "internal", "domain", "models", "permission.go"),
				Name:            "auth permission",
			},
		)
	}
	for _, file := range files {
		if err := file.renderToFile(data); err != nil {
			return err
		}
	}
	return nil
}
