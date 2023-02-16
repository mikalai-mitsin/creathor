package models

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
)

type Interceptor struct {
	Path   string
	Name   string
	Model  *Model
	Params []*Param
}

func (i Interceptor) AstStruct() *ast.TypeSpec {
	structure := &ast.TypeSpec{
		Doc:        nil,
		Name:       ast.NewIdent(i.Name),
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
	for _, param := range i.Params {
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

func (i Interceptor) SyncStruct() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, i.Path, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var structureExists bool
	var structure *ast.TypeSpec
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok && t.Name.String() == i.Name {
			structure = t
			structureExists = true
			return false
		}
		return true
	})
	if structure == nil {
		structure = i.AstStruct()
	}
	for _, param := range i.Params {
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
	if err := os.WriteFile(i.Path, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (i Interceptor) AstConstructor() *ast.FuncDecl {
	var args []*ast.Field
	cl := &ast.CompositeLit{
		Type:       ast.NewIdent(i.Name),
		Lbrace:     0,
		Elts:       nil,
		Rbrace:     0,
		Incomplete: false,
	}
	for _, param := range i.Params {
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
		Name: ast.NewIdent(fmt.Sprintf("New%s", i.Name)),
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
						Type:    ast.NewIdent(fmt.Sprintf("interceptors.%s", i.Name)),
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

func (i Interceptor) SyncConstructor() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, i.Path, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var structureConstructorExists bool
	var structureConstructor *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == fmt.Sprintf("New%s", i.Name) {
			structureConstructorExists = true
			structureConstructor = t
			return false
		}
		return true
	})
	if structureConstructor == nil {
		structureConstructor = i.AstConstructor()
	}
	for _, param := range i.Params {
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
				if t, ok := cl.Type.(*ast.Ident); ok && t.String() == i.Name {
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
	if err := os.WriteFile(i.Path, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (i Interceptor) AstCreateMethod() *ast.FuncDecl {
	return &ast.FuncDecl{
		Doc: nil,
		Recv: &ast.FieldList{
			Opening: 0,
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("i"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(i.Name),
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
								Sel: ast.NewIdent(i.Model.CreateTypeName()),
							},
						},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("requestUser")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("models"),
								Sel: ast.NewIdent("User"),
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
								Sel: ast.NewIdent(i.Model.ModelName()),
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
				// Check permission
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
										X:   ast.NewIdent("i"),
										Sel: ast.NewIdent("authUseCase"),
									},
									Sel: ast.NewIdent("HasPermission"),
								},
								Lparen: 0,
								Args: []ast.Expr{
									ast.NewIdent("ctx"),
									ast.NewIdent("requestUser"),
									&ast.SelectorExpr{
										X:   ast.NewIdent("models"),
										Sel: ast.NewIdent(i.Model.PermissionIDCreate()),
									},
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
				},
				// Check permission
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
										X:   ast.NewIdent("i"),
										Sel: ast.NewIdent("authUseCase"),
									},
									Sel: ast.NewIdent("HasObjectPermission"),
								},
								Lparen: 0,
								Args: []ast.Expr{
									ast.NewIdent("ctx"),
									ast.NewIdent("requestUser"),
									&ast.SelectorExpr{
										X:   ast.NewIdent("models"),
										Sel: ast.NewIdent(i.Model.PermissionIDCreate()),
									},
									ast.NewIdent("create"),
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
				},
				// Try to create model at use case
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent(i.Model.Variable()),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("i"),
									Sel: ast.NewIdent(i.Model.UseCaseVariableName()),
								},
								Sel: ast.NewIdent("Create"),
							},
							Args: []ast.Expr{
								ast.NewIdent("ctx"),
								ast.NewIdent("create"),
							},
						},
					},
				},
				// Check error
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
				// Return created model and nil error
				&ast.ReturnStmt{
					Results: []ast.Expr{
						ast.NewIdent(i.Model.Variable()),
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
}

func (i Interceptor) SyncCreateMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, i.Path, nil, parser.ParseComments)
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
		method = i.AstCreateMethod()
	}
	for _, param := range i.Model.Params {
		param := param
		ast.Inspect(method, func(node ast.Node) bool {
			if cl, ok := node.(*ast.CompositeLit); ok {
				if t, ok := cl.Type.(*ast.SelectorExpr); ok && t.Sel.String() == i.Model.ModelName() {
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
	if err := os.WriteFile(i.Path, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (i Interceptor) AstListMethod() *ast.FuncDecl {
	return &ast.FuncDecl{
		Doc: nil,
		Recv: &ast.FieldList{
			Opening: 0,
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("i"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(i.Name),
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
								Sel: ast.NewIdent(i.Model.FilterTypeName()),
							},
						},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("requestUser")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("models"),
								Sel: ast.NewIdent("User"),
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
									Sel: ast.NewIdent(i.Model.ModelName()),
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
				// Check permission
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
										X:   ast.NewIdent("i"),
										Sel: ast.NewIdent("authUseCase"),
									},
									Sel: ast.NewIdent("HasPermission"),
								},
								Lparen: 0,
								Args: []ast.Expr{
									ast.NewIdent("ctx"),
									ast.NewIdent("requestUser"),
									&ast.SelectorExpr{
										X:   ast.NewIdent("models"),
										Sel: ast.NewIdent(i.Model.PermissionIDList()),
									},
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
									ast.NewIdent("0"),
									ast.NewIdent("err"),
								},
							},
						},
					},
				},
				// Check filter permission
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
										X:   ast.NewIdent("i"),
										Sel: ast.NewIdent("authUseCase"),
									},
									Sel: ast.NewIdent("HasObjectPermission"),
								},
								Lparen: 0,
								Args: []ast.Expr{
									ast.NewIdent("ctx"),
									ast.NewIdent("requestUser"),
									&ast.SelectorExpr{
										X:   ast.NewIdent("models"),
										Sel: ast.NewIdent(i.Model.PermissionIDList()),
									},
									ast.NewIdent("filter"),
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
									ast.NewIdent("0"),
									ast.NewIdent("err"),
								},
							},
						},
					},
				},
				// Try to update model at use case
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent(i.Model.ListVariable()),
						ast.NewIdent("count"),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("i"),
									Sel: ast.NewIdent(i.Model.UseCaseVariableName()),
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
				// Check error
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
				// Return created model and nil error
				&ast.ReturnStmt{
					Results: []ast.Expr{
						ast.NewIdent(i.Model.ListVariable()),
						ast.NewIdent("count"),
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
}

func (i Interceptor) SyncListMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, i.Path, nil, parser.ParseComments)
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
		method = i.AstListMethod()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(i.Path, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (i Interceptor) AstGetMethod() *ast.FuncDecl {
	return &ast.FuncDecl{
		Doc: nil,
		Recv: &ast.FieldList{
			Opening: 0,
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("i"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(i.Name),
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
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("models"),
							Sel: ast.NewIdent("UUID"),
						},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("requestUser")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("models"),
								Sel: ast.NewIdent("User"),
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
								Sel: ast.NewIdent(i.Model.ModelName()),
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
				// Check permission
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
										X:   ast.NewIdent("i"),
										Sel: ast.NewIdent("authUseCase"),
									},
									Sel: ast.NewIdent("HasPermission"),
								},
								Lparen: 0,
								Args: []ast.Expr{
									ast.NewIdent("ctx"),
									ast.NewIdent("requestUser"),
									&ast.SelectorExpr{
										X:   ast.NewIdent("models"),
										Sel: ast.NewIdent(i.Model.PermissionIDDetail()),
									},
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
				},
				// Try to get model from use case
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent(i.Model.Variable()),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("i"),
									Sel: ast.NewIdent(i.Model.UseCaseVariableName()),
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
				// Check error
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
				// Check object permission
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
										X:   ast.NewIdent("i"),
										Sel: ast.NewIdent("authUseCase"),
									},
									Sel: ast.NewIdent("HasObjectPermission"),
								},
								Lparen: 0,
								Args: []ast.Expr{
									ast.NewIdent("ctx"),
									ast.NewIdent("requestUser"),
									&ast.SelectorExpr{
										X:   ast.NewIdent("models"),
										Sel: ast.NewIdent(i.Model.PermissionIDDetail()),
									},
									ast.NewIdent(i.Model.Variable()),
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
				},
				// Return created model and nil error
				&ast.ReturnStmt{
					Results: []ast.Expr{
						ast.NewIdent(i.Model.Variable()),
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
}

func (i Interceptor) SyncGetMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, i.Path, nil, parser.ParseComments)
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
		method = i.AstGetMethod()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(i.Path, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (i Interceptor) AstUpdateMethod() *ast.FuncDecl {
	return &ast.FuncDecl{
		Doc: nil,
		Recv: &ast.FieldList{
			Opening: 0,
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("i"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(i.Name),
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
								Sel: ast.NewIdent(i.Model.UpdateTypeName()),
							},
						},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("requestUser")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("models"),
								Sel: ast.NewIdent("User"),
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
								Sel: ast.NewIdent(i.Model.ModelName()),
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
				// Check permission
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
										X:   ast.NewIdent("i"),
										Sel: ast.NewIdent("authUseCase"),
									},
									Sel: ast.NewIdent("HasPermission"),
								},
								Lparen: 0,
								Args: []ast.Expr{
									ast.NewIdent("ctx"),
									ast.NewIdent("requestUser"),
									&ast.SelectorExpr{
										X:   ast.NewIdent("models"),
										Sel: ast.NewIdent(i.Model.PermissionIDUpdate()),
									},
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
				},
				// Try to get model from use case
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent(i.Model.Variable()),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("i"),
									Sel: ast.NewIdent(i.Model.UseCaseVariableName()),
								},
								Sel: ast.NewIdent("Get"),
							},
							Args: []ast.Expr{
								ast.NewIdent("ctx"),
								&ast.SelectorExpr{
									X:   ast.NewIdent("update"),
									Sel: ast.NewIdent("ID"),
								},
							},
						},
					},
				},
				// Check error
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
				// Check object permission
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
										X:   ast.NewIdent("i"),
										Sel: ast.NewIdent("authUseCase"),
									},
									Sel: ast.NewIdent("HasObjectPermission"),
								},
								Lparen: 0,
								Args: []ast.Expr{
									ast.NewIdent("ctx"),
									ast.NewIdent("requestUser"),
									&ast.SelectorExpr{
										X:   ast.NewIdent("models"),
										Sel: ast.NewIdent(i.Model.PermissionIDUpdate()),
									},
									ast.NewIdent(i.Model.Variable()),
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
				},
				// Try to update model at use case
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("updated"),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("i"),
									Sel: ast.NewIdent(i.Model.UseCaseVariableName()),
								},
								Sel: ast.NewIdent("Update"),
							},
							Args: []ast.Expr{
								ast.NewIdent("ctx"),
								ast.NewIdent("update"),
							},
						},
					},
				},
				// Check error
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
				// Return created model and nil error
				&ast.ReturnStmt{
					Results: []ast.Expr{
						ast.NewIdent("updated"),
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
}

func (i Interceptor) SyncUpdateMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, i.Path, nil, parser.ParseComments)
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
		method = i.AstUpdateMethod()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(i.Path, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (i Interceptor) AstDeleteMethod() *ast.FuncDecl {
	return &ast.FuncDecl{
		Doc: nil,
		Recv: &ast.FieldList{
			Opening: 0,
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("i"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(i.Name),
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
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("models"),
							Sel: ast.NewIdent("UUID"),
						},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("requestUser")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("models"),
								Sel: ast.NewIdent("User"),
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
				// Check permission
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
										X:   ast.NewIdent("i"),
										Sel: ast.NewIdent("authUseCase"),
									},
									Sel: ast.NewIdent("HasPermission"),
								},
								Args: []ast.Expr{
									ast.NewIdent("ctx"),
									ast.NewIdent("requestUser"),
									&ast.SelectorExpr{
										X:   ast.NewIdent("models"),
										Sel: ast.NewIdent(i.Model.PermissionIDDelete()),
									},
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
				},
				// Try to get model from use case
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent(i.Model.Variable()),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("i"),
									Sel: ast.NewIdent(i.Model.UseCaseVariableName()),
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
				// Check error
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
									ast.NewIdent("err"),
								},
							},
						},
					},
					Else: nil,
				},
				// Check object permission
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
										X:   ast.NewIdent("i"),
										Sel: ast.NewIdent("authUseCase"),
									},
									Sel: ast.NewIdent("HasObjectPermission"),
								},
								Lparen: 0,
								Args: []ast.Expr{
									ast.NewIdent("ctx"),
									ast.NewIdent("requestUser"),
									&ast.SelectorExpr{
										X:   ast.NewIdent("models"),
										Sel: ast.NewIdent(i.Model.PermissionIDDelete()),
									},
									ast.NewIdent(i.Model.Variable()),
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
				},
				// Try to delete model at use case
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
										X:   ast.NewIdent("i"),
										Sel: ast.NewIdent(i.Model.UseCaseVariableName()),
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
				},
				// Return created model and nil error
				&ast.ReturnStmt{
					Results: []ast.Expr{
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
}

func (i Interceptor) SyncDeleteMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, i.Path, nil, parser.ParseComments)
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
		method = i.AstDeleteMethod()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(i.Path, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}
