package services

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
	filename := filepath.Join(
		"internal",
		"app",
		i.domain.AppConfig.AppName(),
		"services",
		i.domain.DirName(),
		fmt.Sprintf("%s_interfaces.go", i.domain.SnakeName()),
	)
	err := os.MkdirAll(path.Dir(filename), 0777)
	if err != nil {
		return err
	}
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = i.file()
	}
	repositoryExists := false
	loggerExists := false
	clockExists := false
	uuidGeneratorExists := false
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok {
			if t.Name.String() == i.domain.GetRepositoryInterfaceName() {
				repositoryExists = true
			}
			if t.Name.String() == "logger" {
				loggerExists = true
			}
			if t.Name.String() == "clock" {
				clockExists = true
			}
			if t.Name.String() == "uuidGenerator" {
				uuidGeneratorExists = true
			}
			return true
		}
		return true
	})
	if !repositoryExists {
		file.Decls = append(file.Decls, i.repositoryInterface())
	}
	if !clockExists {
		file.Decls = append(file.Decls, i.clockInterface())
	}
	if !loggerExists {
		file.Decls = append(file.Decls, i.loggerInterface())
	}
	if !uuidGeneratorExists {
		file.Decls = append(file.Decls, i.uuidGeneratorInterface())
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
		Name: ast.NewIdent("services"),
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
					Text: fmt.Sprintf(
						"//go:generate mockgen -source=%s_interfaces.go -package=services -destination=%s_interfaces_mock.go",
						i.domain.SnakeName(),
						i.domain.SnakeName(),
					),
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
					Value: `"time"`,
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
					Value: i.domain.AppConfig.ProjectConfig.LogImportPath(),
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
					Value: i.domain.AppConfig.ProjectConfig.DTXImportPath(),
				},
			},
		},
	}
}

func (i InterfacesGenerator) repositoryInterface() *ast.GenDecl {
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
								Sel: ast.NewIdent(i.domain.GetMainModel().Name),
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
							Type: ast.NewIdent("error"),
						},
					},
				},
			},
		},
		{
			Names: []*ast.Ident{ast.NewIdent("Count")},
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
								Sel: ast.NewIdent(i.domain.GetMainModel().Name),
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
				Name: ast.NewIdent(i.domain.GetRepositoryInterfaceName()),
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

func (i InterfacesGenerator) clockInterface() *ast.GenDecl {
	return &ast.GenDecl{
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{
					Text: "// clock - clock interface",
				},
			},
		},
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent("clock"),
				Type: &ast.InterfaceType{
					Methods: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("Now"),
								},
								Type: &ast.FuncType{
									Results: &ast.FieldList{
										List: []*ast.Field{
											{
												Type: ast.NewIdent("time.Time"),
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

func (i InterfacesGenerator) uuidGeneratorInterface() *ast.GenDecl {
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent("uuidGenerator"),
				Type: &ast.InterfaceType{
					Methods: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("NewUUID"),
								},
								Type: &ast.FuncType{
									Results: &ast.FieldList{
										List: []*ast.Field{
											{
												Type: ast.NewIdent("uuid.UUID"),
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
