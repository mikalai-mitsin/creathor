package grpc

import (
	"path"

	"github.com/018bf/creathor/internal/pkg/configs"
	"github.com/018bf/creathor/internal/pkg/tmpl"
)

var destinationPath = "."

type Handler struct {
	project *configs.Project
}

func NewHandler(project *configs.Project) *Handler {
	return &Handler{project: project}
}

func (c *Handler) Sync() error {
	files := []*tmpl.Template{
		{
			SourcePath: "templates/internal/auth/handlers/grpc/auth.go.tmpl",
			DestinationPath: path.Join(
				destinationPath,
				"internal",
				"app",
				"auth",
				"handlers",
				"grpc",
				"auth.go",
			),
			Name: "grpc auth",
		},
	}
	for _, file := range files {
		if err := file.RenderToFile(c.project); err != nil {
			return err
		}
	}
	return nil
}
