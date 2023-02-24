package domain

import (
	"bytes"
	"github.com/018bf/creathor/internal/configs"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"
)

type UseCaseInterface struct {
	model *configs.ModelConfig
}

func NewUseCaseInterface(config *configs.ModelConfig) *UseCaseInterface {
	return &UseCaseInterface{model: config}
}

func (i UseCaseInterface) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "domain", "usecases", i.model.FileName())
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var structureExists bool
	var structure *ast.TypeSpec
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok && t.Name.String() == i.model.UseCaseTypeName() {
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
			Doc:    nil,
			TokPos: 0,
			Tok:    token.TYPE,
			Lparen: 0,
			Specs:  []ast.Spec{structure},
			Rparen: 0,
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

func (i UseCaseInterface) astInterface() *ast.TypeSpec {
	return &ast.TypeSpec{
		Name: &ast.Ident{
			Name: i.model.UseCaseTypeName(),
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
				},
			},
		},
	}
}
