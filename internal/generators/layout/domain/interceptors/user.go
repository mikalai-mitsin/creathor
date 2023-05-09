package interceptors

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"

	"github.com/018bf/creathor/internal/configs"
)

type InterceptorInterfaceUser struct {
	project *configs.Project
}

func NewInterceptorInterfaceUser(project *configs.Project) *InterceptorInterfaceUser {
	return &InterceptorInterfaceUser{project: project}
}

func (i InterceptorInterfaceUser) file() *ast.File {
	return &ast.File{
		Name: &ast.Ident{
			Name: "interceptors",
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
							Value: fmt.Sprintf(`"%s/internal/domain/models"`, i.project.Module),
						},
					},
				},
			},
			&ast.GenDecl{
				Doc: &ast.CommentGroup{
					List: []*ast.Comment{
						{
							Text: "//UserInterceptor - domain layer interceptor interface",
						},
						{
							Text: "//go:generate mockgen -build_flags=-mod=mod -destination mock/user.go . UserInterceptor",
						},
					},
				},
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{
							Name: "UserInterceptor",
						},
						Type: &ast.InterfaceType{
							Methods: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											{
												Name: "Get",
											},
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
													{
														Names: []*ast.Ident{
															{
																Name: "id",
															},
														},
														Type: &ast.SelectorExpr{
															X: &ast.Ident{
																Name: "models",
															},
															Sel: &ast.Ident{
																Name: "UUID",
															},
														},
													},
													{
														Names: []*ast.Ident{
															{
																Name: "requestUser",
															},
														},
														Type: &ast.StarExpr{
															X: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "models",
																},
																Sel: &ast.Ident{
																	Name: "User",
																},
															},
														},
													},
												},
											},
											Results: &ast.FieldList{
												List: []*ast.Field{
													{
														Type: &ast.StarExpr{
															X: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "models",
																},
																Sel: &ast.Ident{
																	Name: "User",
																},
															},
														},
													},
													{
														Type: &ast.Ident{
															Name: "error",
														},
													},
												},
											},
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "List",
											},
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
													{
														Names: []*ast.Ident{
															{
																Name: "filter",
															},
														},
														Type: &ast.StarExpr{
															X: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "models",
																},
																Sel: &ast.Ident{
																	Name: "UserFilter",
																},
															},
														},
													},
													{
														Names: []*ast.Ident{
															{
																Name: "requestUser",
															},
														},
														Type: &ast.StarExpr{
															X: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "models",
																},
																Sel: &ast.Ident{
																	Name: "User",
																},
															},
														},
													},
												},
											},
											Results: &ast.FieldList{
												List: []*ast.Field{
													{
														Type: &ast.ArrayType{
															Elt: &ast.StarExpr{
																X: &ast.SelectorExpr{
																	X: &ast.Ident{
																		Name: "models",
																	},
																	Sel: &ast.Ident{
																		Name: "User",
																	},
																},
															},
														},
													},
													{
														Type: &ast.Ident{
															Name: "uint64",
														},
													},
													{
														Type: &ast.Ident{
															Name: "error",
														},
													},
												},
											},
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "Create",
											},
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
													{
														Names: []*ast.Ident{
															{
																Name: "create",
															},
														},
														Type: &ast.StarExpr{
															X: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "models",
																},
																Sel: &ast.Ident{
																	Name: "UserCreate",
																},
															},
														},
													},
													{
														Names: []*ast.Ident{
															{
																Name: "requestUser",
															},
														},
														Type: &ast.StarExpr{
															X: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "models",
																},
																Sel: &ast.Ident{
																	Name: "User",
																},
															},
														},
													},
												},
											},
											Results: &ast.FieldList{
												List: []*ast.Field{
													{
														Type: &ast.StarExpr{
															X: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "models",
																},
																Sel: &ast.Ident{
																	Name: "User",
																},
															},
														},
													},
													{
														Type: &ast.Ident{
															Name: "error",
														},
													},
												},
											},
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "Update",
											},
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
													{
														Names: []*ast.Ident{
															{
																Name: "update",
															},
														},
														Type: &ast.StarExpr{
															X: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "models",
																},
																Sel: &ast.Ident{
																	Name: "UserUpdate",
																},
															},
														},
													},
													{
														Names: []*ast.Ident{
															{
																Name: "requestUser",
															},
														},
														Type: &ast.StarExpr{
															X: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "models",
																},
																Sel: &ast.Ident{
																	Name: "User",
																},
															},
														},
													},
												},
											},
											Results: &ast.FieldList{
												List: []*ast.Field{
													{
														Type: &ast.StarExpr{
															X: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "models",
																},
																Sel: &ast.Ident{
																	Name: "User",
																},
															},
														},
													},
													{
														Type: &ast.Ident{
															Name: "error",
														},
													},
												},
											},
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "Delete",
											},
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
													{
														Names: []*ast.Ident{
															{
																Name: "id",
															},
														},
														Type: &ast.SelectorExpr{
															X: &ast.Ident{
																Name: "models",
															},
															Sel: &ast.Ident{
																Name: "UUID",
															},
														},
													},
													{
														Names: []*ast.Ident{
															{
																Name: "requestUser",
															},
														},
														Type: &ast.StarExpr{
															X: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "models",
																},
																Sel: &ast.Ident{
																	Name: "User",
																},
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

func (i InterceptorInterfaceUser) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "domain", "interceptors", "user.go")
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
