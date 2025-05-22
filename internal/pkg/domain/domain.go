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

type Domain struct {
	Config      *configs.DomainConfig
	Name        string
	Module      string
	ProtoModule string
	Entities    []*Model
	Auth        bool
}

func (m *Domain) SnakeName() string {
	return strcase.ToSnake(m.Name)
}

func (m *Domain) FileName() string {
	return fmt.Sprintf("%s.go", m.SnakeName())
}

func (m *Domain) TestFileName() string {
	return fmt.Sprintf("%s_test.go", m.SnakeName())
}

func (m *Domain) MigrationUpFileName() string {
	last, err := lastMigration()
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%06d_%s.up.sql", last+1, m.TableName())
}

func (m *Domain) MigrationDownFileName() string {
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

func (m *Domain) CamelName() string {
	return strcase.ToCamel(m.Name)
}

func (m *Domain) LowerCamelName() string {
	return strcase.ToLowerCamel(m.Name)
}

func (m *Domain) DirName() string {
	return m.SnakeName()
}

func (m *Domain) EntitiesImportPath() string {
	return fmt.Sprintf(`"%s/internal/app/%s/entities"`, m.Module, m.DirName())
}

func (m *Domain) GetMainModel() *Model {
	index := slices.IndexFunc(
		m.Entities,
		func(model *Model) bool { return model.Type == ModelTypeMain },
	)
	if index >= 0 {
		return m.Entities[index]
	}
	return nil
}

func (m *Domain) TableName() string {
	return strcase.ToSnake(inflection.Plural(m.Name))
}

func (m *Domain) SearchEnabled() bool {
	return slices.ContainsFunc(
		m.GetMainModel().Params,
		func(param *configs.Param) bool { return param.Search },
	)
}

func (m *Domain) GetCreateModel() *Model {
	index := slices.IndexFunc(
		m.Entities,
		func(model *Model) bool { return model.Type == ModelTypeCreate },
	)
	if index >= 0 {
		return m.Entities[index]
	}
	return nil
}

func (m *Domain) GetUpdateModel() *Model {
	index := slices.IndexFunc(
		m.Entities,
		func(model *Model) bool { return model.Type == ModelTypeUpdate },
	)
	if index > 0 {
		return m.Entities[index]
	}
	return nil
}

func (m *Domain) GetFilterModel() *Model {
	index := slices.IndexFunc(
		m.Entities,
		func(model *Model) bool { return model.Type == ModelTypeFilter },
	)
	if index > 0 {
		return m.Entities[index]
	}
	return nil
}

func (m *Domain) PermissionIDCreate() string {
	return fmt.Sprintf("PermissionID%sCreate", strcase.ToCamel(m.CamelName()))
}

func (m *Domain) PermissionIDUpdate() string {
	return fmt.Sprintf("PermissionID%sUpdate", m.CamelName())
}

func (m *Domain) PermissionIDDelete() string {
	return fmt.Sprintf("PermissionID%sDelete", m.CamelName())
}

func (m *Domain) PermissionIDDetail() string {
	return fmt.Sprintf("PermissionID%sDetail", m.CamelName())
}

func (m *Domain) PermissionIDList() string {
	return fmt.Sprintf("PermissionID%sList", m.CamelName())
}

func (m *Domain) GetOneVariableName() string {
	return strcase.ToLowerCamel(m.Config.Model)
}

func (m *Domain) GetManyVariableName() string {
	return inflection.Plural(m.GetOneVariableName())
}

func (m *Domain) GetHTTPPath() string {
	return strcase.ToSnake(inflection.Plural(m.GetOneVariableName()))
}

func (m *Domain) GetGRPCHandlerPrivateVariableName() string {
	return fmt.Sprintf("grpc%sHandler", strcase.ToCamel(m.Config.Model))
}

func (m *Domain) GetGRPCHandlerPublicVariableName() string {
	return fmt.Sprintf("%sHandler", strcase.ToCamel(m.Config.Model))
}

func (m *Domain) GetGRPCHandlerTypeName() string {
	return fmt.Sprintf("%sServiceServer", strcase.ToCamel(m.Config.Model))
}

func (m *Domain) GetGRPCHandlerConstructorName() string {
	return fmt.Sprintf("New%s", m.GetGRPCHandlerTypeName())
}

func (m *Domain) GetGRPCServiceDescriptionName() string {
	return fmt.Sprintf("%sService_ServiceDesc", strcase.ToCamel(m.Config.Model))
}

func (m *Domain) GetHTTPHandlerConstructorName() string {
	return fmt.Sprintf("New%s", m.GetHTTPHandlerTypeName())
}

func (m *Domain) GetHTTPHandlerTypeName() string {
	return fmt.Sprintf("%sHandler", strcase.ToCamel(m.Config.Model))
}

func (m *Domain) GetHTTPHandlerPrivateVariableName() string {
	return fmt.Sprintf("http%sHandler", strcase.ToCamel(m.Config.Model))
}

func (m *Domain) GetHTTPItemDTOName() string {
	return fmt.Sprintf("%sDTO", strcase.ToCamel(m.GetMainModel().Name))
}
func (m *Domain) GetHTTPItemDTOConstructorName() string {
	return fmt.Sprintf("New%s", m.GetHTTPItemDTOName())
}

func (m *Domain) GetHTTPUpdateDTOName() string {
	return fmt.Sprintf("%sDTO", strcase.ToCamel(m.GetUpdateModel().Name))
}
func (m *Domain) GetHTTPUpdateDTOConstructorName() string {
	return fmt.Sprintf("New%s", m.GetHTTPUpdateDTOName())
}

func (m *Domain) GetHTTPCreateDTOName() string {
	return fmt.Sprintf("%sDTO", strcase.ToCamel(m.GetCreateModel().Name))
}
func (m *Domain) GetHTTPCreateDTOConstructorName() string {
	return fmt.Sprintf("New%s", m.GetHTTPCreateDTOName())
}

func (m *Domain) GetHTTPListDTOName() string {
	return fmt.Sprintf("%sListDTO", strcase.ToCamel(m.GetMainModel().Name))
}

func (m *Domain) GetHTTPListDTOConstructorName() string {
	return fmt.Sprintf("New%s", strcase.ToCamel(m.GetHTTPListDTOName()))
}

func (m *Domain) GetUseCasePrivateVariableName() string {
	return fmt.Sprintf("%sUseCase", strcase.ToLowerCamel(m.Config.Model))
}

func (m *Domain) GetUseCasePublicVariableName() string {
	return fmt.Sprintf("%sUseCase", strcase.ToCamel(m.Config.Model))
}

func (m *Domain) GetUseCaseTypeName() string {
	return fmt.Sprintf("%sUseCase", strcase.ToCamel(m.Config.Model))
}

func (m *Domain) GetUseCaseInterfaceName() string {
	return fmt.Sprintf("%sUseCase", strcase.ToLowerCamel(m.Config.Model))
}

func (m *Domain) GetUseCaseConstructorName() string {
	return fmt.Sprintf("New%s", m.GetUseCaseTypeName())
}

func (m *Domain) GetServicePrivateVariableName() string {
	return fmt.Sprintf("%sService", strcase.ToLowerCamel(m.Config.Model))
}

func (m *Domain) GetServicePublicVariableName() string {
	return fmt.Sprintf("%sService", strcase.ToCamel(m.Config.Model))
}

func (m *Domain) GetServiceTypeName() string {
	return fmt.Sprintf("%sService", strcase.ToCamel(m.Config.Model))
}

func (m *Domain) GetServiceInterfaceName() string {
	return fmt.Sprintf("%sService", strcase.ToLowerCamel(m.Config.Model))
}

func (m *Domain) GetServiceConstructorName() string {
	return fmt.Sprintf("New%s", m.GetServiceTypeName())
}

func (m *Domain) GetRepositoryPrivateVariableName() string {
	return fmt.Sprintf("%sRepository", strcase.ToLowerCamel(m.Config.Model))
}

func (m *Domain) GetRepositoryPublicVariableName() string {
	return fmt.Sprintf("%sRepository", strcase.ToCamel(m.Config.Model))
}

func (m *Domain) GetRepositoryTypeName() string {
	return fmt.Sprintf("%sRepository", strcase.ToCamel(m.Config.Model))
}

func (m *Domain) GetRepositoryInterfaceName() string {
	return fmt.Sprintf("%sRepository", strcase.ToLowerCamel(m.Config.Model))
}

func (m *Domain) GetRepositoryConstructorName() string {
	return fmt.Sprintf("New%s", m.GetRepositoryTypeName())
}

func (m *Domain) GetHTTPFilterDTOName() string {
	return fmt.Sprintf("%sFilterDTO", strcase.ToCamel(m.Config.Model))
}

func (m *Domain) GetHTTPFilterDTOConstructorName() string {
	return fmt.Sprintf("New%s", m.GetHTTPFilterDTOName())
}
