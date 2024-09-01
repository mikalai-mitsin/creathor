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

type Server struct {
	project *configs.Project
}

func NewServer(project *configs.Project) *Server {
	return &Server{project: project}
}

func (u Server) file() *ast.File {
	return &ast.File{
		Name: &ast.Ident{
			Name: "grpc",
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
						Name: ast.NewIdent("grpc_zap"),
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"google.golang.org/grpc/health"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"google.golang.org/grpc/health/grpc_health_v1"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"google.golang.org/grpc/reflection"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"net"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/pkg/configs"`, u.project.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/pkg/log"`, u.project.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"google.golang.org/grpc"`,
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{
							Name: "Server",
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
												Name: "server",
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
												Name: "handlers",
											},
										},
										Type: &ast.MapType{
											Key: &ast.StarExpr{
												X: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "grpc",
													},
													Sel: &ast.Ident{
														Name: "ServiceDesc",
													},
												},
											},
											Value: &ast.Ident{
												Name: "any",
											},
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "unaryUseCases",
											},
										},
										Type: &ast.ArrayType{
											Elt: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "grpc",
												},
												Sel: &ast.Ident{
													Name: "UnaryServerUseCase",
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
					Name: "NewServer",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
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
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.StarExpr{
									X: &ast.Ident{
										Name: "Server",
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
											Name: "Server",
										},
										Elts: []ast.Expr{
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
													Name: "server",
												},
												Value: &ast.Ident{
													Name: "nil",
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
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "handlers",
												},
												Value: &ast.CompositeLit{
													Type: &ast.MapType{
														Key: &ast.StarExpr{
															X: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "grpc",
																},
																Sel: &ast.Ident{
																	Name: "ServiceDesc",
																},
															},
														},
														Value: &ast.Ident{
															Name: "any",
														},
													},
												},
											},
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "unaryUseCases",
												},
												Value: &ast.CompositeLit{
													Type: &ast.ArrayType{
														Elt: &ast.SelectorExpr{
															X: &ast.Ident{
																Name: "grpc",
															},
															Sel: &ast.Ident{
																Name: "UnaryServerUseCase",
															},
														},
													},
													Elts: []ast.Expr{
														&ast.Ident{
															Name: "unaryErrorServerUseCase",
														},
														&ast.CallExpr{
															Fun: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "otelgrpc",
																},
																Sel: &ast.Ident{
																	Name: "UnaryServerUseCase",
																},
															},
														},
														&ast.CallExpr{
															Fun: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "grpc_zap",
																},
																Sel: &ast.Ident{
																	Name: "UnaryServerUseCase",
																},
															},
															Args: []ast.Expr{
																&ast.CallExpr{
																	Fun: &ast.SelectorExpr{
																		X: &ast.Ident{
																			Name: "logger",
																		},
																		Sel: &ast.Ident{
																			Name: "Logger",
																		},
																	},
																},
																&ast.CallExpr{
																	Fun: &ast.SelectorExpr{
																		X: &ast.Ident{
																			Name: "grpc_zap",
																		},
																		Sel: &ast.Ident{
																			Name: "WithMessageProducer",
																		},
																	},
																	Args: []ast.Expr{
																		&ast.Ident{
																			Name: "defaultMessageProducer",
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
									Name: "s",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "Server",
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
										Name: "_",
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
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.SelectorExpr{
									X: &ast.Ident{
										Name: "s",
									},
									Sel: &ast.Ident{
										Name: "server",
									},
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "grpc",
										},
										Sel: &ast.Ident{
											Name: "NewServer",
										},
									},
									Args: []ast.Expr{
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "grpc",
												},
												Sel: &ast.Ident{
													Name: "ChainUnaryUseCase",
												},
											},
											Args: []ast.Expr{
												&ast.SelectorExpr{
													X: &ast.Ident{
														Name: "s",
													},
													Sel: &ast.Ident{
														Name: "unaryUseCases",
													},
												},
											},
											Ellipsis: 1180,
										},
									},
								},
							},
						},
						&ast.RangeStmt{
							Key: &ast.Ident{
								Name: "sd",
							},
							Value: &ast.Ident{
								Name: "ss",
							},
							Tok: token.DEFINE,
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "s",
								},
								Sel: &ast.Ident{
									Name: "handlers",
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "s",
													},
													Sel: &ast.Ident{
														Name: "server",
													},
												},
												Sel: &ast.Ident{
													Name: "RegisterService",
												},
											},
											Args: []ast.Expr{
												&ast.Ident{
													Name: "sd",
												},
												&ast.Ident{
													Name: "ss",
												},
											},
										},
									},
								},
							},
						},
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "reflection",
									},
									Sel: &ast.Ident{
										Name: "Register",
									},
								},
								Args: []ast.Expr{
									&ast.SelectorExpr{
										X: &ast.Ident{
											Name: "s",
										},
										Sel: &ast.Ident{
											Name: "server",
										},
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "healthServer",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "health",
										},
										Sel: &ast.Ident{
											Name: "NewServer",
										},
									},
								},
							},
						},
						&ast.RangeStmt{
							Key: &ast.Ident{
								Name: "service",
							},
							Tok: token.DEFINE,
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "s",
										},
										Sel: &ast.Ident{
											Name: "server",
										},
									},
									Sel: &ast.Ident{
										Name: "GetServiceInfo",
									},
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "healthServer",
												},
												Sel: &ast.Ident{
													Name: "SetServingStatus",
												},
											},
											Args: []ast.Expr{
												&ast.Ident{
													Name: "service",
												},
												&ast.SelectorExpr{
													X: &ast.Ident{
														Name: "grpc_health_v1",
													},
													Sel: &ast.Ident{
														Name: "HealthCheckResponse_SERVING",
													},
												},
											},
										},
									},
								},
							},
						},
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "grpc_health_v1",
									},
									Sel: &ast.Ident{
										Name: "RegisterHealthServer",
									},
								},
								Args: []ast.Expr{
									&ast.SelectorExpr{
										X: &ast.Ident{
											Name: "s",
										},
										Sel: &ast.Ident{
											Name: "server",
										},
									},
									&ast.Ident{
										Name: "healthServer",
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "listener",
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
											Name: "net",
										},
										Sel: &ast.Ident{
											Name: "Listen",
										},
									},
									Args: []ast.Expr{
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: `"tcp"`,
										},
										&ast.SelectorExpr{
											X: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "s",
												},
												Sel: &ast.Ident{
													Name: "config",
												},
											},
											Sel: &ast.Ident{
												Name: "BindAddr",
											},
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
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "s",
											},
											Sel: &ast.Ident{
												Name: "server",
											},
										},
										Sel: &ast.Ident{
											Name: "Serve",
										},
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "listener",
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
									Name: "s",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "Server",
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
										Name: "_",
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
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "s",
										},
										Sel: &ast.Ident{
											Name: "server",
										},
									},
									Sel: &ast.Ident{
										Name: "GracefulStop",
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
									Name: "s",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "Server",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "AddHandler",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{
										Name: "sd",
									},
								},
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "grpc",
										},
										Sel: &ast.Ident{
											Name: "ServiceDesc",
										},
									},
								},
							},
							{
								Names: []*ast.Ident{
									{
										Name: "ss",
									},
								},
								Type: &ast.Ident{
									Name: "any",
								},
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.IndexExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "s",
										},
										Sel: &ast.Ident{
											Name: "handlers",
										},
									},
									Index: &ast.Ident{
										Name: "sd",
									},
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.Ident{
									Name: "ss",
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
									Name: "s",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "Server",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "AddUseCase",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{
										Name: "usecase",
									},
								},
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "grpc",
									},
									Sel: &ast.Ident{
										Name: "UnaryServerUseCase",
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
								&ast.SelectorExpr{
									X: &ast.Ident{
										Name: "s",
									},
									Sel: &ast.Ident{
										Name: "unaryUseCases",
									},
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.Ident{
										Name: "append",
									},
									Args: []ast.Expr{
										&ast.SelectorExpr{
											X: &ast.Ident{
												Name: "s",
											},
											Sel: &ast.Ident{
												Name: "unaryUseCases",
											},
										},
										&ast.Ident{
											Name: "usecase",
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

func (u Server) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "pkg", "grpc", "server.go")
	if err := os.MkdirAll(path.Dir(filename), 0777); err != nil {
		return err
	}
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = u.file()
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
