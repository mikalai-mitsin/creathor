package services

import (
	"bytes"
	"fmt"
	"github.com/mikalai-mitsin/creathor/internal/pkg/tmpl"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"
	"path/filepath"

	"github.com/mikalai-mitsin/creathor/internal/pkg/domain"
)

type ServiceCrud struct {
	domain *domain.Domain
}

func NewServiceCrud(domain *domain.Domain) *ServiceCrud {
	return &ServiceCrud{domain: domain}
}

func (u ServiceCrud) Sync() error {
	err := os.MkdirAll(path.Dir(u.filename()), 0777)
	if err != nil {
		return err
	}
	if err := u.syncStruct(); err != nil {
		return err
	}
	if err := u.syncConstructor(); err != nil {
		return err
	}
	if err := u.syncCreateMethod(); err != nil {
		return err
	}
	if err := u.syncGetMethod(); err != nil {
		return err
	}
	if err := u.syncListMethod(); err != nil {
		return err
	}
	if err := u.syncUpdateMethod(); err != nil {
		return err
	}
	if err := u.syncDeleteMethod(); err != nil {
		return err
	}
	if u.domain.SnakeName() == "user" {
		if err := u.syncGetByEmailMethod(); err != nil {
			return err
		}
	}
	if err := u.syncTest(); err != nil {
		return err
	}
	return nil
}

func (u ServiceCrud) filename() string {
	return filepath.Join("internal", "app", u.domain.DirName(), "services", "service.go")
}

func (u ServiceCrud) file() *ast.File {
	return &ast.File{
		Name: ast.NewIdent("services"),
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
							Value: u.domain.EntitiesImportPath(),
						},
					},
				},
			},
		},
	}
}

func (u ServiceCrud) structure() *ast.TypeSpec {
	structure := &ast.TypeSpec{
		Name: ast.NewIdent(u.domain.Service.Name),
		Type: &ast.StructType{
			Fields: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent(u.domain.Repository.Variable)},
						Type:  ast.NewIdent(u.domain.Repository.Name),
					},
					{
						Names: []*ast.Ident{ast.NewIdent("clock")},
						Type:  ast.NewIdent("Clock"),
					},
					{
						Names: []*ast.Ident{ast.NewIdent("logger")},
						Type:  ast.NewIdent("Logger"),
					},
					{
						Names: []*ast.Ident{ast.NewIdent("uuid")},
						Type:  ast.NewIdent("UUIDGenerator"),
					},
				},
			},
		},
	}
	return structure
}

func (u ServiceCrud) syncStruct() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, u.filename(), nil, parser.ParseComments)
	if err != nil {
		file = u.file()
	}
	var structureExists bool
	var structure *ast.TypeSpec
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok && t.Name.String() == u.domain.Service.Name {
			structure = t
			structureExists = true
			return false
		}
		return true
	})
	if structure == nil {
		structure = u.structure()
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
	if err := os.WriteFile(u.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (u ServiceCrud) constructor() *ast.FuncDecl {
	constructor := &ast.FuncDecl{
		Name: ast.NewIdent(fmt.Sprintf("New%s", u.domain.Service.Name)),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent(u.domain.Repository.Variable)},
						Type:  ast.NewIdent(u.domain.Repository.Name),
					},
					{
						Names: []*ast.Ident{ast.NewIdent("clock")},
						Type:  ast.NewIdent("Clock"),
					},
					{
						Names: []*ast.Ident{ast.NewIdent("logger")},
						Type:  ast.NewIdent("Logger"),
					},
					{
						Names: []*ast.Ident{ast.NewIdent("uuid")},
						Type:  ast.NewIdent("UUIDGenerator"),
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent(fmt.Sprintf("*%s", u.domain.Service.Name)),
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
								Type: ast.NewIdent(u.domain.Service.Name),
								Elts: []ast.Expr{
									&ast.KeyValueExpr{
										Key:   ast.NewIdent(u.domain.Repository.Variable),
										Value: ast.NewIdent(u.domain.Repository.Variable),
									},
									&ast.KeyValueExpr{
										Key:   ast.NewIdent("clock"),
										Value: ast.NewIdent("clock"),
									},
									&ast.KeyValueExpr{
										Key:   ast.NewIdent("logger"),
										Value: ast.NewIdent("logger"),
									},
									&ast.KeyValueExpr{
										Key:   ast.NewIdent("uuid"),
										Value: ast.NewIdent("uuid"),
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return constructor
}

func (u ServiceCrud) syncConstructor() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, u.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var structureConstructorExists bool
	var structureConstructor *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok &&
			t.Name.String() == fmt.Sprintf("New%s", u.domain.Service.Name) {
			structureConstructorExists = true
			structureConstructor = t
			return false
		}
		return true
	})
	if structureConstructor == nil {
		structureConstructor = u.constructor()
	}
	if !structureConstructorExists {
		file.Decls = append(file.Decls, structureConstructor)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(u.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (u ServiceCrud) create() *ast.FuncDecl {
	params := []ast.Expr{
		&ast.KeyValueExpr{
			Key:   ast.NewIdent("ID"),
			Value: ast.NewIdent(`u.uuid.NewUUID()`),
		},
		&ast.KeyValueExpr{
			Key:   ast.NewIdent("UpdatedAt"),
			Value: ast.NewIdent("now"),
		},
		&ast.KeyValueExpr{
			Key:   ast.NewIdent("CreatedAt"),
			Value: ast.NewIdent("now"),
		},
	}
	for _, param := range u.domain.GetCreateModel().Params {
		params = append(params, &ast.KeyValueExpr{
			Key: ast.NewIdent(param.GetName()),
			Value: &ast.SelectorExpr{
				X:   ast.NewIdent("create"),
				Sel: ast.NewIdent(param.GetName()),
			},
		})
	}
	fun := &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("u"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(u.domain.Service.Name),
					},
				},
			},
		},
		Name: ast.NewIdent("Create"),
		Type: &ast.FuncType{
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
								X:   ast.NewIdent("entities"),
								Sel: ast.NewIdent(u.domain.GetCreateModel().Name),
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
								X:   ast.NewIdent("entities"),
								Sel: ast.NewIdent(u.domain.GetMainModel().Name),
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
					Lhs: []ast.Expr{ast.NewIdent(u.domain.GetMainModel().Variable)},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X: &ast.CompositeLit{
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("entities"),
									Sel: ast.NewIdent(u.domain.GetMainModel().Name),
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
										Sel: ast.NewIdent(u.domain.Repository.Variable),
									},
									Sel: ast.NewIdent("Create"),
								},
								Args: []ast.Expr{
									ast.NewIdent("ctx"),
									ast.NewIdent(u.domain.GetMainModel().Variable),
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
						ast.NewIdent(u.domain.GetMainModel().Variable),
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
	return fun
}

func (u ServiceCrud) syncCreateMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, u.filename(), nil, parser.ParseComments)
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
		method = u.create()
	}
	for _, param := range u.domain.GetCreateModel().Params {
		param := param
		ast.Inspect(method, func(node ast.Node) bool {
			if cl, ok := node.(*ast.CompositeLit); ok {
				if t, ok := cl.Type.(*ast.SelectorExpr); ok &&
					t.Sel.String() == u.domain.GetMainModel().Name {
					for _, elt := range cl.Elts {
						if kv, ok := elt.(*ast.KeyValueExpr); ok {
							if key, ok := kv.Key.(*ast.Ident); ok &&
								key.String() == param.GetName() {
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
	if err := os.WriteFile(u.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (u ServiceCrud) list() *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			Opening: 0,
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("u"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(u.domain.Service.Name),
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
						Names: []*ast.Ident{ast.NewIdent(u.domain.GetFilterModel().Variable)},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("entities"),
								Sel: ast.NewIdent(u.domain.GetFilterModel().Name),
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
									X:   ast.NewIdent("entities"),
									Sel: ast.NewIdent(u.domain.GetMainModel().Name),
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
						ast.NewIdent(u.domain.GetMainModel().Variable),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("u"),
									Sel: ast.NewIdent(u.domain.Repository.Variable),
								},
								Sel: ast.NewIdent("List"),
							},
							Args: []ast.Expr{
								ast.NewIdent("ctx"),
								ast.NewIdent(u.domain.GetFilterModel().Variable),
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
									Sel: ast.NewIdent(u.domain.Repository.Variable),
								},
								Sel: ast.NewIdent("Count"),
							},
							Args: []ast.Expr{
								ast.NewIdent("ctx"),
								ast.NewIdent(u.domain.GetFilterModel().Variable),
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
						ast.NewIdent(u.domain.GetMainModel().Variable),
						ast.NewIdent("count"),
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
}

func (u ServiceCrud) syncListMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, u.filename(), nil, parser.ParseComments)
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
		method = u.list()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(u.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (u ServiceCrud) get() *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			Opening: 0,
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("u"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(u.domain.Service.Name),
					},
				},
			},
		},
		Name: ast.NewIdent("Get"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("ctx")},
						Type:  ast.NewIdent("context.Context"),
					},
					{
						Names: []*ast.Ident{ast.NewIdent("id")},
						Type:  ast.NewIdent("uuid.UUID"),
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("entities"),
								Sel: ast.NewIdent(u.domain.GetMainModel().Name),
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
						ast.NewIdent(u.domain.GetMainModel().Variable),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("u"),
									Sel: ast.NewIdent(u.domain.Repository.Variable),
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
						ast.NewIdent(u.domain.GetMainModel().Variable),
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
}

func (u ServiceCrud) syncGetMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, u.filename(), nil, parser.ParseComments)
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
		method = u.get()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(u.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (u ServiceCrud) update() *ast.FuncDecl {
	block := &ast.BlockStmt{
		List: []ast.Stmt{},
	}
	for _, param := range u.domain.GetUpdateModel().Params {
		if param.Name == "ID" {
			continue
		}
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
								X:   ast.NewIdent(u.domain.GetMainModel().Variable),
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
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("u"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(u.domain.Service.Name),
					},
				},
			},
		},
		Name: ast.NewIdent("Update"),
		Type: &ast.FuncType{
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
								X:   ast.NewIdent("entities"),
								Sel: ast.NewIdent(u.domain.GetUpdateModel().Name),
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
								X:   ast.NewIdent("entities"),
								Sel: ast.NewIdent(u.domain.GetMainModel().Name),
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
						ast.NewIdent(u.domain.GetMainModel().Variable),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("u"),
									Sel: ast.NewIdent(u.domain.Repository.Variable),
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
				// Block of updated fields
				block,
				// Set updated at
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X:   ast.NewIdent(u.domain.GetMainModel().Variable),
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
										Sel: ast.NewIdent(u.domain.Repository.Variable),
									},
									Sel: ast.NewIdent("Update"),
								},
								Args: []ast.Expr{
									ast.NewIdent("ctx"),
									ast.NewIdent(u.domain.GetMainModel().Variable),
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
						ast.NewIdent(u.domain.GetMainModel().Variable),
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
	return fun
}

func (u ServiceCrud) syncUpdateMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, u.filename(), nil, parser.ParseComments)
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
		method = u.update()
	}
	for _, param := range u.domain.GetUpdateModel().Params {
		param := param
		if param.Name == "ID" {
			continue
		}
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
											X:   ast.NewIdent(u.domain.GetMainModel().Variable),
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
	if err := os.WriteFile(u.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (u ServiceCrud) delete() *ast.FuncDecl {
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
						X: ast.NewIdent(u.domain.Service.Name),
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
						Type:  ast.NewIdent("uuid.UUID"),
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
										Sel: ast.NewIdent(u.domain.Repository.Variable),
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

func (u ServiceCrud) syncDeleteMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, u.filename(), nil, parser.ParseComments)
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
		method = u.delete()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(u.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (u ServiceCrud) getByEmail() *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			Opening: 0,
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("u"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(u.domain.Service.Name),
					},
				},
			},
		},
		Name: ast.NewIdent("GetByEmail"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("ctx")},
						Type:  ast.NewIdent("context.Context"),
					},
					{
						Names: []*ast.Ident{ast.NewIdent("email")},
						Type:  ast.NewIdent("string"),
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("entities"),
								Sel: ast.NewIdent(u.domain.GetMainModel().Name),
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
						ast.NewIdent(u.domain.GetMainModel().Variable),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("u"),
									Sel: ast.NewIdent(u.domain.Repository.Variable),
								},
								Sel: ast.NewIdent("GetByEmail"),
							},
							Args: []ast.Expr{
								ast.NewIdent("ctx"),
								ast.NewIdent("email"),
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
						ast.NewIdent(u.domain.GetMainModel().Variable),
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
}

func (u ServiceCrud) syncGetByEmailMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, u.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == "GetByEmail" {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = u.getByEmail()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(u.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

var destinationPath = "."

func (u ServiceCrud) syncTest() error {
	test := tmpl.Template{
		SourcePath: "templates/internal/domain/services/crud_test.go.tmpl",
		DestinationPath: filepath.Join(
			destinationPath,
			"internal",
			"app",
			u.domain.DirName(),
			"services",
			u.domain.TestFileName(),
		),
		Name: "service test",
	}
	if err := test.RenderToFile(u.domain); err != nil {
		return err
	}
	return nil
}
