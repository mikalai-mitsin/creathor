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

type UseCaseInterfaceAuth struct {
	project *configs.Project
}

func NewUseCaseInterfaceAuth(project *configs.Project) *UseCaseInterfaceAuth {
	return &UseCaseInterfaceAuth{project: project}
}

func (i UseCaseInterfaceAuth) file() *ast.File {
	return &ast.File{
		Name: ast.NewIdent("grpc"),
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
							Value: fmt.Sprintf(`"%s/internal/app/auth/entities"`, i.project.Module),
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
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: ast.NewIdent("Logger"),
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
				Doc: &ast.CommentGroup{
					List: []*ast.Comment{
						{
							Text: "//AuthUseCase - domain layer usecase interface",
						},
						{
							Text: "//go:generate mockgen -build_flags=-mod=mod -destination mock/auth_usecase.go . AuthUseCase",
						},
					},
				},
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: ast.NewIdent("AuthUseCase"),
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
															ast.NewIdent("token"),
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
														Type: &ast.StarExpr{
															X: &ast.SelectorExpr{
																X:   ast.NewIdent("userEntities"),
																Sel: ast.NewIdent("User"),
															},
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

func (i UseCaseInterfaceAuth) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "pkg", "grpc", "auth_interfaces.go")
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
