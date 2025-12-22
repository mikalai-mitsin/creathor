package configs

import (
	"fmt"

	"github.com/iancoleman/strcase"
)

func (m *EntityConfig) ImportPathServices() string {
	return fmt.Sprintf(`"%s/internal/app/%s/services/%s"`, m.Module, m.AppName(), m.DirName())
}

func (m *EntityConfig) ImportAliasServices() string {
	return fmt.Sprintf("%sServices", m.LowerCamelName())
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

func (m *EntityConfig) ServiceTypeName() string {
	return fmt.Sprintf("%sService", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) ServiceVariableName() string {
	return fmt.Sprintf("%sService", strcase.ToLowerCamel(m.Name))
}

func (m *EntityConfig) EventServicePrivateVariableName() string {
	return fmt.Sprintf("%sEventService", m.LowerCamelName())
}

func (m *EntityConfig) EventServiceInterfaceName() string {
	return fmt.Sprintf("%sEventService", m.LowerCamelName())
}

func (m *EntityConfig) EventServiceName() string {
	return fmt.Sprintf("%sEventService", m.CamelName())
}

func (m *EntityConfig) EventServiceConstructorName() string {
	return fmt.Sprintf("New%s", m.EventServiceName())
}
