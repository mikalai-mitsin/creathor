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

func (u Middlewares) file() *ast.File {
	return &ast.File{
		Package: 1,
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
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"errors"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/pkg/errs"`, u.project.Module),
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
				},
			},
			&ast.FuncDecl{
				Name: &ast.Ident{
					Name: "defaultMessageProducer",
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
										Name: "msg",
									},
								},
								Type: &ast.Ident{
									Name: "string",
								},
							},
							{
								Names: []*ast.Ident{
									{
										Name: "level",
									},
								},
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "zapcore",
									},
									Sel: &ast.Ident{
										Name: "Level",
									},
								},
							},
							{
								Names: []*ast.Ident{
									{
										Name: "code",
									},
								},
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "codes",
									},
									Sel: &ast.Ident{
										Name: "Code",
									},
								},
							},
							{
								Names: []*ast.Ident{
									{
										Name: "err",
									},
								},
								Type: &ast.Ident{
									Name: "error",
								},
							},
							{
								Names: []*ast.Ident{
									{
										Name: "duration",
									},
								},
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "zapcore",
									},
									Sel: &ast.Ident{
										Name: "Field",
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
									Name: "logger",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "ctxzap",
										},
										Sel: &ast.Ident{
											Name: "Extract",
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
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "params",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CompositeLit{
									Type: &ast.ArrayType{
										Elt: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "zap",
											},
											Sel: &ast.Ident{
												Name: "Field",
											},
										},
									},
									Elts: []ast.Expr{
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
													Value: `"grpc.code"`,
												},
												&ast.CallExpr{
													Fun: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "code",
														},
														Sel: &ast.Ident{
															Name: "String",
														},
													},
												},
											},
										},
										&ast.Ident{
											Name: "duration",
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "zap",
												},
												Sel: &ast.Ident{
													Name: "Any",
												},
											},
											Args: []ast.Expr{
												&ast.BasicLit{
													Kind:  token.STRING,
													Value: `"request_id"`,
												},
												&ast.CallExpr{
													Fun: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "ctx",
														},
														Sel: &ast.Ident{
															Name: "Value",
														},
													},
													Args: []ast.Expr{
														&ast.SelectorExpr{
															X: &ast.Ident{
																Name: "log",
															},
															Sel: &ast.Ident{
																Name: "RequestIDKey",
															},
														},
													},
												},
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
									&ast.AssignStmt{
										Lhs: []ast.Expr{
											&ast.Ident{
												Name: "sts",
											},
										},
										Tok: token.DEFINE,
										Rhs: []ast.Expr{
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "status",
													},
													Sel: &ast.Ident{
														Name: "Convert",
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
									&ast.AssignStmt{
										Lhs: []ast.Expr{
											&ast.Ident{
												Name: "msg",
											},
										},
										Tok: token.ASSIGN,
										Rhs: []ast.Expr{
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "sts",
													},
													Sel: &ast.Ident{
														Name: "Message",
													},
												},
											},
										},
									},
									&ast.RangeStmt{
										Key: &ast.Ident{
											Name: "_",
										},
										Value: &ast.Ident{
											Name: "v",
										},
										Tok: token.DEFINE,
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "sts",
												},
												Sel: &ast.Ident{
													Name: "Details",
												},
											},
										},
										Body: &ast.BlockStmt{
											List: []ast.Stmt{
												&ast.AssignStmt{
													Lhs: []ast.Expr{
														&ast.Ident{
															Name: "errParams",
														},
													},
													Tok: token.DEFINE,
													Rhs: []ast.Expr{
														&ast.CompositeLit{
															Type: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "errs",
																},
																Sel: &ast.Ident{
																	Name: "Params",
																},
															},
														},
													},
												},
												&ast.AssignStmt{
													Lhs: []ast.Expr{
														&ast.Ident{
															Name: "badRequest",
														},
														&ast.Ident{
															Name: "ok",
														},
													},
													Tok: token.DEFINE,
													Rhs: []ast.Expr{
														&ast.TypeAssertExpr{
															X: &ast.Ident{
																Name: "v",
															},
															Type: &ast.StarExpr{
																X: &ast.SelectorExpr{
																	X: &ast.Ident{
																		Name: "errdetails",
																	},
																	Sel: &ast.Ident{
																		Name: "BadRequest",
																	},
																},
															},
														},
													},
												},
												&ast.IfStmt{
													Cond: &ast.Ident{
														Name: "ok",
													},
													Body: &ast.BlockStmt{
														List: []ast.Stmt{
															&ast.RangeStmt{
																Key: &ast.Ident{
																	Name: "_",
																},
																Value: &ast.Ident{
																	Name: "violation",
																},
																Tok: token.DEFINE,
																X: &ast.CallExpr{
																	Fun: &ast.SelectorExpr{
																		X: &ast.Ident{
																			Name: "badRequest",
																		},
																		Sel: &ast.Ident{
																			Name: "GetFieldViolations",
																		},
																	},
																},
																Body: &ast.BlockStmt{
																	List: []ast.Stmt{
																		&ast.AssignStmt{
																			Lhs: []ast.Expr{
																				&ast.Ident{
																					Name: "errParams",
																				},
																			},
																			Tok: token.ASSIGN,
																			Rhs: []ast.Expr{
																				&ast.CallExpr{
																					Fun: &ast.Ident{
																						Name: "append",
																					},
																					Args: []ast.Expr{
																						&ast.Ident{
																							Name: "errParams",
																						},
																						&ast.CompositeLit{
																							Type: &ast.SelectorExpr{
																								X: &ast.Ident{
																									Name: "errs",
																								},
																								Sel: &ast.Ident{
																									Name: "Param",
																								},
																							},
																							Elts: []ast.Expr{
																								&ast.KeyValueExpr{
																									Key: &ast.Ident{
																										Name: "Key",
																									},
																									Value: &ast.CallExpr{
																										Fun: &ast.SelectorExpr{
																											X: &ast.Ident{
																												Name: "violation",
																											},
																											Sel: &ast.Ident{
																												Name: "GetField",
																											},
																										},
																									},
																								},
																								&ast.KeyValueExpr{
																									Key: &ast.Ident{
																										Name: "Value",
																									},
																									Value: &ast.CallExpr{
																										Fun: &ast.SelectorExpr{
																											X: &ast.Ident{
																												Name: "violation",
																											},
																											Sel: &ast.Ident{
																												Name: "GetDescription",
																											},
																										},
																									},
																								},
																							},
																						},
																					},
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
															Name: "errorInfo",
														},
														&ast.Ident{
															Name: "ok",
														},
													},
													Tok: token.DEFINE,
													Rhs: []ast.Expr{
														&ast.TypeAssertExpr{
															X: &ast.Ident{
																Name: "v",
															},
															Type: &ast.StarExpr{
																X: &ast.SelectorExpr{
																	X: &ast.Ident{
																		Name: "errdetails",
																	},
																	Sel: &ast.Ident{
																		Name: "ErrorInfo",
																	},
																},
															},
														},
													},
												},
												&ast.IfStmt{
													Cond: &ast.Ident{
														Name: "ok",
													},
													Body: &ast.BlockStmt{
														List: []ast.Stmt{
															&ast.RangeStmt{
																Key: &ast.Ident{
																	Name: "key",
																},
																Value: &ast.Ident{
																	Name: "value",
																},
																Tok: token.DEFINE,
																X: &ast.CallExpr{
																	Fun: &ast.SelectorExpr{
																		X: &ast.Ident{
																			Name: "errorInfo",
																		},
																		Sel: &ast.Ident{
																			Name: "GetMetadata",
																		},
																	},
																},
																Body: &ast.BlockStmt{
																	List: []ast.Stmt{
																		&ast.AssignStmt{
																			Lhs: []ast.Expr{
																				&ast.Ident{
																					Name: "errParams",
																				},
																			},
																			Tok: token.ASSIGN,
																			Rhs: []ast.Expr{
																				&ast.CallExpr{
																					Fun: &ast.Ident{
																						Name: "append",
																					},
																					Args: []ast.Expr{
																						&ast.Ident{
																							Name: "errParams",
																						},
																						&ast.CompositeLit{
																							Type: &ast.SelectorExpr{
																								X: &ast.Ident{
																									Name: "errs",
																								},
																								Sel: &ast.Ident{
																									Name: "Param",
																								},
																							},
																							Elts: []ast.Expr{
																								&ast.KeyValueExpr{
																									Key: &ast.Ident{
																										Name: "Key",
																									},
																									Value: &ast.Ident{
																										Name: "key",
																									},
																								},
																								&ast.KeyValueExpr{
																									Key: &ast.Ident{
																										Name: "Value",
																									},
																									Value: &ast.Ident{
																										Name: "value",
																									},
																								},
																							},
																						},
																					},
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
															Name: "params",
														},
													},
													Tok: token.ASSIGN,
													Rhs: []ast.Expr{
														&ast.CallExpr{
															Fun: &ast.Ident{
																Name: "append",
															},
															Args: []ast.Expr{
																&ast.Ident{
																	Name: "params",
																},
																&ast.CallExpr{
																	Fun: &ast.SelectorExpr{
																		X: &ast.Ident{
																			Name: "zap",
																		},
																		Sel: &ast.Ident{
																			Name: "Object",
																		},
																	},
																	Args: []ast.Expr{
																		&ast.BasicLit{
																			Kind:  token.STRING,
																			Value: `"params"`,
																		},
																		&ast.Ident{
																			Name: "errParams",
																		},
																	},
																},
															},
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
											X: &ast.Ident{
												Name: "logger",
											},
											Sel: &ast.Ident{
												Name: "Check",
											},
										},
										Args: []ast.Expr{
											&ast.Ident{
												Name: "level",
											},
											&ast.Ident{
												Name: "msg",
											},
										},
									},
									Sel: &ast.Ident{
										Name: "Write",
									},
								},
								Args: []ast.Expr{
									&ast.Ident{
										Name: "params",
									},
								},
								Ellipsis: 1405,
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Name: &ast.Ident{
					Name: "unaryErrorServerUseCase",
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
								Type: &ast.InterfaceType{
									Methods: &ast.FieldList{},
								},
							},
							{
								Names: []*ast.Ident{
									{
										Name: "info",
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
								Type: &ast.InterfaceType{
									Methods: &ast.FieldList{},
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
									Name: "resp",
								},
								&ast.Ident{
									Name: "err",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.Ident{
										Name: "handler",
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "ctx",
										},
										&ast.Ident{
											Name: "req",
										},
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.Ident{
									Name: "resp",
								},
								&ast.CallExpr{
									Fun: &ast.Ident{
										Name: "handleUnaryServerError",
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "ctx",
										},
										&ast.Ident{
											Name: "req",
										},
										&ast.Ident{
											Name: "info",
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
			&ast.FuncDecl{
				Name: &ast.Ident{
					Name: "handleUnaryServerError",
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
							{
								Names: []*ast.Ident{
									{
										Name: "_",
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
										Name: "err",
									},
								},
								Type: &ast.Ident{
									Name: "error",
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
							Cond: &ast.BinaryExpr{
								X: &ast.Ident{
									Name: "err",
								},
								Op: token.EQL,
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
											{
												Name: "domainError",
											},
										},
										Type: &ast.StarExpr{
											X: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "errs",
												},
												Sel: &ast.Ident{
													Name: "Error",
												},
											},
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "errors",
									},
									Sel: &ast.Ident{
										Name: "As",
									},
								},
								Args: []ast.Expr{
									&ast.Ident{
										Name: "err",
									},
									&ast.UnaryExpr{
										Op: token.AND,
										X: &ast.Ident{
											Name: "domainError",
										},
									},
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.AssignStmt{
										Lhs: []ast.Expr{
											&ast.Ident{
												Name: "stat",
											},
										},
										Tok: token.DEFINE,
										Rhs: []ast.Expr{
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "status",
													},
													Sel: &ast.Ident{
														Name: "New",
													},
												},
												Args: []ast.Expr{
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X: &ast.Ident{
																Name: "codes",
															},
															Sel: &ast.Ident{
																Name: "Code",
															},
														},
														Args: []ast.Expr{
															&ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "domainError",
																},
																Sel: &ast.Ident{
																	Name: "Code",
																},
															},
														},
													},
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "domainError",
														},
														Sel: &ast.Ident{
															Name: "Message",
														},
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
														{
															Name: "withDetails",
														},
													},
													Type: &ast.StarExpr{
														X: &ast.SelectorExpr{
															X: &ast.Ident{
																Name: "status",
															},
															Sel: &ast.Ident{
																Name: "Status",
															},
														},
													},
												},
											},
										},
									},
									&ast.SwitchStmt{
										Tag: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "domainError",
											},
											Sel: &ast.Ident{
												Name: "Code",
											},
										},
										Body: &ast.BlockStmt{
											List: []ast.Stmt{
												&ast.CaseClause{
													List: []ast.Expr{
														&ast.SelectorExpr{
															X: &ast.Ident{
																Name: "errs",
															},
															Sel: &ast.Ident{
																Name: "ErrorCodeInvalidArgument",
															},
														},
													},
													Body: []ast.Stmt{
														&ast.AssignStmt{
															Lhs: []ast.Expr{
																&ast.Ident{
																	Name: "d",
																},
															},
															Tok: token.DEFINE,
															Rhs: []ast.Expr{
																&ast.UnaryExpr{
																	Op: token.AND,
																	X: &ast.CompositeLit{
																		Type: &ast.SelectorExpr{
																			X: &ast.Ident{
																				Name: "errdetails",
																			},
																			Sel: &ast.Ident{
																				Name: "BadRequest",
																			},
																		},
																	},
																},
															},
														},
														&ast.RangeStmt{
															Key: &ast.Ident{
																Name: "_",
															},
															Value: &ast.Ident{
																Name: "param",
															},
															Tok: token.DEFINE,
															X: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "domainError",
																},
																Sel: &ast.Ident{
																	Name: "Params",
																},
															},
															Body: &ast.BlockStmt{
																List: []ast.Stmt{
																	&ast.AssignStmt{
																		Lhs: []ast.Expr{
																			&ast.SelectorExpr{
																				X: &ast.Ident{
																					Name: "d",
																				},
																				Sel: &ast.Ident{
																					Name: "FieldViolations",
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
																							Name: "d",
																						},
																						Sel: &ast.Ident{
																							Name: "FieldViolations",
																						},
																					},
																					&ast.UnaryExpr{
																						Op: token.AND,
																						X: &ast.CompositeLit{
																							Type: &ast.SelectorExpr{
																								X: &ast.Ident{
																									Name: "errdetails",
																								},
																								Sel: &ast.Ident{
																									Name: "BadRequest_FieldViolation",
																								},
																							},
																							Elts: []ast.Expr{
																								&ast.KeyValueExpr{
																									Key: &ast.Ident{
																										Name: "Field",
																									},
																									Value: &ast.SelectorExpr{
																										X: &ast.Ident{
																											Name: "param",
																										},
																										Sel: &ast.Ident{
																											Name: "Key",
																										},
																									},
																								},
																								&ast.KeyValueExpr{
																									Key: &ast.Ident{
																										Name: "Description",
																									},
																									Value: &ast.SelectorExpr{
																										X: &ast.Ident{
																											Name: "param",
																										},
																										Sel: &ast.Ident{
																											Name: "Value",
																										},
																									},
																								},
																							},
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
																	Name: "withDetails",
																},
																&ast.Ident{
																	Name: "err",
																},
															},
															Tok: token.ASSIGN,
															Rhs: []ast.Expr{
																&ast.CallExpr{
																	Fun: &ast.SelectorExpr{
																		X: &ast.Ident{
																			Name: "stat",
																		},
																		Sel: &ast.Ident{
																			Name: "WithDetails",
																		},
																	},
																	Args: []ast.Expr{
																		&ast.Ident{
																			Name: "d",
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
																					X: &ast.Ident{
																						Name: "status",
																					},
																					Sel: &ast.Ident{
																						Name: "Error",
																					},
																				},
																				Args: []ast.Expr{
																					&ast.SelectorExpr{
																						X: &ast.Ident{
																							Name: "codes",
																						},
																						Sel: &ast.Ident{
																							Name: "Internal",
																						},
																					},
																					&ast.CallExpr{
																						Fun: &ast.SelectorExpr{
																							X: &ast.Ident{
																								Name: "err",
																							},
																							Sel: &ast.Ident{
																								Name: "Error",
																							},
																						},
																					},
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
																&ast.Ident{
																	Name: "d",
																},
															},
															Tok: token.DEFINE,
															Rhs: []ast.Expr{
																&ast.UnaryExpr{
																	Op: token.AND,
																	X: &ast.CompositeLit{
																		Type: &ast.SelectorExpr{
																			X: &ast.Ident{
																				Name: "errdetails",
																			},
																			Sel: &ast.Ident{
																				Name: "ErrorInfo",
																			},
																		},
																		Elts: []ast.Expr{
																			&ast.KeyValueExpr{
																				Key: &ast.Ident{
																					Name: "Reason",
																				},
																				Value: &ast.SelectorExpr{
																					X: &ast.Ident{
																						Name: "domainError",
																					},
																					Sel: &ast.Ident{
																						Name: "Message",
																					},
																				},
																			},
																			&ast.KeyValueExpr{
																				Key: &ast.Ident{
																					Name: "Domain",
																				},
																				Value: &ast.BasicLit{
																					Kind:  token.STRING,
																					Value: `""`,
																				},
																			},
																			&ast.KeyValueExpr{
																				Key: &ast.Ident{
																					Name: "Metadata",
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
																								Name: "string",
																							},
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
															Key: &ast.Ident{
																Name: "_",
															},
															Value: &ast.Ident{
																Name: "param",
															},
															Tok: token.DEFINE,
															X: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "domainError",
																},
																Sel: &ast.Ident{
																	Name: "Params",
																},
															},
															Body: &ast.BlockStmt{
																List: []ast.Stmt{
																	&ast.AssignStmt{
																		Lhs: []ast.Expr{
																			&ast.IndexExpr{
																				X: &ast.SelectorExpr{
																					X: &ast.Ident{
																						Name: "d",
																					},
																					Sel: &ast.Ident{
																						Name: "Metadata",
																					},
																				},
																				Index: &ast.SelectorExpr{
																					X: &ast.Ident{
																						Name: "param",
																					},
																					Sel: &ast.Ident{
																						Name: "Key",
																					},
																				},
																			},
																		},
																		Tok: token.ASSIGN,
																		Rhs: []ast.Expr{
																			&ast.SelectorExpr{
																				X: &ast.Ident{
																					Name: "param",
																				},
																				Sel: &ast.Ident{
																					Name: "Value",
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
																	Name: "withDetails",
																},
																&ast.Ident{
																	Name: "err",
																},
															},
															Tok: token.ASSIGN,
															Rhs: []ast.Expr{
																&ast.CallExpr{
																	Fun: &ast.SelectorExpr{
																		X: &ast.Ident{
																			Name: "stat",
																		},
																		Sel: &ast.Ident{
																			Name: "WithDetails",
																		},
																	},
																	Args: []ast.Expr{
																		&ast.Ident{
																			Name: "d",
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
																					X: &ast.Ident{
																						Name: "status",
																					},
																					Sel: &ast.Ident{
																						Name: "Error",
																					},
																				},
																				Args: []ast.Expr{
																					&ast.SelectorExpr{
																						X: &ast.Ident{
																							Name: "codes",
																						},
																						Sel: &ast.Ident{
																							Name: "Internal",
																						},
																					},
																					&ast.CallExpr{
																						Fun: &ast.SelectorExpr{
																							X: &ast.Ident{
																								Name: "err",
																							},
																							Sel: &ast.Ident{
																								Name: "Error",
																							},
																						},
																					},
																				},
																			},
																		},
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
														Name: "withDetails",
													},
													Sel: &ast.Ident{
														Name: "Err",
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
											Name: "status",
										},
										Sel: &ast.Ident{
											Name: "Error",
										},
									},
									Args: []ast.Expr{
										&ast.SelectorExpr{
											X: &ast.Ident{
												Name: "codes",
											},
											Sel: &ast.Ident{
												Name: "Internal",
											},
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "err",
												},
												Sel: &ast.Ident{
													Name: "Error",
												},
											},
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
					Value: fmt.Sprintf(`"%s/internal/pkg/errs"`, u.project.Module),
				},
			},
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf(`"%s/internal/pkg/log"`, u.project.Module),
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
