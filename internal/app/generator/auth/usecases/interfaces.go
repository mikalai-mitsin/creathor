package usecases

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

type ServiceInterfaceAuth struct {
	project *configs.Project
}

func NewServiceInterfaceAuth(project *configs.Project) *ServiceInterfaceAuth {
	return &ServiceInterfaceAuth{project: project}
}

func (i ServiceInterfaceAuth) file() *ast.File {
	return &ast.File{
		Name: ast.NewIdent("usecases"),
		Decls: []ast.Decl{
			&ast.GenDecl{
				Doc: &ast.CommentGroup{
					List: []*ast.Comment{
						{
							Slash: token.NoPos,
							Text:  "//go:generate mockgen -source=interfaces.go -package=usecases -destination=interfaces_mock.go",
						},
					},
				},
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
							Value: `"time"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/app/auth/entities"`, i.project.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/pkg/log"`, i.project.Module),
						},
					},
					&ast.ImportSpec{
						Name: ast.NewIdent("userEntities"),
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/app/user/entities"`, i.project.Module),
						},
					},
				},
			},
			&ast.GenDecl{
				Doc: &ast.CommentGroup{
					List: []*ast.Comment{
						{
							Text: "// clock - clock interface",
						},
					},
				},
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: ast.NewIdent("clock"),
						Type: &ast.InterfaceType{
							Methods: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											ast.NewIdent("Now"),
										},
										Type: &ast.FuncType{
											Results: &ast.FieldList{
												List: []*ast.Field{
													{
														Type: ast.NewIdent("time.Time"),
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
						Name: ast.NewIdent("logger"),
						Type: &ast.InterfaceType{
							Methods: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											ast.NewIdent("Debug"),
										},
										Type: &ast.FuncType{
											Params: &ast.FieldList{
												List: []*ast.Field{
													{
														Names: []*ast.Ident{
															ast.NewIdent("msg"),
														},
														Type: ast.NewIdent("string"),
													},
													{
														Names: []*ast.Ident{
															ast.NewIdent("fields"),
														},
														Type: &ast.Ellipsis{
															Elt: &ast.SelectorExpr{
																X:   ast.NewIdent("log"),
																Sel: ast.NewIdent("Field"),
															},
														},
													},
												},
											},
										},
									},
									{
										Names: []*ast.Ident{
											ast.NewIdent("Info"),
										},
										Type: &ast.FuncType{
											Params: &ast.FieldList{
												List: []*ast.Field{
													{
														Names: []*ast.Ident{
															ast.NewIdent("msg"),
														},
														Type: ast.NewIdent("string"),
													},
													{
														Names: []*ast.Ident{
															ast.NewIdent("fields"),
														},
														Type: &ast.Ellipsis{
															Elt: &ast.SelectorExpr{
																X:   ast.NewIdent("log"),
																Sel: ast.NewIdent("Field"),
															},
														},
													},
												},
											},
										},
									},
									{
										Names: []*ast.Ident{
											ast.NewIdent("Print"),
										},
										Type: &ast.FuncType{
											Params: &ast.FieldList{
												List: []*ast.Field{
													{
														Names: []*ast.Ident{
															ast.NewIdent("msg"),
														},
														Type: ast.NewIdent("string"),
													},
													{
														Names: []*ast.Ident{
															ast.NewIdent("fields"),
														},
														Type: &ast.Ellipsis{
															Elt: &ast.SelectorExpr{
																X:   ast.NewIdent("log"),
																Sel: ast.NewIdent("Field"),
															},
														},
													},
												},
											},
										},
									},
									{
										Names: []*ast.Ident{
											ast.NewIdent("Warn"),
										},
										Type: &ast.FuncType{
											Params: &ast.FieldList{
												List: []*ast.Field{
													{
														Names: []*ast.Ident{
															ast.NewIdent("msg"),
														},
														Type: ast.NewIdent("string"),
													},
													{
														Names: []*ast.Ident{
															ast.NewIdent("fields"),
														},
														Type: &ast.Ellipsis{
															Elt: &ast.SelectorExpr{
																X:   ast.NewIdent("log"),
																Sel: ast.NewIdent("Field"),
															},
														},
													},
												},
											},
										},
									},
									{
										Names: []*ast.Ident{
											ast.NewIdent("Error"),
										},
										Type: &ast.FuncType{
											Params: &ast.FieldList{
												List: []*ast.Field{
													{
														Names: []*ast.Ident{
															ast.NewIdent("msg"),
														},
														Type: ast.NewIdent("string"),
													},
													{
														Names: []*ast.Ident{
															ast.NewIdent("fields"),
														},
														Type: &ast.Ellipsis{
															Elt: &ast.SelectorExpr{
																X:   ast.NewIdent("log"),
																Sel: ast.NewIdent("Field"),
															},
														},
													},
												},
											},
										},
									},
									{
										Names: []*ast.Ident{
											ast.NewIdent("Fatal"),
										},
										Type: &ast.FuncType{
											Params: &ast.FieldList{
												List: []*ast.Field{
													{
														Names: []*ast.Ident{
															ast.NewIdent("msg"),
														},
														Type: ast.NewIdent("string"),
													},
													{
														Names: []*ast.Ident{
															ast.NewIdent("fields"),
														},
														Type: &ast.Ellipsis{
															Elt: &ast.SelectorExpr{
																X:   ast.NewIdent("log"),
																Sel: ast.NewIdent("Field"),
															},
														},
													},
												},
											},
										},
									},
									{
										Names: []*ast.Ident{
											ast.NewIdent("Panic"),
										},
										Type: &ast.FuncType{
											Params: &ast.FieldList{
												List: []*ast.Field{
													{
														Names: []*ast.Ident{
															ast.NewIdent("msg"),
														},
														Type: ast.NewIdent("string"),
													},
													{
														Names: []*ast.Ident{
															ast.NewIdent("fields"),
														},
														Type: &ast.Ellipsis{
															Elt: &ast.SelectorExpr{
																X:   ast.NewIdent("log"),
																Sel: ast.NewIdent("Field"),
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
				},
			},
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: ast.NewIdent("authService"),
						Type: &ast.InterfaceType{
							Methods: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											ast.NewIdent("CreateToken"),
										},
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
															ast.NewIdent("login"),
														},
														Type: &ast.SelectorExpr{
															X:   ast.NewIdent("entities"),
															Sel: ast.NewIdent("Login"),
														},
													},
												},
											},
											Results: &ast.FieldList{
												List: []*ast.Field{
													{
														Type: &ast.SelectorExpr{
															X:   ast.NewIdent("entities"),
															Sel: ast.NewIdent("TokenPair"),
														},
													},
													{
														Type: ast.NewIdent("error"),
													},
												},
											},
										},
									},
									{
										Names: []*ast.Ident{
											ast.NewIdent("CreateTokenByUser"),
										},
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
															ast.NewIdent("user"),
														},
														Type: &ast.SelectorExpr{
															X:   ast.NewIdent("userEntities"),
															Sel: ast.NewIdent("User"),
														},
													},
												},
											},
											Results: &ast.FieldList{
												List: []*ast.Field{
													{
														Type: &ast.SelectorExpr{
															X:   ast.NewIdent("entities"),
															Sel: ast.NewIdent("TokenPair"),
														},
													},
													{
														Type: ast.NewIdent("error"),
													},
												},
											},
										},
									},
									{
										Names: []*ast.Ident{
											ast.NewIdent("RefreshToken"),
										},
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
															ast.NewIdent("refresh"),
														},
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
															X:   ast.NewIdent("entities"),
															Sel: ast.NewIdent("TokenPair"),
														},
													},
													{
														Type: ast.NewIdent("error"),
													},
												},
											},
										},
									},
									{
										Names: []*ast.Ident{
											ast.NewIdent("ValidateToken"),
										},
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
															ast.NewIdent("access"),
														},
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
														Type: ast.NewIdent("error"),
													},
												},
											},
										},
									},
									{
										Names: []*ast.Ident{
											ast.NewIdent("Auth"),
										},
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
															ast.NewIdent("access"),
														},
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
		},
	}
}

func (i ServiceInterfaceAuth) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "app", "auth", "usecases", "interfaces.go")
	if err := os.MkdirAll(path.Dir(filename), 0777); err != nil {
		return err
	}
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = i.file()
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
