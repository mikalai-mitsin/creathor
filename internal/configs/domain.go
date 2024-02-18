package configs

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
)

type DomainConfig struct {
	Model          string   `json:"model"         yaml:"model"`
	Module         string   `json:"module"        yaml:"module"`
	ProjectName    string   `json:"project_name"  yaml:"projectName"`
	ProtoPackage   string   `json:"proto_package" yaml:"protoPackage"`
	Auth           bool     `json:"auth"          yaml:"auth"`
	Params         []*Param `json:"params"        yaml:"params"`
	GRPCEnabled    bool     `                     yaml:"gRPC"`
	GatewayEnabled bool     `                     yaml:"gateway"`
	KafkaEnabled   bool     `                     yaml:"kafka"`
}

func (m *DomainConfig) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.Model, validation.Required),
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

func (m *DomainConfig) SearchEnabled() bool {
	for _, param := range m.Params {
		if param.Search {
			return true
		}
	}
	return false
}

func (m *DomainConfig) SearchVector() string {
	var params []string
	for _, param := range m.Params {
		if param.Search {
			params = append(params, param.Tag())
		}
	}
	vector := fmt.Sprintf("to_tsvector('english', %s)", strings.Join(params, " || "))
	return vector
}

func (m *DomainConfig) Variable() string {
	return strcase.ToLowerCamel(m.Model)
}

func (m *DomainConfig) ListVariable() string {
	return strcase.ToLowerCamel(fmt.Sprintf("list%s", strcase.ToCamel(inflection.Plural(m.Model))))
}

func (m *DomainConfig) ModelName() string {
	return strcase.ToCamel(m.Model)
}

func (m *DomainConfig) DomainName() string {
	return strcase.ToSnake(m.Model)
}

func (m *DomainConfig) DomainAlias() string {
	return strcase.ToLowerCamel(m.Model)
}

func (m *DomainConfig) UseCaseTypeName() string {
	return fmt.Sprintf("%sUseCase", strcase.ToCamel(m.Model))
}

func (m *DomainConfig) GRPCHandlerTypeName() string {
	return fmt.Sprintf("%sServiceServer", strcase.ToCamel(m.Model))
}

func (m *DomainConfig) RESTHandlerTypeName() string {
	return fmt.Sprintf("%sHandler", strcase.ToCamel(m.Model))
}

func (m *DomainConfig) GatewayHandlerTypeName() string {
	return fmt.Sprintf("Register%sServiceHandlerFromEndpoint", strcase.ToCamel(m.Model))
}

func (m *DomainConfig) RESTHandlerPath() string {
	return strcase.ToSnake(inflection.Plural(m.Model))
}

func (m *DomainConfig) RESTHandlerVariableName() string {
	return fmt.Sprintf("%sHandler", strcase.ToLowerCamel(m.Model))
}

func (m *DomainConfig) GRPCHandlerVariableName() string {
	return fmt.Sprintf("%sHandler", strcase.ToLowerCamel(m.Model))
}

func (m *DomainConfig) UseCaseVariableName() string {
	return fmt.Sprintf("%sUseCase", strcase.ToLowerCamel(m.Model))
}

func (m *DomainConfig) InterceptorTypeName() string {
	return fmt.Sprintf("%sInterceptor", strcase.ToCamel(m.Model))
}

func (m *DomainConfig) InterceptorVariableName() string {
	return fmt.Sprintf("%sInterceptor", strcase.ToLowerCamel(m.Model))
}

func (m *DomainConfig) RepositoryTypeName() string {
	return fmt.Sprintf("%sRepository", strcase.ToCamel(m.Model))
}

func (m *DomainConfig) RepositoryVariableName() string {
	return fmt.Sprintf("%sRepository", strcase.ToLowerCamel(m.Model))
}

func (m *DomainConfig) FilterTypeName() string {
	return fmt.Sprintf("%sFilter", strcase.ToCamel(m.Model))
}

func (m *DomainConfig) UpdateTypeName() string {
	return fmt.Sprintf("%sUpdate", strcase.ToCamel(m.Model))
}

func (m *DomainConfig) CreateTypeName() string {
	return fmt.Sprintf("%sCreate", strcase.ToCamel(m.Model))
}

func (m *DomainConfig) PostgresDTOTypeName() string {
	return fmt.Sprintf("%sDTO", strcase.ToCamel(m.Model))
}

func (m *DomainConfig) PostgresDTOListTypeName() string {
	return fmt.Sprintf("%sListDTO", strcase.ToCamel(m.Model))
}

func (m *DomainConfig) KeyName() string {
	return strcase.ToSnake(m.Model)
}

func (m *DomainConfig) SnakeName() string {
	return strcase.ToSnake(m.Model)
}

func (m *DomainConfig) FileName() string {
	return fmt.Sprintf("%s.go", m.SnakeName())
}

func (m *DomainConfig) MigrationUpFileName() string {
	last, err := lastMigration()
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%06d_%s.up.sql", last+1, m.TableName())
}

func (m *DomainConfig) MigrationDownFileName() string {
	last, err := lastMigration()
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%06d_%s.down.sql", last+1, m.TableName())
}

func lastMigration() (int, error) {
	dir, err := os.ReadDir(path.Join("internal", "interfaces", "postgres", "migrations"))
	if err != nil {
		return 0, err
	}
	var files []string
	for _, entry := range dir {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	last := files[len(files)-1]
	n, _, _ := strings.Cut(strings.Trim(last, "0"), "_")
	index, err := strconv.Atoi(n)
	if err != nil {
		return 0, err
	}
	return index, nil
}

func (m *DomainConfig) TestFileName() string {
	return fmt.Sprintf("%s_test.go", m.SnakeName())
}

func (m *DomainConfig) ProtoFileName() string {
	return fmt.Sprintf("%s.proto", m.SnakeName())
}

func (m *DomainConfig) MockFileName() string {
	return fmt.Sprintf("%s_mock.go", m.SnakeName())
}

func (m *DomainConfig) TableName() string {
	return strcase.ToSnake(inflection.Plural(m.Model))
}

func (m *DomainConfig) PermissionIDList() string {
	return fmt.Sprintf("PermissionID%sList", m.ModelName())
}

func (m *DomainConfig) PermissionIDDetail() string {
	return fmt.Sprintf("PermissionID%sDetail", m.ModelName())
}

func (m *DomainConfig) PermissionIDCreate() string {
	return fmt.Sprintf("PermissionID%sCreate", m.ModelName())
}

func (m *DomainConfig) PermissionIDUpdate() string {
	return fmt.Sprintf("PermissionID%sUpdate", m.ModelName())
}

func (m *DomainConfig) PermissionIDDelete() string {
	return fmt.Sprintf("PermissionID%sDelete", m.ModelName())
}
