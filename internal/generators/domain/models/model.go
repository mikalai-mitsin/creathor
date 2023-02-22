package models

import (
	"bytes"
	"fmt"
	"github.com/018bf/creathor/internal/generators"
	"github.com/018bf/creathor/internal/models"
	"github.com/iancoleman/strcase"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"
	"strings"
)

type Param struct {
	Name string
	Type string
}

func (p Param) GetPublicName() string {
	if strings.ToLower(p.Name) == "id" {
		return "ID"
	}
	return strcase.ToCamel(p.Name)
}

func (p Param) GetPrivateName() string {
	if strings.ToLower(p.Name) == "id" {
		return "ID"
	}
	return strcase.ToLowerCamel(p.Name)
}

func (p Param) GetTag() string {
	return strcase.ToSnake(p.Name)
}

type Model struct {
	Name        string
	ModelConfig *models.ModelConfig
	Params      []*Param
}

func (m *Model) SyncStruct() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "domain", "models", m.ModelConfig.FileName())
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var structureExists bool
	var structure *ast.TypeSpec
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok && t.Name.String() == m.Name {
			structure = t
			structureExists = true
			return false
		}
		return true
	})
	if structure == nil {
		structure = &ast.TypeSpec{
			Doc:        nil,
			Name:       ast.NewIdent(m.Name),
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
	for _, param := range m.Params {
		ast.Inspect(structure, func(node ast.Node) bool {
			if st, ok := node.(*ast.StructType); ok && st.Fields != nil {
				for _, field := range st.Fields.List {
					for _, fieldName := range field.Names {
						if fieldName.Name == param.GetPublicName() {
							return false
						}
					}
				}
				st.Fields.List = append(st.Fields.List, &ast.Field{
					Doc:   nil,
					Names: []*ast.Ident{ast.NewIdent(param.GetPublicName())},
					Type:  ast.NewIdent(param.Type),
					Tag: &ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf("`json:\"%s\"`", param.GetTag()),
					},
					Comment: nil,
				})
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

func (m *Model) SyncValidate() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "domain", "models", m.ModelConfig.FileName())
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
					if ok && ident.Name == m.Name {
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
		validator = &ast.FuncDecl{
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
							X:    ast.NewIdent(m.Name),
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
								Args: []ast.Expr{
									ast.NewIdent("m"),
								},
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
	for _, param := range m.Params {
		ast.Inspect(validator.Body, func(node ast.Node) bool {
			if call, ok := node.(*ast.CallExpr); ok {
				if fun, ok := call.Fun.(*ast.SelectorExpr); ok && fun.Sel.Name == "ValidateStruct" {
					for _, arg := range call.Args {
						if vf, ok := arg.(*ast.CallExpr); ok {
							for _, vfArgs := range vf.Args {
								if field, ok := vfArgs.(*ast.UnaryExpr); ok {
									if fieldX, ok := field.X.(*ast.SelectorExpr); ok {
										if fieldX.Sel.Name == param.GetPublicName() {
											return false
										}
									}
								}
							}
						}
					}
					vc := &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("validation"),
							Sel: ast.NewIdent("Field"),
						},
						Args: []ast.Expr{
							&ast.UnaryExpr{
								Op: token.AND,
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("m"),
									Sel: ast.NewIdent(param.GetPublicName()),
								},
							},
						},
					}
					if !strings.HasPrefix(param.Type, "*") {
						vc.Args = append(vc.Args, &ast.SelectorExpr{
							X:   ast.NewIdent("validation"),
							Sel: ast.NewIdent("Required"),
						})
					}
					if strings.ToLower(param.Type) == "uuid" {
						vc.Args = append(vc.Args, &ast.SelectorExpr{
							X:   ast.NewIdent("is"),
							Sel: ast.NewIdent("UUID"),
						})
					}
					if strings.Contains(strings.ToLower(param.GetPublicName()), "email") {
						vc.Args = append(vc.Args, &ast.SelectorExpr{
							X:   ast.NewIdent("is"),
							Sel: ast.NewIdent("EmailFormat"),
						})
					}
					call.Args = append(call.Args, vc)
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

func (m *Model) SyncMock() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "domain", "models", "mock", m.ModelConfig.FileName())
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	mockName := fmt.Sprintf("New%s", m.Name)
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
									Sel: ast.NewIdent(m.Name),
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
										Sel: ast.NewIdent(m.Name),
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
	for _, param := range m.Params {
		ast.Inspect(mock.Body, func(node ast.Node) bool {
			if cl, ok := node.(*ast.CompositeLit); ok {
				if sel, ok := cl.Type.(*ast.SelectorExpr); ok {
					if sel.Sel.Name != m.Name {
						return true
					}
				}
				for _, elt := range cl.Elts {
					if kve, ok := elt.(*ast.KeyValueExpr); ok {
						if ident, ok := kve.Key.(*ast.Ident); ok && ident.Name == param.GetPublicName() {
							return false
						}
					}
				}
				cl.Elts = append(cl.Elts, &ast.KeyValueExpr{
					Key:   ast.NewIdent(param.GetPublicName()),
					Colon: 0,
					Value: generators.FakeAst(param.Type),
				})
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
	if err := m.SyncStruct(); err != nil {
		return err
	}
	if err := m.SyncValidate(); err != nil {
		return err
	}
	if err := m.SyncMock(); err != nil {
		return err
	}
	return nil
}
