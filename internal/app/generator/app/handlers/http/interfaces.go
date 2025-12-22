package http

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"

	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
)

type InterfacesGenerator struct {
	entityConfig configs.EntityConfig
}

func NewInterfacesGenerator(entityConfig configs.EntityConfig) *InterfacesGenerator {
	return &InterfacesGenerator{entityConfig: entityConfig}
}

func (i InterfacesGenerator) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join(
		"internal",
		"app",
		i.entityConfig.AppConfig.AppName(),
		"handlers",
		"http",
		i.entityConfig.DirName(),
		fmt.Sprintf("%s_interfaces.go", i.entityConfig.SnakeName()),
	)
	err := os.MkdirAll(path.Dir(filename), 0777)
	if err != nil {
		return err
	}
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = i.file()
	}
	var loggerExists bool
	var usecaseExists bool
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok {
			if t.Name.String() == i.entityConfig.GetUseCaseInterfaceName() {
				usecaseExists = true
			}
			if t.Name.String() == "logger" {
				loggerExists = true
			}
			return true
		}
		return true
	})
	if !usecaseExists {
		file.Decls = append(file.Decls, i.usecaseInterface())
	}
	if !loggerExists {
		file.Decls = append(file.Decls, i.loggerInterface())
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(filename, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (i InterfacesGenerator) file() *ast.File {
	return &ast.File{
		Name: ast.NewIdent("handlers"),
		Decls: []ast.Decl{
			i.imports(),
		},
	}
}

func (i InterfacesGenerator) imports() *ast.GenDecl {
	return &ast.GenDecl{
		Tok: token.IMPORT,
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{
					Slash: token.NoPos,
					Text:  "//go:generate mockgen -package=$GOPACKAGE -source=$GOFILE -destination=mock.go",
				},
			},
		},
		Specs: []ast.Spec{
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"context"`,
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: i.entityConfig.ImportPathEntities(),
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: i.entityConfig.AppConfig.ProjectConfig.UUIDImportPath(),
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: i.entityConfig.AppConfig.ProjectConfig.LogImportPath(),
				},
			},
		},
	}
}

func (i InterfacesGenerator) usecaseInterface() *ast.GenDecl {
	methods := []*ast.Field{
		{
			Names: []*ast.Ident{ast.NewIdent("Create")},
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: &ast.SelectorExpr{
								X:   ast.NewIdent("context"),
								Sel: ast.NewIdent("Context"),
							},
						},
						{
							Type: &ast.SelectorExpr{
								X:   ast.NewIdent("entities"),
								Sel: ast.NewIdent(i.entityConfig.GetCreateModel().Name),
							},
						},
					},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: &ast.SelectorExpr{
								X:   ast.NewIdent("entities"),
								Sel: ast.NewIdent(i.entityConfig.GetMainModel().Name),
							},
						},
						{
							Type: ast.NewIdent("error"),
						},
					},
				},
			},
		},
		{
			Names: []*ast.Ident{ast.NewIdent("Get")},
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: &ast.SelectorExpr{
								X:   ast.NewIdent("context"),
								Sel: ast.NewIdent("Context"),
							},
						},
						{
							Type: &ast.SelectorExpr{
								X:   ast.NewIdent("uuid"),
								Sel: ast.NewIdent("UUID"),
							},
						},
					},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: &ast.SelectorExpr{
								X:   ast.NewIdent("entities"),
								Sel: ast.NewIdent(i.entityConfig.GetMainModel().Name),
							},
						},
						{
							Type: ast.NewIdent("error"),
						},
					},
				},
			},
		},
		{
			Names: []*ast.Ident{ast.NewIdent("List")},
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: &ast.SelectorExpr{
								X:   ast.NewIdent("context"),
								Sel: ast.NewIdent("Context"),
							},
						},
						{
							Type: &ast.SelectorExpr{
								X:   ast.NewIdent("entities"),
								Sel: ast.NewIdent(i.entityConfig.GetFilterModel().Name),
							},
						},
					},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: &ast.ArrayType{
								Elt: &ast.SelectorExpr{
									X:   ast.NewIdent("entities"),
									Sel: ast.NewIdent(i.entityConfig.GetMainModel().Name),
								},
							},
						},
						{
							Type: ast.NewIdent("uint64"),
						},
						{
							Type: ast.NewIdent("error"),
						},
					},
				},
			},
		},
		{
			Names: []*ast.Ident{ast.NewIdent("Update")},
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: &ast.SelectorExpr{
								X:   ast.NewIdent("context"),
								Sel: ast.NewIdent("Context"),
							},
						},
						{
							Type: &ast.SelectorExpr{
								X:   ast.NewIdent("entities"),
								Sel: ast.NewIdent(i.entityConfig.GetUpdateModel().Name),
							},
						},
					},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: &ast.SelectorExpr{
								X:   ast.NewIdent("entities"),
								Sel: ast.NewIdent(i.entityConfig.GetMainModel().Name),
							},
						},
						{
							Type: ast.NewIdent("error"),
						},
					},
				},
			},
		},
		{
			Names: []*ast.Ident{ast.NewIdent("Delete")},
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: &ast.SelectorExpr{
								X:   ast.NewIdent("context"),
								Sel: ast.NewIdent("Context"),
							},
						},
						{
							Type: &ast.SelectorExpr{
								X:   ast.NewIdent("entities"),
								Sel: ast.NewIdent(i.entityConfig.GetDeleteModel().Name),
							},
						},
					},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: &ast.SelectorExpr{
								X:   ast.NewIdent("entities"),
								Sel: ast.NewIdent(i.entityConfig.GetMainModel().Name),
							},
						},
						{
							Type: ast.NewIdent("error"),
						},
					},
				},
			},
		},
	}
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(i.entityConfig.GetUseCaseInterfaceName()),
				Type: &ast.InterfaceType{
					Methods: &ast.FieldList{
						List: methods,
					},
				},
			},
		},
	}
}

func (i InterfacesGenerator) loggerInterface() *ast.GenDecl {
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent("logger"),
				Type: &ast.InterfaceType{
					Methods: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "log",
									},
									Sel: &ast.Ident{
										Name: "Logger",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
