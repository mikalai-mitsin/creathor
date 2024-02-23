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

	"github.com/018bf/creathor/internal/configs"
)

type InterceptorInterfaceAuth struct {
	project *configs.Project
}

func NewInterceptorInterfaceAuth(project *configs.Project) *InterceptorInterfaceAuth {
	return &InterceptorInterfaceAuth{project: project}
}

func (i InterceptorInterfaceAuth) file() *ast.File {
	return &ast.File{
		Name: &ast.Ident{
			Name: "grpc",
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
							Value: fmt.Sprintf(`"%s/internal/app/auth/models"`, i.project.Module),
						},
					},
					&ast.ImportSpec{
						Name: ast.NewIdent("userModels"),
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/app/user/models"`, i.project.Module),
						},
					},
				},
			},
			&ast.GenDecl{
				Doc: &ast.CommentGroup{
					List: []*ast.Comment{
						{
							Text: "//AuthInterceptor - domain layer interceptor interface",
						},
						{
							Text: "//go:generate mockgen -build_flags=-mod=mod -destination mock/auth_interceptor.go . AuthInterceptor",
						},
					},
				},
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{
							Name: "AuthInterceptor",
						},
						Type: &ast.InterfaceType{
							Methods: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											{
												Name: "Auth",
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
								},
							},
						},
					},
				},
			},
		},
	}
}

func (i InterceptorInterfaceAuth) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "pkg", "grpc", "interfaces.go")
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
