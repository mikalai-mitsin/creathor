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
		Name: ast.NewIdent("grpc"),
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
						Name: ast.NewIdent("Server"),
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
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
											ast.NewIdent("server"),
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
											ast.NewIdent("config"),
										},
										Type: &ast.StarExpr{
											X: ast.NewIdent("Config"),
										},
									},
									{
										Names: []*ast.Ident{
											ast.NewIdent("handlers"),
										},
										Type: &ast.MapType{
											Key: &ast.StarExpr{
												X: &ast.SelectorExpr{
													X:   ast.NewIdent("grpc"),
													Sel: ast.NewIdent("ServiceDesc"),
												},
											},
											Value: ast.NewIdent("any"),
										},
									},
									{
										Names: []*ast.Ident{
											ast.NewIdent("unaryInterceptors"),
										},
										Type: &ast.ArrayType{
											Elt: &ast.SelectorExpr{
												X:   ast.NewIdent("grpc"),
												Sel: ast.NewIdent("UnaryServerInterceptor"),
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
				Name: ast.NewIdent("NewServer"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
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
									ast.NewIdent("config"),
								},
								Type: &ast.StarExpr{
									X: ast.NewIdent("Config"),
								},
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.StarExpr{
									X: ast.NewIdent("Server"),
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
										Type: ast.NewIdent("Server"),
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key:   ast.NewIdent("logger"),
												Value: ast.NewIdent("logger"),
											},
											&ast.KeyValueExpr{
												Key:   ast.NewIdent("server"),
												Value: ast.NewIdent("nil"),
											},
											&ast.KeyValueExpr{
												Key:   ast.NewIdent("config"),
												Value: ast.NewIdent("config"),
											},
											&ast.KeyValueExpr{
												Key: ast.NewIdent("handlers"),
												Value: &ast.CompositeLit{
													Type: &ast.MapType{
														Key: &ast.StarExpr{
															X: &ast.SelectorExpr{
																X:   ast.NewIdent("grpc"),
																Sel: ast.NewIdent("ServiceDesc"),
															},
														},
														Value: ast.NewIdent("any"),
													},
												},
											},
											&ast.KeyValueExpr{
												Key: ast.NewIdent("unaryInterceptors"),
												Value: &ast.CompositeLit{
													Type: &ast.ArrayType{
														Elt: &ast.SelectorExpr{
															X:   ast.NewIdent("grpc"),
															Sel: ast.NewIdent("UnaryServerInterceptor"),
														},
													},
													Elts: []ast.Expr{
														ast.NewIdent("unaryErrorServerInterceptor"),
														&ast.CallExpr{
															Fun: &ast.SelectorExpr{
																X:   ast.NewIdent("grpc_zap"),
																Sel: ast.NewIdent("UnaryServerInterceptor"),
															},
															Args: []ast.Expr{
																&ast.CallExpr{
																	Fun: &ast.SelectorExpr{
																		X:   ast.NewIdent("logger"),
																		Sel: ast.NewIdent("Logger"),
																	},
																},
																&ast.CallExpr{
																	Fun: &ast.SelectorExpr{
																		X:   ast.NewIdent("grpc_zap"),
																		Sel: ast.NewIdent("WithMessageProducer"),
																	},
																	Args: []ast.Expr{
																		ast.NewIdent("defaultMessageProducer"),
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
								ast.NewIdent("s"),
							},
							Type: &ast.StarExpr{
								X: ast.NewIdent("Server"),
							},
						},
					},
				},
				Name: ast.NewIdent("Start"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("_"),
								},
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("context"),
									Sel: ast.NewIdent("Context"),
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
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.SelectorExpr{
									X:   ast.NewIdent("s"),
									Sel: ast.NewIdent("server"),
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("grpc"),
										Sel: ast.NewIdent("NewServer"),
									},
									Args: []ast.Expr{
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("grpc"),
												Sel: ast.NewIdent("ChainUnaryInterceptor"),
											},
											Args: []ast.Expr{
												&ast.SelectorExpr{
													X:   ast.NewIdent("s"),
													Sel: ast.NewIdent("unaryInterceptors"),
												},
											},
											Ellipsis: 1180,
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("grpc"),
												Sel: ast.NewIdent("StatsHandler"),
											},
											Args: []ast.Expr{
												&ast.CallExpr{
													Fun: &ast.SelectorExpr{
														X:   ast.NewIdent("otelgrpc"),
														Sel: ast.NewIdent("NewServerHandler"),
													},
													Args: []ast.Expr{},
												},
											},
										},
									},
								},
							},
						},
						&ast.RangeStmt{
							Key:   ast.NewIdent("sd"),
							Value: ast.NewIdent("ss"),
							Tok:   token.DEFINE,
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("s"),
								Sel: ast.NewIdent("handlers"),
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.SelectorExpr{
													X:   ast.NewIdent("s"),
													Sel: ast.NewIdent("server"),
												},
												Sel: ast.NewIdent("RegisterService"),
											},
											Args: []ast.Expr{
												ast.NewIdent("sd"),
												ast.NewIdent("ss"),
											},
										},
									},
								},
							},
						},
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("reflection"),
									Sel: ast.NewIdent("Register"),
								},
								Args: []ast.Expr{
									&ast.SelectorExpr{
										X:   ast.NewIdent("s"),
										Sel: ast.NewIdent("server"),
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("healthServer"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("health"),
										Sel: ast.NewIdent("NewServer"),
									},
								},
							},
						},
						&ast.RangeStmt{
							Key: ast.NewIdent("service"),
							Tok: token.DEFINE,
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.SelectorExpr{
										X:   ast.NewIdent("s"),
										Sel: ast.NewIdent("server"),
									},
									Sel: ast.NewIdent("GetServiceInfo"),
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("healthServer"),
												Sel: ast.NewIdent("SetServingStatus"),
											},
											Args: []ast.Expr{
												ast.NewIdent("service"),
												&ast.SelectorExpr{
													X:   ast.NewIdent("grpc_health_v1"),
													Sel: ast.NewIdent("HealthCheckResponse_SERVING"),
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
									X:   ast.NewIdent("grpc_health_v1"),
									Sel: ast.NewIdent("RegisterHealthServer"),
								},
								Args: []ast.Expr{
									&ast.SelectorExpr{
										X:   ast.NewIdent("s"),
										Sel: ast.NewIdent("server"),
									},
									ast.NewIdent("healthServer"),
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("listener"),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("net"),
										Sel: ast.NewIdent("Listen"),
									},
									Args: []ast.Expr{
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: `"tcp"`,
										},
										&ast.SelectorExpr{
											X: &ast.SelectorExpr{
												X:   ast.NewIdent("s"),
												Sel: ast.NewIdent("config"),
											},
											Sel: ast.NewIdent("Address"),
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.BinaryExpr{
								X:  ast.NewIdent("err"),
								Op: token.NEQ,
								Y:  ast.NewIdent("nil"),
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											ast.NewIdent("err"),
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
											X:   ast.NewIdent("s"),
											Sel: ast.NewIdent("server"),
										},
										Sel: ast.NewIdent("Serve"),
									},
									Args: []ast.Expr{
										ast.NewIdent("listener"),
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
								ast.NewIdent("s"),
							},
							Type: &ast.StarExpr{
								X: ast.NewIdent("Server"),
							},
						},
					},
				},
				Name: ast.NewIdent("Stop"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("_"),
								},
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("context"),
									Sel: ast.NewIdent("Context"),
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
									X: &ast.SelectorExpr{
										X:   ast.NewIdent("s"),
										Sel: ast.NewIdent("server"),
									},
									Sel: ast.NewIdent("GracefulStop"),
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
			},
			&ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								ast.NewIdent("s"),
							},
							Type: &ast.StarExpr{
								X: ast.NewIdent("Server"),
							},
						},
					},
				},
				Name: ast.NewIdent("AddHandler"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("sd"),
								},
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X:   ast.NewIdent("grpc"),
										Sel: ast.NewIdent("ServiceDesc"),
									},
								},
							},
							{
								Names: []*ast.Ident{
									ast.NewIdent("ss"),
								},
								Type: ast.NewIdent("any"),
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
										X:   ast.NewIdent("s"),
										Sel: ast.NewIdent("handlers"),
									},
									Index: ast.NewIdent("sd"),
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								ast.NewIdent("ss"),
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
								ast.NewIdent("s"),
							},
							Type: &ast.StarExpr{
								X: ast.NewIdent("Server"),
							},
						},
					},
				},
				Name: ast.NewIdent("AddInterceptor"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("Interceptor"),
								},
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("grpc"),
									Sel: ast.NewIdent("UnaryServerInterceptor"),
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
									X:   ast.NewIdent("s"),
									Sel: ast.NewIdent("unaryInterceptors"),
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: ast.NewIdent("append"),
									Args: []ast.Expr{
										&ast.SelectorExpr{
											X:   ast.NewIdent("s"),
											Sel: ast.NewIdent("unaryInterceptors"),
										},
										ast.NewIdent("Interceptor"),
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
