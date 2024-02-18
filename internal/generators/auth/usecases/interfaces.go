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

type RepositoryInterfaceAuth struct {
	project *configs.Project
}

func NewRepositoryInterfaceAuth(project *configs.Project) *RepositoryInterfaceAuth {
	return &RepositoryInterfaceAuth{project: project}
}

func (i RepositoryInterfaceAuth) file() *ast.File {
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
							Value: `"context"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/auth/models"`, i.project.Module),
						},
					},
					&ast.ImportSpec{
						Name: ast.NewIdent("userModels"),
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/user/models"`, i.project.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/pkg/uuid"`, i.project.Module),
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.TYPE,
				Doc: &ast.CommentGroup{
					List: []*ast.Comment{
						{
							Text: "//AuthRepository - domain layer repository interface",
						},
						{
							Text: "//go:generate mockgen -build_flags=-mod=mod -destination mock/auth_repository.go . AuthRepository",
						},
					},
				},
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{
							Name: "AuthRepository",
						},
						Type: &ast.InterfaceType{
							Methods: &ast.FieldList{
								List: []*ast.Field{
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
																Name: "user",
															},
														},
														Type: &ast.StarExpr{
															X: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "userModels",
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
																	Name: "TokenPair",
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
												Name: "Validate",
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
																Name: "token",
															},
														},
														Type: &ast.SelectorExpr{
															X: &ast.Ident{
																Name: "models",
															},
															Sel: &ast.Ident{
																Name: "Token",
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
									{
										Names: []*ast.Ident{
											{
												Name: "RefreshToken",
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
																Name: "token",
															},
														},
														Type: &ast.SelectorExpr{
															X: &ast.Ident{
																Name: "models",
															},
															Sel: &ast.Ident{
																Name: "Token",
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
																	Name: "TokenPair",
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
												Name: "GetSubject",
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
																Name: "token",
															},
														},
														Type: &ast.SelectorExpr{
															X: &ast.Ident{
																Name: "models",
															},
															Sel: &ast.Ident{
																Name: "Token",
															},
														},
													},
												},
											},
											Results: &ast.FieldList{
												List: []*ast.Field{
													{
														Type: &ast.Ident{
															Name: "string",
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
								},
							},
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.TYPE,
				Doc: &ast.CommentGroup{
					List: []*ast.Comment{
						{
							Text: "//UserRepository - domain layer repository interface",
						},
						{
							Text: "//go:generate mockgen -build_flags=-mod=mod -destination mock/user_repository.go . UserRepository",
						},
					},
				},
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{
							Name: "UserRepository",
						},
						Type: &ast.InterfaceType{
							Methods: &ast.FieldList{
								List: []*ast.Field{
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
														Type: ast.NewIdent("string"),
													},
												},
											},
											Results: &ast.FieldList{
												List: []*ast.Field{
													{
														Type: &ast.StarExpr{
															X: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "userModels",
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
														Type: ast.NewIdent("uuid.UUID"),
													},
												},
											},
											Results: &ast.FieldList{
												List: []*ast.Field{
													{
														Type: &ast.StarExpr{
															X: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "userModels",
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
								},
							},
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.TYPE,
				Doc: &ast.CommentGroup{
					List: []*ast.Comment{
						{
							Text: "//PermissionRepository - domain layer repository interface",
						},
						{
							Text: "//go:generate mockgen -build_flags=-mod=mod -destination mock/permission_repository.go . PermissionRepository",
						},
					},
				},
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{
							Name: "PermissionRepository",
						},
						Type: &ast.InterfaceType{
							Methods: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											{
												Name: "HasPermission",
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
																Name: "permission",
															},
														},
														Type: ast.NewIdent("userModels.PermissionID"),
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
																	Name: "userModels",
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
									{
										Names: []*ast.Ident{
											{
												Name: "HasObjectPermission",
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
																Name: "permission",
															},
														},
														Type: ast.NewIdent("userModels.PermissionID"),
													},
													{
														Names: []*ast.Ident{
															{
																Name: "user",
															},
														},
														Type: &ast.StarExpr{
															X: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "userModels",
																},
																Sel: &ast.Ident{
																	Name: "User",
																},
															},
														},
													},
													{
														Names: []*ast.Ident{
															{
																Name: "obj",
															},
														},
														Type: &ast.Ident{
															Name: "any",
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

func (i RepositoryInterfaceAuth) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "auth", "usecases", "interfaces.go")
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
