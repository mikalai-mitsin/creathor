package main

import (
	"github.com/018bf/creathor/models"
	"os"
	"path/filepath"
)

func CreateCI(project *models.Project) error {
	var directories []string
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
	switch project.CI {
	case "gitlab":
		files = append(files, &Template{
			SourcePath:      "templates/ci/gitlab/gitlab-ci.yml.tmpl",
			DestinationPath: filepath.Join(destinationPath, ".gitlab-ci.yml"),
			Name:            "gitlab-ci",
		})
	case "github":
		files = append(files, &Template{
			SourcePath:      "templates/ci/github/workflows/tests.yaml.tmpl",
			DestinationPath: filepath.Join(destinationPath, ".github", "workflows", "tests.yaml"),
			Name:            "gitlab-ci",
		})
		directories = append(directories, filepath.Join(destinationPath, ".github", "workflows"))
	}
	for _, directory := range directories {
		if err := os.MkdirAll(directory, 0777); err != nil {
			return NewUnexpectedBehaviorError(err.Error())
		}
	}
	for _, tmpl := range files {
		if err := tmpl.renderToFile(project); err != nil {
			return err
		}
	}
	return nil
}
