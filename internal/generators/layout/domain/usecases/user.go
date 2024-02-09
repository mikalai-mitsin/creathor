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

	"github.com/018bf/creathor/internal/configs"
)

type UseCaseInterfaceUser struct {
	project *configs.Project
}

// NewUseCaseInterfaceUser
// deprecated
func NewUseCaseInterfaceUser(project *configs.Project) *UseCaseInterfaceUser {
	return &UseCaseInterfaceUser{project: project}
}

func (i UseCaseInterfaceUser) file() *ast.File {
	return &ast.File{
		Name: &ast.Ident{
			Name: "usecases",
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
							Text: "//UserUseCase - domain layer interceptor interface",
						},
						{
							Text: "//go:generate mockgen -build_flags=-mod=mod -destination mock/user.go . UserUseCase",
						},
					},
				},
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{
							Name: "UserUseCase",
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
												Name: "GetByEmail",
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
																Name: "email",
															},
														},
														Type: &ast.Ident{
															Name: "string",
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

func (i UseCaseInterfaceUser) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "domain", "usecases", "user.go")
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
