package main

import (
	"encoding/json"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
	"os"
	"path"
	"strconv"
	"strings"
)

type Param struct {
	Name   string `json:"name" yaml:"name"`
	Type   string `json:"type" yaml:"type"`
	Search bool   `json:"search" yaml:"search"`
}

func (p *Param) Validate() error {
	err := validation.ValidateStruct(
		p,
		validation.Field(&p.Name, validation.Required),
		validation.Field(&p.Type, validation.Required, validation.In(
			"int", "int64", "int32", "int16", "int8",
			"[]int", "[]int64", "[]int32", "[]int16", "[]int8",
			"uint", "uint64", "uint32", "uint16", "uint8",
			"[]uint", "[]uint64", "[]uint32", "[]uint16", "[]uint8",
			"string",
			"[]string",
			"time.Time",
			"[]time.Time",
		)),
	)
	if err != nil {
		return err
	}
	return nil
}

func (p *Param) IsSlice() bool {
	return strings.HasPrefix(p.Type, "[]")
}

func (p *Param) SliceType() string {
	return strings.TrimPrefix(p.Type, "[]")
}

func (p *Param) GrpcGetFromListValueAs() string {
	sliceType := p.SliceType()
	switch sliceType {
	case "int", "int32", "int64", "uint", "uint32", "uint64", "float32", "float64":
		return "GetNumberValue"
	case "string":
		return "GetStringValue"
	case "bool":
		return "GetBoolValue"
	case "map[string]interface{}", "map[string]any":
		return "GetStructValue"
	case "[]interface{}", "[]any":
		return "GetListValue"
	default:
		return "AsInterface"
	}
}

func (p Param) Fake() string {
	var fake string
	switch p.Type {
	case "int":
		fake = "faker.RandomInt(2, 100)"
	case "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		fake = fmt.Sprintf("%s(faker.RandomInt(2, 100))", p.Type)
	case "[]int":
		fake = "[]int{faker.RandomInt(2, 100), faker.RandomInt(2, 100)}"
	case "[]int8", "[]int16", "[]int32", "[]int64", "[]uint", "[]uint8", "[]uint16", "[]uint32", "[]uint64":
		fake = fmt.Sprintf("%s{%s(faker.RandomInt(2, 100)), %s(faker.RandomInt(2, 100))}", p.Type, p.SliceType(), p.SliceType())
	case "string":
		fake = "faker.Lorem().String()"
	case "[]string":
		return "faker.Lorem().Words(5)"
	case "uuid":
		fake = "uuid.NewString()"
	case "time.Time":
		fake = "faker.Time().Backward(40 * time.Hour).UTC()"
	default:
		return "/*FIXME*/"
	}
	return fake
}

func (p Param) SQLType() string {
	switch p.Type {
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
		switch p.Name {
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
		switch p.Name {
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

func (p Param) GetGRPCWrapper() string {
	switch p.Type {
	case "int", "int32", "int8", "int16":
		return "wrapperspb.Int32"
	case "int64":
		return "wrapperspb.Int64"
	case "uint8", "uint16", "uint32":
		return "wrapperspb.UInt32"
	case "uint64":
		return "wrapperspb.UInt64"
	case "string":
		return "wrapperspb.String"
	case "bool", "booleand":
		return "wrapperspb.Bool"
	case "float32":
		return "wrapperspb.Float"
	case "float64":
		return "wrapperspb.Double"
	case "time.Time":
		return "timestamppb.New"
	default:
		return "/* FIXME */"
	}
}

func (p Param) GetGRPCWrapperArgumentType() string {
	switch p.Type {
	case "int", "int32", "int8", "int16":
		return "int32"
	case "int64":
		return "int64"
	case "uint8", "uint16", "uint32":
		return "uint32"
	case "uint64":
		return "uint64"
	case "string":
		return "string"
	case "bool", "booleand":
		return "bool"
	case "float32":
		return "float32"
	case "float64":
		return "float64"
	case "time.Time":
		return "time.Time"
	default:
		return "/* FIXME */"
	}
}

func (p Param) GRPCType() string {
	switch p.Type {
	case "int8", "int16", "int32", "int":
		return "int32"
	case "int64":
		return "int64"
	case "float32":
		return "float"
	case "float64":
		return "double"
	case "uint", "uint8", "uint16", "uint32":
		return "uint32"
	case "uint64":
		return "uint64"
	case "[]int8", "[]int16", "[]int32", "[]int":
		return "[]int32"
	case "[]int64":
		return "[]int64"
	case "[]uint", "[]uint8", "[]uint16", "[]uint32":
		return "[]uint32"
	case "[]uint64":
		return "[]uint64"
	case "string", "uuid":
		return "string"
	case "[]string":
		return "[]string"
	case "time.Time":
		return "timestamppb.New"
	case "bool":
		return "bool"
	default:
		return "/* FIXME */"
	}
}

func (p Param) GRPCSliceType() string {
	return strings.TrimPrefix(p.GRPCType(), "[]")
}

func (p Param) ProtoType() string {
	switch p.Type {
	case "int8", "int16", "int32", "int":
		return "int32"
	case "int64":
		return "int64"
	case "float32":
		return "float"
	case "float64":
		return "double"
	case "uint", "uint8", "uint16", "uint32":
		return "uint32"
	case "uint64":
		return "uint64"
	case "[]int8", "[]int16", "[]int32", "[]int":
		return "repeated int32"
	case "[]int64":
		return "repeated int64"
	case "[]uint", "[]uint8", "[]uint16", "[]uint32":
		return "repeated uint32"
	case "[]uint64":
		return "repeated uint64"
	case "string", "uuid":
		return "string"
	case "[]string":
		return "repeated string"
	case "time.Time":
		return "google.protobuf.Timestamp"
	case "bool":
		return "bool"
	default:
		return "/* FIXME */"
	}
}

func (p Param) GRPCGetter() string {
	return fmt.Sprintf("Get%s", p.GRPCParam())
}

func (p Param) GRPCParam() string {
	return strings.ReplaceAll(strcase.ToCamel(p.Name), "ID", "Id")
}

func (p Param) ProtoWrapType() string {
	switch p.Type {
	case "int8", "int16", "int32", "int":
		return "google.protobuf.Int32Value"
	case "int64":
		return "google.protobuf.Int64Value"
	case "float32":
		return "google.protobuf.FloatValue"
	case "float64":
		return "google.protobuf.DoubleValue"
	case "uint", "uint8", "uint16", "uint32":
		return "google.protobuf.UInt32Value"
	case "uint64":
		return "google.protobuf.UInt64Value"
	case "[]int8", "[]int16", "[]int32", "[]int":
		return "google.protobuf.ListValue"
	case "[]int64":
		return "google.protobuf.ListValue"
	case "[]uint", "[]uint8", "[]uint16", "[]uint32":
		return "google.protobuf.ListValue"
	case "[]uint64":
		return "google.protobuf.ListValue"
	case "string", "uuid":
		return "google.protobuf.StringValue"
	case "[]string":
		return "google.protobuf.ListValue"
	case "time.Time":
		return "google.protobuf.Timestamp"
	case "bool":
		return "google.protobuf.BoolValue"
	case "[]bool":
		return "google.protobuf.ListValue"
	default:
		return "/* FIXME */"
	}
}

func (p Param) GetName() string {
	return strcase.ToCamel(p.Name)
}

func (p Param) Tag() string {
	return strcase.ToSnake(p.Name)
}

type Model struct {
	Model        string   `json:"model" yaml:"model"`
	Module       string   `json:"module" yaml:"module"`
	ProjectName  string   `json:"project_name" yaml:"projectName"`
	ProtoPackage string   `json:"proto_package" yaml:"protoPackage"`
	Auth         bool     `json:"auth" yaml:"auth"`
	Params       []*Param `json:"params" yaml:"params"`
}

func (m *Model) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.Model, validation.Required),
		validation.Field(&m.Module, validation.Required),
		validation.Field(&m.ProjectName, validation.Required),
		validation.Field(&m.Auth),
		validation.Field(&m.Params),
	)
	if err != nil {
		return err
	}
	return nil
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
			Params: nil,
		}
	}
	return model
}

func (m Model) Variable() string {
	return strcase.ToLowerCamel(m.Model)
}

func (m Model) ListVariable() string {
	return strcase.ToLowerCamel(fmt.Sprintf("list%s", strcase.ToCamel(inflection.Plural(m.Model))))
}

func (m Model) ModelName() string {
	return strcase.ToCamel(m.Model)
}

func (m Model) UseCaseTypeName() string {
	return fmt.Sprintf("%sUseCase", strcase.ToCamel(m.Model))
}

func (m Model) GRPCHandlerTypeName() string {
	return fmt.Sprintf("%sServiceServer", strcase.ToCamel(m.Model))
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

func (m Model) ProtoFileName() string {
	return fmt.Sprintf("%s.proto", m.SnakeName())
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
