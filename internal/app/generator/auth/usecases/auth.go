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
	"path/filepath"

	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
)

type UseCaseAuth struct {
	project *configs.Project
}

func NewUseCaseAuth(project *configs.Project) *UseCaseAuth {
	return &UseCaseAuth{project: project}
}

func (i UseCaseAuth) Sync() error {
	fileset := token.NewFileSet()
	filename := filepath.Join("internal", "app", "auth", "usecases", "auth.go")
	if err := os.MkdirAll(path.Dir(filename), 0777); err != nil {
		return err
	}
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = i.file()
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

func (i UseCaseAuth) file() *ast.File {
	return &ast.File{
		Name: &ast.Ident{
			Name: "usecases",
		},
		Decls: []ast.Decl{
			&ast.GenDecl{
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
							Value: fmt.Sprintf(`"%s/internal/app/auth/entities"`, i.project.Module),
						},
					},
					&ast.ImportSpec{
						Name: ast.NewIdent("userEntities"),
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/app/user/entities"`, i.project.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/pkg/clock"`, i.project.Module),
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{
							Name: "AuthUseCase",
						},
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											{
												Name: "authService",
											},
										},
										Type: ast.NewIdent("AuthService"),
									},
									{
										Names: []*ast.Ident{
											{
												Name: "clock",
											},
										},
										Type: ast.NewIdent("clock"),
									},
									{
										Names: []*ast.Ident{
											{
												Name: "logger",
											},
										},
										Type: ast.NewIdent("logger"),
									},
								},
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Name: &ast.Ident{
					Name: "NewAuthUseCase",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{
										Name: "authService",
									},
								},
								Type: ast.NewIdent("AuthService"),
							},
							{
								Names: []*ast.Ident{
									{
										Name: "clock",
									},
								},
								Type: ast.NewIdent("clock"),
							},
							{
								Names: []*ast.Ident{
									{
										Name: "logger",
									},
								},
								Type: ast.NewIdent("logger"),
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: ast.NewIdent("*AuthUseCase"),
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.UnaryExpr{
									Op: token.AND,
									X: &ast.CompositeLit{
										Type: &ast.Ident{
											Name: "AuthUseCase",
										},
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "authService",
												},
												Value: &ast.Ident{
													Name: "authService",
												},
											},
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "clock",
												},
												Value: &ast.Ident{
													Name: "clock",
												},
											},
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "logger",
												},
												Value: &ast.Ident{
													Name: "logger",
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
			&ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								{
									Name: "i",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "AuthUseCase",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "CreateToken",
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
										Name: "login",
									},
								},
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "entities",
										},
										Sel: &ast.Ident{
											Name: "Login",
										},
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
										X: &ast.Ident{
											Name: "entities",
										},
										Sel: &ast.Ident{
											Name: "TokenPair",
										},
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
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "pair",
								},
								&ast.Ident{
									Name: "err",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "i",
											},
											Sel: &ast.Ident{
												Name: "authService",
											},
										},
										Sel: &ast.Ident{
											Name: "CreateToken",
										},
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "ctx",
										},
										&ast.Ident{
											Name: "login",
										},
									},
								},
							},
						},
						&ast.IfStmt{
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
												Name: "nil",
											},
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
									Name: "pair",
								},
								&ast.Ident{
									Name: "nil",
								},
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								{
									Name: "i",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "AuthUseCase",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "ValidateToken",
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
										Name: "token",
									},
								},
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "entities",
									},
									Sel: &ast.Ident{
										Name: "Token",
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
											X: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "i",
												},
												Sel: &ast.Ident{
													Name: "authService",
												},
											},
											Sel: &ast.Ident{
												Name: "ValidateToken",
											},
										},
										Args: []ast.Expr{
											&ast.Ident{
												Name: "ctx",
											},
											&ast.Ident{
												Name: "token",
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
			},
			&ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								{
									Name: "i",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "AuthUseCase",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "RefreshToken",
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
										Name: "refresh",
									},
								},
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "entities",
									},
									Sel: &ast.Ident{
										Name: "Token",
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
										X: &ast.Ident{
											Name: "entities",
										},
										Sel: &ast.Ident{
											Name: "TokenPair",
										},
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
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "pair",
								},
								&ast.Ident{
									Name: "err",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "i",
											},
											Sel: &ast.Ident{
												Name: "authService",
											},
										},
										Sel: &ast.Ident{
											Name: "RefreshToken",
										},
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "ctx",
										},
										&ast.Ident{
											Name: "refresh",
										},
									},
								},
							},
						},
						&ast.IfStmt{
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
												Name: "nil",
											},
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
									Name: "pair",
								},
								&ast.Ident{
									Name: "nil",
								},
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								{
									Name: "i",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "AuthUseCase",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "Auth",
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
										Name: "access",
									},
								},
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "entities",
									},
									Sel: &ast.Ident{
										Name: "Token",
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
										X: ast.NewIdent("userEntities"),
										Sel: &ast.Ident{
											Name: "User",
										},
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
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "user",
								},
								&ast.Ident{
									Name: "err",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "i",
											},
											Sel: &ast.Ident{
												Name: "authService",
											},
										},
										Sel: &ast.Ident{
											Name: "Auth",
										},
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "ctx",
										},
										&ast.Ident{
											Name: "access",
										},
									},
								},
							},
						},
						&ast.IfStmt{
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
												Name: "nil",
											},
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
									Name: "user",
								},
								&ast.Ident{
									Name: "nil",
								},
							},
						},
					},
				},
			},
		},
	}
}
