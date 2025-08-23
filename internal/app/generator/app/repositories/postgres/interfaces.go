package postgres

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"
	"path/filepath"

	"github.com/mikalai-mitsin/creathor/internal/pkg/astfile"
	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
)

type InterfacesGenerator struct {
	domain *configs.BaseEntity
}

func NewInterfacesGenerator(domain *configs.BaseEntity) *InterfacesGenerator {
	return &InterfacesGenerator{domain: domain}
}

func (r InterfacesGenerator) Sync() error {
	fileset := token.NewFileSet()
	filename := filepath.Join(
		"internal",
		"app",
		r.domain.AppName(),
		"repositories",
		"postgres",
		r.domain.DirName(),
		fmt.Sprintf("%s_interfaces.go", r.domain.SnakeName()),
	)
	err := os.MkdirAll(path.Dir(filename), 0777)
	if err != nil {
		return err
	}
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = r.file()
	}
	if !astfile.TypeExists(file, "logger") {
		file.Decls = append(file.Decls, r.loggerInterface())
	}
	if !astfile.TypeExists(file, "database") {
		file.Decls = append(file.Decls, r.databaseInterface())
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

func (r InterfacesGenerator) file() *ast.File {
	return &ast.File{
		Name: ast.NewIdent("repositories"),
		Decls: []ast.Decl{
			r.imports(),
		},
	}
}

func (r InterfacesGenerator) imports() *ast.GenDecl {
	return &ast.GenDecl{
		Tok: token.IMPORT,
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{
					Slash: token.NoPos,
					Text:  fmt.Sprintf("//go:generate mockgen -source=%s_interfaces.go -package=repositories -destination=%s_interfaces_mock.go", r.domain.SnakeName(), r.domain.SnakeName()),
				},
			},
		},
		Specs: []ast.Spec{
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: r.domain.AppConfig.ProjectConfig.LogImportPath(),
				},
			},
		},
	}
}

func (r InterfacesGenerator) loggerInterface() *ast.GenDecl {
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

func (r InterfacesGenerator) databaseInterface() *ast.GenDecl {
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: &ast.Ident{
					Name: "database",
				},
				Type: &ast.InterfaceType{
					Methods: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{
										Name: "ExecContext",
									},
								},
								Type: &ast.FuncType{
									Params: &ast.FieldList{
										List: []*ast.Field{
											{
												Names: []*ast.Ident{
													{
														Name: "ctx",
													},
												},
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
												Names: []*ast.Ident{
													{
														Name: "query",
													},
												},
												Type: &ast.Ident{
													Name: "string",
												},
											},
											{
												Names: []*ast.Ident{
													{
														Name: "args",
													},
												},
												Type: &ast.Ellipsis{
													Ellipsis: 93,
													Elt: &ast.InterfaceType{
														Methods: &ast.FieldList{},
													},
												},
											},
										},
									},
									Results: &ast.FieldList{
										List: []*ast.Field{
											{
												Type: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "sql",
													},
													Sel: &ast.Ident{
														Name: "Result",
													},
												},
											},
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
										Name: "GetContext",
									},
								},
								Type: &ast.FuncType{
									Params: &ast.FieldList{
										List: []*ast.Field{
											{
												Names: []*ast.Ident{
													{
														Name: "ctx",
													},
												},
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
												Names: []*ast.Ident{
													{
														Name: "dest",
													},
												},
												Type: &ast.Ident{
													Name: "any",
												},
											},
											{
												Names: []*ast.Ident{
													{
														Name: "query",
													},
												},
												Type: &ast.Ident{
													Name: "string",
												},
											},
											{
												Names: []*ast.Ident{
													{
														Name: "args",
													},
												},
												Type: &ast.Ellipsis{
													Ellipsis: 191,
													Elt: &ast.InterfaceType{
														Methods: &ast.FieldList{},
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
										Name: "SelectContext",
									},
								},
								Type: &ast.FuncType{
									Params: &ast.FieldList{
										List: []*ast.Field{
											{
												Names: []*ast.Ident{
													{
														Name: "ctx",
													},
												},
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
												Names: []*ast.Ident{
													{
														Name: "dest",
													},
												},
												Type: &ast.Ident{
													Name: "any",
												},
											},
											{
												Names: []*ast.Ident{
													{
														Name: "query",
													},
												},
												Type: &ast.Ident{
													Name: "string",
												},
											},
											{
												Names: []*ast.Ident{
													{
														Name: "args",
													},
												},
												Type: &ast.Ellipsis{
													Ellipsis: 278,
													Elt: &ast.InterfaceType{
														Methods: &ast.FieldList{},
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
