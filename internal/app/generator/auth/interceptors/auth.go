package interceptors

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

type InterceptorAuth struct {
	project *configs.Project
}

func NewInterceptorAuth(project *configs.Project) *InterceptorAuth {
	return &InterceptorAuth{project: project}
}

func (i InterceptorAuth) Sync() error {
	fileset := token.NewFileSet()
	filename := filepath.Join("internal", "app", "auth", "interceptors", "auth.go")
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

func (i InterceptorAuth) file() *ast.File {
	return &ast.File{
		Name: &ast.Ident{
			Name: "interceptors",
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
							Value: fmt.Sprintf(`"%s/internal/app/auth/models"`, i.project.Module),
						},
					},
					&ast.ImportSpec{
						Name: ast.NewIdent("userModels"),
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/app/user/models"`, i.project.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/pkg/clock"`, i.project.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/pkg/log"`, i.project.Module),
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{
							Name: "AuthInterceptor",
						},
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											{
												Name: "authUseCase",
											},
										},
										Type: ast.NewIdent("AuthUseCase"),
									},
									{
										Names: []*ast.Ident{
											{
												Name: "clock",
											},
										},
										Type: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "clock",
											},
											Sel: &ast.Ident{
												Name: "Clock",
											},
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "logger",
											},
										},
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
			},
			&ast.FuncDecl{
				Name: &ast.Ident{
					Name: "NewAuthInterceptor",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{
										Name: "authUseCase",
									},
								},
								Type: ast.NewIdent("AuthUseCase"),
							},
							{
								Names: []*ast.Ident{
									{
										Name: "clock",
									},
								},
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "clock",
									},
									Sel: &ast.Ident{
										Name: "Clock",
									},
								},
							},
							{
								Names: []*ast.Ident{
									{
										Name: "logger",
									},
								},
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
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: ast.NewIdent("*AuthInterceptor"),
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
											Name: "AuthInterceptor",
										},
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "authUseCase",
												},
												Value: &ast.Ident{
													Name: "authUseCase",
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
									Name: "AuthInterceptor",
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
											Name: "models",
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
											Name: "models",
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
												Name: "authUseCase",
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
									Name: "AuthInterceptor",
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
										Name: "models",
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
													Name: "authUseCase",
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
									Name: "AuthInterceptor",
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
										Name: "models",
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
											Name: "models",
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
												Name: "authUseCase",
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
									Name: "AuthInterceptor",
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
										Name: "models",
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
										X: ast.NewIdent("userModels"),
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
												Name: "authUseCase",
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
