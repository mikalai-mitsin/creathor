package generators

import (
	"bytes"
	"fmt"
	"github.com/018bf/creathor/internal/models"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
)

type Repository struct {
	Path   string
	Name   string
	Model  *models.ModelConfig
	Params []*Param
}

func (r Repository) AstDTOStruct() *ast.TypeSpec {
	structure := &ast.TypeSpec{
		Name: &ast.Ident{
			Name: r.Model.PostgresDTOTypeName(),
		},
		Type: &ast.StructType{
			Fields: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
								Name: "ID",
							},
						},
						Type: &ast.Ident{
							Name: "string",
						},
						Tag: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "`db:\"id,omitempty\"`",
						},
					},
					{
						Names: []*ast.Ident{
							{
								Name: "UpdatedAt",
							},
						},
						Type: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "time",
							},
							Sel: &ast.Ident{
								Name: "Time",
							},
						},
						Tag: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "`db:\"updated_at,omitempty\"`",
						},
					},
					{
						Names: []*ast.Ident{
							{
								Name: "CreatedAt",
							},
						},
						Type: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "time",
							},
							Sel: &ast.Ident{
								Name: "Time",
							},
						},
						Tag: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "`db:\"created_at,omitempty\"`",
						},
					},
				},
			},
		},
	}
	for _, param := range r.Model.Params {
		ast.Inspect(structure, func(node ast.Node) bool {
			if st, ok := node.(*ast.StructType); ok && st.Fields != nil {
				for _, field := range st.Fields.List {
					for _, fieldName := range field.Names {
						if fieldName.Name == param.GetName() {
							return false
						}
					}
				}
				st.Fields.List = append(st.Fields.List, &ast.Field{
					Doc:   nil,
					Names: []*ast.Ident{ast.NewIdent(param.GetName())},
					Type:  ast.NewIdent(param.PostgresDTOType()),
					Tag: &ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf("`db:\"%s\"`", param.Tag()),
					},
					Comment: nil,
				})
				return false
			}
			return true
		})
	}
	return structure
}

func (r Repository) SyncDTOStruct() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, r.Path, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var structureExists bool
	var structure *ast.TypeSpec
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok && t.Name.String() == r.Model.PostgresDTOTypeName() {
			structure = t
			structureExists = true
			return false
		}
		return true
	})
	if structure == nil {
		structure = r.AstDTOStruct()
	}
	for _, param := range r.Model.Params {
		ast.Inspect(structure, func(node ast.Node) bool {
			if st, ok := node.(*ast.StructType); ok && st.Fields != nil {
				for _, field := range st.Fields.List {
					for _, fieldName := range field.Names {
						if fieldName.Name == param.GetName() {
							return false
						}
					}
				}
				st.Fields.List = append(st.Fields.List, &ast.Field{
					Doc:   nil,
					Names: []*ast.Ident{ast.NewIdent(param.GetName())},
					Type:  ast.NewIdent(param.PostgresDTOType()),
					Tag: &ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf("`db:\"%s\"`", param.Tag()),
					},
					Comment: nil,
				})
				return false
			}
			return true
		})
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
	if err := os.WriteFile(r.Path, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (r Repository) AstDTOConstructor() *ast.FuncDecl {
	dto := &ast.CompositeLit{
		Type: &ast.Ident{
			Name: r.Model.PostgresDTOTypeName(),
		},
		Elts: []ast.Expr{
			&ast.KeyValueExpr{
				Key: &ast.Ident{
					Name: "ID",
				},
				Value: &ast.CallExpr{
					Fun: &ast.Ident{
						Name: "string",
					},
					Args: []ast.Expr{
						&ast.SelectorExpr{
							X: &ast.Ident{
								Name: r.Model.Variable(),
							},
							Sel: &ast.Ident{
								Name: "ID",
							},
						},
					},
				},
			},
			&ast.KeyValueExpr{
				Key: &ast.Ident{
					Name: "UpdatedAt",
				},
				Value: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: r.Model.Variable(),
					},
					Sel: &ast.Ident{
						Name: "UpdatedAt",
					},
				},
			},
			&ast.KeyValueExpr{
				Key: &ast.Ident{
					Name: "CreatedAt",
				},
				Value: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: r.Model.Variable(),
					},
					Sel: &ast.Ident{
						Name: "CreatedAt",
					},
				},
			},
		},
	}
	for _, param := range r.Model.Params {
		elt := &ast.KeyValueExpr{
			Key: &ast.Ident{
				Name: param.GetName(),
			},
			Value: nil,
		}
		if param.IsSlice() {
			elt.Value = &ast.CompositeLit{
				Type: ast.NewIdent(param.PostgresDTOType()),
			}
		} else {
			if param.PostgresDTOType() == param.Type {
				elt.Value = &ast.SelectorExpr{
					X: &ast.Ident{
						Name: r.Model.Variable(),
					},
					Sel: &ast.Ident{
						Name: param.GetName(),
					},
				}
			} else {
				elt.Value = &ast.CallExpr{
					Fun: &ast.Ident{
						Name: param.PostgresDTOType(),
					},
					Args: []ast.Expr{
						&ast.SelectorExpr{
							X: &ast.Ident{
								Name: r.Model.Variable(),
							},
							Sel: &ast.Ident{
								Name: param.GetName(),
							},
						},
					},
				}
			}
		}
		dto.Elts = append(dto.Elts, elt)
	}
	constructor := &ast.FuncDecl{
		Name: &ast.Ident{
			Name: fmt.Sprintf("New%sFromModel", r.Model.PostgresDTOTypeName()),
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
								Name: r.Model.Variable(),
							},
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "models",
								},
								Sel: &ast.Ident{
									Name: r.Model.ModelName(),
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
							X: &ast.Ident{
								Name: r.Model.PostgresDTOTypeName(),
							},
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "dto",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X:  dto,
						},
					},
				},
			},
		},
	}
	for _, param := range r.Model.Params {
		if !param.IsSlice() {
			continue
		}
		var valueToAppend ast.Expr
		if param.SliceType() == param.PostgresDTOSliceType() {
			valueToAppend = ast.NewIdent("param")
		} else {
			valueToAppend = &ast.CallExpr{
				Fun: &ast.Ident{
					Name: param.PostgresDTOSliceType(),
				},
				Args: []ast.Expr{
					&ast.Ident{
						Name: "param",
					},
				},
			}
		}
		rang := &ast.RangeStmt{
			Key: &ast.Ident{
				Name: "_",
			},
			Value: &ast.Ident{
				Name: "param",
			},
			Tok: token.DEFINE,
			X: &ast.SelectorExpr{
				X: &ast.Ident{
					Name: r.Model.Variable(),
				},
				Sel: &ast.Ident{
					Name: param.GetName(),
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.SelectorExpr{
								X: &ast.Ident{
									Name: "dto",
								},
								Sel: &ast.Ident{
									Name: param.GetName(),
								},
							},
						},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.Ident{
									Name: "append",
								},
								Args: []ast.Expr{
									&ast.SelectorExpr{
										X: &ast.Ident{
											Name: "dto",
										},
										Sel: &ast.Ident{
											Name: param.GetName(),
										},
									},
									valueToAppend,
								},
							},
						},
					},
				},
			},
		}
		constructor.Body.List = append(constructor.Body.List, rang)
	}
	constructor.Body.List = append(
		constructor.Body.List,
		&ast.ReturnStmt{
			Results: []ast.Expr{
				ast.NewIdent("dto"),
			},
		},
	)
	return constructor
}

func (r Repository) SyncDTOConstructor() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, r.Path, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var structureConstructorExists bool
	var structureConstructor *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == fmt.Sprintf("New%sFromModel", r.Model.PostgresDTOTypeName()) {
			structureConstructorExists = true
			structureConstructor = t
			return false
		}
		return true
	})
	if structureConstructor == nil {
		structureConstructor = r.AstDTOConstructor()
	}
	for _, param := range r.Model.Params {
		param := param
		ast.Inspect(structureConstructor, func(node ast.Node) bool {
			if cl, ok := node.(*ast.CompositeLit); ok {
				if i, ok := cl.Type.(*ast.Ident); ok && i.String() == r.Model.PostgresDTOTypeName() {
					_ = i
					for _, elt := range cl.Elts {
						elt := elt
						if kv, ok := elt.(*ast.KeyValueExpr); ok {
							if key, ok := kv.Key.(*ast.Ident); ok && key.String() == param.GetName() {
								return false
							}
						}
					}
					elt := &ast.KeyValueExpr{
						Key: &ast.Ident{
							Name: param.GetName(),
						},
						Value: nil,
					}
					if param.IsSlice() {
						elt.Value = &ast.CompositeLit{
							Type: ast.NewIdent(param.PostgresDTOType()),
						}
					} else {
						if param.PostgresDTOType() == param.Type {
							elt.Value = &ast.SelectorExpr{
								X: &ast.Ident{
									Name: r.Model.Variable(),
								},
								Sel: &ast.Ident{
									Name: param.GetName(),
								},
							}
						} else {
							elt.Value = &ast.CallExpr{
								Fun: &ast.Ident{
									Name: param.PostgresDTOType(),
								},
								Args: []ast.Expr{
									&ast.SelectorExpr{
										X: &ast.Ident{
											Name: r.Model.Variable(),
										},
										Sel: &ast.Ident{
											Name: param.GetName(),
										},
									},
								},
							}
						}
					}
					cl.Elts = append(cl.Elts, elt)
					return false
				}
			}
			return true
		})
	}
	// TODO: add range sync
	if !structureConstructorExists {
		file.Decls = append(file.Decls, structureConstructor)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(r.Path, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (r Repository) AstDTOToModel() *ast.FuncDecl {
	model := &ast.CompositeLit{
		Type: &ast.SelectorExpr{
			X: &ast.Ident{
				Name: "models",
			},
			Sel: &ast.Ident{
				Name: r.Model.ModelName(),
			},
		},
		Elts: []ast.Expr{
			&ast.KeyValueExpr{
				Key: &ast.Ident{
					Name: "ID",
				},
				Value: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: "models",
						},
						Sel: &ast.Ident{
							Name: "UUID",
						},
					},
					Args: []ast.Expr{
						&ast.SelectorExpr{
							X: &ast.Ident{
								Name: "dto",
							},
							Sel: &ast.Ident{
								Name: "ID",
							},
						},
					},
				},
			},
			&ast.KeyValueExpr{
				Key: &ast.Ident{
					Name: "UpdatedAt",
				},
				Value: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "dto",
					},
					Sel: &ast.Ident{
						Name: "UpdatedAt",
					},
				},
			},
			&ast.KeyValueExpr{
				Key: &ast.Ident{
					Name: "CreatedAt",
				},
				Value: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "dto",
					},
					Sel: &ast.Ident{
						Name: "CreatedAt",
					},
				},
			},
		},
	}
	for _, param := range r.Model.Params {
		par := &ast.KeyValueExpr{
			Key: &ast.Ident{
				Name: param.GetName(),
			},
			Value: nil,
		}
		if param.IsSlice() {
			par.Value = &ast.CompositeLit{
				Type: &ast.ArrayType{
					Elt: &ast.Ident{
						Name: param.SliceType(),
					},
				},
			}
		} else {
			if param.PostgresDTOType() == param.Type {
				par.Value = &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "dto",
					},
					Sel: &ast.Ident{
						Name: param.GetName(),
					},
				}
			} else {
				par.Value = &ast.CallExpr{
					Fun: &ast.Ident{
						Name: param.Type,
					},
					Args: []ast.Expr{
						&ast.SelectorExpr{
							X: &ast.Ident{
								Name: "dto",
							},
							Sel: &ast.Ident{
								Name: param.GetName(),
							},
						},
					},
				}
			}
		}
		model.Elts = append(model.Elts, par)
	}
	method := &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{
							Name: "dto",
						},
					},
					Type: &ast.StarExpr{
						X: &ast.Ident{
							Name: r.Model.PostgresDTOTypeName(),
						},
					},
				},
			},
		},
		Name: &ast.Ident{
			Name: "ToModel",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "models",
								},
								Sel: &ast.Ident{
									Name: r.Model.ModelName(),
								},
							},
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "model",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X:  model,
						},
					},
				},
			},
		},
	}
	for _, param := range r.Model.Params {
		if !param.IsSlice() {
			continue
		}
		var valueToAppend ast.Expr
		if param.SliceType() == param.PostgresDTOSliceType() {
			valueToAppend = ast.NewIdent("param")
		} else {
			valueToAppend = &ast.CallExpr{
				Fun: &ast.Ident{
					Name: param.SliceType(),
				},
				Args: []ast.Expr{
					&ast.Ident{
						Name: "param",
					},
				},
			}
		}
		rang := &ast.RangeStmt{
			Key: &ast.Ident{
				Name: "_",
			},
			Value: &ast.Ident{
				Name: "param",
			},
			Tok: token.DEFINE,
			X: &ast.SelectorExpr{
				X: &ast.Ident{
					Name: "dto",
				},
				Sel: &ast.Ident{
					Name: param.GetName(),
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.SelectorExpr{
								X: &ast.Ident{
									Name: "model",
								},
								Sel: &ast.Ident{
									Name: param.GetName(),
								},
							},
						},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.Ident{
									Name: "append",
								},
								Args: []ast.Expr{
									&ast.SelectorExpr{
										X: &ast.Ident{
											Name: "model",
										},
										Sel: &ast.Ident{
											Name: param.GetName(),
										},
									},
									valueToAppend,
								},
							},
						},
					},
				},
			},
		}
		method.Body.List = append(method.Body.List, rang)
	}
	method.Body.List = append(
		method.Body.List,
		&ast.ReturnStmt{
			Results: []ast.Expr{
				&ast.Ident{
					Name: "model",
				},
			},
		},
	)
	return method
}

func (r Repository) SyncDTOToModel() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, r.Path, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExists bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == "ToModel" {
			methodExists = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = r.AstDTOToModel()
	}
	// TODO: add range sync
	if !methodExists {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(r.Path, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (r Repository) AstStruct() *ast.TypeSpec {
	structure := &ast.TypeSpec{
		Doc:        nil,
		Name:       ast.NewIdent(r.Name),
		TypeParams: nil,
		Assign:     0,
		Type: &ast.StructType{
			Struct: 0,
			Fields: &ast.FieldList{
				Opening: 0,
				List:    nil,
				Closing: 0,
			},
			Incomplete: false,
		},
		Comment: nil,
	}
	for _, param := range r.Params {
		ast.Inspect(structure, func(node ast.Node) bool {
			if st, ok := node.(*ast.StructType); ok && st.Fields != nil {
				for _, field := range st.Fields.List {
					for _, fieldName := range field.Names {
						if fieldName.Name == param.GetPrivateName() {
							return false
						}
					}
				}
				st.Fields.List = append(st.Fields.List, &ast.Field{
					Doc:     nil,
					Names:   []*ast.Ident{ast.NewIdent(param.GetPrivateName())},
					Type:    ast.NewIdent(param.Type),
					Tag:     nil,
					Comment: nil,
				})
				return false
			}
			return true
		})
	}
	return structure
}

func (r Repository) SyncStruct() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, r.Path, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var structureExists bool
	var structure *ast.TypeSpec
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok && t.Name.String() == r.Name {
			structure = t
			structureExists = true
			return false
		}
		return true
	})
	if structure == nil {
		structure = r.AstStruct()
	}
	for _, param := range r.Params {
		ast.Inspect(structure, func(node ast.Node) bool {
			if st, ok := node.(*ast.StructType); ok && st.Fields != nil {
				for _, field := range st.Fields.List {
					for _, fieldName := range field.Names {
						if fieldName.Name == param.GetPrivateName() {
							return false
						}
					}
				}
				st.Fields.List = append(st.Fields.List, &ast.Field{
					Doc:     nil,
					Names:   []*ast.Ident{ast.NewIdent(param.GetPrivateName())},
					Type:    ast.NewIdent(param.Type),
					Tag:     nil,
					Comment: nil,
				})
				return false
			}
			return true
		})
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
	if err := os.WriteFile(r.Path, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (r Repository) AstConstructor() *ast.FuncDecl {
	var args []*ast.Field
	cl := &ast.CompositeLit{
		Type:       ast.NewIdent(r.Name),
		Lbrace:     0,
		Elts:       nil,
		Rbrace:     0,
		Incomplete: false,
	}
	for _, param := range r.Params {
		args = append(
			args,
			&ast.Field{
				Doc:     nil,
				Names:   []*ast.Ident{ast.NewIdent(param.GetPrivateName())},
				Type:    ast.NewIdent(param.Type),
				Tag:     nil,
				Comment: nil,
			},
		)
		cl.Elts = append(cl.Elts, &ast.KeyValueExpr{
			Key:   ast.NewIdent(param.GetPrivateName()),
			Colon: 0,
			Value: ast.NewIdent(param.GetPrivateName()),
		})
	}
	constructor := &ast.FuncDecl{
		Doc:  nil,
		Recv: nil,
		Name: ast.NewIdent(fmt.Sprintf("New%s", r.Name)),
		Type: &ast.FuncType{
			Func:       0,
			TypeParams: nil,
			Params: &ast.FieldList{
				Opening: 0,
				List:    args,
				Closing: 0,
			},
			Results: &ast.FieldList{
				Opening: 0,
				List: []*ast.Field{
					{
						Doc:   nil,
						Names: nil,
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("repositories"),
							Sel: ast.NewIdent(r.Name),
						},
						Tag:     nil,
						Comment: nil,
					},
				},
				Closing: 0,
			},
		},
		Body: &ast.BlockStmt{
			Lbrace: 0,
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Return: 0,
					Results: []ast.Expr{
						&ast.UnaryExpr{
							OpPos: 0,
							Op:    token.AND,
							X:     cl,
						},
					},
				},
			},
			Rbrace: 0,
		},
	}
	return constructor
}

func (r Repository) SyncConstructor() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, r.Path, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var structureConstructorExists bool
	var structureConstructor *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == fmt.Sprintf("New%s", r.Name) {
			structureConstructorExists = true
			structureConstructor = t
			return false
		}
		return true
	})
	if structureConstructor == nil {
		structureConstructor = r.AstConstructor()
	}
	for _, param := range r.Params {
		param := param
		var argExists bool
		for _, arg := range structureConstructor.Type.Params.List {
			for _, fieldName := range arg.Names {
				if fieldName.Name == param.GetPrivateName() {
					argExists = true
				}
			}
		}
		ast.Inspect(structureConstructor.Body, func(node ast.Node) bool {
			if cl, ok := node.(*ast.CompositeLit); ok {
				if t, ok := cl.Type.(*ast.Ident); ok && t.String() == r.Name {
					for _, elt := range cl.Elts {
						if kv, ok := elt.(*ast.KeyValueExpr); ok {
							if key, ok := kv.Key.(*ast.Ident); ok && key.String() == param.GetPrivateName() {
								return false
							}
						}
					}
					cl.Elts = append(cl.Elts, &ast.KeyValueExpr{
						Key:   ast.NewIdent(param.GetPrivateName()),
						Colon: 0,
						Value: ast.NewIdent(param.GetPrivateName()),
					})
					return false
				}
			}
			return true
		})
		if !argExists {
			structureConstructor.Type.Params.List = append(
				structureConstructor.Type.Params.List,
				&ast.Field{
					Doc:     nil,
					Names:   []*ast.Ident{ast.NewIdent(param.GetPrivateName())},
					Type:    ast.NewIdent(param.Type),
					Tag:     nil,
					Comment: nil,
				},
			)
		}
	}
	if !structureConstructorExists {
		file.Decls = append(file.Decls, structureConstructor)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(r.Path, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (r Repository) AstCreateMethod() *ast.FuncDecl {
	columns := []ast.Expr{
		&ast.BasicLit{
			Kind:  token.STRING,
			Value: "\"updated_at\"",
		},
		&ast.BasicLit{
			Kind:  token.STRING,
			Value: "\"created_at\"",
		},
	}
	values := []ast.Expr{
		&ast.SelectorExpr{
			X:   ast.NewIdent("dto"),
			Sel: ast.NewIdent("UpdatedAt"),
		},
		&ast.SelectorExpr{
			X:   ast.NewIdent("dto"),
			Sel: ast.NewIdent("CreatedAt"),
		},
	}
	for _, param := range r.Model.Params {
		columns = append(columns, &ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf("\"%s\"", param.Tag()),
		})
		values = append(values, &ast.SelectorExpr{
			X:   ast.NewIdent("dto"),
			Sel: ast.NewIdent(param.GetName()),
		})
	}
	fun := &ast.FuncDecl{
		Doc: nil,
		Recv: &ast.FieldList{
			Opening: 0,
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("r"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(r.Name),
					},
				},
			},
			Closing: 0,
		},
		Name: ast.NewIdent("Create"),
		Type: &ast.FuncType{
			Func:       0,
			TypeParams: nil,
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("ctx")},
						Type:  ast.NewIdent("context.Context"),
					},
					{
						Names: []*ast.Ident{ast.NewIdent(r.Model.Variable())},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("models"),
								Sel: ast.NewIdent(r.Model.ModelName()),
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				// Setup timeout
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("ctx"),
						ast.NewIdent("cancel"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("context"),
								Sel: ast.NewIdent("WithTimeout"),
							},
							Args: []ast.Expr{
								ast.NewIdent("ctx"),
								&ast.SelectorExpr{
									X:   ast.NewIdent("time"),
									Sel: ast.NewIdent("Second"),
								},
							},
						},
					},
				},
				// Defer cancel
				&ast.DeferStmt{
					Call: &ast.CallExpr{
						Fun: ast.NewIdent("cancel"),
					},
				},
				// Create DTO from model
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "dto",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: ast.NewIdent(fmt.Sprintf("New%sDTOFromModel", r.Model.ModelName())),
							Args: []ast.Expr{
								ast.NewIdent(r.Model.Variable()),
							},
						},
					},
				},
				// Create sq
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("q"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.CallExpr{
													Fun: &ast.SelectorExpr{
														X:   ast.NewIdent("sq"),
														Sel: ast.NewIdent("Insert"),
													},
													Args: []ast.Expr{
														&ast.BasicLit{
															Kind:  token.STRING,
															Value: fmt.Sprintf("\"public.%s\"", r.Model.TableName()),
														},
													},
												},
												Sel: ast.NewIdent("Columns"),
											},
											Args: columns,
										},
										Sel: ast.NewIdent("Values"),
									},
									Args: values,
								},
								Sel: ast.NewIdent("Suffix"),
							},
							Args: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: "\"RETURNING id\"",
								},
							},
						},
					},
				},
				// Build query from sq
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("query"),
						ast.NewIdent("args"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("q"),
										Sel: ast.NewIdent("PlaceholderFormat"),
									},
									Args: []ast.Expr{
										&ast.SelectorExpr{
											X:   ast.NewIdent("sq"),
											Sel: ast.NewIdent("Dollar"),
										},
									},
								},
								Sel: ast.NewIdent("MustSql"),
							},
						},
					},
				},
				// Run query at DB
				&ast.IfStmt{
					Init: &ast.AssignStmt{
						Lhs: []ast.Expr{
							ast.NewIdent("err"),
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.SelectorExpr{
												X:   ast.NewIdent("r"),
												Sel: ast.NewIdent("database"),
											},
											Sel: ast.NewIdent("QueryRowxContext"),
										},
										Args: []ast.Expr{
											ast.NewIdent("ctx"),
											ast.NewIdent("query"),
											ast.NewIdent("args"),
										},
										Ellipsis: 3467,
									},
									Sel: ast.NewIdent("StructScan"),
								},
								Args: []ast.Expr{
									ast.NewIdent("dto"),
								},
							},
						},
					},
					Cond: &ast.BinaryExpr{
						X:  ast.NewIdent("err"),
						Op: token.NEQ,
						Y:  ast.NewIdent("nil"),
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									ast.NewIdent("e"),
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   ast.NewIdent("errs"),
											Sel: ast.NewIdent("FromPostgresError"),
										},
										Args: []ast.Expr{
											ast.NewIdent("err"),
										},
									},
								},
							},
							&ast.ReturnStmt{
								Results: []ast.Expr{
									ast.NewIdent("e"),
								},
							},
						},
					},
				},
				// Set model ID from DTO
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X:   ast.NewIdent(r.Model.Variable()),
							Sel: ast.NewIdent("ID"),
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("models"),
								Sel: ast.NewIdent("UUID"),
							},
							Args: []ast.Expr{
								&ast.SelectorExpr{
									X:   ast.NewIdent("dto"),
									Sel: ast.NewIdent("ID"),
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
	return fun
}

func (r Repository) SyncCreateMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, r.Path, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == "Create" {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = r.AstCreateMethod()
	}
	for _, param := range r.Model.Params {
		param := param
		ast.Inspect(method, func(node ast.Node) bool {
			if call, ok := node.(*ast.CallExpr); ok {
				if fun, ok := call.Fun.(*ast.SelectorExpr); ok && fun.Sel.String() == "Columns" {
					for _, arg := range call.Args {
						arg := arg
						if bl, ok := arg.(*ast.BasicLit); ok && bl.Value == fmt.Sprintf("\"%s\"", param.Tag()) {
							return false
						}
					}
					call.Args = append(call.Args, &ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf("\"%s\"", param.Tag()),
					})
					return false
				}
			}
			return true
		})
		ast.Inspect(method, func(node ast.Node) bool {
			if call, ok := node.(*ast.CallExpr); ok {
				if fun, ok := call.Fun.(*ast.SelectorExpr); ok && fun.Sel.String() == "Values" {
					for _, arg := range call.Args {
						arg := arg
						if bl, ok := arg.(*ast.SelectorExpr); ok && bl.Sel.String() == param.GetName() {
							return false
						}
					}
					call.Args = append(call.Args, &ast.SelectorExpr{
						X:   ast.NewIdent("dto"),
						Sel: ast.NewIdent(param.GetName()),
					})
					return false
				}
			}
			return true
		})
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}

	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(r.Path, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (r Repository) AstListMethod() *ast.FuncDecl {
	tableName := r.Model.TableName()
	columns := []ast.Expr{
		&ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf("\"%s.id\"", tableName),
		},
		&ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf("\"%s.updated_at\"", tableName),
		},
		&ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf("\"%s.created_at\"", tableName),
		},
	}
	for _, param := range r.Model.Params {
		columns = append(
			columns,
			&ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf("\"%s.%s\"", tableName, param.Tag()),
			},
		)
	}
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{
							Name: "r",
						},
					},
					Type: &ast.StarExpr{
						X: &ast.Ident{
							Name: r.Name,
						},
					},
				},
			},
		},
		Name: &ast.Ident{
			Name: "List",
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
									Name: r.Model.FilterTypeName(),
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
										Name: r.Model.ModelName(),
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
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "ctx",
						},
						&ast.Ident{
							Name: "cancel",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "context",
								},
								Sel: &ast.Ident{
									Name: "WithTimeout",
								},
							},
							Args: []ast.Expr{
								&ast.Ident{
									Name: "ctx",
								},
								&ast.SelectorExpr{
									X: &ast.Ident{
										Name: "time",
									},
									Sel: &ast.Ident{
										Name: "Second",
									},
								},
							},
						},
					},
				},
				&ast.DeferStmt{
					Call: &ast.CallExpr{
						Fun: &ast.Ident{
							Name: "cancel",
						},
					},
				},
				&ast.DeclStmt{
					Decl: &ast.GenDecl{
						Tok: token.VAR,
						Specs: []ast.Spec{
							&ast.ValueSpec{
								Names: []*ast.Ident{
									{
										Name: "dto",
									},
								},
								Type: &ast.Ident{
									Name: r.Model.PostgresDTOListTypeName(),
								},
							},
						},
					},
				},
				&ast.DeclStmt{
					Decl: &ast.GenDecl{
						Tok: token.CONST,
						Specs: []ast.Spec{
							&ast.ValueSpec{
								Names: []*ast.Ident{
									{
										Name: "pageSize",
									},
								},
								Values: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.Ident{
											Name: "uint64",
										},
										Args: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.INT,
												Value: "10",
											},
										},
									},
								},
							},
						},
					},
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "filter",
							},
							Sel: &ast.Ident{
								Name: "PageSize",
							},
						},
						Op: token.EQL,
						Y: &ast.Ident{
							Name: "nil",
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.SelectorExpr{
										X: &ast.Ident{
											Name: "filter",
										},
										Sel: &ast.Ident{
											Name: "PageSize",
										},
									},
								},
								Tok: token.ASSIGN,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "utils",
											},
											Sel: &ast.Ident{
												Name: "Pointer",
											},
										},
										Args: []ast.Expr{
											&ast.Ident{
												Name: "pageSize",
											},
										},
									},
								},
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "q",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "sq",
												},
												Sel: &ast.Ident{
													Name: "Select",
												},
											},
											Args: columns,
										},
										Sel: &ast.Ident{
											Name: "From",
										},
									},
									Args: []ast.Expr{
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: fmt.Sprintf("\"public.%s\"", tableName),
										},
									},
								},
								Sel: &ast.Ident{
									Name: "Limit",
								},
							},
							Args: []ast.Expr{
								&ast.Ident{
									Name: "pageSize",
								},
							},
						},
					},
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "filter",
							},
							Sel: &ast.Ident{
								Name: "Search",
							},
						},
						Op: token.NEQ,
						Y: &ast.Ident{
							Name: "nil",
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.Ident{
										Name: "q",
									},
								},
								Tok: token.ASSIGN,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "q",
											},
											Sel: &ast.Ident{
												Name: "Where",
											},
										},
										Args: []ast.Expr{
											&ast.CompositeLit{
												Type: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "postgresql",
													},
													Sel: &ast.Ident{
														Name: "Search",
													},
												},
												Elts: []ast.Expr{
													&ast.KeyValueExpr{
														Key: &ast.Ident{
															Name: "Lang",
														},
														Value: &ast.BasicLit{
															Kind:  token.STRING,
															Value: "\"english\"",
														},
													},
													&ast.KeyValueExpr{
														Key: &ast.Ident{
															Name: "Query",
														},
														Value: &ast.StarExpr{
															X: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "filter",
																},
																Sel: &ast.Ident{
																	Name: "Search",
																},
															},
														},
													},
													&ast.KeyValueExpr{
														Key: &ast.Ident{
															Name: "Fields",
														},
														Value: &ast.CompositeLit{
															Type: &ast.ArrayType{
																Elt: &ast.Ident{
																	Name: "string",
																},
															},
															Elts: []ast.Expr{
																&ast.BasicLit{
																	Kind:  token.STRING,
																	Value: "\"name\"",
																},
																&ast.BasicLit{
																	Kind:  token.STRING,
																	Value: "\"subtitle\"",
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
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X: &ast.BinaryExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "filter",
								},
								Sel: &ast.Ident{
									Name: "PageNumber",
								},
							},
							Op: token.NEQ,
							Y: &ast.Ident{
								Name: "nil",
							},
						},
						Op: token.LAND,
						Y: &ast.BinaryExpr{
							X: &ast.StarExpr{
								X: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "filter",
									},
									Sel: &ast.Ident{
										Name: "PageNumber",
									},
								},
							},
							Op: token.GTR,
							Y: &ast.BasicLit{
								Kind:  token.INT,
								Value: "1",
							},
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.Ident{
										Name: "q",
									},
								},
								Tok: token.ASSIGN,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "q",
											},
											Sel: &ast.Ident{
												Name: "Offset",
											},
										},
										Args: []ast.Expr{
											&ast.BinaryExpr{
												X: &ast.ParenExpr{
													X: &ast.BinaryExpr{
														X: &ast.StarExpr{
															X: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "filter",
																},
																Sel: &ast.Ident{
																	Name: "PageNumber",
																},
															},
														},
														Op: token.SUB,
														Y: &ast.BasicLit{
															Kind:  token.INT,
															Value: "1",
														},
													},
												},
												Op: token.MUL,
												Y: &ast.StarExpr{
													X: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "filter",
														},
														Sel: &ast.Ident{
															Name: "PageSize",
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
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "q",
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "q",
								},
								Sel: &ast.Ident{
									Name: "Limit",
								},
							},
							Args: []ast.Expr{
								&ast.StarExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "filter",
										},
										Sel: &ast.Ident{
											Name: "PageSize",
										},
									},
								},
							},
						},
					},
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X: &ast.CallExpr{
							Fun: &ast.Ident{
								Name: "len",
							},
							Args: []ast.Expr{
								&ast.SelectorExpr{
									X: &ast.Ident{
										Name: "filter",
									},
									Sel: &ast.Ident{
										Name: "OrderBy",
									},
								},
							},
						},
						Op: token.GTR,
						Y: &ast.BasicLit{
							Kind:  token.INT,
							Value: "0",
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.Ident{
										Name: "q",
									},
								},
								Tok: token.ASSIGN,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "q",
											},
											Sel: &ast.Ident{
												Name: "OrderBy",
											},
										},
										Args: []ast.Expr{
											&ast.SelectorExpr{
												X: &ast.Ident{
													Name: "filter",
												},
												Sel: &ast.Ident{
													Name: "OrderBy",
												},
											},
										},
										Ellipsis: 5337,
									},
								},
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "query",
						},
						&ast.Ident{
							Name: "args",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "q",
										},
										Sel: &ast.Ident{
											Name: "PlaceholderFormat",
										},
									},
									Args: []ast.Expr{
										&ast.SelectorExpr{
											X: &ast.Ident{
												Name: "sq",
											},
											Sel: &ast.Ident{
												Name: "Dollar",
											},
										},
									},
								},
								Sel: &ast.Ident{
									Name: "MustSql",
								},
							},
						},
					},
				},
				&ast.IfStmt{
					Init: &ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.Ident{
								Name: "err",
							},
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "r",
										},
										Sel: &ast.Ident{
											Name: "database",
										},
									},
									Sel: &ast.Ident{
										Name: "SelectContext",
									},
								},
								Args: []ast.Expr{
									&ast.Ident{
										Name: "ctx",
									},
									&ast.UnaryExpr{
										Op: token.AND,
										X: &ast.Ident{
											Name: "dto",
										},
									},
									&ast.Ident{
										Name: "query",
									},
									&ast.Ident{
										Name: "args",
									},
								},
								Ellipsis: 5460,
							},
						},
					},
					Cond: &ast.BinaryExpr{
						X: &ast.Ident{
							Name: "err",
						},
						Op: token.NEQ,
						Y: &ast.Ident{
							Name: "nil",
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.Ident{
										Name: "e",
									},
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "errs",
											},
											Sel: &ast.Ident{
												Name: "FromPostgresError",
											},
										},
										Args: []ast.Expr{
											&ast.Ident{
												Name: "err",
											},
										},
									},
								},
							},
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.Ident{
										Name: "nil",
									},
									&ast.Ident{
										Name: "e",
									},
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "dto",
								},
								Sel: &ast.Ident{
									Name: "ToModels",
								},
							},
						},
						&ast.Ident{
							Name: "nil",
						},
					},
				},
			},
		},
	}
}

func (r Repository) SyncListMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, r.Path, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == "List" {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = r.AstListMethod()
	}
	for _, param := range r.Model.Params {
		param := param
		column := fmt.Sprintf("\"%s.%s\"", r.Model.TableName(), param.Tag())
		ast.Inspect(method, func(node ast.Node) bool {
			if call, ok := node.(*ast.CallExpr); ok {
				if fun, ok := call.Fun.(*ast.SelectorExpr); ok && fun.Sel.String() == "Select" {
					for _, arg := range call.Args {
						arg := arg
						if bl, ok := arg.(*ast.BasicLit); ok && bl.Value == column {
							return false
						}
					}
					call.Args = append(call.Args, &ast.BasicLit{
						Kind:  token.STRING,
						Value: column,
					})
					return false
				}
			}
			return true
		})
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(r.Path, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (r Repository) AstCountMethod() *ast.FuncDecl {
	tableName := r.Model.TableName()
	columns := []ast.Expr{
		&ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf("\"%s.id\"", tableName),
		},
		&ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf("\"%s.updated_at\"", tableName),
		},
		&ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf("\"%s.created_at\"", tableName),
		},
	}
	for _, param := range r.Model.Params {
		columns = append(
			columns,
			&ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf("\"%s.%s\"", tableName, param.Tag()),
			},
		)
	}
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{
							Name: "r",
						},
					},
					Type: &ast.StarExpr{
						X: &ast.Ident{
							Name: r.Name,
						},
					},
				},
			},
		},
		Name: &ast.Ident{
			Name: "Count",
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
									Name: r.Model.FilterTypeName(),
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
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "ctx",
						},
						&ast.Ident{
							Name: "cancel",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "context",
								},
								Sel: &ast.Ident{
									Name: "WithTimeout",
								},
							},
							Args: []ast.Expr{
								&ast.Ident{
									Name: "ctx",
								},
								&ast.SelectorExpr{
									X: &ast.Ident{
										Name: "time",
									},
									Sel: &ast.Ident{
										Name: "Second",
									},
								},
							},
						},
					},
				},
				&ast.DeferStmt{
					Call: &ast.CallExpr{
						Fun: &ast.Ident{
							Name: "cancel",
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "q",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "sq",
										},
										Sel: &ast.Ident{
											Name: "Select",
										},
									},
									Args: []ast.Expr{
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: "\"count(id)\"",
										},
									},
								},
								Sel: &ast.Ident{
									Name: "From",
								},
							},
							Args: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: fmt.Sprintf("\"public.%s\"", r.Model.TableName()),
								},
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "query",
						},
						&ast.Ident{
							Name: "args",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "q",
										},
										Sel: &ast.Ident{
											Name: "PlaceholderFormat",
										},
									},
									Args: []ast.Expr{
										&ast.SelectorExpr{
											X: &ast.Ident{
												Name: "sq",
											},
											Sel: &ast.Ident{
												Name: "Dollar",
											},
										},
									},
								},
								Sel: &ast.Ident{
									Name: "MustSql",
								},
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "result",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "r",
									},
									Sel: &ast.Ident{
										Name: "database",
									},
								},
								Sel: &ast.Ident{
									Name: "QueryRowxContext",
								},
							},
							Args: []ast.Expr{
								&ast.Ident{
									Name: "ctx",
								},
								&ast.Ident{
									Name: "query",
								},
								&ast.Ident{
									Name: "args",
								},
							},
							Ellipsis: 7757,
						},
					},
				},
				&ast.IfStmt{
					Init: &ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.Ident{
								Name: "err",
							},
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "result",
									},
									Sel: &ast.Ident{
										Name: "Err",
									},
								},
							},
						},
					},
					Cond: &ast.BinaryExpr{
						X: &ast.Ident{
							Name: "err",
						},
						Op: token.NEQ,
						Y: &ast.Ident{
							Name: "nil",
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.Ident{
										Name: "e",
									},
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "errs",
											},
											Sel: &ast.Ident{
												Name: "FromPostgresError",
											},
										},
										Args: []ast.Expr{
											&ast.Ident{
												Name: "err",
											},
										},
									},
								},
							},
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.BasicLit{
										Kind:  token.INT,
										Value: "0",
									},
									&ast.Ident{
										Name: "e",
									},
								},
							},
						},
					},
				},
				&ast.DeclStmt{
					Decl: &ast.GenDecl{
						Tok: token.VAR,
						Specs: []ast.Spec{
							&ast.ValueSpec{
								Names: []*ast.Ident{
									{
										Name: "count",
									},
								},
								Type: &ast.Ident{
									Name: "uint64",
								},
							},
						},
					},
				},
				&ast.IfStmt{
					Init: &ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.Ident{
								Name: "err",
							},
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "result",
									},
									Sel: &ast.Ident{
										Name: "Scan",
									},
								},
								Args: []ast.Expr{
									&ast.UnaryExpr{
										Op: token.AND,
										X: &ast.Ident{
											Name: "count",
										},
									},
								},
							},
						},
					},
					Cond: &ast.BinaryExpr{
						X: &ast.Ident{
							Name: "err",
						},
						Op: token.NEQ,
						Y: &ast.Ident{
							Name: "nil",
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.Ident{
										Name: "e",
									},
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "errs",
											},
											Sel: &ast.Ident{
												Name: "FromPostgresError",
											},
										},
										Args: []ast.Expr{
											&ast.Ident{
												Name: "err",
											},
										},
									},
								},
							},
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.BasicLit{
										Kind:  token.INT,
										Value: "0",
									},
									&ast.Ident{
										Name: "e",
									},
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.Ident{
							Name: "count",
						},
						&ast.Ident{
							Name: "nil",
						},
					},
				},
			},
		},
	}
}

func (r Repository) SyncCountMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, r.Path, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == "Count" {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = r.AstCountMethod()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(r.Path, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (r Repository) AstGetMethod() *ast.FuncDecl {
	tableName := r.Model.TableName()
	columns := []ast.Expr{
		&ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf("\"%s.id\"", tableName),
		},
		&ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf("\"%s.updated_at\"", tableName),
		},
		&ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf("\"%s.created_at\"", tableName),
		},
	}
	for _, param := range r.Model.Params {
		columns = append(
			columns,
			&ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf("\"%s.%s\"", tableName, param.Tag()),
			},
		)
	}
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{
							Name: "r",
						},
					},
					Type: &ast.StarExpr{
						X: &ast.Ident{
							Name: r.Name,
						},
					},
				},
			},
		},
		Name: &ast.Ident{
			Name: "Get",
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
									Name: r.Model.ModelName(),
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
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "ctx",
						},
						&ast.Ident{
							Name: "cancel",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "context",
								},
								Sel: &ast.Ident{
									Name: "WithTimeout",
								},
							},
							Args: []ast.Expr{
								&ast.Ident{
									Name: "ctx",
								},
								&ast.SelectorExpr{
									X: &ast.Ident{
										Name: "time",
									},
									Sel: &ast.Ident{
										Name: "Second",
									},
								},
							},
						},
					},
				},
				&ast.DeferStmt{
					Call: &ast.CallExpr{
						Fun: &ast.Ident{
							Name: "cancel",
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "dto",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X: &ast.CompositeLit{
								Type: &ast.Ident{
									Name: r.Model.PostgresDTOTypeName(),
								},
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "q",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.CallExpr{
													Fun: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "sq",
														},
														Sel: &ast.Ident{
															Name: "Select",
														},
													},
													Args: columns,
												},
												Sel: &ast.Ident{
													Name: "From",
												},
											},
											Args: []ast.Expr{
												&ast.BasicLit{
													Kind:  token.STRING,
													Value: fmt.Sprintf("\"public.%s\"", tableName),
												},
											},
										},
										Sel: &ast.Ident{
											Name: "Where",
										},
									},
									Args: []ast.Expr{
										&ast.CompositeLit{
											Type: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "sq",
												},
												Sel: &ast.Ident{
													Name: "Eq",
												},
											},
											Elts: []ast.Expr{
												&ast.KeyValueExpr{
													Key: &ast.BasicLit{
														Kind:  token.STRING,
														Value: "\"id\"",
													},
													Value: &ast.Ident{
														Name: "id",
													},
												},
											},
										},
									},
								},
								Sel: &ast.Ident{
									Name: "Limit",
								},
							},
							Args: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.INT,
									Value: "1",
								},
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "query",
						},
						&ast.Ident{
							Name: "args",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "q",
										},
										Sel: &ast.Ident{
											Name: "PlaceholderFormat",
										},
									},
									Args: []ast.Expr{
										&ast.SelectorExpr{
											X: &ast.Ident{
												Name: "sq",
											},
											Sel: &ast.Ident{
												Name: "Dollar",
											},
										},
									},
								},
								Sel: &ast.Ident{
									Name: "MustSql",
								},
							},
						},
					},
				},
				&ast.IfStmt{
					Init: &ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.Ident{
								Name: "err",
							},
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "r",
										},
										Sel: &ast.Ident{
											Name: "database",
										},
									},
									Sel: &ast.Ident{
										Name: "GetContext",
									},
								},
								Args: []ast.Expr{
									&ast.Ident{
										Name: "ctx",
									},
									&ast.Ident{
										Name: "dto",
									},
									&ast.Ident{
										Name: "query",
									},
									&ast.Ident{
										Name: "args",
									},
								},
								Ellipsis: 4211,
							},
						},
					},
					Cond: &ast.BinaryExpr{
						X: &ast.Ident{
							Name: "err",
						},
						Op: token.NEQ,
						Y: &ast.Ident{
							Name: "nil",
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.Ident{
										Name: "e",
									},
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "errs",
													},
													Sel: &ast.Ident{
														Name: "FromPostgresError",
													},
												},
												Args: []ast.Expr{
													&ast.Ident{
														Name: "err",
													},
												},
											},
											Sel: &ast.Ident{
												Name: "WithParam",
											},
										},
										Args: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: fmt.Sprintf("\"%s_id\"", r.Model.KeyName()),
											},
											&ast.CallExpr{
												Fun: &ast.Ident{
													Name: "string",
												},
												Args: []ast.Expr{
													&ast.Ident{
														Name: "id",
													},
												},
											},
										},
									},
								},
							},
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.Ident{
										Name: "nil",
									},
									&ast.Ident{
										Name: "e",
									},
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "dto",
								},
								Sel: &ast.Ident{
									Name: "ToModel",
								},
							},
						},
						&ast.Ident{
							Name: "nil",
						},
					},
				},
			},
		},
	}
}

func (r Repository) SyncGetMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, r.Path, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == "Get" {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = r.AstGetMethod()
	}
	for _, param := range r.Model.Params {
		param := param
		column := fmt.Sprintf("\"%s.%s\"", r.Model.TableName(), param.Tag())
		ast.Inspect(method, func(node ast.Node) bool {
			if call, ok := node.(*ast.CallExpr); ok {
				if fun, ok := call.Fun.(*ast.SelectorExpr); ok && fun.Sel.String() == "Select" {
					for _, arg := range call.Args {
						arg := arg
						if bl, ok := arg.(*ast.BasicLit); ok && bl.Value == column {
							return false
						}
					}
					call.Args = append(call.Args, &ast.BasicLit{
						Kind:  token.STRING,
						Value: column,
					})
					return false
				}
			}
			return true
		})
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(r.Path, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (r Repository) AstUpdateMethod() *ast.FuncDecl {
	tableName := r.Model.TableName()
	updateBlock := &ast.BlockStmt{
		List: []ast.Stmt{
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.Ident{
						Name: "q",
					},
				},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "q",
							},
							Sel: &ast.Ident{
								Name: "Set",
							},
						},
						Args: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: fmt.Sprintf("\"%s.%s\"", tableName, "updated_at"),
							},
							&ast.SelectorExpr{
								X: &ast.Ident{
									Name: "dto",
								},
								Sel: &ast.Ident{
									Name: "UpdatedAt",
								},
							},
						},
					},
				},
			},
		},
	}
	for _, param := range r.Model.Params {
		updateBlock.List = append(updateBlock.List, &ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: "q",
				},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: "q",
						},
						Sel: &ast.Ident{
							Name: "Set",
						},
					},
					Args: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf("\"%s.%s\"", tableName, param.Tag()),
						},
						&ast.SelectorExpr{
							X: &ast.Ident{
								Name: "dto",
							},
							Sel: &ast.Ident{
								Name: param.GetName(),
							},
						},
					},
				},
			},
		})
	}
	method := &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{
							Name: "r",
						},
					},
					Type: &ast.StarExpr{
						X: &ast.Ident{
							Name: r.Name,
						},
					},
				},
			},
		},
		Name: &ast.Ident{
			Name: "Update",
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
								Name: r.Model.Variable(),
							},
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "models",
								},
								Sel: &ast.Ident{
									Name: r.Model.ModelName(),
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
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "ctx",
						},
						&ast.Ident{
							Name: "cancel",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "context",
								},
								Sel: &ast.Ident{
									Name: "WithTimeout",
								},
							},
							Args: []ast.Expr{
								&ast.Ident{
									Name: "ctx",
								},
								&ast.SelectorExpr{
									X: &ast.Ident{
										Name: "time",
									},
									Sel: &ast.Ident{
										Name: "Second",
									},
								},
							},
						},
					},
				},
				&ast.DeferStmt{
					Call: &ast.CallExpr{
						Fun: &ast.Ident{
							Name: "cancel",
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "dto",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.Ident{
								Name: fmt.Sprintf("New%sFromModel", r.Model.PostgresDTOTypeName()),
							},
							Args: []ast.Expr{
								&ast.Ident{
									Name: r.Model.Variable(),
								},
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "q",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "sq",
										},
										Sel: &ast.Ident{
											Name: "Update",
										},
									},
									Args: []ast.Expr{
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: fmt.Sprintf("\"public.%s\"", tableName),
										},
									},
								},
								Sel: &ast.Ident{
									Name: "Where",
								},
							},
							Args: []ast.Expr{
								&ast.CompositeLit{
									Type: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "sq",
										},
										Sel: &ast.Ident{
											Name: "Eq",
										},
									},
									Elts: []ast.Expr{
										&ast.KeyValueExpr{
											Key: &ast.BasicLit{
												Kind:  token.STRING,
												Value: "\"id\"",
											},
											Value: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: r.Model.Variable(),
												},
												Sel: &ast.Ident{
													Name: "ID",
												},
											},
										},
									},
								},
							},
						},
					},
				},
				updateBlock,
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "query",
						},
						&ast.Ident{
							Name: "args",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "q",
										},
										Sel: &ast.Ident{
											Name: "PlaceholderFormat",
										},
									},
									Args: []ast.Expr{
										&ast.SelectorExpr{
											X: &ast.Ident{
												Name: "sq",
											},
											Sel: &ast.Ident{
												Name: "Dollar",
											},
										},
									},
								},
								Sel: &ast.Ident{
									Name: "MustSql",
								},
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "result",
						},
						&ast.Ident{
							Name: "err",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "r",
									},
									Sel: &ast.Ident{
										Name: "database",
									},
								},
								Sel: &ast.Ident{
									Name: "ExecContext",
								},
							},
							Args: []ast.Expr{
								&ast.Ident{
									Name: "ctx",
								},
								&ast.Ident{
									Name: "query",
								},
								&ast.Ident{
									Name: "args",
								},
							},
							Ellipsis: 6334,
						},
					},
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X: &ast.Ident{
							Name: "err",
						},
						Op: token.NEQ,
						Y: &ast.Ident{
							Name: "nil",
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.Ident{
										Name: "e",
									},
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "errs",
													},
													Sel: &ast.Ident{
														Name: "FromPostgresError",
													},
												},
												Args: []ast.Expr{
													&ast.Ident{
														Name: "err",
													},
												},
											},
											Sel: &ast.Ident{
												Name: "WithParam",
											},
										},
										Args: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: fmt.Sprintf("\"%s_id\"", r.Model.KeyName()),
											},
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "fmt",
													},
													Sel: &ast.Ident{
														Name: "Sprint",
													},
												},
												Args: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: r.Model.Variable(),
														},
														Sel: &ast.Ident{
															Name: "ID",
														},
													},
												},
											},
										},
									},
								},
							},
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.Ident{
										Name: "e",
									},
								},
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "affected",
						},
						&ast.Ident{
							Name: "err",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "result",
								},
								Sel: &ast.Ident{
									Name: "RowsAffected",
								},
							},
						},
					},
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X: &ast.Ident{
							Name: "err",
						},
						Op: token.NEQ,
						Y: &ast.Ident{
							Name: "nil",
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "errs",
													},
													Sel: &ast.Ident{
														Name: "FromPostgresError",
													},
												},
												Args: []ast.Expr{
													&ast.Ident{
														Name: "err",
													},
												},
											},
											Sel: &ast.Ident{
												Name: "WithParam",
											},
										},
										Args: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: fmt.Sprintf("\"%s_id\"", r.Model.KeyName()),
											},
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "fmt",
													},
													Sel: &ast.Ident{
														Name: "Sprint",
													},
												},
												Args: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: r.Model.Variable(),
														},
														Sel: &ast.Ident{
															Name: "ID",
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
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X: &ast.Ident{
							Name: "affected",
						},
						Op: token.EQL,
						Y: &ast.BasicLit{
							Kind:  token.INT,
							Value: "0",
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.Ident{
										Name: "e",
									},
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "errs",
													},
													Sel: &ast.Ident{
														Name: "NewEntityNotFound",
													},
												},
											},
											Sel: &ast.Ident{
												Name: "WithParam",
											},
										},
										Args: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: fmt.Sprintf("\"%s_id\"", r.Model.KeyName()),
											},
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "fmt",
													},
													Sel: &ast.Ident{
														Name: "Sprint",
													},
												},
												Args: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: r.Model.Variable(),
														},
														Sel: &ast.Ident{
															Name: "ID",
														},
													},
												},
											},
										},
									},
								},
							},
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.Ident{
										Name: "e",
									},
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.Ident{
							Name: "nil",
						},
					},
				},
			},
		},
	}
	return method
}

func (r Repository) SyncUpdateMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, r.Path, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == "Update" {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = r.AstUpdateMethod()
	}
	tableName := r.Model.TableName()
	for _, param := range r.Model.Params {
		param := param
		exists := false
		_ = param
		for _, stmt := range method.Body.List {
			if update, ok := stmt.(*ast.BlockStmt); ok {
				for _, updateStmt := range update.List {
					ast.Inspect(updateStmt, func(node ast.Node) bool {
						if call, ok := node.(*ast.CallExpr); ok {
							if callSelector, ok := call.Fun.(*ast.SelectorExpr); ok && callSelector.Sel.String() == "Set" {
								for _, arg := range call.Args {
									if bl, ok := arg.(*ast.BasicLit); ok && bl.Value == fmt.Sprintf("\"%s.%s\"", tableName, param.Tag()) {
										exists = true
										return false
									}
								}
							}
						}
						return true
					})
				}
				if !exists {
					update.List = append(
						update.List,
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "q",
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "q",
										},
										Sel: &ast.Ident{
											Name: "Set",
										},
									},
									Args: []ast.Expr{
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: fmt.Sprintf("\"%s.%s\"", tableName, param.Tag()),
										},
										&ast.SelectorExpr{
											X: &ast.Ident{
												Name: "dto",
											},
											Sel: &ast.Ident{
												Name: param.GetName(),
											},
										},
									},
								},
							},
						},
					)
				}
			}
		}
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(r.Path, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (r Repository) AstDeleteMethod() *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{
							Name: "r",
						},
					},
					Type: &ast.StarExpr{
						X: &ast.Ident{
							Name: r.Name,
						},
					},
				},
			},
		},
		Name: &ast.Ident{
			Name: "Delete",
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
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "ctx",
						},
						&ast.Ident{
							Name: "cancel",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "context",
								},
								Sel: &ast.Ident{
									Name: "WithTimeout",
								},
							},
							Args: []ast.Expr{
								&ast.Ident{
									Name: "ctx",
								},
								&ast.SelectorExpr{
									X: &ast.Ident{
										Name: "time",
									},
									Sel: &ast.Ident{
										Name: "Second",
									},
								},
							},
						},
					},
				},
				&ast.DeferStmt{
					Call: &ast.CallExpr{
						Fun: &ast.Ident{
							Name: "cancel",
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "q",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "sq",
										},
										Sel: &ast.Ident{
											Name: "Delete",
										},
									},
									Args: []ast.Expr{
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: fmt.Sprintf("\"public.%s\"", r.Model.TableName()),
										},
									},
								},
								Sel: &ast.Ident{
									Name: "Where",
								},
							},
							Args: []ast.Expr{
								&ast.CompositeLit{
									Type: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "sq",
										},
										Sel: &ast.Ident{
											Name: "Eq",
										},
									},
									Elts: []ast.Expr{
										&ast.KeyValueExpr{
											Key: &ast.BasicLit{
												Kind:  token.STRING,
												Value: "\"id\"",
											},
											Value: &ast.Ident{
												Name: "id",
											},
										},
									},
								},
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "query",
						},
						&ast.Ident{
							Name: "args",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "q",
										},
										Sel: &ast.Ident{
											Name: "PlaceholderFormat",
										},
									},
									Args: []ast.Expr{
										&ast.SelectorExpr{
											X: &ast.Ident{
												Name: "sq",
											},
											Sel: &ast.Ident{
												Name: "Dollar",
											},
										},
									},
								},
								Sel: &ast.Ident{
									Name: "MustSql",
								},
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "result",
						},
						&ast.Ident{
							Name: "err",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "r",
									},
									Sel: &ast.Ident{
										Name: "database",
									},
								},
								Sel: &ast.Ident{
									Name: "ExecContext",
								},
							},
							Args: []ast.Expr{
								&ast.Ident{
									Name: "ctx",
								},
								&ast.Ident{
									Name: "query",
								},
								&ast.Ident{
									Name: "args",
								},
							},
							Ellipsis: 7041,
						},
					},
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X: &ast.Ident{
							Name: "err",
						},
						Op: token.NEQ,
						Y: &ast.Ident{
							Name: "nil",
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.Ident{
										Name: "e",
									},
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "errs",
													},
													Sel: &ast.Ident{
														Name: "FromPostgresError",
													},
												},
												Args: []ast.Expr{
													&ast.Ident{
														Name: "err",
													},
												},
											},
											Sel: &ast.Ident{
												Name: "WithParam",
											},
										},
										Args: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: fmt.Sprintf("\"%s_id\"", r.Model.KeyName()),
											},
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "fmt",
													},
													Sel: &ast.Ident{
														Name: "Sprint",
													},
												},
												Args: []ast.Expr{
													&ast.Ident{
														Name: "id",
													},
												},
											},
										},
									},
								},
							},
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.Ident{
										Name: "e",
									},
								},
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "affected",
						},
						&ast.Ident{
							Name: "err",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "result",
								},
								Sel: &ast.Ident{
									Name: "RowsAffected",
								},
							},
						},
					},
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X: &ast.Ident{
							Name: "err",
						},
						Op: token.NEQ,
						Y: &ast.Ident{
							Name: "nil",
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.Ident{
										Name: "e",
									},
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "errs",
													},
													Sel: &ast.Ident{
														Name: "FromPostgresError",
													},
												},
												Args: []ast.Expr{
													&ast.Ident{
														Name: "err",
													},
												},
											},
											Sel: &ast.Ident{
												Name: "WithParam",
											},
										},
										Args: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: fmt.Sprintf("\"%s_id\"", r.Model.KeyName()),
											},
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "fmt",
													},
													Sel: &ast.Ident{
														Name: "Sprint",
													},
												},
												Args: []ast.Expr{
													&ast.Ident{
														Name: "id",
													},
												},
											},
										},
									},
								},
							},
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.Ident{
										Name: "e",
									},
								},
							},
						},
					},
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X: &ast.Ident{
							Name: "affected",
						},
						Op: token.EQL,
						Y: &ast.BasicLit{
							Kind:  token.INT,
							Value: "0",
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.Ident{
										Name: "e",
									},
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "errs",
													},
													Sel: &ast.Ident{
														Name: "NewEntityNotFound",
													},
												},
											},
											Sel: &ast.Ident{
												Name: "WithParam",
											},
										},
										Args: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: fmt.Sprintf("\"%s_id\"", r.Model.KeyName()),
											},
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "fmt",
													},
													Sel: &ast.Ident{
														Name: "Sprint",
													},
												},
												Args: []ast.Expr{
													&ast.Ident{
														Name: "id",
													},
												},
											},
										},
									},
								},
							},
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.Ident{
										Name: "e",
									},
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.Ident{
							Name: "nil",
						},
					},
				},
			},
		},
	}
}

func (r Repository) SyncDeleteMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, r.Path, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == "Delete" {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = r.AstDeleteMethod()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(r.Path, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (r Repository) AstDTOListType() *ast.TypeSpec {
	return &ast.TypeSpec{
		Name: &ast.Ident{
			Name: r.Model.PostgresDTOListTypeName(),
		},
		Type: &ast.ArrayType{
			Elt: &ast.StarExpr{
				X: &ast.Ident{
					Name: r.Model.PostgresDTOTypeName(),
				},
			},
		},
	}
}

func (r Repository) SyncDTOListType() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, r.Path, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var structureExists bool
	var dtoListType *ast.TypeSpec
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok && t.Name.String() == r.Model.PostgresDTOListTypeName() {
			dtoListType = t
			structureExists = true
			return false
		}
		return true
	})
	if dtoListType == nil {
		dtoListType = r.AstDTOListType()
	}
	if !structureExists {
		gd := &ast.GenDecl{
			Doc:    nil,
			TokPos: 0,
			Tok:    token.TYPE,
			Lparen: 0,
			Specs:  []ast.Spec{dtoListType},
			Rparen: 0,
		}
		file.Decls = append(file.Decls, gd)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(r.Path, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (r Repository) AstDTOToModels() *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{
							Name: "list",
						},
					},
					Type: &ast.Ident{
						Name: r.Model.PostgresDTOListTypeName(),
					},
				},
			},
		},
		Name: &ast.Ident{
			Name: "ToModels",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
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
										Name: r.Model.ModelName(),
									},
								},
							},
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: r.Model.ListVariable(),
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.Ident{
								Name: "make",
							},
							Args: []ast.Expr{
								&ast.ArrayType{
									Elt: &ast.StarExpr{
										X: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "models",
											},
											Sel: &ast.Ident{
												Name: r.Model.ModelName(),
											},
										},
									},
								},
								&ast.CallExpr{
									Fun: &ast.Ident{
										Name: "len",
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "list",
										},
									},
								},
							},
						},
					},
				},
				&ast.RangeStmt{
					Key: &ast.Ident{
						Name: "i",
					},
					Tok: token.DEFINE,
					X: &ast.Ident{
						Name: "list",
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.IndexExpr{
										X: &ast.Ident{
											Name: r.Model.ListVariable(),
										},
										Index: &ast.Ident{
											Name: "i",
										},
									},
								},
								Tok: token.ASSIGN,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.IndexExpr{
												X: &ast.Ident{
													Name: "list",
												},
												Index: &ast.Ident{
													Name: "i",
												},
											},
											Sel: &ast.Ident{
												Name: "ToModel",
											},
										},
									},
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.Ident{
							Name: r.Model.ListVariable(),
						},
					},
				},
			},
		},
	}
}

func (r Repository) SyncDTOListToModels() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, r.Path, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExists bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == "ToModels" {
			methodExists = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = r.AstDTOToModels()
	}
	if !methodExists {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(r.Path, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}