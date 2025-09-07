package usecases

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"

	"github.com/mikalai-mitsin/creathor/internal/pkg/astfile"
	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
)

type InterfacesGenerator struct {
	domain *configs.EntityConfig
}

func NewInterfacesGenerator(domain *configs.EntityConfig) *InterfacesGenerator {
	return &InterfacesGenerator{domain: domain}
}

func (i InterfacesGenerator) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join(
		"internal",
		"app",
		i.domain.AppConfig.AppName(),
		"usecases",
		i.domain.DirName(),
		"interfaces.go",
	)
	err := os.MkdirAll(path.Dir(filename), 0777)
	if err != nil {
		return err
	}
	file, err := parser.ParseFile(fileset, filename, nil, parser.AllErrors)
	if err != nil {
		file = i.file()
	}
	if !astfile.TypeExists(file, i.domain.GetServiceInterfaceName()) {
		file.Decls = append(file.Decls, i.appServiceInterface())
	}
	if !astfile.TypeExists(file, i.domain.EventProducerInterfaceName()) &&
		i.domain.AppConfig.ProjectConfig.KafkaEnabled {
		file.Decls = append(file.Decls, i.appEventServiceInterface())
	}
	if !astfile.TypeExists(file, "logger") {
		file.Decls = append(file.Decls, i.loggerInterface())
	}
	if !astfile.TypeExists(file, "dtxManager") {
		file.Decls = append(file.Decls, i.dtxManagerInterface())
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
					Text:  "//go:generate mockgen -source=interfaces.go -package=usecases -destination=interfaces_mock.go",
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
					Value: i.domain.AppConfig.ProjectConfig.UUIDImportPath(),
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: i.domain.AppConfig.ProjectConfig.LogImportPath(),
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: i.domain.AppConfig.ProjectConfig.DTXImportPath(),
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
								X:   ast.NewIdent("dtx"),
								Sel: ast.NewIdent("TX"),
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
								X:   ast.NewIdent("dtx"),
								Sel: ast.NewIdent("TX"),
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
								X:   ast.NewIdent("dtx"),
								Sel: ast.NewIdent("TX"),
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

func (i InterfacesGenerator) appEventServiceInterface() *ast.GenDecl {
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: &ast.Ident{
					Name: i.domain.EventServiceInterfaceName(),
				},
				Type: &ast.InterfaceType{
					Methods: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{
										Name: "Created",
									},
								},
								Type: &ast.FuncType{
									Params: &ast.FieldList{
										List: []*ast.Field{
											{
												Type: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "context",
													},
													Sel: &ast.Ident{
														Name: "Context",
													},
												},
											},
											{
												Type: &ast.SelectorExpr{
													X:   ast.NewIdent("dtx"),
													Sel: ast.NewIdent("TX"),
												},
											},
											{
												Type: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "entities",
													},
													Sel: &ast.Ident{
														Name: i.domain.GetMainModel().Name,
													},
												},
											},
										},
									},
									Results: &ast.FieldList{
										List: []*ast.Field{
											{
												Type: &ast.Ident{
													Name: "error",
												},
											},
										},
									},
								},
							},
							{
								Names: []*ast.Ident{
									{
										Name: "Updated",
									},
								},
								Type: &ast.FuncType{
									Params: &ast.FieldList{
										List: []*ast.Field{
											{
												Type: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "context",
													},
													Sel: &ast.Ident{
														Name: "Context",
													},
												},
											},
											{
												Type: &ast.SelectorExpr{
													X:   ast.NewIdent("dtx"),
													Sel: ast.NewIdent("TX"),
												},
											},
											{
												Type: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "entities",
													},
													Sel: &ast.Ident{
														Name: i.domain.GetMainModel().Name,
													},
												},
											},
										},
									},
									Results: &ast.FieldList{
										List: []*ast.Field{
											{
												Type: &ast.Ident{
													Name: "error",
												},
											},
										},
									},
								},
							},
							{
								Names: []*ast.Ident{
									{
										Name: "Deleted",
									},
								},
								Type: &ast.FuncType{
									Params: &ast.FieldList{
										List: []*ast.Field{
											{
												Type: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "context",
													},
													Sel: &ast.Ident{
														Name: "Context",
													},
												},
											},
											{
												Type: &ast.SelectorExpr{
													X:   ast.NewIdent("dtx"),
													Sel: ast.NewIdent("TX"),
												},
											},
											{
												Type: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "uuid",
													},
													Sel: &ast.Ident{
														Name: "UUID",
													},
												},
											},
										},
									},
									Results: &ast.FieldList{
										List: []*ast.Field{
											{
												Type: &ast.Ident{
													Name: "error",
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

func (i InterfacesGenerator) dtxManagerInterface() *ast.GenDecl {
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: &ast.Ident{
					Name: "dtxManager",
				},
				Type: &ast.InterfaceType{
					Methods: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{
										Name: "NewTx",
									},
								},
								Type: &ast.FuncType{
									Params: &ast.FieldList{},
									Results: &ast.FieldList{
										List: []*ast.Field{
											{
												Type: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "dtx",
													},
													Sel: &ast.Ident{
														Name: "TX",
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
