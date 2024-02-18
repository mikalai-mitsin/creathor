package domain

import (
	"fmt"

	"github.com/018bf/creathor/internal/configs"
)

type ModelType uint8

const (
	ModelTypeMain = iota
	ModelTypeCreate
	ModelTypeUpdate
	ModelTypeFilter
)

type Model struct {
	Type       ModelType
	Name       string
	Variable   string
	Params     []*configs.Param // FIXME: replace with own type
	Validation bool
	Mock       bool
}

func NewCreateModel(modelConfig *configs.DomainConfig) *Model {
	return &Model{
		Type:       ModelTypeCreate,
		Name:       modelConfig.CreateTypeName(),
		Variable:   "create",
		Params:     modelConfig.Params,
		Validation: true,
		Mock:       true,
	}
}

func NewUpdateModel(modelConfig *configs.DomainConfig) *Model {
	model := &Model{
		Type:     ModelTypeUpdate,
		Name:     modelConfig.UpdateTypeName(),
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
	for _, param := range modelConfig.Params {
		model.Params = append(model.Params, &configs.Param{
			Name: param.GetName(),
			Type: fmt.Sprintf("*%s", param.Type),
		})
	}
	return model
}

func NewMainModel(modelConfig *configs.DomainConfig) *Model {
	model := &Model{
		Type:     ModelTypeMain,
		Name:     modelConfig.ModelName(),
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

func NewFilterModel(modelConfig *configs.DomainConfig) *Model {
	model := &Model{
		Type:     ModelTypeFilter,
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
