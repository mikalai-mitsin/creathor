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

	"github.com/mikalai-mitsin/creathor/internal/pkg/domain"
)

type InterfacesGenerator struct {
	domain *domain.Domain
}

func NewInterfacesGenerator(domain *domain.Domain) *InterfacesGenerator {
	return &InterfacesGenerator{domain: domain}
}

func (r InterfacesGenerator) Sync() error {
	fileset := token.NewFileSet()
	filename := filepath.Join(
		"internal",
		"app",
		r.domain.DirName(),
		"repositories",
		"postgres",
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
	var loggerExists bool
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok {
			if t.Name.String() == "logger" {
				loggerExists = true
			}
			return true
		}
		return true
	})
	if !loggerExists {
		file.Decls = append(file.Decls, r.loggerInterface())
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
		Name: ast.NewIdent("postgres"),
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
					Text:  fmt.Sprintf("//go:generate mockgen -source=%s_interfaces.go -package=postgres -destination=%s_interfaces_mock.go", r.domain.SnakeName(), r.domain.SnakeName()),
				},
			},
		},
		Specs: []ast.Spec{
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf(`"%s/internal/pkg/log"`, r.domain.Module),
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
