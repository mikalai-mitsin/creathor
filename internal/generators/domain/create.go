package domain

import (
	"bytes"
	"fmt"
	"github.com/018bf/creathor/internal/configs"
	"github.com/018bf/creathor/internal/fake"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type CreateModel struct {
	model *configs.ModelConfig
}

func NewCreateModel(modelConfig *configs.ModelConfig) *CreateModel {
	return &CreateModel{model: modelConfig}
}

func (m *CreateModel) filename() string {
	return filepath.Join("internal", "domain", "models", m.model.FileName())
}

func (m *CreateModel) file() *ast.File {
	return &ast.File{
		Package: 1,
		Name: &ast.Ident{
			Name: "models",
		},
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"time"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/domain/errs"`, m.model.Module),
						},
					},
					&ast.ImportSpec{
						Name: ast.NewIdent("validation"),
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"github.com/go-ozzo/ozzo-validation/v4"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"github.com/go-ozzo/ozzo-validation/v4/is"`,
						},
					},
				},
			},
		},
		Imports:  nil,
		Comments: nil,
	}
}

func (m *CreateModel) params() []*ast.Field {
	var fields []*ast.Field
	for _, param := range m.model.Params {
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

func (m *CreateModel) toValidate() []*ast.CallExpr {
	var fields []*ast.CallExpr
	for _, param := range m.model.Params {
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

func (m *CreateModel) syncStruct() error {
	fileset := token.NewFileSet()
	filename := m.filename()
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = m.file()
	}
	var structureExists bool
	var structure *ast.TypeSpec
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok && t.Name.String() == m.model.CreateTypeName() {
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

func (m *CreateModel) astStruct() *ast.TypeSpec {
	return &ast.TypeSpec{
		Name:       ast.NewIdent(m.model.CreateTypeName()),
		Doc:        nil,
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

func (m *CreateModel) astValidate() *ast.FuncDecl {
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
						X:    ast.NewIdent(m.model.CreateTypeName()),
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

func (m *CreateModel) syncValidate() error {
	fileset := token.NewFileSet()
	filename := m.filename()
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = m.file()
	}
	var validatorExists bool
	var validator *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if fun, ok := node.(*ast.FuncDecl); ok {
			for _, field := range fun.Recv.List {
				if expr, ok := field.Type.(*ast.StarExpr); ok {
					ident, ok := expr.X.(*ast.Ident)
					if ok && ident.Name == m.model.CreateTypeName() {
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

func (m *CreateModel) astFakeValues() []*ast.KeyValueExpr {
	var kvs []*ast.KeyValueExpr
	for _, param := range m.model.Params {
		kvs = append(kvs, &ast.KeyValueExpr{
			Key:   ast.NewIdent(param.GetName()),
			Colon: 0,
			Value: fake.FakeAst(param.Type),
		})
	}
	return kvs
}

func (m *CreateModel) mockFile() *ast.File {
	return &ast.File{
		Package: 1,
		Name: &ast.Ident{
			Name: "mock_models",
		},
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/domain/models"`, m.model.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/pkg/utils"`, m.model.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"github.com/google/uuid"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"github.com/jaswdr/faker"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"testing"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"time"`,
						},
					},
				},
			},
		},
		Imports:  nil,
		Comments: nil,
	}
}

func (m *CreateModel) syncMock() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "domain", "models", "mock", m.model.FileName())
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = m.mockFile()
	}
	mockName := fmt.Sprintf("New%s", m.model.CreateTypeName())
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
									Sel: ast.NewIdent(m.model.CreateTypeName()),
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
										Sel: ast.NewIdent(m.model.CreateTypeName()),
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
					if sel.Sel.Name != m.model.CreateTypeName() {
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

func (m *CreateModel) Sync() error {
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