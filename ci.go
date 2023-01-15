package main

import (
	"path/filepath"
)

func CreateCI() error {
	files := []*Template{
		{
			SourcePath:      "templates/ci/golangci.yml.tmpl",
			DestinationPath: filepath.Join(destinationPath, ".golangci.yml"),
			Name:            "golangci-lint",
		},
		{
			SourcePath:      "templates/ci/pre-commit-config.yaml.tmpl",
			DestinationPath: filepath.Join(destinationPath, ".pre-commit-config.yaml"),
			Name:            "pre-commit",
		},
	}
	switch ci {
	case "gitlab":
		files = append(files, &Template{
			SourcePath:      "templates/ci/gitlab-ci.yml.tmpl",
			DestinationPath: filepath.Join(destinationPath, ".gitlab-ci.yml"),
			Name:            "gitlab-ci",
		})
	}
	for _, tmpl := range files {
		if err := tmpl.renderToFile(nil); err != nil {
			return err
		}
	}
	return nil
}
