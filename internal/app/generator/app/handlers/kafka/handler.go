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

type HandlerGenerator struct {
	domain *configs.EntityConfig
}

func NewHandlerGenerator(domain *configs.EntityConfig) *HandlerGenerator {
	return &HandlerGenerator{
		domain: domain,
	}
}

func (h *HandlerGenerator) Sync() error {
	fileset := token.NewFileSet()
	filename := h.filename()
	if err := os.MkdirAll(path.Dir(filename), 0777); err != nil {
		return err
	}
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = h.file()
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

func (h *HandlerGenerator) filename() string {
	return path.Join(
		"internal",
		"app",
		h.domain.AppConfig.AppName(),
		"handlers",
		"kafka",
		h.domain.DirName(),
		h.domain.FileName(),
	)
}

func (h *HandlerGenerator) file() *ast.File {
	return &ast.File{
		Package: 1,
		Name: &ast.Ident{
			Name: "handlers",
		},
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"github.com/IBM/sarama\"",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: h.domain.AppConfig.ProjectConfig.LogImportPath(),
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{
							Name: h.domain.KafkaHandlerTypeName(),
						},
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{

									{
										Names: []*ast.Ident{
											{
												Name: h.domain.GetUseCasePrivateVariableName(),
											},
										},
										Type: &ast.Ident{
											Name: h.domain.GetUseCaseInterfaceName(),
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "logger",
											},
										},
										Type: &ast.Ident{
											Name: "logger",
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
					Name: h.domain.KafkaHandlerConstructorName(),
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{
										Name: h.domain.GetUseCasePrivateVariableName(),
									},
								},
								Type: &ast.Ident{
									Name: h.domain.GetUseCaseInterfaceName(),
								},
							},
							{
								Names: []*ast.Ident{
									{
										Name: "logger",
									},
								},
								Type: &ast.Ident{
									Name: "logger",
								},
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.StarExpr{
									X: &ast.Ident{
										Name: h.domain.KafkaHandlerTypeName(),
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
								&ast.UnaryExpr{
									Op: token.AND,
									X: &ast.CompositeLit{
										Type: &ast.Ident{
											Name: h.domain.KafkaHandlerTypeName(),
										},
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: h.domain.GetUseCasePrivateVariableName(),
												},
												Value: &ast.Ident{
													Name: h.domain.GetUseCasePrivateVariableName(),
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
									Name: "h",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: h.domain.KafkaHandlerTypeName(),
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "Created",
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
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "h",
											},
											Sel: &ast.Ident{
												Name: "logger",
											},
										},
										Sel: &ast.Ident{
											Name: "WithContext",
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
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "logger",
									},
									Sel: &ast.Ident{
										Name: "Info",
									},
								},
								Args: []ast.Expr{
									&ast.BasicLit{
										Kind:  token.STRING,
										Value: "\"received created message\"",
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
													Name: "msg",
												},
												Sel: &ast.Ident{
													Name: "Topic",
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
												Name: "Int32",
											},
										},
										Args: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: "\"partition\"",
											},
											&ast.SelectorExpr{
												X: &ast.Ident{
													Name: "msg",
												},
												Sel: &ast.Ident{
													Name: "Partition",
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
												Value: "\"offset\"",
											},
											&ast.SelectorExpr{
												X: &ast.Ident{
													Name: "msg",
												},
												Sel: &ast.Ident{
													Name: "Offset",
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
												Value: "\"key\"",
											},
											&ast.CallExpr{
												Fun: &ast.Ident{
													Name: "string",
												},
												Args: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "msg",
														},
														Sel: &ast.Ident{
															Name: "Key",
														},
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
												Value: "\"value\"",
											},
											&ast.CallExpr{
												Fun: &ast.Ident{
													Name: "string",
												},
												Args: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "msg",
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
									Name: "h",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: h.domain.KafkaHandlerTypeName(),
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "Updated",
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
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "h",
											},
											Sel: &ast.Ident{
												Name: "logger",
											},
										},
										Sel: &ast.Ident{
											Name: "WithContext",
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
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "logger",
									},
									Sel: &ast.Ident{
										Name: "Info",
									},
								},
								Args: []ast.Expr{
									&ast.BasicLit{
										Kind:  token.STRING,
										Value: "\"received updated message\"",
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
													Name: "msg",
												},
												Sel: &ast.Ident{
													Name: "Topic",
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
												Name: "Int32",
											},
										},
										Args: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: "\"partition\"",
											},
											&ast.SelectorExpr{
												X: &ast.Ident{
													Name: "msg",
												},
												Sel: &ast.Ident{
													Name: "Partition",
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
												Value: "\"offset\"",
											},
											&ast.SelectorExpr{
												X: &ast.Ident{
													Name: "msg",
												},
												Sel: &ast.Ident{
													Name: "Offset",
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
												Value: "\"key\"",
											},
											&ast.CallExpr{
												Fun: &ast.Ident{
													Name: "string",
												},
												Args: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "msg",
														},
														Sel: &ast.Ident{
															Name: "Key",
														},
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
												Value: "\"value\"",
											},
											&ast.CallExpr{
												Fun: &ast.Ident{
													Name: "string",
												},
												Args: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "msg",
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
									Name: "h",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: h.domain.KafkaHandlerTypeName(),
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "Deleted",
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
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "h",
											},
											Sel: &ast.Ident{
												Name: "logger",
											},
										},
										Sel: &ast.Ident{
											Name: "WithContext",
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
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "logger",
									},
									Sel: &ast.Ident{
										Name: "Info",
									},
								},
								Args: []ast.Expr{
									&ast.BasicLit{
										Kind:  token.STRING,
										Value: "\"received deleted message\"",
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
													Name: "msg",
												},
												Sel: &ast.Ident{
													Name: "Topic",
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
												Name: "Int32",
											},
										},
										Args: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: "\"partition\"",
											},
											&ast.SelectorExpr{
												X: &ast.Ident{
													Name: "msg",
												},
												Sel: &ast.Ident{
													Name: "Partition",
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
												Value: "\"offset\"",
											},
											&ast.SelectorExpr{
												X: &ast.Ident{
													Name: "msg",
												},
												Sel: &ast.Ident{
													Name: "Offset",
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
												Value: "\"key\"",
											},
											&ast.CallExpr{
												Fun: &ast.Ident{
													Name: "string",
												},
												Args: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "msg",
														},
														Sel: &ast.Ident{
															Name: "Key",
														},
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
												Value: "\"value\"",
											},
											&ast.CallExpr{
												Fun: &ast.Ident{
													Name: "string",
												},
												Args: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "msg",
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
	}
}
