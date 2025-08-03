package usecases

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"

	"github.com/mikalai-mitsin/creathor/internal/pkg/app"
)

type InterfacesGenerator struct {
	domain *app.BaseEntity
}

func NewInterfacesGenerator(domain *app.BaseEntity) *InterfacesGenerator {
	return &InterfacesGenerator{domain: domain}
}

func (i InterfacesGenerator) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "app", i.domain.AppName(), "usecases", i.domain.DirName(), fmt.Sprintf("%s_interfaces.go", i.domain.SnakeName()))
	err := os.MkdirAll(path.Dir(filename), 0777)
	if err != nil {
		return err
	}
	file, err := parser.ParseFile(fileset, filename, nil, parser.AllErrors)
	if err != nil {
		file = i.file()
	}
	appServiceExists := false
	loggerExists := false
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok {
			if t.Name.String() == i.domain.GetServiceInterfaceName() {
				appServiceExists = true
			}
			if t.Name.String() == "logger" {
				loggerExists = true
			}
			return true
		}
		return true
	})
	if !appServiceExists {
		file.Decls = append(file.Decls, i.appServiceInterface())
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
	file := &ast.File{
		Name: ast.NewIdent("usecases"),
		Decls: []ast.Decl{
			i.imports(),
		},
	}
	return file
}

func (i InterfacesGenerator) imports() *ast.GenDecl {
	imports := &ast.GenDecl{
		Tok: token.IMPORT,
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{
					Slash: token.NoPos,
					Text:  fmt.Sprintf("//go:generate mockgen -source=%s_interfaces.go -package=usecases -destination=%s_interfaces_mock.go", i.domain.SnakeName(), i.domain.SnakeName()),
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
					Value: i.domain.EntitiesImportPath(),
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf(`"%s/internal/pkg/uuid"`, i.domain.Module),
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf(`"%s/internal/pkg/log"`, i.domain.Module),
				},
			},
		},
	}
	return imports
}

func (i InterfacesGenerator) appServiceInterface() *ast.GenDecl {
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
								Sel: ast.NewIdent(i.domain.GetCreateModel().Name),
							},
						},
					},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: &ast.SelectorExpr{
								X:   ast.NewIdent("entities"),
								Sel: ast.NewIdent(i.domain.GetMainModel().Name),
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
								Sel: ast.NewIdent(i.domain.GetMainModel().Name),
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
								Sel: ast.NewIdent(i.domain.GetFilterModel().Name),
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
									Sel: ast.NewIdent(i.domain.GetMainModel().Name),
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
								Sel: ast.NewIdent(i.domain.GetUpdateModel().Name),
							},
						},
					},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: &ast.SelectorExpr{
								X:   ast.NewIdent("entities"),
								Sel: ast.NewIdent(i.domain.GetMainModel().Name),
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
								X:   ast.NewIdent("uuid"),
								Sel: ast.NewIdent("UUID"),
							},
						},
					},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
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
				Name: ast.NewIdent(i.domain.GetServiceInterfaceName()),
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
								Names: []*ast.Ident{
									ast.NewIdent("Debug"),
								},
								Type: &ast.FuncType{
									Params: &ast.FieldList{
										List: []*ast.Field{
											{
												Names: []*ast.Ident{
													ast.NewIdent("msg"),
												},
												Type: ast.NewIdent("string"),
											},
											{
												Names: []*ast.Ident{
													ast.NewIdent("fields"),
												},
												Type: &ast.Ellipsis{
													Elt: &ast.SelectorExpr{
														X:   ast.NewIdent("log"),
														Sel: ast.NewIdent("Field"),
													},
												},
											},
										},
									},
								},
							},
							{
								Names: []*ast.Ident{
									ast.NewIdent("Info"),
								},
								Type: &ast.FuncType{
									Params: &ast.FieldList{
										List: []*ast.Field{
											{
												Names: []*ast.Ident{
													ast.NewIdent("msg"),
												},
												Type: ast.NewIdent("string"),
											},
											{
												Names: []*ast.Ident{
													ast.NewIdent("fields"),
												},
												Type: &ast.Ellipsis{
													Elt: &ast.SelectorExpr{
														X:   ast.NewIdent("log"),
														Sel: ast.NewIdent("Field"),
													},
												},
											},
										},
									},
								},
							},
							{
								Names: []*ast.Ident{
									ast.NewIdent("Print"),
								},
								Type: &ast.FuncType{
									Params: &ast.FieldList{
										List: []*ast.Field{
											{
												Names: []*ast.Ident{
													ast.NewIdent("msg"),
												},
												Type: ast.NewIdent("string"),
											},
											{
												Names: []*ast.Ident{
													ast.NewIdent("fields"),
												},
												Type: &ast.Ellipsis{
													Elt: &ast.SelectorExpr{
														X:   ast.NewIdent("log"),
														Sel: ast.NewIdent("Field"),
													},
												},
											},
										},
									},
								},
							},
							{
								Names: []*ast.Ident{
									ast.NewIdent("Warn"),
								},
								Type: &ast.FuncType{
									Params: &ast.FieldList{
										List: []*ast.Field{
											{
												Names: []*ast.Ident{
													ast.NewIdent("msg"),
												},
												Type: ast.NewIdent("string"),
											},
											{
												Names: []*ast.Ident{
													ast.NewIdent("fields"),
												},
												Type: &ast.Ellipsis{
													Elt: &ast.SelectorExpr{
														X:   ast.NewIdent("log"),
														Sel: ast.NewIdent("Field"),
													},
												},
											},
										},
									},
								},
							},
							{
								Names: []*ast.Ident{
									ast.NewIdent("Error"),
								},
								Type: &ast.FuncType{
									Params: &ast.FieldList{
										List: []*ast.Field{
											{
												Names: []*ast.Ident{
													ast.NewIdent("msg"),
												},
												Type: ast.NewIdent("string"),
											},
											{
												Names: []*ast.Ident{
													ast.NewIdent("fields"),
												},
												Type: &ast.Ellipsis{
													Elt: &ast.SelectorExpr{
														X:   ast.NewIdent("log"),
														Sel: ast.NewIdent("Field"),
													},
												},
											},
										},
									},
								},
							},
							{
								Names: []*ast.Ident{
									ast.NewIdent("Fatal"),
								},
								Type: &ast.FuncType{
									Params: &ast.FieldList{
										List: []*ast.Field{
											{
												Names: []*ast.Ident{
													ast.NewIdent("msg"),
												},
												Type: ast.NewIdent("string"),
											},
											{
												Names: []*ast.Ident{
													ast.NewIdent("fields"),
												},
												Type: &ast.Ellipsis{
													Elt: &ast.SelectorExpr{
														X:   ast.NewIdent("log"),
														Sel: ast.NewIdent("Field"),
													},
												},
											},
										},
									},
								},
							},
							{
								Names: []*ast.Ident{
									ast.NewIdent("Panic"),
								},
								Type: &ast.FuncType{
									Params: &ast.FieldList{
										List: []*ast.Field{
											{
												Names: []*ast.Ident{
													ast.NewIdent("msg"),
												},
												Type: ast.NewIdent("string"),
											},
											{
												Names: []*ast.Ident{
													ast.NewIdent("fields"),
												},
												Type: &ast.Ellipsis{
													Elt: &ast.SelectorExpr{
														X:   ast.NewIdent("log"),
														Sel: ast.NewIdent("Field"),
													},
												},
											},
										},
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
