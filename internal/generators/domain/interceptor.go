package domain

import (
	"bytes"
	"fmt"
	"github.com/018bf/creathor/internal/configs"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"
)

type InterceptorInterface struct {
	model *configs.ModelConfig
}

func NewInterceptorInterface(config *configs.ModelConfig) *InterceptorInterface {
	return &InterceptorInterface{model: config}
}

func (i InterceptorInterface) file() *ast.File {
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
							Value: `"context"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/domain/models"`, i.model.Module),
						},
					},
				},
			},
		},
		Imports:  nil,
		Comments: nil,
	}
}

func (i InterceptorInterface) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "domain", "interceptors", i.model.FileName())
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = i.file()
	}
	var structureExists bool
	var structure *ast.TypeSpec
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok && t.Name.String() == i.model.InterceptorTypeName() {
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
						Text: fmt.Sprintf("//%s - domain layer interceptor interface", i.model.InterceptorTypeName()),
					},
					{
						Text: fmt.Sprintf("//go:generate mockgen -build_flags=-mod=mod -destination mock/%s %s/internal/domain/interceptors %s", i.model.FileName(), i.model.Module, i.model.InterceptorTypeName()),
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

func (i InterceptorInterface) astInterface() *ast.TypeSpec {
	requestUser := &ast.Field{
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
	}
	methods := []*ast.Field{
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
										Name: i.model.ModelName(),
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
										Name: i.model.FilterTypeName(),
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
											Name: i.model.ModelName(),
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
										Name: i.model.UpdateTypeName(),
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
										Name: i.model.ModelName(),
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
										Name: i.model.CreateTypeName(),
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
										Name: i.model.ModelName(),
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
	}
	if i.model.Auth {
		for _, method := range methods {
			ast.Inspect(method, func(node ast.Node) bool {
				if fun, ok := node.(*ast.FuncType); ok {
					fun.Params.List = append(fun.Params.List, requestUser)
					return false
				}
				return true
			})
		}
	}
	return &ast.TypeSpec{
		Name: &ast.Ident{
			Name: i.model.InterceptorTypeName(),
		},
		Type: &ast.InterfaceType{
			Methods: &ast.FieldList{
				List: methods,
			},
		},
	}
}
