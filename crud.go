package main

import (
	"bytes"
	"fmt"
	"github.com/018bf/creathor/internal/generators"
	"github.com/018bf/creathor/internal/models"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

func SyncModelStruct(m *models.ModelConfig) error {
	model := &generators.Model{
		Name:        m.ModelName(),
		ModelConfig: m,
		Params: []*generators.Param{
			{
				Name: "ID",
				Type: "UUID",
			},
			{
				Name: "UpdatedAt",
				Type: "time.Time",
			},
			{
				Name: "CreatedAt",
				Type: "time.Time",
			},
		},
	}
	for _, param := range m.Params {
		model.Params = append(
			model.Params,
			&generators.Param{
				Name: param.GetName(),
				Type: param.Type,
			},
		)
	}
	if err := model.Sync(); err != nil {
		return err
	}
	return nil
}

func SyncCreateStruct(m *models.ModelConfig) error {
	create := &generators.Model{
		Name:        m.CreateTypeName(),
		ModelConfig: m,
		Params:      []*generators.Param{},
	}
	for _, param := range m.Params {
		create.Params = append(create.Params, &generators.Param{
			Name: param.GetName(),
			Type: param.Type,
		})
	}
	if err := create.Sync(); err != nil {
		return err
	}
	return nil
}

func SyncUpdateStruct(m *models.ModelConfig) error {
	update := &generators.Model{
		Name:        m.UpdateTypeName(),
		ModelConfig: m,
		Params: []*generators.Param{
			{
				Name: "ID",
				Type: "UUID",
			},
		},
	}
	for _, param := range m.Params {
		update.Params = append(update.Params, &generators.Param{
			Name: param.GetName(),
			Type: fmt.Sprintf("*%s", param.Type),
		})
	}
	if err := update.Sync(); err != nil {
		return err
	}
	return nil
}

func SyncRepositoryInterface(m *models.ModelConfig) error {
	usecase := &generators.Interface{
		Path:     filepath.Join("internal", "domain", "repositories", m.FileName()),
		Name:     m.RepositoryTypeName(),
		Comments: nil,
		Methods: []*generators.Method{
			{
				Name: "Get",
				Args: []*generators.Param{
					{
						Name: "ctx",
						Type: "context.Context",
					},
					{
						Name: "id",
						Type: "models.UUID",
					},
				},
				Results: []*generators.Param{
					{
						Name: "",
						Type: fmt.Sprintf("*models.%s", m.ModelName()),
					},
					{
						Name: "",
						Type: "error",
					},
				},
			},
			{
				Name: "List",
				Args: []*generators.Param{
					{
						Name: "ctx",
						Type: "context.Context",
					},
					{
						Name: "filter",
						Type: fmt.Sprintf("*models.%s", m.FilterTypeName()),
					},
				},
				Results: []*generators.Param{
					{
						Name: "",
						Type: fmt.Sprintf("[]*models.%s", m.ModelName()),
					},
					{
						Name: "",
						Type: "error",
					},
				},
			},
			{
				Name: "Count",
				Args: []*generators.Param{
					{
						Name: "ctx",
						Type: "context.Context",
					},
					{
						Name: "filter",
						Type: fmt.Sprintf("*models.%s", m.FilterTypeName()),
					},
				},
				Results: []*generators.Param{
					{
						Name: "",
						Type: "uint64",
					},
					{
						Name: "",
						Type: "error",
					},
				},
			},
			{
				Name: "Update",
				Args: []*generators.Param{
					{
						Name: "ctx",
						Type: "context.Context",
					},
					{
						Name: "update",
						Type: fmt.Sprintf("*models.%s", m.ModelName()),
					},
				},
				Results: []*generators.Param{
					{
						Name: "",
						Type: "error",
					},
				},
			},
			{
				Name: "Create",
				Args: []*generators.Param{
					{
						Name: "ctx",
						Type: "context.Context",
					},
					{
						Name: "create",
						Type: fmt.Sprintf("*models.%s", m.ModelName()),
					},
				},
				Results: []*generators.Param{
					{
						Name: "",
						Type: "error",
					},
				},
			},
			{
				Name: "Delete",
				Args: []*generators.Param{
					{
						Name: "ctx",
						Type: "context.Context",
					},
					{
						Name: "id",
						Type: "models.UUID",
					},
				},
				Results: []*generators.Param{
					{
						Name: "",
						Type: "error",
					},
				},
			},
		},
	}
	if err := usecase.SyncInterface(); err != nil {
		return err
	}
	return nil
}

func SyncUsecaseInterface(m *models.ModelConfig) error {
	usecase := &generators.Interface{
		Path:     filepath.Join("internal", "domain", "usecases", m.FileName()),
		Name:     m.UseCaseTypeName(),
		Comments: nil,
		Methods: []*generators.Method{
			{
				Name: "Get",
				Args: []*generators.Param{
					{
						Name: "ctx",
						Type: "context.Context",
					},
					{
						Name: "id",
						Type: "models.UUID",
					},
				},
				Results: []*generators.Param{
					{
						Name: "",
						Type: fmt.Sprintf("*models.%s", m.ModelName()),
					},
					{
						Name: "",
						Type: "error",
					},
				},
			},
			{
				Name: "List",
				Args: []*generators.Param{
					{
						Name: "ctx",
						Type: "context.Context",
					},
					{
						Name: "filter",
						Type: fmt.Sprintf("*models.%s", m.FilterTypeName()),
					},
				},
				Results: []*generators.Param{
					{
						Name: "",
						Type: fmt.Sprintf("[]*models.%s", m.ModelName()),
					},
					{
						Name: "",
						Type: "uint64",
					},
					{
						Name: "",
						Type: "error",
					},
				},
			},
			{
				Name: "Update",
				Args: []*generators.Param{
					{
						Name: "ctx",
						Type: "context.Context",
					},
					{
						Name: "update",
						Type: fmt.Sprintf("*models.%s", m.UpdateTypeName()),
					},
				},
				Results: []*generators.Param{
					{
						Name: "",
						Type: fmt.Sprintf("*models.%s", m.ModelName()),
					},
					{
						Name: "",
						Type: "error",
					},
				},
			},
			{
				Name: "Create",
				Args: []*generators.Param{
					{
						Name: "ctx",
						Type: "context.Context",
					},
					{
						Name: "create",
						Type: fmt.Sprintf("*models.%s", m.CreateTypeName()),
					},
				},
				Results: []*generators.Param{
					{
						Name: "",
						Type: fmt.Sprintf("*models.%s", m.ModelName()),
					},
					{
						Name: "",
						Type: "error",
					},
				},
			},
			{
				Name: "Delete",
				Args: []*generators.Param{
					{
						Name: "ctx",
						Type: "context.Context",
					},
					{
						Name: "id",
						Type: "models.UUID",
					},
				},
				Results: []*generators.Param{
					{
						Name: "",
						Type: "error",
					},
				},
			},
		},
	}
	if err := usecase.SyncInterface(); err != nil {
		return err
	}
	return nil
}

func SyncInterceptorInterface(m *models.ModelConfig) error {
	interceptor := &generators.Interface{
		Path:     filepath.Join("internal", "domain", "interceptors", m.FileName()),
		Name:     m.InterceptorTypeName(),
		Comments: nil,
		Methods: []*generators.Method{
			{
				Name: "Get",
				Args: []*generators.Param{
					{
						Name: "ctx",
						Type: "context.Context",
					},
					{
						Name: "id",
						Type: "models.UUID",
					},
				},
				Results: []*generators.Param{
					{
						Name: "",
						Type: fmt.Sprintf("*models.%s", m.ModelName()),
					},
					{
						Name: "",
						Type: "error",
					},
				},
			},
			{
				Name: "List",
				Args: []*generators.Param{
					{
						Name: "ctx",
						Type: "context.Context",
					},
					{
						Name: "filter",
						Type: fmt.Sprintf("*models.%s", m.FilterTypeName()),
					},
				},
				Results: []*generators.Param{
					{
						Name: "",
						Type: fmt.Sprintf("[]*models.%s", m.ModelName()),
					},
					{
						Name: "",
						Type: "uint64",
					},
					{
						Name: "",
						Type: "error",
					},
				},
			},
			{
				Name: "Update",
				Args: []*generators.Param{
					{
						Name: "ctx",
						Type: "context.Context",
					},
					{
						Name: "update",
						Type: fmt.Sprintf("*models.%s", m.UpdateTypeName()),
					},
				},
				Results: []*generators.Param{
					{
						Name: "",
						Type: fmt.Sprintf("*models.%s", m.ModelName()),
					},
					{
						Name: "",
						Type: "error",
					},
				},
			},
			{
				Name: "Create",
				Args: []*generators.Param{
					{
						Name: "ctx",
						Type: "context.Context",
					},
					{
						Name: "create",
						Type: fmt.Sprintf("*models.%s", m.CreateTypeName()),
					},
				},
				Results: []*generators.Param{
					{
						Name: "",
						Type: fmt.Sprintf("*models.%s", m.ModelName()),
					},
					{
						Name: "",
						Type: "error",
					},
				},
			},
			{
				Name: "Delete",
				Args: []*generators.Param{
					{
						Name: "ctx",
						Type: "context.Context",
					},
					{
						Name: "id",
						Type: "models.UUID",
					},
				},
				Results: []*generators.Param{
					{
						Name: "",
						Type: "error",
					},
				},
			},
		},
	}
	if m.Auth {
		for _, method := range interceptor.Methods {
			method.Args = append(method.Args, &generators.Param{
				Name: "requestUser",
				Type: "*models.User",
			})
		}
	}
	if err := interceptor.SyncInterface(); err != nil {
		return err
	}
	return nil
}

func SyncUseCaseImplementation(m *models.ModelConfig) error {
	useCase := &generators.UseCase{
		Path:  filepath.Join("internal", "usecases", m.FileName()),
		Name:  m.UseCaseTypeName(),
		Model: m,
		Params: []*generators.Param{
			{
				Name: m.RepositoryVariableName(),
				Type: fmt.Sprintf("repositories.%s", m.RepositoryTypeName()),
			},
			{
				Name: "clock",
				Type: "clock.Clock",
			},
			{
				Name: "logger",
				Type: "log.Logger",
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

func SyncRepositoryImplementation(m *models.ModelConfig) error {
	repository := &generators.Repository{
		Path:  filepath.Join("internal", "repositories", "postgres", m.FileName()),
		Name:  m.RepositoryTypeName(),
		Model: m,
		Params: []*generators.Param{
			{
				Name: "database",
				Type: "*sqlx.DB",
			},
			{
				Name: "logger",
				Type: "log.Logger",
			},
		},
	}
	if err := repository.SyncDTOStruct(); err != nil {
		return err
	}
	if err := repository.SyncDTOListType(); err != nil {
		return err
	}
	if err := repository.SyncDTOListToModels(); err != nil {
		return err
	}
	if err := repository.SyncDTOConstructor(); err != nil {
		return err
	}
	if err := repository.SyncDTOToModel(); err != nil {
		return err
	}
	if err := repository.SyncStruct(); err != nil {
		return err
	}
	if err := repository.SyncConstructor(); err != nil {
		return err
	}
	if err := repository.SyncCreateMethod(); err != nil {
		return err
	}
	if err := repository.SyncGetMethod(); err != nil {
		return err
	}
	if err := repository.SyncListMethod(); err != nil {
		return err
	}
	if err := repository.SyncCountMethod(); err != nil {
		return err
	}
	if err := repository.SyncUpdateMethod(); err != nil {
		return err
	}
	if err := repository.SyncDeleteMethod(); err != nil {
		return err
	}
	return nil
}

func SyncInterceptorImplementation(m *models.ModelConfig) error {
	interceptor := &generators.Interceptor{
		Path:  filepath.Join("internal", "interceptors", m.FileName()),
		Name:  m.InterceptorTypeName(),
		Model: m,
		Params: []*generators.Param{
			{
				Name: m.UseCaseTypeName(),
				Type: fmt.Sprintf("usecases.%s", m.UseCaseTypeName()),
			},
			{
				Name: "logger",
				Type: "log.Logger",
			},
		},
	}
	if m.Auth {
		interceptor.Params = append(
			interceptor.Params,
			&generators.Param{
				Name: "authUseCase",
				Type: "usecases.AuthUseCase",
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

func SyncModel(m *models.ModelConfig) error {
	if err := SyncModelStruct(m); err != nil {
		return err
	}
	if err := SyncCreateStruct(m); err != nil {
		return err
	}
	if err := SyncUpdateStruct(m); err != nil {
		return err
	}
	if err := SyncRepositoryInterface(m); err != nil {
		return err
	}
	if err := SyncUsecaseInterface(m); err != nil {
		return err
	}
	if err := SyncInterceptorInterface(m); err != nil {
		return err
	}
	if err := SyncUseCaseImplementation(m); err != nil {
		return err
	}
	if err := SyncRepositoryImplementation(m); err != nil {
		return err
	}
	if err := SyncInterceptorImplementation(m); err != nil {
		return err
	}
	return nil
}

func CreateCRUD(model *models.ModelConfig) error {
	if err := model.Validate(); err != nil {
		fmt.Printf("invalid model %s: %s\n", model.Model, err)
		return err
	}
	files := []*Template{
		{
			SourcePath:      "templates/internal/domain/models/crud.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "domain", "models", model.FileName()),
			Name:            "model",
		},
		{
			SourcePath:      "templates/internal/domain/models/crud_mock.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "domain", "models", "mock", model.FileName()),
			Name:            "model_mock",
		},
		{
			SourcePath:      "templates/internal/domain/repositories/crud.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "domain", "repositories", model.FileName()),
			Name:            "repository",
		},
		{
			SourcePath:      "templates/internal/domain/usecases/crud.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "domain", "usecases", model.FileName()),
			Name:            "usecase",
		},
		{
			SourcePath:      "templates/internal/domain/interceptors/crud.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "domain", "interceptors", model.FileName()),
			Name:            "interceptor",
		},
		{
			SourcePath:      "templates/internal/usecases/crud.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "usecases", model.FileName()),
			Name:            "usecase",
		},
		{
			SourcePath:      "templates/internal/usecases/crud_test.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "usecases", model.TestFileName()),
			Name:            "usecase test",
		},
		{
			SourcePath:      "templates/internal/interceptors/crud.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "interceptors", model.FileName()),
			Name:            "interceptor",
		},
		{
			SourcePath:      "templates/internal/interceptors/crud_test.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "interceptors", model.TestFileName()),
			Name:            "interceptor test",
		},
		{
			SourcePath:      "templates/internal/repositories/postgres/crud.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "repositories", "postgres", model.FileName()),
			Name:            "repository",
		},
		{
			SourcePath:      "templates/internal/repositories/postgres/crud_test.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "repositories", "postgres", model.TestFileName()),
			Name:            "repository test",
		},
		{
			SourcePath:      "templates/internal/interfaces/postgres/migrations/crud.up.sql.tmpl",
			DestinationPath: path.Join(destinationPath, "internal", "interfaces", "postgres", "migrations", model.MigrationUpFileName()),
			Name:            "migration up",
		},
		{
			SourcePath:      "templates/internal/interfaces/postgres/migrations/crud.down.sql.tmpl",
			DestinationPath: path.Join(destinationPath, "internal", "interfaces", "postgres", "migrations", model.MigrationDownFileName()),
			Name:            "migration down",
		},
	}
	if model.RESTEnabled {
		files = append(
			files,
			&Template{
				SourcePath:      "templates/internal/interfaces/rest/crud.go.tmpl",
				DestinationPath: path.Join(destinationPath, "internal", "interfaces", "rest", model.FileName()),
				Name:            "rest crud",
			},
		)
	}
	if model.GRPCEnabled {
		files = append(
			files,
			&Template{
				SourcePath:      "templates/internal/interfaces/grpc/crud.go.tmpl",
				DestinationPath: path.Join(destinationPath, "internal", "interfaces", "grpc", model.FileName()),
				Name:            "grpc service server",
			},
			&Template{
				SourcePath:      "templates/internal/interfaces/grpc/crud_test.go.tmpl",
				DestinationPath: path.Join(destinationPath, "internal", "interfaces", "grpc", model.TestFileName()),
				Name:            "test grpc service server",
			},
			&Template{
				SourcePath:      "templates/api/proto/service/v1/crud.proto.tmpl",
				DestinationPath: path.Join(destinationPath, "api", "proto", model.ProtoPackage, "v1", model.ProtoFileName()),
				Name:            "proto def",
			},
		)
	}
	for _, tmpl := range files {
		if err := tmpl.renderToFile(model); err != nil {
			return err
		}
	}
	if err := SyncModel(model); err != nil {
		return err
	}
	if err := addToDI("usecases", fmt.Sprintf("New%s", model.UseCaseTypeName())); err != nil {
		return err
	}
	if err := addToDI("interceptors", fmt.Sprintf("New%s", model.InterceptorTypeName())); err != nil {
		return err
	}
	if err := addToDI("postgresRepositories", fmt.Sprintf("New%s", model.RepositoryTypeName())); err != nil {
		return err
	}

	if model.RESTEnabled {
		if err := addToDI("restInterface", fmt.Sprintf("New%s", model.RESTHandlerTypeName())); err != nil {
			return err
		}
		if err := registerRESTHandler(model.RESTHandlerVariableName(), model.RESTHandlerTypeName()); err != nil {
			return err
		}
	}
	if model.GRPCEnabled {
		if err := addToDI("grpcInterface", fmt.Sprintf("New%s", model.GRPCHandlerTypeName())); err != nil {
			return err
		}
		if err := registerGRPCHandler(model.GRPCHandlerVariableName(), model.ProtoPackage, model.GRPCHandlerTypeName()); err != nil {
			return err
		}
	}
	if model.GRPCEnabled && model.GatewayEnabled {
		if err := registerGatewayHandler(model.ProtoPackage, model.GatewayHandlerTypeName()); err != nil {
			return err
		}
	}
	if model.Auth && model.ModelName() != "User" {
		if err := addPermission(model.PermissionIDList(), "objectAnybody"); err != nil {
			return err
		}
		if err := addPermission(model.PermissionIDDetail(), "objectAnybody"); err != nil {
			return err
		}
		if err := addPermission(model.PermissionIDCreate(), "objectAnybody"); err != nil {
			return err
		}
		if err := addPermission(model.PermissionIDUpdate(), "objectAnybody"); err != nil {
			return err
		}
		if err := addPermission(model.PermissionIDDelete(), "objectAnybody"); err != nil {
			return err
		}
	}
	return nil
}

func registerRESTHandler(variableName, typeName string) error {
	packagePath := filepath.Join(destinationPath, "internal", "interfaces", "rest")
	fileset := token.NewFileSet()
	tree, err := parser.ParseDir(fileset, packagePath, func(info fs.FileInfo) bool {
		return true
	}, parser.ParseComments)
	if err != nil {
		return err
	}
	for _, p := range tree {
		for filePath, file := range p.Files {
			for _, decl := range file.Decls {
				funcDecl, ok := decl.(*ast.FuncDecl)
				if ok {
					if funcDecl.Name.String() == "NewServer" {
						var exists bool
						for _, existedParam := range funcDecl.Type.Params.List {
							selector, ok := existedParam.Type.(*ast.StarExpr)
							if ok {
								t, ok := selector.X.(*ast.Ident)
								if ok && t.Name == typeName {
									exists = true
									break
								}
							}
						}
						if exists {
							continue
						}
						field := &ast.Field{
							Doc: &ast.CommentGroup{
								List: nil,
							},
							Names: []*ast.Ident{
								{
									NamePos: 0,
									Name:    variableName,
									Obj:     nil,
								},
							},
							Type: &ast.StarExpr{
								Star: 0,
								X: &ast.Ident{
									Name: typeName,
								},
							},
							Tag: nil,
							Comment: &ast.CommentGroup{
								List: nil,
							},
						}
						funcDecl.Type.Params.List = append(funcDecl.Type.Params.List, field)
						registerCall := &ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										NamePos: 0,
										Name:    variableName,
										Obj:     nil,
									},
									Sel: &ast.Ident{
										NamePos: 0,
										Name:    "Register",
										Obj:     nil,
									},
								},
								Lparen: 0,
								Args: []ast.Expr{
									&ast.Ident{
										NamePos: 0,
										Name:    "apiV1",
										Obj:     nil,
									},
								},
								Ellipsis: 0,
								Rparen:   0,
							},
						}
						le := len(funcDecl.Body.List)
						newBody := append(funcDecl.Body.List[:le-1], registerCall, funcDecl.Body.List[le-1])
						funcDecl.Body.List = newBody
						a := &bytes.Buffer{}
						if err := printer.Fprint(a, fileset, file); err != nil {
							return err
						}
						if err := os.WriteFile(filePath, a.Bytes(), 0777); err != nil {
							return err
						}
					}
				}
			}
		}
	}
	return nil
}

func registerGatewayHandler(typePackage, typeName string) error {
	packagePath := filepath.Join(destinationPath, "internal", "interfaces", "gateway")
	fileset := token.NewFileSet()
	tree, err := parser.ParseDir(fileset, packagePath, func(info fs.FileInfo) bool {
		return true
	}, parser.ParseComments)
	if err != nil {
		return err
	}
	for _, p := range tree {
		for filePath, file := range p.Files {
			for _, decl := range file.Decls {
				funcDecl, ok := decl.(*ast.FuncDecl)
				if ok {
					if funcDecl.Name.String() == "Start" {
						registerCall := &ast.AssignStmt{
							Lhs: []ast.Expr{ast.NewIdent("_")},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											NamePos: 0,
											Name:    typePackage,
											Obj:     nil,
										},
										Sel: &ast.Ident{
											NamePos: 0,
											Name:    typeName,
											Obj:     nil,
										},
									},
									Lparen: 0,
									Args: []ast.Expr{
										ast.NewIdent("ctx"),
										ast.NewIdent("mux"),
										ast.NewIdent("s.config.BindAddr"),
										ast.NewIdent("opts"),
									},
									Ellipsis: 0,
									Rparen:   0,
								},
							},
						}
						le := len(funcDecl.Body.List)
						newBody := append(funcDecl.Body.List[:le-2], registerCall, funcDecl.Body.List[le-2], funcDecl.Body.List[le-1])
						funcDecl.Body.List = newBody
						buff := &bytes.Buffer{}
						if err := printer.Fprint(buff, fileset, file); err != nil {
							return err
						}
						if err := os.WriteFile(filePath, buff.Bytes(), 0777); err != nil {
							return err
						}
					}
				}
			}
		}
	}
	return nil
}

func registerGRPCHandler(variableName, typePackage, typeName string) error {
	packagePath := filepath.Join(destinationPath, "internal", "interfaces", "grpc")
	fileset := token.NewFileSet()
	tree, err := parser.ParseDir(fileset, packagePath, func(info fs.FileInfo) bool {
		return true
	}, parser.ParseComments)
	if err != nil {
		return err
	}
	for _, p := range tree {
		for filePath, file := range p.Files {
			for _, decl := range file.Decls {
				funcDecl, ok := decl.(*ast.FuncDecl)
				if ok {
					if funcDecl.Name.String() == "NewServer" {
						var exists bool
						for _, existedParam := range funcDecl.Type.Params.List {
							selector, ok := existedParam.Type.(*ast.SelectorExpr)
							if ok && selector.Sel.Name == typeName {
								exists = true
								break
							}
						}
						if exists {
							continue
						}
						field := &ast.Field{
							Doc: &ast.CommentGroup{
								List: nil,
							},
							Names: []*ast.Ident{
								ast.NewIdent(variableName),
							},
							Type: &ast.SelectorExpr{
								X:   ast.NewIdent(typePackage),
								Sel: ast.NewIdent(typeName),
							},
							Tag: nil,
							Comment: &ast.CommentGroup{
								List: nil,
							},
						}
						_ = field
						funcDecl.Type.Params.List = append(funcDecl.Type.Params.List, field)
						registerCall := &ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										NamePos: 0,
										Name:    typePackage,
										Obj:     nil,
									},
									Sel: &ast.Ident{
										NamePos: 0,
										Name:    fmt.Sprintf("Register%s", typeName),
										Obj:     nil,
									},
								},
								Lparen: 0,
								Args: []ast.Expr{
									ast.NewIdent("server"),
									ast.NewIdent(variableName),
								},
								Ellipsis: 0,
								Rparen:   0,
							},
						}
						le := len(funcDecl.Body.List)
						newBody := append(funcDecl.Body.List[:le-1], registerCall, funcDecl.Body.List[le-1])
						funcDecl.Body.List = newBody
						buff := &bytes.Buffer{}
						if err := printer.Fprint(buff, fileset, file); err != nil {
							return err
						}
						if err := os.WriteFile(filePath, buff.Bytes(), 0777); err != nil {
							return err
						}
					}
				}
			}
		}
	}
	return nil
}

func addPermission(permission, check string) error {
	packagePath := filepath.Join(destinationPath, "internal", "repositories", "postgres")
	fileset := token.NewFileSet()
	tree, err := parser.ParseDir(fileset, packagePath, func(info fs.FileInfo) bool {
		return true
	}, parser.SkipObjectResolution)
	if err != nil {
		return err
	}
	for _, p := range tree {
		for filePath, file := range p.Files {
			for _, decl := range file.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if ok {
					for _, spec := range genDecl.Specs {
						variable, ok := spec.(*ast.ValueSpec)
						if ok {
							for _, name := range variable.Names {
								if name.Name == "hasObjectPermission" {
									for _, values := range variable.Values {
										lit, ok := values.(*ast.CompositeLit)
										if ok {
											var exists bool
											for _, elt := range lit.Elts {
												kv, ok := elt.(*ast.KeyValueExpr)
												if ok {
													selector, ok := kv.Key.(*ast.SelectorExpr)
													if ok && selector.Sel.Name == permission {
														exists = true
														break
													}
												}
											}
											if exists {
												continue
											}
											lit.Elts = append(lit.Elts, &ast.KeyValueExpr{
												Key: &ast.SelectorExpr{
													X:   ast.NewIdent("models"),
													Sel: ast.NewIdent(permission),
												},
												Colon: 0,
												Value: &ast.CompositeLit{
													Type:   nil,
													Lbrace: 0,
													Elts: []ast.Expr{
														ast.NewIdent(check),
													},
													Rbrace:     0,
													Incomplete: false,
												},
											})
											a := &bytes.Buffer{}
											if err := printer.Fprint(a, fileset, file); err != nil {
												return err
											}
											if err := os.WriteFile(filePath, a.Bytes(), 0777); err != nil {
												return err
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return nil
}
