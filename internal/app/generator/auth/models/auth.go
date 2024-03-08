package models

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
		Name: &ast.Ident{
			Name: "models",
		},
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
						Name: ast.NewIdent("userModels"),
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/app/user/models"`, m.project.Module),
						},
					},
					&ast.ImportSpec{
						Name: &ast.Ident{
							Name: "validation",
						},
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
						Name: &ast.Ident{
							Name: "Token",
						},
						Type: &ast.Ident{
							Name: "string",
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
									Name: "t",
								},
							},
							Type: &ast.Ident{
								Name: "Token",
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "String",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.Ident{
									Name: "string",
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
									Fun: &ast.Ident{
										Name: "string",
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "t",
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
						Name: &ast.Ident{
							Name: "TokenPair",
						},
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											{
												Name: "Access",
											},
										},
										Type: &ast.Ident{
											Name: "Token",
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"access\"  form:\"access\"`",
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "Refresh",
											},
										},
										Type: &ast.Ident{
											Name: "Token",
										},
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
								{
									Name: "c",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "TokenPair",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "Validate",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
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
						&ast.AssignStmt{
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
											Name: "validation",
										},
										Sel: &ast.Ident{
											Name: "ValidateStruct",
										},
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "c",
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "validation",
												},
												Sel: &ast.Ident{
													Name: "Field",
												},
											},
											Args: []ast.Expr{
												&ast.UnaryExpr{
													Op: token.AND,
													X: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "c",
														},
														Sel: &ast.Ident{
															Name: "Access",
														},
													},
												},
											},
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "validation",
												},
												Sel: &ast.Ident{
													Name: "Field",
												},
											},
											Args: []ast.Expr{
												&ast.UnaryExpr{
													Op: token.AND,
													X: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "c",
														},
														Sel: &ast.Ident{
															Name: "Refresh",
														},
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
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "errs",
													},
													Sel: &ast.Ident{
														Name: "FromValidationError",
													},
												},
												Args: []ast.Expr{
													&ast.Ident{
														Name: "err",
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
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{
							Name: "Login",
						},
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											{
												Name: "Email",
											},
										},
										Type: &ast.Ident{
											Name: "string",
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"email\"    form:\"email\"`",
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "Password",
											},
										},
										Type: &ast.Ident{
											Name: "string",
										},
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
								{
									Name: "c",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "Login",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "Validate",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
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
						&ast.AssignStmt{
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
											Name: "validation",
										},
										Sel: &ast.Ident{
											Name: "ValidateStruct",
										},
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "c",
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "validation",
												},
												Sel: &ast.Ident{
													Name: "Field",
												},
											},
											Args: []ast.Expr{
												&ast.UnaryExpr{
													Op: token.AND,
													X: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "c",
														},
														Sel: &ast.Ident{
															Name: "Email",
														},
													},
												},
												&ast.SelectorExpr{
													X: &ast.Ident{
														Name: "is",
													},
													Sel: &ast.Ident{
														Name: "Email",
													},
												},
											},
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "validation",
												},
												Sel: &ast.Ident{
													Name: "Field",
												},
											},
											Args: []ast.Expr{
												&ast.UnaryExpr{
													Op: token.AND,
													X: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "c",
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
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "errs",
													},
													Sel: &ast.Ident{
														Name: "FromValidationError",
													},
												},
												Args: []ast.Expr{
													&ast.Ident{
														Name: "err",
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
			&ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{
							{
								Name: "Guest",
							},
						},
						Values: []ast.Expr{
							&ast.UnaryExpr{
								Op: token.AND,
								X: &ast.CompositeLit{
									Type: ast.NewIdent("userModels.User"),
									Elts: []ast.Expr{
										&ast.KeyValueExpr{
											Key: &ast.Ident{
												Name: "ID",
											},
											Value: &ast.BasicLit{
												Kind:  token.STRING,
												Value: `""`,
											},
										},
										&ast.KeyValueExpr{
											Key: &ast.Ident{
												Name: "FirstName",
											},
											Value: &ast.BasicLit{
												Kind:  token.STRING,
												Value: `""`,
											},
										},
										&ast.KeyValueExpr{
											Key: &ast.Ident{
												Name: "LastName",
											},
											Value: &ast.BasicLit{
												Kind:  token.STRING,
												Value: `""`,
											},
										},
										&ast.KeyValueExpr{
											Key: &ast.Ident{
												Name: "Password",
											},
											Value: &ast.BasicLit{
												Kind:  token.STRING,
												Value: `""`,
											},
										},
										&ast.KeyValueExpr{
											Key: &ast.Ident{
												Name: "Email",
											},
											Value: &ast.BasicLit{
												Kind:  token.STRING,
												Value: `""`,
											},
										},
										&ast.KeyValueExpr{
											Key: &ast.Ident{
												Name: "CreatedAt",
											},
											Value: &ast.CompositeLit{
												Type: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "time",
													},
													Sel: &ast.Ident{
														Name: "Time",
													},
												},
											},
										},
										&ast.KeyValueExpr{
											Key: &ast.Ident{
												Name: "UpdatedAt",
											},
											Value: &ast.CompositeLit{
												Type: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "time",
													},
													Sel: &ast.Ident{
														Name: "Time",
													},
												},
											},
										},
										&ast.KeyValueExpr{
											Key: &ast.Ident{
												Name: "GroupID",
											},
											Value: ast.NewIdent("userModels.GroupIDGuest"),
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

var destinationPath = "."

func (m ModelAuth) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "app", "auth", "models", "auth.go")
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
		SourcePath: "templates/internal/auth/models/auth_mock.go.tmpl",
		DestinationPath: path.Join(
			destinationPath,
			"internal",
			"app",
			"auth",
			"models",
			"mock",
			"auth.go",
		),
		Name: "auth mock models",
	}
	if err := mock.RenderToFile(m.project); err != nil {
		return err
	}
	return nil
}
