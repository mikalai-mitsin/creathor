package main

import (
	"bytes"
	"fmt"
	"github.com/018bf/creathor/internal/configs"
	generatorsIntercepstorInterfaces "github.com/018bf/creathor/internal/generators/domain/interceptors"
	generatorsModels "github.com/018bf/creathor/internal/generators/domain/models"
	generatorsRepositoriesInterfaces "github.com/018bf/creathor/internal/generators/domain/repositories"
	generatorsUseCasesInterfaces "github.com/018bf/creathor/internal/generators/domain/usecases"
	"github.com/018bf/creathor/internal/generators/interceptors"
	"github.com/018bf/creathor/internal/generators/repositories/postgres"
	"github.com/018bf/creathor/internal/generators/usecases"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

func SyncModelStruct(m *configs.ModelConfig) error {
	model := &generatorsModels.Model{
		Name:        m.ModelName(),
		ModelConfig: m,
		Params: []*generatorsModels.Param{
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
			&generatorsModels.Param{
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

func SyncFilterStruct(m *configs.ModelConfig) error {
	create := &generatorsModels.Model{
		Name:        m.FilterTypeName(),
		ModelConfig: m,
		Params: []*generatorsModels.Param{
			{
				Name: "IDs",
				Type: "[]UUID",
			},
			{
				Name: "PageSize",
				Type: "*uint64",
			},
			{
				Name: "PageNumber",
				Type: "*uint64",
			},
			{
				Name: "OrderBy",
				Type: "[]string",
			},
		},
	}
	if m.SearchEnabled() {
		create.Params = append(
			create.Params,
			&generatorsModels.Param{
				Name: "Search",
				Type: "*string",
			},
		)
	}
	if err := create.Sync(); err != nil {
		return err
	}
	return nil
}

func SyncCreateStruct(m *configs.ModelConfig) error {
	create := &generatorsModels.Model{
		Name:        m.CreateTypeName(),
		ModelConfig: m,
		Params:      []*generatorsModels.Param{},
	}
	for _, param := range m.Params {
		create.Params = append(create.Params, &generatorsModels.Param{
			Name: param.GetName(),
			Type: param.Type,
		})
	}
	if err := create.Sync(); err != nil {
		return err
	}
	return nil
}

func SyncUpdateStruct(m *configs.ModelConfig) error {
	update := &generatorsModels.Model{
		Name:        m.UpdateTypeName(),
		ModelConfig: m,
		Params: []*generatorsModels.Param{
			{
				Name: "ID",
				Type: "UUID",
			},
		},
	}
	for _, param := range m.Params {
		update.Params = append(update.Params, &generatorsModels.Param{
			Name: param.GetName(),
			Type: fmt.Sprintf("*%s", param.Type),
		})
	}
	if err := update.Sync(); err != nil {
		return err
	}
	return nil
}

func SyncUseCaseImplementation(m *configs.ModelConfig) error {
	useCase := &usecases.UseCase{
		Model: m,
	}
	if err := useCase.Sync(); err != nil {
		return err
	}
	return nil
}

func SyncRepositoryImplementation(m *configs.ModelConfig) error {
	repository := &postgres.Repository{
		Path:  filepath.Join("internal", "repositories", "postgres", m.FileName()),
		Name:  m.RepositoryTypeName(),
		Model: m,
		Params: []*generatorsModels.Param{
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

func SyncInterceptorImplementation(m *configs.ModelConfig) error {
	interceptor := &interceptors.Interceptor{
		Path:  filepath.Join("internal", "interceptors", m.FileName()),
		Name:  m.InterceptorTypeName(),
		Model: m,
		Params: []*generatorsModels.Param{
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
			&generatorsModels.Param{
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

func SyncModel(m *configs.ModelConfig) error {
	if err := SyncModelStruct(m); err != nil {
		return err
	}
	if err := SyncFilterStruct(m); err != nil {
		return err
	}
	if err := SyncCreateStruct(m); err != nil {
		return err
	}
	if err := SyncUpdateStruct(m); err != nil {
		return err
	}
	repositoryInterface := generatorsRepositoriesInterfaces.RepositoryInterface{Config: m}
	if err := repositoryInterface.Sync(); err != nil {
		return err
	}
	useCaseInterface := generatorsUseCasesInterfaces.UseCaseInterface{Config: m}
	if err := useCaseInterface.Sync(); err != nil {
		return err
	}
	interceptorInterface := generatorsIntercepstorInterfaces.InterceptorInterface{Config: m}
	if err := interceptorInterface.Sync(); err != nil {
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

func CreateCRUD(model *configs.ModelConfig) error {
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
