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

type TraceGenerator struct {
	project *configs.Project
}

func NewTraceGenerator(project *configs.Project) *TraceGenerator {
	return &TraceGenerator{project: project}
}

func (u TraceGenerator) file() *ast.File {
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
							Value: "\"github.com/IBM/sarama\"",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"go.opentelemetry.io/otel\"",
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{
							Name: "consumerHeadersCarrier",
						},
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											{
												Name: "headers",
											},
										},
										Type: &ast.ArrayType{
											Elt: &ast.StarExpr{
												X: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "sarama",
													},
													Sel: &ast.Ident{
														Name: "RecordHeader",
													},
												},
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
							Type: &ast.Ident{
								Name: "consumerHeadersCarrier",
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "Get",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{
										Name: "key",
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
								Type: &ast.Ident{
									Name: "string",
								},
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.RangeStmt{
							Key: &ast.Ident{
								Name: "_",
							},
							Value: &ast.Ident{
								Name: "h",
							},
							Tok: token.DEFINE,
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "c",
								},
								Sel: &ast.Ident{
									Name: "headers",
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.IfStmt{
										Cond: &ast.BinaryExpr{
											X: &ast.CallExpr{
												Fun: &ast.Ident{
													Name: "string",
												},
												Args: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "h",
														},
														Sel: &ast.Ident{
															Name: "Key",
														},
													},
												},
											},
											Op: token.EQL,
											Y: &ast.Ident{
												Name: "key",
											},
										},
										Body: &ast.BlockStmt{
											List: []ast.Stmt{
												&ast.ReturnStmt{
													Results: []ast.Expr{
														&ast.CallExpr{
															Fun: &ast.Ident{
																Name: "string",
															},
															Args: []ast.Expr{
																&ast.SelectorExpr{
																	X: &ast.Ident{
																		Name: "h",
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
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: "\"\"",
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
							Type: &ast.Ident{
								Name: "consumerHeadersCarrier",
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "Set",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{
										Name: "key",
									},
									{
										Name: "val",
									},
								},
								Type: &ast.Ident{
									Name: "string",
								},
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.RangeStmt{
							Key: &ast.Ident{
								Name: "_",
							},
							Value: &ast.Ident{
								Name: "h",
							},
							Tok: token.DEFINE,
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "c",
								},
								Sel: &ast.Ident{
									Name: "headers",
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.IfStmt{
										Cond: &ast.BinaryExpr{
											X: &ast.CallExpr{
												Fun: &ast.Ident{
													Name: "string",
												},
												Args: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "h",
														},
														Sel: &ast.Ident{
															Name: "Key",
														},
													},
												},
											},
											Op: token.EQL,
											Y: &ast.Ident{
												Name: "key",
											},
										},
										Body: &ast.BlockStmt{
											List: []ast.Stmt{
												&ast.AssignStmt{
													Lhs: []ast.Expr{
														&ast.SelectorExpr{
															X: &ast.Ident{
																Name: "h",
															},
															Sel: &ast.Ident{
																Name: "Value",
															},
														},
													},
													Tok: token.ASSIGN,
													Rhs: []ast.Expr{
														&ast.CallExpr{
															Fun: &ast.ArrayType{
																Elt: &ast.Ident{
																	Name: "byte",
																},
															},
															Args: []ast.Expr{
																&ast.Ident{
																	Name: "val",
																},
															},
														},
													},
												},
												&ast.ReturnStmt{},
											},
										},
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.SelectorExpr{
									X: &ast.Ident{
										Name: "c",
									},
									Sel: &ast.Ident{
										Name: "headers",
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
												Name: "c",
											},
											Sel: &ast.Ident{
												Name: "headers",
											},
										},
										&ast.UnaryExpr{
											Op: token.AND,
											X: &ast.CompositeLit{
												Type: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "sarama",
													},
													Sel: &ast.Ident{
														Name: "RecordHeader",
													},
												},
												Elts: []ast.Expr{
													&ast.KeyValueExpr{
														Key: &ast.Ident{
															Name: "Key",
														},
														Value: &ast.CallExpr{
															Fun: &ast.ArrayType{
																Elt: &ast.Ident{
																	Name: "byte",
																},
															},
															Args: []ast.Expr{
																&ast.Ident{
																	Name: "key",
																},
															},
														},
													},
													&ast.KeyValueExpr{
														Key: &ast.Ident{
															Name: "Value",
														},
														Value: &ast.CallExpr{
															Fun: &ast.ArrayType{
																Elt: &ast.Ident{
																	Name: "byte",
																},
															},
															Args: []ast.Expr{
																&ast.Ident{
																	Name: "val",
																},
															},
														},
													},
												},
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
							Type: &ast.Ident{
								Name: "consumerHeadersCarrier",
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "Keys",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.ArrayType{
									Elt: &ast.Ident{
										Name: "string",
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
									Name: "keys",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.Ident{
										Name: "make",
									},
									Args: []ast.Expr{
										&ast.ArrayType{
											Elt: &ast.Ident{
												Name: "string",
											},
										},
										&ast.BasicLit{
											Kind:  token.INT,
											Value: "0",
										},
										&ast.CallExpr{
											Fun: &ast.Ident{
												Name: "len",
											},
											Args: []ast.Expr{
												&ast.SelectorExpr{
													X: &ast.Ident{
														Name: "c",
													},
													Sel: &ast.Ident{
														Name: "headers",
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
								Name: "h",
							},
							Tok: token.DEFINE,
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "c",
								},
								Sel: &ast.Ident{
									Name: "headers",
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.AssignStmt{
										Lhs: []ast.Expr{
											&ast.Ident{
												Name: "keys",
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
														Name: "keys",
													},
													&ast.CallExpr{
														Fun: &ast.Ident{
															Name: "string",
														},
														Args: []ast.Expr{
															&ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "h",
																},
																Sel: &ast.Ident{
																	Name: "Key",
																},
															},
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
									Name: "keys",
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
							Name: "producerHeadersCarrier",
						},
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											{
												Name: "headers",
											},
										},
										Type: &ast.ArrayType{
											Elt: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "sarama",
												},
												Sel: &ast.Ident{
													Name: "RecordHeader",
												},
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
							Type: &ast.Ident{
								Name: "producerHeadersCarrier",
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "Get",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{
										Name: "key",
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
								Type: &ast.Ident{
									Name: "string",
								},
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.RangeStmt{
							Key: &ast.Ident{
								Name: "_",
							},
							Value: &ast.Ident{
								Name: "h",
							},
							Tok: token.DEFINE,
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "c",
								},
								Sel: &ast.Ident{
									Name: "headers",
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.IfStmt{
										Cond: &ast.BinaryExpr{
											X: &ast.CallExpr{
												Fun: &ast.Ident{
													Name: "string",
												},
												Args: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "h",
														},
														Sel: &ast.Ident{
															Name: "Key",
														},
													},
												},
											},
											Op: token.EQL,
											Y: &ast.Ident{
												Name: "key",
											},
										},
										Body: &ast.BlockStmt{
											List: []ast.Stmt{
												&ast.ReturnStmt{
													Results: []ast.Expr{
														&ast.CallExpr{
															Fun: &ast.Ident{
																Name: "string",
															},
															Args: []ast.Expr{
																&ast.SelectorExpr{
																	X: &ast.Ident{
																		Name: "h",
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
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: "\"\"",
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
									Name: "producerHeadersCarrier",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "Set",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{
										Name: "key",
									},
									{
										Name: "val",
									},
								},
								Type: &ast.Ident{
									Name: "string",
								},
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.RangeStmt{
							Key: &ast.Ident{
								Name: "i",
							},
							Tok: token.DEFINE,
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "c",
								},
								Sel: &ast.Ident{
									Name: "headers",
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.IfStmt{
										Cond: &ast.BinaryExpr{
											X: &ast.CallExpr{
												Fun: &ast.Ident{
													Name: "string",
												},
												Args: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.IndexExpr{
															X: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "c",
																},
																Sel: &ast.Ident{
																	Name: "headers",
																},
															},
															Index: &ast.Ident{
																Name: "i",
															},
														},
														Sel: &ast.Ident{
															Name: "Key",
														},
													},
												},
											},
											Op: token.EQL,
											Y: &ast.Ident{
												Name: "key",
											},
										},
										Body: &ast.BlockStmt{
											List: []ast.Stmt{
												&ast.AssignStmt{
													Lhs: []ast.Expr{
														&ast.SelectorExpr{
															X: &ast.IndexExpr{
																X: &ast.SelectorExpr{
																	X: &ast.Ident{
																		Name: "c",
																	},
																	Sel: &ast.Ident{
																		Name: "headers",
																	},
																},
																Index: &ast.Ident{
																	Name: "i",
																},
															},
															Sel: &ast.Ident{
																Name: "Value",
															},
														},
													},
													Tok: token.ASSIGN,
													Rhs: []ast.Expr{
														&ast.CallExpr{
															Fun: &ast.ArrayType{
																Elt: &ast.Ident{
																	Name: "byte",
																},
															},
															Args: []ast.Expr{
																&ast.Ident{
																	Name: "val",
																},
															},
														},
													},
												},
												&ast.ReturnStmt{},
											},
										},
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.SelectorExpr{
									X: &ast.Ident{
										Name: "c",
									},
									Sel: &ast.Ident{
										Name: "headers",
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
												Name: "c",
											},
											Sel: &ast.Ident{
												Name: "headers",
											},
										},
										&ast.CompositeLit{
											Type: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "sarama",
												},
												Sel: &ast.Ident{
													Name: "RecordHeader",
												},
											},
											Elts: []ast.Expr{
												&ast.KeyValueExpr{
													Key: &ast.Ident{
														Name: "Key",
													},
													Value: &ast.CallExpr{
														Fun: &ast.ArrayType{
															Elt: &ast.Ident{
																Name: "byte",
															},
														},
														Args: []ast.Expr{
															&ast.Ident{
																Name: "key",
															},
														},
													},
												},
												&ast.KeyValueExpr{
													Key: &ast.Ident{
														Name: "Value",
													},
													Value: &ast.CallExpr{
														Fun: &ast.ArrayType{
															Elt: &ast.Ident{
																Name: "byte",
															},
														},
														Args: []ast.Expr{
															&ast.Ident{
																Name: "val",
															},
														},
													},
												},
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
							Type: &ast.Ident{
								Name: "producerHeadersCarrier",
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "Keys",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.ArrayType{
									Elt: &ast.Ident{
										Name: "string",
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
									Name: "keys",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.Ident{
										Name: "make",
									},
									Args: []ast.Expr{
										&ast.ArrayType{
											Elt: &ast.Ident{
												Name: "string",
											},
										},
										&ast.BasicLit{
											Kind:  token.INT,
											Value: "0",
										},
										&ast.CallExpr{
											Fun: &ast.Ident{
												Name: "len",
											},
											Args: []ast.Expr{
												&ast.SelectorExpr{
													X: &ast.Ident{
														Name: "c",
													},
													Sel: &ast.Ident{
														Name: "headers",
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
								Name: "h",
							},
							Tok: token.DEFINE,
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "c",
								},
								Sel: &ast.Ident{
									Name: "headers",
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.AssignStmt{
										Lhs: []ast.Expr{
											&ast.Ident{
												Name: "keys",
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
														Name: "keys",
													},
													&ast.CallExpr{
														Fun: &ast.Ident{
															Name: "string",
														},
														Args: []ast.Expr{
															&ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "h",
																},
																Sel: &ast.Ident{
																	Name: "Key",
																},
															},
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
									Name: "keys",
								},
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Name: &ast.Ident{
					Name: "ExtractTraceContext",
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
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "sarama",
										},
										Sel: &ast.Ident{
											Name: "ConsumerMessage",
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
										Name: "context",
									},
									Sel: &ast.Ident{
										Name: "Context",
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
									Name: "propagator",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "otel",
										},
										Sel: &ast.Ident{
											Name: "GetTextMapPropagator",
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
											Name: "propagator",
										},
										Sel: &ast.Ident{
											Name: "Extract",
										},
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "ctx",
										},
										&ast.CompositeLit{
											Type: &ast.Ident{
												Name: "consumerHeadersCarrier",
											},
											Elts: []ast.Expr{
												&ast.KeyValueExpr{
													Key: &ast.Ident{
														Name: "headers",
													},
													Value: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "msg",
														},
														Sel: &ast.Ident{
															Name: "Headers",
														},
													},
												},
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
					Name: "InjectTraceContext",
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
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "sarama",
										},
										Sel: &ast.Ident{
											Name: "ProducerMessage",
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
									Name: "propagator",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "otel",
										},
										Sel: &ast.Ident{
											Name: "GetTextMapPropagator",
										},
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "carrier",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.UnaryExpr{
									Op: token.AND,
									X: &ast.CompositeLit{
										Type: &ast.Ident{
											Name: "producerHeadersCarrier",
										},
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "headers",
												},
												Value: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "msg",
													},
													Sel: &ast.Ident{
														Name: "Headers",
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
										Name: "propagator",
									},
									Sel: &ast.Ident{
										Name: "Inject",
									},
								},
								Args: []ast.Expr{
									&ast.Ident{
										Name: "ctx",
									},
									&ast.Ident{
										Name: "carrier",
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.SelectorExpr{
									X: &ast.Ident{
										Name: "msg",
									},
									Sel: &ast.Ident{
										Name: "Headers",
									},
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.SelectorExpr{
									X: &ast.Ident{
										Name: "carrier",
									},
									Sel: &ast.Ident{
										Name: "headers",
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
					Value: "\"context\"",
				},
			},
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"github.com/IBM/sarama\"",
				},
			},
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"go.opentelemetry.io/otel\"",
				},
			},
		},
	}
}

func (u TraceGenerator) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "pkg", "kafka", "trace.go")
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
