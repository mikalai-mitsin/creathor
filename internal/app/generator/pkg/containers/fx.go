package containers

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
	if f.project.GRPCEnabled || f.project.HTTPEnabled {
		if err := f.syncServerContainer(); err != nil {
			return err
		}
	}
	if f.project.GatewayEnabled {
		if err := f.syncGatewayContainer(); err != nil {
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
				Value: f.project.PostgresImportPath(),
			},
		},

		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: f.project.ClockImportPath(),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: f.project.UUIDImportPath(),
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
				Value: f.project.LogImportPath(),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: f.project.ConfigsImportPath(),
			},
		},
	}
	if f.project.GRPCEnabled {
		imports = append(imports, &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: f.project.GRPCImportPath(),
			},
		})
	}
	if f.project.KafkaEnabled {
		imports = append(imports, &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: f.project.KafkaImportPath(),
			},
		})
	}
	if f.project.HTTPEnabled {
		imports = append(imports, &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: f.project.HTTPImportPath(),
			},
		})
	}
	if f.project.UptraceEnabled {
		imports = append(imports, &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: f.project.UptraceImportPath(),
			},
		})
	}
	if f.project.GRPCEnabled && f.project.GatewayEnabled {
		imports = append(imports, &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: f.project.GatewayImportPath(),
			},
		})
	}
	for _, modelConfig := range f.project.Apps {
		imports = append(
			imports,
			&ast.ImportSpec{
				Name: ast.NewIdent(modelConfig.AppAlias()),
				Path: &ast.BasicLit{
					Kind: token.STRING,
					Value: fmt.Sprintf(
						`"%s/internal/app/%s"`,
						f.project.Module,
						modelConfig.AppName(),
					),
				},
			},
		)
	}
	return &ast.File{
		Name: ast.NewIdent("containers"),
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
			X:   ast.NewIdent("context"),
			Sel: ast.NewIdent("Background"),
		},
		&ast.SelectorExpr{
			X:   ast.NewIdent("configs"),
			Sel: ast.NewIdent("ParseConfig"),
		},
		&ast.SelectorExpr{
			X:   ast.NewIdent("clock"),
			Sel: ast.NewIdent("NewClock"),
		},
		&ast.SelectorExpr{
			X:   ast.NewIdent("uuid"),
			Sel: ast.NewIdent("NewUUIDv7Generator"),
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
									X:   ast.NewIdent("postgres"),
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
								Sel: ast.NewIdent("Database"),
							},
						},
					},
				},
			},
		},
		&ast.SelectorExpr{
			X:   ast.NewIdent("postgres"),
			Sel: ast.NewIdent("NewDatabase"),
		},
		&ast.SelectorExpr{
			X:   ast.NewIdent("postgres"),
			Sel: ast.NewIdent("NewMigrateManager"),
		},
	}
	if f.project.KafkaEnabled {
		toProvide = append(
			toProvide,
			&ast.SelectorExpr{
				X:   ast.NewIdent("kafka"),
				Sel: ast.NewIdent("NewConsumer"),
			},
			&ast.SelectorExpr{
				X:   ast.NewIdent("kafka"),
				Sel: ast.NewIdent("NewProducer"),
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
										X:   ast.NewIdent("kafka"),
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
									Sel: ast.NewIdent("Kafka"),
								},
							},
						},
					},
				},
			},
		)
	}
	if f.project.UptraceEnabled {
		toProvide = append(
			toProvide,
			&ast.SelectorExpr{
				X:   ast.NewIdent("uptrace"),
				Sel: ast.NewIdent("NewProvider"),
			},
		)
	}
	for _, model := range f.project.Apps {
		toProvide = append(
			toProvide,
			&ast.SelectorExpr{
				X:   ast.NewIdent(model.AppAlias()),
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
									X:   ast.NewIdent("log"),
									Sel: ast.NewIdent("Log"),
								},
							},
						},
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
									X:   ast.NewIdent("log"),
									Sel: ast.NewIdent("NewLog"),
								},
								Args: []ast.Expr{
									&ast.SelectorExpr{
										X:   ast.NewIdent("config"),
										Sel: ast.NewIdent("LogLevel"),
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
				X:   ast.NewIdent("fx"),
				Sel: ast.NewIdent("WithLogger"),
			},
			Args: []ast.Expr{
				&ast.FuncLit{
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
							},
						},
						Results: &ast.FieldList{
							List: []*ast.Field{
								{
									Type: &ast.SelectorExpr{
										X:   ast.NewIdent("fxevent"),
										Sel: ast.NewIdent("Logger"),
									},
								},
							},
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{
									ast.NewIdent("logger"),
								},
							},
						},
					},
				},
			},
		},
		&ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("fx"),
				Sel: ast.NewIdent("Provide"),
			},
			Args: toProvide,
		},
	}
	if f.project.UptraceEnabled {
		exprs = append(exprs, &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("fx"),
				Sel: ast.NewIdent("Invoke"),
			},
			Args: []ast.Expr{
				&ast.FuncLit{
					Type: &ast.FuncType{
						Params: &ast.FieldList{
							List: []*ast.Field{
								{
									Names: []*ast.Ident{
										ast.NewIdent("lifecycle"),
									},
									Type: &ast.SelectorExpr{
										X:   ast.NewIdent("fx"),
										Sel: ast.NewIdent("Lifecycle"),
									},
								},
								{
									Names: []*ast.Ident{
										ast.NewIdent("server"),
									},
									Type: &ast.StarExpr{
										X: &ast.SelectorExpr{
											X:   ast.NewIdent("uptrace"),
											Sel: ast.NewIdent("Provider"),
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
							},
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("lifecycle"),
										Sel: ast.NewIdent("Append"),
									},
									Args: []ast.Expr{
										&ast.CompositeLit{
											Type: &ast.SelectorExpr{
												X:   ast.NewIdent("fx"),
												Sel: ast.NewIdent("Hook"),
											},
											Elts: []ast.Expr{
												&ast.KeyValueExpr{
													Key: ast.NewIdent("OnStart"),
													Value: &ast.SelectorExpr{
														X:   ast.NewIdent("server"),
														Sel: ast.NewIdent("Start"),
													},
												},
												&ast.KeyValueExpr{
													Key: ast.NewIdent("OnStop"),
													Value: &ast.SelectorExpr{
														X:   ast.NewIdent("server"),
														Sel: ast.NewIdent("Stop"),
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
	if f.project.KafkaEnabled {
		exprs = append(exprs,
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("fx"),
					Sel: ast.NewIdent("Invoke"),
				},
				Args: []ast.Expr{
					&ast.FuncLit{
						Type: &ast.FuncType{
							Params: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											ast.NewIdent("lifecycle"),
										},
										Type: &ast.SelectorExpr{
											X:   ast.NewIdent("fx"),
											Sel: ast.NewIdent("Lifecycle"),
										},
									},
									{
										Names: []*ast.Ident{
											ast.NewIdent("producer"),
										},
										Type: &ast.StarExpr{
											X: &ast.SelectorExpr{
												X:   ast.NewIdent("kafka"),
												Sel: ast.NewIdent("Producer"),
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
											X:   ast.NewIdent("lifecycle"),
											Sel: ast.NewIdent("Append"),
										},
										Args: []ast.Expr{
											&ast.CompositeLit{
												Type: &ast.SelectorExpr{
													X:   ast.NewIdent("fx"),
													Sel: ast.NewIdent("Hook"),
												},
												Elts: []ast.Expr{
													&ast.KeyValueExpr{
														Key: ast.NewIdent("OnStart"),
														Value: &ast.SelectorExpr{
															X:   ast.NewIdent("producer"),
															Sel: ast.NewIdent("Start"),
														},
													},
													&ast.KeyValueExpr{
														Key: ast.NewIdent("OnStop"),
														Value: &ast.SelectorExpr{
															X:   ast.NewIdent("producer"),
															Sel: ast.NewIdent("Stop"),
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
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("fx"),
					Sel: ast.NewIdent("Invoke"),
				},
				Args: []ast.Expr{
					&ast.FuncLit{
						Type: &ast.FuncType{
							Params: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											ast.NewIdent("lifecycle"),
										},
										Type: &ast.SelectorExpr{
											X:   ast.NewIdent("fx"),
											Sel: ast.NewIdent("Lifecycle"),
										},
									},
									{
										Names: []*ast.Ident{
											ast.NewIdent("consumer"),
										},
										Type: &ast.StarExpr{
											X: &ast.SelectorExpr{
												X:   ast.NewIdent("kafka"),
												Sel: ast.NewIdent("Consumer"),
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
											X:   ast.NewIdent("lifecycle"),
											Sel: ast.NewIdent("Append"),
										},
										Args: []ast.Expr{
											&ast.CompositeLit{
												Type: &ast.SelectorExpr{
													X:   ast.NewIdent("fx"),
													Sel: ast.NewIdent("Hook"),
												},
												Elts: []ast.Expr{
													&ast.KeyValueExpr{
														Key: ast.NewIdent("OnStart"),
														Value: &ast.SelectorExpr{
															X:   ast.NewIdent("consumer"),
															Sel: ast.NewIdent("Start"),
														},
													},
													&ast.KeyValueExpr{
														Key: ast.NewIdent("OnStop"),
														Value: &ast.SelectorExpr{
															X:   ast.NewIdent("consumer"),
															Sel: ast.NewIdent("Stop"),
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
		)
	}
	return &ast.ValueSpec{
		Names: []*ast.Ident{
			ast.NewIdent("FXModule"),
		},
		Values: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("fx"),
					Sel: ast.NewIdent("Options"),
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

func (f FxContainer) astServerContainer() *ast.FuncDecl {
	args := []ast.Expr{
		&ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("fx"),
				Sel: ast.NewIdent("Provide"),
			},
			Args: []ast.Expr{
				&ast.FuncLit{
					Type: &ast.FuncType{
						Params: &ast.FieldList{},
						Results: &ast.FieldList{
							List: []*ast.Field{
								{
									Type: ast.NewIdent("string"),
								},
							},
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{
									ast.NewIdent("config"),
								},
							},
						},
					},
				},
			},
		},
		ast.NewIdent("FXModule"),
	}
	if f.project.GRPCEnabled {
		args = append(
			args,
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("fx"),
					Sel: ast.NewIdent("Provide"),
				},
				Args: []ast.Expr{
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
												X:   ast.NewIdent("grpc"),
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
											Sel: ast.NewIdent("GRPC"),
										},
									},
								},
							},
						},
					},
					&ast.SelectorExpr{
						X:   ast.NewIdent("grpc"),
						Sel: ast.NewIdent("NewServer"),
					},
				},
			},
		)
		for _, domain := range f.project.Apps {
			args = append(args, &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("fx"),
					Sel: ast.NewIdent("Invoke"),
				},
				Args: []ast.Expr{
					&ast.FuncLit{
						Type: &ast.FuncType{
							Params: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											ast.NewIdent("lifecycle"),
										},
										Type: &ast.SelectorExpr{
											X:   ast.NewIdent("fx"),
											Sel: ast.NewIdent("Lifecycle"),
										},
									},
									{
										Names: []*ast.Ident{
											ast.NewIdent("app"),
										},
										Type: &ast.StarExpr{
											X: &ast.SelectorExpr{
												X:   ast.NewIdent(domain.AppAlias()),
												Sel: ast.NewIdent("App"),
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
								},
							},
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.ExprStmt{
									X: &ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   ast.NewIdent("lifecycle"),
											Sel: ast.NewIdent("Append"),
										},
										Args: []ast.Expr{
											&ast.CompositeLit{
												Type: &ast.SelectorExpr{
													X:   ast.NewIdent("fx"),
													Sel: ast.NewIdent("Hook"),
												},
												Elts: []ast.Expr{
													&ast.KeyValueExpr{
														Key: ast.NewIdent("OnStart"),
														Value: &ast.FuncLit{
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
																	&ast.IfStmt{
																		Init: &ast.AssignStmt{
																			Lhs: []ast.Expr{
																				ast.NewIdent("err"),
																			},
																			Tok: token.DEFINE,
																			Rhs: []ast.Expr{
																				&ast.CallExpr{
																					Fun: &ast.SelectorExpr{
																						X:   ast.NewIdent("app"),
																						Sel: ast.NewIdent("RegisterGRPC"),
																					},
																					Args: []ast.Expr{
																						ast.NewIdent("server"),
																					},
																				},
																			},
																		},
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
																			ast.NewIdent("nil"),
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
				X:   ast.NewIdent("fx"),
				Sel: ast.NewIdent("Invoke"),
			},
			Args: []ast.Expr{
				&ast.FuncLit{
					Type: &ast.FuncType{
						Params: &ast.FieldList{
							List: []*ast.Field{
								{
									Names: []*ast.Ident{
										ast.NewIdent("lifecycle"),
									},
									Type: &ast.SelectorExpr{
										X:   ast.NewIdent("fx"),
										Sel: ast.NewIdent("Lifecycle"),
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
										ast.NewIdent("shutdowner"),
									},
									Type: &ast.SelectorExpr{
										X:   ast.NewIdent("fx"),
										Sel: ast.NewIdent("Shutdowner"),
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
										X:   ast.NewIdent("lifecycle"),
										Sel: ast.NewIdent("Append"),
									},
									Args: []ast.Expr{
										&ast.CompositeLit{
											Type: &ast.SelectorExpr{
												X:   ast.NewIdent("fx"),
												Sel: ast.NewIdent("Hook"),
											},
											Elts: []ast.Expr{
												&ast.KeyValueExpr{
													Key: ast.NewIdent("OnStart"),
													Value: &ast.FuncLit{
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
																							ast.NewIdent("err"),
																						},
																						Tok: token.DEFINE,
																						Rhs: []ast.Expr{
																							&ast.CallExpr{
																								Fun: &ast.SelectorExpr{
																									X:   ast.NewIdent("server"),
																									Sel: ast.NewIdent("Start"),
																								},
																								Args: []ast.Expr{
																									ast.NewIdent("ctx"),
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
																								&ast.ExprStmt{
																									X: &ast.CallExpr{
																										Fun: &ast.SelectorExpr{
																											X:   ast.NewIdent("logger"),
																											Sel: ast.NewIdent("Error"),
																										},
																										Args: []ast.Expr{
																											&ast.BasicLit{
																												Kind:  token.STRING,
																												Value: `"shutdown"`,
																											},
																											&ast.CallExpr{
																												Fun: &ast.SelectorExpr{
																													X:   ast.NewIdent("log"),
																													Sel: ast.NewIdent("Any"),
																												},
																												Args: []ast.Expr{
																													&ast.BasicLit{
																														Kind:  token.STRING,
																														Value: `"error"`,
																													},
																													ast.NewIdent("err"),
																												},
																											},
																										},
																									},
																								},
																								&ast.AssignStmt{
																									Lhs: []ast.Expr{
																										ast.NewIdent("_"),
																									},
																									Tok: token.ASSIGN,
																									Rhs: []ast.Expr{
																										&ast.CallExpr{
																											Fun: &ast.SelectorExpr{
																												X:   ast.NewIdent("shutdowner"),
																												Sel: ast.NewIdent("Shutdown"),
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
																		ast.NewIdent("nil"),
																	},
																},
															},
														},
													},
												},
												&ast.KeyValueExpr{
													Key: ast.NewIdent("OnStop"),
													Value: &ast.SelectorExpr{
														X:   ast.NewIdent("server"),
														Sel: ast.NewIdent("Stop"),
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
	if f.project.HTTPEnabled {
		args = append(
			args,
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("fx"),
					Sel: ast.NewIdent("Provide"),
				},
				Args: []ast.Expr{
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
						X:   ast.NewIdent("http"),
						Sel: ast.NewIdent("NewServer"),
					},
				},
			},
		)

		for _, domain := range f.project.Apps {
			args = append(args, &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("fx"),
					Sel: ast.NewIdent("Invoke"),
				},
				Args: []ast.Expr{
					&ast.FuncLit{
						Type: &ast.FuncType{
							Params: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											ast.NewIdent("lifecycle"),
										},
										Type: &ast.SelectorExpr{
											X:   ast.NewIdent("fx"),
											Sel: ast.NewIdent("Lifecycle"),
										},
									},
									{
										Names: []*ast.Ident{
											ast.NewIdent("app"),
										},
										Type: &ast.StarExpr{
											X: &ast.SelectorExpr{
												X:   ast.NewIdent(domain.AppAlias()),
												Sel: ast.NewIdent("App"),
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
								},
							},
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.ExprStmt{
									X: &ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   ast.NewIdent("lifecycle"),
											Sel: ast.NewIdent("Append"),
										},
										Args: []ast.Expr{
											&ast.CompositeLit{
												Type: &ast.SelectorExpr{
													X:   ast.NewIdent("fx"),
													Sel: ast.NewIdent("Hook"),
												},
												Elts: []ast.Expr{
													&ast.KeyValueExpr{
														Key: ast.NewIdent("OnStart"),
														Value: &ast.FuncLit{
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
																	&ast.IfStmt{
																		Init: &ast.AssignStmt{
																			Lhs: []ast.Expr{
																				ast.NewIdent("err"),
																			},
																			Tok: token.DEFINE,
																			Rhs: []ast.Expr{
																				&ast.CallExpr{
																					Fun: &ast.SelectorExpr{
																						X:   ast.NewIdent("app"),
																						Sel: ast.NewIdent("RegisterHTTP"),
																					},
																					Args: []ast.Expr{
																						ast.NewIdent("server"),
																					},
																				},
																			},
																		},
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
																			ast.NewIdent("nil"),
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
				X:   ast.NewIdent("fx"),
				Sel: ast.NewIdent("Invoke"),
			},
			Args: []ast.Expr{
				&ast.FuncLit{
					Type: &ast.FuncType{
						Params: &ast.FieldList{
							List: []*ast.Field{
								{
									Names: []*ast.Ident{
										ast.NewIdent("lifecycle"),
									},
									Type: &ast.SelectorExpr{
										X:   ast.NewIdent("fx"),
										Sel: ast.NewIdent("Lifecycle"),
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
										ast.NewIdent("shutdowner"),
									},
									Type: &ast.SelectorExpr{
										X:   ast.NewIdent("fx"),
										Sel: ast.NewIdent("Shutdowner"),
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
										X:   ast.NewIdent("lifecycle"),
										Sel: ast.NewIdent("Append"),
									},
									Args: []ast.Expr{
										&ast.CompositeLit{
											Type: &ast.SelectorExpr{
												X:   ast.NewIdent("fx"),
												Sel: ast.NewIdent("Hook"),
											},
											Elts: []ast.Expr{
												&ast.KeyValueExpr{
													Key: ast.NewIdent("OnStart"),
													Value: &ast.FuncLit{
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
																							ast.NewIdent("err"),
																						},
																						Tok: token.DEFINE,
																						Rhs: []ast.Expr{
																							&ast.CallExpr{
																								Fun: &ast.SelectorExpr{
																									X:   ast.NewIdent("server"),
																									Sel: ast.NewIdent("Start"),
																								},
																								Args: []ast.Expr{
																									ast.NewIdent("ctx"),
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
																								&ast.ExprStmt{
																									X: &ast.CallExpr{
																										Fun: &ast.SelectorExpr{
																											X:   ast.NewIdent("logger"),
																											Sel: ast.NewIdent("Error"),
																										},
																										Args: []ast.Expr{
																											&ast.BasicLit{
																												Kind:  token.STRING,
																												Value: `"shutdown"`,
																											},
																											&ast.CallExpr{
																												Fun: &ast.SelectorExpr{
																													X:   ast.NewIdent("log"),
																													Sel: ast.NewIdent("Any"),
																												},
																												Args: []ast.Expr{
																													&ast.BasicLit{
																														Kind:  token.STRING,
																														Value: `"error"`,
																													},
																													ast.NewIdent("err"),
																												},
																											},
																										},
																									},
																								},
																								&ast.AssignStmt{
																									Lhs: []ast.Expr{
																										ast.NewIdent("_"),
																									},
																									Tok: token.ASSIGN,
																									Rhs: []ast.Expr{
																										&ast.CallExpr{
																											Fun: &ast.SelectorExpr{
																												X:   ast.NewIdent("shutdowner"),
																												Sel: ast.NewIdent("Shutdown"),
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
																		ast.NewIdent("nil"),
																	},
																},
															},
														},
													},
												},
												&ast.KeyValueExpr{
													Key: ast.NewIdent("OnStop"),
													Value: &ast.SelectorExpr{
														X:   ast.NewIdent("server"),
														Sel: ast.NewIdent("Stop"),
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
	return &ast.FuncDecl{
		Name: ast.NewIdent("NewServerContainer"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("config"),
						},
						Type: ast.NewIdent("string"),
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("fx"),
								Sel: ast.NewIdent("App"),
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
						ast.NewIdent("app"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("fx"),
								Sel: ast.NewIdent("New"),
							},
							Args: args,
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						ast.NewIdent("app"),
					},
				},
			},
		},
	}
}

func (f FxContainer) syncServerContainer() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "pkg", "containers", "fx.go")
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var functionExists bool
	var function *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == "NewServerContainer" {
			functionExists = true
			function = t
			return false
		}
		return true
	})
	if function == nil {
		function = f.astServerContainer()
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
		Name: ast.NewIdent("NewGatewayContainer"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("config"),
						},
						Type: ast.NewIdent("string"),
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("fx"),
								Sel: ast.NewIdent("App"),
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
						ast.NewIdent("app"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("fx"),
								Sel: ast.NewIdent("New"),
							},
							Args: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("fx"),
										Sel: ast.NewIdent("Provide"),
									},
									Args: []ast.Expr{
										&ast.FuncLit{
											Type: &ast.FuncType{
												Params: &ast.FieldList{},
												Results: &ast.FieldList{
													List: []*ast.Field{
														{
															Type: ast.NewIdent("string"),
														},
													},
												},
											},
											Body: &ast.BlockStmt{
												List: []ast.Stmt{
													&ast.ReturnStmt{
														Results: []ast.Expr{
															ast.NewIdent("config"),
														},
													},
												},
											},
										},
									},
								},
								&ast.SelectorExpr{
									X:   ast.NewIdent("gateway"),
									Sel: ast.NewIdent("NewServer"),
								},
								ast.NewIdent("FXModule"),
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("fx"),
										Sel: ast.NewIdent("Invoke"),
									},
									Args: []ast.Expr{
										&ast.FuncLit{
											Type: &ast.FuncType{
												Params: &ast.FieldList{
													List: []*ast.Field{
														{
															Names: []*ast.Ident{
																ast.NewIdent("lifecycle"),
															},
															Type: &ast.SelectorExpr{
																X:   ast.NewIdent("fx"),
																Sel: ast.NewIdent("Lifecycle"),
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
																ast.NewIdent("server"),
															},
															Type: &ast.StarExpr{
																X: &ast.SelectorExpr{
																	X:   ast.NewIdent("gateway"),
																	Sel: ast.NewIdent("Server"),
																},
															},
														},
														{
															Names: []*ast.Ident{
																ast.NewIdent("shutdowner"),
															},
															Type: &ast.SelectorExpr{
																X:   ast.NewIdent("fx"),
																Sel: ast.NewIdent("Shutdowner"),
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
																X:   ast.NewIdent("lifecycle"),
																Sel: ast.NewIdent("Append"),
															},
															Args: []ast.Expr{
																&ast.CompositeLit{
																	Type: &ast.SelectorExpr{
																		X:   ast.NewIdent("fx"),
																		Sel: ast.NewIdent("Hook"),
																	},
																	Elts: []ast.Expr{
																		&ast.KeyValueExpr{
																			Key: ast.NewIdent("OnStart"),
																			Value: &ast.FuncLit{
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
																													ast.NewIdent("err"),
																												},
																												Tok: token.DEFINE,
																												Rhs: []ast.Expr{
																													&ast.CallExpr{
																														Fun: &ast.SelectorExpr{
																															X:   ast.NewIdent("server"),
																															Sel: ast.NewIdent("Start"),
																														},
																														Args: []ast.Expr{
																															ast.NewIdent("ctx"),
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
																														&ast.ExprStmt{
																															X: &ast.CallExpr{
																																Fun: &ast.SelectorExpr{
																																	X:   ast.NewIdent("logger"),
																																	Sel: ast.NewIdent("Error"),
																																},
																																Args: []ast.Expr{
																																	&ast.BasicLit{
																																		Kind:  token.STRING,
																																		Value: `"shutdown"`,
																																	},
																																	&ast.CallExpr{
																																		Fun: &ast.SelectorExpr{
																																			X:   ast.NewIdent("log"),
																																			Sel: ast.NewIdent("Any"),
																																		},
																																		Args: []ast.Expr{
																																			&ast.BasicLit{
																																				Kind:  token.STRING,
																																				Value: `"error"`,
																																			},
																																			ast.NewIdent("err"),
																																		},
																																	},
																																},
																															},
																														},
																														&ast.AssignStmt{
																															Lhs: []ast.Expr{
																																ast.NewIdent("_"),
																															},
																															Tok: token.ASSIGN,
																															Rhs: []ast.Expr{
																																&ast.CallExpr{
																																	Fun: &ast.SelectorExpr{
																																		X:   ast.NewIdent("shutdowner"),
																																		Sel: ast.NewIdent("Shutdown"),
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
																								ast.NewIdent("nil"),
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
						ast.NewIdent("app"),
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
		Name: ast.NewIdent("NewMigrateContainer"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("config"),
						},
						Type: ast.NewIdent("string"),
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("fx"),
								Sel: ast.NewIdent("App"),
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
						ast.NewIdent("app"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("fx"),
								Sel: ast.NewIdent("New"),
							},
							Args: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("fx"),
										Sel: ast.NewIdent("Provide"),
									},
									Args: []ast.Expr{
										&ast.FuncLit{
											Type: &ast.FuncType{
												Params: &ast.FieldList{},
												Results: &ast.FieldList{
													List: []*ast.Field{
														{
															Type: ast.NewIdent("string"),
														},
													},
												},
											},
											Body: &ast.BlockStmt{
												List: []ast.Stmt{
													&ast.ReturnStmt{
														Results: []ast.Expr{
															ast.NewIdent("config"),
														},
													},
												},
											},
										},
									},
								},
								ast.NewIdent("FXModule"),
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("fx"),
										Sel: ast.NewIdent("Invoke"),
									},
									Args: []ast.Expr{
										&ast.FuncLit{
											Type: &ast.FuncType{
												Params: &ast.FieldList{
													List: []*ast.Field{
														{
															Names: []*ast.Ident{
																ast.NewIdent("lifecycle"),
															},
															Type: &ast.SelectorExpr{
																X:   ast.NewIdent("fx"),
																Sel: ast.NewIdent("Lifecycle"),
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
																ast.NewIdent("manager"),
															},
															Type: &ast.StarExpr{
																X: &ast.SelectorExpr{
																	X:   ast.NewIdent("postgres"),
																	Sel: ast.NewIdent("MigrateManager"),
																},
															},
														},
														{
															Names: []*ast.Ident{
																ast.NewIdent("shutdowner"),
															},
															Type: &ast.SelectorExpr{
																X:   ast.NewIdent("fx"),
																Sel: ast.NewIdent("Shutdowner"),
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
																X:   ast.NewIdent("lifecycle"),
																Sel: ast.NewIdent("Append"),
															},
															Args: []ast.Expr{
																&ast.CompositeLit{
																	Type: &ast.SelectorExpr{
																		X:   ast.NewIdent("fx"),
																		Sel: ast.NewIdent("Hook"),
																	},
																	Elts: []ast.Expr{
																		&ast.KeyValueExpr{
																			Key: ast.NewIdent("OnStart"),
																			Value: &ast.FuncLit{
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
																													ast.NewIdent("err"),
																												},
																												Tok: token.DEFINE,
																												Rhs: []ast.Expr{
																													&ast.CallExpr{
																														Fun: &ast.SelectorExpr{
																															X:   ast.NewIdent("manager"),
																															Sel: ast.NewIdent("Up"),
																														},
																														Args: []ast.Expr{
																															ast.NewIdent("ctx"),
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
																														&ast.ExprStmt{
																															X: &ast.CallExpr{
																																Fun: &ast.SelectorExpr{
																																	X:   ast.NewIdent("logger"),
																																	Sel: ast.NewIdent("Error"),
																																},
																																Args: []ast.Expr{
																																	&ast.BasicLit{
																																		Kind:  token.STRING,
																																		Value: `"shutdown"`,
																																	},
																																	&ast.CallExpr{
																																		Fun: &ast.SelectorExpr{
																																			X:   ast.NewIdent("log"),
																																			Sel: ast.NewIdent("Any"),
																																		},
																																		Args: []ast.Expr{
																																			&ast.BasicLit{
																																				Kind:  token.STRING,
																																				Value: `"error"`,
																																			},
																																			ast.NewIdent("err"),
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
																								ast.NewIdent("_"),
																							},
																							Tok: token.ASSIGN,
																							Rhs: []ast.Expr{
																								&ast.CallExpr{
																									Fun: &ast.SelectorExpr{
																										X:   ast.NewIdent("shutdowner"),
																										Sel: ast.NewIdent("Shutdown"),
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
						ast.NewIdent("app"),
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
