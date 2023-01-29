package main

import (
	"bytes"
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

func CreateDI(data *models.Project) error {
	directories := []string{
		path.Join(destinationPath, "internal", "containers"),
	}
	for _, directory := range directories {
		if err := os.MkdirAll(directory, 0777); err != nil {
			return NewUnexpectedBehaviorError(err.Error())
		}
	}
	files := []*Template{
		{
			SourcePath:      "templates/internal/containers/fx.go.tmpl",
			DestinationPath: path.Join(destinationPath, "internal", "containers", "fx.go"),
			Name:            "Uber FX DI container",
		},
	}
	for _, tmpl := range files {
		if err := tmpl.renderToFile(data); err != nil {
			return err
		}
	}
	return nil
}

func addToDI(packageName string, constructor string) error {
	packagePath := filepath.Join(destinationPath, "internal", "containers")
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
								if name.Name == "FXModule" {
									for _, values := range variable.Values {
										optionsFunc, ok := values.(*ast.CallExpr)
										if ok {
											for _, arg := range optionsFunc.Args {
												provideFunc, ok := arg.(*ast.CallExpr)
												if ok {
													fun, ok := provideFunc.Fun.(*ast.SelectorExpr)
													if ok && fun.Sel.Name == "Provide" {
														var exists bool
														for _, existedArg := range provideFunc.Args {
															selector, sOk := existedArg.(*ast.SelectorExpr)
															if sOk {
																ident, iOk := selector.X.(*ast.Ident)
																if iOk {
																	if ident.Name == packageName && selector.Sel.Name == constructor {
																		exists = true
																		break
																	}
																}
															}
														}
														if !exists {
															provideFunc.Args = append(provideFunc.Args, &ast.SelectorExpr{
																X:   ast.NewIdent(packageName),
																Sel: ast.NewIdent(constructor),
															})
														}
													} else {
														continue
													}
													break
												}
											}
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
					}
				}
			}
		}
	}
	return nil
}
