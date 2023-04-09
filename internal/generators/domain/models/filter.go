package models

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"

	"github.com/018bf/creathor/internal/configs"
)

type FilterModel struct {
	model *configs.ModelConfig
}

func NewFilterModel(modelConfig *configs.ModelConfig) *FilterModel {
	return &FilterModel{model: modelConfig}
}

func (m *FilterModel) filename() string {
	return filepath.Join("internal", "domain", "models", m.model.FileName())
}

func (m *FilterModel) file() *ast.File {
	return &ast.File{
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

func (m *FilterModel) params() []*ast.Field {
	fields := []*ast.Field{
		{
			Names: []*ast.Ident{ast.NewIdent("IDs")},
			Type: &ast.ArrayType{
				Elt: ast.NewIdent("UUID"),
			},
			Tag: &ast.BasicLit{
				Kind:  token.STRING,
				Value: "`json:\"ids\"`",
			},
		},
		{
			Names: []*ast.Ident{ast.NewIdent("PageSize")},
			Type:  &ast.StarExpr{X: ast.NewIdent("uint64")},
			Tag: &ast.BasicLit{
				Kind:  token.STRING,
				Value: "`json:\"page_size\"`",
			},
		},
		{
			Names: []*ast.Ident{ast.NewIdent("PageNumber")},
			Type:  &ast.StarExpr{X: ast.NewIdent("uint64")},
			Tag: &ast.BasicLit{
				Kind:  token.STRING,
				Value: "`json:\"page_number\"`",
			},
		},
		{
			Names: []*ast.Ident{ast.NewIdent("OrderBy")},
			Type: &ast.ArrayType{
				Elt: ast.NewIdent("string"),
			},
			Tag: &ast.BasicLit{
				Kind:  token.STRING,
				Value: "`json:\"order_by\"`",
			},
		},
	}
	if m.model.SearchEnabled() {
		fields = append(fields, &ast.Field{
			Names: []*ast.Ident{ast.NewIdent("Search")},
			Type:  &ast.StarExpr{X: ast.NewIdent("string")},
			Tag: &ast.BasicLit{
				Kind:  token.STRING,
				Value: "`json:\"search\"`",
			},
		})
	}
	return fields
}

func (m *FilterModel) toValidate() []*ast.CallExpr {
	callExprs := []*ast.CallExpr{
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
						Sel: ast.NewIdent("IDs"),
					},
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
						Sel: ast.NewIdent("PageNumber"),
					},
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
						Sel: ast.NewIdent("PageSize"),
					},
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
						Sel: ast.NewIdent("OrderBy"),
					},
				},
			},
		},
	}
	if m.model.Auth {
		callExprs = append(callExprs, &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("validation"),
				Sel: ast.NewIdent("Field"),
			},
			Args: []ast.Expr{
				&ast.UnaryExpr{
					Op: token.AND,
					X: &ast.SelectorExpr{
						X:   ast.NewIdent("m"),
						Sel: ast.NewIdent("Search"),
					},
				},
			},
		})
	}
	return callExprs
}

func (m *FilterModel) syncStruct() error {
	fileset := token.NewFileSet()
	filename := m.filename()
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = m.file()
	}
	var structureExists bool
	var structure *ast.TypeSpec
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok && t.Name.String() == m.model.FilterTypeName() {
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

func (m *FilterModel) astStruct() *ast.TypeSpec {
	return &ast.TypeSpec{
		Name: ast.NewIdent(m.model.FilterTypeName()),
		Type: &ast.StructType{
			Fields: &ast.FieldList{
				List:    m.params(),
				Closing: 0,
			},
			Incomplete: false,
		},
		Comment: nil,
	}
}

func (m *FilterModel) astValidate() *ast.FuncDecl {
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
						X:    ast.NewIdent(m.model.FilterTypeName()),
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

func (m *FilterModel) syncValidate() error {
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
					if ok && ident.Name == m.model.FilterTypeName() {
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

func (m *FilterModel) Sync() error {
	if err := m.syncStruct(); err != nil {
		return err
	}
	if err := m.syncValidate(); err != nil {
		return err
	}
	mock := NewMock(m.astStruct(), m.model.FileName())
	if err := mock.Sync(); err != nil {
		return err
	}
	return nil
}
