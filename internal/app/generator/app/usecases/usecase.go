package usecases

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

	"github.com/mikalai-mitsin/creathor/internal/pkg/domain"
)

type UseCaseGenerator struct {
	domain *domain.Domain
}

func NewUseCaseGenerator(domain *domain.Domain) *UseCaseGenerator {
	return &UseCaseGenerator{domain: domain}
}

func (i UseCaseGenerator) Sync() error {
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
	if err := i.syncCreateMethod(); err != nil {
		return err
	}
	if err := i.syncGetMethod(); err != nil {
		return err
	}
	if err := i.syncListMethod(); err != nil {
		return err
	}
	if err := i.syncUpdateMethod(); err != nil {
		return err
	}
	if err := i.syncDeleteMethod(); err != nil {
		return err
	}
	return nil
}

func (i UseCaseGenerator) filename() string {
	return filepath.Join("internal", "app", i.domain.DirName(), "usecases", i.domain.FileName())
}

func (i UseCaseGenerator) structure() *ast.TypeSpec {
	fields := []*ast.Field{
		{
			Names: []*ast.Ident{ast.NewIdent(i.domain.GetServicePrivateVariableName())},
			Type:  ast.NewIdent(i.domain.GetServiceTypeName()),
		},
		{
			Names: []*ast.Ident{ast.NewIdent("logger")},
			Type:  ast.NewIdent("logger"),
		},
	}
	structure := &ast.TypeSpec{
		Name: ast.NewIdent(i.domain.GetUseCaseTypeName()),
		Type: &ast.StructType{
			Fields: &ast.FieldList{
				List: fields,
			},
		},
	}
	return structure
}

func (i UseCaseGenerator) syncStruct() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, i.filename(), nil, parser.ParseComments)
	if err != nil {
		file = i.file()
	}
	var structureExists bool
	var structure *ast.TypeSpec
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok && t.Name.String() == i.domain.GetUseCaseTypeName() {
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

func (i UseCaseGenerator) constructor() *ast.FuncDecl {
	fields := []*ast.Field{
		{
			Names: []*ast.Ident{ast.NewIdent(i.domain.GetServicePrivateVariableName())},
			Type:  ast.NewIdent(i.domain.GetServiceTypeName()),
		},
		{
			Names: []*ast.Ident{ast.NewIdent("logger")},
			Type:  ast.NewIdent("logger"),
		},
	}
	exprs := []ast.Expr{
		&ast.KeyValueExpr{
			Key:   ast.NewIdent(i.domain.GetServicePrivateVariableName()),
			Value: ast.NewIdent(i.domain.GetServicePrivateVariableName()),
		},
		&ast.KeyValueExpr{
			Key:   ast.NewIdent("logger"),
			Value: ast.NewIdent("logger"),
		},
	}
	constructor := &ast.FuncDecl{
		Name: ast.NewIdent(i.domain.GetUseCaseConstructorName()),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: fields,
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent(
							fmt.Sprintf("*%s", i.domain.GetUseCaseTypeName()),
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
								Type: ast.NewIdent(i.domain.GetUseCaseTypeName()),
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

func (i UseCaseGenerator) syncConstructor() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, i.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var structureConstructorExists bool
	var structureConstructor *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok &&
			t.Name.String() == i.domain.GetUseCaseConstructorName() {
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

func (i UseCaseGenerator) createMethod() *ast.FuncDecl {
	var body []ast.Stmt
	body = append(body,
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(i.domain.GetMainModel().Variable),
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.SelectorExpr{
							X:   ast.NewIdent("i"),
							Sel: ast.NewIdent(i.domain.GetServicePrivateVariableName()),
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
				ast.NewIdent(i.domain.GetMainModel().Variable),
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
						X: ast.NewIdent(i.domain.GetUseCaseTypeName()),
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
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("context"),
							Sel: ast.NewIdent("Context"),
						},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("create")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("entities"),
								Sel: ast.NewIdent(i.domain.GetCreateModel().Name),
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
								Sel: ast.NewIdent(i.domain.GetMainModel().Name),
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
			List: body,
		},
	}
}

func (i UseCaseGenerator) syncCreateMethod() error {
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
		method = i.createMethod()
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

func (i UseCaseGenerator) astListMethod() *ast.FuncDecl {
	var body []ast.Stmt
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
							Sel: ast.NewIdent(i.domain.GetServicePrivateVariableName()),
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
						X: ast.NewIdent(i.domain.GetUseCaseTypeName()),
					},
				},
			},
		},
		Name: ast.NewIdent("List"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("ctx")},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("context"),
							Sel: ast.NewIdent("Context"),
						},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("filter")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("entities"),
								Sel: ast.NewIdent(i.domain.GetFilterModel().Name),
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
									Sel: ast.NewIdent(i.domain.GetMainModel().Name),
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
			List: body,
		},
	}
}

func (i UseCaseGenerator) syncListMethod() error {
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
		method = i.astListMethod()
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

func (i UseCaseGenerator) astGetMethod() *ast.FuncDecl {
	var body []ast.Stmt
	body = append(
		body,
		// Try to get model from use case
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(i.domain.GetMainModel().Variable),
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.SelectorExpr{
							X:   ast.NewIdent("i"),
							Sel: ast.NewIdent(i.domain.GetServicePrivateVariableName()),
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
	body = append(
		body,
		// Return created model and nil error
		&ast.ReturnStmt{
			Results: []ast.Expr{
				ast.NewIdent(i.domain.GetMainModel().Variable),
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
						X: ast.NewIdent(i.domain.GetUseCaseTypeName()),
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
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("context"),
							Sel: ast.NewIdent("Context"),
						},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("id")},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("uuid"),
							Sel: ast.NewIdent("UUID"),
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
								Sel: ast.NewIdent(i.domain.GetMainModel().Name),
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
			List: body,
		},
	}
}

func (i UseCaseGenerator) syncGetMethod() error {
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
		method = i.astGetMethod()
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

func (i UseCaseGenerator) updateMethod() *ast.FuncDecl {
	var body []ast.Stmt
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
							Sel: ast.NewIdent(i.domain.GetServicePrivateVariableName()),
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
						X: ast.NewIdent(i.domain.GetUseCaseTypeName()),
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
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("context"),
							Sel: ast.NewIdent("Context"),
						},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("update")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("entities"),
								Sel: ast.NewIdent(i.domain.GetUpdateModel().Name),
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
								Sel: ast.NewIdent(i.domain.GetMainModel().Name),
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
			List: body,
		},
	}
}

func (i UseCaseGenerator) syncUpdateMethod() error {
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
		method = i.updateMethod()
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

func (i UseCaseGenerator) deleteMethod() *ast.FuncDecl {
	var body []ast.Stmt
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
								Sel: ast.NewIdent(i.domain.GetServicePrivateVariableName()),
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
						X: ast.NewIdent(i.domain.GetUseCaseTypeName()),
					},
				},
			},
		},
		Name: ast.NewIdent("Delete"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("ctx")},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("context"),
							Sel: ast.NewIdent("Context"),
						},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("id")},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("uuid"),
							Sel: ast.NewIdent("UUID"),
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
			List: body,
		},
	}
}

func (i UseCaseGenerator) syncDeleteMethod() error {
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
		method = i.deleteMethod()
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

func (i UseCaseGenerator) file() *ast.File {
	return &ast.File{
		Name: ast.NewIdent("usecases"),
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
							Value: i.domain.EntitiesImportPath(),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/pkg/uuid"`, i.domain.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/pkg/log"`, i.domain.Module),
						},
					},
					&ast.ImportSpec{
						Name: ast.NewIdent("userEntities"),
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/app/user/entities"`, i.domain.Module),
						},
					},
				},
			},
		},
	}
}

var destinationPath = "."
