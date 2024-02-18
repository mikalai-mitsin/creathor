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

	"github.com/018bf/creathor/internal/domain"
)

type UseCaseInterfaceCrud struct {
	domain *domain.Domain
}

func NewUseCaseInterfaceCrud(domain *domain.Domain) *UseCaseInterfaceCrud {
	return &UseCaseInterfaceCrud{domain: domain}
}

func (i UseCaseInterfaceCrud) file() *ast.File {
	file := &ast.File{
		Name: ast.NewIdent("interceptors"),
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
							Value: i.domain.ModelsImportPath(),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/pkg/uuid"`, i.domain.Module),
						},
					},
					&ast.ImportSpec{
						Name: ast.NewIdent("userModels"),
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/user/models"`, i.domain.Module),
						},
					},
				},
			},
		},
	}
	if i.domain.Auth {
		file.Decls = append(file.Decls, &ast.GenDecl{

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
		})
	}
	return file
}

func (i UseCaseInterfaceCrud) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", i.domain.DirName(), "interceptors", "interfaces.go")
	err := os.MkdirAll(path.Dir(filename), 0777)
	if err != nil {
		return err
	}
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = i.file()
	}
	var structureExists bool
	var structure *ast.TypeSpec
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok && t.Name.String() == i.domain.UseCase.Name {
			structure = t
			structureExists = true
			return false
		}
		return true
	})
	if structure == nil {
		structure = i.astInterface()
	}
	if !structureExists {
		gd := &ast.GenDecl{
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
			Tok:   token.TYPE,
			Specs: []ast.Spec{structure},
		}
		file.Decls = append(file.Decls, gd)
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

func (i UseCaseInterfaceCrud) astInterface() *ast.TypeSpec {
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
	return &ast.TypeSpec{
		Name: &ast.Ident{
			Name: i.domain.UseCase.Name,
		},
		Type: &ast.InterfaceType{
			Methods: &ast.FieldList{
				List: methods,
			},
		},
	}
}
