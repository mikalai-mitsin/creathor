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

func (m *Model) SyncModelStruct() error {
	model := &Struct{
		Path: filepath.Join("internal", "domain", "models", m.FileName()),
		Name: m.ModelName(),
		Params: []*Param{
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
		},
	}
	model.Params = append(model.Params, m.Params...)
	if err := SyncStruct(model); err != nil {
		return err
	}
	if err := SyncValidate(model); err != nil {
		return err
	}
	mockFilePath := filepath.Join("internal", "domain", "models", "mock", m.FileName())
	if err := SyncMock(mockFilePath, model); err != nil {
		return err
	}
	return nil
}

func (m *Model) SyncCreateStruct() error {
	create := &Struct{
		Path:   filepath.Join("internal", "domain", "models", m.FileName()),
		Name:   m.CreateTypeName(),
		Params: []*Param{},
	}
	for _, param := range m.Params {
		create.Params = append(create.Params, &Param{
			Name:   param.GetName(),
			Type:   param.Type,
			Search: false,
		})
	}
	if err := SyncStruct(create); err != nil {
		return err
	}

	if err := SyncValidate(create); err != nil {
		return err
	}
	mockFilePath := filepath.Join("internal", "domain", "models", "mock", m.FileName())
	if err := SyncMock(mockFilePath, create); err != nil {
		return err
	}
	return nil
}

func (m *Model) SyncUpdateStruct() error {
	update := &Struct{
		Path: filepath.Join("internal", "domain", "models", m.FileName()),
		Name: m.UpdateTypeName(),
		Params: []*Param{
			{
				Name:   "ID",
				Type:   "UUID",
				Search: false,
			},
		},
	}
	for _, param := range m.Params {
		update.Params = append(update.Params, &Param{
			Name:   param.GetName(),
			Type:   fmt.Sprintf("*%s", param.Type),
			Search: false,
		})
	}
	if err := SyncStruct(update); err != nil {
		return err
	}
	if err := SyncValidate(update); err != nil {
		return err
	}
	mockFilePath := filepath.Join("internal", "domain", "models", "mock", m.FileName())
	if err := SyncMock(mockFilePath, update); err != nil {
		return err
	}
	return nil
}

func (m *Model) SyncRepositoryInterface() error {
	usecase := &Interface{
		Path:     filepath.Join("internal", "domain", "repositories", m.FileName()),
		Name:     m.RepositoryTypeName(),
		Comments: nil,
		Methods: []*Method{
			{
				Name: "Get",
				Args: []*Param{
					{
						Name:   "ctx",
						Type:   "context.Context",
						Search: false,
					},
					{
						Name:   "id",
						Type:   "models.UUID",
						Search: false,
					},
				},
				Results: []*Param{
					{
						Name:   "",
						Type:   fmt.Sprintf("*models.%s", m.ModelName()),
						Search: false,
					},
					{
						Name:   "",
						Type:   "error",
						Search: false,
					},
				},
			},
			{
				Name: "List",
				Args: []*Param{
					{
						Name:   "ctx",
						Type:   "context.Context",
						Search: false,
					},
					{
						Name:   "filter",
						Type:   fmt.Sprintf("*models.%s", m.FilterTypeName()),
						Search: false,
					},
				},
				Results: []*Param{
					{
						Name:   "",
						Type:   fmt.Sprintf("[]*models.%s", m.ModelName()),
						Search: false,
					},
					{
						Name:   "",
						Type:   "error",
						Search: false,
					},
				},
			},
			{
				Name: "Count",
				Args: []*Param{
					{
						Name:   "ctx",
						Type:   "context.Context",
						Search: false,
					},
					{
						Name:   "filter",
						Type:   fmt.Sprintf("*models.%s", m.FilterTypeName()),
						Search: false,
					},
				},
				Results: []*Param{
					{
						Name:   "",
						Type:   "uint64",
						Search: false,
					},
					{
						Name:   "",
						Type:   "error",
						Search: false,
					},
				},
			},
			{
				Name: "Update",
				Args: []*Param{
					{
						Name:   "ctx",
						Type:   "context.Context",
						Search: false,
					},
					{
						Name:   "update",
						Type:   fmt.Sprintf("*models.%s", m.ModelName()),
						Search: false,
					},
				},
				Results: []*Param{
					{
						Name:   "",
						Type:   "error",
						Search: false,
					},
				},
			},
			{
				Name: "Create",
				Args: []*Param{
					{
						Name:   "ctx",
						Type:   "context.Context",
						Search: false,
					},
					{
						Name:   "create",
						Type:   fmt.Sprintf("*models.%s", m.ModelName()),
						Search: false,
					},
				},
				Results: []*Param{
					{
						Name:   "",
						Type:   "error",
						Search: false,
					},
				},
			},
			{
				Name: "Delete",
				Args: []*Param{
					{
						Name:   "ctx",
						Type:   "context.Context",
						Search: false,
					},
					{
						Name:   "id",
						Type:   "models.UUID",
						Search: false,
					},
				},
				Results: []*Param{
					{
						Name:   "",
						Type:   "error",
						Search: false,
					},
				},
			},
		},
	}
	if err := SyncInterface(usecase); err != nil {
		return err
	}
	return nil
}

func (m *Model) SyncUsecaseInterface() error {
	usecase := &Interface{
		Path:     filepath.Join("internal", "domain", "usecases", m.FileName()),
		Name:     m.UseCaseTypeName(),
		Comments: nil,
		Methods: []*Method{
			{
				Name: "Get",
				Args: []*Param{
					{
						Name:   "ctx",
						Type:   "context.Context",
						Search: false,
					},
					{
						Name:   "id",
						Type:   "models.UUID",
						Search: false,
					},
				},
				Results: []*Param{
					{
						Name:   "",
						Type:   fmt.Sprintf("*models.%s", m.ModelName()),
						Search: false,
					},
					{
						Name:   "",
						Type:   "error",
						Search: false,
					},
				},
			},
			{
				Name: "List",
				Args: []*Param{
					{
						Name:   "ctx",
						Type:   "context.Context",
						Search: false,
					},
					{
						Name:   "filter",
						Type:   fmt.Sprintf("*models.%s", m.FilterTypeName()),
						Search: false,
					},
				},
				Results: []*Param{
					{
						Name:   "",
						Type:   fmt.Sprintf("[]*models.%s", m.ModelName()),
						Search: false,
					},
					{
						Name:   "",
						Type:   "uint64",
						Search: false,
					},
					{
						Name:   "",
						Type:   "error",
						Search: false,
					},
				},
			},
			{
				Name: "Update",
				Args: []*Param{
					{
						Name:   "ctx",
						Type:   "context.Context",
						Search: false,
					},
					{
						Name:   "update",
						Type:   fmt.Sprintf("*models.%s", m.UpdateTypeName()),
						Search: false,
					},
				},
				Results: []*Param{
					{
						Name:   "",
						Type:   fmt.Sprintf("*models.%s", m.ModelName()),
						Search: false,
					},
					{
						Name:   "",
						Type:   "error",
						Search: false,
					},
				},
			},
			{
				Name: "Create",
				Args: []*Param{
					{
						Name:   "ctx",
						Type:   "context.Context",
						Search: false,
					},
					{
						Name:   "create",
						Type:   fmt.Sprintf("*models.%s", m.CreateTypeName()),
						Search: false,
					},
				},
				Results: []*Param{
					{
						Name:   "",
						Type:   fmt.Sprintf("*models.%s", m.ModelName()),
						Search: false,
					},
					{
						Name:   "",
						Type:   "error",
						Search: false,
					},
				},
			},
			{
				Name: "Delete",
				Args: []*Param{
					{
						Name:   "ctx",
						Type:   "context.Context",
						Search: false,
					},
					{
						Name:   "id",
						Type:   "models.UUID",
						Search: false,
					},
				},
				Results: []*Param{
					{
						Name:   "",
						Type:   "error",
						Search: false,
					},
				},
			},
		},
	}
	if err := SyncInterface(usecase); err != nil {
		return err
	}
	return nil
}

func (m *Model) SyncInterceptorInterface() error {
	interceptor := &Interface{
		Path:     filepath.Join("internal", "domain", "interceptors", m.FileName()),
		Name:     m.InterceptorTypeName(),
		Comments: nil,
		Methods: []*Method{
			{
				Name: "Get",
				Args: []*Param{
					{
						Name:   "ctx",
						Type:   "context.Context",
						Search: false,
					},
					{
						Name:   "id",
						Type:   "models.UUID",
						Search: false,
					},
				},
				Results: []*Param{
					{
						Name:   "",
						Type:   fmt.Sprintf("*models.%s", m.ModelName()),
						Search: false,
					},
					{
						Name:   "",
						Type:   "error",
						Search: false,
					},
				},
			},
			{
				Name: "List",
				Args: []*Param{
					{
						Name:   "ctx",
						Type:   "context.Context",
						Search: false,
					},
					{
						Name:   "filter",
						Type:   fmt.Sprintf("*models.%s", m.FilterTypeName()),
						Search: false,
					},
				},
				Results: []*Param{
					{
						Name:   "",
						Type:   fmt.Sprintf("[]*models.%s", m.ModelName()),
						Search: false,
					},
					{
						Name:   "",
						Type:   "uint64",
						Search: false,
					},
					{
						Name:   "",
						Type:   "error",
						Search: false,
					},
				},
			},
			{
				Name: "Update",
				Args: []*Param{
					{
						Name:   "ctx",
						Type:   "context.Context",
						Search: false,
					},
					{
						Name:   "update",
						Type:   fmt.Sprintf("*models.%s", m.UpdateTypeName()),
						Search: false,
					},
				},
				Results: []*Param{
					{
						Name:   "",
						Type:   fmt.Sprintf("*models.%s", m.ModelName()),
						Search: false,
					},
					{
						Name:   "",
						Type:   "error",
						Search: false,
					},
				},
			},
			{
				Name: "Create",
				Args: []*Param{
					{
						Name:   "ctx",
						Type:   "context.Context",
						Search: false,
					},
					{
						Name:   "create",
						Type:   fmt.Sprintf("*models.%s", m.CreateTypeName()),
						Search: false,
					},
				},
				Results: []*Param{
					{
						Name:   "",
						Type:   fmt.Sprintf("*models.%s", m.ModelName()),
						Search: false,
					},
					{
						Name:   "",
						Type:   "error",
						Search: false,
					},
				},
			},
			{
				Name: "Delete",
				Args: []*Param{
					{
						Name:   "ctx",
						Type:   "context.Context",
						Search: false,
					},
					{
						Name:   "id",
						Type:   "models.UUID",
						Search: false,
					},
				},
				Results: []*Param{
					{
						Name:   "",
						Type:   "error",
						Search: false,
					},
				},
			},
		},
	}
	if m.Auth {
		for _, method := range interceptor.Methods {
			method.Args = append(method.Args, &Param{
				Name:   "requestUser",
				Type:   "*models.User",
				Search: false,
			})
		}
	}
	if err := SyncInterface(interceptor); err != nil {
		return err
	}
	return nil
}

func (m *Model) SyncUseCaseImplementation() error {
	useCase := &UseCase{
		Path:  filepath.Join("internal", "usecases", m.FileName()),
		Name:  m.UseCaseTypeName(),
		Model: m,
		Params: []*Param{
			{
				Name:   m.RepositoryVariableName(),
				Type:   fmt.Sprintf("repositories.%s", m.RepositoryTypeName()),
				Search: false,
			},
			{
				Name:   "clock",
				Type:   "clock.Clock",
				Search: false,
			},
			{
				Name:   "logger",
				Type:   "log.Logger",
				Search: false,
			},
		},
	}
	if err := useCase.SyncStruct(); err != nil {
		return err
	}
	if err := useCase.SyncConstructor(); err != nil {
		return err
	}
	if err := useCase.SyncCreateMethod(); err != nil {
		return err
	}
	if err := useCase.SyncGetMethod(); err != nil {
		return err
	}
	if err := useCase.SyncListMethod(); err != nil {
		return err
	}
	if err := useCase.SyncUpdateMethod(); err != nil {
		return err
	}
	if err := useCase.SyncDeleteMethod(); err != nil {
		return err
	}
	return nil
}

func (m *Model) SyncInterceptorImplementation() error {
	interceptor := &Interceptor{
		Path:  filepath.Join("internal", "interceptors", m.FileName()),
		Name:  m.InterceptorTypeName(),
		Model: m,
		Params: []*Param{
			{
				Name:   m.UseCaseTypeName(),
				Type:   fmt.Sprintf("usecases.%s", m.UseCaseTypeName()),
				Search: false,
			},
			{
				Name:   "logger",
				Type:   "log.Logger",
				Search: false,
			},
		},
	}
	if m.Auth {
		interceptor.Params = append(
			interceptor.Params,
			&Param{
				Name:   "authUseCase",
				Type:   "usecases.AuthUseCase",
				Search: false,
			},
		)
	}
	if err := interceptor.SyncStruct(); err != nil {
		return err
	}
	if err := interceptor.SyncConstructor(); err != nil {
		return err
	}
	if err := interceptor.SyncCreateMethod(); err != nil {
		return err
	}
	if err := interceptor.SyncGetMethod(); err != nil {
		return err
	}
	if err := interceptor.SyncListMethod(); err != nil {
		return err
	}
	if err := interceptor.SyncUpdateMethod(); err != nil {
		return err
	}
	if err := interceptor.SyncDeleteMethod(); err != nil {
		return err
	}
	return nil
}

func (m *Model) SyncModels() error {
	if err := m.SyncModelStruct(); err != nil {
		return err
	}
	if err := m.SyncCreateStruct(); err != nil {
		return err
	}
	if err := m.SyncUpdateStruct(); err != nil {
		return err
	}
	if err := m.SyncRepositoryInterface(); err != nil {
		return err
	}
	if err := m.SyncUsecaseInterface(); err != nil {
		return err
	}
	if err := m.SyncInterceptorInterface(); err != nil {
		return err
	}
	if err := m.SyncUseCaseImplementation(); err != nil {
		return err
	}
	if err := m.SyncInterceptorImplementation(); err != nil {
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

func (m *Model) UpdateTypeName() string {
	return fmt.Sprintf("%sUpdate", strcase.ToCamel(m.Model))
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
