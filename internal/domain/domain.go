package domain

import (
	"fmt"

	"github.com/018bf/creathor/internal/configs"

	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
	"golang.org/x/exp/slices"
)

type Domain struct {
	Name        string
	Module      string
	ProtoModule string
	Filename    string
	Models      []*Model
	UseCase     *Layer
	Repository  *Layer
	Interceptor *Layer
	GRPCHandler *Layer
	Auth        bool
}

func (m *Domain) GetMainModel() *Model {
	index := slices.IndexFunc(
		m.Models,
		func(model *Model) bool { return model.Type == ModelTypeMain },
	)
	if index >= 0 {
		return m.Models[index]
	}
	return nil
}

func (m *Domain) TableName() string {
	return strcase.ToSnake(inflection.Plural(m.Name))
}

func (m *Domain) SearchEnabled() bool {
	return slices.ContainsFunc(
		m.GetMainModel().Params,
		func(param *configs.Param) bool { return param.Search },
	)
}

func (m *Domain) GetCreateModel() *Model {
	index := slices.IndexFunc(
		m.Models,
		func(model *Model) bool { return model.Type == ModelTypeCreate },
	)
	if index >= 0 {
		return m.Models[index]
	}
	return nil
}

func (m *Domain) GetUpdateModel() *Model {
	index := slices.IndexFunc(
		m.Models,
		func(model *Model) bool { return model.Type == ModelTypeUpdate },
	)
	if index > 0 {
		return m.Models[index]
	}
	return nil
}

func (m *Domain) GetFilterModel() *Model {
	index := slices.IndexFunc(
		m.Models,
		func(model *Model) bool { return model.Type == ModelTypeFilter },
	)
	if index > 0 {
		return m.Models[index]
	}
	return nil
}

func (m *Domain) PermissionIDCreate() string {
	return fmt.Sprintf("PermissionID%sCreate", strcase.ToCamel(m.Name))
}

func (m *Domain) PermissionIDUpdate() string {
	return fmt.Sprintf("PermissionID%sUpdate", strcase.ToCamel(m.Name))
}

func (m *Domain) PermissionIDDelete() string {
	return fmt.Sprintf("PermissionID%sDelete", strcase.ToCamel(m.Name))
}

func (m *Domain) PermissionIDDetail() string {
	return fmt.Sprintf("PermissionID%sDetail", strcase.ToCamel(m.Name))
}

func (m *Domain) PermissionIDList() string {
	return fmt.Sprintf("PermissionID%sList", strcase.ToCamel(m.Name))
}
