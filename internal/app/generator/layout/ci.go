package layout

import (
	"os"
	"path/filepath"

	"github.com/018bf/creathor/internal/pkg/errs"
	"github.com/018bf/creathor/internal/pkg/tmpl"

	"github.com/018bf/creathor/internal/pkg/configs"
)

type CIGenerator struct {
	project *configs.Project
}

func NewCIGenerator(project *configs.Project) *CIGenerator {
	return &CIGenerator{project: project}
}

func (c *CIGenerator) Sync() error {
	var directories []string
	files := []*tmpl.Template{
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
	switch c.project.CI {
	case "gitlab":
		files = append(files, &tmpl.Template{
			SourcePath:      "templates/ci/gitlab/gitlab-ci.yml.tmpl",
			DestinationPath: filepath.Join(destinationPath, ".gitlab-ci.yml"),
			Name:            "gitlab-ci",
		})
	case "github":
		files = append(files, &tmpl.Template{
			SourcePath:      "templates/ci/github/workflows/tests.yaml.tmpl",
			DestinationPath: filepath.Join(destinationPath, ".github", "workflows", "tests.yaml"),
			Name:            "gitlab-ci",
		})
		directories = append(directories, filepath.Join(destinationPath, ".github", "workflows"))
	}
	for _, directory := range directories {
		if err := os.MkdirAll(directory, 0777); err != nil {
			return errs.NewUnexpectedBehaviorError(err.Error())
		}
	}
	for _, file := range files {
		if err := file.RenderToFile(c.project); err != nil {
			return err
		}
	}
	return nil
}
