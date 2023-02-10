package models

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type Model struct {
	Model          string   `json:"model" yaml:"model"`
	Module         string   `json:"module" yaml:"module"`
	ProjectName    string   `json:"project_name" yaml:"projectName"`
	ProtoPackage   string   `json:"proto_package" yaml:"protoPackage"`
	Auth           bool     `json:"auth" yaml:"auth"`
	Params         []*Param `json:"params" yaml:"params"`
	GRPCEnabled    bool     `yaml:"gRPC"`
	GatewayEnabled bool     `yaml:"gateway"`
	RESTEnabled    bool     `yaml:"REST"`
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

func (m *Model) IsExists() bool {
	packagePath := filepath.Join("internal", "domain", "models")
	fileset := token.NewFileSet()
	tree, err := parser.ParseDir(fileset, packagePath, func(info fs.FileInfo) bool {
		return true
	}, parser.ParseComments)
	if err != nil {
		return false
	}
	for _, p := range tree {
		for _, file := range p.Files {
			for _, decl := range file.Decls {
				gen, ok := decl.(*ast.GenDecl)
				if ok {
					for _, spec := range gen.Specs {
						t, ok := spec.(*ast.TypeSpec)
						if ok && t.Name.String() == m.ModelName() {
							return true
						}
					}
				}
			}
		}
	}
	return false
}

func (m *Model) SyncModel() error {
	modelFilePath := filepath.Join("internal", "domain", "models", m.FileName())
	modelParams := []*Param{
		{
			Name:   "ID",
			Type:   "UUID",
			Search: false,
		},
		{
			Name:   "UpdatedAt",
			Type:   "time.Time",
			Search: false,
		},
		{
			Name:   "CreatedAt",
			Type:   "time.Time",
			Search: false,
		},
	}
	modelParams = append(modelParams, m.Params...)
	if err := SyncStruct(modelFilePath, m.ModelName(), modelParams); err != nil {
		return err
	}
	if err := SyncValidate(modelFilePath, m.ModelName(), modelParams); err != nil {
		return err
	}
	mockFilePath := filepath.Join("internal", "domain", "models", "mock", m.FileName())
	if err := SyncMock(mockFilePath, m.ModelName(), modelParams); err != nil {
		return err
	}
	return nil
}

func (m *Model) SyncCreate() error {
	filePath := filepath.Join("internal", "domain", "models", m.FileName())
	var createParams []*Param
	for _, param := range m.Params {
		createParams = append(createParams, &Param{
			Name:   param.GetName(),
			Type:   param.Type,
			Search: false,
		})
	}
	if err := SyncStruct(filePath, m.CreateTypeName(), createParams); err != nil {
		return err
	}

	if err := SyncValidate(filePath, m.CreateTypeName(), createParams); err != nil {
		return err
	}
	mockFilePath := filepath.Join("internal", "domain", "models", "mock", m.FileName())
	if err := SyncMock(mockFilePath, m.CreateTypeName(), createParams); err != nil {
		return err
	}
	return nil
}

func (m *Model) SyncUpdate() error {
	filePath := filepath.Join("internal", "domain", "models", m.FileName())
	updateParams := []*Param{
		{
			Name:   "ID",
			Type:   "UUID",
			Search: false,
		},
	}
	for _, param := range m.Params {
		updateParams = append(updateParams, &Param{
			Name:   param.GetName(),
			Type:   fmt.Sprintf("*%s", param.Type),
			Search: false,
		})
	}
	if err := SyncStruct(filePath, m.UpdateTypeName(), updateParams); err != nil {
		return err
	}
	if err := SyncValidate(filePath, m.UpdateTypeName(), updateParams); err != nil {
		return err
	}
	mockFilePath := filepath.Join("internal", "domain", "models", "mock", m.FileName())
	if err := SyncMock(mockFilePath, m.UpdateTypeName(), updateParams); err != nil {
		return err
	}
	return nil
}

func (m *Model) SyncModels() error {
	if err := m.SyncModel(); err != nil {
		return err
	}
	if err := m.SyncCreate(); err != nil {
		return err
	}
	if err := m.SyncUpdate(); err != nil {
		return err
	}
	return nil
}

func (m *Model) SearchEnabled() bool {
	for _, param := range m.Params {
		if param.Search {
			return true
		}
	}
	return false
}

func (m *Model) SearchVector() string {
	var params []string
	for _, param := range m.Params {
		if param.Search {
			params = append(params, param.Tag())
		}
	}
	vector := fmt.Sprintf("to_tsvector('english', %s)", strings.Join(params, " || "))
	return vector
}

func (m *Model) Variable() string {
	return strcase.ToLowerCamel(m.Model)
}

func (m *Model) ListVariable() string {
	return strcase.ToLowerCamel(fmt.Sprintf("list%s", strcase.ToCamel(inflection.Plural(m.Model))))
}

func (m *Model) ModelName() string {
	return strcase.ToCamel(m.Model)
}

func (m *Model) UseCaseTypeName() string {
	return fmt.Sprintf("%sUseCase", strcase.ToCamel(m.Model))
}

func (m *Model) GRPCHandlerTypeName() string {
	return fmt.Sprintf("%sServiceServer", strcase.ToCamel(m.Model))
}

func (m *Model) RESTHandlerTypeName() string {
	return fmt.Sprintf("%sHandler", strcase.ToCamel(m.Model))
}

func (m *Model) GatewayHandlerTypeName() string {
	return fmt.Sprintf("Register%sServiceHandlerFromEndpoint", strcase.ToCamel(m.Model))
}

func (m *Model) RESTHandlerPath() string {
	return strcase.ToSnake(inflection.Plural(m.Model))
}

func (m *Model) RESTHandlerVariableName() string {
	return fmt.Sprintf("%sHandler", strcase.ToLowerCamel(m.Model))
}

func (m *Model) GRPCHandlerVariableName() string {
	return fmt.Sprintf("%sHandler", strcase.ToLowerCamel(m.Model))
}

func (m *Model) UseCaseVariableName() string {
	return fmt.Sprintf("%sUseCase", strcase.ToLowerCamel(m.Model))
}

func (m *Model) InterceptorTypeName() string {
	return fmt.Sprintf("%sInterceptor", strcase.ToCamel(m.Model))
}

func (m *Model) InterceptorVariableName() string {
	return fmt.Sprintf("%sInterceptor", strcase.ToLowerCamel(m.Model))
}

func (m *Model) RepositoryTypeName() string {
	return fmt.Sprintf("%sRepository", strcase.ToCamel(m.Model))
}

func (m *Model) RepositoryVariableName() string {
	return fmt.Sprintf("%sRepository", strcase.ToLowerCamel(m.Model))
}

func (m *Model) FilterTypeName() string {
	return fmt.Sprintf("%sFilter", strcase.ToCamel(m.Model))
}

func (m *Model) FilterVariableName() string {
	return fmt.Sprintf("%sFilter", strcase.ToLowerCamel(m.Model))
}

func (m *Model) UpdateTypeName() string {
	return fmt.Sprintf("%sUpdate", strcase.ToCamel(m.Model))
}

func (m *Model) UpdateVariableName() string {
	return fmt.Sprintf("%sUpdate", strcase.ToLowerCamel(m.Model))
}

func (m *Model) CreateTypeName() string {
	return fmt.Sprintf("%sCreate", strcase.ToCamel(m.Model))
}

func (m *Model) PostgresDTOTypeName() string {
	return fmt.Sprintf("%sDTO", strcase.ToCamel(m.Model))
}

func (m *Model) PostgresDTOListTypeName() string {
	return fmt.Sprintf("%sListDTO", strcase.ToCamel(m.Model))
}

func (m *Model) CreateVariableName() string {
	return fmt.Sprintf("%sCreate", strcase.ToLowerCamel(m.Model))
}

func (m *Model) KeyName() string {
	return strcase.ToSnake(m.Model)
}

func (m *Model) SnakeName() string {
	return strcase.ToSnake(m.Model)
}

func (m *Model) FileName() string {
	return fmt.Sprintf("%s.go", m.SnakeName())
}

func (m *Model) MigrationUpFileName() string {
	last, err := lastMigration()
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%06d_%s.up.sql", last+1, m.TableName())
}

func (m *Model) MigrationDownFileName() string {
	last, err := lastMigration()
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%06d_%s.down.sql", last+1, m.TableName())
}

func lastMigration() (int, error) {
	dir, err := os.ReadDir(path.Join("internal", "interfaces", "postgres", "migrations"))
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

func (m *Model) TestFileName() string {
	return fmt.Sprintf("%s_test.go", m.SnakeName())
}

func (m *Model) ProtoFileName() string {
	return fmt.Sprintf("%s.proto", m.SnakeName())
}

func (m *Model) MockFileName() string {
	return fmt.Sprintf("%s_mock.go", m.SnakeName())
}

func (m *Model) TableName() string {
	return strcase.ToSnake(inflection.Plural(m.Model))
}

func (m *Model) PermissionIDList() string {
	return fmt.Sprintf("PermissionID%sList", m.ModelName())
}

func (m *Model) PermissionIDDetail() string {
	return fmt.Sprintf("PermissionID%sDetail", m.ModelName())
}

func (m *Model) PermissionIDCreate() string {
	return fmt.Sprintf("PermissionID%sCreate", m.ModelName())
}

func (m *Model) PermissionIDUpdate() string {
	return fmt.Sprintf("PermissionID%sUpdate", m.ModelName())
}

func (m *Model) PermissionIDDelete() string {
	return fmt.Sprintf("PermissionID%sDelete", m.ModelName())
}
