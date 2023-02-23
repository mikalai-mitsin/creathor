package interceptors

import (
	"bytes"
	"fmt"
	"github.com/018bf/creathor/internal/configs"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
)

type Interceptor struct {
	model *configs.ModelConfig
}

func NewInterceptor(model *configs.ModelConfig) *Interceptor {
	return &Interceptor{model: model}
}

func (i Interceptor) Sync() error {
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

func (i Interceptor) filename() string {
	return filepath.Join("internal", "interceptors", i.model.FileName())
}

func (i Interceptor) astStruct() *ast.TypeSpec {
	fields := []*ast.Field{
		{
			Doc:   nil,
			Names: []*ast.Ident{ast.NewIdent(i.model.UseCaseVariableName())},
			Type: &ast.SelectorExpr{
				X:   ast.NewIdent("usecases"),
				Sel: ast.NewIdent(i.model.UseCaseTypeName()),
			},
			Tag:     nil,
			Comment: nil,
		},
		{
			Doc:   nil,
			Names: []*ast.Ident{ast.NewIdent("logger")},
			Type: &ast.SelectorExpr{
				X:   ast.NewIdent("log"),
				Sel: ast.NewIdent("Logger"),
			},
			Tag:     nil,
			Comment: nil,
		},
	}
	if i.model.Auth {
		fields = append(
			fields,
			&ast.Field{
				Names: []*ast.Ident{ast.NewIdent("authUseCase")},
				Type: &ast.SelectorExpr{
					X:   ast.NewIdent("usecases"),
					Sel: ast.NewIdent("AuthUseCase"),
				},
			},
		)
	}
	structure := &ast.TypeSpec{
		Name: ast.NewIdent(i.model.InterceptorTypeName()),
		Type: &ast.StructType{
			Fields: &ast.FieldList{
				List: fields,
			},
		},
	}
	return structure
}

func (i Interceptor) syncStruct() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, i.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
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
		structure = i.astStruct()
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
	if err := os.WriteFile(i.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (i Interceptor) astConstructor() *ast.FuncDecl {
	fields := []*ast.Field{
		{
			Doc:   nil,
			Names: []*ast.Ident{ast.NewIdent(i.model.UseCaseVariableName())},
			Type: &ast.SelectorExpr{
				X:   ast.NewIdent("usecases"),
				Sel: ast.NewIdent(i.model.UseCaseTypeName()),
			},
			Tag:     nil,
			Comment: nil,
		},
		{
			Doc:   nil,
			Names: []*ast.Ident{ast.NewIdent("logger")},
			Type: &ast.SelectorExpr{
				X:   ast.NewIdent("log"),
				Sel: ast.NewIdent("Logger"),
			},
			Tag:     nil,
			Comment: nil,
		},
	}
	exprs := []ast.Expr{
		&ast.KeyValueExpr{
			Key: &ast.Ident{
				Name: i.model.UseCaseVariableName(),
			},
			Value: &ast.Ident{
				Name: i.model.UseCaseVariableName(),
			},
		},
		&ast.KeyValueExpr{
			Key: &ast.Ident{
				Name: "logger",
			},
			Value: &ast.Ident{
				Name: "logger",
			},
		},
	}
	if i.model.Auth {
		fields = append(
			fields,
			&ast.Field{
				Names: []*ast.Ident{ast.NewIdent("authUseCase")},
				Type: &ast.SelectorExpr{
					X:   ast.NewIdent("usecases"),
					Sel: ast.NewIdent("AuthUseCase"),
				},
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
		Doc:  nil,
		Recv: nil,
		Name: ast.NewIdent(fmt.Sprintf("New%s", i.model.InterceptorTypeName())),
		Type: &ast.FuncType{
			Func:       0,
			TypeParams: nil,
			Params: &ast.FieldList{
				Opening: 0,
				List:    fields,
				Closing: 0,
			},
			Results: &ast.FieldList{
				Opening: 0,
				List: []*ast.Field{
					{
						Doc:     nil,
						Names:   nil,
						Type:    ast.NewIdent(fmt.Sprintf("interceptors.%s", i.model.InterceptorTypeName())),
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
							X: &ast.CompositeLit{
								Type:       ast.NewIdent(i.model.InterceptorTypeName()),
								Lbrace:     0,
								Elts:       exprs,
								Rbrace:     0,
								Incomplete: false,
							},
						},
					},
				},
			},
			Rbrace: 0,
		},
	}
	return constructor
}

func (i Interceptor) syncConstructor() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, i.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var structureConstructorExists bool
	var structureConstructor *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == fmt.Sprintf("New%s", i.model.InterceptorTypeName()) {
			structureConstructorExists = true
			structureConstructor = t
			return false
		}
		return true
	})
	if structureConstructor == nil {
		structureConstructor = i.astConstructor()
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

func (i Interceptor) astCreateMethod() *ast.FuncDecl {
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
						X: ast.NewIdent(i.model.InterceptorTypeName()),
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
								Sel: ast.NewIdent(i.model.CreateTypeName()),
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
								Sel: ast.NewIdent(i.model.ModelName()),
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
										Sel: ast.NewIdent(i.model.PermissionIDCreate()),
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
										Sel: ast.NewIdent(i.model.PermissionIDCreate()),
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
						ast.NewIdent(i.model.Variable()),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("i"),
									Sel: ast.NewIdent(i.model.UseCaseVariableName()),
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
						ast.NewIdent(i.model.Variable()),
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
}

func (i Interceptor) syncCreateMethod() error {
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
		method = i.astCreateMethod()
	}
	for _, param := range i.model.Params {
		param := param
		ast.Inspect(method, func(node ast.Node) bool {
			if cl, ok := node.(*ast.CompositeLit); ok {
				if t, ok := cl.Type.(*ast.SelectorExpr); ok && t.Sel.String() == i.model.ModelName() {
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
	if err := os.WriteFile(i.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (i Interceptor) astListMethod() *ast.FuncDecl {
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
						X: ast.NewIdent(i.model.InterceptorTypeName()),
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
								Sel: ast.NewIdent(i.model.FilterTypeName()),
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
									Sel: ast.NewIdent(i.model.ModelName()),
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
										Sel: ast.NewIdent(i.model.PermissionIDList()),
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
										Sel: ast.NewIdent(i.model.PermissionIDList()),
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
						ast.NewIdent(i.model.ListVariable()),
						ast.NewIdent("count"),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("i"),
									Sel: ast.NewIdent(i.model.UseCaseVariableName()),
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
						ast.NewIdent(i.model.ListVariable()),
						ast.NewIdent("count"),
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
}

func (i Interceptor) syncListMethod() error {
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

func (i Interceptor) astGetMethod() *ast.FuncDecl {
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
						X: ast.NewIdent(i.model.InterceptorTypeName()),
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
								Sel: ast.NewIdent(i.model.ModelName()),
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
										Sel: ast.NewIdent(i.model.PermissionIDDetail()),
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
						ast.NewIdent(i.model.Variable()),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("i"),
									Sel: ast.NewIdent(i.model.UseCaseVariableName()),
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
										Sel: ast.NewIdent(i.model.PermissionIDDetail()),
									},
									ast.NewIdent(i.model.Variable()),
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
						ast.NewIdent(i.model.Variable()),
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
}

func (i Interceptor) syncGetMethod() error {
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

func (i Interceptor) astUpdateMethod() *ast.FuncDecl {
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
						X: ast.NewIdent(i.model.InterceptorTypeName()),
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
								Sel: ast.NewIdent(i.model.UpdateTypeName()),
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
								Sel: ast.NewIdent(i.model.ModelName()),
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
										Sel: ast.NewIdent(i.model.PermissionIDUpdate()),
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
						ast.NewIdent(i.model.Variable()),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("i"),
									Sel: ast.NewIdent(i.model.UseCaseVariableName()),
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
										Sel: ast.NewIdent(i.model.PermissionIDUpdate()),
									},
									ast.NewIdent(i.model.Variable()),
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
									Sel: ast.NewIdent(i.model.UseCaseVariableName()),
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

func (i Interceptor) syncUpdateMethod() error {
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
		method = i.astUpdateMethod()
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

func (i Interceptor) astDeleteMethod() *ast.FuncDecl {
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
						X: ast.NewIdent(i.model.InterceptorTypeName()),
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
										Sel: ast.NewIdent(i.model.PermissionIDDelete()),
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
						ast.NewIdent(i.model.Variable()),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("i"),
									Sel: ast.NewIdent(i.model.UseCaseVariableName()),
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
										Sel: ast.NewIdent(i.model.PermissionIDDelete()),
									},
									ast.NewIdent(i.model.Variable()),
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
										Sel: ast.NewIdent(i.model.UseCaseVariableName()),
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

func (i Interceptor) syncDeleteMethod() error {
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
		method = i.astDeleteMethod()
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
