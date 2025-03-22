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
						Name: ast.NewIdent("TokenPairDTO"),
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											ast.NewIdent("AccessToken"),
										},
										Type: ast.NewIdent("string"),
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"access_token\"`",
										},
									},
									{
										Names: []*ast.Ident{
											ast.NewIdent("RefreshToken"),
										},
										Type: ast.NewIdent("string"),
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
				Name: ast.NewIdent("NewTokenPairDTO"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("entity"),
								},
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("entities"),
									Sel: ast.NewIdent("TokenPair"),
								},
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: ast.NewIdent("TokenPairDTO"),
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
								ast.NewIdent("dto"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CompositeLit{
									Type: ast.NewIdent("TokenPairDTO"),
									Elts: []ast.Expr{
										&ast.KeyValueExpr{
											Key: ast.NewIdent("AccessToken"),
											Value: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.SelectorExpr{
														X:   ast.NewIdent("entity"),
														Sel: ast.NewIdent("Access"),
													},
													Sel: ast.NewIdent("String"),
												},
											},
										},
										&ast.KeyValueExpr{
											Key: ast.NewIdent("RefreshToken"),
											Value: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.SelectorExpr{
														X:   ast.NewIdent("entity"),
														Sel: ast.NewIdent("Refresh"),
													},
													Sel: ast.NewIdent("String"),
												},
											},
										},
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								ast.NewIdent("dto"),
								ast.NewIdent("nil"),
							},
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: ast.NewIdent("ObtainTokenDTO"),
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											ast.NewIdent("Email"),
										},
										Type: ast.NewIdent("string"),
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"email\"`",
										},
									},
									{
										Names: []*ast.Ident{
											ast.NewIdent("Password"),
										},
										Type: ast.NewIdent("string"),
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
				Name: ast.NewIdent("NewObtainTokenDTO"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
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
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: ast.NewIdent("ObtainTokenDTO"),
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
								ast.NewIdent("update"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CompositeLit{
									Type: ast.NewIdent("ObtainTokenDTO"),
								},
							},
						},
						&ast.IfStmt{
							Init: &ast.AssignStmt{
								Lhs: []ast.Expr{
									ast.NewIdent("err"),
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   ast.NewIdent("render"),
											Sel: ast.NewIdent("DecodeJSON"),
										},
										Args: []ast.Expr{
											&ast.SelectorExpr{
												X:   ast.NewIdent("r"),
												Sel: ast.NewIdent("Body"),
											},
											&ast.UnaryExpr{
												Op: token.AND,
												X:  ast.NewIdent("update"),
											},
										},
									},
								},
							},
							Cond: &ast.BinaryExpr{
								X:  ast.NewIdent("err"),
								Op: token.NEQ,
								Y:  ast.NewIdent("nil"),
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.CompositeLit{
												Type: ast.NewIdent("ObtainTokenDTO"),
											},
											ast.NewIdent("err"),
										},
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								ast.NewIdent("update"),
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
								ast.NewIdent("dto"),
							},
							Type: ast.NewIdent("ObtainTokenDTO"),
						},
					},
				},
				Name: ast.NewIdent("toEntity"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("entities"),
									Sel: ast.NewIdent("Login"),
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
								ast.NewIdent("login"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CompositeLit{
									Type: &ast.SelectorExpr{
										X:   ast.NewIdent("entities"),
										Sel: ast.NewIdent("Login"),
									},
									Elts: []ast.Expr{
										&ast.KeyValueExpr{
											Key: ast.NewIdent("Email"),
											Value: &ast.SelectorExpr{
												X:   ast.NewIdent("dto"),
												Sel: ast.NewIdent("Email"),
											},
										},
										&ast.KeyValueExpr{
											Key: ast.NewIdent("Password"),
											Value: &ast.SelectorExpr{
												X:   ast.NewIdent("dto"),
												Sel: ast.NewIdent("Password"),
											},
										},
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								ast.NewIdent("login"),
								ast.NewIdent("nil"),
							},
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: ast.NewIdent("RefreshTokenDTO"),
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											ast.NewIdent("RefreshToken"),
										},
										Type: ast.NewIdent("string"),
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
				Name: ast.NewIdent("NewRefreshTokenDTO"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
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
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: ast.NewIdent("RefreshTokenDTO"),
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
								ast.NewIdent("dto"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CompositeLit{
									Type: ast.NewIdent("RefreshTokenDTO"),
								},
							},
						},
						&ast.IfStmt{
							Init: &ast.AssignStmt{
								Lhs: []ast.Expr{
									ast.NewIdent("err"),
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   ast.NewIdent("render"),
											Sel: ast.NewIdent("DecodeJSON"),
										},
										Args: []ast.Expr{
											&ast.SelectorExpr{
												X:   ast.NewIdent("r"),
												Sel: ast.NewIdent("Body"),
											},
											&ast.UnaryExpr{
												Op: token.AND,
												X:  ast.NewIdent("dto"),
											},
										},
									},
								},
							},
							Cond: &ast.BinaryExpr{
								X:  ast.NewIdent("err"),
								Op: token.NEQ,
								Y:  ast.NewIdent("nil"),
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.CompositeLit{
												Type: ast.NewIdent("RefreshTokenDTO"),
											},
											ast.NewIdent("err"),
										},
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								ast.NewIdent("dto"),
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
								ast.NewIdent("dto"),
							},
							Type: ast.NewIdent("RefreshTokenDTO"),
						},
					},
				},
				Name: ast.NewIdent("toEntity"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
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
								ast.NewIdent("token"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("entities"),
										Sel: ast.NewIdent("Token"),
									},
									Args: []ast.Expr{
										&ast.SelectorExpr{
											X:   ast.NewIdent("dto"),
											Sel: ast.NewIdent("RefreshToken"),
										},
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								ast.NewIdent("token"),
								ast.NewIdent("nil"),
							},
						},
					},
				},
			},
		},
		FileStart: 1,
		FileEnd:   1334,
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
					Value: "\"github.com/go-chi/render\"",
				},
			},
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf(`"%s/internal/app/auth/entities"`, g.project.Module),
				},
			},
		},
	}
}
