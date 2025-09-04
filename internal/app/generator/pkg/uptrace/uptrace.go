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

	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
)

type Provider struct {
	project *configs.Project
}

func NewProvider(project *configs.Project) *Provider {
	return &Provider{project: project}
}

func (u Provider) file() *ast.File {
	return &ast.File{
		Name: ast.NewIdent("uptrace"),
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
							Value: u.project.ConfigsImportPath(),
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
						Name: ast.NewIdent("Provider"),
						Type: &ast.StructType{
							Fields: &ast.FieldList{
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
						},
					},
				},
			},
			&ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								ast.NewIdent("p"),
							},
							Type: ast.NewIdent("Provider"),
						},
					},
				},
				Name: ast.NewIdent("Stop"),
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
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("uptrace"),
										Sel: ast.NewIdent("Shutdown"),
									},
									Args: []ast.Expr{
										ast.NewIdent("ctx"),
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
								ast.NewIdent("p"),
							},
							Type: ast.NewIdent("Provider"),
						},
					},
				},
				Name: ast.NewIdent("Start"),
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
							Cond: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X: &ast.SelectorExpr{
										X:   ast.NewIdent("p"),
										Sel: ast.NewIdent("config"),
									},
									Sel: ast.NewIdent("Otel"),
								},
								Sel: ast.NewIdent("Enabled"),
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("uptrace"),
												Sel: ast.NewIdent("ConfigureOpentelemetry"),
											},
											Args: []ast.Expr{
												&ast.CallExpr{
													Fun: &ast.SelectorExpr{
														X:   ast.NewIdent("uptrace"),
														Sel: ast.NewIdent("WithDSN"),
													},
													Args: []ast.Expr{
														&ast.SelectorExpr{
															X: &ast.SelectorExpr{
																X: &ast.SelectorExpr{
																	X:   ast.NewIdent("p"),
																	Sel: ast.NewIdent("config"),
																},
																Sel: ast.NewIdent("Otel"),
															},
															Sel: ast.NewIdent("URL"),
														},
													},
												},
												&ast.CallExpr{
													Fun: &ast.SelectorExpr{
														X:   ast.NewIdent("uptrace"),
														Sel: ast.NewIdent("WithServiceName"),
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
														X:   ast.NewIdent("uptrace"),
														Sel: ast.NewIdent("WithServiceVersion"),
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
														X: ast.NewIdent("uptrace"),
														Sel: ast.NewIdent(
															"WithDeploymentEnvironment",
														),
													},
													Args: []ast.Expr{
														&ast.SelectorExpr{
															X: &ast.SelectorExpr{
																X: &ast.SelectorExpr{
																	X:   ast.NewIdent("p"),
																	Sel: ast.NewIdent("config"),
																},
																Sel: ast.NewIdent("Otel"),
															},
															Sel: ast.NewIdent("Environment"),
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
			&ast.FuncDecl{
				Name: ast.NewIdent("NewProvider"),
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
									X: ast.NewIdent("Provider"),
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
										Type: ast.NewIdent("Provider"),
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key:   ast.NewIdent("config"),
												Value: ast.NewIdent("config"),
											},
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
