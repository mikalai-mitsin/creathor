package services

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

	"github.com/mikalai-mitsin/creathor/internal/pkg/app"
	"github.com/mikalai-mitsin/creathor/internal/pkg/astfile"
)

type ServiceGenerator struct {
	domain *app.BaseEntity
}

func NewServiceGenerator(domain *app.BaseEntity) *ServiceGenerator {
	return &ServiceGenerator{domain: domain}
}

func (u ServiceGenerator) Sync() error {
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
	return nil
}

func (u ServiceGenerator) filename() string {
	return filepath.Join("internal", "app", u.domain.AppName(), "services", u.domain.DirName(), u.domain.FileName())
}

func (u ServiceGenerator) file() *ast.File {
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
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/pkg/uuid"`, u.domain.Module),
						},
					},
				},
			},
		},
	}
}

func (u ServiceGenerator) structure() *ast.TypeSpec {
	structure := &ast.TypeSpec{
		Name: ast.NewIdent(u.domain.GetServiceTypeName()),
		Type: &ast.StructType{
			Fields: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent(u.domain.GetRepositoryPrivateVariableName()),
						},
						Type: ast.NewIdent(u.domain.GetRepositoryInterfaceName()),
					},
					{
						Names: []*ast.Ident{ast.NewIdent("clock")},
						Type:  ast.NewIdent("clock"),
					},
					{
						Names: []*ast.Ident{ast.NewIdent("logger")},
						Type:  ast.NewIdent("logger"),
					},
					{
						Names: []*ast.Ident{ast.NewIdent("uuid")},
						Type:  ast.NewIdent("uuidGenerator"),
					},
				},
			},
		},
	}
	return structure
}

func (u ServiceGenerator) syncStruct() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, u.filename(), nil, parser.ParseComments)
	if err != nil {
		file = u.file()
	}
	structure, structureExists := astfile.FindType(file, u.domain.GetServiceTypeName())
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

func (u ServiceGenerator) constructor() *ast.FuncDecl {
	constructor := &ast.FuncDecl{
		Name: ast.NewIdent(u.domain.GetServiceConstructorName()),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent(u.domain.GetRepositoryPrivateVariableName()),
						},
						Type: ast.NewIdent(u.domain.GetRepositoryInterfaceName()),
					},
					{
						Names: []*ast.Ident{ast.NewIdent("clock")},
						Type:  ast.NewIdent("clock"),
					},
					{
						Names: []*ast.Ident{ast.NewIdent("logger")},
						Type:  ast.NewIdent("logger"),
					},
					{
						Names: []*ast.Ident{ast.NewIdent("uuid")},
						Type:  ast.NewIdent("uuidGenerator"),
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent(fmt.Sprintf("*%s", u.domain.GetServiceTypeName())),
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
								Type: ast.NewIdent(u.domain.GetServiceTypeName()),
								Elts: []ast.Expr{
									&ast.KeyValueExpr{
										Key: ast.NewIdent(
											u.domain.GetRepositoryPrivateVariableName(),
										),
										Value: ast.NewIdent(
											u.domain.GetRepositoryPrivateVariableName(),
										),
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

func (u ServiceGenerator) syncConstructor() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, u.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	method, methodExist := astfile.FindFunc(file, u.domain.GetServiceConstructorName())
	if method == nil {
		method = u.constructor()
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

func (u ServiceGenerator) create() *ast.FuncDecl {
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
						X: ast.NewIdent(u.domain.GetServiceTypeName()),
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
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("entities"),
							Sel: ast.NewIdent(u.domain.GetCreateModel().Name),
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("entities"),
							Sel: ast.NewIdent(u.domain.GetMainModel().Name),
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
									&ast.CompositeLit{
										Type: &ast.SelectorExpr{
											X:   ast.NewIdent("entities"),
											Sel: ast.NewIdent(u.domain.GetMainModel().Name),
										},
									},
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
						&ast.CompositeLit{
							Type: &ast.SelectorExpr{
								X:   ast.NewIdent("entities"),
								Sel: ast.NewIdent(u.domain.GetMainModel().Name),
							},
							Elts: params,
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
										X: ast.NewIdent("u"),
										Sel: ast.NewIdent(
											u.domain.GetRepositoryPrivateVariableName(),
										),
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
									&ast.CompositeLit{
										Type: &ast.SelectorExpr{
											X:   ast.NewIdent("entities"),
											Sel: ast.NewIdent(u.domain.GetMainModel().Name),
										},
									},
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

func (u ServiceGenerator) syncCreateMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, u.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	method, methodExist := astfile.FindFunc(file, "Create")
	if method == nil {
		method = u.create()
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

func (u ServiceGenerator) list() *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			Opening: 0,
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("u"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(u.domain.GetServiceTypeName()),
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
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("entities"),
							Sel: ast.NewIdent(u.domain.GetFilterModel().Name),
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.ArrayType{
							Elt: &ast.SelectorExpr{
								X:   ast.NewIdent("entities"),
								Sel: ast.NewIdent(u.domain.GetMainModel().Name),
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
									Sel: ast.NewIdent(u.domain.GetRepositoryPrivateVariableName()),
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
									Sel: ast.NewIdent(u.domain.GetRepositoryPrivateVariableName()),
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

func (u ServiceGenerator) syncListMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, u.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	method, methodExist := astfile.FindFunc(file, "List")
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

func (u ServiceGenerator) get() *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			Opening: 0,
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("u"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(u.domain.GetServiceTypeName()),
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
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("entities"),
							Sel: ast.NewIdent(u.domain.GetMainModel().Name),
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
									Sel: ast.NewIdent(u.domain.GetRepositoryPrivateVariableName()),
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
									&ast.CompositeLit{
										Type: &ast.SelectorExpr{
											X:   ast.NewIdent("entities"),
											Sel: ast.NewIdent(u.domain.GetMainModel().Name),
										},
									},
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

func (u ServiceGenerator) syncGetMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, u.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	method, methodExist := astfile.FindFunc(file, "Get")
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

func (u ServiceGenerator) update() *ast.FuncDecl {
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
						X: ast.NewIdent(u.domain.GetServiceTypeName()),
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
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("entities"),
							Sel: ast.NewIdent(u.domain.GetUpdateModel().Name),
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("entities"),
							Sel: ast.NewIdent(u.domain.GetMainModel().Name),
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
									&ast.CompositeLit{
										Type: &ast.SelectorExpr{
											X:   ast.NewIdent("entities"),
											Sel: ast.NewIdent(u.domain.GetMainModel().Name),
										},
									},
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
									Sel: ast.NewIdent(u.domain.GetRepositoryPrivateVariableName()),
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
									&ast.CompositeLit{
										Type: &ast.SelectorExpr{
											X:   ast.NewIdent("entities"),
											Sel: ast.NewIdent(u.domain.GetMainModel().Name),
										},
									},
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
										X: ast.NewIdent("u"),
										Sel: ast.NewIdent(
											u.domain.GetRepositoryPrivateVariableName(),
										),
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
									&ast.CompositeLit{
										Type: &ast.SelectorExpr{
											X:   ast.NewIdent("entities"),
											Sel: ast.NewIdent(u.domain.GetMainModel().Name),
										},
									},
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

func (u ServiceGenerator) syncUpdateMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, u.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	method, methodExist := astfile.FindFunc(file, "Update")
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
										&ast.SelectorExpr{
											X:   ast.NewIdent("update"),
											Sel: ast.NewIdent(param.GetName()),
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

func (u ServiceGenerator) delete() *ast.FuncDecl {
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
						X: ast.NewIdent(u.domain.GetServiceTypeName()),
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
										X: ast.NewIdent("u"),
										Sel: ast.NewIdent(
											u.domain.GetRepositoryPrivateVariableName(),
										),
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

func (u ServiceGenerator) syncDeleteMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, u.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	method, methodExist := astfile.FindFunc(file, "Delete")
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

var destinationPath = "."
