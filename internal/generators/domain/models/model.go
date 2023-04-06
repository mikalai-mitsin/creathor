package models

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
	"strings"

	"github.com/018bf/creathor/internal/configs"
	"github.com/018bf/creathor/internal/fake"
)

type MainModel struct {
	model *configs.ModelConfig
}

func NewMainModel(modelConfig *configs.ModelConfig) *MainModel {
	return &MainModel{model: modelConfig}
}

func (m *MainModel) file() *ast.File {
	decls := []ast.Decl{
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
	}
	if m.model.Auth {
		decls = append(decls, &ast.GenDecl{
			Tok: token.CONST,
			Specs: []ast.Spec{
				&ast.ValueSpec{
					Names: []*ast.Ident{
						{
							Name: m.model.PermissionIDList(),
						},
					},
					Type: &ast.Ident{
						Name: "PermissionID",
					},
					Values: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s_list"`, m.model.KeyName()),
						},
					},
				},
				&ast.ValueSpec{
					Names: []*ast.Ident{
						{
							Name: m.model.PermissionIDDetail(),
						},
					},
					Type: &ast.Ident{
						Name: "PermissionID",
					},
					Values: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s_detail"`, m.model.KeyName()),
						},
					},
				},
				&ast.ValueSpec{
					Names: []*ast.Ident{
						{
							Name: m.model.PermissionIDCreate(),
						},
					},
					Type: &ast.Ident{
						Name: "PermissionID",
					},
					Values: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s_create"`, m.model.KeyName()),
						},
					},
				},
				&ast.ValueSpec{
					Names: []*ast.Ident{
						{
							Name: m.model.PermissionIDUpdate(),
						},
					},
					Type: &ast.Ident{
						Name: "PermissionID",
					},
					Values: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s_update"`, m.model.KeyName()),
						},
					},
				},
				&ast.ValueSpec{
					Names: []*ast.Ident{
						{
							Name: m.model.PermissionIDDelete(),
						},
					},
					Type: &ast.Ident{
						Name: "PermissionID",
					},
					Values: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s_delete"`, m.model.KeyName()),
						},
					},
				},
			},
		})
	}
	return &ast.File{
		Name:  ast.NewIdent("models"),
		Decls: decls,
	}
}

func (m *MainModel) filename() string {
	return filepath.Join("internal", "domain", "models", m.model.FileName())
}

func (m *MainModel) params() []*ast.Field {
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

func (m *MainModel) toValidate() []*ast.CallExpr {
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

func (m *MainModel) syncStruct() error {
	fileset := token.NewFileSet()
	filename := m.filename()
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = m.file()
	}
	var structureExists bool
	var structure *ast.TypeSpec
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok && t.Name.String() == m.model.ModelName() {
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

func (m *MainModel) astStruct() *ast.TypeSpec {
	return &ast.TypeSpec{
		Doc:        nil,
		Name:       ast.NewIdent(m.model.ModelName()),
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

func (m *MainModel) astValidate() *ast.FuncDecl {
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
						X:    ast.NewIdent(m.model.ModelName()),
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

func (m *MainModel) syncValidate() error {
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
					if ok && ident.Name == m.model.ModelName() {
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

func (m *MainModel) astFakeValues() []*ast.KeyValueExpr {
	kvs := []*ast.KeyValueExpr{
		{
			Key:   ast.NewIdent("ID"),
			Value: fake.FakeAst("UUID"),
		},
		{
			Key:   ast.NewIdent("UpdatedAt"),
			Value: fake.FakeAst("time.Time"),
		},
		{
			Key:   ast.NewIdent("CreatedAt"),
			Value: fake.FakeAst("time.Time"),
		},
	}
	for _, param := range m.model.Params {
		kvs = append(kvs, &ast.KeyValueExpr{
			Key:   ast.NewIdent(param.GetName()),
			Colon: 0,
			Value: fake.FakeAst(param.Type),
		})
	}
	return kvs
}

func (m *MainModel) mockFile() *ast.File {
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

func (m *MainModel) syncMock() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "domain", "models", "mock", m.model.FileName())
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = m.mockFile()
	}
	mockName := fmt.Sprintf("New%s", m.model.ModelName())
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
									Sel: ast.NewIdent(m.model.ModelName()),
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
										Sel: ast.NewIdent(m.model.ModelName()),
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
					if sel.Sel.Name != m.model.ModelName() {
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

func (m *MainModel) Sync() error {
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