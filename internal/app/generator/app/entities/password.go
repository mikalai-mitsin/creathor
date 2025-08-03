package entities

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"

	mods "github.com/mikalai-mitsin/creathor/internal/pkg/app"
)

type Password struct {
	domain *mods.BaseEntity
}

func NewPassword(domain *mods.BaseEntity) *Password {
	return &Password{domain: domain}
}

func (m *Password) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "app", m.domain.AppName(), "entities", m.domain.DirName(), m.domain.FileName())
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
						ast.NewIdent("m"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent("User"),
					},
				},
			},
		},
		Name: ast.NewIdent("SetPassword"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("password"),
						},
						Type: ast.NewIdent("string"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("fromPassword"),
						ast.NewIdent("_"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("bcrypt"),
								Sel: ast.NewIdent("GenerateFromPassword"),
							},
							Args: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.ArrayType{
										Elt: ast.NewIdent("byte"),
									},
									Args: []ast.Expr{
										ast.NewIdent("password"),
									},
								},
								&ast.SelectorExpr{
									X:   ast.NewIdent("bcrypt"),
									Sel: ast.NewIdent("DefaultCost"),
								},
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X:   ast.NewIdent("m"),
							Sel: ast.NewIdent("Password"),
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: ast.NewIdent("string"),
							Args: []ast.Expr{
								ast.NewIdent("fromPassword"),
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
						ast.NewIdent("m"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent("User"),
					},
				},
			},
		},
		Name: ast.NewIdent("CheckPassword"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("password"),
						},
						Type: ast.NewIdent("string"),
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
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.IfStmt{
					Init: &ast.AssignStmt{
						Lhs: []ast.Expr{
							ast.NewIdent("err"),
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("bcrypt"),
									Sel: ast.NewIdent("CompareHashAndPassword"),
								},
								Args: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.ArrayType{
											Elt: ast.NewIdent("byte"),
										},
										Args: []ast.Expr{
											&ast.SelectorExpr{
												X:   ast.NewIdent("m"),
												Sel: ast.NewIdent("Password"),
											},
										},
									},
									&ast.CallExpr{
										Fun: &ast.ArrayType{
											Elt: ast.NewIdent("byte"),
										},
										Args: []ast.Expr{
											ast.NewIdent("password"),
										},
									},
								},
							},
						},
					},
					Cond: &ast.BinaryExpr{
						X:  ast.NewIdent("err"),
						Op: token.NEQ,
						Y:  ast.NewIdent("nil"),
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   ast.NewIdent("errs"),
											Sel: ast.NewIdent("NewInvalidParameter"),
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
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
}
