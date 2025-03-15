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
	"path/filepath"
)

type DTOGenerator struct {
	project *configs.Project
}

func NewDTOGenerator(project *configs.Project) *DTOGenerator {
	return &DTOGenerator{project: project}
}

func (g *DTOGenerator) filename() string {
	return filepath.Join(
		"internal",
		"app",
		"auth",
		"handlers",
		"http",
		"dto.go",
	)
}

func (g *DTOGenerator) Sync() error {
	fileset := token.NewFileSet()
	filename := g.filename()
	if err := os.MkdirAll(path.Dir(filename), 0777); err != nil {
		return err
	}
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = g.file()
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

func (g *DTOGenerator) file() *ast.File {
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
							Value: "\"github.com/go-chi/render\"",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/app/auth/entities"`, g.project.Module),
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{
							Name: "TokenPairDTO",
						},
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									&ast.Field{
										Names: []*ast.Ident{
											&ast.Ident{
												Name: "AccessToken",
											},
										},
										Type: &ast.Ident{
											Name: "string",
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"access_token\"`",
										},
									},
									&ast.Field{
										Names: []*ast.Ident{
											&ast.Ident{
												Name: "RefreshToken",
											},
										},
										Type: &ast.Ident{
											Name: "string",
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"refresh_token\"`",
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
					Name: "NewTokenPairDTO",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Names: []*ast.Ident{
									&ast.Ident{
										Name: "entity",
									},
								},
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "entities",
									},
									Sel: &ast.Ident{
										Name: "TokenPair",
									},
								},
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Type: &ast.Ident{
									Name: "TokenPairDTO",
								},
							},
							&ast.Field{
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
									Name: "dto",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CompositeLit{
									Type: &ast.Ident{
										Name: "TokenPairDTO",
									},
									Elts: []ast.Expr{
										&ast.KeyValueExpr{
											Key: &ast.Ident{
												Name: "AccessToken",
											},
											Value: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "entity",
														},
														Sel: &ast.Ident{
															Name: "Access",
														},
													},
													Sel: &ast.Ident{
														Name: "String",
													},
												},
											},
										},
										&ast.KeyValueExpr{
											Key: &ast.Ident{
												Name: "RefreshToken",
											},
											Value: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "entity",
														},
														Sel: &ast.Ident{
															Name: "Refresh",
														},
													},
													Sel: &ast.Ident{
														Name: "String",
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
									Name: "dto",
								},
								&ast.Ident{
									Name: "nil",
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
							Name: "ObtainTokenDTO",
						},
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									&ast.Field{
										Names: []*ast.Ident{
											&ast.Ident{
												Name: "Email",
											},
										},
										Type: &ast.Ident{
											Name: "string",
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"email\"`",
										},
									},
									&ast.Field{
										Names: []*ast.Ident{
											&ast.Ident{
												Name: "Password",
											},
										},
										Type: &ast.Ident{
											Name: "string",
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"password\"`",
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
					Name: "NewObtainTokenDTO",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
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
					Results: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Type: &ast.Ident{
									Name: "ObtainTokenDTO",
								},
							},
							&ast.Field{
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
									Name: "update",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CompositeLit{
									Type: &ast.Ident{
										Name: "ObtainTokenDTO",
									},
								},
							},
						},
						&ast.IfStmt{
							Init: &ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.Ident{
										Name: "err",
									},
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "render",
											},
											Sel: &ast.Ident{
												Name: "DecodeJSON",
											},
										},
										Args: []ast.Expr{
											&ast.SelectorExpr{
												X: &ast.Ident{
													Name: "r",
												},
												Sel: &ast.Ident{
													Name: "Body",
												},
											},
											&ast.UnaryExpr{
												Op: token.AND,
												X: &ast.Ident{
													Name: "update",
												},
											},
										},
									},
								},
							},
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
											&ast.CompositeLit{
												Type: &ast.Ident{
													Name: "ObtainTokenDTO",
												},
											},
											&ast.Ident{
												Name: "err",
											},
										},
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.Ident{
									Name: "update",
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
						&ast.Field{
							Names: []*ast.Ident{
								&ast.Ident{
									Name: "dto",
								},
							},
							Type: &ast.Ident{
								Name: "ObtainTokenDTO",
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "toEntity",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
					Results: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "entities",
									},
									Sel: &ast.Ident{
										Name: "Login",
									},
								},
							},
							&ast.Field{
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
									Name: "login",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CompositeLit{
									Type: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "entities",
										},
										Sel: &ast.Ident{
											Name: "Login",
										},
									},
									Elts: []ast.Expr{
										&ast.KeyValueExpr{
											Key: &ast.Ident{
												Name: "Email",
											},
											Value: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "dto",
												},
												Sel: &ast.Ident{
													Name: "Email",
												},
											},
										},
										&ast.KeyValueExpr{
											Key: &ast.Ident{
												Name: "Password",
											},
											Value: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "dto",
												},
												Sel: &ast.Ident{
													Name: "Password",
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
									Name: "login",
								},
								&ast.Ident{
									Name: "nil",
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
							Name: "RefreshTokenDTO",
						},
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									&ast.Field{
										Names: []*ast.Ident{
											&ast.Ident{
												Name: "RefreshToken",
											},
										},
										Type: &ast.Ident{
											Name: "string",
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"refresh_token\"`",
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
					Name: "NewRefreshTokenDTO",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
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
					Results: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Type: &ast.Ident{
									Name: "RefreshTokenDTO",
								},
							},
							&ast.Field{
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
									Name: "dto",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CompositeLit{
									Type: &ast.Ident{
										Name: "RefreshTokenDTO",
									},
								},
							},
						},
						&ast.IfStmt{
							Init: &ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.Ident{
										Name: "err",
									},
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "render",
											},
											Sel: &ast.Ident{
												Name: "DecodeJSON",
											},
										},
										Args: []ast.Expr{
											&ast.SelectorExpr{
												X: &ast.Ident{
													Name: "r",
												},
												Sel: &ast.Ident{
													Name: "Body",
												},
											},
											&ast.UnaryExpr{
												Op: token.AND,
												X: &ast.Ident{
													Name: "dto",
												},
											},
										},
									},
								},
							},
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
											&ast.CompositeLit{
												Type: &ast.Ident{
													Name: "RefreshTokenDTO",
												},
											},
											&ast.Ident{
												Name: "err",
											},
										},
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.Ident{
									Name: "dto",
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
						&ast.Field{
							Names: []*ast.Ident{
								&ast.Ident{
									Name: "dto",
								},
							},
							Type: &ast.Ident{
								Name: "RefreshTokenDTO",
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "toEntity",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
					Results: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "entities",
									},
									Sel: &ast.Ident{
										Name: "Token",
									},
								},
							},
							&ast.Field{
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
									Name: "token",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "entities",
										},
										Sel: &ast.Ident{
											Name: "Token",
										},
									},
									Args: []ast.Expr{
										&ast.SelectorExpr{
											X: &ast.Ident{
												Name: "dto",
											},
											Sel: &ast.Ident{
												Name: "RefreshToken",
											},
										},
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.Ident{
									Name: "token",
								},
								&ast.Ident{
									Name: "nil",
								},
							},
						},
					},
				},
			},
		},
		FileStart: 1,
		FileEnd:   1334,
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
					Value: "\"github.com/go-chi/render\"",
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf(`"%s/internal/app/auth/entities"`, g.project.Module),
				},
			},
		},
	}
}
