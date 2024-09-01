package usecases

import (
	"bytes"
	"fmt"
	"github.com/mikalai-mitsin/creathor/internal/pkg/domain"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"
)

type UseCaseInterfaces struct {
	domain *domain.Domain
}

func NewUseCaseInterfaces(domain *domain.Domain) *UseCaseInterfaces {
	return &UseCaseInterfaces{domain: domain}
}

func (i UseCaseInterfaces) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "app", i.domain.DirName(), "usecases", "interfaces.go")
	err := os.MkdirAll(path.Dir(filename), 0777)
	if err != nil {
		return err
	}
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = i.file()
	}
	appServiceExists := false
	loggerExists := false
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok {
			if t.Name.String() == i.domain.Service.Name {
				appServiceExists = true
			}
			if t.Name.String() == "Logger" {
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

func (i UseCaseInterfaces) file() *ast.File {
	file := &ast.File{
		Name: ast.NewIdent("usecases"),
		Decls: []ast.Decl{
			i.imports(),
		},
	}
	return file
}

func (i UseCaseInterfaces) imports() *ast.GenDecl {
	imports := &ast.GenDecl{
		Tok: token.IMPORT,
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
	if i.domain.Auth {
		imports.Specs = append(imports.Specs, &ast.ImportSpec{
			Name: ast.NewIdent("userEntities"),
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s/internal/app/user/entities"`, i.domain.Module),
			},
		})
	}
	return imports
}

func (i UseCaseInterfaces) appServiceInterface() *ast.GenDecl {
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
							Type: &ast.StarExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("entities"),
									Sel: ast.NewIdent(i.domain.GetCreateModel().Name),
								},
							},
						},
					},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: &ast.StarExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("entities"),
									Sel: ast.NewIdent(i.domain.GetMainModel().Name),
								},
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
							Type: &ast.StarExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("entities"),
									Sel: ast.NewIdent(i.domain.GetMainModel().Name),
								},
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
							Type: &ast.StarExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("entities"),
									Sel: ast.NewIdent(i.domain.GetFilterModel().Name),
								},
							},
						},
					},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: &ast.ArrayType{
								Elt: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X:   ast.NewIdent("entities"),
										Sel: ast.NewIdent(i.domain.GetMainModel().Name),
									},
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
							Type: &ast.StarExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("entities"),
									Sel: ast.NewIdent(i.domain.GetUpdateModel().Name),
								},
							},
						},
					},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: &ast.StarExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("entities"),
									Sel: ast.NewIdent(i.domain.GetMainModel().Name),
								},
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
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{
					Text: fmt.Sprintf(
						"//%s - domain layer use case interface",
						i.domain.Service.Name,
					),
				},
				{
					Text: fmt.Sprintf(
						"//go:generate mockgen -build_flags=-mod=mod -destination mock/service.go . %s",
						i.domain.Service.Name,
					),
				},
			},
		},
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: &ast.Ident{
					Name: i.domain.Service.Name,
				},
				Type: &ast.InterfaceType{
					Methods: &ast.FieldList{
						List: methods,
					},
				},
			},
		},
	}
}

func (i UseCaseInterfaces) loggerInterface() *ast.GenDecl {
	return &ast.GenDecl{
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{
					Text: "//Logger - base logger interface",
				},
				{
					Text: "//go:generate mockgen -build_flags=-mod=mod -destination mock/logger.go . Logger",
				},
			},
		},
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent("Logger"),
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
