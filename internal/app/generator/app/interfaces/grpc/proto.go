package grpc

import (
	"fmt"
	"github.com/mikalai-mitsin/creathor/internal/pkg/domain"
	"path"

	"github.com/mikalai-mitsin/creathor/internal/pkg/tmpl"
)

type Proto struct {
	domain *domain.Domain
}

func NewProto(domain *domain.Domain) *Proto {
	return &Proto{domain: domain}
}

const destinationPath = "."

func (c *Proto) Sync() error {
	files := []*tmpl.Template{
		{
			SourcePath: "templates/internal/domain/handlers/grpc/crud_test.go.tmpl",
			DestinationPath: path.Join(
				destinationPath,
				"internal",
				"app",
				c.domain.DirName(),
				"handlers",
				"grpc",
				c.domain.TestFileName(),
			),
			Name: "test grpc service server",
		},
	}
	for _, file := range files {
		if err := file.RenderToFile(c.domain); err != nil {
			return err
		}
	}
	proto := &tmpl.Template{
		SourcePath: "templates/api/proto/service/v1/crud.proto.tmpl",
		DestinationPath: path.Join(
			destinationPath,
			"api",
			"proto",
			c.domain.ProtoModule,
			"v1",
			fmt.Sprintf("%s.proto", c.domain.SnakeName()),
		),
		Name: "proto def",
	}
	if err := proto.RenderToFile(c.domain); err != nil {
		return err
	}
	return nil
}
