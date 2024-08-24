package auth

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

type AppAuth struct {
	project *configs.Project
}

func NewAppAuth(project *configs.Project) *AppAuth {
	return &AppAuth{project: project}
}

func (i AppAuth) Sync() error {
	fileset := token.NewFileSet()
	filename := filepath.Join("internal", "app", "auth", "app.go")
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

func (i AppAuth) file() *ast.File {
	return &ast.File{
		Package: 1,
		Name: &ast.Ident{
			Name: "auth",
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
							Value: fmt.Sprintf(`"%s/internal/app/auth/handlers/grpc"`, i.project.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/app/auth/interceptors"`, i.project.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/app/auth/repositories/jwt"`, i.project.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/app/auth/usecases"`, i.project.Module),
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
							Value: fmt.Sprintf(`"%s/internal/pkg/configs"`, i.project.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/pkg/grpc"`, i.project.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/pkg/log"`, i.project.Module),
						},
					},
					&ast.ImportSpec{
						Name: &ast.Ident{
							Name: i.project.ProtoPackage(),
						},
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/pkg/%s/v1"`, i.project.Module, i.project.ProtoPackage()),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"github.com/jmoiron/sqlx"`,
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{
							Name: "App",
						},
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											{
												Name: "db",
											},
										},
										Type: &ast.StarExpr{
											X: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "sqlx",
												},
												Sel: &ast.Ident{
													Name: "DB",
												},
											},
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "grpcServer",
											},
										},
										Type: &ast.StarExpr{
											X: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "grpc",
												},
												Sel: &ast.Ident{
													Name: "Server",
												},
											},
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "logger",
											},
										},
										Type: &ast.StarExpr{
											X: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "log",
												},
												Sel: &ast.Ident{
													Name: "Log",
												},
											},
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "authRepository",
											},
										},
										Type: &ast.StarExpr{
											X: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "jwt",
												},
												Sel: &ast.Ident{
													Name: "AuthRepository",
												},
											},
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "authUseCase",
											},
										},
										Type: &ast.StarExpr{
											X: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "usecases",
												},
												Sel: &ast.Ident{
													Name: "AuthUseCase",
												},
											},
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "authInterceptor",
											},
										},
										Type: &ast.StarExpr{
											X: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "interceptors",
												},
												Sel: &ast.Ident{
													Name: "AuthInterceptor",
												},
											},
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "authHandler",
											},
										},
										Type: &ast.StarExpr{
											X: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "handlers",
												},
												Sel: &ast.Ident{
													Name: "AuthServiceServer",
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
				Name: &ast.Ident{
					Name: "NewApp",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{
										Name: "db",
									},
								},
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "sqlx",
										},
										Sel: &ast.Ident{
											Name: "DB",
										},
									},
								},
							},
							{
								Names: []*ast.Ident{
									{
										Name: "config",
									},
								},
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "configs",
										},
										Sel: &ast.Ident{
											Name: "Config",
										},
									},
								},
							},
							{
								Names: []*ast.Ident{
									{
										Name: "grpcServer",
									},
								},
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "grpc",
										},
										Sel: &ast.Ident{
											Name: "Server",
										},
									},
								},
							},
							{
								Names: []*ast.Ident{
									{
										Name: "logger",
									},
								},
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "log",
										},
										Sel: &ast.Ident{
											Name: "Log",
										},
									},
								},
							},
							{
								Names: []*ast.Ident{
									{
										Name: "clock",
									},
								},
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "clock",
										},
										Sel: &ast.Ident{
											Name: "Clock",
										},
									},
								},
							},
							{
								Names: []*ast.Ident{
									{
										Name: "userRepository",
									},
								},
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "usecases",
									},
									Sel: &ast.Ident{
										Name: "UserRepository",
									},
								},
							},
							{
								Names: []*ast.Ident{
									{
										Name: "permissionRepository",
									},
								},
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "usecases",
									},
									Sel: &ast.Ident{
										Name: "PermissionRepository",
									},
								},
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.StarExpr{
									X: &ast.Ident{
										Name: "App",
									},
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
									Name: "authRepository",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "jwt",
										},
										Sel: &ast.Ident{
											Name: "NewAuthRepository",
										},
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "config",
										},
										&ast.Ident{
											Name: "clock",
										},
										&ast.Ident{
											Name: "logger",
										},
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "authUseCase",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "usecases",
										},
										Sel: &ast.Ident{
											Name: "NewAuthUseCase",
										},
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "authRepository",
										},
										&ast.Ident{
											Name: "userRepository",
										},
										&ast.Ident{
											Name: "permissionRepository",
										},
										&ast.Ident{
											Name: "logger",
										},
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "authInterceptor",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "interceptors",
										},
										Sel: &ast.Ident{
											Name: "NewAuthInterceptor",
										},
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "authUseCase",
										},
										&ast.Ident{
											Name: "clock",
										},
										&ast.Ident{
											Name: "logger",
										},
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "authHandler",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "handlers",
										},
										Sel: &ast.Ident{
											Name: "NewAuthServiceServer",
										},
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "authInterceptor",
										},
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.UnaryExpr{
									Op: token.AND,
									X: &ast.CompositeLit{
										Type: &ast.Ident{
											Name: "App",
										},
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "db",
												},
												Value: &ast.Ident{
													Name: "db",
												},
											},
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "grpcServer",
												},
												Value: &ast.Ident{
													Name: "grpcServer",
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
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "authRepository",
												},
												Value: &ast.Ident{
													Name: "authRepository",
												},
											},
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
													Name: "authInterceptor",
												},
												Value: &ast.Ident{
													Name: "authInterceptor",
												},
											},
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "authHandler",
												},
												Value: &ast.Ident{
													Name: "authHandler",
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
									Name: "a",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "App",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "Start",
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
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "examplepb",
									},
									Sel: &ast.Ident{
										Name: "RegisterAuthServiceServer",
									},
								},
								Args: []ast.Expr{
									&ast.SelectorExpr{
										X: &ast.Ident{
											Name: "a",
										},
										Sel: &ast.Ident{
											Name: "grpcServer",
										},
									},
									&ast.SelectorExpr{
										X: &ast.Ident{
											Name: "a",
										},
										Sel: &ast.Ident{
											Name: "authHandler",
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
									Name: "a",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "App",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "Stop",
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
		},
	}
}
