package app

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"

	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
	"golang.org/x/exp/slices"
)

type App struct {
	Config      configs.AppConfig
	Name        string
	Module      string
	ProtoModule string
	Entities    []*BaseEntity
}

func (m *App) AppName() string {
	return m.Config.AppName()
}

type BaseEntity struct {
	Config      configs.EntityConfig
	Name        string
	Module      string
	ProtoModule string
	Entities    []*Entity
	AppConfig   *configs.AppConfig
}

func (m *BaseEntity) SnakeName() string {
	return strcase.ToSnake(m.Name)
}

func (m *BaseEntity) FileName() string {
	return fmt.Sprintf("%s.go", m.SnakeName())
}

func (m *BaseEntity) TestFileName() string {
	return fmt.Sprintf("%s_test.go", m.SnakeName())
}

func (m *BaseEntity) MigrationUpFileName() string {
	last, err := lastMigration()
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%06d_%s.up.sql", last+1, m.TableName())
}

func (m *BaseEntity) MigrationDownFileName() string {
	last, err := lastMigration()
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%06d_%s.down.sql", last+1, m.TableName())
}

func lastMigration() (int, error) {
	dir, err := os.ReadDir(path.Join("internal", "pkg", "postgres", "migrations"))
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

func (m *BaseEntity) CamelName() string {
	return strcase.ToCamel(m.Name)
}

func (m *BaseEntity) LowerCamelName() string {
	return strcase.ToLowerCamel(m.Name)
}

func (m *BaseEntity) AppName() string {
	return m.AppConfig.AppName()
}

func (m *BaseEntity) DirName() string {
	return strcase.ToSnake(m.Name)
}

func (m *BaseEntity) EntitiesImportPath() string {
	return fmt.Sprintf(`"%s/internal/app/%s/entities/%s"`, m.Module, m.AppName(), m.DirName())
}

func (m *BaseEntity) GetMainModel() *Entity {
	index := slices.IndexFunc(
		m.Entities,
		func(model *Entity) bool { return model.Type == EntityTypeMain },
	)
	if index >= 0 {
		return m.Entities[index]
	}
	return nil
}

func (m *BaseEntity) TableName() string {
	return strcase.ToSnake(inflection.Plural(m.Name))
}

func (m *BaseEntity) SearchEnabled() bool {
	return slices.ContainsFunc(
		m.GetMainModel().Params,
		func(param *configs.Param) bool { return param.Search },
	)
}

func (m *BaseEntity) GetCreateModel() *Entity {
	index := slices.IndexFunc(
		m.Entities,
		func(model *Entity) bool { return model.Type == EntityTypeCreate },
	)
	if index >= 0 {
		return m.Entities[index]
	}
	return nil
}

func (m *BaseEntity) GetUpdateModel() *Entity {
	index := slices.IndexFunc(
		m.Entities,
		func(model *Entity) bool { return model.Type == EntityTypeUpdate },
	)
	if index > 0 {
		return m.Entities[index]
	}
	return nil
}

func (m *BaseEntity) GetFilterModel() *Entity {
	index := slices.IndexFunc(
		m.Entities,
		func(model *Entity) bool { return model.Type == EntityTypeFilter },
	)
	if index > 0 {
		return m.Entities[index]
	}
	return nil
}

func (m *BaseEntity) PermissionIDCreate() string {
	return fmt.Sprintf("PermissionID%sCreate", strcase.ToCamel(m.CamelName()))
}

func (m *BaseEntity) PermissionIDUpdate() string {
	return fmt.Sprintf("PermissionID%sUpdate", m.CamelName())
}

func (m *BaseEntity) PermissionIDDelete() string {
	return fmt.Sprintf("PermissionID%sDelete", m.CamelName())
}

func (m *BaseEntity) PermissionIDDetail() string {
	return fmt.Sprintf("PermissionID%sDetail", m.CamelName())
}

func (m *BaseEntity) PermissionIDList() string {
	return fmt.Sprintf("PermissionID%sList", m.CamelName())
}

func (m *BaseEntity) GetOneVariableName() string {
	return strcase.ToLowerCamel(m.Config.Name)
}

func (m *BaseEntity) GetManyVariableName() string {
	return inflection.Plural(m.GetOneVariableName())
}

func (m *BaseEntity) GetHTTPPath() string {
	return strcase.ToSnake(inflection.Plural(m.GetOneVariableName()))
}

func (m *BaseEntity) GetGRPCHandlerPrivateVariableName() string {
	return fmt.Sprintf("grpc%sHandler", strcase.ToCamel(m.Config.Name))
}

func (m *BaseEntity) GetGRPCHandlerPublicVariableName() string {
	return fmt.Sprintf("%sHandler", strcase.ToCamel(m.Config.Name))
}

func (m *BaseEntity) GetGRPCHandlerTypeName() string {
	return fmt.Sprintf("%sServiceServer", strcase.ToCamel(m.Config.Name))
}

func (m *BaseEntity) GetGRPCHandlerConstructorName() string {
	return fmt.Sprintf("New%s", m.GetGRPCHandlerTypeName())
}

func (m *BaseEntity) GetGRPCServiceDescriptionName() string {
	return fmt.Sprintf("%sService_ServiceDesc", strcase.ToCamel(m.Config.Name))
}

func (m *BaseEntity) GetGRPCCreateDTOEncodeName() string {
	return fmt.Sprintf("encode%s", m.GetCreateModel().Name)
}

func (m *BaseEntity) GetGRPCUpdateDTOEncodeName() string {
	return fmt.Sprintf("encode%s", m.GetUpdateModel().Name)
}

func (m *BaseEntity) GetGRPCFilterDTOEncodeName() string {
	return fmt.Sprintf("encode%s", m.GetFilterModel().Name)
}

func (m *BaseEntity) GetGRPCMainDecodeName() string {
	return fmt.Sprintf("decode%s", m.GetMainModel().Name)
}
func (m *BaseEntity) GetGRPCMainListDecodeName() string {
	return fmt.Sprintf("decodeList%s", m.GetMainModel().Name)
}
func (m *BaseEntity) GetGRPCUpdateDecodeName() string {
	return fmt.Sprintf("decode%s", m.GetUpdateModel().Name)
}

func (m *BaseEntity) GetHTTPHandlerConstructorName() string {
	return fmt.Sprintf("New%s", m.GetHTTPHandlerTypeName())
}

func (m *BaseEntity) GetHTTPHandlerTypeName() string {
	return fmt.Sprintf("%sHandler", strcase.ToCamel(m.Config.Name))
}

func (m *BaseEntity) GetHTTPHandlerPrivateVariableName() string {
	return fmt.Sprintf("http%sHandler", strcase.ToCamel(m.Config.Name))
}

func (m *BaseEntity) GetHTTPItemDTOName() string {
	return fmt.Sprintf("%sDTO", strcase.ToCamel(m.GetMainModel().Name))
}
func (m *BaseEntity) GetHTTPItemDTOConstructorName() string {
	return fmt.Sprintf("New%s", m.GetHTTPItemDTOName())
}

func (m *BaseEntity) GetHTTPUpdateDTOName() string {
	return fmt.Sprintf("%sDTO", strcase.ToCamel(m.GetUpdateModel().Name))
}

func (m *BaseEntity) GetHTTPUpdateDTOConstructorName() string {
	return fmt.Sprintf("New%s", m.GetHTTPUpdateDTOName())
}

func (m *BaseEntity) GetHTTPCreateDTOName() string {
	return fmt.Sprintf("%sDTO", strcase.ToCamel(m.GetCreateModel().Name))
}
func (m *BaseEntity) GetHTTPCreateDTOConstructorName() string {
	return fmt.Sprintf("New%s", m.GetHTTPCreateDTOName())
}

func (m *BaseEntity) GetHTTPListDTOName() string {
	return fmt.Sprintf("%sListDTO", strcase.ToCamel(m.GetMainModel().Name))
}

func (m *BaseEntity) GetHTTPListDTOConstructorName() string {
	return fmt.Sprintf("New%s", strcase.ToCamel(m.GetHTTPListDTOName()))
}

func (m *BaseEntity) GetUseCasePrivateVariableName() string {
	return fmt.Sprintf("%sUseCase", strcase.ToLowerCamel(m.Config.Name))
}

func (m *BaseEntity) GetUseCasePublicVariableName() string {
	return fmt.Sprintf("%sUseCase", strcase.ToCamel(m.Config.Name))
}

func (m *BaseEntity) GetUseCaseTypeName() string {
	return fmt.Sprintf("%sUseCase", strcase.ToCamel(m.Config.Name))
}

func (m *BaseEntity) GetUseCaseInterfaceName() string {
	return fmt.Sprintf("%sUseCase", strcase.ToLowerCamel(m.Config.Name))
}

func (m *BaseEntity) GetUseCaseConstructorName() string {
	return fmt.Sprintf("New%s", m.GetUseCaseTypeName())
}

func (m *BaseEntity) GetServicePrivateVariableName() string {
	return fmt.Sprintf("%sService", strcase.ToLowerCamel(m.Config.Name))
}

func (m *BaseEntity) GetServicePublicVariableName() string {
	return fmt.Sprintf("%sService", strcase.ToCamel(m.Config.Name))
}

func (m *BaseEntity) GetServiceTypeName() string {
	return fmt.Sprintf("%sService", strcase.ToCamel(m.Config.Name))
}

func (m *BaseEntity) GetServiceInterfaceName() string {
	return fmt.Sprintf("%sService", strcase.ToLowerCamel(m.Config.Name))
}

func (m *BaseEntity) GetServiceConstructorName() string {
	return fmt.Sprintf("New%s", m.GetServiceTypeName())
}

func (m *BaseEntity) GetRepositoryPrivateVariableName() string {
	return fmt.Sprintf("%sRepository", strcase.ToLowerCamel(m.Config.Name))
}

func (m *BaseEntity) GetRepositoryPublicVariableName() string {
	return fmt.Sprintf("%sRepository", strcase.ToCamel(m.Config.Name))
}

func (m *BaseEntity) GetRepositoryTypeName() string {
	return fmt.Sprintf("%sRepository", strcase.ToCamel(m.Config.Name))
}

func (m *BaseEntity) GetRepositoryInterfaceName() string {
	return fmt.Sprintf("%sRepository", strcase.ToLowerCamel(m.Config.Name))
}

func (m *BaseEntity) GetRepositoryConstructorName() string {
	return fmt.Sprintf("New%s", m.GetRepositoryTypeName())
}

func (m *BaseEntity) GetHTTPFilterDTOName() string {
	return fmt.Sprintf("%sFilterDTO", strcase.ToCamel(m.Config.Name))
}

func (m *BaseEntity) GetHTTPFilterDTOConstructorName() string {
	return fmt.Sprintf("New%s", m.GetHTTPFilterDTOName())
}
