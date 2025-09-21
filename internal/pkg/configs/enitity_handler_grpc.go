package configs

import (
	"fmt"

	"github.com/iancoleman/strcase"
)

func (m *EntityConfig) ImportPathGrpcHandlers() string {
	return fmt.Sprintf(`"%s/internal/app/%s/handlers/grpc/%s"`, m.Module, m.AppName(), m.DirName())
}

func (m *EntityConfig) ImportAliasGrpcHandlers() string {
	return fmt.Sprintf("%sGrpcHandlers", m.LowerCamelName())
}

func (m *EntityConfig) GetGRPCHandlerPrivateVariableName() string {
	return fmt.Sprintf("grpc%sHandler", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) GetGRPCHandlerPublicVariableName() string {
	return fmt.Sprintf("%sHandler", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) GetGRPCHandlerTypeName() string {
	return fmt.Sprintf("%sServiceServer", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) GetGRPCHandlerConstructorName() string {
	return fmt.Sprintf("New%s", m.GetGRPCHandlerTypeName())
}

func (m *EntityConfig) GetGRPCServiceDescriptionName() string {
	return fmt.Sprintf("%sService_ServiceDesc", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) GetGRPCCreateDTOEncodeName() string {
	return fmt.Sprintf("encode%s", m.GetCreateModel().Name)
}

func (m *EntityConfig) GetGRPCUpdateDTOEncodeName() string {
	return fmt.Sprintf("encode%s", m.GetUpdateModel().Name)
}

func (m *EntityConfig) GetGRPCDeleteDTOEncodeName() string {
	return fmt.Sprintf("encode%s", m.GetDeleteModel().Name)
}

func (m *EntityConfig) GetGRPCFilterDTOEncodeName() string {
	return fmt.Sprintf("encode%s", m.GetFilterModel().Name)
}

func (m *EntityConfig) GetGRPCMainDecodeName() string {
	return fmt.Sprintf("decode%s", m.GetMainModel().Name)
}

func (m *EntityConfig) GetGRPCMainListDecodeName() string {
	return fmt.Sprintf("decodeList%s", m.GetMainModel().Name)
}

func (m *EntityConfig) GetGRPCUpdateDecodeName() string {
	return fmt.Sprintf("decode%s", m.GetUpdateModel().Name)
}

func (m *EntityConfig) GRPCHandlerTypeName() string {
	return fmt.Sprintf("%sServiceServer", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) GRPCHandlerVariableName() string {
	return fmt.Sprintf("%sHandler", strcase.ToLowerCamel(m.Name))
}

func (m *EntityConfig) ProtoFileName() string {
	return fmt.Sprintf("%s.proto", m.SnakeName())
}
