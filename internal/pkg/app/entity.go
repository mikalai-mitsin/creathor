package app

import (
	"fmt"

	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
)

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
	Params     []*configs.Param // FIXME: replace with own type
	Validation bool
	Mock       bool
}

func NewCreateEntity(entityConfig configs.EntityConfig) *Entity {
	return &Entity{
		Type:       EntityTypeCreate,
		Name:       entityConfig.CreateTypeName(),
		Variable:   "create",
		Params:     entityConfig.Params,
		Validation: true,
		Mock:       true,
	}
}

func NewUpdateEntity(entityConfig configs.EntityConfig) *Entity {
	model := &Entity{
		Type:     EntityTypeUpdate,
		Name:     entityConfig.UpdateTypeName(),
		Variable: "update",
		Params: []*configs.Param{
			{
				Name: "ID",
				Type: "uuid.UUID",
			},
		},
		Validation: true,
		Mock:       true,
	}
	for _, param := range entityConfig.Params {
		model.Params = append(model.Params, &configs.Param{
			Name: param.GetName(),
			Type: fmt.Sprintf("*%s", param.Type),
		})
	}
	return model
}

func NewMainEntity(modelConfig configs.EntityConfig) *Entity {
	model := &Entity{
		Type:     EntityTypeMain,
		Name:     modelConfig.EntityName(),
		Variable: modelConfig.Variable(),
		Params: []*configs.Param{
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

func NewFilterEntity(modelConfig configs.EntityConfig) *Entity {
	model := &Entity{
		Type:     EntityTypeFilter,
		Name:     modelConfig.FilterTypeName(),
		Variable: "filter",
		Params: []*configs.Param{
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
			{
				Name:   "IDs",
				Type:   "[]uuid.UUID",
				Search: false,
			},
		},
		Validation: true,
		Mock:       true,
	}
	return model
}
