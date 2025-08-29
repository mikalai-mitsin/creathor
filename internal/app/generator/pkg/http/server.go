package http

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
		Package: 1,
		Name:    ast.NewIdent("http"),
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"context\"",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"net/http\"",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"github.com/go-chi/chi/v5\"",
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
											ast.NewIdent("config"),
										},
										Type: &ast.StarExpr{
											X: ast.NewIdent("Config"),
										},
									},
									{
										Names: []*ast.Ident{
											ast.NewIdent("router"),
										},
										Type: &ast.StarExpr{
											X: &ast.SelectorExpr{
												X:   ast.NewIdent("chi"),
												Sel: ast.NewIdent("Mux"),
											},
										},
									},
									{
										Names: []*ast.Ident{
											ast.NewIdent("server"),
										},
										Type: &ast.StarExpr{
											X: &ast.SelectorExpr{
												X:   ast.NewIdent("http"),
												Sel: ast.NewIdent("Server"),
											},
										},
									},
									{
										Names: []*ast.Ident{
											ast.NewIdent("logger"),
										},
										Type: &ast.SelectorExpr{
											X:   ast.NewIdent("log"),
											Sel: ast.NewIdent("Logger"),
										},
									},
								},
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Doc: &ast.CommentGroup{
					List: []*ast.Comment{
						{
							Text: "// NewServer - provide http server",
						},
						{
							Text: "//",
						},
						{
							Text: fmt.Sprintf("// @title %s", u.project.Name),
						},
						{
							Text: "// @host http://127.0.0.1:8000",
						},
						{
							Text: "// @BasePath /",
						},
						{
							Text: "// @version 0.0.0",
						},
						{
							Text: "// @securitydefinitions.BearerAuth BearerAuth",
						},
					},
				},
				Name: ast.NewIdent("NewServer"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
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
									ast.NewIdent("logger"),
								},
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("log"),
									Sel: ast.NewIdent("Logger"),
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
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("router"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("chi"),
										Sel: ast.NewIdent("NewRouter"),
									},
								},
							},
						},
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "router",
									},
									Sel: &ast.Ident{
										Name: "Use",
									},
								},
								Args: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "otelchi",
											},
											Sel: &ast.Ident{
												Name: "Middleware",
											},
										},
										Args: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: fmt.Sprintf(`"%s"`, u.project.Name),
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
										Name: "router",
									},
									Sel: &ast.Ident{
										Name: "Use",
									},
								},
								Args: []ast.Expr{
									&ast.CallExpr{
										Fun: ast.NewIdent("loggerMiddleware"),
										Args: []ast.Expr{
											ast.NewIdent("logger"),
										},
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("server"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.UnaryExpr{
									Op: token.AND,
									X: &ast.CompositeLit{
										Type: &ast.SelectorExpr{
											X:   ast.NewIdent("http"),
											Sel: ast.NewIdent("Server"),
										},
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key: ast.NewIdent("Addr"),
												Value: &ast.SelectorExpr{
													X:   ast.NewIdent("config"),
													Sel: ast.NewIdent("Address"),
												},
											},
											&ast.KeyValueExpr{
												Key:   ast.NewIdent("Handler"),
												Value: ast.NewIdent("router"),
											},
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
										Type: ast.NewIdent("Server"),
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key:   ast.NewIdent("server"),
												Value: ast.NewIdent("server"),
											},
											&ast.KeyValueExpr{
												Key:   ast.NewIdent("config"),
												Value: ast.NewIdent("config"),
											},
											&ast.KeyValueExpr{
												Key:   ast.NewIdent("router"),
												Value: ast.NewIdent("router"),
											},
											&ast.KeyValueExpr{
												Key:   ast.NewIdent("logger"),
												Value: ast.NewIdent("logger"),
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
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X:   ast.NewIdent("s"),
											Sel: ast.NewIdent("server"),
										},
										Sel: ast.NewIdent("ListenAndServe"),
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
									ast.NewIdent("ctx"),
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
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X:   ast.NewIdent("s"),
											Sel: ast.NewIdent("server"),
										},
										Sel: ast.NewIdent("Shutdown"),
									},
									Args: []ast.Expr{
										ast.NewIdent("ctx"),
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
				Name: ast.NewIdent("Mount"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("path"),
								},
								Type: ast.NewIdent("string"),
							},
							{
								Names: []*ast.Ident{
									ast.NewIdent("handler"),
								},
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("http"),
									Sel: ast.NewIdent("Handler"),
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
										X:   ast.NewIdent("s"),
										Sel: ast.NewIdent("router"),
									},
									Sel: ast.NewIdent("Mount"),
								},
								Args: []ast.Expr{
									ast.NewIdent("path"),
									ast.NewIdent("handler"),
								},
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Name: &ast.Ident{
					Name: "loggerMiddleware",
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
								Type: &ast.FuncType{
									Params: &ast.FieldList{
										List: []*ast.Field{
											{
												Names: []*ast.Ident{
													{
														Name: "next",
													},
												},
												Type: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "http",
													},
													Sel: &ast.Ident{
														Name: "Handler",
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
														Name: "http",
													},
													Sel: &ast.Ident{
														Name: "Handler",
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
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.FuncLit{
									Type: &ast.FuncType{
										Params: &ast.FieldList{
											List: []*ast.Field{
												{
													Names: []*ast.Ident{
														{
															Name: "next",
														},
													},
													Type: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "http",
														},
														Sel: &ast.Ident{
															Name: "Handler",
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
															Name: "http",
														},
														Sel: &ast.Ident{
															Name: "Handler",
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
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X: &ast.Ident{
																Name: "http",
															},
															Sel: &ast.Ident{
																Name: "HandlerFunc",
															},
														},
														Args: []ast.Expr{
															&ast.FuncLit{
																Type: &ast.FuncType{
																	Params: &ast.FieldList{
																		List: []*ast.Field{
																			{
																				Names: []*ast.Ident{
																					{
																						Name: "w",
																					},
																				},
																				Type: &ast.SelectorExpr{
																					X: &ast.Ident{
																						Name: "http",
																					},
																					Sel: &ast.Ident{
																						Name: "ResponseWriter",
																					},
																				},
																			},
																			{
																				Names: []*ast.Ident{
																					{
																						Name: "r",
																					},
																				},
																				Type: &ast.StarExpr{
																					X: &ast.SelectorExpr{
																						X: &ast.Ident{
																							Name: "http",
																						},
																						Sel: &ast.Ident{
																							Name: "Request",
																						},
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
																					Name: "start",
																				},
																			},
																			Tok: token.DEFINE,
																			Rhs: []ast.Expr{
																				&ast.CallExpr{
																					Fun: &ast.SelectorExpr{
																						X: &ast.Ident{
																							Name: "time",
																						},
																						Sel: &ast.Ident{
																							Name: "Now",
																						},
																					},
																				},
																			},
																		},
																		&ast.AssignStmt{
																			Lhs: []ast.Expr{
																				&ast.Ident{
																					Name: "ww",
																				},
																			},
																			Tok: token.DEFINE,
																			Rhs: []ast.Expr{
																				&ast.CallExpr{
																					Fun: &ast.SelectorExpr{
																						X: &ast.Ident{
																							Name: "middleware",
																						},
																						Sel: &ast.Ident{
																							Name: "NewWrapResponseWriter",
																						},
																					},
																					Args: []ast.Expr{
																						&ast.Ident{
																							Name: "w",
																						},
																						&ast.SelectorExpr{
																							X: &ast.Ident{
																								Name: "r",
																							},
																							Sel: &ast.Ident{
																								Name: "ProtoMajor",
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
																						Name: "next",
																					},
																					Sel: &ast.Ident{
																						Name: "ServeHTTP",
																					},
																				},
																				Args: []ast.Expr{
																					&ast.Ident{
																						Name: "ww",
																					},
																					&ast.Ident{
																						Name: "r",
																					},
																				},
																			},
																		},
																		&ast.ExprStmt{
																			X: &ast.CallExpr{
																				Fun: &ast.SelectorExpr{
																					X: &ast.CallExpr{
																						Fun: &ast.SelectorExpr{
																							X: &ast.Ident{
																								Name: "logger",
																							},
																							Sel: &ast.Ident{
																								Name: "WithContext",
																							},
																						},
																						Args: []ast.Expr{
																							&ast.CallExpr{
																								Fun: &ast.SelectorExpr{
																									X: &ast.Ident{
																										Name: "r",
																									},
																									Sel: &ast.Ident{
																										Name: "Context",
																									},
																								},
																							},
																						},
																					},
																					Sel: &ast.Ident{
																						Name: "Info",
																					},
																				},
																				Args: []ast.Expr{
																					&ast.BasicLit{
																						Kind:  token.STRING,
																						Value: "\"finished http request\"",
																					},
																					&ast.CallExpr{
																						Fun: &ast.SelectorExpr{
																							X: &ast.Ident{
																								Name: "log",
																							},
																							Sel: &ast.Ident{
																								Name: "String",
																							},
																						},
																						Args: []ast.Expr{
																							&ast.BasicLit{
																								Kind:  token.STRING,
																								Value: "\"system\"",
																							},
																							&ast.BasicLit{
																								Kind:  token.STRING,
																								Value: "\"http\"",
																							},
																						},
																					},
																					&ast.CallExpr{
																						Fun: &ast.SelectorExpr{
																							X: &ast.Ident{
																								Name: "log",
																							},
																							Sel: &ast.Ident{
																								Name: "String",
																							},
																						},
																						Args: []ast.Expr{
																							&ast.BasicLit{
																								Kind:  token.STRING,
																								Value: "\"http.method\"",
																							},
																							&ast.SelectorExpr{
																								X: &ast.Ident{
																									Name: "r",
																								},
																								Sel: &ast.Ident{
																									Name: "Method",
																								},
																							},
																						},
																					},
																					&ast.CallExpr{
																						Fun: &ast.SelectorExpr{
																							X: &ast.Ident{
																								Name: "log",
																							},
																							Sel: &ast.Ident{
																								Name: "String",
																							},
																						},
																						Args: []ast.Expr{
																							&ast.BasicLit{
																								Kind:  token.STRING,
																								Value: "\"http.path\"",
																							},
																							&ast.SelectorExpr{
																								X: &ast.SelectorExpr{
																									X: &ast.Ident{
																										Name: "r",
																									},
																									Sel: &ast.Ident{
																										Name: "URL",
																									},
																								},
																								Sel: &ast.Ident{
																									Name: "Path",
																								},
																							},
																						},
																					},
																					&ast.CallExpr{
																						Fun: &ast.SelectorExpr{
																							X: &ast.Ident{
																								Name: "log",
																							},
																							Sel: &ast.Ident{
																								Name: "String",
																							},
																						},
																						Args: []ast.Expr{
																							&ast.BasicLit{
																								Kind:  token.STRING,
																								Value: "\"http.remote_addr\"",
																							},
																							&ast.SelectorExpr{
																								X: &ast.Ident{
																									Name: "r",
																								},
																								Sel: &ast.Ident{
																									Name: "RemoteAddr",
																								},
																							},
																						},
																					},
																					&ast.CallExpr{
																						Fun: &ast.SelectorExpr{
																							X: &ast.Ident{
																								Name: "log",
																							},
																							Sel: &ast.Ident{
																								Name: "Int",
																							},
																						},
																						Args: []ast.Expr{
																							&ast.BasicLit{
																								Kind:  token.STRING,
																								Value: "\"http.status\"",
																							},
																							&ast.CallExpr{
																								Fun: &ast.SelectorExpr{
																									X: &ast.Ident{
																										Name: "ww",
																									},
																									Sel: &ast.Ident{
																										Name: "Status",
																									},
																								},
																							},
																						},
																					},
																					&ast.CallExpr{
																						Fun: &ast.SelectorExpr{
																							X: &ast.Ident{
																								Name: "log",
																							},
																							Sel: &ast.Ident{
																								Name: "Int64",
																							},
																						},
																						Args: []ast.Expr{
																							&ast.BasicLit{
																								Kind:  token.STRING,
																								Value: "\"http.time_ms\"",
																							},
																							&ast.CallExpr{
																								Fun: &ast.SelectorExpr{
																									X: &ast.CallExpr{
																										Fun: &ast.SelectorExpr{
																											X: &ast.Ident{
																												Name: "time",
																											},
																											Sel: &ast.Ident{
																												Name: "Since",
																											},
																										},
																										Args: []ast.Expr{
																											&ast.Ident{
																												Name: "start",
																											},
																										},
																									},
																									Sel: &ast.Ident{
																										Name: "Milliseconds",
																									},
																								},
																							},
																						},
																					},
																					&ast.CallExpr{
																						Fun: &ast.SelectorExpr{
																							X: &ast.Ident{
																								Name: "log",
																							},
																							Sel: &ast.Ident{
																								Name: "String",
																							},
																						},
																						Args: []ast.Expr{
																							&ast.BasicLit{
																								Kind:  token.STRING,
																								Value: "\"http.start_time\"",
																							},
																							&ast.CallExpr{
																								Fun: &ast.SelectorExpr{
																									X: &ast.Ident{
																										Name: "start",
																									},
																									Sel: &ast.Ident{
																										Name: "Format",
																									},
																								},
																								Args: []ast.Expr{
																									&ast.SelectorExpr{
																										X: &ast.Ident{
																											Name: "time",
																										},
																										Sel: &ast.Ident{
																											Name: "RFC3339",
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
	filename := path.Join("internal", "pkg", "http", "server.go")
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
