package main

import (
	"encoding/json"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
	"os"
	"path"
	"strconv"
	"strings"
)

type Param struct {
	Name   string `yaml:"name"`
	Type   string `yaml:"type"`
	Search bool   `yaml:"search"`
}

func (n Param) Fake() string {
	typeName := strings.TrimPrefix(n.Type, "*")
	var fake string
	switch typeName {
	case "int":
		fake = "faker.RandomInt(2, 100)"
	case "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		fake = fmt.Sprintf("%s(faker.RandomInt(2, 100))", n)
	case "[]int":
		fake = fmt.Sprintf("[]%s{faker.RandomInt(2, 100), %s(faker.RandomInt(2, 100))}", n, n)
	case "[]int8", "[]int16", "[]int32", "[]int64", "[]uint", "[]uint8", "[]uint16", "[]uint32", "[]uint64":
		fake = fmt.Sprintf("[]%s{%s(faker.RandomInt(2, 100)), %s(faker.RandomInt(2, 100))}", n, n, n)
	case "string":
		fake = "faker.Lorem().String()"
	case "[]string":
		return "faker.Lorem().Words(5)"
	case "uuid":
		fake = "uuid.NewString()"
	default:
		return "/*FIXME*/"
	}
	if strings.HasPrefix(n.Type, "*") {
		return fmt.Sprintf("utils.Pointer(%s)", fake)
	}
	return fake
}

func (n Param) SQLType() string {
	typeName := strings.TrimPrefix(n.Type, "*")
	switch typeName {
	case "int8", "int16", "int32", "int":
		return "int"
	case "float32":
		return "float"
	case "float64":
		return "double precision"
	case "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return "bigint"
	case "[]int8", "[]int16", "[]int32", "[]int":
		return "int[]"
	case "[]int64", "[]uint", "[]uint8", "[]uint16", "[]uint32", "[]uint64":
		return "bigint[]"
	case "string":
		switch n.Name {
		case "title", "name":
			return "varchar"
		default:
			return "text"
		}
	case "[]string":
		return "varchar[]"
	case "uuid":
		return "uuid"
	case "time.Time":
		switch n.Name {
		case "date":
			return "date"
		case "time":
			return "time"
		default:
			return "timestamp"
		}
	case "time.Duration":
		return "interval"
	case "bool":
		return "boolean"
	default:
		return "/* FIXME */"
	}
}

func (n Param) Required() bool {
	return !strings.HasPrefix(n.Type, "*")
}

func (n Param) GetName() string {
	return strcase.ToCamel(n.Name)
}

func (n Param) Tag() string {
	return strcase.ToSnake(n.Name)
}

type Params []Param

type Model struct {
	Model  string `json:"model" yaml:"model"`
	Module string `json:"module" yaml:"module"`
	Auth   bool   `json:"auth" yaml:"auth"`
	Params Params `json:"params" yaml:"params"`
}

func (m Model) SearchEnabled() bool {
	for _, param := range m.Params {
		if param.Search {
			return true
		}
	}
	return false
}

func (m Model) SearchVector() string {
	var params []string
	for _, param := range m.Params {
		if param.Search {
			params = append(params, param.Tag())
		}
	}
	vector := fmt.Sprintf("to_tsvector('english', %s)", strings.Join(params, " || "))
	return vector
}

func ParseModel(s string) *Model {
	model := &Model{}
	if err := json.Unmarshal([]byte(s), model); err != nil {
		model = &Model{
			Model:  s,
			Module: "",
			Auth:   false,
			Params: Params{},
		}
	}
	return model
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

func (m Model) RESTHandlerTypeName() string {
	return fmt.Sprintf("%sHandler", strcase.ToCamel(m.Model))
}

func (m Model) RESTHandlerPath() string {
	return strcase.ToSnake(inflection.Plural(m.Model))
}

func (m Model) RESTHandlerVariableName() string {
	return fmt.Sprintf("%sHandler", strcase.ToLowerCamel(m.Model))
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

func (m Model) MigrationUpFileName() string {
	last, err := lastMigration()
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%06d_%s.up.sql", last+1, m.TableName())
}

func (m Model) MigrationDownFileName() string {
	last, err := lastMigration()
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%06d_%s.down.sql", last+1, m.TableName())
}

func lastMigration() (int, error) {
	dir, err := os.ReadDir(path.Join(destinationPath, "internal", "interfaces", "postgres", "migrations"))
	if err != nil {
		return 0, err
	}
	var files []string
	for _, entry := range dir {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	last := files[len(files)-1]
	n, _, _ := strings.Cut(strings.Trim(last, "0"), "_")
	index, err := strconv.Atoi(n)
	if err != nil {
		return 0, err
	}
	return index, nil
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

func (m Model) PermissionIDList() string {
	return fmt.Sprintf("PermissionID%sList", m.ModelName())
}

func (m Model) PermissionIDDetail() string {
	return fmt.Sprintf("PermissionID%sDetail", m.ModelName())
}

func (m Model) PermissionIDCreate() string {
	return fmt.Sprintf("PermissionID%sCreate", m.ModelName())
}

func (m Model) PermissionIDUpdate() string {
	return fmt.Sprintf("PermissionID%sUpdate", m.ModelName())
}

func (m Model) PermissionIDDelete() string {
	return fmt.Sprintf("PermissionID%sDelete", m.ModelName())
}
