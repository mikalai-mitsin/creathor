package configs

import (
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
)

type AppConfig struct {
	Name           string   `json:"name" yaml:"name"`
	Module         string   `json:"module"        yaml:"module"`
	ProjectName    string   `json:"project_name"  yaml:"projectName"`
	ProtoPackage   string   `json:"proto_package" yaml:"protoPackage"`
	Auth           bool     `json:"auth"          yaml:"auth"`
	Params         []*Param `json:"params"        yaml:"params"`
	HTTPEnabled    bool     `                     yaml:"http"`
	GRPCEnabled    bool     `                     yaml:"gRPC"`
	GatewayEnabled bool     `                     yaml:"gateway"`
	KafkaEnabled   bool     `                     yaml:"kafka"`
}

func (m *AppConfig) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.Name, validation.Required),
		validation.Field(&m.Module, validation.Required),
		validation.Field(&m.ProjectName, validation.Required),
		validation.Field(&m.Auth),
		validation.Field(&m.Params),
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *AppConfig) SearchEnabled() bool {
	for _, param := range m.Params {
		if param.Search {
			return true
		}
	}
	return false
}

func (m *AppConfig) SearchVector() string {
	var params []string
	for _, param := range m.Params {
		if param.Search {
			params = append(params, param.Tag())
		}
	}
	vector := fmt.Sprintf("to_tsvector('english', %s)", strings.Join(params, " || "))
	return vector
}

func (m *AppConfig) Variable() string {
	return strcase.ToLowerCamel(m.Name)
}

func (m *AppConfig) ListVariable() string {
	return strcase.ToLowerCamel(fmt.Sprintf("list%s", strcase.ToCamel(inflection.Plural(m.Name))))
}

func (m *AppConfig) EntityName() string {
	return strcase.ToCamel(m.Name)
}

func (m *AppConfig) AppName() string {
	return strcase.ToSnake(m.Name)
}

func (m *AppConfig) AppAlias() string {
	return strcase.ToLowerCamel(m.Name)
}

func (m *AppConfig) CamelCase() string {
	return strcase.ToCamel(m.Name)
}

func (m *AppConfig) ServiceTypeName() string {
	return fmt.Sprintf("%sService", strcase.ToCamel(m.Name))
}

func (m *AppConfig) GRPCHandlerTypeName() string {
	return fmt.Sprintf("%sServiceServer", strcase.ToCamel(m.Name))
}

func (m *AppConfig) RESTHandlerTypeName() string {
	return fmt.Sprintf("%sHandler", strcase.ToCamel(m.Name))
}

func (m *AppConfig) GatewayHandlerTypeName() string {
	return fmt.Sprintf("Register%sServiceHandlerFromEndpoint", strcase.ToCamel(m.Name))
}

func (m *AppConfig) RESTHandlerPath() string {
	return strcase.ToSnake(inflection.Plural(m.Name))
}

func (m *AppConfig) RESTHandlerVariableName() string {
	return fmt.Sprintf("%sHandler", strcase.ToLowerCamel(m.Name))
}

func (m *AppConfig) GRPCHandlerVariableName() string {
	return fmt.Sprintf("%sHandler", strcase.ToLowerCamel(m.Name))
}

func (m *AppConfig) ServiceVariableName() string {
	return fmt.Sprintf("%sService", strcase.ToLowerCamel(m.Name))
}

func (m *AppConfig) UseCaseTypeName() string {
	return fmt.Sprintf("%sUseCase", strcase.ToCamel(m.Name))
}

func (m *AppConfig) UseCaseVariableName() string {
	return fmt.Sprintf("%sUseCase", strcase.ToLowerCamel(m.Name))
}

func (m *AppConfig) RepositoryTypeName() string {
	return fmt.Sprintf("%sRepository", strcase.ToCamel(m.Name))
}

func (m *AppConfig) RepositoryVariableName() string {
	return fmt.Sprintf("%sRepository", strcase.ToLowerCamel(m.Name))
}

func (m *AppConfig) FilterTypeName() string {
	return fmt.Sprintf("%sFilter", strcase.ToCamel(m.Name))
}

func (m *AppConfig) UpdateTypeName() string {
	return fmt.Sprintf("%sUpdate", strcase.ToCamel(m.Name))
}

func (m *AppConfig) CreateTypeName() string {
	return fmt.Sprintf("%sCreate", strcase.ToCamel(m.Name))
}

func (m *AppConfig) PostgresDTOTypeName() string {
	return fmt.Sprintf("%sDTO", strcase.ToCamel(m.Name))
}

func (m *AppConfig) PostgresDTOListTypeName() string {
	return fmt.Sprintf("%sListDTO", strcase.ToCamel(m.Name))
}

func (m *AppConfig) KeyName() string {
	return strcase.ToSnake(m.Name)
}

func (m *AppConfig) SnakeName() string {
	return strcase.ToSnake(m.Name)
}

func (m *AppConfig) ProtoFileName() string {
	return fmt.Sprintf("%s.proto", m.SnakeName())
}

func (m *AppConfig) MockFileName() string {
	return fmt.Sprintf("%s_mock.go", m.SnakeName())
}

func (m *AppConfig) TableName() string {
	return strcase.ToSnake(inflection.Plural(m.Name))
}

func (m *AppConfig) PermissionIDList() string {
	return fmt.Sprintf("PermissionID%sList", m.EntityName())
}

func (m *AppConfig) PermissionIDDetail() string {
	return fmt.Sprintf("PermissionID%sDetail", m.EntityName())
}

func (m *AppConfig) PermissionIDCreate() string {
	return fmt.Sprintf("PermissionID%sCreate", m.EntityName())
}

func (m *AppConfig) PermissionIDUpdate() string {
	return fmt.Sprintf("PermissionID%sUpdate", m.EntityName())
}

func (m *AppConfig) PermissionIDDelete() string {
	return fmt.Sprintf("PermissionID%sDelete", m.EntityName())
}
