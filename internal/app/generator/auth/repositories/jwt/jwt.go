package jwt

import (
	"path"

	"github.com/018bf/creathor/internal/pkg/configs"
	"github.com/018bf/creathor/internal/pkg/tmpl"
)

var destinationPath = "."

type Repository struct {
	project *configs.Project
}

func NewRepository(project *configs.Project) *Repository {
	return &Repository{project: project}
}

func (c *Repository) Sync() error {
	files := []*tmpl.Template{
		{
			SourcePath: "templates/internal/auth/repositories/jwt/auth.go.tmpl",
			DestinationPath: path.Join(
				destinationPath,
				"internal",
				"app",
				"auth",
				"repositories",
				"jwt",
				"auth.go",
			),
			Name: "jwt auth repository implementation",
		},
		{
			SourcePath: "templates/internal/auth/repositories/jwt/auth_test.go.tmpl",
			DestinationPath: path.Join(
				destinationPath,
				"internal",
				"app",
				"auth",
				"repositories",
				"jwt",
				"auth_test.go",
			),
			Name: "test auth repository implementation",
		},
	}
	for _, file := range files {
		if err := file.RenderToFile(c.project); err != nil {
			return err
		}
	}
	return nil
}
