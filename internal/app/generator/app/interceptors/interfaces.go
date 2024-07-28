package interceptors

import (
	"bytes"
	"fmt"
	"github.com/mikalai-mitsin/creathor/internal/pkg/domain"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"
)

type InterceptorInterfaces struct {
	domain *domain.Domain
}

func NewInterceptorInterfaces(domain *domain.Domain) *InterceptorInterfaces {
	return &InterceptorInterfaces{domain: domain}
}

func (i InterceptorInterfaces) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "app", i.domain.DirName(), "interceptors", "interfaces.go")
	err := os.MkdirAll(path.Dir(filename), 0777)
	if err != nil {
		return err
	}
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = i.file()
	}
	authUseCaseExists := false
	appUseCaseExists := false
	loggerExists := false
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok {
			if t.Name.String() == i.domain.UseCase.Name {
				appUseCaseExists = true
			}
			if t.Name.String() == "AuthUseCase" {
				authUseCaseExists = true
			}
			if t.Name.String() == "Logger" {
				loggerExists = true
			}
			return true
		}
		return true
	})
	if !appUseCaseExists {
		file.Decls = append(file.Decls, i.appUseCaseInterface())
	}
	if !loggerExists {
		file.Decls = append(file.Decls, i.loggerInterface())
	}
	if i.domain.Auth && !authUseCaseExists {
		file.Decls = append(file.Decls, i.authUseCaseInterface())
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

func (i InterceptorInterfaces) file() *ast.File {
	file := &ast.File{
		Name: ast.NewIdent("interceptors"),
		Decls: []ast.Decl{
			i.imports(),
		},
	}
	return file
}

func (i InterceptorInterfaces) imports() *ast.GenDecl {
	imports := &ast.GenDecl{
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
					Value: i.domain.ModelsImportPath(),
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf(`"%s/internal/pkg/uuid"`, i.domain.Module),
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf(`"%s/internal/pkg/log"`, i.domain.Module),
				},
			},
		},
	}
	if i.domain.Auth {
		imports.Specs = append(imports.Specs, &ast.ImportSpec{
			Name: ast.NewIdent("userModels"),
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s/internal/app/user/models"`, i.domain.Module),
			},
		})
	}
	return imports
}

func (i InterceptorInterfaces) authUseCaseInterface() *ast.GenDecl {
	return &ast.GenDecl{
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{
					Text: "//AuthUseCase - domain layer interceptor interface",
				},
				{
					Text: "//go:generate mockgen -build_flags=-mod=mod -destination mock/auth.go . AuthUseCase",
				},
			},
		},
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: &ast.Ident{
					Name: "AuthUseCase",
				},
				Type: &ast.InterfaceType{
					Methods: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{
										Name: "GetUser",
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
										},
									},
									Results: &ast.FieldList{
										List: []*ast.Field{
											{
												Type: &ast.StarExpr{
													X: &ast.SelectorExpr{
														X:   ast.NewIdent("userModels"),
														Sel: ast.NewIdent("User"),
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
														Name: "user",
													},
												},
												Type: &ast.StarExpr{
													X: &ast.SelectorExpr{
														X:   ast.NewIdent("userModels"),
														Sel: ast.NewIdent("User"),
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
														Name: "user",
													},
												},
												Type: &ast.StarExpr{
													X: &ast.SelectorExpr{
														X:   ast.NewIdent("userModels"),
														Sel: ast.NewIdent("User"),
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
														Name: "object",
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
	}
}

func (i InterceptorInterfaces) appUseCaseInterface() *ast.GenDecl {
	methods := make([]*ast.Field, len(i.domain.UseCase.Methods))
	for i, method := range i.domain.UseCase.Methods {
		methods[i] = &ast.Field{
			Names: []*ast.Ident{
				{
					Name: method.Name,
				},
			},
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: method.Args,
				},
				Results: &ast.FieldList{
					List: method.Return,
				},
			},
		}
	}
	return &ast.GenDecl{
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{
					Text: fmt.Sprintf(
						"//%s - domain layer use case interface",
						i.domain.UseCase.Name,
					),
				},
				{
					Text: fmt.Sprintf(
						"//go:generate mockgen -build_flags=-mod=mod -destination mock/usecase.go . %s",
						i.domain.UseCase.Name,
					),
				},
			},
		},
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: &ast.Ident{
					Name: i.domain.UseCase.Name,
				},
				Type: &ast.InterfaceType{
					Methods: &ast.FieldList{
						List: methods,
					},
				},
			},
		},
	}
}

func (i InterceptorInterfaces) loggerInterface() *ast.GenDecl {
	return &ast.GenDecl{
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{
					Text: "//Logger - base logger interface",
				},
				{
					Text: "//go:generate mockgen -build_flags=-mod=mod -destination mock/logger.go . Logger",
				},
			},
		},
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
	}
}
