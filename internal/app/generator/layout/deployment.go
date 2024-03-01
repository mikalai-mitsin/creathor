package layout

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/018bf/creathor/internal/pkg/errs"
	"github.com/018bf/creathor/internal/pkg/tmpl"
	"github.com/iancoleman/strcase"

	"github.com/018bf/creathor/internal/pkg/configs"
)

var destinationPath = "."

type DeploymentGenerator struct {
	project *configs.Project
}

func NewDeploymentGenerator(project *configs.Project) *DeploymentGenerator {
	return &DeploymentGenerator{project: project}
}

func (d *DeploymentGenerator) Sync() error {
	directories := []string{
		path.Join(destinationPath, "deployments", "helm_vars", "staging"),
		path.Join(destinationPath, "deployments", "helm_vars", "development"),
		path.Join(destinationPath, "deployments", "helm_vars", "production"),
	}
	for _, directory := range directories {
		if err := os.MkdirAll(directory, 0777); err != nil {
			return errs.NewUnexpectedBehaviorError(err.Error())
		}
	}
	helmCreate := exec.Command("helm", "create", strcase.ToKebab(d.project.Name))
	helmCreate.Dir = "deployments"
	var errb bytes.Buffer
	helmCreate.Stderr = &errb
	if err := helmCreate.Run(); err != nil {
		fmt.Println(strings.Join(helmCreate.Args, " "))
		fmt.Println(errb.String())
	}
	files := []*tmpl.Template{
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
	for _, file := range files {
		if err := file.RenderToFile(nil); err != nil {
			return err
		}
	}
	return nil
}
