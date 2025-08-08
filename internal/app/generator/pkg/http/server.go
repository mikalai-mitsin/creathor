package http

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

type Server struct {
	project *configs.Project
}

func NewServer(project *configs.Project) *Server {
	return &Server{project: project}
}

func (u Server) file() *ast.File {
	return &ast.File{
		Package: 1,
		Name:    ast.NewIdent("http"),
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
						Name: ast.NewIdent("Server"),
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											ast.NewIdent("config"),
										},
										Type: &ast.StarExpr{
											X: ast.NewIdent("Config"),
										},
									},
									{
										Names: []*ast.Ident{
											ast.NewIdent("router"),
										},
										Type: &ast.StarExpr{
											X: &ast.SelectorExpr{
												X:   ast.NewIdent("chi"),
												Sel: ast.NewIdent("Mux"),
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
					},
				},
			},
			&ast.FuncDecl{
				Doc: &ast.CommentGroup{
					List: []*ast.Comment{
						{
							Text: "// NewServer - provide http server",
						},
						{
							Text: "//",
						},
						{
							Text: fmt.Sprintf("// @title %s", u.project.Name),
						},
						{
							Text: "// @host http://127.0.0.1:8000",
						},
						{
							Text: "// @BasePath /",
						},
						{
							Text: "// @version 0.0.0",
						},
						{
							Text: "// @securitydefinitions.BearerAuth BearerAuth",
						},
					},
				},
				Name: ast.NewIdent("NewServer"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("config"),
								},
								Type: &ast.StarExpr{
									X: ast.NewIdent("Config"),
								},
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.StarExpr{
									X: ast.NewIdent("Server"),
								},
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("router"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("chi"),
										Sel: ast.NewIdent("NewRouter"),
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("server"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.UnaryExpr{
									Op: token.AND,
									X: &ast.CompositeLit{
										Type: &ast.SelectorExpr{
											X:   ast.NewIdent("http"),
											Sel: ast.NewIdent("Server"),
										},
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key: ast.NewIdent("Addr"),
												Value: &ast.SelectorExpr{
													X:   ast.NewIdent("config"),
													Sel: ast.NewIdent("Address"),
												},
											},
											&ast.KeyValueExpr{
												Key:   ast.NewIdent("Handler"),
												Value: ast.NewIdent("router"),
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
										Type: ast.NewIdent("Server"),
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key:   ast.NewIdent("server"),
												Value: ast.NewIdent("server"),
											},
											&ast.KeyValueExpr{
												Key:   ast.NewIdent("config"),
												Value: ast.NewIdent("config"),
											},
											&ast.KeyValueExpr{
												Key:   ast.NewIdent("router"),
												Value: ast.NewIdent("router"),
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
								ast.NewIdent("s"),
							},
							Type: &ast.StarExpr{
								X: ast.NewIdent("Server"),
							},
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
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X:   ast.NewIdent("s"),
											Sel: ast.NewIdent("server"),
										},
										Sel: ast.NewIdent("ListenAndServe"),
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
								ast.NewIdent("s"),
							},
							Type: &ast.StarExpr{
								X: ast.NewIdent("Server"),
							},
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
										X: &ast.SelectorExpr{
											X:   ast.NewIdent("s"),
											Sel: ast.NewIdent("server"),
										},
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
								ast.NewIdent("s"),
							},
							Type: &ast.StarExpr{
								X: ast.NewIdent("Server"),
							},
						},
					},
				},
				Name: ast.NewIdent("Mount"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("path"),
								},
								Type: ast.NewIdent("string"),
							},
							{
								Names: []*ast.Ident{
									ast.NewIdent("handler"),
								},
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("http"),
									Sel: ast.NewIdent("Handler"),
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
										X:   ast.NewIdent("s"),
										Sel: ast.NewIdent("router"),
									},
									Sel: ast.NewIdent("Mount"),
								},
								Args: []ast.Expr{
									ast.NewIdent("path"),
									ast.NewIdent("handler"),
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
