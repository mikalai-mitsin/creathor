package app

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

	"github.com/018bf/creathor/internal/app/generator"
	"github.com/018bf/creathor/internal/app/generator/app/interceptors"
	"github.com/018bf/creathor/internal/app/generator/app/interfaces/grpc"
	"github.com/018bf/creathor/internal/app/generator/app/models"
	"github.com/018bf/creathor/internal/app/generator/app/repositories/postgres"
	"github.com/018bf/creathor/internal/app/generator/app/usecases"
	"github.com/018bf/creathor/internal/pkg/domain"
	"github.com/018bf/creathor/internal/pkg/tmpl"
)

type Generator struct {
	domain *domain.Domain
}

func NewGenerator(d *domain.Domain) *Generator {
	return &Generator{domain: d}
}

func (g *Generator) Sync() error {
	domainGenerators := []generator.Generator{
		interceptors.NewInterceptorCrud(g.domain),
		interceptors.NewUseCaseInterfaceCrud(g.domain),

		usecases.NewUseCaseCrud(g.domain),
		usecases.NewRepositoryInterfaceCrud(g.domain),

		postgres.NewRepositoryCrud(g.domain),

		grpc.NewHandler(g.domain),
		grpc.NewInterceptorInterfaceCrud(g.domain),
	}
	for _, model := range g.domain.Models {
		domainGenerators = append(domainGenerators, models.NewModel(model, g.domain))
	}
	for _, domainGenerator := range domainGenerators {
		if err := domainGenerator.Sync(); err != nil {
			return err
		}
	}
	if err := renderTemplates(g.domain); err != nil {
		return err
	}
	if g.domain.Auth && g.domain.CamelName() != "User" {
		if err := addPermission(g.domain.PermissionIDList(), "objectAnybody"); err != nil {
			return err
		}
		if err := addPermission(g.domain.PermissionIDDetail(), "objectAnybody"); err != nil {
			return err
		}
		if err := addPermission(g.domain.PermissionIDCreate(), "objectAnybody"); err != nil {
			return err
		}
		if err := addPermission(g.domain.PermissionIDUpdate(), "objectAnybody"); err != nil {
			return err
		}
		if err := addPermission(g.domain.PermissionIDDelete(), "objectAnybody"); err != nil {
			return err
		}
	}
	return nil
}

func addPermission(permission, check string) error {
	packagePath := filepath.Join(
		destinationPath,
		"internal",
		"app",
		"user",
		"repositories",
		"postgres",
	)
	if err := os.MkdirAll(packagePath, 0777); err != nil {
		return err
	}
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

var destinationPath = "."

func renderTemplates(domain *domain.Domain) error {
	if err := renderMigrations(domain); err != nil {
		return err
	}
	if err := renderTests(domain); err != nil {
		return err
	}
	if domain.Config.GRPCEnabled {
		if err := renderGrpc(domain); err != nil {
			return err
		}
	}
	return nil
}

func renderGrpc(domain *domain.Domain) error {
	files := []*tmpl.Template{
		{
			SourcePath: "templates/internal/domain/handlers/grpc/crud_test.go.tmpl",
			DestinationPath: path.Join(
				destinationPath,
				"internal",
				"app",
				domain.DirName(),
				"handlers",
				"grpc",
				domain.TestFileName(),
			),
			Name: "test grpc service server",
		},
		{
			SourcePath: "templates/api/proto/service/v1/crud.proto.tmpl",
			DestinationPath: path.Join(
				destinationPath,
				"api",
				"proto",
				domain.ProtoModule,
				"v1",
				fmt.Sprintf("%s.proto", domain.SnakeName()),
			),
			Name: "proto def",
		},
	}
	for _, template := range files {
		if err := template.RenderToFile(domain); err != nil {
			return err
		}
	}
	return nil
}

func renderMigrations(domain *domain.Domain) error {
	files := []*tmpl.Template{
		{
			SourcePath: "templates/internal/pkg/postgres/migrations/crud.up.sql.tmpl",
			DestinationPath: path.Join(
				destinationPath,
				"internal",
				"pkg",
				"postgres",
				"migrations",
				domain.MigrationUpFileName(),
			),
			Name: "migration up",
		},
		{
			SourcePath: "templates/internal/pkg/postgres/migrations/crud.down.sql.tmpl",
			DestinationPath: path.Join(
				destinationPath,
				"internal",
				"pkg",
				"postgres",
				"migrations",
				domain.MigrationDownFileName(),
			),
			Name: "migration down",
		},
	}
	for _, template := range files {
		if err := template.RenderToFile(domain); err != nil {
			return err
		}
	}
	return nil
}

func renderTests(domain *domain.Domain) error {
	files := []*tmpl.Template{
		{
			SourcePath: "templates/internal/domain/usecases/crud_test.go.tmpl",
			DestinationPath: filepath.Join(
				destinationPath,
				"internal",
				"app",
				domain.DirName(),
				"usecases",
				domain.TestFileName(),
			),
			Name: "usecase test",
		},
		{
			SourcePath: "templates/internal/domain/interceptors/crud_test.go.tmpl",
			DestinationPath: filepath.Join(
				destinationPath,
				"internal",
				"app",
				domain.DirName(),
				"interceptors",
				domain.TestFileName(),
			),
			Name: "interceptor test",
		},
		{
			SourcePath: "templates/internal/domain/repositories/postgres/crud_test.go.tmpl",
			DestinationPath: filepath.Join(
				destinationPath,
				"internal",
				"app",
				domain.DirName(),
				"repositories",
				"postgres",
				domain.TestFileName(),
			),
			Name: "repository test",
		},
		{
			SourcePath: "templates/internal/pkg/postgres/migrations/crud.up.sql.tmpl",
			DestinationPath: path.Join(
				destinationPath,
				"internal",
				"pkg",
				"postgres",
				"migrations",
				domain.MigrationUpFileName(),
			),
			Name: "migration up",
		},
		{
			SourcePath: "templates/internal/pkg/postgres/migrations/crud.down.sql.tmpl",
			DestinationPath: path.Join(
				destinationPath,
				"internal",
				"pkg",
				"postgres",
				"migrations",
				domain.MigrationDownFileName(),
			),
			Name: "migration down",
		},
	}
	for _, template := range files {
		if err := template.RenderToFile(domain); err != nil {
			return err
		}
	}
	return nil
}
