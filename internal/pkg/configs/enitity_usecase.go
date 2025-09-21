package configs

import (
	"fmt"

	"github.com/iancoleman/strcase"
)

func (m *EntityConfig) ImportPathUseCases() string {
	return fmt.Sprintf(`"%s/internal/app/%s/usecases/%s"`, m.Module, m.AppName(), m.DirName())
}

func (m *EntityConfig) ImportAliasUseCases() string {
	return fmt.Sprintf("%sUseCases", m.LowerCamelName())
}

func (m *EntityConfig) UseCaseTypeName() string {
	return fmt.Sprintf("%sUseCase", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) UseCaseVariableName() string {
	return fmt.Sprintf("%sUseCase", strcase.ToLowerCamel(m.Name))
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
