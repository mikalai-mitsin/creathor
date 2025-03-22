package entities

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"

	"github.com/mikalai-mitsin/creathor/internal/pkg/tmpl"

	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
)

type ModelAuth struct {
	project *configs.Project
}

func NewModelAuth(project *configs.Project) *ModelAuth {
	return &ModelAuth{project: project}
}

func (m ModelAuth) file() *ast.File {
	return &ast.File{
		Name: ast.NewIdent("entities"),
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"time"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/pkg/errs"`, m.project.Module),
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
						Name: ast.NewIdent("validation"),
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"github.com/go-ozzo/ozzo-validation/v4"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"github.com/go-ozzo/ozzo-validation/v4/is"`,
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: ast.NewIdent("Token"),
						Type: ast.NewIdent("string"),
					},
				},
			},
			&ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								ast.NewIdent("t"),
							},
							Type: ast.NewIdent("Token"),
						},
					},
				},
				Name: ast.NewIdent("String"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: ast.NewIdent("string"),
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.CallExpr{
									Fun: ast.NewIdent("string"),
									Args: []ast.Expr{
										ast.NewIdent("t"),
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
						Name: ast.NewIdent("TokenPair"),
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											ast.NewIdent("Access"),
										},
										Type: ast.NewIdent("Token"),
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"access\"  form:\"access\"`",
										},
									},
									{
										Names: []*ast.Ident{
											ast.NewIdent("Refresh"),
										},
										Type: ast.NewIdent("Token"),
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"refresh\" form:\"refresh\"`",
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
								ast.NewIdent("c"),
							},
							Type: &ast.StarExpr{
								X: ast.NewIdent("TokenPair"),
							},
						},
					},
				},
				Name: ast.NewIdent("Validate"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
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
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("validation"),
										Sel: ast.NewIdent("ValidateStruct"),
									},
									Args: []ast.Expr{
										ast.NewIdent("c"),
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("validation"),
												Sel: ast.NewIdent("Field"),
											},
											Args: []ast.Expr{
												&ast.UnaryExpr{
													Op: token.AND,
													X: &ast.SelectorExpr{
														X:   ast.NewIdent("c"),
														Sel: ast.NewIdent("Access"),
													},
												},
											},
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("validation"),
												Sel: ast.NewIdent("Field"),
											},
											Args: []ast.Expr{
												&ast.UnaryExpr{
													Op: token.AND,
													X: &ast.SelectorExpr{
														X:   ast.NewIdent("c"),
														Sel: ast.NewIdent("Refresh"),
													},
												},
											},
										},
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
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X:   ast.NewIdent("errs"),
													Sel: ast.NewIdent("NewFromValidationError"),
												},
												Args: []ast.Expr{
													ast.NewIdent("err"),
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
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: ast.NewIdent("Login"),
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
											Value: "`json:\"email\"    form:\"email\"`",
										},
									},
									{
										Names: []*ast.Ident{
											ast.NewIdent("Password"),
										},
										Type: ast.NewIdent("string"),
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"password\" form:\"password\"`",
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
								ast.NewIdent("c"),
							},
							Type: &ast.StarExpr{
								X: ast.NewIdent("Login"),
							},
						},
					},
				},
				Name: ast.NewIdent("Validate"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
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
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("validation"),
										Sel: ast.NewIdent("ValidateStruct"),
									},
									Args: []ast.Expr{
										ast.NewIdent("c"),
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("validation"),
												Sel: ast.NewIdent("Field"),
											},
											Args: []ast.Expr{
												&ast.UnaryExpr{
													Op: token.AND,
													X: &ast.SelectorExpr{
														X:   ast.NewIdent("c"),
														Sel: ast.NewIdent("Email"),
													},
												},
												&ast.SelectorExpr{
													X:   ast.NewIdent("is"),
													Sel: ast.NewIdent("Email"),
												},
											},
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("validation"),
												Sel: ast.NewIdent("Field"),
											},
											Args: []ast.Expr{
												&ast.UnaryExpr{
													Op: token.AND,
													X: &ast.SelectorExpr{
														X:   ast.NewIdent("c"),
														Sel: ast.NewIdent("Password"),
													},
												},
											},
										},
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
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X:   ast.NewIdent("errs"),
													Sel: ast.NewIdent("NewFromValidationError"),
												},
												Args: []ast.Expr{
													ast.NewIdent("err"),
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
			&ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("Guest"),
						},
						Values: []ast.Expr{
							&ast.CompositeLit{
								Type: ast.NewIdent("userEntities.User"),
								Elts: []ast.Expr{
									&ast.KeyValueExpr{
										Key: ast.NewIdent("ID"),
										Value: &ast.BasicLit{
											Kind:  token.STRING,
											Value: `""`,
										},
									},
									&ast.KeyValueExpr{
										Key: ast.NewIdent("FirstName"),
										Value: &ast.BasicLit{
											Kind:  token.STRING,
											Value: `""`,
										},
									},
									&ast.KeyValueExpr{
										Key: ast.NewIdent("LastName"),
										Value: &ast.BasicLit{
											Kind:  token.STRING,
											Value: `""`,
										},
									},
									&ast.KeyValueExpr{
										Key: ast.NewIdent("Password"),
										Value: &ast.BasicLit{
											Kind:  token.STRING,
											Value: `""`,
										},
									},
									&ast.KeyValueExpr{
										Key: ast.NewIdent("Email"),
										Value: &ast.BasicLit{
											Kind:  token.STRING,
											Value: `""`,
										},
									},
									&ast.KeyValueExpr{
										Key: ast.NewIdent("CreatedAt"),
										Value: &ast.CompositeLit{
											Type: &ast.SelectorExpr{
												X:   ast.NewIdent("time"),
												Sel: ast.NewIdent("Time"),
											},
										},
									},
									&ast.KeyValueExpr{
										Key: ast.NewIdent("UpdatedAt"),
										Value: &ast.CompositeLit{
											Type: &ast.SelectorExpr{
												X:   ast.NewIdent("time"),
												Sel: ast.NewIdent("Time"),
											},
										},
									},
									&ast.KeyValueExpr{
										Key:   ast.NewIdent("GroupID"),
										Value: ast.NewIdent("userEntities.GroupIDGuest"),
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

var destinationPath = "."

func (m ModelAuth) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "app", "auth", "entities", "auth.go")
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
	mock := &tmpl.Template{
		SourcePath: "templates/internal/auth/entities/auth_mock.go.tmpl",
		DestinationPath: path.Join(
			destinationPath,
			"internal",
			"app",
			"auth",
			"entities",
			"mock",
			"auth.go",
		),
		Name: "auth mock entities",
	}
	if err := mock.RenderToFile(m.project); err != nil {
		return err
	}
	return nil
}
