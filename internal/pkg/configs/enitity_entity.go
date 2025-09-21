package configs

import (
	"fmt"
	"slices"

	"github.com/iancoleman/strcase"
)

func (m *EntityConfig) ImportPathEntities() string {
	return fmt.Sprintf(`"%s/internal/app/%s/entities/%s"`, m.Module, m.AppName(), m.DirName())
}
func (m *EntityConfig) ImportAliasEntities() string {
	return fmt.Sprintf("%sEnitites", m.LowerCamelName())
}

func (m *EntityConfig) CreateTypeName() string {
	return fmt.Sprintf("%sCreate", m.CamelName())
}

func (m *EntityConfig) FilterTypeName() string {
	return fmt.Sprintf("%sFilter", m.CamelName())
}

func (m *EntityConfig) UpdateTypeName() string {
	return fmt.Sprintf("%sUpdate", m.CamelName())
}

func (m *EntityConfig) DeleteTypeName() string {
	return fmt.Sprintf("%sDelete", m.CamelName())
}

func (m *EntityConfig) OrderingTypeName() string {
	return fmt.Sprintf("%sOrdering", m.CamelName())
}

func (m *EntityConfig) OrderingConsts() map[string]string {
	consts := map[string]string{}
	for _, param := range m.GetMainModel().Params {
		consts[fmt.Sprintf("%s%sASC", m.OrderingTypeName(), strcase.ToCamel(param.Name))] = fmt.Sprintf(
			`"%s"`,
			param.Tag(),
		)
		consts[fmt.Sprintf("%s%sDESC", m.OrderingTypeName(), strcase.ToCamel(param.Name))] = fmt.Sprintf(
			`"-%s"`,
			param.Tag(),
		)
	}
	return consts
}

func (m *EntityConfig) OrderingMap() map[string]string {
	consts := map[string]string{}
	for _, param := range m.GetMainModel().Params {
		consts[fmt.Sprintf("%s%sASC", m.OrderingTypeName(), strcase.ToCamel(param.Name))] = fmt.Sprintf(
			`"%s.%s ASC"`,
			m.TableName(),
			param.Tag(),
		)
		consts[fmt.Sprintf("%s%sDESC", m.OrderingTypeName(), strcase.ToCamel(param.Name))] = fmt.Sprintf(
			`"%s.%s DESC"`,
			m.TableName(),
			param.Tag(),
		)
	}
	return consts
}

func (m *EntityConfig) MockFileName() string {
	return fmt.Sprintf("%s_mock.go", m.SnakeName())
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

func (m *EntityConfig) GetDeleteModel() *Entity {
	index := slices.IndexFunc(
		m.Entities,
		func(model *Entity) bool { return model.Type == EntityTypeDelete },
	)
	if index > 0 {
		return m.Entities[index]
	}
	return nil
}
