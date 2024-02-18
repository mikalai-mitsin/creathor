package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/018bf/creathor/internal/configs"
)

func CreateCRUD(model *configs.DomainConfig) error {
	if err := model.Validate(); err != nil {
		fmt.Printf("invalid model %s: %s\n", model.Model, err)
		return err
	}
	files := []*Template{

		{
			SourcePath: "templates/internal/domain/usecases/crud_test.go.tmpl",
			DestinationPath: filepath.Join(
				destinationPath,
				"internal",
				model.Model,
				"usecases",
				model.TestFileName(),
			),
			Name: "usecase test",
		},

		{
			SourcePath: "templates/internal/domain/interceptors/crud_test.go.tmpl",
			DestinationPath: filepath.Join(
				destinationPath,
				"internal",
				model.Model,
				"interceptors",
				model.TestFileName(),
			),
			Name: "interceptor test",
		},

		{
			SourcePath: "templates/internal/domain/repositories/postgres/crud_test.go.tmpl",
			DestinationPath: filepath.Join(
				destinationPath,
				"internal",
				model.Model,
				"repositories",
				"postgres",
				model.TestFileName(),
			),
			Name: "repository test",
		},
		{
			SourcePath: "templates/internal/interfaces/postgres/migrations/crud.up.sql.tmpl",
			DestinationPath: path.Join(
				destinationPath,
				"internal",
				"interfaces",
				"postgres",
				"migrations",
				model.MigrationUpFileName(),
			),
			Name: "migration up",
		},
		{
			SourcePath: "templates/internal/interfaces/postgres/migrations/crud.down.sql.tmpl",
			DestinationPath: path.Join(
				destinationPath,
				"internal",
				"interfaces",
				"postgres",
				"migrations",
				model.MigrationDownFileName(),
			),
			Name: "migration down",
		},
	}
	if model.GRPCEnabled {
		files = append(
			files,
			&Template{
				SourcePath: "templates/internal/domain/handlers/grpc/crud_test.go.tmpl",
				DestinationPath: path.Join(
					destinationPath,
					"internal",
					model.Model,
					"handlers",
					"grpc",
					model.TestFileName(),
				),
				Name: "test grpc service server",
			},
			&Template{
				SourcePath: "templates/api/proto/service/v1/crud.proto.tmpl",
				DestinationPath: path.Join(
					destinationPath,
					"api",
					"proto",
					model.ProtoPackage,
					"v1",
					model.ProtoFileName(),
				),
				Name: "proto def",
			},
		)
	}
	for _, tmpl := range files {
		if err := tmpl.renderToFile(model); err != nil {
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

func addPermission(permission, check string) error {
	packagePath := filepath.Join(destinationPath, "internal", "user", "repositories", "postgres")
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
