package models

import (
	"bytes"
	"fmt"
	"github.com/018bf/creathor/internal/configs"
	"github.com/018bf/creathor/internal/generators"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Model struct {
	ModelConfig *configs.ModelConfig
}

func NewModel(modelConfig *configs.ModelConfig) *Model {
	return &Model{ModelConfig: modelConfig}
}

func (m *Model) filename() string {
	return filepath.Join("internal", "domain", "models", m.ModelConfig.FileName())
}

func (m *Model) params() []*ast.Field {
	fields := []*ast.Field{
		{
			Doc:   nil,
			Names: []*ast.Ident{ast.NewIdent("ID")},
			Type:  ast.NewIdent("UUID"),
			Tag: &ast.BasicLit{
				Kind:  token.STRING,
				Value: "`json:\"id\"`",
			},
			Comment: nil,
		},
		{
			Doc:   nil,
			Names: []*ast.Ident{ast.NewIdent("UpdatedAt")},
			Type: &ast.SelectorExpr{
				X:   ast.NewIdent("time"),
				Sel: ast.NewIdent("Time"),
			},
			Tag: &ast.BasicLit{
				Kind:  token.STRING,
				Value: "`json:\"updated_at\"`",
			},
			Comment: nil,
		},
		{
			Doc:   nil,
			Names: []*ast.Ident{ast.NewIdent("CreatedAt")},
			Type: &ast.SelectorExpr{
				X:   ast.NewIdent("time"),
				Sel: ast.NewIdent("Time"),
			},
			Tag: &ast.BasicLit{
				Kind:  token.STRING,
				Value: "`json:\"created_at\"`",
			},
			Comment: nil,
		},
	}
	for _, param := range m.ModelConfig.Params {
		fields = append(fields, &ast.Field{
			Doc:   nil,
			Names: []*ast.Ident{ast.NewIdent(param.GetName())},
			Type:  ast.NewIdent(param.Type),
			Tag: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf("`json:\"%s\"`", param.Tag()),
			},
			Comment: nil,
		})
	}
	return fields
}

func (m *Model) toValidate() []*ast.CallExpr {
	fields := []*ast.CallExpr{
		{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("validation"),
				Sel: ast.NewIdent("Field"),
			},
			Args: []ast.Expr{
				&ast.UnaryExpr{
					Op: token.AND,
					X: &ast.SelectorExpr{
						X:   ast.NewIdent("m"),
						Sel: ast.NewIdent("ID"),
					},
				},
				&ast.SelectorExpr{
					X:   ast.NewIdent("validation"),
					Sel: ast.NewIdent("Required"),
				},
				&ast.SelectorExpr{
					X:   ast.NewIdent("is"),
					Sel: ast.NewIdent("UUID"),
				},
			},
		},
		{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("validation"),
				Sel: ast.NewIdent("Field"),
			},
			Args: []ast.Expr{
				&ast.UnaryExpr{
					Op: token.AND,
					X: &ast.SelectorExpr{
						X:   ast.NewIdent("m"),
						Sel: ast.NewIdent("UpdatedAt"),
					},
				},
				&ast.SelectorExpr{
					X:   ast.NewIdent("validation"),
					Sel: ast.NewIdent("Required"),
				},
			},
		},
		{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("validation"),
				Sel: ast.NewIdent("Field"),
			},
			Args: []ast.Expr{
				&ast.UnaryExpr{
					Op: token.AND,
					X: &ast.SelectorExpr{
						X:   ast.NewIdent("m"),
						Sel: ast.NewIdent("CreatedAt"),
					},
				},
				&ast.SelectorExpr{
					X:   ast.NewIdent("validation"),
					Sel: ast.NewIdent("Required"),
				},
			},
		},
	}

	for _, param := range m.ModelConfig.Params {
		call := &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("validation"),
				Sel: ast.NewIdent("Field"),
			},
			Args: []ast.Expr{
				&ast.UnaryExpr{
					Op: token.AND,
					X: &ast.SelectorExpr{
						X:   ast.NewIdent("m"),
						Sel: ast.NewIdent(param.GetName()),
					},
				},
			},
		}
		if !strings.HasPrefix(param.Type, "*") {
			call.Args = append(call.Args, &ast.SelectorExpr{
				X:   ast.NewIdent("validation"),
				Sel: ast.NewIdent("Required"),
			})
		}
		if strings.ToLower(param.Type) == "uuid" {
			call.Args = append(call.Args, &ast.SelectorExpr{
				X:   ast.NewIdent("is"),
				Sel: ast.NewIdent("UUID"),
			})
		}
		if strings.Contains(strings.ToLower(param.GetName()), "email") {
			call.Args = append(call.Args, &ast.SelectorExpr{
				X:   ast.NewIdent("is"),
				Sel: ast.NewIdent("EmailFormat"),
			})
		}
		fields = append(fields, call)
	}
	return fields
}

func (m *Model) syncStruct() error {
	fileset := token.NewFileSet()
	filename := m.filename()
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var structureExists bool
	var structure *ast.TypeSpec
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok && t.Name.String() == m.ModelConfig.ModelName() {
			structure = t
			structureExists = true
			return false
		}
		return true
	})
	if structure == nil {
		structure = m.astStruct()
	}
	for _, param := range m.params() {
		ast.Inspect(structure, func(node ast.Node) bool {
			if st, ok := node.(*ast.StructType); ok && st.Fields != nil {
				for _, field := range st.Fields.List {
					for _, fieldName := range field.Names {
						for _, name := range param.Names {
							if fieldName.Name == name.String() {
								return false
							}
						}

					}
				}
				st.Fields.List = append(st.Fields.List, param)
				return true
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
	if err := os.WriteFile(filename, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (m *Model) astStruct() *ast.TypeSpec {
	return &ast.TypeSpec{
		Doc:        nil,
		Name:       ast.NewIdent(m.ModelConfig.ModelName()),
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
}

func (m *Model) astValidate() *ast.FuncDecl {
	exprs := []ast.Expr{
		ast.NewIdent("m"),
	}
	for _, expr := range m.toValidate() {
		exprs = append(exprs, expr)
	}
	return &ast.FuncDecl{
		Doc: nil,
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Doc: nil,
					Names: []*ast.Ident{
						ast.NewIdent("m"),
					},
					Type: &ast.StarExpr{
						Star: 0,
						X:    ast.NewIdent(m.ModelConfig.ModelName()),
					},
					Tag:     nil,
					Comment: nil,
				},
			},
		},
		Name: ast.NewIdent("Validate"),
		Type: &ast.FuncType{
			Func:       0,
			TypeParams: nil,
			Params:     nil,
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Doc:     nil,
						Names:   nil,
						Type:    ast.NewIdent("error"),
						Tag:     nil,
						Comment: nil,
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			Lbrace: 0,
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs:    []ast.Expr{ast.NewIdent("err")},
					TokPos: 0,
					Tok:    token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("validation"),
								Sel: ast.NewIdent("ValidateStruct"),
							},
							Args: exprs,
						},
					},
				},
				&ast.IfStmt{
					If:   0,
					Init: nil,
					Cond: &ast.BinaryExpr{
						X:     ast.NewIdent("err"),
						OpPos: 0,
						Op:    token.NEQ,
						Y:     ast.NewIdent("nil"),
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Return: 0,
								Results: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   ast.NewIdent("errs"),
											Sel: ast.NewIdent("FromValidationError"),
										},
										Lparen: 0,
										Args: []ast.Expr{
											ast.NewIdent("err"),
										},
										Ellipsis: 0,
										Rparen:   0,
									},
								},
							},
						},
					},
					Else: nil,
				},
				&ast.ReturnStmt{
					Return:  0,
					Results: []ast.Expr{ast.NewIdent("nil")},
				},
			},
			Rbrace: 0,
		},
	}
}

func (m *Model) syncValidate() error {
	fileset := token.NewFileSet()
	filename := m.filename()
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var validatorExists bool
	var validator *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if fun, ok := node.(*ast.FuncDecl); ok {
			for _, field := range fun.Recv.List {
				if expr, ok := field.Type.(*ast.StarExpr); ok {
					ident, ok := expr.X.(*ast.Ident)
					if ok && ident.Name == m.ModelConfig.ModelName() {
						validator = fun
						validatorExists = true
						return false
					}
				}
			}
		}
		return true
	})
	if validator == nil {
		validator = m.astValidate()
	}
	for _, param := range m.toValidate() {
		ast.Inspect(validator.Body, func(node ast.Node) bool {
			if call, ok := node.(*ast.CallExpr); ok {
				if fun, ok := call.Fun.(*ast.SelectorExpr); ok && fun.Sel.Name == "ValidateStruct" {
					for _, arg := range call.Args {
						if vf, ok := arg.(*ast.CallExpr); ok {
							for _, vfArgs := range vf.Args {
								if field, ok := vfArgs.(*ast.UnaryExpr); ok {
									if fieldX, ok := field.X.(*ast.SelectorExpr); ok {
										for _, arg := range param.Args {
											if ue, ok := arg.(*ast.UnaryExpr); ok {
												if sel, ok := ue.X.(*ast.SelectorExpr); ok {
													if fieldX.Sel.String() == sel.Sel.String() {
														return false
													}
												}
											}
										}
									}
								}
							}
						}
					}
					call.Args = append(call.Args, param)
					return false
				}
			}
			return true
		})
	}
	if !validatorExists {
		file.Decls = append(file.Decls, validator)
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

func (m *Model) astFakeValues() []*ast.KeyValueExpr {
	kvs := []*ast.KeyValueExpr{
		{
			Key:   ast.NewIdent("ID"),
			Value: generators.FakeAst("UUID"),
		},
		{
			Key:   ast.NewIdent("UpdatedAt"),
			Value: generators.FakeAst("time.Time"),
		},
		{
			Key:   ast.NewIdent("CreatedAt"),
			Value: generators.FakeAst("time.Time"),
		},
	}
	for _, param := range m.ModelConfig.Params {
		kvs = append(kvs, &ast.KeyValueExpr{
			Key:   ast.NewIdent(param.GetName()),
			Colon: 0,
			Value: generators.FakeAst(param.Type),
		})
	}
	return kvs
}

func (m *Model) syncMock() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "domain", "models", "mock", m.ModelConfig.FileName())
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	mockName := fmt.Sprintf("New%s", m.ModelConfig.ModelName())
	var mockExists bool
	var mock *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if fun, ok := node.(*ast.FuncDecl); ok && fun.Name.Name == mockName {
			mock = fun
			mockExists = true
			return false
		}
		return true
	})
	if mock == nil {
		mock = &ast.FuncDecl{
			Doc:  nil,
			Recv: nil,
			Name: ast.NewIdent(mockName),
			Type: &ast.FuncType{
				Func:       0,
				TypeParams: nil,
				Params: &ast.FieldList{
					Opening: 0,
					List: []*ast.Field{
						{
							Doc:   nil,
							Names: []*ast.Ident{ast.NewIdent("t")},
							Type: &ast.StarExpr{
								Star: 0,
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("testing"),
									Sel: ast.NewIdent("T"),
								},
							},
							Tag:     nil,
							Comment: nil,
						},
					},
					Closing: 0,
				},
				Results: &ast.FieldList{
					Opening: 0,
					List: []*ast.Field{
						{
							Doc:   nil,
							Names: nil,
							Type: &ast.StarExpr{
								Star: 0,
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("models"),
									Sel: ast.NewIdent(m.ModelConfig.ModelName()),
								},
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
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("t"),
								Sel: ast.NewIdent("Helper"),
							},
							Lparen:   0,
							Args:     nil,
							Ellipsis: 0,
							Rparen:   0,
						},
					},
					&ast.AssignStmt{
						Lhs:    []ast.Expr{ast.NewIdent("m")},
						TokPos: 0,
						Tok:    token.DEFINE,
						Rhs: []ast.Expr{
							&ast.UnaryExpr{
								OpPos: 0,
								Op:    token.AND,
								X: &ast.CompositeLit{
									Type: &ast.SelectorExpr{
										X:   ast.NewIdent("models"),
										Sel: ast.NewIdent(m.ModelConfig.ModelName()),
									},
									Lbrace:     0,
									Elts:       nil,
									Rbrace:     0,
									Incomplete: false,
								},
							},
						},
					},
					&ast.ReturnStmt{
						Return:  0,
						Results: []ast.Expr{ast.NewIdent("m")},
					},
				},
				Rbrace: 0,
			},
		}
	}
	if !mockExists {
		file.Decls = append(file.Decls, mock)
	}
	for _, param := range m.astFakeValues() {
		ast.Inspect(mock.Body, func(node ast.Node) bool {
			if cl, ok := node.(*ast.CompositeLit); ok {
				if sel, ok := cl.Type.(*ast.SelectorExpr); ok {
					if sel.Sel.Name != m.ModelConfig.ModelName() {
						return true
					}
				}
				for _, elt := range cl.Elts {
					if kve, ok := elt.(*ast.KeyValueExpr); ok {
						if ident, ok := kve.Key.(*ast.Ident); ok {
							if n, ok := param.Key.(*ast.Ident); ok && ident.String() == n.String() {
								return false
							}
						}
					}
				}
				cl.Elts = append(cl.Elts, param)
				return false
			}
			return true
		})
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

func (m *Model) Sync() error {
	if err := m.syncStruct(); err != nil {
		return err
	}
	if err := m.syncValidate(); err != nil {
		return err
	}
	if err := m.syncMock(); err != nil {
		return err
	}
	return nil
}
