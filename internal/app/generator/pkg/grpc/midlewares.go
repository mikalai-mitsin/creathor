package grpc

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

type Middlewares struct {
	project *configs.Project
}

func NewMiddlewares(project *configs.Project) *Middlewares {
	return &Middlewares{project: project}
}

func (u Middlewares) file() *ast.File {
	return &ast.File{
		Package: 1,
		Name:    ast.NewIdent("grpc"),
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
							Value: `"errors"`,
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
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"go.uber.org/zap"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"go.uber.org/zap/zapcore"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"google.golang.org/genproto/googleapis/rpc/errdetails"`,
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
							Value: `"go.opentelemetry.io/otel/trace"`,
						},
					},
				},
			},
			&ast.FuncDecl{
				Name: ast.NewIdent("defaultMessageProducer"),
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
							{
								Names: []*ast.Ident{
									ast.NewIdent("msg"),
								},
								Type: ast.NewIdent("string"),
							},
							{
								Names: []*ast.Ident{
									ast.NewIdent("level"),
								},
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("zapcore"),
									Sel: ast.NewIdent("Level"),
								},
							},
							{
								Names: []*ast.Ident{
									ast.NewIdent("code"),
								},
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("codes"),
									Sel: ast.NewIdent("Code"),
								},
							},
							{
								Names: []*ast.Ident{
									ast.NewIdent("err"),
								},
								Type: ast.NewIdent("error"),
							},
							{
								Names: []*ast.Ident{
									ast.NewIdent("duration"),
								},
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("zapcore"),
									Sel: ast.NewIdent("Field"),
								},
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("logger"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("ctxzap"),
										Sel: ast.NewIdent("Extract"),
									},
									Args: []ast.Expr{
										ast.NewIdent("ctx"),
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "span",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "trace",
										},
										Sel: &ast.Ident{
											Name: "SpanFromContext",
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
							Cond: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "span",
											},
											Sel: &ast.Ident{
												Name: "SpanContext",
											},
										},
									},
									Sel: &ast.Ident{
										Name: "IsValid",
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
										Tok: token.ASSIGN,
										Rhs: []ast.Expr{
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "logger",
													},
													Sel: &ast.Ident{
														Name: "With",
													},
												},
												Args: []ast.Expr{
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X: &ast.Ident{
																Name: "zap",
															},
															Sel: &ast.Ident{
																Name: "String",
															},
														},
														Args: []ast.Expr{
															&ast.BasicLit{
																Kind:  token.STRING,
																Value: "\"trace_id\"",
															},
															&ast.CallExpr{
																Fun: &ast.SelectorExpr{
																	X: &ast.CallExpr{
																		Fun: &ast.SelectorExpr{
																			X: &ast.CallExpr{
																				Fun: &ast.SelectorExpr{
																					X: &ast.Ident{
																						Name: "span",
																					},
																					Sel: &ast.Ident{
																						Name: "SpanContext",
																					},
																				},
																			},
																			Sel: &ast.Ident{
																				Name: "TraceID",
																			},
																		},
																	},
																	Sel: &ast.Ident{
																		Name: "String",
																	},
																},
															},
														},
													},
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X: &ast.Ident{
																Name: "zap",
															},
															Sel: &ast.Ident{
																Name: "String",
															},
														},
														Args: []ast.Expr{
															&ast.BasicLit{
																Kind:  token.STRING,
																Value: "\"span_id\"",
															},
															&ast.CallExpr{
																Fun: &ast.SelectorExpr{
																	X: &ast.CallExpr{
																		Fun: &ast.SelectorExpr{
																			X: &ast.CallExpr{
																				Fun: &ast.SelectorExpr{
																					X: &ast.Ident{
																						Name: "span",
																					},
																					Sel: &ast.Ident{
																						Name: "SpanContext",
																					},
																				},
																			},
																			Sel: &ast.Ident{
																				Name: "SpanID",
																			},
																		},
																	},
																	Sel: &ast.Ident{
																		Name: "String",
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
								ast.NewIdent("params"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CompositeLit{
									Type: &ast.ArrayType{
										Elt: &ast.SelectorExpr{
											X:   ast.NewIdent("zap"),
											Sel: ast.NewIdent("Field"),
										},
									},
									Elts: []ast.Expr{
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("zap"),
												Sel: ast.NewIdent("String"),
											},
											Args: []ast.Expr{
												&ast.BasicLit{
													Kind:  token.STRING,
													Value: `"grpc.code"`,
												},
												&ast.CallExpr{
													Fun: &ast.SelectorExpr{
														X:   ast.NewIdent("code"),
														Sel: ast.NewIdent("String"),
													},
												},
											},
										},
										ast.NewIdent("duration"),
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
									&ast.AssignStmt{
										Lhs: []ast.Expr{
											ast.NewIdent("sts"),
										},
										Tok: token.DEFINE,
										Rhs: []ast.Expr{
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X:   ast.NewIdent("status"),
													Sel: ast.NewIdent("Convert"),
												},
												Args: []ast.Expr{
													ast.NewIdent("err"),
												},
											},
										},
									},
									&ast.AssignStmt{
										Lhs: []ast.Expr{
											ast.NewIdent("msg"),
										},
										Tok: token.ASSIGN,
										Rhs: []ast.Expr{
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X:   ast.NewIdent("sts"),
													Sel: ast.NewIdent("Message"),
												},
											},
										},
									},
									&ast.RangeStmt{
										Key:   ast.NewIdent("_"),
										Value: ast.NewIdent("v"),
										Tok:   token.DEFINE,
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("sts"),
												Sel: ast.NewIdent("Details"),
											},
										},
										Body: &ast.BlockStmt{
											List: []ast.Stmt{
												&ast.AssignStmt{
													Lhs: []ast.Expr{
														ast.NewIdent("errParams"),
													},
													Tok: token.DEFINE,
													Rhs: []ast.Expr{
														&ast.CompositeLit{
															Type: &ast.SelectorExpr{
																X:   ast.NewIdent("errs"),
																Sel: ast.NewIdent("Params"),
															},
														},
													},
												},
												&ast.AssignStmt{
													Lhs: []ast.Expr{
														ast.NewIdent("badRequest"),
														ast.NewIdent("ok"),
													},
													Tok: token.DEFINE,
													Rhs: []ast.Expr{
														&ast.TypeAssertExpr{
															X: ast.NewIdent("v"),
															Type: &ast.StarExpr{
																X: &ast.SelectorExpr{
																	X:   ast.NewIdent("errdetails"),
																	Sel: ast.NewIdent("BadRequest"),
																},
															},
														},
													},
												},
												&ast.IfStmt{
													Cond: ast.NewIdent("ok"),
													Body: &ast.BlockStmt{
														List: []ast.Stmt{
															&ast.RangeStmt{
																Key:   ast.NewIdent("_"),
																Value: ast.NewIdent("violation"),
																Tok:   token.DEFINE,
																X: &ast.CallExpr{
																	Fun: &ast.SelectorExpr{
																		X: ast.NewIdent(
																			"badRequest",
																		),
																		Sel: ast.NewIdent(
																			"GetFieldViolations",
																		),
																	},
																},
																Body: &ast.BlockStmt{
																	List: []ast.Stmt{
																		&ast.AssignStmt{
																			Lhs: []ast.Expr{
																				ast.NewIdent(
																					"errParams",
																				),
																			},
																			Tok: token.ASSIGN,
																			Rhs: []ast.Expr{
																				&ast.CallExpr{
																					Fun: ast.NewIdent(
																						"append",
																					),
																					Args: []ast.Expr{
																						ast.NewIdent(
																							"errParams",
																						),
																						&ast.CompositeLit{
																							Type: &ast.SelectorExpr{
																								X: ast.NewIdent(
																									"errs",
																								),
																								Sel: ast.NewIdent(
																									"Param",
																								),
																							},
																							Elts: []ast.Expr{
																								&ast.KeyValueExpr{
																									Key: ast.NewIdent(
																										"Key",
																									),
																									Value: &ast.CallExpr{
																										Fun: &ast.SelectorExpr{
																											X: ast.NewIdent(
																												"violation",
																											),
																											Sel: ast.NewIdent(
																												"GetField",
																											),
																										},
																									},
																								},
																								&ast.KeyValueExpr{
																									Key: ast.NewIdent(
																										"Value",
																									),
																									Value: &ast.CallExpr{
																										Fun: &ast.SelectorExpr{
																											X: ast.NewIdent(
																												"violation",
																											),
																											Sel: ast.NewIdent(
																												"GetDescription",
																											),
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
												&ast.AssignStmt{
													Lhs: []ast.Expr{
														ast.NewIdent("errorInfo"),
														ast.NewIdent("ok"),
													},
													Tok: token.DEFINE,
													Rhs: []ast.Expr{
														&ast.TypeAssertExpr{
															X: ast.NewIdent("v"),
															Type: &ast.StarExpr{
																X: &ast.SelectorExpr{
																	X:   ast.NewIdent("errdetails"),
																	Sel: ast.NewIdent("ErrorInfo"),
																},
															},
														},
													},
												},
												&ast.IfStmt{
													Cond: ast.NewIdent("ok"),
													Body: &ast.BlockStmt{
														List: []ast.Stmt{
															&ast.RangeStmt{
																Key:   ast.NewIdent("key"),
																Value: ast.NewIdent("value"),
																Tok:   token.DEFINE,
																X: &ast.CallExpr{
																	Fun: &ast.SelectorExpr{
																		X: ast.NewIdent(
																			"errorInfo",
																		),
																		Sel: ast.NewIdent(
																			"GetMetadata",
																		),
																	},
																},
																Body: &ast.BlockStmt{
																	List: []ast.Stmt{
																		&ast.AssignStmt{
																			Lhs: []ast.Expr{
																				ast.NewIdent(
																					"errParams",
																				),
																			},
																			Tok: token.ASSIGN,
																			Rhs: []ast.Expr{
																				&ast.CallExpr{
																					Fun: ast.NewIdent(
																						"append",
																					),
																					Args: []ast.Expr{
																						ast.NewIdent(
																							"errParams",
																						),
																						&ast.CompositeLit{
																							Type: &ast.SelectorExpr{
																								X: ast.NewIdent(
																									"errs",
																								),
																								Sel: ast.NewIdent(
																									"Param",
																								),
																							},
																							Elts: []ast.Expr{
																								&ast.KeyValueExpr{
																									Key: ast.NewIdent(
																										"Key",
																									),
																									Value: ast.NewIdent(
																										"key",
																									),
																								},
																								&ast.KeyValueExpr{
																									Key: ast.NewIdent(
																										"Value",
																									),
																									Value: ast.NewIdent(
																										"value",
																									),
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
														ast.NewIdent("params"),
													},
													Tok: token.ASSIGN,
													Rhs: []ast.Expr{
														&ast.CallExpr{
															Fun: ast.NewIdent("append"),
															Args: []ast.Expr{
																ast.NewIdent("params"),
																&ast.CallExpr{
																	Fun: &ast.SelectorExpr{
																		X:   ast.NewIdent("zap"),
																		Sel: ast.NewIdent("Object"),
																	},
																	Args: []ast.Expr{
																		&ast.BasicLit{
																			Kind:  token.STRING,
																			Value: `"params"`,
																		},
																		ast.NewIdent("errParams"),
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
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   ast.NewIdent("logger"),
											Sel: ast.NewIdent("Check"),
										},
										Args: []ast.Expr{
											ast.NewIdent("level"),
											ast.NewIdent("msg"),
										},
									},
									Sel: ast.NewIdent("Write"),
								},
								Args: []ast.Expr{
									ast.NewIdent("params"),
								},
								Ellipsis: 1405,
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Name: ast.NewIdent("unaryErrorServerInterceptor"),
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
							{
								Names: []*ast.Ident{
									ast.NewIdent("req"),
								},
								Type: &ast.InterfaceType{
									Methods: &ast.FieldList{},
								},
							},
							{
								Names: []*ast.Ident{
									ast.NewIdent("info"),
								},
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X:   ast.NewIdent("grpc"),
										Sel: ast.NewIdent("UnaryServerInfo"),
									},
								},
							},
							{
								Names: []*ast.Ident{
									ast.NewIdent("handler"),
								},
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("grpc"),
									Sel: ast.NewIdent("UnaryHandler"),
								},
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.InterfaceType{
									Methods: &ast.FieldList{},
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
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("resp"),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: ast.NewIdent("handler"),
									Args: []ast.Expr{
										ast.NewIdent("ctx"),
										ast.NewIdent("req"),
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								ast.NewIdent("resp"),
								&ast.CallExpr{
									Fun: ast.NewIdent("handleUnaryServerError"),
									Args: []ast.Expr{
										ast.NewIdent("ctx"),
										ast.NewIdent("req"),
										ast.NewIdent("info"),
										ast.NewIdent("err"),
									},
								},
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Name: ast.NewIdent("handleUnaryServerError"),
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
							{
								Names: []*ast.Ident{
									ast.NewIdent("_"),
								},
								Type: ast.NewIdent("any"),
							},
							{
								Names: []*ast.Ident{
									ast.NewIdent("_"),
								},
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X:   ast.NewIdent("grpc"),
										Sel: ast.NewIdent("UnaryServerInfo"),
									},
								},
							},
							{
								Names: []*ast.Ident{
									ast.NewIdent("err"),
								},
								Type: ast.NewIdent("error"),
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
							Cond: &ast.BinaryExpr{
								X:  ast.NewIdent("err"),
								Op: token.EQL,
								Y:  ast.NewIdent("nil"),
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											ast.NewIdent("nil"),
										},
									},
								},
							},
						},
						&ast.DeclStmt{
							Decl: &ast.GenDecl{
								Tok: token.VAR,
								Specs: []ast.Spec{
									&ast.ValueSpec{
										Names: []*ast.Ident{
											ast.NewIdent("domainError"),
										},
										Type: &ast.StarExpr{
											X: &ast.SelectorExpr{
												X:   ast.NewIdent("errs"),
												Sel: ast.NewIdent("Error"),
											},
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("errors"),
									Sel: ast.NewIdent("As"),
								},
								Args: []ast.Expr{
									ast.NewIdent("err"),
									&ast.UnaryExpr{
										Op: token.AND,
										X:  ast.NewIdent("domainError"),
									},
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.AssignStmt{
										Lhs: []ast.Expr{
											ast.NewIdent("stat"),
										},
										Tok: token.DEFINE,
										Rhs: []ast.Expr{
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X:   ast.NewIdent("status"),
													Sel: ast.NewIdent("New"),
												},
												Args: []ast.Expr{
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent("codes"),
															Sel: ast.NewIdent("Code"),
														},
														Args: []ast.Expr{
															&ast.SelectorExpr{
																X:   ast.NewIdent("domainError"),
																Sel: ast.NewIdent("Code"),
															},
														},
													},
													&ast.SelectorExpr{
														X:   ast.NewIdent("domainError"),
														Sel: ast.NewIdent("Message"),
													},
												},
											},
										},
									},
									&ast.DeclStmt{
										Decl: &ast.GenDecl{
											Tok: token.VAR,
											Specs: []ast.Spec{
												&ast.ValueSpec{
													Names: []*ast.Ident{
														ast.NewIdent("withDetails"),
													},
													Type: &ast.StarExpr{
														X: &ast.SelectorExpr{
															X:   ast.NewIdent("status"),
															Sel: ast.NewIdent("Status"),
														},
													},
												},
											},
										},
									},
									&ast.SwitchStmt{
										Tag: &ast.SelectorExpr{
											X:   ast.NewIdent("domainError"),
											Sel: ast.NewIdent("Code"),
										},
										Body: &ast.BlockStmt{
											List: []ast.Stmt{
												&ast.CaseClause{
													List: []ast.Expr{
														&ast.SelectorExpr{
															X: ast.NewIdent("errs"),
															Sel: ast.NewIdent(
																"ErrorCodeInvalidArgument",
															),
														},
													},
													Body: []ast.Stmt{
														&ast.AssignStmt{
															Lhs: []ast.Expr{
																ast.NewIdent("d"),
															},
															Tok: token.DEFINE,
															Rhs: []ast.Expr{
																&ast.UnaryExpr{
																	Op: token.AND,
																	X: &ast.CompositeLit{
																		Type: &ast.SelectorExpr{
																			X: ast.NewIdent(
																				"errdetails",
																			),
																			Sel: ast.NewIdent(
																				"BadRequest",
																			),
																		},
																	},
																},
															},
														},
														&ast.RangeStmt{
															Key:   ast.NewIdent("_"),
															Value: ast.NewIdent("param"),
															Tok:   token.DEFINE,
															X: &ast.SelectorExpr{
																X:   ast.NewIdent("domainError"),
																Sel: ast.NewIdent("Params"),
															},
															Body: &ast.BlockStmt{
																List: []ast.Stmt{
																	&ast.AssignStmt{
																		Lhs: []ast.Expr{
																			&ast.SelectorExpr{
																				X: ast.NewIdent(
																					"d",
																				),
																				Sel: ast.NewIdent(
																					"FieldViolations",
																				),
																			},
																		},
																		Tok: token.ASSIGN,
																		Rhs: []ast.Expr{
																			&ast.CallExpr{
																				Fun: ast.NewIdent(
																					"append",
																				),
																				Args: []ast.Expr{
																					&ast.SelectorExpr{
																						X: ast.NewIdent(
																							"d",
																						),
																						Sel: ast.NewIdent(
																							"FieldViolations",
																						),
																					},
																					&ast.UnaryExpr{
																						Op: token.AND,
																						X: &ast.CompositeLit{
																							Type: &ast.SelectorExpr{
																								X: ast.NewIdent(
																									"errdetails",
																								),
																								Sel: ast.NewIdent(
																									"BadRequest_FieldViolation",
																								),
																							},
																							Elts: []ast.Expr{
																								&ast.KeyValueExpr{
																									Key: ast.NewIdent(
																										"Field",
																									),
																									Value: &ast.SelectorExpr{
																										X: ast.NewIdent(
																											"param",
																										),
																										Sel: ast.NewIdent(
																											"Key",
																										),
																									},
																								},
																								&ast.KeyValueExpr{
																									Key: ast.NewIdent(
																										"Description",
																									),
																									Value: &ast.SelectorExpr{
																										X: ast.NewIdent(
																											"param",
																										),
																										Sel: ast.NewIdent(
																											"Value",
																										),
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
																ast.NewIdent("withDetails"),
																ast.NewIdent("err"),
															},
															Tok: token.ASSIGN,
															Rhs: []ast.Expr{
																&ast.CallExpr{
																	Fun: &ast.SelectorExpr{
																		X: ast.NewIdent("stat"),
																		Sel: ast.NewIdent(
																			"WithDetails",
																		),
																	},
																	Args: []ast.Expr{
																		ast.NewIdent("d"),
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
																			&ast.CallExpr{
																				Fun: &ast.SelectorExpr{
																					X: ast.NewIdent(
																						"status",
																					),
																					Sel: ast.NewIdent(
																						"Error",
																					),
																				},
																				Args: []ast.Expr{
																					&ast.SelectorExpr{
																						X: ast.NewIdent(
																							"codes",
																						),
																						Sel: ast.NewIdent(
																							"Internal",
																						),
																					},
																					&ast.CallExpr{
																						Fun: &ast.SelectorExpr{
																							X: ast.NewIdent(
																								"err",
																							),
																							Sel: ast.NewIdent(
																								"Error",
																							),
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
												&ast.CaseClause{
													Body: []ast.Stmt{
														&ast.AssignStmt{
															Lhs: []ast.Expr{
																ast.NewIdent("d"),
															},
															Tok: token.DEFINE,
															Rhs: []ast.Expr{
																&ast.UnaryExpr{
																	Op: token.AND,
																	X: &ast.CompositeLit{
																		Type: &ast.SelectorExpr{
																			X: ast.NewIdent(
																				"errdetails",
																			),
																			Sel: ast.NewIdent(
																				"ErrorInfo",
																			),
																		},
																		Elts: []ast.Expr{
																			&ast.KeyValueExpr{
																				Key: ast.NewIdent(
																					"Reason",
																				),
																				Value: &ast.SelectorExpr{
																					X: ast.NewIdent(
																						"domainError",
																					),
																					Sel: ast.NewIdent(
																						"Message",
																					),
																				},
																			},
																			&ast.KeyValueExpr{
																				Key: ast.NewIdent(
																					"Domain",
																				),
																				Value: &ast.BasicLit{
																					Kind:  token.STRING,
																					Value: `""`,
																				},
																			},
																			&ast.KeyValueExpr{
																				Key: ast.NewIdent(
																					"Metadata",
																				),
																				Value: &ast.CallExpr{
																					Fun: ast.NewIdent(
																						"make",
																					),
																					Args: []ast.Expr{
																						&ast.MapType{
																							Key: ast.NewIdent(
																								"string",
																							),
																							Value: ast.NewIdent(
																								"string",
																							),
																						},
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
														&ast.RangeStmt{
															Key:   ast.NewIdent("_"),
															Value: ast.NewIdent("param"),
															Tok:   token.DEFINE,
															X: &ast.SelectorExpr{
																X:   ast.NewIdent("domainError"),
																Sel: ast.NewIdent("Params"),
															},
															Body: &ast.BlockStmt{
																List: []ast.Stmt{
																	&ast.AssignStmt{
																		Lhs: []ast.Expr{
																			&ast.IndexExpr{
																				X: &ast.SelectorExpr{
																					X: ast.NewIdent(
																						"d",
																					),
																					Sel: ast.NewIdent(
																						"Metadata",
																					),
																				},
																				Index: &ast.SelectorExpr{
																					X: ast.NewIdent(
																						"param",
																					),
																					Sel: ast.NewIdent(
																						"Key",
																					),
																				},
																			},
																		},
																		Tok: token.ASSIGN,
																		Rhs: []ast.Expr{
																			&ast.SelectorExpr{
																				X: ast.NewIdent(
																					"param",
																				),
																				Sel: ast.NewIdent(
																					"Value",
																				),
																			},
																		},
																	},
																},
															},
														},
														&ast.AssignStmt{
															Lhs: []ast.Expr{
																ast.NewIdent("withDetails"),
																ast.NewIdent("err"),
															},
															Tok: token.ASSIGN,
															Rhs: []ast.Expr{
																&ast.CallExpr{
																	Fun: &ast.SelectorExpr{
																		X: ast.NewIdent("stat"),
																		Sel: ast.NewIdent(
																			"WithDetails",
																		),
																	},
																	Args: []ast.Expr{
																		ast.NewIdent("d"),
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
																			&ast.CallExpr{
																				Fun: &ast.SelectorExpr{
																					X: ast.NewIdent(
																						"status",
																					),
																					Sel: ast.NewIdent(
																						"Error",
																					),
																				},
																				Args: []ast.Expr{
																					&ast.SelectorExpr{
																						X: ast.NewIdent(
																							"codes",
																						),
																						Sel: ast.NewIdent(
																							"Internal",
																						),
																					},
																					&ast.CallExpr{
																						Fun: &ast.SelectorExpr{
																							X: ast.NewIdent(
																								"err",
																							),
																							Sel: ast.NewIdent(
																								"Error",
																							),
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
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X:   ast.NewIdent("withDetails"),
													Sel: ast.NewIdent("Err"),
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
										X:   ast.NewIdent("status"),
										Sel: ast.NewIdent("Error"),
									},
									Args: []ast.Expr{
										&ast.SelectorExpr{
											X:   ast.NewIdent("codes"),
											Sel: ast.NewIdent("Internal"),
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("err"),
												Sel: ast.NewIdent("Error"),
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
		Imports: []*ast.ImportSpec{
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"context"`,
				},
			},
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"errors"`,
				},
			},
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: u.project.ErrsImportPath(),
				},
			},
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: u.project.LogImportPath(),
				},
			},
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"`,
				},
			},
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"go.uber.org/zap"`,
				},
			},
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"go.uber.org/zap/zapcore"`,
				},
			},
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"google.golang.org/genproto/googleapis/rpc/errdetails"`,
				},
			},
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"google.golang.org/grpc"`,
				},
			},
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"google.golang.org/grpc/codes"`,
				},
			},
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"google.golang.org/grpc/status"`,
				},
			},
		},
	}
}

func (u Middlewares) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "pkg", "grpc", "middlewares.go")
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
