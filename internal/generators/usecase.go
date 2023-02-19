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

type UseCase struct {
	Path   string
	Name   string
	Model  *models.ModelConfig
	Params []*Param
}

func (u UseCase) AstStruct() *ast.TypeSpec {
	structure := &ast.TypeSpec{
		Doc:        nil,
		Name:       ast.NewIdent(u.Name),
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
	for _, param := range u.Params {
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

func (u UseCase) SyncStruct() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, u.Path, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var structureExists bool
	var structure *ast.TypeSpec
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok && t.Name.String() == u.Name {
			structure = t
			structureExists = true
			return false
		}
		return true
	})
	if structure == nil {
		structure = u.AstStruct()
	}
	for _, param := range u.Params {
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
	if err := os.WriteFile(u.Path, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (u UseCase) AstConstructor() *ast.FuncDecl {
	var args []*ast.Field
	cl := &ast.CompositeLit{
		Type:       ast.NewIdent(u.Name),
		Lbrace:     0,
		Elts:       nil,
		Rbrace:     0,
		Incomplete: false,
	}
	for _, param := range u.Params {
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
		Name: ast.NewIdent(fmt.Sprintf("New%s", u.Name)),
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
						Doc:     nil,
						Names:   nil,
						Type:    ast.NewIdent(fmt.Sprintf("usecases.%s", u.Name)),
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

func (u UseCase) SyncConstructor() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, u.Path, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var structureConstructorExists bool
	var structureConstructor *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == fmt.Sprintf("New%s", u.Name) {
			structureConstructorExists = true
			structureConstructor = t
			return false
		}
		return true
	})
	if structureConstructor == nil {
		structureConstructor = u.AstConstructor()
	}
	for _, param := range u.Params {
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
				if t, ok := cl.Type.(*ast.Ident); ok && t.String() == u.Name {
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
	if err := os.WriteFile(u.Path, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (u UseCase) AstCreateMethod() *ast.FuncDecl {
	params := []ast.Expr{
		&ast.KeyValueExpr{
			Key:   ast.NewIdent("ID"),
			Colon: 0,
			Value: ast.NewIdent("\"\""),
		},
		&ast.KeyValueExpr{
			Key:   ast.NewIdent("UpdatedAt"),
			Colon: 0,
			Value: ast.NewIdent("now"),
		},
		&ast.KeyValueExpr{
			Key:   ast.NewIdent("CreatedAt"),
			Colon: 0,
			Value: ast.NewIdent("now"),
		},
	}
	for _, param := range u.Model.Params {
		params = append(params, &ast.KeyValueExpr{
			Key:   ast.NewIdent(param.GetName()),
			Colon: 0,
			Value: &ast.SelectorExpr{
				X:   ast.NewIdent("create"),
				Sel: ast.NewIdent(param.GetName()),
			},
		})
	}
	fun := &ast.FuncDecl{
		Doc: nil,
		Recv: &ast.FieldList{
			Opening: 0,
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("u"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(u.Name),
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
						Names: []*ast.Ident{ast.NewIdent("create")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("models"),
								Sel: ast.NewIdent(u.Model.CreateTypeName()),
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
								X:   ast.NewIdent("models"),
								Sel: ast.NewIdent(u.Model.ModelName()),
							},
						},
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				// Create validation
				&ast.IfStmt{
					Init: &ast.AssignStmt{
						Lhs: []ast.Expr{
							ast.NewIdent("err"),
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("create"),
									Sel: ast.NewIdent("Validate"),
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
							&ast.ReturnStmt{
								Results: []ast.Expr{
									ast.NewIdent("nil"),
									ast.NewIdent("err"),
								},
							},
						},
					},
					Else: nil,
				},
				// Now from clock
				&ast.AssignStmt{
					Lhs:    []ast.Expr{ast.NewIdent("now")},
					TokPos: 0,
					Tok:    token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X:   ast.NewIdent("u"),
											Sel: ast.NewIdent("clock"),
										},
										Sel: ast.NewIdent("Now"),
									},
								},
								Sel: ast.NewIdent("UTC"),
							},
						},
					},
				},
				// Fill model struct from create form
				&ast.AssignStmt{
					Lhs: []ast.Expr{ast.NewIdent(u.Model.Variable())},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X: &ast.CompositeLit{
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("models"),
									Sel: ast.NewIdent(u.Model.ModelName()),
								},
								Elts: params,
							},
						},
					},
				},
				// Try to create model at repository
				&ast.IfStmt{
					Init: &ast.AssignStmt{
						Lhs: []ast.Expr{
							ast.NewIdent("err"),
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.SelectorExpr{
										X:   ast.NewIdent("u"),
										Sel: ast.NewIdent(u.Model.RepositoryVariableName()),
									},
									Sel: ast.NewIdent("Create"),
								},
								Args: []ast.Expr{
									ast.NewIdent("ctx"),
									ast.NewIdent(u.Model.Variable()),
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
							&ast.ReturnStmt{
								Results: []ast.Expr{
									ast.NewIdent("nil"),
									ast.NewIdent("err"),
								},
							},
						},
					},
					Else: nil,
				},
				// Return created model and nil error
				&ast.ReturnStmt{
					Results: []ast.Expr{
						ast.NewIdent(u.Model.Variable()),
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
	return fun
}

func (u UseCase) SyncCreateMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, u.Path, nil, parser.ParseComments)
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
		method = u.AstCreateMethod()
	}
	for _, param := range u.Model.Params {
		param := param
		ast.Inspect(method, func(node ast.Node) bool {
			if cl, ok := node.(*ast.CompositeLit); ok {
				if t, ok := cl.Type.(*ast.SelectorExpr); ok && t.Sel.String() == u.Model.ModelName() {
					for _, elt := range cl.Elts {
						if kv, ok := elt.(*ast.KeyValueExpr); ok {
							if key, ok := kv.Key.(*ast.Ident); ok && key.String() == param.GetName() {
								return false
							}
						}
					}
					cl.Elts = append(cl.Elts, &ast.KeyValueExpr{
						Key:   ast.NewIdent(param.GetName()),
						Colon: 0,
						Value: &ast.SelectorExpr{
							X:   ast.NewIdent("create"),
							Sel: ast.NewIdent(param.GetName()),
						},
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
	if err := os.WriteFile(u.Path, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (u UseCase) AstListMethod() *ast.FuncDecl {
	return &ast.FuncDecl{
		Doc: nil,
		Recv: &ast.FieldList{
			Opening: 0,
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("u"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(u.Name),
					},
				},
			},
			Closing: 0,
		},
		Name: ast.NewIdent("List"),
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
						Names: []*ast.Ident{ast.NewIdent("filter")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("models"),
								Sel: ast.NewIdent(u.Model.FilterTypeName()),
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
									X:   ast.NewIdent("models"),
									Sel: ast.NewIdent(u.Model.ModelName()),
								},
							},
						},
					},
					{
						Type: ast.NewIdent("uint64"),
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent(u.Model.Variable()),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("u"),
									Sel: ast.NewIdent(u.Model.RepositoryVariableName()),
								},
								Sel: ast.NewIdent("List"),
							},
							Args: []ast.Expr{
								ast.NewIdent("ctx"),
								ast.NewIdent("filter"),
							},
						},
					},
				},
				&ast.IfStmt{
					Init: nil,
					Cond: &ast.BinaryExpr{
						X:  ast.NewIdent("err"),
						Op: token.NEQ,
						Y:  ast.NewIdent("nil"),
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{
									ast.NewIdent("nil"),
									ast.NewIdent("0"),
									ast.NewIdent("err"),
								},
							},
						},
					},
					Else: nil,
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("count"),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("u"),
									Sel: ast.NewIdent(u.Model.RepositoryVariableName()),
								},
								Sel: ast.NewIdent("Count"),
							},
							Args: []ast.Expr{
								ast.NewIdent("ctx"),
								ast.NewIdent("filter"),
							},
						},
					},
				},
				&ast.IfStmt{
					Init: nil,
					Cond: &ast.BinaryExpr{
						X:  ast.NewIdent("err"),
						Op: token.NEQ,
						Y:  ast.NewIdent("nil"),
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{
									ast.NewIdent("nil"),
									ast.NewIdent("0"),
									ast.NewIdent("err"),
								},
							},
						},
					},
					Else: nil,
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						ast.NewIdent(u.Model.Variable()),
						ast.NewIdent("count"),
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
}

func (u UseCase) SyncListMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, u.Path, nil, parser.ParseComments)
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
		method = u.AstListMethod()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(u.Path, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (u UseCase) AstGetMethod() *ast.FuncDecl {
	return &ast.FuncDecl{
		Doc: nil,
		Recv: &ast.FieldList{
			Opening: 0,
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("u"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(u.Name),
					},
				},
			},
			Closing: 0,
		},
		Name: ast.NewIdent("Get"),
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
						Names: []*ast.Ident{ast.NewIdent("id")},
						Type:  ast.NewIdent("models.UUID"),
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("models"),
								Sel: ast.NewIdent(u.Model.ModelName()),
							},
						},
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent(u.Model.Variable()),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("u"),
									Sel: ast.NewIdent(u.Model.RepositoryVariableName()),
								},
								Sel: ast.NewIdent("Get"),
							},
							Args: []ast.Expr{
								ast.NewIdent("ctx"),
								ast.NewIdent("id"),
							},
						},
					},
				},
				&ast.IfStmt{
					Init: nil,
					Cond: &ast.BinaryExpr{
						X:  ast.NewIdent("err"),
						Op: token.NEQ,
						Y:  ast.NewIdent("nil"),
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{
									ast.NewIdent("nil"),
									ast.NewIdent("err"),
								},
							},
						},
					},
					Else: nil,
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						ast.NewIdent(u.Model.Variable()),
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
}

func (u UseCase) SyncGetMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, u.Path, nil, parser.ParseComments)
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
		method = u.AstGetMethod()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(u.Path, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (u UseCase) AstUpdateMethod() *ast.FuncDecl {
	block := &ast.BlockStmt{
		Lbrace: 0,
		List:   []ast.Stmt{},
		Rbrace: 0,
	}
	for _, param := range u.Model.Params {
		block.List = append(block.List, &ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X: &ast.SelectorExpr{
					X:   ast.NewIdent("update"),
					Sel: ast.NewIdent(param.GetName()),
				},
				Op: token.NEQ,
				Y:  ast.NewIdent("nil"),
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.SelectorExpr{
								X:   ast.NewIdent(u.Model.Variable()),
								Sel: ast.NewIdent(param.GetName()),
							},
						},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							&ast.StarExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("update"),
									Sel: ast.NewIdent(param.GetName()),
								},
							},
						},
					},
				},
			},
		})
	}
	fun := &ast.FuncDecl{
		Doc: nil,
		Recv: &ast.FieldList{
			Opening: 0,
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("u"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(u.Name),
					},
				},
			},
			Closing: 0,
		},
		Name: ast.NewIdent("Update"),
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
						Names: []*ast.Ident{ast.NewIdent("update")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("models"),
								Sel: ast.NewIdent(u.Model.UpdateTypeName()),
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
								X:   ast.NewIdent("models"),
								Sel: ast.NewIdent(u.Model.ModelName()),
							},
						},
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				// Update from validation
				&ast.IfStmt{
					Init: &ast.AssignStmt{
						Lhs: []ast.Expr{
							ast.NewIdent("err"),
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("update"),
									Sel: ast.NewIdent("Validate"),
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
							&ast.ReturnStmt{
								Results: []ast.Expr{
									ast.NewIdent("nil"),
									ast.NewIdent("err"),
								},
							},
						},
					},
					Else: nil,
				},
				// Get model to update
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent(u.Model.Variable()),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("u"),
									Sel: ast.NewIdent(u.Model.RepositoryVariableName()),
								},
								Sel: ast.NewIdent("Get"),
							},
							Args: []ast.Expr{
								ast.NewIdent("ctx"),
								ast.NewIdent("update.ID"),
							},
						},
					},
				},
				&ast.IfStmt{
					Init: nil,
					Cond: &ast.BinaryExpr{
						X:  ast.NewIdent("err"),
						Op: token.NEQ,
						Y:  ast.NewIdent("nil"),
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{
									ast.NewIdent("nil"),
									ast.NewIdent("err"),
								},
							},
						},
					},
					Else: nil,
				},
				// Block of updated fields
				block,
				// Set updated at
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X:   ast.NewIdent(u.Model.Variable()),
							Sel: ast.NewIdent("UpdatedAt"),
						},
					},
					TokPos: 0,
					Tok:    token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X:   ast.NewIdent("u"),
											Sel: ast.NewIdent("clock"),
										},
										Sel: ast.NewIdent("Now"),
									},
								},
								Sel: ast.NewIdent("UTC"),
							},
						},
					},
				},
				// Try to update model at repository
				&ast.IfStmt{
					Init: &ast.AssignStmt{
						Lhs: []ast.Expr{
							ast.NewIdent("err"),
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.SelectorExpr{
										X:   ast.NewIdent("u"),
										Sel: ast.NewIdent(u.Model.RepositoryVariableName()),
									},
									Sel: ast.NewIdent("Update"),
								},
								Args: []ast.Expr{
									ast.NewIdent("ctx"),
									ast.NewIdent(u.Model.Variable()),
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
							&ast.ReturnStmt{
								Results: []ast.Expr{
									ast.NewIdent("nil"),
									ast.NewIdent("err"),
								},
							},
						},
					},
					Else: nil,
				},
				// Return updated model and nil error
				&ast.ReturnStmt{
					Results: []ast.Expr{
						ast.NewIdent(u.Model.Variable()),
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
	return fun
}

func (u UseCase) SyncUpdateMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, u.Path, nil, parser.ParseComments)
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
		method = u.AstUpdateMethod()
	}
	for _, param := range u.Model.Params {
		param := param
		exists := false
		for _, stmt := range method.Body.List {
			if update, ok := stmt.(*ast.BlockStmt); ok {
				for _, updateStmt := range update.List {
					ast.Inspect(updateStmt, func(node ast.Node) bool {
						if ifStmt, ok := node.(*ast.IfStmt); ok {
							if binaryExpr, ok := ifStmt.Cond.(*ast.BinaryExpr); ok {
								if selectorExpr, ok := binaryExpr.X.(*ast.SelectorExpr); ok {
									if selectorExpr.Sel.String() == param.GetName() {
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
					update.List = append(update.List, &ast.IfStmt{
						Cond: &ast.BinaryExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("update"),
								Sel: ast.NewIdent(param.GetName()),
							},
							Op: token.NEQ,
							Y:  ast.NewIdent("nil"),
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.AssignStmt{
									Lhs: []ast.Expr{
										&ast.SelectorExpr{
											X:   ast.NewIdent(u.Model.Variable()),
											Sel: ast.NewIdent(param.GetName()),
										},
									},
									Tok: token.ASSIGN,
									Rhs: []ast.Expr{
										&ast.StarExpr{
											X: &ast.SelectorExpr{
												X:   ast.NewIdent("update"),
												Sel: ast.NewIdent(param.GetName()),
											},
										},
									},
								},
							},
						},
					})
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
	if err := os.WriteFile(u.Path, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (u UseCase) AstDeleteMethod() *ast.FuncDecl {
	return &ast.FuncDecl{
		Doc: nil,
		Recv: &ast.FieldList{
			Opening: 0,
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("u"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(u.Name),
					},
				},
			},
			Closing: 0,
		},
		Name: ast.NewIdent("Delete"),
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
						Names: []*ast.Ident{ast.NewIdent("id")},
						Type:  ast.NewIdent("models.UUID"),
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
				&ast.IfStmt{
					Init: &ast.AssignStmt{
						Lhs: []ast.Expr{
							ast.NewIdent("err"),
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.SelectorExpr{
										X:   ast.NewIdent("u"),
										Sel: ast.NewIdent(u.Model.RepositoryVariableName()),
									},
									Sel: ast.NewIdent("Delete"),
								},
								Args: []ast.Expr{
									ast.NewIdent("ctx"),
									ast.NewIdent("id"),
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
							&ast.ReturnStmt{
								Results: []ast.Expr{
									ast.NewIdent("err"),
								},
							},
						},
					},
					Else: nil,
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
}

func (u UseCase) SyncDeleteMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, u.Path, nil, parser.ParseComments)
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
		method = u.AstDeleteMethod()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(u.Path, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}
