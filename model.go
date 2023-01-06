package main

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
)

type Model struct {
	Model  string
	Module string
	Auth   bool
}

func (m Model) Variable() string {
	return strcase.ToLowerCamel(m.Model)
}

func (m Model) ListVariable() string {
	return strcase.ToLowerCamel(inflection.Plural(m.Model))
}

func (m Model) ModelName() string {
	return strcase.ToCamel(m.Model)
}

func (m Model) UseCaseTypeName() string {
	return fmt.Sprintf("%sUseCase", strcase.ToCamel(m.Model))
}

func (m Model) UseCaseVariableName() string {
	return fmt.Sprintf("%sUseCase", strcase.ToLowerCamel(m.Model))
}

func (m Model) InterceptorTypeName() string {
	return fmt.Sprintf("%sInterceptor", strcase.ToCamel(m.Model))
}

func (m Model) InterceptorVariableName() string {
	return fmt.Sprintf("%sInterceptor", strcase.ToLowerCamel(m.Model))
}

func (m Model) RepositoryTypeName() string {
	return fmt.Sprintf("%sRepository", strcase.ToCamel(m.Model))
}

func (m Model) RepositoryVariableName() string {
	return fmt.Sprintf("%sRepository", strcase.ToLowerCamel(m.Model))
}

func (m Model) FilterTypeName() string {
	return fmt.Sprintf("%sFilter", strcase.ToCamel(m.Model))
}

func (m Model) FilterVariableName() string {
	return fmt.Sprintf("%sFilter", strcase.ToLowerCamel(m.Model))
}

func (m Model) UpdateTypeName() string {
	return fmt.Sprintf("%sUpdate", strcase.ToCamel(m.Model))
}

func (m Model) UpdateVariableName() string {
	return fmt.Sprintf("%sUpdate", strcase.ToLowerCamel(m.Model))
}

func (m Model) CreateTypeName() string {
	return fmt.Sprintf("%sCreate", strcase.ToCamel(m.Model))
}

func (m Model) CreateVariableName() string {
	return fmt.Sprintf("%sCreate", strcase.ToLowerCamel(m.Model))
}

func (m Model) KeyName() string {
	return strcase.ToSnake(m.Model)
}

func (m Model) SnakeName() string {
	return strcase.ToSnake(m.Model)
}

func (m Model) FileName() string {
	return fmt.Sprintf("%s.go", m.SnakeName())
}

func (m Model) TestFileName() string {
	return fmt.Sprintf("%s_test.go", m.SnakeName())
}

func (m Model) MockFileName() string {
	return fmt.Sprintf("%s_mock.go", m.SnakeName())
}

func (m Model) TableName() string {
	return strcase.ToSnake(inflection.Plural(m.Model))
}
