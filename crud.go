package main

import (
	"bytes"
	"fmt"
	"github.com/018bf/creathor/models"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

func CreateCRUD(model *models.Model) error {
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
			SourcePath:      "templates/internal/interfaces/rest/crud.go.tmpl",
			DestinationPath: path.Join(destinationPath, "internal", "interfaces", "rest", model.FileName()),
			Name:            "rest crud",
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
		{
			SourcePath:      "templates/internal/interfaces/grpc/crud.go.tmpl",
			DestinationPath: path.Join(destinationPath, "internal", "interfaces", "grpc", model.FileName()),
			Name:            "grpc service server",
		},
		{
			SourcePath:      "templates/internal/interfaces/grpc/crud_test.go.tmpl",
			DestinationPath: path.Join(destinationPath, "internal", "interfaces", "grpc", model.TestFileName()),
			Name:            "test grpc service server",
		},
		{
			SourcePath:      "templates/api/proto/crud.proto.tmpl",
			DestinationPath: path.Join(destinationPath, "api", "proto", model.ProtoFileName()),
			Name:            "proto def",
		},
	}
	for _, tmpl := range files {
		if err := tmpl.renderToFile(model); err != nil {
			return err
		}
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
	if err := addToDI("restInterface", fmt.Sprintf("New%s", model.RESTHandlerTypeName())); err != nil {
		return err
	}
	if err := addToDI("grpcInterface", fmt.Sprintf("New%s", model.GRPCHandlerTypeName())); err != nil {
		return err
	}
	if err := registerRESTHandler(model.RESTHandlerVariableName(), model.RESTHandlerTypeName()); err != nil {
		return err
	}
	if err := registerGRPCHandler(model.RESTHandlerVariableName(), model.ProtoPackage, model.GRPCHandlerTypeName()); err != nil {
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
	packagePath := filepath.Join(destinationPath, "internal", "repositories")
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
