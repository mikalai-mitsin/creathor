package configs

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
)

func (m *EntityConfig) ImportPathHttpHandlers() string {
	return fmt.Sprintf(`"%s/internal/app/%s/handlers/http/%s"`, m.Module, m.AppName(), m.DirName())
}

func (m *EntityConfig) ImportAliasHttpHandlers() string {
	return fmt.Sprintf("%sHttpHandlers", m.LowerCamelName())
}

func (m *EntityConfig) GetHTTPPath() string {
	return strcase.ToSnake(inflection.Plural(m.GetOneVariableName()))
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
func (m *EntityConfig) GetHTTPDeleteDTOName() string {
	return fmt.Sprintf("%sDTO", strcase.ToCamel(m.GetDeleteModel().Name))
}

func (m *EntityConfig) GetHTTPUpdateDTOConstructorName() string {
	return fmt.Sprintf("New%s", m.GetHTTPUpdateDTOName())
}

func (m *EntityConfig) GetHTTPDeleteDTOConstructorName() string {
	return fmt.Sprintf("New%s", m.GetHTTPDeleteDTOName())
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

func (m *EntityConfig) GetHTTPFilterDTOName() string {
	return fmt.Sprintf("%sFilterDTO", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) GetHTTPFilterDTOConstructorName() string {
	return fmt.Sprintf("New%s", m.GetHTTPFilterDTOName())
}

func (m *EntityConfig) RESTHandlerTypeName() string {
	return fmt.Sprintf("%sHandler", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) RESTHandlerPath() string {
	return strcase.ToSnake(inflection.Plural(m.Name))
}

func (m *EntityConfig) RESTHandlerVariableName() string {
	return fmt.Sprintf("%sHandler", strcase.ToLowerCamel(m.Name))
}
