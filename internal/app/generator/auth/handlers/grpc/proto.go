package grpc

import (
	"path"

	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
	"github.com/mikalai-mitsin/creathor/internal/pkg/tmpl"
)

type Proto struct {
	project *configs.Project
}

func NewProto(project *configs.Project) *Proto {
	return &Proto{project: project}
}

func (c *Proto) Sync() error {
	files := []*tmpl.Template{
		{
			SourcePath: "templates/api/proto/service/v1/auth.proto.tmpl",
			DestinationPath: path.Join(
				destinationPath,
				"api",
				"proto",
				c.project.ProtoPackage(),
				"v1",
				"auth.proto",
			),
			Name: "auth.proto",
		},
	}
	for _, file := range files {
		if err := file.RenderToFile(c.project); err != nil {
			return err
		}
	}
	return nil
}
