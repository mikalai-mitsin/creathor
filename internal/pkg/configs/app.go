package configs

import (
	"github.com/iancoleman/strcase"
)

type AppConfig struct {
	Name          string         `json:"name" yaml:"name"`
	Module        string         `json:"module"        yaml:"module"`
	ProjectName   string         `json:"project_name"  yaml:"projectName"`
	ProtoPackage  string         `json:"proto_package" yaml:"protoPackage"`
	HTTPEnabled   bool           `                     yaml:"http"`
	GRPCEnabled   bool           `                     yaml:"gRPC"`
	KafkaEnabled  bool           `                     yaml:"kafka"`
	Entities      []EntityConfig `json:"entities" yaml:"entities"`
	ProjectConfig *Project       `json:"-" yaml:"-"`
}

func (m *AppConfig) AppName() string {
	return strcase.ToSnake(m.Name)
}

func (m *AppConfig) AppAlias() string {
	return strcase.ToLowerCamel(m.Name)
}
