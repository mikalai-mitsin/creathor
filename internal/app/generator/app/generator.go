package app

import (
	"bytes"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/app/repositories/postgres"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/app/services"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/app/usecases"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/mikalai-mitsin/creathor/internal/app/generator"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/app/entities"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/app/handlers/grpc"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/app/handlers/http"
	"github.com/mikalai-mitsin/creathor/internal/pkg/app"
)

type Generator struct {
	domain *app.App
}

func NewGenerator(d *app.App) *Generator {
	return &Generator{domain: d}
}

func (g *Generator) Sync() error {
	domainGenerators := []generator.Generator{NewApp(g.domain)}
	for _, entity := range g.domain.Entities {
		domainGenerators = append(domainGenerators,
			usecases.NewInterfacesGenerator(entity),
			usecases.NewUseCaseGenerator(entity),
			usecases.NewTestGenerator(entity),

			services.NewInterfacesGenerator(entity),
			services.NewServiceGenerator(entity),
			services.NewTestGenerator(entity),

			postgres.NewInterfacesGenerator(entity),
			postgres.NewRepositoryGenerator(entity),
			postgres.NewTestGenerator(entity),
		)
		if g.domain.Config.HTTPEnabled {
			domainGenerators = append(
				domainGenerators,
				http.NewDTOGenerator(entity),
				http.NewHandlerGenerator(entity),
				http.NewInterfacesGenerator(entity),
			)
		}
		if g.domain.Config.GRPCEnabled {
			domainGenerators = append(
				domainGenerators,
				grpc.NewProtoGenerator(entity),
				grpc.NewInterfacesGenerator(entity),
				grpc.NewHandlerGenerator(entity),
				grpc.NewTestGenerator(entity),
			)
		}
		for _, baseEntity := range entity.Entities {
			domainGenerators = append(domainGenerators, entities.NewModel(baseEntity, entity))
		}
		if g.domain.Auth && entity.CamelName() != "User" {
			if err := addPermission(entity.PermissionIDList(), "objectAnybody"); err != nil {
				return err
			}
			if err := addPermission(entity.PermissionIDDetail(), "objectAnybody"); err != nil {
				return err
			}
			if err := addPermission(entity.PermissionIDCreate(), "objectAnybody"); err != nil {
				return err
			}
			if err := addPermission(entity.PermissionIDUpdate(), "objectAnybody"); err != nil {
				return err
			}
			if err := addPermission(entity.PermissionIDDelete(), "objectAnybody"); err != nil {
				return err
			}
		}
	}
	for _, domainGenerator := range domainGenerators {
		if err := domainGenerator.Sync(); err != nil {
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
													X:   ast.NewIdent("entities"),
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
