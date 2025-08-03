package configs

import (
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
)

type AppConfig struct {
	Name           string         `json:"name" yaml:"name"`
	Module         string         `json:"module"        yaml:"module"`
	ProjectName    string         `json:"project_name"  yaml:"projectName"`
	ProtoPackage   string         `json:"proto_package" yaml:"protoPackage"`
	HTTPEnabled    bool           `                     yaml:"http"`
	GRPCEnabled    bool           `                     yaml:"gRPC"`
	GatewayEnabled bool           `                     yaml:"gateway"`
	KafkaEnabled   bool           `                     yaml:"kafka"`
	Entities       []EntityConfig `json:"entities" yaml:"entities"`
}

func (m *AppConfig) AppName() string {
	return strcase.ToSnake(m.Name)
}

func (m *AppConfig) AppAlias() string {
	return strcase.ToLowerCamel(m.Name)
}

type EntityConfig struct {
	Name           string   `json:"name" yaml:"name"`
	Module         string   `json:"module"        yaml:"module"`
	ProjectName    string   `json:"project_name"  yaml:"projectName"`
	ProtoPackage   string   `json:"proto_package" yaml:"protoPackage"`
	Params         []*Param `json:"params"        yaml:"params"`
	HTTPEnabled    bool     `                     yaml:"http"`
	GRPCEnabled    bool     `                     yaml:"gRPC"`
	GatewayEnabled bool     `                     yaml:"gateway"`
	KafkaEnabled   bool     `                     yaml:"kafka"`
}

func (m *EntityConfig) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.Name, validation.Required),
		validation.Field(&m.Module, validation.Required),
		validation.Field(&m.ProjectName, validation.Required),
		validation.Field(&m.Params),
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *EntityConfig) SearchEnabled() bool {
	for _, param := range m.Params {
		if param.Search {
			return true
		}
	}
	return false
}

func (m *EntityConfig) SearchVector() string {
	var params []string
	for _, param := range m.Params {
		if param.Search {
			params = append(params, param.Tag())
		}
	}
	vector := fmt.Sprintf("to_tsvector('english', %s)", strings.Join(params, " || "))
	return vector
}

func (m *EntityConfig) Variable() string {
	return strcase.ToLowerCamel(m.Name)
}

func (m *EntityConfig) ListVariable() string {
	return strcase.ToLowerCamel(fmt.Sprintf("list%s", strcase.ToCamel(inflection.Plural(m.Name))))
}

func (m *EntityConfig) EntityName() string {
	return strcase.ToCamel(m.Name)
}

func (m *EntityConfig) AppName() string {
	return strcase.ToSnake(m.Name)
}

func (m *EntityConfig) AppAlias() string {
	return strcase.ToLowerCamel(m.Name)
}

func (m *EntityConfig) CamelCase() string {
	return strcase.ToCamel(m.Name)
}

func (m *EntityConfig) ServiceTypeName() string {
	return fmt.Sprintf("%sService", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) GRPCHandlerTypeName() string {
	return fmt.Sprintf("%sServiceServer", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) RESTHandlerTypeName() string {
	return fmt.Sprintf("%sHandler", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) GatewayHandlerTypeName() string {
	return fmt.Sprintf("Register%sServiceHandlerFromEndpoint", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) RESTHandlerPath() string {
	return strcase.ToSnake(inflection.Plural(m.Name))
}

func (m *EntityConfig) RESTHandlerVariableName() string {
	return fmt.Sprintf("%sHandler", strcase.ToLowerCamel(m.Name))
}

func (m *EntityConfig) GRPCHandlerVariableName() string {
	return fmt.Sprintf("%sHandler", strcase.ToLowerCamel(m.Name))
}

func (m *EntityConfig) ServiceVariableName() string {
	return fmt.Sprintf("%sService", strcase.ToLowerCamel(m.Name))
}

func (m *EntityConfig) UseCaseTypeName() string {
	return fmt.Sprintf("%sUseCase", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) UseCaseVariableName() string {
	return fmt.Sprintf("%sUseCase", strcase.ToLowerCamel(m.Name))
}

func (m *EntityConfig) RepositoryTypeName() string {
	return fmt.Sprintf("%sRepository", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) RepositoryVariableName() string {
	return fmt.Sprintf("%sRepository", strcase.ToLowerCamel(m.Name))
}

func (m *EntityConfig) FilterTypeName() string {
	return fmt.Sprintf("%sFilter", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) UpdateTypeName() string {
	return fmt.Sprintf("%sUpdate", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) CreateTypeName() string {
	return fmt.Sprintf("%sCreate", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) PostgresDTOTypeName() string {
	return fmt.Sprintf("%sDTO", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) PostgresDTOListTypeName() string {
	return fmt.Sprintf("%sListDTO", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) KeyName() string {
	return strcase.ToSnake(m.Name)
}

func (m *EntityConfig) SnakeName() string {
	return strcase.ToSnake(m.Name)
}

func (m *EntityConfig) ProtoFileName() string {
	return fmt.Sprintf("%s.proto", m.SnakeName())
}

func (m *EntityConfig) MockFileName() string {
	return fmt.Sprintf("%s_mock.go", m.SnakeName())
}

func (m *EntityConfig) TableName() string {
	return strcase.ToSnake(inflection.Plural(m.Name))
}

func (m *EntityConfig) PermissionIDList() string {
	return fmt.Sprintf("PermissionID%sList", m.EntityName())
}

func (m *EntityConfig) PermissionIDDetail() string {
	return fmt.Sprintf("PermissionID%sDetail", m.EntityName())
}

func (m *EntityConfig) PermissionIDCreate() string {
	return fmt.Sprintf("PermissionID%sCreate", m.EntityName())
}

func (m *EntityConfig) PermissionIDUpdate() string {
	return fmt.Sprintf("PermissionID%sUpdate", m.EntityName())
}

func (m *EntityConfig) PermissionIDDelete() string {
	return fmt.Sprintf("PermissionID%sDelete", m.EntityName())
}
