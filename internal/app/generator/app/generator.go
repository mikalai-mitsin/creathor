package app

import (
	"bytes"
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
	"github.com/mikalai-mitsin/creathor/internal/app/generator/app/repositories/postgres"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/app/services"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/app/usecases"
	"github.com/mikalai-mitsin/creathor/internal/pkg/domain"
)

type Generator struct {
	domain *domain.Domain
}

func NewGenerator(d *domain.Domain) *Generator {
	return &Generator{domain: d}
}

func (g *Generator) Sync() error {
	domainGenerators := []generator.Generator{
		usecases.NewInterfacesGenerator(g.domain),
		usecases.NewUseCaseGenerator(g.domain),
		usecases.NewTestGenerator(g.domain),

		services.NewInterfacesGenerator(g.domain),
		services.NewServiceGenerator(g.domain),
		services.NewTestGenerator(g.domain),

		postgres.NewInterfacesGenerator(g.domain),
		postgres.NewRepositoryGenerator(g.domain),
		postgres.NewTestGenerator(g.domain),

		NewApp(g.domain),
	}
	if g.domain.Config.HTTPEnabled {
		domainGenerators = append(
			domainGenerators,
			http.NewDTOGenerator(g.domain),
			http.NewHandlerGenerator(g.domain),
			http.NewInterfacesGenerator(g.domain),
		)
	}
	if g.domain.Config.GRPCEnabled {
		domainGenerators = append(
			domainGenerators,
			grpc.NewProtoGenerator(g.domain),
			grpc.NewInterfacesGenerator(g.domain),
			grpc.NewHandlerGenerator(g.domain),
			grpc.NewTestGenerator(g.domain),
		)
	}
	for _, model := range g.domain.Entities {
		domainGenerators = append(domainGenerators, entities.NewModel(model, g.domain))
	}
	for _, domainGenerator := range domainGenerators {
		if err := domainGenerator.Sync(); err != nil {
			return err
		}
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
