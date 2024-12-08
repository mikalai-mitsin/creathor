package http

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

type Server struct {
	project *configs.Project
}

func NewServer(project *configs.Project) *Server {
	return &Server{project: project}
}

func (u Server) file() *ast.File {
	return &ast.File{
		Package: 1,
		Name: &ast.Ident{
			Name: "http",
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
							Value: "\"net/http\"",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"github.com/go-chi/chi/v5\"",
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{
							Name: "Server",
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
												Name: "router",
											},
										},
										Type: &ast.StarExpr{
											X: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "chi",
												},
												Sel: &ast.Ident{
													Name: "Mux",
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
								},
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Name: &ast.Ident{
					Name: "NewServer",
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
									X: &ast.Ident{
										Name: "Config",
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
										Name: "Server",
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
									Name: "router",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "chi",
										},
										Sel: &ast.Ident{
											Name: "NewRouter",
										},
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "server",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.UnaryExpr{
									Op: token.AND,
									X: &ast.CompositeLit{
										Type: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "http",
											},
											Sel: &ast.Ident{
												Name: "Server",
											},
										},
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "Addr",
												},
												Value: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "config",
													},
													Sel: &ast.Ident{
														Name: "Address",
													},
												},
											},
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "Handler",
												},
												Value: &ast.Ident{
													Name: "router",
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
											Name: "Server",
										},
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "server",
												},
												Value: &ast.Ident{
													Name: "server",
												},
											},
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "config",
												},
												Value: &ast.Ident{
													Name: "config",
												},
											},
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "router",
												},
												Value: &ast.Ident{
													Name: "router",
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
									Name: "s",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "Server",
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
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "s",
											},
											Sel: &ast.Ident{
												Name: "server",
											},
										},
										Sel: &ast.Ident{
											Name: "ListenAndServe",
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
									Name: "s",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "Server",
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
												Name: "s",
											},
											Sel: &ast.Ident{
												Name: "server",
											},
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
									Name: "s",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "Server",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "Mount",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{
										Name: "path",
									},
								},
								Type: &ast.Ident{
									Name: "string",
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
										Name: "http",
									},
									Sel: &ast.Ident{
										Name: "Handler",
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
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "s",
										},
										Sel: &ast.Ident{
											Name: "router",
										},
									},
									Sel: &ast.Ident{
										Name: "Mount",
									},
								},
								Args: []ast.Expr{
									&ast.Ident{
										Name: "path",
									},
									&ast.Ident{
										Name: "handler",
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
					Value: "\"net/http\"",
				},
			},
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"github.com/go-chi/chi/v5\"",
				},
			},
		},
	}
}

func (u Server) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "pkg", "http", "server.go")
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
