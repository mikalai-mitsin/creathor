package uptrace

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"

	"github.com/018bf/creathor/internal/configs"
)

type Provider struct {
	project *configs.Project
}

func NewProvider(project *configs.Project) *Provider {
	return &Provider{project: project}
}

func (u Provider) file() *ast.File {
	return &ast.File{
		Name: &ast.Ident{
			Name: "uptrace",
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
							Value: fmt.Sprintf(`"%s"`, u.project.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/pkg/configs"`, u.project.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"github.com/uptrace/uptrace-go/uptrace"`,
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{
							Name: "Provider",
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
							Type: &ast.Ident{
								Name: "Provider",
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
										X: &ast.Ident{
											Name: "uptrace",
										},
										Sel: &ast.Ident{
											Name: "Shutdown",
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
							Type: &ast.Ident{
								Name: "Provider",
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
							Cond: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "p",
										},
										Sel: &ast.Ident{
											Name: "config",
										},
									},
									Sel: &ast.Ident{
										Name: "Otel",
									},
								},
								Sel: &ast.Ident{
									Name: "Enabled",
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "uptrace",
												},
												Sel: &ast.Ident{
													Name: "ConfigureOpentelemetry",
												},
											},
											Args: []ast.Expr{
												&ast.CallExpr{
													Fun: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "uptrace",
														},
														Sel: &ast.Ident{
															Name: "WithDSN",
														},
													},
													Args: []ast.Expr{
														&ast.SelectorExpr{
															X: &ast.SelectorExpr{
																X: &ast.SelectorExpr{
																	X: &ast.Ident{
																		Name: "p",
																	},
																	Sel: &ast.Ident{
																		Name: "config",
																	},
																},
																Sel: &ast.Ident{
																	Name: "Otel",
																},
															},
															Sel: &ast.Ident{
																Name: "URL",
															},
														},
													},
												},
												&ast.CallExpr{
													Fun: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "uptrace",
														},
														Sel: &ast.Ident{
															Name: "WithServiceName",
														},
													},
													Args: []ast.Expr{
														&ast.SelectorExpr{
															X:   ast.NewIdent(u.project.Name),
															Sel: ast.NewIdent("Name"),
														},
													},
												},
												&ast.CallExpr{
													Fun: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "uptrace",
														},
														Sel: &ast.Ident{
															Name: "WithServiceVersion",
														},
													},
													Args: []ast.Expr{
														&ast.SelectorExpr{
															X:   ast.NewIdent(u.project.Name),
															Sel: ast.NewIdent("Version"),
														},
													},
												},
												&ast.CallExpr{
													Fun: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "uptrace",
														},
														Sel: &ast.Ident{
															Name: "WithDeploymentEnvironment",
														},
													},
													Args: []ast.Expr{
														&ast.SelectorExpr{
															X: &ast.SelectorExpr{
																X: &ast.SelectorExpr{
																	X: &ast.Ident{
																		Name: "p",
																	},
																	Sel: &ast.Ident{
																		Name: "config",
																	},
																},
																Sel: &ast.Ident{
																	Name: "Otel",
																},
															},
															Sel: &ast.Ident{
																Name: "Environment",
															},
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
				Name: &ast.Ident{
					Name: "NewProvider",
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
									X: &ast.Ident{
										Name: "Provider",
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
											Name: "Provider",
										},
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "config",
												},
												Value: &ast.Ident{
													Name: "config",
												},
											},
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

func (u Provider) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "pkg", "uptrace", "uptrace.go")
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
