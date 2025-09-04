package entities

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

type Ordering struct {
	entityConfig *configs.EntityConfig
}

func NewOrdering(entityConfig *configs.EntityConfig) *Ordering {
	return &Ordering{entityConfig: entityConfig}
}

func (r Ordering) Sync() error {
	fileset := token.NewFileSet()
	filename := filepath.Join(
		"internal",
		"app",
		r.entityConfig.AppConfig.AppName(),
		"entities",
		r.entityConfig.DirName(),
		r.entityConfig.FileName(),
	)
	err := os.MkdirAll(path.Dir(filename), 0777)
	if err != nil {
		return err
	}
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = r.file()
	}
	if !astfile.TypeExists(file, r.entityConfig.OrderingTypeName()) {
		file.Decls = append(file.Decls, r.typeDecl())
	}
	if _, ok := astfile.FindMethod(file, r.entityConfig.OrderingTypeName(), "Validate"); !ok {
		file.Decls = append(file.Decls, r.validateFunc())
	}
	if _, ok := astfile.FindMethod(file, r.entityConfig.OrderingTypeName(), "String"); !ok {
		file.Decls = append(file.Decls, r.stringerFunc())
	}
	for name, value := range r.entityConfig.OrderingConsts() {
		if !astfile.ConstExists(file, name) {
			file.Decls = append(file.Decls, &ast.GenDecl{
				Tok: token.CONST,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Doc: nil,
						Names: []*ast.Ident{
							{
								Name: name,
							},
						},
						Type: ast.NewIdent(r.entityConfig.OrderingTypeName()),
						Values: []ast.Expr{
							&ast.BasicLit{
								Value: value,
							},
						},
					},
				},
			})
		}
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

func (r Ordering) file() *ast.File {
	return &ast.File{
		Name: ast.NewIdent("entities"),
		Decls: []ast.Decl{
			r.imports(),
		},
	}
}

func (r Ordering) imports() *ast.GenDecl {
	return &ast.GenDecl{
		Tok: token.IMPORT,
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{
					Slash: token.NoPos,
					Text: fmt.Sprintf(
						"//go:generate mockgen -source=%s_interfaces.go -package=repositories -destination=%s_interfaces_mock.go",
						r.entityConfig.SnakeName(),
						r.entityConfig.SnakeName(),
					),
				},
			},
		},
		Specs: []ast.Spec{
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: r.entityConfig.AppConfig.ProjectConfig.LogImportPath(),
				},
			},
		},
	}
}

func (r Ordering) typeDecl() *ast.GenDecl {
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: &ast.Ident{
					Name: r.entityConfig.OrderingTypeName(),
				},
				Type: &ast.Ident{
					Name: "string",
				},
			},
		},
	}

}

func (r Ordering) validateFunc() *ast.FuncDecl {
	validateIn := make([]ast.Expr, 0, len(r.entityConfig.GetMainModel().Params))
	for k := range r.entityConfig.OrderingConsts() {
		validateIn = append(
			validateIn,
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent(k),
					Sel: ast.NewIdent("String"),
				},
			},
		)
	}
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{
							Name: "o",
						},
					},
					Type: &ast.Ident{
						Name: r.entityConfig.OrderingTypeName(),
					},
				},
			},
		},
		Name: &ast.Ident{
			Name: "Validate",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
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
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.IfStmt{
					Init: &ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.Ident{
								Name: "err",
							},
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "validation",
									},
									Sel: &ast.Ident{
										Name: "Validate",
									},
								},
								Args: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   ast.NewIdent("o"),
											Sel: ast.NewIdent("String"),
										},
									},
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "validation",
											},
											Sel: &ast.Ident{
												Name: "In",
											},
										},
										Args: validateIn,
									},
								},
							},
						},
					},
					Cond: &ast.BinaryExpr{
						X: &ast.Ident{
							Name: "err",
						},
						Op: token.NEQ,
						Y: &ast.Ident{
							Name: "nil",
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.Ident{
										Name: "err",
									},
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.Ident{
							Name: "nil",
						},
					},
				},
			},
		},
	}
}

func (r Ordering) stringerFunc() *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{
							Name: "o",
						},
					},
					Type: &ast.Ident{
						Name: r.entityConfig.OrderingTypeName(),
					},
				},
			},
		},
		Name: &ast.Ident{
			Name: "String",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.Ident{
							Name: "string",
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.Ident{
								Name: "string",
							},
							Args: []ast.Expr{
								&ast.Ident{
									Name: "o",
								},
							},
						},
					},
				},
			},
		},
	}
}
