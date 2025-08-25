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

type ProducerGenerator struct {
	project *configs.Project
}

func NewProducerGenerator(project *configs.Project) *ProducerGenerator {
	return &ProducerGenerator{project: project}
}

func (u ProducerGenerator) file() *ast.File {
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
							Name: "Message",
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
												Name: "Value",
											},
										},
										Type: &ast.ArrayType{
											Elt: &ast.Ident{
												Name: "byte",
											},
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "Key",
											},
										},
										Type: &ast.Ident{
											Name: "string",
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
							Name: "Producer",
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
												Name: "producer",
											},
										},
										Type: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "sarama",
											},
											Sel: &ast.Ident{
												Name: "SyncProducer",
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
					Name: "NewProducer",
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
										Name: "Producer",
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
									X: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "config",
											},
											Sel: &ast.Ident{
												Name: "Producer",
											},
										},
										Sel: &ast.Ident{
											Name: "Return",
										},
									},
									Sel: &ast.Ident{
										Name: "Successes",
									},
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.Ident{
									Name: "true",
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.SelectorExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "config",
										},
										Sel: &ast.Ident{
											Name: "Producer",
										},
									},
									Sel: &ast.Ident{
										Name: "RequiredAcks",
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
										Name: "WaitForAll",
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.SelectorExpr{
									X: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "config",
											},
											Sel: &ast.Ident{
												Name: "Producer",
											},
										},
										Sel: &ast.Ident{
											Name: "Retry",
										},
									},
									Sel: &ast.Ident{
										Name: "Max",
									},
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.INT,
									Value: "5",
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
									Name: "producer",
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
											Name: "NewSyncProducer",
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
																Value: "\"cant build kafka producer\"",
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
											Name: "Producer",
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
													Name: "producer",
												},
												Value: &ast.Ident{
													Name: "producer",
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
									Name: "p",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "Producer",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "Send",
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
										Name: "message",
									},
								},
								Type: &ast.StarExpr{
									X: &ast.Ident{
										Name: "Message",
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
									Name: "msg",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.UnaryExpr{
									Op: token.AND,
									X: &ast.CompositeLit{
										Type: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "sarama",
											},
											Sel: &ast.Ident{
												Name: "ProducerMessage",
											},
										},
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "Topic",
												},
												Value: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "message",
													},
													Sel: &ast.Ident{
														Name: "Topic",
													},
												},
											},
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "Key",
												},
												Value: &ast.CallExpr{
													Fun: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "sarama",
														},
														Sel: &ast.Ident{
															Name: "StringEncoder",
														},
													},
													Args: []ast.Expr{
														&ast.SelectorExpr{
															X: &ast.Ident{
																Name: "message",
															},
															Sel: &ast.Ident{
																Name: "Key",
															},
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
															Name: "sarama",
														},
														Sel: &ast.Ident{
															Name: "ByteEncoder",
														},
													},
													Args: []ast.Expr{
														&ast.SelectorExpr{
															X: &ast.Ident{
																Name: "message",
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
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "_",
								},
								&ast.Ident{
									Name: "_",
								},
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
												Name: "p",
											},
											Sel: &ast.Ident{
												Name: "producer",
											},
										},
										Sel: &ast.Ident{
											Name: "SendMessage",
										},
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "msg",
										},
									},
								},
							},
						},
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
			&ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								{
									Name: "p",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "Producer",
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
									Name: "p",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "Producer",
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
												Name: "p",
											},
											Sel: &ast.Ident{
												Name: "producer",
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

func (u ProducerGenerator) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "pkg", "kafka", "producer.go")
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
