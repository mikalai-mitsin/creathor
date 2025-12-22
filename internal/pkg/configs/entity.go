package configs

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
)

type EntityConfig struct {
	Name         string   `json:"name"          yaml:"name"`
	Module       string   `json:"module"        yaml:"module"`
	ProjectName  string   `json:"project_name"  yaml:"projectName"`
	ProtoPackage string   `json:"proto_package" yaml:"protoPackage"`
	Params       []*Param `json:"params"        yaml:"params"`
	HTTPEnabled  bool     `                     yaml:"http"`
	GRPCEnabled  bool     `                     yaml:"gRPC"`
	KafkaEnabled bool     `                     yaml:"kafka"`
	AppConfig    *AppConfig
	Entities     []*Entity
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

func (m *EntityConfig) CamelName() string {
	return strcase.ToCamel(m.Name)
}

func (m *EntityConfig) AppName() string {
	return m.AppConfig.AppName()
}

func (m *EntityConfig) AppAlias() string {
	return strcase.ToLowerCamel(m.Name)
}

func (m *EntityConfig) KeyName() string {
	return strcase.ToSnake(m.Name)
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

func (m *EntityConfig) LowerCamelName() string {
	return strcase.ToLowerCamel(m.Name)
}

func (m *EntityConfig) DirName() string {
	return strcase.ToSnake(m.Name)
}

func (m *EntityConfig) GetOneVariableName() string {
	return strcase.ToLowerCamel(m.Name)
}

func (m *EntityConfig) GetManyVariableName() string {
	return inflection.Plural(m.GetOneVariableName())
}

type EntityType uint8

const (
	EntityTypeMain = iota
	EntityTypeCreate
	EntityTypeUpdate
	EntityTypeFilter
	EntityTypeDelete
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
		Name:     modelConfig.CamelName(),
		Variable: modelConfig.GetOneVariableName(),
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
			{
				Name:   "DeletedAt",
				Type:   "*time.Time",
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
				Type:   fmt.Sprintf("[]%s", modelConfig.OrderingTypeName()),
				Search: false,
			},
			{
				Name:   "IsDeleted",
				Type:   "*bool",
				Search: false,
			},
		},
		Validation: true,
		Mock:       true,
	}
	return model
}

func NewDeleteEntity(entityConfig EntityConfig) *Entity {
	model := &Entity{
		Type:     EntityTypeDelete,
		Name:     entityConfig.DeleteTypeName(),
		Variable: "del",
		Params: []*Param{
			{
				Name: "ID",
				Type: "uuid.UUID",
			},
		},
		Validation: true,
		Mock:       true,
	}
	return model
}
