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

	"github.com/mikalai-mitsin/creathor/internal/pkg/astfile"
	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
)

type UseCaseGenerator struct {
	domain configs.EntityConfig
}

func NewUseCaseGenerator(domain configs.EntityConfig) *UseCaseGenerator {
	return &UseCaseGenerator{domain: domain}
}

func (u UseCaseGenerator) Sync() error {
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

func (u UseCaseGenerator) filename() string {
	return filepath.Join(
		"internal",
		"app",
		u.domain.AppConfig.AppName(),
		"usecases",
		u.domain.DirName(),
		u.domain.FileName(),
	)
}

func (u UseCaseGenerator) structure() *ast.TypeSpec {
	fields := []*ast.Field{
		{
			Names: []*ast.Ident{ast.NewIdent(u.domain.GetServicePrivateVariableName())},
			Type:  ast.NewIdent(u.domain.GetServiceInterfaceName()),
		},
	}
	if u.domain.AppConfig.ProjectConfig.KafkaEnabled {
		fields = append(fields, &ast.Field{
			Names: []*ast.Ident{ast.NewIdent(u.domain.EventServicePrivateVariableName())},
			Type:  ast.NewIdent(u.domain.EventServiceInterfaceName()),
		})
	}
	fields = append(fields, &ast.Field{
		Names: []*ast.Ident{ast.NewIdent("dtxManager")},
		Type:  ast.NewIdent("dtxManager"),
	}, &ast.Field{
		Names: []*ast.Ident{ast.NewIdent("logger")},
		Type:  ast.NewIdent("logger"),
	})
	structure := &ast.TypeSpec{
		Name: ast.NewIdent(u.domain.GetUseCaseTypeName()),
		Type: &ast.StructType{
			Fields: &ast.FieldList{
				List: fields,
			},
		},
	}
	return structure
}

func (u UseCaseGenerator) syncStruct() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, u.filename(), nil, parser.ParseComments)
	if err != nil {
		file = u.file()
	}
	structure, structureExists := astfile.FindType(file, u.domain.GetUseCaseTypeName())
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

func (u UseCaseGenerator) constructor() *ast.FuncDecl {
	fields := []*ast.Field{
		{
			Names: []*ast.Ident{ast.NewIdent(u.domain.GetServicePrivateVariableName())},
			Type:  ast.NewIdent(u.domain.GetServiceInterfaceName()),
		},
	}
	if u.domain.AppConfig.ProjectConfig.KafkaEnabled {
		fields = append(fields, &ast.Field{
			Names: []*ast.Ident{ast.NewIdent(u.domain.EventServicePrivateVariableName())},
			Type:  ast.NewIdent(u.domain.EventServiceInterfaceName()),
		})
	}
	fields = append(fields, &ast.Field{
		Names: []*ast.Ident{ast.NewIdent("dtxManager")},
		Type:  ast.NewIdent("dtxManager"),
	}, &ast.Field{
		Names: []*ast.Ident{ast.NewIdent("logger")},
		Type:  ast.NewIdent("logger"),
	})
	exprs := []ast.Expr{
		&ast.KeyValueExpr{
			Key:   ast.NewIdent(u.domain.GetServicePrivateVariableName()),
			Value: ast.NewIdent(u.domain.GetServicePrivateVariableName()),
		},
	}
	if u.domain.AppConfig.ProjectConfig.KafkaEnabled {
		exprs = append(exprs, &ast.KeyValueExpr{
			Key:   ast.NewIdent(u.domain.EventServicePrivateVariableName()),
			Value: ast.NewIdent(u.domain.EventServicePrivateVariableName()),
		})
	}
	exprs = append(exprs, &ast.KeyValueExpr{
		Key:   ast.NewIdent("dtxManager"),
		Value: ast.NewIdent("dtxManager"),
	}, &ast.KeyValueExpr{
		Key:   ast.NewIdent("logger"),
		Value: ast.NewIdent("logger"),
	})
	constructor := &ast.FuncDecl{
		Name: ast.NewIdent(u.domain.GetUseCaseConstructorName()),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: fields,
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent(
							fmt.Sprintf("*%s", u.domain.GetUseCaseTypeName()),
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
								Type: ast.NewIdent(u.domain.GetUseCaseTypeName()),
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

func (u UseCaseGenerator) syncConstructor() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, u.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	constructor, constructorExists := astfile.FindFunc(file, u.domain.GetUseCaseConstructorName())
	if constructor == nil {
		constructor = u.constructor()
	}
	if !constructorExists {
		file.Decls = append(file.Decls, constructor)
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

func (u UseCaseGenerator) createMethod() *ast.FuncDecl {
	var body []ast.Stmt
	body = append(body,
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: "logger",
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "u",
							},
							Sel: &ast.Ident{
								Name: "logger",
							},
						},
						Sel: &ast.Ident{
							Name: "WithContext",
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
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: "tx",
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "u",
							},
							Sel: &ast.Ident{
								Name: "dtxManager",
							},
						},
						Sel: &ast.Ident{
							Name: "NewTx",
						},
					},
				},
			},
		},
		&ast.DeferStmt{
			Call: &ast.CallExpr{
				Fun: &ast.FuncLit{
					Type: &ast.FuncType{
						Params: &ast.FieldList{
							List: []*ast.Field{
								{
									Names: []*ast.Ident{
										{
											Name: "tx",
										},
									},
									Type: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "dtx",
										},
										Sel: &ast.Ident{
											Name: "TX",
										},
									},
								},
							},
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
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
													Name: "tx",
												},
												Sel: &ast.Ident{
													Name: "Rollback",
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
										&ast.ExprStmt{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "logger",
													},
													Sel: &ast.Ident{
														Name: "Error",
													},
												},
												Args: []ast.Expr{
													&ast.BasicLit{
														Kind:  token.STRING,
														Value: "\"cant rollback transaction\"",
													},
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X: &ast.Ident{
																Name: "log",
															},
															Sel: &ast.Ident{
																Name: "Error",
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
										},
									},
								},
							},
						},
					},
				},
				Args: []ast.Expr{
					&ast.Ident{
						Name: "tx",
					},
				},
			},
		},
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(u.domain.GetOneVariableName()),
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.SelectorExpr{
							X:   ast.NewIdent("u"),
							Sel: ast.NewIdent(u.domain.GetServicePrivateVariableName()),
						},
						Sel: ast.NewIdent("Create"),
					},
					Args: []ast.Expr{
						ast.NewIdent("ctx"),
						ast.NewIdent("tx"),
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
	)
	if u.domain.AppConfig.ProjectConfig.KafkaEnabled {
		body = append(body, &ast.IfStmt{
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
									Name: "u",
								},
								Sel: &ast.Ident{
									Name: u.domain.EventServicePrivateVariableName(),
								},
							},
							Sel: &ast.Ident{
								Name: "Send",
							},
						},
						Args: []ast.Expr{
							&ast.Ident{
								Name: "ctx",
							},
							ast.NewIdent("tx"),
							&ast.Ident{
								Name: u.domain.GetOneVariableName(),
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
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.CompositeLit{
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "entities",
									},
									Sel: &ast.Ident{
										Name: u.domain.GetMainModel().Name,
									},
								},
							},
							&ast.Ident{
								Name: "err",
							},
						},
					},
				},
			},
		})
	}
	body = append(body,
		// Commit transaction
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
								Name: "tx",
							},
							Sel: &ast.Ident{
								Name: "Commit",
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
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.CompositeLit{
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "entities",
									},
									Sel: &ast.Ident{
										Name: u.domain.GetMainModel().Name,
									},
								},
							},
							&ast.Ident{
								Name: "err",
							},
						},
					},
				},
			},
		},
		// Return created model and nil error
		&ast.ReturnStmt{
			Results: []ast.Expr{
				ast.NewIdent(u.domain.GetOneVariableName()),
				ast.NewIdent("nil"),
			},
		},
	)
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("u"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(u.domain.GetUseCaseTypeName()),
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
			List: body,
		},
	}
}

func (u UseCaseGenerator) syncCreateMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, u.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	method, methodExist := astfile.FindFunc(file, "Create")
	if method == nil {
		method = u.createMethod()
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

func (u UseCaseGenerator) astListMethod() *ast.FuncDecl {
	var body []ast.Stmt
	body = append(body,
		// Try to update model at use case
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(u.domain.GetManyVariableName()),
				ast.NewIdent("count"),
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.SelectorExpr{
							X:   ast.NewIdent("u"),
							Sel: ast.NewIdent(u.domain.GetServicePrivateVariableName()),
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
				ast.NewIdent(u.domain.GetManyVariableName()),
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
						ast.NewIdent("u"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(u.domain.GetUseCaseTypeName()),
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
			List: body,
		},
	}
}

func (u UseCaseGenerator) syncListMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, u.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	method, methodExist := astfile.FindFunc(file, "List")
	if method == nil {
		method = u.astListMethod()
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

func (u UseCaseGenerator) astGetMethod() *ast.FuncDecl {
	var body []ast.Stmt
	body = append(
		body,
		// Try to get model from use case
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(u.domain.GetOneVariableName()),
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.SelectorExpr{
							X:   ast.NewIdent("u"),
							Sel: ast.NewIdent(u.domain.GetServicePrivateVariableName()),
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
	)
	body = append(
		body,
		// Return created model and nil error
		&ast.ReturnStmt{
			Results: []ast.Expr{
				ast.NewIdent(u.domain.GetOneVariableName()),
				ast.NewIdent("nil"),
			},
		},
	)
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("u"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(u.domain.GetUseCaseTypeName()),
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
			List: body,
		},
	}
}

func (u UseCaseGenerator) syncGetMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, u.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	method, methodExist := astfile.FindFunc(file, "Get")
	if method == nil {
		method = u.astGetMethod()
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

func (u UseCaseGenerator) updateMethod() *ast.FuncDecl {
	var body []ast.Stmt
	body = append(body,
		// Setup logger
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: "logger",
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "u",
							},
							Sel: &ast.Ident{
								Name: "logger",
							},
						},
						Sel: &ast.Ident{
							Name: "WithContext",
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
		// Begin transaction
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: "tx",
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "u",
							},
							Sel: &ast.Ident{
								Name: "dtxManager",
							},
						},
						Sel: &ast.Ident{
							Name: "NewTx",
						},
					},
				},
			},
		},
		// Defer rollback transaction
		&ast.DeferStmt{
			Call: &ast.CallExpr{
				Fun: &ast.FuncLit{
					Type: &ast.FuncType{
						Params: &ast.FieldList{
							List: []*ast.Field{
								{
									Names: []*ast.Ident{
										{
											Name: "tx",
										},
									},
									Type: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "dtx",
										},
										Sel: &ast.Ident{
											Name: "TX",
										},
									},
								},
							},
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
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
													Name: "tx",
												},
												Sel: &ast.Ident{
													Name: "Rollback",
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
										&ast.ExprStmt{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "logger",
													},
													Sel: &ast.Ident{
														Name: "Error",
													},
												},
												Args: []ast.Expr{
													&ast.BasicLit{
														Kind:  token.STRING,
														Value: "\"cant rollback transaction\"",
													},
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X: &ast.Ident{
																Name: "log",
															},
															Sel: &ast.Ident{
																Name: "Error",
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
										},
									},
								},
							},
						},
					},
				},
				Args: []ast.Expr{
					&ast.Ident{
						Name: "tx",
					},
				},
			},
		},
		// Try to update model at use case
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(u.domain.GetOneVariableName()),
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.SelectorExpr{
							X:   ast.NewIdent("u"),
							Sel: ast.NewIdent(u.domain.GetServicePrivateVariableName()),
						},
						Sel: ast.NewIdent("Update"),
					},
					Args: []ast.Expr{
						ast.NewIdent("ctx"),
						ast.NewIdent("tx"),
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
	)
	if u.domain.AppConfig.ProjectConfig.KafkaEnabled {
		body = append(body, &ast.IfStmt{
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
									Name: "u",
								},
								Sel: &ast.Ident{
									Name: u.domain.EventServicePrivateVariableName(),
								},
							},
							Sel: &ast.Ident{
								Name: "Send",
							},
						},
						Args: []ast.Expr{
							&ast.Ident{
								Name: "ctx",
							},
							ast.NewIdent("tx"),
							&ast.Ident{
								Name: u.domain.GetOneVariableName(),
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
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.CompositeLit{
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "entities",
									},
									Sel: &ast.Ident{
										Name: u.domain.GetMainModel().Name,
									},
								},
							},
							&ast.Ident{
								Name: "err",
							},
						},
					},
				},
			},
		})
	}
	body = append(body,
		// Commit transaction
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
								Name: "tx",
							},
							Sel: &ast.Ident{
								Name: "Commit",
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
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.CompositeLit{
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "entities",
									},
									Sel: &ast.Ident{
										Name: u.domain.GetMainModel().Name,
									},
								},
							},
							&ast.Ident{
								Name: "err",
							},
						},
					},
				},
			},
		},
		// Return created model and nil error
		&ast.ReturnStmt{
			Results: []ast.Expr{
				ast.NewIdent(u.domain.GetOneVariableName()),
				ast.NewIdent("nil"),
			},
		},
	)
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("u"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(u.domain.GetUseCaseTypeName()),
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
			List: body,
		},
	}
}

func (u UseCaseGenerator) syncUpdateMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, u.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	method, methodExist := astfile.FindFunc(file, "Update")
	if method == nil {
		method = u.updateMethod()
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

func (u UseCaseGenerator) deleteMethod() *ast.FuncDecl {
	var body []ast.Stmt
	body = append(body,
		// Setup logger
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: "logger",
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "u",
							},
							Sel: &ast.Ident{
								Name: "logger",
							},
						},
						Sel: &ast.Ident{
							Name: "WithContext",
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
		// Begin transaction
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: "tx",
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "u",
							},
							Sel: &ast.Ident{
								Name: "dtxManager",
							},
						},
						Sel: &ast.Ident{
							Name: "NewTx",
						},
					},
				},
			},
		},
		// Defer rollback transaction
		&ast.DeferStmt{
			Call: &ast.CallExpr{
				Fun: &ast.FuncLit{
					Type: &ast.FuncType{
						Params: &ast.FieldList{
							List: []*ast.Field{
								{
									Names: []*ast.Ident{
										{
											Name: "tx",
										},
									},
									Type: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "dtx",
										},
										Sel: &ast.Ident{
											Name: "TX",
										},
									},
								},
							},
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
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
													Name: "tx",
												},
												Sel: &ast.Ident{
													Name: "Rollback",
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
										&ast.ExprStmt{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "logger",
													},
													Sel: &ast.Ident{
														Name: "Error",
													},
												},
												Args: []ast.Expr{
													&ast.BasicLit{
														Kind:  token.STRING,
														Value: "\"cant rollback transaction\"",
													},
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X: &ast.Ident{
																Name: "log",
															},
															Sel: &ast.Ident{
																Name: "Error",
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
										},
									},
								},
							},
						},
					},
				},
				Args: []ast.Expr{
					&ast.Ident{
						Name: "tx",
					},
				},
			},
		},
		// Try to delete model at use case
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(u.domain.GetOneVariableName()),
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.SelectorExpr{
							X:   ast.NewIdent("u"),
							Sel: ast.NewIdent(u.domain.GetServicePrivateVariableName()),
						},
						Sel: ast.NewIdent("Delete"),
					},
					Args: []ast.Expr{
						ast.NewIdent("ctx"),
						ast.NewIdent("tx"),
						ast.NewIdent(u.domain.GetDeleteModel().Variable),
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
	)
	if u.domain.AppConfig.ProjectConfig.KafkaEnabled {
		body = append(body, &ast.IfStmt{
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
									Name: "u",
								},
								Sel: &ast.Ident{
									Name: u.domain.EventServicePrivateVariableName(),
								},
							},
							Sel: &ast.Ident{
								Name: "Send",
							},
						},
						Args: []ast.Expr{
							&ast.Ident{
								Name: "ctx",
							},
							ast.NewIdent("tx"),
							&ast.Ident{
								Name: u.domain.GetOneVariableName(),
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
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.CompositeLit{
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("entities"),
									Sel: ast.NewIdent(u.domain.GetMainModel().Name),
								},
							},
							&ast.Ident{
								Name: "err",
							},
						},
					},
				},
			},
		})
	}
	body = append(body,
		// Commit transaction
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
								Name: "tx",
							},
							Sel: &ast.Ident{
								Name: "Commit",
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
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.CompositeLit{
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("entities"),
									Sel: ast.NewIdent(u.domain.GetMainModel().Name),
								},
							},
							&ast.Ident{
								Name: "err",
							},
						},
					},
				},
			},
		},
		// Return created model and nil error
		&ast.ReturnStmt{
			Results: []ast.Expr{
				ast.NewIdent(u.domain.GetOneVariableName()),
				ast.NewIdent("nil"),
			},
		},
	)
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("u"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(u.domain.GetUseCaseTypeName()),
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
						Names: []*ast.Ident{ast.NewIdent(u.domain.GetDeleteModel().Variable)},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("entities"),
							Sel: ast.NewIdent(u.domain.GetDeleteModel().Name),
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
			List: body,
		},
	}
}

func (u UseCaseGenerator) syncDeleteMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, u.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	method, methodExist := astfile.FindFunc(file, "Delete")
	if method == nil {
		method = u.deleteMethod()
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

func (u UseCaseGenerator) file() *ast.File {
	specs := []ast.Spec{
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
				Value: u.domain.AppConfig.ProjectConfig.UUIDImportPath(),
			},
		},
	}
	return &ast.File{
		Name: ast.NewIdent("usecases"),
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok:   token.IMPORT,
				Specs: specs,
			},
		},
	}
}

var destinationPath = "."
