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

func (a AppAuth) Sync() error {
	fileset := token.NewFileSet()
	filename := filepath.Join("internal", "app", "auth", "app.go")
	if err := os.MkdirAll(path.Dir(filename), 0777); err != nil {
		return err
	}
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = a.file()
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

func (a AppAuth) file() *ast.File {
	importSpecs := []ast.Spec{
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s/internal/app/auth/usecases"`, a.project.Module),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind: token.STRING,
				Value: fmt.Sprintf(
					`"%s/internal/app/auth/repositories/jwt"`,
					a.project.Module,
				),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s/internal/app/auth/services"`, a.project.Module),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s/internal/app/user/repositories/postgres"`, a.project.Module),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s/internal/pkg/clock"`, a.project.Module),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s/internal/pkg/configs"`, a.project.Module),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s/internal/pkg/log"`, a.project.Module),
			},
		},

		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: `"github.com/jmoiron/sqlx"`,
			},
		},
	}
	if a.project.HTTPEnabled {
		importSpecs = append(importSpecs, &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind: token.STRING,
				Value: fmt.Sprintf(
					`"%s/internal/pkg/http"`,
					a.project.Module,
				),
			},
		},
			&ast.ImportSpec{
				Name: ast.NewIdent("httpHandlers"),
				Path: &ast.BasicLit{
					Kind: token.STRING,
					Value: fmt.Sprintf(
						`"%s/internal/app/auth/handlers/http"`,
						a.project.Module,
					),
				},
			})
	}
	if a.project.GRPCEnabled {
		importSpecs = append(importSpecs, &ast.ImportSpec{
			Name: ast.NewIdent("grpcHandlers"),
			Path: &ast.BasicLit{
				Kind: token.STRING,
				Value: fmt.Sprintf(
					`"%s/internal/app/auth/handlers/grpc"`,
					a.project.Module,
				),
			},
		}, &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s/internal/pkg/grpc"`, a.project.Module),
			},
		}, &ast.ImportSpec{
			Name: ast.NewIdent(a.project.ProtoPackage()),
			Path: &ast.BasicLit{
				Kind: token.STRING,
				Value: fmt.Sprintf(
					`"%s/pkg/%s/v1"`,
					a.project.Module,
					a.project.ProtoPackage(),
				),
			},
		})
	}
	decls := []ast.Decl{
		&ast.GenDecl{
			Tok:   token.IMPORT,
			Specs: importSpecs,
		},
		&ast.GenDecl{
			Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{
					Name: ast.NewIdent("App"),
					Type: &ast.StructType{
						Fields: &ast.FieldList{
							List: []*ast.Field{
								{
									Names: []*ast.Ident{
										ast.NewIdent("db"),
									},
									Type: &ast.StarExpr{
										X: &ast.SelectorExpr{
											X:   ast.NewIdent("sqlx"),
											Sel: ast.NewIdent("DB"),
										},
									},
								},
								{
									Names: []*ast.Ident{
										ast.NewIdent("grpcServer"),
									},
									Type: &ast.StarExpr{
										X: &ast.SelectorExpr{
											X:   ast.NewIdent("grpc"),
											Sel: ast.NewIdent("Server"),
										},
									},
								},
								{
									Names: []*ast.Ident{
										ast.NewIdent("logger"),
									},
									Type: &ast.StarExpr{
										X: &ast.SelectorExpr{
											X:   ast.NewIdent("log"),
											Sel: ast.NewIdent("Log"),
										},
									},
								},
								{
									Names: []*ast.Ident{
										ast.NewIdent("authRepository"),
									},
									Type: &ast.StarExpr{
										X: &ast.SelectorExpr{
											X:   ast.NewIdent("jwt"),
											Sel: ast.NewIdent("AuthRepository"),
										},
									},
								},
								{
									Names: []*ast.Ident{
										ast.NewIdent("authService"),
									},
									Type: &ast.StarExpr{
										X: &ast.SelectorExpr{
											X:   ast.NewIdent("services"),
											Sel: ast.NewIdent("AuthService"),
										},
									},
								},
								{
									Names: []*ast.Ident{
										ast.NewIdent("authUseCase"),
									},
									Type: &ast.StarExpr{
										X: &ast.SelectorExpr{
											X:   ast.NewIdent("usecases"),
											Sel: ast.NewIdent("AuthUseCase"),
										},
									},
								},
								{
									Names: []*ast.Ident{
										ast.NewIdent("grpcAuthHandler"),
									},
									Type: &ast.StarExpr{
										X: &ast.SelectorExpr{
											X:   ast.NewIdent("grpcHandlers"),
											Sel: ast.NewIdent("AuthServiceServer"),
										},
									},
								},
								{
									Names: []*ast.Ident{
										ast.NewIdent("httpAuthHandler"),
									},
									Type: &ast.StarExpr{
										X: &ast.SelectorExpr{
											X:   ast.NewIdent("httpHandlers"),
											Sel: ast.NewIdent("AuthHandler"),
										},
									},
								},
								{
									Names: []*ast.Ident{
										ast.NewIdent("grpcAuthMiddleware"),
									},
									Type: &ast.StarExpr{
										X: &ast.SelectorExpr{
											X:   ast.NewIdent("grpcHandlers"),
											Sel: ast.NewIdent("AuthMiddleware"),
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
			Name: ast.NewIdent("NewApp"),
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								ast.NewIdent("db"),
							},
							Type: &ast.StarExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("sqlx"),
									Sel: ast.NewIdent("DB"),
								},
							},
						},
						{
							Names: []*ast.Ident{
								ast.NewIdent("config"),
							},
							Type: &ast.StarExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("configs"),
									Sel: ast.NewIdent("Config"),
								},
							},
						},
						{
							Names: []*ast.Ident{
								ast.NewIdent("logger"),
							},
							Type: &ast.StarExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("log"),
									Sel: ast.NewIdent("Log"),
								},
							},
						},
						{
							Names: []*ast.Ident{
								ast.NewIdent("clock"),
							},
							Type: &ast.StarExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("clock"),
									Sel: ast.NewIdent("Clock"),
								},
							},
						},
					},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: &ast.StarExpr{
								X: ast.NewIdent("App"),
							},
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							ast.NewIdent("userRepository"),
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("postgres"),
									Sel: ast.NewIdent("NewUserRepository"),
								},
								Args: []ast.Expr{
									ast.NewIdent("db"),
									ast.NewIdent("logger"),
								},
							},
						},
					},
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							ast.NewIdent("authRepository"),
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("jwt"),
									Sel: ast.NewIdent("NewAuthRepository"),
								},
								Args: []ast.Expr{
									ast.NewIdent("config"),
									ast.NewIdent("clock"),
									ast.NewIdent("logger"),
								},
							},
						},
					},
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							ast.NewIdent("authService"),
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("services"),
									Sel: ast.NewIdent("NewAuthService"),
								},
								Args: []ast.Expr{
									ast.NewIdent("authRepository"),
									ast.NewIdent("userRepository"),
									ast.NewIdent("logger"),
								},
							},
						},
					},
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							ast.NewIdent("authUseCase"),
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("usecases"),
									Sel: ast.NewIdent("NewAuthUseCase"),
								},
								Args: []ast.Expr{
									ast.NewIdent("authService"),
									ast.NewIdent("clock"),
									ast.NewIdent("logger"),
								},
							},
						},
					},
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							ast.NewIdent("grpcAuthHandler"),
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("grpcHandlers"),
									Sel: ast.NewIdent("NewAuthServiceServer"),
								},
								Args: []ast.Expr{
									ast.NewIdent("authUseCase"),
									ast.NewIdent("logger"),
								},
							},
						},
					},
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							ast.NewIdent("httpAuthHandler"),
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("httpHandlers"),
									Sel: ast.NewIdent("NewAuthHandler"),
								},
								Args: []ast.Expr{
									ast.NewIdent("authUseCase"),
									ast.NewIdent("logger"),
								},
							},
						},
					},
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							ast.NewIdent("grpcAuthMiddleware"),
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("grpcHandlers"),
									Sel: ast.NewIdent("NewAuthMiddleware"),
								},
								Args: []ast.Expr{
									ast.NewIdent("authService"),
									ast.NewIdent("logger"),
									ast.NewIdent("config"),
								},
							},
						},
					},
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.UnaryExpr{
								Op: token.AND,
								X: &ast.CompositeLit{
									Type: ast.NewIdent("App"),
									Elts: []ast.Expr{
										&ast.KeyValueExpr{
											Key:   ast.NewIdent("db"),
											Value: ast.NewIdent("db"),
										},
										&ast.KeyValueExpr{
											Key:   ast.NewIdent("logger"),
											Value: ast.NewIdent("logger"),
										},
										&ast.KeyValueExpr{
											Key:   ast.NewIdent("authRepository"),
											Value: ast.NewIdent("authRepository"),
										},
										&ast.KeyValueExpr{
											Key:   ast.NewIdent("authService"),
											Value: ast.NewIdent("authService"),
										},
										&ast.KeyValueExpr{
											Key:   ast.NewIdent("authUseCase"),
											Value: ast.NewIdent("authUseCase"),
										},
										&ast.KeyValueExpr{
											Key:   ast.NewIdent("grpcAuthHandler"),
											Value: ast.NewIdent("grpcAuthHandler"),
										},
										&ast.KeyValueExpr{
											Key:   ast.NewIdent("grpcAuthMiddleware"),
											Value: ast.NewIdent("grpcAuthMiddleware"),
										},
										&ast.KeyValueExpr{
											Key:   ast.NewIdent("httpAuthHandler"),
											Value: ast.NewIdent("httpAuthHandler"),
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
	if a.project.HTTPEnabled {
		decls = append(decls, a.registerHTTP())
	}
	if a.project.GRPCEnabled {
		decls = append(decls, a.registerGRPC())
	}
	return &ast.File{
		Name:  ast.NewIdent("auth"),
		Decls: decls,
	}
}

func (a AppAuth) registerGRPC() *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("a"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent("App"),
					},
				},
			},
		},
		Name: ast.NewIdent("RegisterGRPC"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("grpcServer"),
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("grpc"),
								Sel: ast.NewIdent("Server"),
							},
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
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("grpcServer"),
							Sel: ast.NewIdent("AddHandler"),
						},
						Args: []ast.Expr{
							&ast.UnaryExpr{
								Op: token.AND,
								X: &ast.SelectorExpr{
									X:   ast.NewIdent(a.project.ProtoPackage()),
									Sel: ast.NewIdent("AuthService_ServiceDesc"),
								},
							},
							&ast.SelectorExpr{
								X:   ast.NewIdent("a"),
								Sel: ast.NewIdent("grpcAuthHandler"),
							},
						},
					},
				},
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("grpcServer"),
							Sel: ast.NewIdent("AddInterceptor"),
						},
						Args: []ast.Expr{
							&ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("a"),
									Sel: ast.NewIdent("grpcAuthMiddleware"),
								},
								Sel: ast.NewIdent("UnaryServerInterceptor"),
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

func (a AppAuth) registerHTTP() *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("a"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent("App"),
					},
				},
			},
		},
		Name: ast.NewIdent("RegisterHTTP"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("httpServer"),
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("http"),
								Sel: ast.NewIdent("Server"),
							},
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
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("httpServer"),
							Sel: ast.NewIdent("Mount"),
						},
						Args: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: `"/api/v1/auth/"`,
							},
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.SelectorExpr{
										X:   ast.NewIdent("a"),
										Sel: ast.NewIdent("httpAuthHandler"),
									},
									Sel: ast.NewIdent("ChiRouter"),
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
