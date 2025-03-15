package http

import (
	"bytes"
	"fmt"
	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"
)

const destinationPath = "."

type HandlerGenerator struct {
	project *configs.Project
}

func NewHandler(project *configs.Project) *HandlerGenerator {
	return &HandlerGenerator{
		project: project,
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
	return path.Join("internal", "app", "auth", "handlers", "http", "auth.go")
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
							Value: "\"net/http\"",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"github.com/go-chi/chi/v5\"",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"github.com/go-chi/render\"",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/pkg/errs"`, h.project.Module),
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{
							Name: "AuthHandler",
						},
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									&ast.Field{
										Names: []*ast.Ident{
											&ast.Ident{
												Name: "authUseCase",
											},
										},
										Type: &ast.Ident{
											Name: "authUseCase",
										},
									},
									&ast.Field{
										Names: []*ast.Ident{
											&ast.Ident{
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
					Name: "NewAuthHandler",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Names: []*ast.Ident{
									&ast.Ident{
										Name: "authUseCase",
									},
								},
								Type: &ast.Ident{
									Name: "authUseCase",
								},
							},
							&ast.Field{
								Names: []*ast.Ident{
									&ast.Ident{
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
							&ast.Field{
								Type: &ast.StarExpr{
									X: &ast.Ident{
										Name: "AuthHandler",
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
											Name: "AuthHandler",
										},
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "authUseCase",
												},
												Value: &ast.Ident{
													Name: "authUseCase",
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
						&ast.Field{
							Names: []*ast.Ident{
								&ast.Ident{
									Name: "h",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "AuthHandler",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "ObtainTokenPair",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Names: []*ast.Ident{
									&ast.Ident{
										Name: "w",
									},
								},
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "http",
									},
									Sel: &ast.Ident{
										Name: "ResponseWriter",
									},
								},
							},
							&ast.Field{
								Names: []*ast.Ident{
									&ast.Ident{
										Name: "r",
									},
								},
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "http",
										},
										Sel: &ast.Ident{
											Name: "Request",
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
									Name: "createDTO",
								},
								&ast.Ident{
									Name: "err",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.Ident{
										Name: "NewObtainTokenDTO",
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "r",
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
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "errs",
												},
												Sel: &ast.Ident{
													Name: "RenderToHTTPResponse",
												},
											},
											Args: []ast.Expr{
												&ast.Ident{
													Name: "err",
												},
												&ast.Ident{
													Name: "w",
												},
												&ast.Ident{
													Name: "r",
												},
											},
										},
									},
									&ast.ReturnStmt{},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "create",
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
											Name: "createDTO",
										},
										Sel: &ast.Ident{
											Name: "toEntity",
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
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "errs",
												},
												Sel: &ast.Ident{
													Name: "RenderToHTTPResponse",
												},
											},
											Args: []ast.Expr{
												&ast.Ident{
													Name: "err",
												},
												&ast.Ident{
													Name: "w",
												},
												&ast.Ident{
													Name: "r",
												},
											},
										},
									},
									&ast.ReturnStmt{},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "tokenPair",
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
												Name: "h",
											},
											Sel: &ast.Ident{
												Name: "authUseCase",
											},
										},
										Sel: &ast.Ident{
											Name: "CreateToken",
										},
									},
									Args: []ast.Expr{
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "r",
												},
												Sel: &ast.Ident{
													Name: "Context",
												},
											},
										},
										&ast.Ident{
											Name: "create",
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
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "errs",
												},
												Sel: &ast.Ident{
													Name: "RenderToHTTPResponse",
												},
											},
											Args: []ast.Expr{
												&ast.Ident{
													Name: "err",
												},
												&ast.Ident{
													Name: "w",
												},
												&ast.Ident{
													Name: "r",
												},
											},
										},
									},
									&ast.ReturnStmt{},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "response",
								},
								&ast.Ident{
									Name: "err",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.Ident{
										Name: "NewTokenPairDTO",
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "tokenPair",
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
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "errs",
												},
												Sel: &ast.Ident{
													Name: "RenderToHTTPResponse",
												},
											},
											Args: []ast.Expr{
												&ast.Ident{
													Name: "err",
												},
												&ast.Ident{
													Name: "w",
												},
												&ast.Ident{
													Name: "r",
												},
											},
										},
									},
									&ast.ReturnStmt{},
								},
							},
						},
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "render",
									},
									Sel: &ast.Ident{
										Name: "Status",
									},
								},
								Args: []ast.Expr{
									&ast.Ident{
										Name: "r",
									},
									&ast.SelectorExpr{
										X: &ast.Ident{
											Name: "http",
										},
										Sel: &ast.Ident{
											Name: "StatusOK",
										},
									},
								},
							},
						},
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "render",
									},
									Sel: &ast.Ident{
										Name: "JSON",
									},
								},
								Args: []ast.Expr{
									&ast.Ident{
										Name: "w",
									},
									&ast.Ident{
										Name: "r",
									},
									&ast.Ident{
										Name: "response",
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
						&ast.Field{
							Names: []*ast.Ident{
								&ast.Ident{
									Name: "h",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "AuthHandler",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "RefreshTokenPair",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Names: []*ast.Ident{
									&ast.Ident{
										Name: "w",
									},
								},
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "http",
									},
									Sel: &ast.Ident{
										Name: "ResponseWriter",
									},
								},
							},
							&ast.Field{
								Names: []*ast.Ident{
									&ast.Ident{
										Name: "r",
									},
								},
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "http",
										},
										Sel: &ast.Ident{
											Name: "Request",
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
									Name: "refreshTokenDTO",
								},
								&ast.Ident{
									Name: "err",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.Ident{
										Name: "NewRefreshTokenDTO",
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "r",
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
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "errs",
												},
												Sel: &ast.Ident{
													Name: "RenderToHTTPResponse",
												},
											},
											Args: []ast.Expr{
												&ast.Ident{
													Name: "err",
												},
												&ast.Ident{
													Name: "w",
												},
												&ast.Ident{
													Name: "r",
												},
											},
										},
									},
									&ast.ReturnStmt{},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "refreshToken",
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
											Name: "refreshTokenDTO",
										},
										Sel: &ast.Ident{
											Name: "toEntity",
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
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "errs",
												},
												Sel: &ast.Ident{
													Name: "RenderToHTTPResponse",
												},
											},
											Args: []ast.Expr{
												&ast.Ident{
													Name: "err",
												},
												&ast.Ident{
													Name: "w",
												},
												&ast.Ident{
													Name: "r",
												},
											},
										},
									},
									&ast.ReturnStmt{},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "tokenPair",
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
												Name: "h",
											},
											Sel: &ast.Ident{
												Name: "authUseCase",
											},
										},
										Sel: &ast.Ident{
											Name: "RefreshToken",
										},
									},
									Args: []ast.Expr{
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "r",
												},
												Sel: &ast.Ident{
													Name: "Context",
												},
											},
										},
										&ast.Ident{
											Name: "refreshToken",
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
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "errs",
												},
												Sel: &ast.Ident{
													Name: "RenderToHTTPResponse",
												},
											},
											Args: []ast.Expr{
												&ast.Ident{
													Name: "err",
												},
												&ast.Ident{
													Name: "w",
												},
												&ast.Ident{
													Name: "r",
												},
											},
										},
									},
									&ast.ReturnStmt{},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "response",
								},
								&ast.Ident{
									Name: "err",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.Ident{
										Name: "NewTokenPairDTO",
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "tokenPair",
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
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "errs",
												},
												Sel: &ast.Ident{
													Name: "RenderToHTTPResponse",
												},
											},
											Args: []ast.Expr{
												&ast.Ident{
													Name: "err",
												},
												&ast.Ident{
													Name: "w",
												},
												&ast.Ident{
													Name: "r",
												},
											},
										},
									},
									&ast.ReturnStmt{},
								},
							},
						},
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "render",
									},
									Sel: &ast.Ident{
										Name: "Status",
									},
								},
								Args: []ast.Expr{
									&ast.Ident{
										Name: "r",
									},
									&ast.SelectorExpr{
										X: &ast.Ident{
											Name: "http",
										},
										Sel: &ast.Ident{
											Name: "StatusOK",
										},
									},
								},
							},
						},
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "render",
									},
									Sel: &ast.Ident{
										Name: "JSON",
									},
								},
								Args: []ast.Expr{
									&ast.Ident{
										Name: "w",
									},
									&ast.Ident{
										Name: "r",
									},
									&ast.Ident{
										Name: "response",
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
						&ast.Field{
							Names: []*ast.Ident{
								&ast.Ident{
									Name: "h",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "AuthHandler",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "ChiRouter",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
					Results: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "chi",
									},
									Sel: &ast.Ident{
										Name: "Router",
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
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "router",
									},
									Sel: &ast.Ident{
										Name: "Route",
									},
								},
								Args: []ast.Expr{
									&ast.BasicLit{
										Kind:  token.STRING,
										Value: "\"/\"",
									},
									&ast.FuncLit{
										Type: &ast.FuncType{
											Params: &ast.FieldList{
												List: []*ast.Field{
													&ast.Field{
														Names: []*ast.Ident{
															&ast.Ident{
																Name: "g",
															},
														},
														Type: &ast.SelectorExpr{
															X: &ast.Ident{
																Name: "chi",
															},
															Sel: &ast.Ident{
																Name: "Router",
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
															X: &ast.Ident{
																Name: "g",
															},
															Sel: &ast.Ident{
																Name: "Post",
															},
														},
														Args: []ast.Expr{
															&ast.BasicLit{
																Kind:  token.STRING,
																Value: "\"/obtain\"",
															},
															&ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "h",
																},
																Sel: &ast.Ident{
																	Name: "ObtainTokenPair",
																},
															},
														},
													},
												},
												&ast.ExprStmt{
													X: &ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X: &ast.Ident{
																Name: "g",
															},
															Sel: &ast.Ident{
																Name: "Post",
															},
														},
														Args: []ast.Expr{
															&ast.BasicLit{
																Kind:  token.STRING,
																Value: "\"/refresh\"",
															},
															&ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "h",
																},
																Sel: &ast.Ident{
																	Name: "RefreshTokenPair",
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
									Name: "router",
								},
							},
						},
					},
				},
			},
		},
		FileStart: 1,
		FileEnd:   2737,
		Imports: []*ast.ImportSpec{
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
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"github.com/go-chi/render\"",
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf(`"%s/internal/pkg/errs"`, h.project.Module),
				},
			},
		},
		Comments: []*ast.CommentGroup{
			&ast.CommentGroup{
				List: []*ast.Comment{
					&ast.Comment{
						Slash: 371,
						Text:  "// ObtainTokenPair",
					},
					&ast.Comment{
						Slash: 390,
						Text:  "//",
					},
					&ast.Comment{
						Slash: 393,
						Text:  "// @Tags auth",
					},
					&ast.Comment{
						Slash: 407,
						Text:  "// @Accept json",
					},
					&ast.Comment{
						Slash: 423,
						Text:  "// @Produce json",
					},
					&ast.Comment{
						Slash: 440,
						Text:  "// @Param form body ObtainTokenDTO true \"Obtain token pair\"",
					},
					&ast.Comment{
						Slash: 500,
						Text:  "// @Success 201 {object} TokenPairDTO \"Token pair\"",
					},
					&ast.Comment{
						Slash: 551,
						Text:  "// @Failure 400 {object} errs.Error \"Invalid request body or validation error\"",
					},
					&ast.Comment{
						Slash: 630,
						Text:  "// @Failure 401 {object} errs.Error \"Unauthorized\"",
					},
					&ast.Comment{
						Slash: 681,
						Text:  "// @Failure 404 {object} errs.Error \"Not found\"",
					},
					&ast.Comment{
						Slash: 729,
						Text:  "// @Failure 500 {object} errs.Error \"Internal server error\"",
					},
					&ast.Comment{
						Slash: 789,
						Text:  "// @Router /api/v1/auth/obtain [POST]",
					},
				},
			},
		},
	}
}
