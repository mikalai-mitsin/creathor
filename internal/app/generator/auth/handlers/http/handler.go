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
		Name:    ast.NewIdent("handlers"),
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
						Name: ast.NewIdent("AuthHandler"),
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											ast.NewIdent("authUseCase"),
										},
										Type: ast.NewIdent("authUseCase"),
									},
									{
										Names: []*ast.Ident{
											ast.NewIdent("logger"),
										},
										Type: ast.NewIdent("logger"),
									},
								},
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Name: ast.NewIdent("NewAuthHandler"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("authUseCase"),
								},
								Type: ast.NewIdent("authUseCase"),
							},
							{
								Names: []*ast.Ident{
									ast.NewIdent("logger"),
								},
								Type: ast.NewIdent("logger"),
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.StarExpr{
									X: ast.NewIdent("AuthHandler"),
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
										Type: ast.NewIdent("AuthHandler"),
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key:   ast.NewIdent("authUseCase"),
												Value: ast.NewIdent("authUseCase"),
											},
											&ast.KeyValueExpr{
												Key:   ast.NewIdent("logger"),
												Value: ast.NewIdent("logger"),
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
								ast.NewIdent("h"),
							},
							Type: &ast.StarExpr{
								X: ast.NewIdent("AuthHandler"),
							},
						},
					},
				},
				Doc: &ast.CommentGroup{
					List: []*ast.Comment{
						{
							Text: "// ObtainTokenPair",
						},
						{
							Text: "//",
						},
						{
							Text: "// @Tags auth",
						},
						{
							Text: "// @Accept json",
						},
						{
							Text: "// @Produce json",
						},
						{
							Text: "// @Param form body ObtainTokenDTO true \"Obtain token pair\"",
						},
						{
							Text: "// @Success 200 {object} TokenPairDTO \"Token pair\"",
						},
						{
							Text: "// @Failure 400 {object} errs.Error \"Invalid request body or validation error\"",
						},
						{
							Text: "// @Failure 401 {object} errs.Error \"Unauthorized\"",
						},
						{
							Text: "// @Failure 404 {object} errs.Error \"Not found\"",
						},
						{
							Text: "// @Failure 500 {object} errs.Error \"Internal server error\"",
						},
						{
							Text: "// @Router /api/v1/auth/obtain [POST]",
						},
					},
				},
				Name: ast.NewIdent("ObtainTokenPair"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("w"),
								},
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("http"),
									Sel: ast.NewIdent("ResponseWriter"),
								},
							},
							{
								Names: []*ast.Ident{
									ast.NewIdent("r"),
								},
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X:   ast.NewIdent("http"),
										Sel: ast.NewIdent("Request"),
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
								ast.NewIdent("createDTO"),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: ast.NewIdent("NewObtainTokenDTO"),
									Args: []ast.Expr{
										ast.NewIdent("r"),
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
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("errs"),
												Sel: ast.NewIdent("RenderToHTTPResponse"),
											},
											Args: []ast.Expr{
												ast.NewIdent("err"),
												ast.NewIdent("w"),
												ast.NewIdent("r"),
											},
										},
									},
									&ast.ReturnStmt{},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("create"),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("createDTO"),
										Sel: ast.NewIdent("toEntity"),
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
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("errs"),
												Sel: ast.NewIdent("RenderToHTTPResponse"),
											},
											Args: []ast.Expr{
												ast.NewIdent("err"),
												ast.NewIdent("w"),
												ast.NewIdent("r"),
											},
										},
									},
									&ast.ReturnStmt{},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("tokenPair"),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X:   ast.NewIdent("h"),
											Sel: ast.NewIdent("authUseCase"),
										},
										Sel: ast.NewIdent("CreateToken"),
									},
									Args: []ast.Expr{
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("r"),
												Sel: ast.NewIdent("Context"),
											},
										},
										ast.NewIdent("create"),
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
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("errs"),
												Sel: ast.NewIdent("RenderToHTTPResponse"),
											},
											Args: []ast.Expr{
												ast.NewIdent("err"),
												ast.NewIdent("w"),
												ast.NewIdent("r"),
											},
										},
									},
									&ast.ReturnStmt{},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("response"),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: ast.NewIdent("NewTokenPairDTO"),
									Args: []ast.Expr{
										ast.NewIdent("tokenPair"),
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
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("errs"),
												Sel: ast.NewIdent("RenderToHTTPResponse"),
											},
											Args: []ast.Expr{
												ast.NewIdent("err"),
												ast.NewIdent("w"),
												ast.NewIdent("r"),
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
									X:   ast.NewIdent("render"),
									Sel: ast.NewIdent("Status"),
								},
								Args: []ast.Expr{
									ast.NewIdent("r"),
									&ast.SelectorExpr{
										X:   ast.NewIdent("http"),
										Sel: ast.NewIdent("StatusOK"),
									},
								},
							},
						},
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("render"),
									Sel: ast.NewIdent("JSON"),
								},
								Args: []ast.Expr{
									ast.NewIdent("w"),
									ast.NewIdent("r"),
									ast.NewIdent("response"),
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
								ast.NewIdent("h"),
							},
							Type: &ast.StarExpr{
								X: ast.NewIdent("AuthHandler"),
							},
						},
					},
				},
				Doc: &ast.CommentGroup{
					List: []*ast.Comment{
						{
							Text: "// RefreshTokenPair",
						},
						{
							Text: "//",
						},
						{
							Text: "// @Tags auth",
						},
						{
							Text: "// @Accept json",
						},
						{
							Text: "// @Produce json",
						},
						{
							Text: "// @Param form body RefreshTokenDTO true \"Refresh token pair\"",
						},
						{
							Text: "// @Success 200 {object} TokenPairDTO \"Token pair\"",
						},
						{
							Text: "// @Failure 400 {object} errs.Error \"Invalid request body or validation error\"",
						},
						{
							Text: "// @Failure 401 {object} errs.Error \"Unauthorized\"",
						},
						{
							Text: "// @Failure 404 {object} errs.Error \"Not found\"",
						},
						{
							Text: "// @Failure 500 {object} errs.Error \"Internal server error\"",
						},
						{
							Text: "// @Router /api/v1/auth/refresh [POST]",
						},
					},
				},
				Name: ast.NewIdent("RefreshTokenPair"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("w"),
								},
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("http"),
									Sel: ast.NewIdent("ResponseWriter"),
								},
							},
							{
								Names: []*ast.Ident{
									ast.NewIdent("r"),
								},
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X:   ast.NewIdent("http"),
										Sel: ast.NewIdent("Request"),
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
								ast.NewIdent("refreshTokenDTO"),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: ast.NewIdent("NewRefreshTokenDTO"),
									Args: []ast.Expr{
										ast.NewIdent("r"),
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
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("errs"),
												Sel: ast.NewIdent("RenderToHTTPResponse"),
											},
											Args: []ast.Expr{
												ast.NewIdent("err"),
												ast.NewIdent("w"),
												ast.NewIdent("r"),
											},
										},
									},
									&ast.ReturnStmt{},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("refreshToken"),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("refreshTokenDTO"),
										Sel: ast.NewIdent("toEntity"),
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
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("errs"),
												Sel: ast.NewIdent("RenderToHTTPResponse"),
											},
											Args: []ast.Expr{
												ast.NewIdent("err"),
												ast.NewIdent("w"),
												ast.NewIdent("r"),
											},
										},
									},
									&ast.ReturnStmt{},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("tokenPair"),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X:   ast.NewIdent("h"),
											Sel: ast.NewIdent("authUseCase"),
										},
										Sel: ast.NewIdent("RefreshToken"),
									},
									Args: []ast.Expr{
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("r"),
												Sel: ast.NewIdent("Context"),
											},
										},
										ast.NewIdent("refreshToken"),
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
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("errs"),
												Sel: ast.NewIdent("RenderToHTTPResponse"),
											},
											Args: []ast.Expr{
												ast.NewIdent("err"),
												ast.NewIdent("w"),
												ast.NewIdent("r"),
											},
										},
									},
									&ast.ReturnStmt{},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("response"),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: ast.NewIdent("NewTokenPairDTO"),
									Args: []ast.Expr{
										ast.NewIdent("tokenPair"),
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
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("errs"),
												Sel: ast.NewIdent("RenderToHTTPResponse"),
											},
											Args: []ast.Expr{
												ast.NewIdent("err"),
												ast.NewIdent("w"),
												ast.NewIdent("r"),
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
									X:   ast.NewIdent("render"),
									Sel: ast.NewIdent("Status"),
								},
								Args: []ast.Expr{
									ast.NewIdent("r"),
									&ast.SelectorExpr{
										X:   ast.NewIdent("http"),
										Sel: ast.NewIdent("StatusOK"),
									},
								},
							},
						},
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("render"),
									Sel: ast.NewIdent("JSON"),
								},
								Args: []ast.Expr{
									ast.NewIdent("w"),
									ast.NewIdent("r"),
									ast.NewIdent("response"),
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
								ast.NewIdent("h"),
							},
							Type: &ast.StarExpr{
								X: ast.NewIdent("AuthHandler"),
							},
						},
					},
				},
				Name: ast.NewIdent("ChiRouter"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("chi"),
									Sel: ast.NewIdent("Router"),
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
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("router"),
									Sel: ast.NewIdent("Route"),
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
													{
														Names: []*ast.Ident{
															ast.NewIdent("g"),
														},
														Type: &ast.SelectorExpr{
															X:   ast.NewIdent("chi"),
															Sel: ast.NewIdent("Router"),
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
															X:   ast.NewIdent("g"),
															Sel: ast.NewIdent("Post"),
														},
														Args: []ast.Expr{
															&ast.BasicLit{
																Kind:  token.STRING,
																Value: "\"/obtain\"",
															},
															&ast.SelectorExpr{
																X:   ast.NewIdent("h"),
																Sel: ast.NewIdent("ObtainTokenPair"),
															},
														},
													},
												},
												&ast.ExprStmt{
													X: &ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent("g"),
															Sel: ast.NewIdent("Post"),
														},
														Args: []ast.Expr{
															&ast.BasicLit{
																Kind:  token.STRING,
																Value: "\"/refresh\"",
															},
															&ast.SelectorExpr{
																X:   ast.NewIdent("h"),
																Sel: ast.NewIdent("RefreshTokenPair"),
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
								ast.NewIdent("router"),
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
					Value: "\"net/http\"",
				},
			},
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"github.com/go-chi/chi/v5\"",
				},
			},
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"github.com/go-chi/render\"",
				},
			},
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf(`"%s/internal/pkg/errs"`, h.project.Module),
				},
			},
		},
	}
}
