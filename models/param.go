package models

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/iancoleman/strcase"
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

func (p *Param) Fake() string {
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

func (p *Param) SQLType() string {
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

func (p *Param) GetGRPCWrapper() string {
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

func (p *Param) GetGRPCWrapperArgumentType() string {
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

func (p *Param) GRPCType() string {
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

func (p *Param) GRPCSliceType() string {
	return strings.TrimPrefix(p.GRPCType(), "[]")
}

func (p *Param) ProtoType() string {
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

func (p *Param) PostgresDTOType() string {
	switch p.Type {
	case "int8", "int16", "int32", "int":
		return "int"
	case "float32":
		return "float"
	case "byte":
		return "byte"
	case "[]byte":
		return "pq.ByteaArray"
	case "float64":
		return "float64"
	case "[]float32":
		return "pq.Float32Array"
	case "[]float64":
		return "pq.Float64Array"
	case "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return "int64"
	case "[]int8", "[]int16", "[]int32", "[]int":
		return "pq.Int32Array"
	case "[]int64", "[]uint", "[]uint8", "[]uint16", "[]uint32", "[]uint64":
		return "pq.Int64Array"
	case "string":
		return "string"
	case "[]string":
		return "pq.StringArray"
	case "uuid":
		return "string"
	case "time.Time":
		return "time.Time"
	case "time.Duration":
		return "time.Duration"
	case "bool":
		return "bool"
	case "[]bool":
		return "pq.BoolArray"
	default:
		return "/* FIXME */"
	}
}
func (p *Param) PostgresDTOSliceType() string {
	switch p.Type {
	case "[]byte":
		return "byte"
	case "[]float32":
		return "float32"
	case "[]float64":
		return "float64"
	case "[]int8", "[]int16", "[]int32", "[]int":
		return "int32"
	case "[]int64", "[]uint", "[]uint8", "[]uint16", "[]uint32", "[]uint64":
		return "int64"
	case "[]string":
		return "string"
	case "[]bool":
		return "bool"
	default:
		return "/* FIXME */"
	}
}

func (p *Param) GRPCGetter() string {
	return fmt.Sprintf("Get%s", p.GRPCParam())
}

func (p *Param) GRPCParam() string {
	return strings.ReplaceAll(strcase.ToCamel(p.Name), "ID", "Id")
}

func (p *Param) ProtoWrapType() string {
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

func (p *Param) GetName() string {
	return strcase.ToCamel(p.Name)
}

func (p *Param) Tag() string {
	return strcase.ToSnake(p.Name)
}
