package grpc

import (
	"fmt"
	"path"

	"github.com/mikalai-mitsin/creathor/internal/pkg/app"

	"github.com/mikalai-mitsin/creathor/internal/pkg/tmpl"
)

type ProtoGenerator struct {
	domain *app.BaseEntity
}

func NewProtoGenerator(domain *app.BaseEntity) *ProtoGenerator {
	return &ProtoGenerator{domain: domain}
}

const destinationPath = "."

func (c *ProtoGenerator) Sync() error {
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
