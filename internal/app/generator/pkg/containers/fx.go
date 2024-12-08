package containers

import (
	"bytes"
	"fmt"
	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"
	"path/filepath"
)

type FxContainer struct {
	project *configs.Project
}

func NewFxContainer(project *configs.Project) *FxContainer {
	return &FxContainer{project: project}
}

func (f FxContainer) Sync() error {
	if err := f.syncFxModule(); err != nil {
		return err
	}
	if err := f.syncMigrateContainer(); err != nil {
		return err
	}
	if f.project.GRPCEnabled {
		if err := f.syncGrpcContainer(); err != nil {
			return err
		}
	}
	if f.project.GatewayEnabled {
		if err := f.syncGatewayContainer(); err != nil {
			return err
		}
	}
	if f.project.HTTPEnabled {
		if err := f.syncHttpContainer(); err != nil {
			return err
		}
	}
	return nil
}

func (f FxContainer) filename() string {
	return filepath.Join("internal", "pkg", "containers", "fx.go")
}

func (f FxContainer) file() *ast.File {
	imports := []ast.Spec{
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: `"context"`,
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s/internal/pkg/postgres"`, f.project.Module),
			},
		},

		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s/internal/pkg/clock"`, f.project.Module),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s/internal/pkg/uuid"`, f.project.Module),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: `"go.uber.org/fx/fxevent"`,
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: `"go.uber.org/fx"`,
			},
		},

		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s/internal/pkg/clock"`, f.project.Module),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s/internal/pkg/configs"`, f.project.Module),
			},
		},
	}
	if f.project.Auth {
		imports = append(
			imports,
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind: token.STRING,
					Value: fmt.Sprintf(
						`"%s/internal/app/auth"`,
						f.project.Module,
					),
				},
			},
		)
	}
	if f.project.GRPCEnabled {
		imports = append(imports, &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s/internal/pkg/grpc"`, f.project.Module),
			},
		})
	}
	if f.project.HTTPEnabled {
		imports = append(imports, &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s/internal/pkg/http"`, f.project.Module),
			},
		})
	}
	if f.project.UptraceEnabled {
		imports = append(imports, &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s/internal/pkg/uptrace"`, f.project.Module),
			},
		})
	}
	if f.project.GRPCEnabled && f.project.GatewayEnabled {
		imports = append(imports, &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s/internal/pkg/gateway"`, f.project.Module),
			},
		})
	}
	for _, modelConfig := range f.project.Domains {
		imports = append(
			imports,
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind: token.STRING,
					Value: fmt.Sprintf(
						`"%s/internal/app/%s"`,
						f.project.Module,
						modelConfig.DomainName(),
					),
				},
			},
		)
	}
	return &ast.File{
		Name: &ast.Ident{
			Name: "containers",
		},
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok:   token.IMPORT,
				Specs: imports,
			},
		},
		Imports:  nil,
		Comments: nil,
	}
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
				Name: "NewClock",
			},
		},
		&ast.SelectorExpr{
			X: &ast.Ident{
				Name: "uuid",
			},
			Sel: &ast.Ident{
				Name: "NewUUIDv4Generator",
			},
		},
		&ast.SelectorExpr{
			X: &ast.Ident{
				Name: "postgres",
			},
			Sel: &ast.Ident{
				Name: "NewDatabase",
			},
		},
		&ast.SelectorExpr{
			X: &ast.Ident{
				Name: "postgres",
			},
			Sel: &ast.Ident{
				Name: "NewMigrateManager",
			},
		},
	}
	if f.project.UptraceEnabled {
		toProvide = append(
			toProvide,
			&ast.SelectorExpr{
				X: &ast.Ident{
					Name: "uptrace",
				},
				Sel: &ast.Ident{
					Name: "NewProvider",
				},
			},
		)
	}
	if f.project.Auth {
		toProvide = append(toProvide, &ast.SelectorExpr{
			X:   ast.NewIdent("auth"),
			Sel: ast.NewIdent("NewApp"),
		})
	}
	for _, model := range f.project.Domains {
		toProvide = append(
			toProvide,
			&ast.SelectorExpr{
				X:   ast.NewIdent(model.DomainAlias()),
				Sel: ast.NewIdent("NewApp"),
			},
		)
	}
	return toProvide
}

func (f FxContainer) astFxModule() *ast.ValueSpec {
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
	exprs := []ast.Expr{
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
	}
	if f.project.UptraceEnabled {
		exprs = append(exprs, &ast.CallExpr{
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
								{
									Names: []*ast.Ident{
										{
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
								{
									Names: []*ast.Ident{
										{
											Name: "server",
										},
									},
									Type: &ast.StarExpr{
										X: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "uptrace",
											},
											Sel: &ast.Ident{
												Name: "Provider",
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
													Value: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "server",
														},
														Sel: &ast.Ident{
															Name: "Start",
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
		})
	}
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
				Args: exprs,
			},
		},
	}
}

func (f FxContainer) syncFxModule() error {
	fileset := token.NewFileSet()
	filename := f.filename()
	if err := os.MkdirAll(path.Dir(filename), 0777); err != nil {
		return err
	}
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = f.file()
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
		fxModule = f.astFxModule()
	}
	for _, expr := range f.toProvide() {
		expr, ok := expr.(*ast.SelectorExpr)
		if ok {
			ast.Inspect(fxModule, func(node ast.Node) bool {
				if call, ok := node.(*ast.CallExpr); ok {
					if fun, ok := call.Fun.(*ast.SelectorExpr); ok &&
						fun.Sel.String() == "Provide" {
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

func (f FxContainer) astGrpcContainer() *ast.FuncDecl {
	args := []ast.Expr{
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
								{
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
				&ast.SelectorExpr{
					X: &ast.Ident{
						Name: "grpc",
					},
					Sel: &ast.Ident{
						Name: "NewServer",
					},
				},
			},
		},
		ast.NewIdent("FXModule"),
	}
	if f.project.Auth {
		args = append(args, &ast.CallExpr{
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
								{
									Names: []*ast.Ident{
										{
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
								{
									Names: []*ast.Ident{
										{
											Name: "app",
										},
									},
									Type: &ast.StarExpr{
										X: &ast.SelectorExpr{
											X:   ast.NewIdent("auth"),
											Sel: ast.NewIdent("App"),
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
																					X: &ast.Ident{
																						Name: "app",
																					},
																					Sel: &ast.Ident{
																						Name: "RegisterGRPC",
																					},
																				},
																				Args: []ast.Expr{
																					&ast.Ident{
																						Name: "server",
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
		})
	}
	for _, domain := range f.project.Domains {
		args = append(args, &ast.CallExpr{
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
								{
									Names: []*ast.Ident{
										{
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
								{
									Names: []*ast.Ident{
										{
											Name: "app",
										},
									},
									Type: &ast.StarExpr{
										X: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: domain.DomainName(),
											},
											Sel: ast.NewIdent("App"),
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
																					X: &ast.Ident{
																						Name: "app",
																					},
																					Sel: &ast.Ident{
																						Name: "RegisterGRPC",
																					},
																				},
																				Args: []ast.Expr{
																					&ast.Ident{
																						Name: "server",
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
		})
	}
	args = append(args, &ast.CallExpr{
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
							{
								Names: []*ast.Ident{
									{
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
																											Value: `"shutdown"`,
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
																													Value: `"error"`,
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
	})
	return &ast.FuncDecl{
		Name: &ast.Ident{
			Name: "NewGRPCContainer",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
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
					{
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
							Args: args,
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

func (f FxContainer) syncGrpcContainer() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "pkg", "containers", "fx.go")
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
		function = f.astGrpcContainer()
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

func (f FxContainer) astGatewayContainer() *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: &ast.Ident{
			Name: "NewGatewayContainer",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
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
					{
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
														{
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
								&ast.SelectorExpr{
									X: &ast.Ident{
										Name: "gateway",
									},
									Sel: &ast.Ident{
										Name: "NewServer",
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
														{
															Names: []*ast.Ident{
																{
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
																		Name: "gateway",
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
																																		Value: `"shutdown"`,
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
																																				Value: `"error"`,
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

func (f FxContainer) syncGatewayContainer() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "pkg", "containers", "fx.go")
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
		function = f.astGatewayContainer()
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

func (f FxContainer) astMigrateContainer() *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: &ast.Ident{
			Name: "NewMigrateContainer",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
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
					{
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
														{
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
														{
															Names: []*ast.Ident{
																{
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
																	Name: "manager",
																},
															},
															Type: &ast.StarExpr{
																X: &ast.SelectorExpr{
																	X: &ast.Ident{
																		Name: "postgres",
																	},
																	Sel: &ast.Ident{
																		Name: "MigrateManager",
																	},
																},
															},
														},
														{
															Names: []*ast.Ident{
																{
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
																																		Value: `"shutdown"`,
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
																																				Value: `"error"`,
																																			},
																																			&ast.Ident{
																																				Name: "err",
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

func (f FxContainer) syncMigrateContainer() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "pkg", "containers", "fx.go")
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
		function = f.astMigrateContainer()
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

func (f FxContainer) astHttpContainer() *ast.FuncDecl {
	args := []ast.Expr{
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
								{
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
				&ast.FuncLit{
					Type: &ast.FuncType{
						Params: &ast.FieldList{
							List: []*ast.Field{
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
							},
						},
						Results: &ast.FieldList{
							List: []*ast.Field{
								{
									Type: &ast.StarExpr{
										X: &ast.SelectorExpr{
											X:   ast.NewIdent("http"),
											Sel: ast.NewIdent("Config"),
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
									&ast.SelectorExpr{
										X:   ast.NewIdent("config"),
										Sel: ast.NewIdent("HTTP"),
									},
								},
							},
						},
					},
				},
				&ast.SelectorExpr{
					X: &ast.Ident{
						Name: "http",
					},
					Sel: &ast.Ident{
						Name: "NewServer",
					},
				},
			},
		},
		ast.NewIdent("FXModule"),
	}
	if f.project.Auth {
		args = append(args, &ast.CallExpr{
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
								{
									Names: []*ast.Ident{
										{
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
								{
									Names: []*ast.Ident{
										{
											Name: "app",
										},
									},
									Type: &ast.StarExpr{
										X: &ast.SelectorExpr{
											X:   ast.NewIdent("auth"),
											Sel: ast.NewIdent("App"),
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
												Name: "http",
											},
											Sel: &ast.Ident{
												Name: "Server",
											},
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
																					X: &ast.Ident{
																						Name: "app",
																					},
																					Sel: &ast.Ident{
																						Name: "RegisterHTTP",
																					},
																				},
																				Args: []ast.Expr{
																					&ast.Ident{
																						Name: "server",
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
		})
	}
	for _, domain := range f.project.Domains {
		args = append(args, &ast.CallExpr{
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
								{
									Names: []*ast.Ident{
										{
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
								{
									Names: []*ast.Ident{
										{
											Name: "app",
										},
									},
									Type: &ast.StarExpr{
										X: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: domain.DomainName(),
											},
											Sel: ast.NewIdent("App"),
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
												Name: "http",
											},
											Sel: &ast.Ident{
												Name: "Server",
											},
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
																					X: &ast.Ident{
																						Name: "app",
																					},
																					Sel: &ast.Ident{
																						Name: "RegisterHTTP",
																					},
																				},
																				Args: []ast.Expr{
																					&ast.Ident{
																						Name: "server",
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
		})
	}
	args = append(args, &ast.CallExpr{
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
							{
								Names: []*ast.Ident{
									{
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
											Name: "http",
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
																											Value: `"shutdown"`,
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
																													Value: `"error"`,
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
	})
	return &ast.FuncDecl{
		Name: ast.NewIdent("NewHTTPContainer"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
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
					{
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
							Args: args,
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

func (f FxContainer) syncHttpContainer() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "pkg", "containers", "fx.go")
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var functionExists bool
	var function *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == "NewHTTPContainer" {
			functionExists = true
			function = t
			return false
		}
		return true
	})
	if function == nil {
		function = f.astHttpContainer()
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
