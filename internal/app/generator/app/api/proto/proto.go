package proto

import (
	"fmt"
	"path"

	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"

	"github.com/mikalai-mitsin/creathor/internal/pkg/tmpl"
)

type ProtoGenerator struct {
	entityConfig configs.EntityConfig
}

func NewProtoGenerator(entityConfig configs.EntityConfig) *ProtoGenerator {
	return &ProtoGenerator{entityConfig: entityConfig}
}

const destinationPath = "."

func (c *ProtoGenerator) Sync() error {
	proto := &tmpl.Template{
		SourcePath: "templates/api/proto/service/v1/crud.proto.tmpl",
		DestinationPath: path.Join(
			destinationPath,
			"api",
			"proto",
			c.entityConfig.ProtoPackage,
			"v1",
			fmt.Sprintf("%s.proto", c.entityConfig.SnakeName()),
		),
		Name: "proto def",
	}
	if err := proto.RenderToFile(&c.entityConfig); err != nil {
		return err
	}
	return nil
}
