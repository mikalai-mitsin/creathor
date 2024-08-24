package grpc

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"

	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
)

type Middlewares struct {
	project *configs.Project
}

func NewMiddlewares(project *configs.Project) *Middlewares {
	return &Middlewares{project: project}
}

func (m *Middlewares) file() *ast.File {
	return &ast.File{
		Name: &ast.Ident{
			Name: "handlers",
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
							Value: `"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/app/auth/models"`, m.project.Module),
						},
					},
					&ast.ImportSpec{
						Name: &ast.Ident{
							Name: "userModels",
						},
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/app/user/models"`, m.project.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/pkg/auth"`, m.project.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/pkg/configs"`, m.project.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/pkg/log"`, m.project.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"google.golang.org/grpc"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"google.golang.org/grpc/codes"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"google.golang.org/grpc/status"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"strings"`,
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.CONST,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{
							{
								Name: "headerAuthorize",
							},
						},
						Values: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: `"authorization"`,
							},
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							{
								Name: "expectedScheme",
							},
						},
						Values: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: `"bearer"`,
							},
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
						Type: &ast.InterfaceType{
							Methods: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											{
												Name: "Auth",
											},
										},
										Type: &ast.FuncType{
											Params: &ast.FieldList{
												List: []*ast.Field{
													{
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
																	Name: "userModels",
																},
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
									},
								},
							},
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{
							Name: "AuthMiddleware",
						},
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
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
												Name: "authUseCase",
											},
										},
										Type: &ast.Ident{
											Name: "AuthUseCase",
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
					Name: "NewAuthMiddleware",
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
								Type: &ast.Ident{
									Name: "AuthUseCase",
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
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.StarExpr{
									X: &ast.Ident{
										Name: "AuthMiddleware",
									},
								},
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
											Name: "AuthMiddleware",
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
													Name: "logger",
												},
												Value: &ast.Ident{
													Name: "logger",
												},
											},
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "config",
												},
												Value: &ast.Ident{
													Name: "config",
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
									Name: "m",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "AuthMiddleware",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "UnaryServerInterceptor",
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
										Name: "req",
									},
								},
								Type: &ast.Ident{
									Name: "any",
								},
							},
							{
								Names: []*ast.Ident{
									{
										Name: "_",
									},
								},
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "grpc",
										},
										Sel: &ast.Ident{
											Name: "UnaryServerInfo",
										},
									},
								},
							},
							{
								Names: []*ast.Ident{
									{
										Name: "handler",
									},
								},
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "grpc",
									},
									Sel: &ast.Ident{
										Name: "UnaryHandler",
									},
								},
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.Ident{
									Name: "any",
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
									Name: "newCtx",
								},
								&ast.Ident{
									Name: "err",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "m",
										},
										Sel: &ast.Ident{
											Name: "auth",
										},
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "ctx",
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
								&ast.CallExpr{
									Fun: &ast.Ident{
										Name: "handler",
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "newCtx",
										},
										&ast.Ident{
											Name: "req",
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
									Name: "m",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "AuthMiddleware",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "auth",
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
								Type: &ast.Ident{
									Name: "error",
								},
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.DeclStmt{
							Decl: &ast.GenDecl{
								Tok: token.VAR,
								Specs: []ast.Spec{
									&ast.ValueSpec{
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
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "token",
								},
								&ast.Ident{
									Name: "err",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "m",
										},
										Sel: &ast.Ident{
											Name: "authFromMD",
										},
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "ctx",
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
												Name: "ctx",
											},
											&ast.Ident{
												Name: "err",
											},
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.BinaryExpr{
								X: &ast.Ident{
									Name: "token",
								},
								Op: token.EQL,
								Y: &ast.BasicLit{
									Kind:  token.STRING,
									Value: `""`,
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "auth",
													},
													Sel: &ast.Ident{
														Name: "PutUser",
													},
												},
												Args: []ast.Expr{
													&ast.Ident{
														Name: "ctx",
													},
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "models",
														},
														Sel: &ast.Ident{
															Name: "Guest",
														},
													},
												},
											},
											&ast.Ident{
												Name: "nil",
											},
										},
									},
								},
							},
						},
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
												Name: "m",
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
											Name: "token",
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
												Name: "ctx",
											},
											&ast.Ident{
												Name: "err",
											},
										},
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "newCtx",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "auth",
										},
										Sel: &ast.Ident{
											Name: "PutUser",
										},
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "ctx",
										},
										&ast.Ident{
											Name: "user",
										},
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.Ident{
									Name: "newCtx",
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
									Name: "m",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "AuthMiddleware",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "authFromMD",
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
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "models",
									},
									Sel: &ast.Ident{
										Name: "Token",
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
									Name: "val",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "metautils",
												},
												Sel: &ast.Ident{
													Name: "ExtractIncoming",
												},
											},
											Args: []ast.Expr{
												&ast.Ident{
													Name: "ctx",
												},
											},
										},
										Sel: &ast.Ident{
											Name: "Get",
										},
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "headerAuthorize",
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.BinaryExpr{
								X: &ast.Ident{
									Name: "val",
								},
								Op: token.EQL,
								Y: &ast.BasicLit{
									Kind:  token.STRING,
									Value: `""`,
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: `""`,
											},
											&ast.Ident{
												Name: "nil",
											},
										},
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "splits",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "strings",
										},
										Sel: &ast.Ident{
											Name: "SplitN",
										},
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "val",
										},
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: `" "`,
										},
										&ast.BasicLit{
											Kind:  token.INT,
											Value: "2",
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.BinaryExpr{
								X: &ast.CallExpr{
									Fun: &ast.Ident{
										Name: "len",
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "splits",
										},
									},
								},
								Op: token.LSS,
								Y: &ast.BasicLit{
									Kind:  token.INT,
									Value: "2",
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: `""`,
											},
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "status",
													},
													Sel: &ast.Ident{
														Name: "Errorf",
													},
												},
												Args: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "codes",
														},
														Sel: &ast.Ident{
															Name: "Unauthenticated",
														},
													},
													&ast.BasicLit{
														Kind:  token.STRING,
														Value: `"Bad authorization string"`,
													},
												},
											},
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.UnaryExpr{
								Op: token.NOT,
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "strings",
										},
										Sel: &ast.Ident{
											Name: "EqualFold",
										},
									},
									Args: []ast.Expr{
										&ast.IndexExpr{
											X: &ast.Ident{
												Name: "splits",
											},
											Index: &ast.BasicLit{
												Kind:  token.INT,
												Value: "0",
											},
										},
										&ast.Ident{
											Name: "expectedScheme",
										},
									},
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: `""`,
											},
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "status",
													},
													Sel: &ast.Ident{
														Name: "Errorf",
													},
												},
												Args: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "codes",
														},
														Sel: &ast.Ident{
															Name: "Unauthenticated",
														},
													},
													&ast.BinaryExpr{
														X: &ast.BasicLit{
															Kind:  token.STRING,
															Value: `"Request unauthenticated with "`,
														},
														Op: token.ADD,
														Y: &ast.Ident{
															Name: "expectedScheme",
														},
													},
												},
											},
										},
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "bearerToken",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "strings",
										},
										Sel: &ast.Ident{
											Name: "TrimSpace",
										},
									},
									Args: []ast.Expr{
										&ast.IndexExpr{
											X: &ast.Ident{
												Name: "splits",
											},
											Index: &ast.BasicLit{
												Kind:  token.INT,
												Value: "1",
											},
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.BinaryExpr{
								X: &ast.Ident{
									Name: "bearerToken",
								},
								Op: token.EQL,
								Y: &ast.BasicLit{
									Kind:  token.STRING,
									Value: `""`,
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: `""`,
											},
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "status",
													},
													Sel: &ast.Ident{
														Name: "Errorf",
													},
												},
												Args: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "codes",
														},
														Sel: &ast.Ident{
															Name: "Unauthenticated",
														},
													},
													&ast.BinaryExpr{
														X: &ast.BasicLit{
															Kind:  token.STRING,
															Value: `"Request unauthenticated with "`,
														},
														Op: token.ADD,
														Y: &ast.Ident{
															Name: "expectedScheme",
														},
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
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "models",
										},
										Sel: &ast.Ident{
											Name: "Token",
										},
									},
									Args: []ast.Expr{
										&ast.IndexExpr{
											X: &ast.Ident{
												Name: "splits",
											},
											Index: &ast.BasicLit{
												Kind:  token.INT,
												Value: "1",
											},
										},
									},
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

func (m *Middlewares) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "app", "auth", "handlers", "grpc", "middlewares.go")
	if err := os.MkdirAll(path.Dir(filename), 0777); err != nil {
		return err
	}
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = m.file()
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
