package mods

import (
	"fmt"

	"github.com/018bf/creathor/internal/configs"

	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
	"golang.org/x/exp/slices"
)

type Mod struct {
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

func (m *Mod) GetMainModel() *Model {
	index := slices.IndexFunc(
		m.Models,
		func(model *Model) bool { return model.Type == ModelTypeMain },
	)
	if index >= 0 {
		return m.Models[index]
	}
	return nil
}

func (m *Mod) TableName() string {
	return strcase.ToSnake(inflection.Plural(m.Name))
}

func (m *Mod) SearchEnabled() bool {
	return slices.ContainsFunc(
		m.GetMainModel().Params,
		func(param *configs.Param) bool { return param.Search },
	)
}

func (m *Mod) GetCreateModel() *Model {
	index := slices.IndexFunc(
		m.Models,
		func(model *Model) bool { return model.Type == ModelTypeCreate },
	)
	if index >= 0 {
		return m.Models[index]
	}
	return nil
}

func (m *Mod) GetUpdateModel() *Model {
	index := slices.IndexFunc(
		m.Models,
		func(model *Model) bool { return model.Type == ModelTypeUpdate },
	)
	if index > 0 {
		return m.Models[index]
	}
	return nil
}

func (m *Mod) GetFilterModel() *Model {
	index := slices.IndexFunc(
		m.Models,
		func(model *Model) bool { return model.Type == ModelTypeFilter },
	)
	if index > 0 {
		return m.Models[index]
	}
	return nil
}

func (m *Mod) PermissionIDCreate() string {
	return fmt.Sprintf("PermissionID%sCreate", strcase.ToCamel(m.Name))
}

func (m *Mod) PermissionIDUpdate() string {
	return fmt.Sprintf("PermissionID%sUpdate", strcase.ToCamel(m.Name))
}

func (m *Mod) PermissionIDDelete() string {
	return fmt.Sprintf("PermissionID%sDelete", strcase.ToCamel(m.Name))
}

func (m *Mod) PermissionIDDetail() string {
	return fmt.Sprintf("PermissionID%sDetail", strcase.ToCamel(m.Name))
}

func (m *Mod) PermissionIDList() string {
	return fmt.Sprintf("PermissionID%sList", strcase.ToCamel(m.Name))
}
