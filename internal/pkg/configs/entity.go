package configs

import (
	"fmt"
	"os"
	"path"
	"slices"
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
)

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
	AppConfig      *AppConfig
	Entities       []*Entity
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
	return m.AppConfig.AppName()
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

func (m *EntityConfig) ProtoFileName() string {
	return fmt.Sprintf("%s.proto", m.SnakeName())
}

func (m *EntityConfig) MockFileName() string {
	return fmt.Sprintf("%s_mock.go", m.SnakeName())
}

func (m *EntityConfig) SnakeName() string {
	return strcase.ToSnake(m.Name)
}

func (m *EntityConfig) FileName() string {
	return fmt.Sprintf("%s.go", m.SnakeName())
}

func (m *EntityConfig) TestFileName() string {
	return fmt.Sprintf("%s_test.go", m.SnakeName())
}

func (m *EntityConfig) MigrationUpFileName() string {
	last, err := lastMigration()
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%06d_%s.up.sql", last+1, m.TableName())
}

func (m *EntityConfig) MigrationDownFileName() string {
	last, err := lastMigration()
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%06d_%s.down.sql", last+1, m.TableName())
}

func (m *EntityConfig) CamelName() string {
	return strcase.ToCamel(m.Name)
}

func (m *EntityConfig) LowerCamelName() string {
	return strcase.ToLowerCamel(m.Name)
}

func (m *EntityConfig) DirName() string {
	return strcase.ToSnake(m.Name)
}

func (m *EntityConfig) EventProducerConstructorName() string {
	return fmt.Sprintf("New%s", m.EventProducerTypeName())
}

func (m *EntityConfig) EventProducerTypeName() string {
	return fmt.Sprintf("%sEventProducer", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) EventProducerInterfaceName() string {
	return fmt.Sprintf("%sEventProducer", strcase.ToLowerCamel(m.Name))
}
func (m *EntityConfig) GetEventProducerPrivateVariableName() string {
	return fmt.Sprintf("%sEventProducer", strcase.ToLowerCamel(m.Name))
}

func (m *EntityConfig) CreatedTopicName() string {
	return fmt.Sprintf("%s.created", strcase.ToSnake(m.Name))
}
func (m *EntityConfig) UpdatedTopicName() string {
	return fmt.Sprintf("%s.updated", strcase.ToSnake(m.Name))
}
func (m *EntityConfig) DeletedTopicName() string {
	return fmt.Sprintf("%s.deleted", strcase.ToSnake(m.Name))
}

func (m *EntityConfig) EntitiesImportPath() string {
	return fmt.Sprintf(`"%s/internal/app/%s/entities/%s"`, m.Module, m.AppName(), m.DirName())
}

func (m *EntityConfig) GetMainModel() *Entity {
	index := slices.IndexFunc(
		m.Entities,
		func(model *Entity) bool { return model.Type == EntityTypeMain },
	)
	if index >= 0 {
		return m.Entities[index]
	}
	return nil
}

func (m *EntityConfig) TableName() string {
	return strcase.ToSnake(inflection.Plural(m.Name))
}

func (m *EntityConfig) SearchEnabled() bool {
	return slices.ContainsFunc(
		m.GetMainModel().Params,
		func(param *Param) bool { return param.Search },
	)
}

func (m *EntityConfig) GetCreateModel() *Entity {
	index := slices.IndexFunc(
		m.Entities,
		func(model *Entity) bool { return model.Type == EntityTypeCreate },
	)
	if index >= 0 {
		return m.Entities[index]
	}
	return nil
}

func (m *EntityConfig) GetUpdateModel() *Entity {
	index := slices.IndexFunc(
		m.Entities,
		func(model *Entity) bool { return model.Type == EntityTypeUpdate },
	)
	if index > 0 {
		return m.Entities[index]
	}
	return nil
}

func (m *EntityConfig) GetFilterModel() *Entity {
	index := slices.IndexFunc(
		m.Entities,
		func(model *Entity) bool { return model.Type == EntityTypeFilter },
	)
	if index > 0 {
		return m.Entities[index]
	}
	return nil
}

func (m *EntityConfig) PermissionIDCreate() string {
	return fmt.Sprintf("PermissionID%sCreate", strcase.ToCamel(m.CamelName()))
}

func (m *EntityConfig) PermissionIDUpdate() string {
	return fmt.Sprintf("PermissionID%sUpdate", m.CamelName())
}

func (m *EntityConfig) PermissionIDDelete() string {
	return fmt.Sprintf("PermissionID%sDelete", m.CamelName())
}

func (m *EntityConfig) PermissionIDDetail() string {
	return fmt.Sprintf("PermissionID%sDetail", m.CamelName())
}

func (m *EntityConfig) PermissionIDList() string {
	return fmt.Sprintf("PermissionID%sList", m.CamelName())
}

func (m *EntityConfig) GetOneVariableName() string {
	return strcase.ToLowerCamel(m.Name)
}

func (m *EntityConfig) GetManyVariableName() string {
	return inflection.Plural(m.GetOneVariableName())
}

func (m *EntityConfig) GetHTTPPath() string {
	return strcase.ToSnake(inflection.Plural(m.GetOneVariableName()))
}

func (m *EntityConfig) GetGRPCHandlerPrivateVariableName() string {
	return fmt.Sprintf("grpc%sHandler", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) GetGRPCHandlerPublicVariableName() string {
	return fmt.Sprintf("%sHandler", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) GetGRPCHandlerTypeName() string {
	return fmt.Sprintf("%sServiceServer", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) GetGRPCHandlerConstructorName() string {
	return fmt.Sprintf("New%s", m.GetGRPCHandlerTypeName())
}

func (m *EntityConfig) GetGRPCServiceDescriptionName() string {
	return fmt.Sprintf("%sService_ServiceDesc", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) GetGRPCCreateDTOEncodeName() string {
	return fmt.Sprintf("encode%s", m.GetCreateModel().Name)
}

func (m *EntityConfig) GetGRPCUpdateDTOEncodeName() string {
	return fmt.Sprintf("encode%s", m.GetUpdateModel().Name)
}

func (m *EntityConfig) GetGRPCFilterDTOEncodeName() string {
	return fmt.Sprintf("encode%s", m.GetFilterModel().Name)
}

func (m *EntityConfig) GetGRPCMainDecodeName() string {
	return fmt.Sprintf("decode%s", m.GetMainModel().Name)
}
func (m *EntityConfig) GetGRPCMainListDecodeName() string {
	return fmt.Sprintf("decodeList%s", m.GetMainModel().Name)
}
func (m *EntityConfig) GetGRPCUpdateDecodeName() string {
	return fmt.Sprintf("decode%s", m.GetUpdateModel().Name)
}

func (m *EntityConfig) GetHTTPHandlerConstructorName() string {
	return fmt.Sprintf("New%s", m.GetHTTPHandlerTypeName())
}

func (m *EntityConfig) GetHTTPHandlerTypeName() string {
	return fmt.Sprintf("%sHandler", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) GetHTTPHandlerPrivateVariableName() string {
	return fmt.Sprintf("http%sHandler", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) GetHTTPItemDTOName() string {
	return fmt.Sprintf("%sDTO", strcase.ToCamel(m.GetMainModel().Name))
}
func (m *EntityConfig) GetHTTPItemDTOConstructorName() string {
	return fmt.Sprintf("New%s", m.GetHTTPItemDTOName())
}

func (m *EntityConfig) GetHTTPUpdateDTOName() string {
	return fmt.Sprintf("%sDTO", strcase.ToCamel(m.GetUpdateModel().Name))
}

func (m *EntityConfig) GetHTTPUpdateDTOConstructorName() string {
	return fmt.Sprintf("New%s", m.GetHTTPUpdateDTOName())
}

func (m *EntityConfig) GetHTTPCreateDTOName() string {
	return fmt.Sprintf("%sDTO", strcase.ToCamel(m.GetCreateModel().Name))
}
func (m *EntityConfig) GetHTTPCreateDTOConstructorName() string {
	return fmt.Sprintf("New%s", m.GetHTTPCreateDTOName())
}

func (m *EntityConfig) GetHTTPListDTOName() string {
	return fmt.Sprintf("%sListDTO", strcase.ToCamel(m.GetMainModel().Name))
}

func (m *EntityConfig) GetHTTPListDTOConstructorName() string {
	return fmt.Sprintf("New%s", strcase.ToCamel(m.GetHTTPListDTOName()))
}

func (m *EntityConfig) GetUseCasePrivateVariableName() string {
	return fmt.Sprintf("%sUseCase", strcase.ToLowerCamel(m.Name))
}

func (m *EntityConfig) GetUseCasePublicVariableName() string {
	return fmt.Sprintf("%sUseCase", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) GetUseCaseTypeName() string {
	return fmt.Sprintf("%sUseCase", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) GetUseCaseInterfaceName() string {
	return fmt.Sprintf("%sUseCase", strcase.ToLowerCamel(m.Name))
}

func (m *EntityConfig) GetUseCaseConstructorName() string {
	return fmt.Sprintf("New%s", m.GetUseCaseTypeName())
}

func (m *EntityConfig) GetServicePrivateVariableName() string {
	return fmt.Sprintf("%sService", strcase.ToLowerCamel(m.Name))
}

func (m *EntityConfig) GetServicePublicVariableName() string {
	return fmt.Sprintf("%sService", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) GetServiceTypeName() string {
	return fmt.Sprintf("%sService", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) GetServiceInterfaceName() string {
	return fmt.Sprintf("%sService", strcase.ToLowerCamel(m.Name))
}

func (m *EntityConfig) GetServiceConstructorName() string {
	return fmt.Sprintf("New%s", m.GetServiceTypeName())
}

func (m *EntityConfig) GetRepositoryPrivateVariableName() string {
	return fmt.Sprintf("%sRepository", strcase.ToLowerCamel(m.Name))
}

func (m *EntityConfig) GetRepositoryPublicVariableName() string {
	return fmt.Sprintf("%sRepository", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) GetRepositoryTypeName() string {
	return fmt.Sprintf("%sRepository", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) GetRepositoryInterfaceName() string {
	return fmt.Sprintf("%sRepository", strcase.ToLowerCamel(m.Name))
}

func (m *EntityConfig) GetRepositoryConstructorName() string {
	return fmt.Sprintf("New%s", m.GetRepositoryTypeName())
}

func (m *EntityConfig) GetHTTPFilterDTOName() string {
	return fmt.Sprintf("%sFilterDTO", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) GetHTTPFilterDTOConstructorName() string {
	return fmt.Sprintf("New%s", m.GetHTTPFilterDTOName())
}

func lastMigration() (int, error) {
	dir, err := os.ReadDir(path.Join("internal", "pkg", "postgres", "migrations"))
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

type EntityType uint8

const (
	EntityTypeMain = iota
	EntityTypeCreate
	EntityTypeUpdate
	EntityTypeFilter
)

type Entity struct {
	Type       EntityType
	Name       string
	Variable   string
	Params     []*Param // FIXME: replace with own type
	Validation bool
	Mock       bool
}

func NewCreateEntity(entityConfig EntityConfig) *Entity {
	return &Entity{
		Type:       EntityTypeCreate,
		Name:       entityConfig.CreateTypeName(),
		Variable:   "create",
		Params:     entityConfig.Params,
		Validation: true,
		Mock:       true,
	}
}

func NewUpdateEntity(entityConfig EntityConfig) *Entity {
	model := &Entity{
		Type:     EntityTypeUpdate,
		Name:     entityConfig.UpdateTypeName(),
		Variable: "update",
		Params: []*Param{
			{
				Name: "ID",
				Type: "uuid.UUID",
			},
		},
		Validation: true,
		Mock:       true,
	}
	for _, param := range entityConfig.Params {
		model.Params = append(model.Params, &Param{
			Name: param.GetName(),
			Type: fmt.Sprintf("*%s", param.Type),
		})
	}
	return model
}

func NewMainEntity(modelConfig EntityConfig) *Entity {
	model := &Entity{
		Type:     EntityTypeMain,
		Name:     modelConfig.EntityName(),
		Variable: modelConfig.Variable(),
		Params: []*Param{
			{
				Name:   "ID",
				Type:   "uuid.UUID",
				Search: false,
			},
			{
				Name:   "CreatedAt",
				Type:   "time.Time",
				Search: false,
			},
			{
				Name:   "UpdatedAt",
				Type:   "time.Time",
				Search: false,
			},
		},
		Validation: true,
		Mock:       true,
	}
	model.Params = append(model.Params, modelConfig.Params...)
	return model
}

func NewFilterEntity(modelConfig EntityConfig) *Entity {
	model := &Entity{
		Type:     EntityTypeFilter,
		Name:     modelConfig.FilterTypeName(),
		Variable: "filter",
		Params: []*Param{
			{
				Name:   "PageSize",
				Type:   "*uint64",
				Search: false,
			},
			{
				Name:   "PageNumber",
				Type:   "*uint64",
				Search: false,
			},
			{
				Name:   "Search",
				Type:   "*string",
				Search: false,
			},
			{
				Name:   "OrderBy",
				Type:   "[]string",
				Search: false,
			},
		},
		Validation: true,
		Mock:       true,
	}
	return model
}
