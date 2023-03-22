package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/018bf/creathor/internal/configs"
	"github.com/iancoleman/strcase"
)

func CreateDeployment(data *configs.Project) error {
	directories := []string{
		path.Join(destinationPath, "deployments", "helm_vars", "staging"),
		path.Join(destinationPath, "deployments", "helm_vars", "development"),
		path.Join(destinationPath, "deployments", "helm_vars", "production"),
	}
	for _, directory := range directories {
		if err := os.MkdirAll(directory, 0777); err != nil {
			return NewUnexpectedBehaviorError(err.Error())
		}
	}
	helmCreate := exec.Command("helm", "create", strcase.ToKebab(data.Name))
	helmCreate.Dir = "deployments"
	var errb bytes.Buffer
	helmCreate.Stderr = &errb
	if err := helmCreate.Run(); err != nil {
		fmt.Println(strings.Join(helmCreate.Args, " "))
		fmt.Println(errb.String())
	}
	files := []*Template{
		{
			SourcePath: "templates/deployments/helm_vars/development/values.yaml.tmpl",
			DestinationPath: filepath.Join(
				destinationPath,
				"deployments",
				"helm_vars",
				"development",
				"values.yaml",
			),
			Name: "development values",
		},
		{
			SourcePath: "templates/deployments/helm_vars/staging/values.yaml.tmpl",
			DestinationPath: filepath.Join(
				destinationPath,
				"deployments",
				"helm_vars",
				"staging",
				"values.yaml",
			),
			Name: "staging values",
		},
		{
			SourcePath: "templates/deployments/helm_vars/production/values.yaml.tmpl",
			DestinationPath: filepath.Join(
				destinationPath,
				"deployments",
				"helm_vars",
				"production",
				"values.yaml",
			),
			Name: "production values",
		},
	}
	for _, tmpl := range files {
		if err := tmpl.renderToFile(nil); err != nil {
			return err
		}
	}
	return nil
}
