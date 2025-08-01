package domain

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
	Entities    []*Entity
	Auth        bool
}

func (m *App) SnakeName() string {
	return strcase.ToSnake(m.Name)
}

func (m *App) FileName() string {
	return fmt.Sprintf("%s.go", m.SnakeName())
}

func (m *App) TestFileName() string {
	return fmt.Sprintf("%s_test.go", m.SnakeName())
}

func (m *App) MigrationUpFileName() string {
	last, err := lastMigration()
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%06d_%s.up.sql", last+1, m.TableName())
}

func (m *App) MigrationDownFileName() string {
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

func (m *App) CamelName() string {
	return strcase.ToCamel(m.Name)
}

func (m *App) LowerCamelName() string {
	return strcase.ToLowerCamel(m.Name)
}

func (m *App) DirName() string {
	return m.SnakeName()
}

func (m *App) EntitiesImportPath() string {
	return fmt.Sprintf(`"%s/internal/app/%s/entities"`, m.Module, m.DirName())
}

func (m *App) GetMainModel() *Entity {
	index := slices.IndexFunc(
		m.Entities,
		func(model *Entity) bool { return model.Type == EntityTypeMain },
	)
	if index >= 0 {
		return m.Entities[index]
	}
	return nil
}

func (m *App) TableName() string {
	return strcase.ToSnake(inflection.Plural(m.Name))
}

func (m *App) SearchEnabled() bool {
	return slices.ContainsFunc(
		m.GetMainModel().Params,
		func(param *configs.Param) bool { return param.Search },
	)
}

func (m *App) GetCreateModel() *Entity {
	index := slices.IndexFunc(
		m.Entities,
		func(model *Entity) bool { return model.Type == EntityTypeCreate },
	)
	if index >= 0 {
		return m.Entities[index]
	}
	return nil
}

func (m *App) GetUpdateModel() *Entity {
	index := slices.IndexFunc(
		m.Entities,
		func(model *Entity) bool { return model.Type == EntityTypeUpdate },
	)
	if index > 0 {
		return m.Entities[index]
	}
	return nil
}

func (m *App) GetFilterModel() *Entity {
	index := slices.IndexFunc(
		m.Entities,
		func(model *Entity) bool { return model.Type == EntityTypeFilter },
	)
	if index > 0 {
		return m.Entities[index]
	}
	return nil
}

func (m *App) PermissionIDCreate() string {
	return fmt.Sprintf("PermissionID%sCreate", strcase.ToCamel(m.CamelName()))
}

func (m *App) PermissionIDUpdate() string {
	return fmt.Sprintf("PermissionID%sUpdate", m.CamelName())
}

func (m *App) PermissionIDDelete() string {
	return fmt.Sprintf("PermissionID%sDelete", m.CamelName())
}

func (m *App) PermissionIDDetail() string {
	return fmt.Sprintf("PermissionID%sDetail", m.CamelName())
}

func (m *App) PermissionIDList() string {
	return fmt.Sprintf("PermissionID%sList", m.CamelName())
}

func (m *App) GetOneVariableName() string {
	return strcase.ToLowerCamel(m.Config.Name)
}

func (m *App) GetManyVariableName() string {
	return inflection.Plural(m.GetOneVariableName())
}

func (m *App) GetHTTPPath() string {
	return strcase.ToSnake(inflection.Plural(m.GetOneVariableName()))
}

func (m *App) GetGRPCHandlerPrivateVariableName() string {
	return fmt.Sprintf("grpc%sHandler", strcase.ToCamel(m.Config.Name))
}

func (m *App) GetGRPCHandlerPublicVariableName() string {
	return fmt.Sprintf("%sHandler", strcase.ToCamel(m.Config.Name))
}

func (m *App) GetGRPCHandlerTypeName() string {
	return fmt.Sprintf("%sServiceServer", strcase.ToCamel(m.Config.Name))
}

func (m *App) GetGRPCHandlerConstructorName() string {
	return fmt.Sprintf("New%s", m.GetGRPCHandlerTypeName())
}

func (m *App) GetGRPCServiceDescriptionName() string {
	return fmt.Sprintf("%sService_ServiceDesc", strcase.ToCamel(m.Config.Name))
}

func (m *App) GetHTTPHandlerConstructorName() string {
	return fmt.Sprintf("New%s", m.GetHTTPHandlerTypeName())
}

func (m *App) GetHTTPHandlerTypeName() string {
	return fmt.Sprintf("%sHandler", strcase.ToCamel(m.Config.Name))
}

func (m *App) GetHTTPHandlerPrivateVariableName() string {
	return fmt.Sprintf("http%sHandler", strcase.ToCamel(m.Config.Name))
}

func (m *App) GetHTTPItemDTOName() string {
	return fmt.Sprintf("%sDTO", strcase.ToCamel(m.GetMainModel().Name))
}
func (m *App) GetHTTPItemDTOConstructorName() string {
	return fmt.Sprintf("New%s", m.GetHTTPItemDTOName())
}

func (m *App) GetHTTPUpdateDTOName() string {
	return fmt.Sprintf("%sDTO", strcase.ToCamel(m.GetUpdateModel().Name))
}
func (m *App) GetHTTPUpdateDTOConstructorName() string {
	return fmt.Sprintf("New%s", m.GetHTTPUpdateDTOName())
}

func (m *App) GetHTTPCreateDTOName() string {
	return fmt.Sprintf("%sDTO", strcase.ToCamel(m.GetCreateModel().Name))
}
func (m *App) GetHTTPCreateDTOConstructorName() string {
	return fmt.Sprintf("New%s", m.GetHTTPCreateDTOName())
}

func (m *App) GetHTTPListDTOName() string {
	return fmt.Sprintf("%sListDTO", strcase.ToCamel(m.GetMainModel().Name))
}

func (m *App) GetHTTPListDTOConstructorName() string {
	return fmt.Sprintf("New%s", strcase.ToCamel(m.GetHTTPListDTOName()))
}

func (m *App) GetUseCasePrivateVariableName() string {
	return fmt.Sprintf("%sUseCase", strcase.ToLowerCamel(m.Config.Name))
}

func (m *App) GetUseCasePublicVariableName() string {
	return fmt.Sprintf("%sUseCase", strcase.ToCamel(m.Config.Name))
}

func (m *App) GetUseCaseTypeName() string {
	return fmt.Sprintf("%sUseCase", strcase.ToCamel(m.Config.Name))
}

func (m *App) GetUseCaseInterfaceName() string {
	return fmt.Sprintf("%sUseCase", strcase.ToLowerCamel(m.Config.Name))
}

func (m *App) GetUseCaseConstructorName() string {
	return fmt.Sprintf("New%s", m.GetUseCaseTypeName())
}

func (m *App) GetServicePrivateVariableName() string {
	return fmt.Sprintf("%sService", strcase.ToLowerCamel(m.Config.Name))
}

func (m *App) GetServicePublicVariableName() string {
	return fmt.Sprintf("%sService", strcase.ToCamel(m.Config.Name))
}

func (m *App) GetServiceTypeName() string {
	return fmt.Sprintf("%sService", strcase.ToCamel(m.Config.Name))
}

func (m *App) GetServiceInterfaceName() string {
	return fmt.Sprintf("%sService", strcase.ToLowerCamel(m.Config.Name))
}

func (m *App) GetServiceConstructorName() string {
	return fmt.Sprintf("New%s", m.GetServiceTypeName())
}

func (m *App) GetRepositoryPrivateVariableName() string {
	return fmt.Sprintf("%sRepository", strcase.ToLowerCamel(m.Config.Name))
}

func (m *App) GetRepositoryPublicVariableName() string {
	return fmt.Sprintf("%sRepository", strcase.ToCamel(m.Config.Name))
}

func (m *App) GetRepositoryTypeName() string {
	return fmt.Sprintf("%sRepository", strcase.ToCamel(m.Config.Name))
}

func (m *App) GetRepositoryInterfaceName() string {
	return fmt.Sprintf("%sRepository", strcase.ToLowerCamel(m.Config.Name))
}

func (m *App) GetRepositoryConstructorName() string {
	return fmt.Sprintf("New%s", m.GetRepositoryTypeName())
}

func (m *App) GetHTTPFilterDTOName() string {
	return fmt.Sprintf("%sFilterDTO", strcase.ToCamel(m.Config.Name))
}

func (m *App) GetHTTPFilterDTOConstructorName() string {
	return fmt.Sprintf("New%s", m.GetHTTPFilterDTOName())
}
