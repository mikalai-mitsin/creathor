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

type RepositoryInterface struct {
	model *configs.ModelConfig
}

func NewRepositoryInterface(config *configs.ModelConfig) *RepositoryInterface {
	return &RepositoryInterface{model: config}
}

func (i RepositoryInterface) file() *ast.File {
	return &ast.File{
		Name: &ast.Ident{
			Name: "repositories",
		},
		Decls: []ast.Decl{
			&ast.GenDecl{
				Doc:    nil,
				TokPos: 0,
				Tok:    token.IMPORT,
				Lparen: 0,
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
				Rparen: 0,
			},
		},
		Imports:  nil,
		Comments: nil,
	}
}

func (i RepositoryInterface) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "domain", "repositories", i.model.FileName())
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = i.file()
	}
	var structureExists bool
	var structure *ast.TypeSpec
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok && t.Name.String() == i.model.RepositoryTypeName() {
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
						Text: fmt.Sprintf("//%s - domain layer repository interface", i.model.RepositoryTypeName()),
					},
					{
						Text: fmt.Sprintf("//go:generate mockgen -build_flags=-mod=mod -destination mock/%s %s/internal/domain/repositories %s", i.model.FileName(), i.model.Module, i.model.RepositoryTypeName()),
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

func (i RepositoryInterface) astInterface() *ast.TypeSpec {
	return &ast.TypeSpec{
		Name: &ast.Ident{
			Name: i.model.RepositoryTypeName(),
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
								Name: "Count",
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
													Name: i.model.ModelName(),
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
													Name: i.model.ModelName(),
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
	}
}
