package containers

import (
	"bytes"
	"fmt"
	"github.com/018bf/creathor/internal/models"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"
)

type FxContainer struct {
	Project *models.Project
}

func (f FxContainer) toProvide() []ast.Expr {
	toProvide := []ast.Expr{
		&ast.SelectorExpr{
			X: &ast.Ident{
				Name: "context",
			},
			Sel: &ast.Ident{
				Name: "Background",
			},
		},
		&ast.SelectorExpr{
			X: &ast.Ident{
				Name: "configs",
			},
			Sel: &ast.Ident{
				Name: "ParseConfig",
			},
		},
		&ast.SelectorExpr{
			X: &ast.Ident{
				Name: "clock",
			},
			Sel: &ast.Ident{
				Name: "NewRealClock",
			},
		},
		&ast.SelectorExpr{
			X: &ast.Ident{
				Name: "postgresInterface",
			},
			Sel: &ast.Ident{
				Name: "NewDatabase",
			},
		},
		&ast.SelectorExpr{
			X: &ast.Ident{
				Name: "postgresInterface",
			},
			Sel: &ast.Ident{
				Name: "NewMigrateManager",
			},
		},
	}
	if f.Project.GRPCEnabled {
		toProvide = append(toProvide, &ast.SelectorExpr{
			X: &ast.Ident{
				Name: "grpcInterface",
			},
			Sel: &ast.Ident{
				Name: "NewServer",
			},
		})
		if f.Project.Auth {
			toProvide = append(
				toProvide,
				&ast.SelectorExpr{
					X: &ast.Ident{
						Name: "grpcInterface",
					},
					Sel: &ast.Ident{
						Name: "NewAuthMiddleware",
					},
				},
				&ast.SelectorExpr{
					X: &ast.Ident{
						Name: "grpcInterface",
					},
					Sel: &ast.Ident{
						Name: "NewAuthServiceServer",
					},
				},
				&ast.SelectorExpr{
					X: &ast.Ident{
						Name: "grpcInterface",
					},
					Sel: &ast.Ident{
						Name: "NewUserServiceServer",
					},
				},
			)
		}
	}
	if f.Project.RESTEnabled {
		toProvide = append(toProvide, &ast.SelectorExpr{
			X: &ast.Ident{
				Name: "restInterface",
			},
			Sel: &ast.Ident{
				Name: "NewServer",
			},
		})
		if f.Project.Auth {
			toProvide = append(
				toProvide,
				&ast.SelectorExpr{
					X: &ast.Ident{
						Name: "restInterface",
					},
					Sel: &ast.Ident{
						Name: "NewAuthMiddleware",
					},
				},
				&ast.SelectorExpr{
					X: &ast.Ident{
						Name: "restInterface",
					},
					Sel: &ast.Ident{
						Name: "NewAuthHandler",
					},
				},
				&ast.SelectorExpr{
					X: &ast.Ident{
						Name: "restInterface",
					},
					Sel: &ast.Ident{
						Name: "NewUserHandler",
					},
				},
			)
		}
	}
	if f.Project.GatewayEnabled {
		toProvide = append(toProvide, &ast.SelectorExpr{
			X: &ast.Ident{
				Name: "gatewayInterface",
			},
			Sel: &ast.Ident{
				Name: "NewServer",
			},
		})
	}
	if f.Project.Auth {
		toProvide = append(
			toProvide,
			&ast.SelectorExpr{
				X: &ast.Ident{
					Name: "interceptors",
				},
				Sel: &ast.Ident{
					Name: "NewAuthInterceptor",
				},
			},
			&ast.SelectorExpr{
				X: &ast.Ident{
					Name: "usecases",
				},
				Sel: &ast.Ident{
					Name: "NewAuthUseCase",
				},
			},
			&ast.SelectorExpr{
				X: &ast.Ident{
					Name: "jwtRepositories",
				},
				Sel: &ast.Ident{
					Name: "NewAuthRepository",
				},
			},
			&ast.SelectorExpr{
				X: &ast.Ident{
					Name: "postgresRepositories",
				},
				Sel: &ast.Ident{
					Name: "NewPermissionRepository",
				},
			},
			&ast.SelectorExpr{
				X: &ast.Ident{
					Name: "interceptors",
				},
				Sel: &ast.Ident{
					Name: "NewUserInterceptor",
				},
			},
			&ast.SelectorExpr{
				X: &ast.Ident{
					Name: "usecases",
				},
				Sel: &ast.Ident{
					Name: "NewUserUseCase",
				},
			},
			&ast.SelectorExpr{
				X: &ast.Ident{
					Name: "postgresRepositories",
				},
				Sel: &ast.Ident{
					Name: "NewPostgresUserRepository",
				},
			},
		)
	}
	for _, model := range f.Project.Models {
		toProvide = append(
			toProvide,
			&ast.SelectorExpr{
				X: &ast.Ident{
					Name: "grpcInterface",
				},
				Sel: &ast.Ident{
					Name: fmt.Sprintf("New%s", model.GRPCHandlerTypeName()),
				},
			},
			&ast.SelectorExpr{
				X: &ast.Ident{
					Name: "restInterface",
				},
				Sel: &ast.Ident{
					Name: fmt.Sprintf("New%s", model.RESTHandlerTypeName()),
				},
			},
			&ast.SelectorExpr{
				X: &ast.Ident{
					Name: "interceptors",
				},
				Sel: &ast.Ident{
					Name: fmt.Sprintf("New%s", model.InterceptorTypeName()),
				},
			},
			&ast.SelectorExpr{
				X: &ast.Ident{
					Name: "usecases",
				},
				Sel: &ast.Ident{
					Name: fmt.Sprintf("New%s", model.UseCaseTypeName()),
				},
			},
			&ast.SelectorExpr{
				X: &ast.Ident{
					Name: "postgresRepositories",
				},
				Sel: &ast.Ident{
					Name: fmt.Sprintf("New%s", model.RepositoryTypeName()),
				},
			},
		)
	}
	return toProvide
}

func (f FxContainer) AstFxModule() *ast.ValueSpec {
	toProvide := []ast.Expr{
		&ast.FuncLit{
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
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
							Type: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "log",
								},
								Sel: &ast.Ident{
									Name: "Logger",
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
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "log",
									},
									Sel: &ast.Ident{
										Name: "NewLog",
									},
								},
								Args: []ast.Expr{
									&ast.SelectorExpr{
										X: &ast.Ident{
											Name: "config",
										},
										Sel: &ast.Ident{
											Name: "LogLevel",
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
	toProvide = append(toProvide, f.toProvide()...)
	return &ast.ValueSpec{
		Names: []*ast.Ident{
			{
				Name: "FXModule",
			},
		},
		Values: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "fx",
					},
					Sel: &ast.Ident{
						Name: "Options",
					},
				},
				Args: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "fx",
							},
							Sel: &ast.Ident{
								Name: "WithLogger",
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
												Type: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "fxevent",
													},
													Sel: &ast.Ident{
														Name: "Logger",
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
												&ast.Ident{
													Name: "logger",
												},
											},
										},
									},
								},
							},
						},
					},
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "fx",
							},
							Sel: &ast.Ident{
								Name: "Provide",
							},
						},
						Args: toProvide,
					},
				},
			},
		},
	}
}

func (f FxContainer) SyncFxModule() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "containers", "fx.go")
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var varExists bool
	var fxModule *ast.ValueSpec
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.ValueSpec); ok {
			for _, name := range t.Names {
				if name.String() == "FXModule" {
					fxModule = t
					varExists = true
					return false
				}
			}
		}
		return true
	})
	if fxModule == nil {
		fxModule = f.AstFxModule()
	}
	for _, expr := range f.toProvide() {
		expr, ok := expr.(*ast.SelectorExpr)
		if ok {
			ast.Inspect(fxModule, func(node ast.Node) bool {
				if call, ok := node.(*ast.CallExpr); ok {
					if fun, ok := call.Fun.(*ast.SelectorExpr); ok && fun.Sel.String() == "Provide" {
						for _, arg := range call.Args {
							arg := arg
							if argSelector, ok := arg.(*ast.SelectorExpr); ok {
								if argSelector.Sel.String() == expr.Sel.String() {
									return false
								}
							}
						}
						call.Args = append(call.Args, expr)
						return false
					}
				}
				return true
			})
		}
	}
	if !varExists {
		gd := &ast.GenDecl{
			Doc:    nil,
			TokPos: 0,
			Tok:    token.VAR,
			Lparen: 0,
			Specs:  []ast.Spec{fxModule},
			Rparen: 0,
		}
		file.Decls = append(file.Decls, gd)
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

func (f FxContainer) AstGrpcContainer() *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: &ast.Ident{
			Name: "NewGRPCContainer",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: []*ast.Ident{
							&ast.Ident{
								Name: "config",
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
					&ast.Field{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "fx",
								},
								Sel: &ast.Ident{
									Name: "App",
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
							Name: "app",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "fx",
								},
								Sel: &ast.Ident{
									Name: "New",
								},
							},
							Args: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "fx",
										},
										Sel: &ast.Ident{
											Name: "Provide",
										},
									},
									Args: []ast.Expr{
										&ast.FuncLit{
											Type: &ast.FuncType{
												Params: &ast.FieldList{},
												Results: &ast.FieldList{
													List: []*ast.Field{
														&ast.Field{
															Type: &ast.Ident{
																Name: "string",
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
																Name: "config",
															},
														},
													},
												},
											},
										},
									},
								},
								&ast.Ident{
									Name: "FXModule",
								},
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "fx",
										},
										Sel: &ast.Ident{
											Name: "Invoke",
										},
									},
									Args: []ast.Expr{
										&ast.FuncLit{
											Type: &ast.FuncType{
												Params: &ast.FieldList{
													List: []*ast.Field{
														&ast.Field{
															Names: []*ast.Ident{
																&ast.Ident{
																	Name: "lifecycle",
																},
															},
															Type: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "fx",
																},
																Sel: &ast.Ident{
																	Name: "Lifecycle",
																},
															},
														},
														&ast.Field{
															Names: []*ast.Ident{
																&ast.Ident{
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
														&ast.Field{
															Names: []*ast.Ident{
																&ast.Ident{
																	Name: "server",
																},
															},
															Type: &ast.StarExpr{
																X: &ast.SelectorExpr{
																	X: &ast.Ident{
																		Name: "grpcInterface",
																	},
																	Sel: &ast.Ident{
																		Name: "Server",
																	},
																},
															},
														},
														&ast.Field{
															Names: []*ast.Ident{
																&ast.Ident{
																	Name: "shutdowner",
																},
															},
															Type: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "fx",
																},
																Sel: &ast.Ident{
																	Name: "Shutdowner",
																},
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
																	Name: "lifecycle",
																},
																Sel: &ast.Ident{
																	Name: "Append",
																},
															},
															Args: []ast.Expr{
																&ast.CompositeLit{
																	Type: &ast.SelectorExpr{
																		X: &ast.Ident{
																			Name: "fx",
																		},
																		Sel: &ast.Ident{
																			Name: "Hook",
																		},
																	},
																	Elts: []ast.Expr{
																		&ast.KeyValueExpr{
																			Key: &ast.Ident{
																				Name: "OnStart",
																			},
																			Value: &ast.FuncLit{
																				Type: &ast.FuncType{
																					Params: &ast.FieldList{
																						List: []*ast.Field{
																							&ast.Field{
																								Names: []*ast.Ident{
																									&ast.Ident{
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
																							&ast.Field{
																								Type: &ast.Ident{
																									Name: "error",
																								},
																							},
																						},
																					},
																				},
																				Body: &ast.BlockStmt{
																					List: []ast.Stmt{
																						&ast.GoStmt{
																							Call: &ast.CallExpr{
																								Fun: &ast.FuncLit{
																									Type: &ast.FuncType{
																										Params: &ast.FieldList{},
																									},
																									Body: &ast.BlockStmt{
																										List: []ast.Stmt{
																											&ast.AssignStmt{
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
																																Name: "server",
																															},
																															Sel: &ast.Ident{
																																Name: "Start",
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
																														&ast.ExprStmt{
																															X: &ast.CallExpr{
																																Fun: &ast.SelectorExpr{
																																	X: &ast.Ident{
																																		Name: "logger",
																																	},
																																	Sel: &ast.Ident{
																																		Name: "Error",
																																	},
																																},
																																Args: []ast.Expr{
																																	&ast.BasicLit{
																																		Kind:  token.STRING,
																																		Value: "\"shutdown\"",
																																	},
																																	&ast.CallExpr{
																																		Fun: &ast.SelectorExpr{
																																			X: &ast.Ident{
																																				Name: "log",
																																			},
																																			Sel: &ast.Ident{
																																				Name: "Any",
																																			},
																																		},
																																		Args: []ast.Expr{
																																			&ast.BasicLit{
																																				Kind:  token.STRING,
																																				Value: "\"error\"",
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
																																	Name: "_",
																																},
																															},
																															Tok: token.ASSIGN,
																															Rhs: []ast.Expr{
																																&ast.CallExpr{
																																	Fun: &ast.SelectorExpr{
																																		X: &ast.Ident{
																																			Name: "shutdowner",
																																		},
																																		Sel: &ast.Ident{
																																			Name: "Shutdown",
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
																		&ast.KeyValueExpr{
																			Key: &ast.Ident{
																				Name: "OnStop",
																			},
																			Value: &ast.SelectorExpr{
																				X: &ast.Ident{
																					Name: "server",
																				},
																				Sel: &ast.Ident{
																					Name: "Stop",
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
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.Ident{
							Name: "app",
						},
					},
				},
			},
		},
	}
}

func (f FxContainer) SyncGrpcContainer() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "containers", "fx.go")
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var functionExists bool
	var function *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == "NewGRPCContainer" {
			functionExists = true
			function = t
			return false
		}
		return true
	})
	if function == nil {
		function = f.AstGrpcContainer()
	}
	if !functionExists {
		file.Decls = append(file.Decls, function)
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

func (f FxContainer) AstGatewayContainer() *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: &ast.Ident{
			Name: "NewGatewayContainer",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: []*ast.Ident{
							&ast.Ident{
								Name: "config",
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
					&ast.Field{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "fx",
								},
								Sel: &ast.Ident{
									Name: "App",
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
							Name: "app",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "fx",
								},
								Sel: &ast.Ident{
									Name: "New",
								},
							},
							Args: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "fx",
										},
										Sel: &ast.Ident{
											Name: "Provide",
										},
									},
									Args: []ast.Expr{
										&ast.FuncLit{
											Type: &ast.FuncType{
												Params: &ast.FieldList{},
												Results: &ast.FieldList{
													List: []*ast.Field{
														&ast.Field{
															Type: &ast.Ident{
																Name: "string",
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
																Name: "config",
															},
														},
													},
												},
											},
										},
									},
								},
								&ast.Ident{
									Name: "FXModule",
								},
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "fx",
										},
										Sel: &ast.Ident{
											Name: "Invoke",
										},
									},
									Args: []ast.Expr{
										&ast.FuncLit{
											Type: &ast.FuncType{
												Params: &ast.FieldList{
													List: []*ast.Field{
														&ast.Field{
															Names: []*ast.Ident{
																&ast.Ident{
																	Name: "lifecycle",
																},
															},
															Type: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "fx",
																},
																Sel: &ast.Ident{
																	Name: "Lifecycle",
																},
															},
														},
														&ast.Field{
															Names: []*ast.Ident{
																&ast.Ident{
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
														&ast.Field{
															Names: []*ast.Ident{
																&ast.Ident{
																	Name: "server",
																},
															},
															Type: &ast.StarExpr{
																X: &ast.SelectorExpr{
																	X: &ast.Ident{
																		Name: "gatewayInterface",
																	},
																	Sel: &ast.Ident{
																		Name: "Server",
																	},
																},
															},
														},
														&ast.Field{
															Names: []*ast.Ident{
																&ast.Ident{
																	Name: "shutdowner",
																},
															},
															Type: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "fx",
																},
																Sel: &ast.Ident{
																	Name: "Shutdowner",
																},
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
																	Name: "lifecycle",
																},
																Sel: &ast.Ident{
																	Name: "Append",
																},
															},
															Args: []ast.Expr{
																&ast.CompositeLit{
																	Type: &ast.SelectorExpr{
																		X: &ast.Ident{
																			Name: "fx",
																		},
																		Sel: &ast.Ident{
																			Name: "Hook",
																		},
																	},
																	Elts: []ast.Expr{
																		&ast.KeyValueExpr{
																			Key: &ast.Ident{
																				Name: "OnStart",
																			},
																			Value: &ast.FuncLit{
																				Type: &ast.FuncType{
																					Params: &ast.FieldList{
																						List: []*ast.Field{
																							&ast.Field{
																								Names: []*ast.Ident{
																									&ast.Ident{
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
																							&ast.Field{
																								Type: &ast.Ident{
																									Name: "error",
																								},
																							},
																						},
																					},
																				},
																				Body: &ast.BlockStmt{
																					List: []ast.Stmt{
																						&ast.GoStmt{
																							Call: &ast.CallExpr{
																								Fun: &ast.FuncLit{
																									Type: &ast.FuncType{
																										Params: &ast.FieldList{},
																									},
																									Body: &ast.BlockStmt{
																										List: []ast.Stmt{
																											&ast.AssignStmt{
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
																																Name: "server",
																															},
																															Sel: &ast.Ident{
																																Name: "Start",
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
																														&ast.ExprStmt{
																															X: &ast.CallExpr{
																																Fun: &ast.SelectorExpr{
																																	X: &ast.Ident{
																																		Name: "logger",
																																	},
																																	Sel: &ast.Ident{
																																		Name: "Error",
																																	},
																																},
																																Args: []ast.Expr{
																																	&ast.BasicLit{
																																		Kind:  token.STRING,
																																		Value: "\"shutdown\"",
																																	},
																																	&ast.CallExpr{
																																		Fun: &ast.SelectorExpr{
																																			X: &ast.Ident{
																																				Name: "log",
																																			},
																																			Sel: &ast.Ident{
																																				Name: "Any",
																																			},
																																		},
																																		Args: []ast.Expr{
																																			&ast.BasicLit{
																																				Kind:  token.STRING,
																																				Value: "\"error\"",
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
																																	Name: "_",
																																},
																															},
																															Tok: token.ASSIGN,
																															Rhs: []ast.Expr{
																																&ast.CallExpr{
																																	Fun: &ast.SelectorExpr{
																																		X: &ast.Ident{
																																			Name: "shutdowner",
																																		},
																																		Sel: &ast.Ident{
																																			Name: "Shutdown",
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
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.Ident{
							Name: "app",
						},
					},
				},
			},
		},
	}
}

func (f FxContainer) SyncGatewayContainer() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "containers", "fx.go")
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var functionExists bool
	var function *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == "NewGatewayContainer" {
			functionExists = true
			function = t
			return false
		}
		return true
	})
	if function == nil {
		function = f.AstGatewayContainer()
	}
	if !functionExists {
		file.Decls = append(file.Decls, function)
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

func (f FxContainer) AstRestContainer() *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: &ast.Ident{
			Name: "NewRESTContainer",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: []*ast.Ident{
							&ast.Ident{
								Name: "config",
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
					&ast.Field{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "fx",
								},
								Sel: &ast.Ident{
									Name: "App",
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
							Name: "app",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "fx",
								},
								Sel: &ast.Ident{
									Name: "New",
								},
							},
							Args: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "fx",
										},
										Sel: &ast.Ident{
											Name: "Provide",
										},
									},
									Args: []ast.Expr{
										&ast.FuncLit{
											Type: &ast.FuncType{
												Params: &ast.FieldList{},
												Results: &ast.FieldList{
													List: []*ast.Field{
														&ast.Field{
															Type: &ast.Ident{
																Name: "string",
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
																Name: "config",
															},
														},
													},
												},
											},
										},
									},
								},
								&ast.Ident{
									Name: "FXModule",
								},
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "fx",
										},
										Sel: &ast.Ident{
											Name: "Invoke",
										},
									},
									Args: []ast.Expr{
										&ast.FuncLit{
											Type: &ast.FuncType{
												Params: &ast.FieldList{
													List: []*ast.Field{
														&ast.Field{
															Names: []*ast.Ident{
																&ast.Ident{
																	Name: "lifecycle",
																},
															},
															Type: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "fx",
																},
																Sel: &ast.Ident{
																	Name: "Lifecycle",
																},
															},
														},
														&ast.Field{
															Names: []*ast.Ident{
																&ast.Ident{
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
														&ast.Field{
															Names: []*ast.Ident{
																&ast.Ident{
																	Name: "server",
																},
															},
															Type: &ast.StarExpr{
																X: &ast.SelectorExpr{
																	X: &ast.Ident{
																		Name: "restInterface",
																	},
																	Sel: &ast.Ident{
																		Name: "Server",
																	},
																},
															},
														},
														&ast.Field{
															Names: []*ast.Ident{
																&ast.Ident{
																	Name: "shutdowner",
																},
															},
															Type: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "fx",
																},
																Sel: &ast.Ident{
																	Name: "Shutdowner",
																},
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
																	Name: "lifecycle",
																},
																Sel: &ast.Ident{
																	Name: "Append",
																},
															},
															Args: []ast.Expr{
																&ast.CompositeLit{
																	Type: &ast.SelectorExpr{
																		X: &ast.Ident{
																			Name: "fx",
																		},
																		Sel: &ast.Ident{
																			Name: "Hook",
																		},
																	},
																	Elts: []ast.Expr{
																		&ast.KeyValueExpr{
																			Key: &ast.Ident{
																				Name: "OnStart",
																			},
																			Value: &ast.FuncLit{
																				Type: &ast.FuncType{
																					Params: &ast.FieldList{
																						List: []*ast.Field{
																							&ast.Field{
																								Names: []*ast.Ident{
																									&ast.Ident{
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
																							&ast.Field{
																								Type: &ast.Ident{
																									Name: "error",
																								},
																							},
																						},
																					},
																				},
																				Body: &ast.BlockStmt{
																					List: []ast.Stmt{
																						&ast.GoStmt{
																							Call: &ast.CallExpr{
																								Fun: &ast.FuncLit{
																									Type: &ast.FuncType{
																										Params: &ast.FieldList{},
																									},
																									Body: &ast.BlockStmt{
																										List: []ast.Stmt{
																											&ast.AssignStmt{
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
																																Name: "server",
																															},
																															Sel: &ast.Ident{
																																Name: "Start",
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
																														&ast.ExprStmt{
																															X: &ast.CallExpr{
																																Fun: &ast.SelectorExpr{
																																	X: &ast.Ident{
																																		Name: "logger",
																																	},
																																	Sel: &ast.Ident{
																																		Name: "Error",
																																	},
																																},
																																Args: []ast.Expr{
																																	&ast.BasicLit{
																																		Kind:  token.STRING,
																																		Value: "\"shutdown\"",
																																	},
																																	&ast.CallExpr{
																																		Fun: &ast.SelectorExpr{
																																			X: &ast.Ident{
																																				Name: "log",
																																			},
																																			Sel: &ast.Ident{
																																				Name: "Any",
																																			},
																																		},
																																		Args: []ast.Expr{
																																			&ast.BasicLit{
																																				Kind:  token.STRING,
																																				Value: "\"error\"",
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
																																	Name: "_",
																																},
																															},
																															Tok: token.ASSIGN,
																															Rhs: []ast.Expr{
																																&ast.CallExpr{
																																	Fun: &ast.SelectorExpr{
																																		X: &ast.Ident{
																																			Name: "shutdowner",
																																		},
																																		Sel: &ast.Ident{
																																			Name: "Shutdown",
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
																		&ast.KeyValueExpr{
																			Key: &ast.Ident{
																				Name: "OnStop",
																			},
																			Value: &ast.SelectorExpr{
																				X: &ast.Ident{
																					Name: "server",
																				},
																				Sel: &ast.Ident{
																					Name: "Stop",
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
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.Ident{
							Name: "app",
						},
					},
				},
			},
		},
	}
}

func (f FxContainer) SyncRestContainer() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "containers", "fx.go")
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var functionExists bool
	var function *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == "NewRESTContainer" {
			functionExists = true
			function = t
			return false
		}
		return true
	})
	if function == nil {
		function = f.AstRestContainer()
	}
	if !functionExists {
		file.Decls = append(file.Decls, function)
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

func (f FxContainer) AstMigrateContainer() *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: &ast.Ident{
			Name: "NewMigrateContainer",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: []*ast.Ident{
							&ast.Ident{
								Name: "config",
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
					&ast.Field{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "fx",
								},
								Sel: &ast.Ident{
									Name: "App",
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
							Name: "app",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "fx",
								},
								Sel: &ast.Ident{
									Name: "New",
								},
							},
							Args: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "fx",
										},
										Sel: &ast.Ident{
											Name: "Provide",
										},
									},
									Args: []ast.Expr{
										&ast.FuncLit{
											Type: &ast.FuncType{
												Params: &ast.FieldList{},
												Results: &ast.FieldList{
													List: []*ast.Field{
														&ast.Field{
															Type: &ast.Ident{
																Name: "string",
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
																Name: "config",
															},
														},
													},
												},
											},
										},
									},
								},
								&ast.Ident{
									Name: "FXModule",
								},
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "fx",
										},
										Sel: &ast.Ident{
											Name: "Invoke",
										},
									},
									Args: []ast.Expr{
										&ast.FuncLit{
											Type: &ast.FuncType{
												Params: &ast.FieldList{
													List: []*ast.Field{
														&ast.Field{
															Names: []*ast.Ident{
																&ast.Ident{
																	Name: "lifecycle",
																},
															},
															Type: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "fx",
																},
																Sel: &ast.Ident{
																	Name: "Lifecycle",
																},
															},
														},
														&ast.Field{
															Names: []*ast.Ident{
																&ast.Ident{
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
														&ast.Field{
															Names: []*ast.Ident{
																&ast.Ident{
																	Name: "manager",
																},
															},
															Type: &ast.StarExpr{
																X: &ast.SelectorExpr{
																	X: &ast.Ident{
																		Name: "postgresInterface",
																	},
																	Sel: &ast.Ident{
																		Name: "MigrateManager",
																	},
																},
															},
														},
														&ast.Field{
															Names: []*ast.Ident{
																&ast.Ident{
																	Name: "shutdowner",
																},
															},
															Type: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "fx",
																},
																Sel: &ast.Ident{
																	Name: "Shutdowner",
																},
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
																	Name: "lifecycle",
																},
																Sel: &ast.Ident{
																	Name: "Append",
																},
															},
															Args: []ast.Expr{
																&ast.CompositeLit{
																	Type: &ast.SelectorExpr{
																		X: &ast.Ident{
																			Name: "fx",
																		},
																		Sel: &ast.Ident{
																			Name: "Hook",
																		},
																	},
																	Elts: []ast.Expr{
																		&ast.KeyValueExpr{
																			Key: &ast.Ident{
																				Name: "OnStart",
																			},
																			Value: &ast.FuncLit{
																				Type: &ast.FuncType{
																					Params: &ast.FieldList{
																						List: []*ast.Field{
																							&ast.Field{
																								Names: []*ast.Ident{
																									&ast.Ident{
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
																							&ast.Field{
																								Type: &ast.Ident{
																									Name: "error",
																								},
																							},
																						},
																					},
																				},
																				Body: &ast.BlockStmt{
																					List: []ast.Stmt{
																						&ast.GoStmt{
																							Call: &ast.CallExpr{
																								Fun: &ast.FuncLit{
																									Type: &ast.FuncType{
																										Params: &ast.FieldList{},
																									},
																									Body: &ast.BlockStmt{
																										List: []ast.Stmt{
																											&ast.AssignStmt{
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
																																Name: "manager",
																															},
																															Sel: &ast.Ident{
																																Name: "Up",
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
																														&ast.ExprStmt{
																															X: &ast.CallExpr{
																																Fun: &ast.SelectorExpr{
																																	X: &ast.Ident{
																																		Name: "logger",
																																	},
																																	Sel: &ast.Ident{
																																		Name: "Error",
																																	},
																																},
																																Args: []ast.Expr{
																																	&ast.BasicLit{
																																		Kind:  token.STRING,
																																		Value: "\"shutdown\"",
																																	},
																																	&ast.CallExpr{
																																		Fun: &ast.SelectorExpr{
																																			X: &ast.Ident{
																																				Name: "log",
																																			},
																																			Sel: &ast.Ident{
																																				Name: "Any",
																																			},
																																		},
																																		Args: []ast.Expr{
																																			&ast.BasicLit{
																																				Kind:  token.STRING,
																																				Value: "\"error\"",
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
																																	Name: "_",
																																},
																															},
																															Tok: token.ASSIGN,
																															Rhs: []ast.Expr{
																																&ast.CallExpr{
																																	Fun: &ast.SelectorExpr{
																																		X: &ast.Ident{
																																			Name: "shutdowner",
																																		},
																																		Sel: &ast.Ident{
																																			Name: "Shutdown",
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
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.Ident{
							Name: "app",
						},
					},
				},
			},
		},
	}
}

func (f FxContainer) SyncMigrateContainer() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "containers", "fx.go")
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var functionExists bool
	var function *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == "NewMigrateContainer" {
			functionExists = true
			function = t
			return false
		}
		return true
	})
	if function == nil {
		function = f.AstMigrateContainer()
	}
	if !functionExists {
		file.Decls = append(file.Decls, function)
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
