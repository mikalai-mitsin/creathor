package kafka

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"

	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
)

type ConsumerGenerator struct {
	project *configs.Project
}

func NewConsumerGenerator(project *configs.Project) *ConsumerGenerator {
	return &ConsumerGenerator{project: project}
}

func (u ConsumerGenerator) file() *ast.File {
	return &ast.File{
		Package: 1,
		Name: &ast.Ident{
			Name: "kafka",
		},
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
							Value: "\"os\"",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"os/signal\"",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"syscall\"",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"github.com/IBM/sarama\"",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: u.project.ErrsImportPath(),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: u.project.LogImportPath(),
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{
							Name: "Handler",
						},
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											{
												Name: "Topic",
											},
										},
										Type: &ast.Ident{
											Name: "string",
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "GroupID",
											},
										},
										Type: &ast.Ident{
											Name: "string",
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "Handler",
											},
										},
										Type: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "sarama",
											},
											Sel: &ast.Ident{
												Name: "ConsumerGroupHandler",
											},
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "group",
											},
										},
										Type: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "sarama",
											},
											Sel: &ast.Ident{
												Name: "ConsumerGroup",
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
							Name: "Consumer",
						},
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											{
												Name: "config",
											},
										},
										Type: &ast.StarExpr{
											X: &ast.Ident{
												Name: "Config",
											},
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "client",
											},
										},
										Type: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "sarama",
											},
											Sel: &ast.Ident{
												Name: "Client",
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
											Key: &ast.Ident{
												Name: "string",
											},
											Value: &ast.Ident{
												Name: "Handler",
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
								},
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Name: &ast.Ident{
					Name: "NewConsumer",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{
										Name: "cfg",
									},
								},
								Type: &ast.StarExpr{
									X: &ast.Ident{
										Name: "Config",
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
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.StarExpr{
									X: &ast.Ident{
										Name: "Consumer",
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
									Name: "config",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "sarama",
										},
										Sel: &ast.Ident{
											Name: "NewConfig",
										},
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.SelectorExpr{
									X: &ast.Ident{
										Name: "config",
									},
									Sel: &ast.Ident{
										Name: "Version",
									},
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.SelectorExpr{
									X: &ast.Ident{
										Name: "sarama",
									},
									Sel: &ast.Ident{
										Name: "V2_1_0_0",
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "client",
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
											Name: "sarama",
										},
										Sel: &ast.Ident{
											Name: "NewClient",
										},
									},
									Args: []ast.Expr{
										&ast.SelectorExpr{
											X: &ast.Ident{
												Name: "cfg",
											},
											Sel: &ast.Ident{
												Name: "Brokers",
											},
										},
										&ast.Ident{
											Name: "config",
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
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X: &ast.Ident{
																Name: "errs",
															},
															Sel: &ast.Ident{
																Name: "NewUnexpectedBehaviorError",
															},
														},
														Args: []ast.Expr{
															&ast.BasicLit{
																Kind:  token.STRING,
																Value: "\"cant build kafka client\"",
															},
														},
													},
													Sel: &ast.Ident{
														Name: "WithCause",
													},
												},
												Args: []ast.Expr{
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
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.UnaryExpr{
									Op: token.AND,
									X: &ast.CompositeLit{
										Type: &ast.Ident{
											Name: "Consumer",
										},
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "config",
												},
												Value: &ast.Ident{
													Name: "cfg",
												},
											},
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "handlers",
												},
												Value: &ast.CallExpr{
													Fun: &ast.Ident{
														Name: "make",
													},
													Args: []ast.Expr{
														&ast.MapType{
															Key: &ast.Ident{
																Name: "string",
															},
															Value: &ast.Ident{
																Name: "Handler",
															},
														},
													},
												},
											},
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "client",
												},
												Value: &ast.Ident{
													Name: "client",
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
									Name: "c",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "Consumer",
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
										Name: "handler",
									},
								},
								Type: &ast.Ident{
									Name: "Handler",
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
											Name: "c",
										},
										Sel: &ast.Ident{
											Name: "handlers",
										},
									},
									Index: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "handler",
										},
										Sel: &ast.Ident{
											Name: "GroupID",
										},
									},
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.Ident{
									Name: "handler",
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
									Name: "c",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "Consumer",
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
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "logger",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.SelectorExpr{
									X: &ast.Ident{
										Name: "c",
									},
									Sel: &ast.Ident{
										Name: "logger",
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "sigterm",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.Ident{
										Name: "make",
									},
									Args: []ast.Expr{
										&ast.ChanType{
											Dir: 3,
											Value: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "os",
												},
												Sel: &ast.Ident{
													Name: "Signal",
												},
											},
										},
										&ast.BasicLit{
											Kind:  token.INT,
											Value: "1",
										},
									},
								},
							},
						},
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "signal",
									},
									Sel: &ast.Ident{
										Name: "Notify",
									},
								},
								Args: []ast.Expr{
									&ast.Ident{
										Name: "sigterm",
									},
									&ast.SelectorExpr{
										X: &ast.Ident{
											Name: "syscall",
										},
										Sel: &ast.Ident{
											Name: "SIGINT",
										},
									},
									&ast.SelectorExpr{
										X: &ast.Ident{
											Name: "syscall",
										},
										Sel: &ast.Ident{
											Name: "SIGTERM",
										},
									},
								},
							},
						},
						&ast.RangeStmt{
							Key: &ast.Ident{
								Name: "id",
							},
							Value: &ast.Ident{
								Name: "handler",
							},
							Tok: token.DEFINE,
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "c",
								},
								Sel: &ast.Ident{
									Name: "handlers",
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.AssignStmt{
										Lhs: []ast.Expr{
											&ast.Ident{
												Name: "consumerGroup",
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
														Name: "sarama",
													},
													Sel: &ast.Ident{
														Name: "NewConsumerGroupFromClient",
													},
												},
												Args: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "handler",
														},
														Sel: &ast.Ident{
															Name: "GroupID",
														},
													},
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "c",
														},
														Sel: &ast.Ident{
															Name: "client",
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
														&ast.CallExpr{
															Fun: &ast.SelectorExpr{
																X: &ast.CallExpr{
																	Fun: &ast.SelectorExpr{
																		X: &ast.Ident{
																			Name: "errs",
																		},
																		Sel: &ast.Ident{
																			Name: "NewUnexpectedBehaviorError",
																		},
																	},
																	Args: []ast.Expr{
																		&ast.BasicLit{
																			Kind:  token.STRING,
																			Value: "\"cant build kafka consumer\"",
																		},
																	},
																},
																Sel: &ast.Ident{
																	Name: "WithCause",
																},
															},
															Args: []ast.Expr{
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
									&ast.DeferStmt{
										Call: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "consumerGroup",
												},
												Sel: &ast.Ident{
													Name: "Close",
												},
											},
										},
									},
									&ast.AssignStmt{
										Lhs: []ast.Expr{
											&ast.SelectorExpr{
												X: &ast.Ident{
													Name: "handler",
												},
												Sel: &ast.Ident{
													Name: "group",
												},
											},
										},
										Tok: token.ASSIGN,
										Rhs: []ast.Expr{
											&ast.Ident{
												Name: "consumerGroup",
											},
										},
									},
									&ast.AssignStmt{
										Lhs: []ast.Expr{
											&ast.IndexExpr{
												X: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "c",
													},
													Sel: &ast.Ident{
														Name: "handlers",
													},
												},
												Index: &ast.Ident{
													Name: "id",
												},
											},
										},
										Tok: token.ASSIGN,
										Rhs: []ast.Expr{
											&ast.Ident{
												Name: "handler",
											},
										},
									},
								},
							},
						},
						&ast.ForStmt{
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.SelectStmt{
										Body: &ast.BlockStmt{
											List: []ast.Stmt{
												&ast.CommClause{
													Comm: &ast.ExprStmt{
														X: &ast.UnaryExpr{
															Op: token.ARROW,
															X: &ast.Ident{
																Name: "sigterm",
															},
														},
													},
													Body: []ast.Stmt{
														&ast.ReturnStmt{
															Results: []ast.Expr{
																&ast.Ident{
																	Name: "nil",
																},
															},
														},
													},
												},
												&ast.CommClause{
													Body: []ast.Stmt{
														&ast.RangeStmt{
															Key: &ast.Ident{
																Name: "_",
															},
															Value: &ast.Ident{
																Name: "handler",
															},
															Tok: token.DEFINE,
															X: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "c",
																},
																Sel: &ast.Ident{
																	Name: "handlers",
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
																								Name: "handler",
																							},
																							Sel: &ast.Ident{
																								Name: "group",
																							},
																						},
																						Sel: &ast.Ident{
																							Name: "Consume",
																						},
																					},
																					Args: []ast.Expr{
																						&ast.CallExpr{
																							Fun: ast.NewIdent("context.Background"),
																						},
																						&ast.CompositeLit{
																							Type: &ast.ArrayType{
																								Elt: &ast.Ident{
																									Name: "string",
																								},
																							},
																							Elts: []ast.Expr{
																								&ast.SelectorExpr{
																									X: &ast.Ident{
																										Name: "handler",
																									},
																									Sel: &ast.Ident{
																										Name: "Topic",
																									},
																								},
																							},
																						},
																						&ast.SelectorExpr{
																							X: &ast.Ident{
																								Name: "handler",
																							},
																							Sel: &ast.Ident{
																								Name: "Handler",
																							},
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
																								Value: "\"consume error\"",
																							},
																							&ast.CallExpr{
																								Fun: &ast.SelectorExpr{
																									X: &ast.Ident{
																										Name: "log",
																									},
																									Sel: &ast.Ident{
																										Name: "Error",
																									},
																								},
																								Args: []ast.Expr{
																									&ast.Ident{
																										Name: "err",
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
																										Value: "\"group\"",
																									},
																									&ast.SelectorExpr{
																										X: &ast.Ident{
																											Name: "handler",
																										},
																										Sel: &ast.Ident{
																											Name: "GroupID",
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
																										Value: "\"topic\"",
																									},
																									&ast.SelectorExpr{
																										X: &ast.Ident{
																											Name: "handler",
																										},
																										Sel: &ast.Ident{
																											Name: "Topic",
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
			&ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								{
									Name: "c",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "Consumer",
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
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "c",
											},
											Sel: &ast.Ident{
												Name: "client",
											},
										},
										Sel: &ast.Ident{
											Name: "Close",
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

func (u ConsumerGenerator) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "pkg", "kafka", "consumer.go")
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
