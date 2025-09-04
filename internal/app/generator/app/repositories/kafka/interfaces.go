package kafka

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
	domain *configs.EntityConfig
}

func NewInterfacesGenerator(domain *configs.EntityConfig) *InterfacesGenerator {
	return &InterfacesGenerator{domain: domain}
}

func (r *InterfacesGenerator) filename() string {
	return filepath.Join(
		".",
		"internal",
		"app",
		r.domain.AppConfig.AppName(),
		"repositories",
		"kafka",
		r.domain.DirName(),
		fmt.Sprintf("%s_interfaces.go", r.domain.SnakeName()),
	)
}

func (r InterfacesGenerator) Sync() error {
	fileset := token.NewFileSet()
	filename := r.filename()
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
	if !astfile.TypeExists(file, "producer") {
		file.Decls = append(file.Decls, r.kafkaInterface())
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
		Name: ast.NewIdent("events"),
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
					Text: fmt.Sprintf(
						"//go:generate mockgen -source=%s_interfaces.go -package=events -destination=%s_interfaces_mock.go",
						r.domain.SnakeName(),
						r.domain.SnakeName(),
					),
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
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: r.domain.AppConfig.ProjectConfig.KafkaImportPath(),
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

func (r InterfacesGenerator) kafkaInterface() *ast.GenDecl {
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: &ast.Ident{
					Name: "producer",
				},
				Type: &ast.InterfaceType{
					Methods: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{
										Name: "Send",
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
														Name: "msg",
													},
												},
												Type: &ast.StarExpr{
													X: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "kafka",
														},
														Sel: &ast.Ident{
															Name: "Message",
														},
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
