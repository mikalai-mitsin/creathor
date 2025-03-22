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

func (m *Middlewares) file() *ast.File {
	return &ast.File{
		Name: ast.NewIdent("handlers"),
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
							Value: `"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/app/auth/entities"`, m.project.Module),
						},
					},
					&ast.ImportSpec{
						Name: ast.NewIdent("userEntities"),
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/app/user/entities"`, m.project.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/pkg/auth"`, m.project.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/pkg/errs"`, m.project.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/pkg/configs"`, m.project.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/pkg/log"`, m.project.Module),
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
							Value: `"strings"`,
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.CONST,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("headerAuthorize"),
						},
						Values: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: `"authorization"`,
							},
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("expectedScheme"),
						},
						Values: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: `"bearer"`,
							},
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: ast.NewIdent("AuthService"),
						Type: &ast.InterfaceType{
							Methods: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											ast.NewIdent("Auth"),
										},
										Type: &ast.FuncType{
											Params: &ast.FieldList{
												List: []*ast.Field{
													{
														Type: &ast.SelectorExpr{
															X:   ast.NewIdent("context"),
															Sel: ast.NewIdent("Context"),
														},
													},
													{
														Type: &ast.SelectorExpr{
															X:   ast.NewIdent("entities"),
															Sel: ast.NewIdent("Token"),
														},
													},
												},
											},
											Results: &ast.FieldList{
												List: []*ast.Field{
													{
														Type: &ast.SelectorExpr{
															X:   ast.NewIdent("userEntities"),
															Sel: ast.NewIdent("User"),
														},
													},
													{
														Type: ast.NewIdent("error"),
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
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: ast.NewIdent("AuthMiddleware"),
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											ast.NewIdent("logger"),
										},
										Type: &ast.StarExpr{
											X: &ast.SelectorExpr{
												X:   ast.NewIdent("log"),
												Sel: ast.NewIdent("Log"),
											},
										},
									},
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
									{
										Names: []*ast.Ident{
											ast.NewIdent("authService"),
										},
										Type: ast.NewIdent("AuthService"),
									},
								},
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Name: ast.NewIdent("NewAuthMiddleware"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("authService"),
								},
								Type: ast.NewIdent("AuthService"),
							},
							{
								Names: []*ast.Ident{
									ast.NewIdent("logger"),
								},
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X:   ast.NewIdent("log"),
										Sel: ast.NewIdent("Log"),
									},
								},
							},
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
									X: ast.NewIdent("AuthMiddleware"),
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
										Type: ast.NewIdent("AuthMiddleware"),
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key:   ast.NewIdent("authService"),
												Value: ast.NewIdent("authService"),
											},
											&ast.KeyValueExpr{
												Key:   ast.NewIdent("logger"),
												Value: ast.NewIdent("logger"),
											},
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
			&ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								ast.NewIdent("m"),
							},
							Type: &ast.StarExpr{
								X: ast.NewIdent("AuthMiddleware"),
							},
						},
					},
				},
				Name: ast.NewIdent("UnaryServerInterceptor"),
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
							{
								Names: []*ast.Ident{
									ast.NewIdent("req"),
								},
								Type: ast.NewIdent("any"),
							},
							{
								Names: []*ast.Ident{
									ast.NewIdent("_"),
								},
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X:   ast.NewIdent("grpc"),
										Sel: ast.NewIdent("UnaryServerInfo"),
									},
								},
							},
							{
								Names: []*ast.Ident{
									ast.NewIdent("handler"),
								},
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("grpc"),
									Sel: ast.NewIdent("UnaryHandler"),
								},
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: ast.NewIdent("any"),
							},
							{
								Type: ast.NewIdent("error"),
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("newCtx"),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("m"),
										Sel: ast.NewIdent("auth"),
									},
									Args: []ast.Expr{
										ast.NewIdent("ctx"),
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
									&ast.ReturnStmt{
										Results: []ast.Expr{
											ast.NewIdent("nil"),
											ast.NewIdent("err"),
										},
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.CallExpr{
									Fun: ast.NewIdent("handler"),
									Args: []ast.Expr{
										ast.NewIdent("newCtx"),
										ast.NewIdent("req"),
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
								ast.NewIdent("m"),
							},
							Type: &ast.StarExpr{
								X: ast.NewIdent("AuthMiddleware"),
							},
						},
					},
				},
				Name: ast.NewIdent("auth"),
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
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("context"),
									Sel: ast.NewIdent("Context"),
								},
							},
							{
								Type: ast.NewIdent("error"),
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.DeclStmt{
							Decl: &ast.GenDecl{
								Tok: token.VAR,
								Specs: []ast.Spec{
									&ast.ValueSpec{
										Names: []*ast.Ident{
											ast.NewIdent("token"),
										},
										Type: &ast.SelectorExpr{
											X:   ast.NewIdent("entities"),
											Sel: ast.NewIdent("Token"),
										},
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("token"),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("m"),
										Sel: ast.NewIdent("authFromMD"),
									},
									Args: []ast.Expr{
										ast.NewIdent("ctx"),
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
									&ast.ReturnStmt{
										Results: []ast.Expr{
											ast.NewIdent("ctx"),
											ast.NewIdent("err"),
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.BinaryExpr{
								X:  ast.NewIdent("token"),
								Op: token.EQL,
								Y: &ast.BasicLit{
									Kind:  token.STRING,
									Value: `""`,
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X:   ast.NewIdent("auth"),
													Sel: ast.NewIdent("PutUser"),
												},
												Args: []ast.Expr{
													ast.NewIdent("ctx"),
													&ast.SelectorExpr{
														X:   ast.NewIdent("entities"),
														Sel: ast.NewIdent("Guest"),
													},
												},
											},
											ast.NewIdent("nil"),
										},
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("user"),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X:   ast.NewIdent("m"),
											Sel: ast.NewIdent("authService"),
										},
										Sel: ast.NewIdent("Auth"),
									},
									Args: []ast.Expr{
										ast.NewIdent("ctx"),
										ast.NewIdent("token"),
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
									&ast.ReturnStmt{
										Results: []ast.Expr{
											ast.NewIdent("ctx"),
											ast.NewIdent("err"),
										},
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("newCtx"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("auth"),
										Sel: ast.NewIdent("PutUser"),
									},
									Args: []ast.Expr{
										ast.NewIdent("ctx"),
										ast.NewIdent("user"),
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								ast.NewIdent("newCtx"),
								ast.NewIdent("nil"),
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
								ast.NewIdent("m"),
							},
							Type: &ast.StarExpr{
								X: ast.NewIdent("AuthMiddleware"),
							},
						},
					},
				},
				Name: ast.NewIdent("authFromMD"),
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
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("entities"),
									Sel: ast.NewIdent("Token"),
								},
							},
							{
								Type: ast.NewIdent("error"),
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("val"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("metautils"),
												Sel: ast.NewIdent("ExtractIncoming"),
											},
											Args: []ast.Expr{
												ast.NewIdent("ctx"),
											},
										},
										Sel: ast.NewIdent("Get"),
									},
									Args: []ast.Expr{
										ast.NewIdent("headerAuthorize"),
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.BinaryExpr{
								X:  ast.NewIdent("val"),
								Op: token.EQL,
								Y: &ast.BasicLit{
									Kind:  token.STRING,
									Value: `""`,
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: `""`,
											},
											ast.NewIdent("nil"),
										},
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("splits"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("strings"),
										Sel: ast.NewIdent("SplitN"),
									},
									Args: []ast.Expr{
										ast.NewIdent("val"),
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: `" "`,
										},
										&ast.BasicLit{
											Kind:  token.INT,
											Value: "2",
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.BinaryExpr{
								X: &ast.CallExpr{
									Fun: ast.NewIdent("len"),
									Args: []ast.Expr{
										ast.NewIdent("splits"),
									},
								},
								Op: token.LSS,
								Y: &ast.BasicLit{
									Kind:  token.INT,
									Value: "2",
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: `""`,
											},
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X:   ast.NewIdent("errs"),
													Sel: ast.NewIdent("NewUnauthenticatedError"),
												},
											},
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.UnaryExpr{
								Op: token.NOT,
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("strings"),
										Sel: ast.NewIdent("EqualFold"),
									},
									Args: []ast.Expr{
										&ast.IndexExpr{
											X: ast.NewIdent("splits"),
											Index: &ast.BasicLit{
												Kind:  token.INT,
												Value: "0",
											},
										},
										ast.NewIdent("expectedScheme"),
									},
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: `""`,
											},
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X:   ast.NewIdent("errs"),
													Sel: ast.NewIdent("NewUnauthenticatedError"),
												},
											},
										},
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("bearerToken"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("strings"),
										Sel: ast.NewIdent("TrimSpace"),
									},
									Args: []ast.Expr{
										&ast.IndexExpr{
											X: ast.NewIdent("splits"),
											Index: &ast.BasicLit{
												Kind:  token.INT,
												Value: "1",
											},
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.BinaryExpr{
								X:  ast.NewIdent("bearerToken"),
								Op: token.EQL,
								Y: &ast.BasicLit{
									Kind:  token.STRING,
									Value: `""`,
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: `""`,
											},
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X:   ast.NewIdent("errs"),
													Sel: ast.NewIdent("NewUnauthenticatedError"),
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
										X:   ast.NewIdent("entities"),
										Sel: ast.NewIdent("Token"),
									},
									Args: []ast.Expr{
										&ast.IndexExpr{
											X: ast.NewIdent("splits"),
											Index: &ast.BasicLit{
												Kind:  token.INT,
												Value: "1",
											},
										},
									},
								},
								ast.NewIdent("nil"),
							},
						},
					},
				},
			},
		},
	}
}

func (m *Middlewares) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "app", "auth", "handlers", "grpc", "middlewares.go")
	if err := os.MkdirAll(path.Dir(filename), 0777); err != nil {
		return err
	}
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = m.file()
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
