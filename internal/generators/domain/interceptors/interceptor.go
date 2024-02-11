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
	"path/filepath"

	mods "github.com/018bf/creathor/internal/domain"
)

type InterceptorCrud struct {
	mod *mods.Domain
}

func NewInterceptorCrud(mod *mods.Domain) *InterceptorCrud {
	return &InterceptorCrud{mod: mod}
}

func (i InterceptorCrud) Sync() error {
	err := os.MkdirAll(path.Dir(i.filename()), 0777)
	if err != nil {
		return err
	}
	if err := i.syncStruct(); err != nil {
		return err
	}
	if err := i.syncConstructor(); err != nil {
		return err
	}
	for _, method := range i.mod.Interceptor.Methods {
		switch method.Name {
		case "Create":
			if err := i.syncCreateMethod(method); err != nil {
				return err
			}
		case "Get":
			if err := i.syncGetMethod(method); err != nil {
				return err
			}
		case "List":
			if err := i.syncListMethod(method); err != nil {
				return err
			}
		case "Update":
			if err := i.syncUpdateMethod(method); err != nil {
				return err
			}
		case "Delete":
			if err := i.syncDeleteMethod(method); err != nil {
				return err
			}
		}
	}
	return nil
}

func (i InterceptorCrud) filename() string {
	return filepath.Join("internal", i.mod.Name, "interceptors", i.mod.Filename)
}

func (i InterceptorCrud) structure() *ast.TypeSpec {
	fields := []*ast.Field{
		{
			Names: []*ast.Ident{ast.NewIdent(i.mod.UseCase.Variable)},
			Type:  ast.NewIdent(i.mod.UseCase.Name),
		},
		{
			Names: []*ast.Ident{ast.NewIdent("logger")},
			Type: &ast.SelectorExpr{
				X:   ast.NewIdent("log"),
				Sel: ast.NewIdent("Logger"),
			},
		},
	}
	if i.mod.Auth {
		fields = append(
			fields,
			&ast.Field{
				Names: []*ast.Ident{ast.NewIdent("authUseCase")},
				Type:  ast.NewIdent("AuthUseCase"),
			},
		)
	}
	structure := &ast.TypeSpec{
		Name: ast.NewIdent(i.mod.Interceptor.Name),
		Type: &ast.StructType{
			Fields: &ast.FieldList{
				List: fields,
			},
		},
	}
	return structure
}

func (i InterceptorCrud) syncStruct() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, i.filename(), nil, parser.ParseComments)
	if err != nil {
		file = i.file()
	}
	var structureExists bool
	var structure *ast.TypeSpec
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok && t.Name.String() == i.mod.Interceptor.Name {
			structure = t
			structureExists = true
			return false
		}
		return true
	})
	if structure == nil {
		structure = i.structure()
	}
	if !structureExists {
		gd := &ast.GenDecl{
			Tok:   token.TYPE,
			Specs: []ast.Spec{structure},
		}
		file.Decls = append(file.Decls, gd)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(i.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (i InterceptorCrud) constructor() *ast.FuncDecl {
	fields := []*ast.Field{
		{
			Names: []*ast.Ident{ast.NewIdent(i.mod.UseCase.Variable)},
			Type:  ast.NewIdent(i.mod.UseCase.Name),
		},
		{
			Names: []*ast.Ident{ast.NewIdent("logger")},
			Type: &ast.SelectorExpr{
				X:   ast.NewIdent("log"),
				Sel: ast.NewIdent("Logger"),
			},
		},
	}
	exprs := []ast.Expr{
		&ast.KeyValueExpr{
			Key:   ast.NewIdent(i.mod.UseCase.Variable),
			Value: ast.NewIdent(i.mod.UseCase.Variable),
		},
		&ast.KeyValueExpr{
			Key:   ast.NewIdent("logger"),
			Value: ast.NewIdent("logger"),
		},
	}
	if i.mod.Auth {
		fields = append(
			fields,
			&ast.Field{
				Names: []*ast.Ident{ast.NewIdent("authUseCase")},
				Type:  ast.NewIdent("AuthUseCase"),
			},
		)
		exprs = append(
			exprs,
			&ast.KeyValueExpr{
				Key:   ast.NewIdent("authUseCase"),
				Value: ast.NewIdent("authUseCase"),
			},
		)
	}
	constructor := &ast.FuncDecl{
		Name: ast.NewIdent(fmt.Sprintf("New%s", i.mod.Interceptor.Name)),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: fields,
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent(
							fmt.Sprintf("*%s", i.mod.Interceptor.Name),
						),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X: &ast.CompositeLit{
								Type: ast.NewIdent(i.mod.Interceptor.Name),
								Elts: exprs,
							},
						},
					},
				},
			},
		},
	}
	return constructor
}

func (i InterceptorCrud) syncConstructor() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, i.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var structureConstructorExists bool
	var structureConstructor *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok &&
			t.Name.String() == fmt.Sprintf("New%s", i.mod.Interceptor.Name) {
			structureConstructorExists = true
			structureConstructor = t
			return false
		}
		return true
	})
	if structureConstructor == nil {
		structureConstructor = i.constructor()
	}
	if !structureConstructorExists {
		file.Decls = append(file.Decls, structureConstructor)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(i.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (i InterceptorCrud) createMethod(method *mods.Method) *ast.FuncDecl {
	var body []ast.Stmt
	if i.mod.Auth {
		body = append(body,
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.Ident{
						Name: "requestUser",
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
									Name: "i",
								},
								Sel: &ast.Ident{
									Name: "authUseCase",
								},
							},
							Sel: &ast.Ident{
								Name: "GetUser",
							},
						},
						Args: []ast.Expr{
							&ast.Ident{
								Name: "ctx",
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
								&ast.Ident{
									Name: "nil",
								},
								&ast.Ident{
									Name: "err",
								},
							},
						},
					},
				},
			},
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
									X:   ast.NewIdent("userModels"),
									Sel: ast.NewIdent(i.mod.PermissionIDCreate()),
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
							Args: []ast.Expr{
								ast.NewIdent("ctx"),
								ast.NewIdent("requestUser"),
								&ast.SelectorExpr{
									X:   ast.NewIdent("userModels"),
									Sel: ast.NewIdent(i.mod.PermissionIDCreate()),
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
		)
	}
	body = append(body,
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(i.mod.GetMainModel().Variable),
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.SelectorExpr{
							X:   ast.NewIdent("i"),
							Sel: ast.NewIdent(i.mod.UseCase.Variable),
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
				ast.NewIdent(i.mod.GetMainModel().Variable),
				ast.NewIdent("nil"),
			},
		},
	)
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("i"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(i.mod.Interceptor.Name),
					},
				},
			},
		},
		Name: ast.NewIdent("Create"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: method.Args,
			},
			Results: &ast.FieldList{
				List: method.Return,
			},
		},
		Body: &ast.BlockStmt{
			List: body,
		},
	}
}

func (i InterceptorCrud) syncCreateMethod(m *mods.Method) error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, i.filename(), nil, parser.ParseComments)
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
		method = i.createMethod(m)
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}

	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(i.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (i InterceptorCrud) astListMethod(m *mods.Method) *ast.FuncDecl {
	var body []ast.Stmt
	if i.mod.Auth {
		body = append(body,
			// Get request user
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.Ident{
						Name: "requestUser",
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
									Name: "i",
								},
								Sel: &ast.Ident{
									Name: "authUseCase",
								},
							},
							Sel: &ast.Ident{
								Name: "GetUser",
							},
						},
						Args: []ast.Expr{
							&ast.Ident{
								Name: "ctx",
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
								ast.NewIdent("nil"),
								ast.NewIdent("0"),
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
								Sel: ast.NewIdent("HasPermission"),
							},
							Lparen: 0,
							Args: []ast.Expr{
								ast.NewIdent("ctx"),
								ast.NewIdent("requestUser"),
								&ast.SelectorExpr{
									X:   ast.NewIdent("userModels"),
									Sel: ast.NewIdent(i.mod.PermissionIDList()),
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
									X:   ast.NewIdent("userModels"),
									Sel: ast.NewIdent(i.mod.PermissionIDList()),
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
		)
	}
	body = append(body,
		// Try to update model at use case
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent("items"),
				ast.NewIdent("count"),
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.SelectorExpr{
							X:   ast.NewIdent("i"),
							Sel: ast.NewIdent(i.mod.UseCase.Variable),
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
				ast.NewIdent("items"),
				ast.NewIdent("count"),
				ast.NewIdent("nil"),
			},
		},
	)
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("i"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(i.mod.Interceptor.Name),
					},
				},
			},
			Closing: 0,
		},
		Name: ast.NewIdent("List"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: m.Args,
			},
			Results: &ast.FieldList{
				List: m.Return,
			},
		},
		Body: &ast.BlockStmt{
			List: body,
		},
	}
}

func (i InterceptorCrud) syncListMethod(m *mods.Method) error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, i.filename(), nil, parser.ParseComments)
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
		method = i.astListMethod(m)
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(i.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (i InterceptorCrud) astGetMethod(m *mods.Method) *ast.FuncDecl {
	var body []ast.Stmt
	if i.mod.Auth {
		body = append(
			body,
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.Ident{
						Name: "requestUser",
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
									Name: "i",
								},
								Sel: &ast.Ident{
									Name: "authUseCase",
								},
							},
							Sel: &ast.Ident{
								Name: "GetUser",
							},
						},
						Args: []ast.Expr{
							&ast.Ident{
								Name: "ctx",
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
								&ast.Ident{
									Name: "nil",
								},
								&ast.Ident{
									Name: "err",
								},
							},
						},
					},
				},
			},
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
									X:   ast.NewIdent("userModels"),
									Sel: ast.NewIdent(i.mod.PermissionIDDetail()),
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
		)
	}
	body = append(
		body,
		// Try to get model from use case
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(i.mod.GetMainModel().Variable),
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.SelectorExpr{
							X:   ast.NewIdent("i"),
							Sel: ast.NewIdent(i.mod.UseCase.Variable),
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
	)
	if i.mod.Auth {
		body = append(
			body,
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
									X:   ast.NewIdent("userModels"),
									Sel: ast.NewIdent(i.mod.PermissionIDDetail()),
								},
								ast.NewIdent(i.mod.GetMainModel().Variable),
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
		)
	}
	body = append(
		body,
		// Return created model and nil error
		&ast.ReturnStmt{
			Results: []ast.Expr{
				ast.NewIdent(i.mod.GetMainModel().Variable),
				ast.NewIdent("nil"),
			},
		},
	)
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("i"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(i.mod.Interceptor.Name),
					},
				},
			},
		},
		Name: ast.NewIdent("Get"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: m.Args,
			},
			Results: &ast.FieldList{
				List: m.Return,
			},
		},
		Body: &ast.BlockStmt{
			List: body,
		},
	}
}

func (i InterceptorCrud) syncGetMethod(m *mods.Method) error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, i.filename(), nil, parser.ParseComments)
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
		method = i.astGetMethod(m)
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(i.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (i InterceptorCrud) updateMethod(m *mods.Method) *ast.FuncDecl {
	var body []ast.Stmt
	if i.mod.Auth {
		body = append(body,
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.Ident{
						Name: "requestUser",
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
									Name: "i",
								},
								Sel: &ast.Ident{
									Name: "authUseCase",
								},
							},
							Sel: &ast.Ident{
								Name: "GetUser",
							},
						},
						Args: []ast.Expr{
							&ast.Ident{
								Name: "ctx",
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
								&ast.Ident{
									Name: "nil",
								},
								&ast.Ident{
									Name: "err",
								},
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
								Sel: ast.NewIdent("HasPermission"),
							},
							Lparen: 0,
							Args: []ast.Expr{
								ast.NewIdent("ctx"),
								ast.NewIdent("requestUser"),
								&ast.SelectorExpr{
									X:   ast.NewIdent("userModels"),
									Sel: ast.NewIdent(i.mod.PermissionIDUpdate()),
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
					ast.NewIdent(i.mod.GetMainModel().Variable),
					ast.NewIdent("err"),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("i"),
								Sel: ast.NewIdent(i.mod.UseCase.Variable),
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
									X:   ast.NewIdent("userModels"),
									Sel: ast.NewIdent(i.mod.PermissionIDUpdate()),
								},
								ast.NewIdent(i.mod.GetMainModel().Variable),
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
		)
	}

	body = append(body,
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
							Sel: ast.NewIdent(i.mod.UseCase.Variable),
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
	)
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("i"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(i.mod.Interceptor.Name),
					},
				},
			},
		},
		Name: ast.NewIdent("Update"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: m.Args,
			},
			Results: &ast.FieldList{
				List: m.Return,
			},
		},
		Body: &ast.BlockStmt{
			List: body,
		},
	}
}

func (i InterceptorCrud) syncUpdateMethod(m *mods.Method) error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, i.filename(), nil, parser.ParseComments)
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
		method = i.updateMethod(m)
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(i.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (i InterceptorCrud) deleteMethod(m *mods.Method) *ast.FuncDecl {
	var body []ast.Stmt
	if i.mod.Auth {
		body = append(body,
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.Ident{
						Name: "requestUser",
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
									Name: "i",
								},
								Sel: &ast.Ident{
									Name: "authUseCase",
								},
							},
							Sel: &ast.Ident{
								Name: "GetUser",
							},
						},
						Args: []ast.Expr{
							&ast.Ident{
								Name: "ctx",
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
								Sel: ast.NewIdent("HasPermission"),
							},
							Lparen: 0,
							Args: []ast.Expr{
								ast.NewIdent("ctx"),
								ast.NewIdent("requestUser"),
								&ast.SelectorExpr{
									X:   ast.NewIdent("userModels"),
									Sel: ast.NewIdent(i.mod.PermissionIDDelete()),
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
					ast.NewIdent(i.mod.GetMainModel().Variable),
					ast.NewIdent("err"),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("i"),
								Sel: ast.NewIdent(i.mod.UseCase.Variable),
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
									X:   ast.NewIdent("userModels"),
									Sel: ast.NewIdent(i.mod.PermissionIDDelete()),
								},
								ast.NewIdent(i.mod.GetMainModel().Variable),
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
		)
	}
	body = append(body,
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
								Sel: ast.NewIdent(i.mod.UseCase.Variable),
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
	)
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("i"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(i.mod.Interceptor.Name),
					},
				},
			},
		},
		Name: ast.NewIdent("Delete"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: m.Args,
			},
			Results: &ast.FieldList{
				List: m.Return,
			},
		},
		Body: &ast.BlockStmt{
			List: body,
		},
	}
}

func (i InterceptorCrud) syncDeleteMethod(m *mods.Method) error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, i.filename(), nil, parser.ParseComments)
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
		method = i.deleteMethod(m)
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(i.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (i InterceptorCrud) file() *ast.File {
	return &ast.File{
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
							Value: fmt.Sprintf(`"%s/internal/%s/models"`, i.mod.Module, i.mod.Name),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/pkg/log"`, i.mod.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/pkg/uuid"`, i.mod.Module),
						},
					},
					&ast.ImportSpec{
						Name: ast.NewIdent("userModels"),
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/user/models"`, i.mod.Module),
						},
					},
				},
			},
		},
	}
}
