package models

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"

	mods "github.com/mikalai-mitsin/creathor/internal/pkg/domain"
)

type Password struct {
	domain *mods.Domain
}

func NewPassword(domain *mods.Domain) *Password {
	return &Password{domain: domain}
}

func (m *Password) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "app", m.domain.DirName(), "models", m.domain.FileName())
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var setPassword *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if fun, ok := node.(*ast.FuncDecl); ok && fun.Name.Name == "SetPassword" {
			setPassword = fun
			return false
		}
		return true
	})
	if setPassword == nil {
		setPassword = m.setPasswordMethod()
		file.Decls = append(file.Decls, setPassword)
	}
	var checkPassword *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if fun, ok := node.(*ast.FuncDecl); ok && fun.Name.Name == "CheckPassword" {
			checkPassword = fun
			return false
		}
		return true
	})
	if checkPassword == nil {
		checkPassword = m.checkPasswordMethod()
		file.Decls = append(file.Decls, checkPassword)
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

func (m *Password) setPasswordMethod() *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{
							Name: "m",
						},
					},
					Type: &ast.StarExpr{
						X: &ast.Ident{
							Name: "User",
						},
					},
				},
			},
		},
		Name: &ast.Ident{
			Name: "SetPassword",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
								Name: "password",
							},
						},
						Type: &ast.Ident{
							Name: "string",
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "fromPassword",
						},
						&ast.Ident{
							Name: "_",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "bcrypt",
								},
								Sel: &ast.Ident{
									Name: "GenerateFromPassword",
								},
							},
							Args: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.ArrayType{
										Elt: &ast.Ident{
											Name: "byte",
										},
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "password",
										},
									},
								},
								&ast.SelectorExpr{
									X: &ast.Ident{
										Name: "bcrypt",
									},
									Sel: &ast.Ident{
										Name: "DefaultCost",
									},
								},
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X: &ast.Ident{
								Name: "m",
							},
							Sel: &ast.Ident{
								Name: "Password",
							},
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.Ident{
								Name: "string",
							},
							Args: []ast.Expr{
								&ast.Ident{
									Name: "fromPassword",
								},
							},
						},
					},
				},
			},
		},
	}
}

func (m *Password) checkPasswordMethod() *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{
							Name: "m",
						},
					},
					Type: &ast.StarExpr{
						X: &ast.Ident{
							Name: "User",
						},
					},
				},
			},
		},
		Name: &ast.Ident{
			Name: "CheckPassword",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
								Name: "password",
							},
						},
						Type: &ast.Ident{
							Name: "string",
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
										Name: "bcrypt",
									},
									Sel: &ast.Ident{
										Name: "CompareHashAndPassword",
									},
								},
								Args: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.ArrayType{
											Elt: &ast.Ident{
												Name: "byte",
											},
										},
										Args: []ast.Expr{
											&ast.SelectorExpr{
												X: &ast.Ident{
													Name: "m",
												},
												Sel: &ast.Ident{
													Name: "Password",
												},
											},
										},
									},
									&ast.CallExpr{
										Fun: &ast.ArrayType{
											Elt: &ast.Ident{
												Name: "byte",
											},
										},
										Args: []ast.Expr{
											&ast.Ident{
												Name: "password",
											},
										},
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
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "errs",
											},
											Sel: &ast.Ident{
												Name: "NewInvalidParameter",
											},
										},
										Args: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: "\"email or password\"",
											},
										},
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
