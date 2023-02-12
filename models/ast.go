package models

import (
	"bytes"
	"fmt"
	"github.com/iancoleman/strcase"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"strings"
)

func SyncStruct(strc *Struct) error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, strc.Path, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var structureExists bool
	var structure *ast.TypeSpec
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok && t.Name.String() == strc.Name {
			structure = t
			structureExists = true
			return false
		}
		return true
	})
	if structure == nil {
		structure = &ast.TypeSpec{
			Doc:        nil,
			Name:       ast.NewIdent(strc.Name),
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
	for _, param := range strc.Params {
		ast.Inspect(structure, func(node ast.Node) bool {
			if st, ok := node.(*ast.StructType); ok && st.Fields != nil {
				for _, field := range st.Fields.List {
					for _, fieldName := range field.Names {
						if fieldName.Name == param.GetName() {
							return false
						}
					}
				}
				st.Fields.List = append(st.Fields.List, &ast.Field{
					Doc:   nil,
					Names: []*ast.Ident{ast.NewIdent(param.GetName())},
					Type:  ast.NewIdent(param.Type),
					Tag: &ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf("`json:\"%s\" form:\"%s\"`", param.Tag(), param.Tag()),
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
	if err := os.WriteFile(strc.Path, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func SyncValidate(strc *Struct) error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, strc.Path, nil, parser.ParseComments)
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
					if ok && ident.Name == strc.Name {
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
							X:    ast.NewIdent(strc.Name),
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
	for _, param := range strc.Params {
		ast.Inspect(validator.Body, func(node ast.Node) bool {
			if call, ok := node.(*ast.CallExpr); ok {
				if fun, ok := call.Fun.(*ast.SelectorExpr); ok && fun.Sel.Name == "ValidateStruct" {
					for _, arg := range call.Args {
						if vf, ok := arg.(*ast.CallExpr); ok {
							for _, vfArgs := range vf.Args {
								if field, ok := vfArgs.(*ast.UnaryExpr); ok {
									if fieldX, ok := field.X.(*ast.SelectorExpr); ok {
										if fieldX.Sel.Name == param.GetName() {
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
									Sel: ast.NewIdent(param.GetName()),
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
					if strings.Contains(strings.ToLower(param.GetName()), "email") {
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
	if err := os.WriteFile(strc.Path, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func SyncMock(filePath string, strc *Struct) error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, filePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	mockName := fmt.Sprintf("New%s", strc.Name)
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
									Sel: ast.NewIdent(strc.Name),
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
										Sel: ast.NewIdent(strc.Name),
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
	for _, param := range strc.Params {
		ast.Inspect(mock.Body, func(node ast.Node) bool {
			if cl, ok := node.(*ast.CompositeLit); ok {
				if sel, ok := cl.Type.(*ast.SelectorExpr); ok {
					if sel.Sel.Name != strc.Name {
						return true
					}
				}
				for _, elt := range cl.Elts {
					if kve, ok := elt.(*ast.KeyValueExpr); ok {
						if ident, ok := kve.Key.(*ast.Ident); ok && ident.Name == param.GetName() {
							return false
						}
					}
				}
				cl.Elts = append(cl.Elts, &ast.KeyValueExpr{
					Key:   ast.NewIdent(param.GetName()),
					Colon: 0,
					Value: FakeAst(param.Type),
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
	if err := os.WriteFile(filePath, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func FakeAst(t string) ast.Expr {
	var fake ast.Expr
	typeName := strings.TrimPrefix(t, "*")
	switch typeName {
	case "int", "int64", "int8", "int16", "int32", "float32", "float64", "uint", "uint8", "uint16", "uint32", "uint64":
		fake = &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("faker"),
						Sel: ast.NewIdent("New"),
					},
					Lparen:   0,
					Args:     nil,
					Ellipsis: 0,
					Rparen:   0,
				},
				Sel: ast.NewIdent(FakeNumberFunc(typeName)),
			},
			Lparen:   0,
			Args:     nil,
			Ellipsis: token.NoPos,
			Rparen:   0,
		}
	case "[]int", "[]int8", "[]int16", "[]int32", "[]int64", "[]float32", "[]float64", "[]uint", "[]uint8", "[]uint16", "[]uint32", "[]uint64":
		fake = &ast.CompositeLit{
			Type: &ast.ArrayType{
				Lbrack: 0,
				Len:    nil,
				Elt:    ast.NewIdent(strings.TrimPrefix(typeName, "[]")),
			},
			Lbrace: 0,
			Elts: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("faker"),
								Sel: ast.NewIdent("New"),
							},
							Lparen:   0,
							Args:     nil,
							Ellipsis: 0,
							Rparen:   0,
						},
						Sel: ast.NewIdent(FakeNumberFunc(strings.TrimPrefix(typeName, "[]"))),
					},
					Lparen:   0,
					Args:     nil,
					Ellipsis: token.NoPos,
					Rparen:   0,
				},
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("faker"),
								Sel: ast.NewIdent("New"),
							},
							Lparen:   0,
							Args:     nil,
							Ellipsis: 0,
							Rparen:   0,
						},
						Sel: ast.NewIdent(FakeNumberFunc(strings.TrimPrefix(typeName, "[]"))),
					},
					Lparen:   0,
					Args:     nil,
					Ellipsis: token.NoPos,
					Rparen:   0,
				},
			},
			Rbrace:     0,
			Incomplete: false,
		}
	case "string":
		fake = &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("faker"),
								Sel: ast.NewIdent("New"),
							},
							Lparen:   0,
							Args:     nil,
							Ellipsis: 0,
							Rparen:   0,
						},
						Sel: ast.NewIdent("Lorem"),
					},
					Lparen:   0,
					Args:     nil,
					Ellipsis: 0,
					Rparen:   0,
				},
				Sel: ast.NewIdent("Text"),
			},
			Lparen:   0,
			Args:     []ast.Expr{&ast.BasicLit{Kind: token.INT, Value: "256"}},
			Ellipsis: 0,
			Rparen:   0,
		}
	case "[]string":
		fake = &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("faker"),
								Sel: ast.NewIdent("New"),
							},
							Lparen:   0,
							Args:     nil,
							Ellipsis: 0,
							Rparen:   0,
						},
						Sel: ast.NewIdent("Lorem"),
					},
					Lparen:   0,
					Args:     nil,
					Ellipsis: 0,
					Rparen:   0,
				},
				Sel: ast.NewIdent("Words"),
			},
			Lparen: 0,
			Args: []ast.Expr{
				&ast.BasicLit{Kind: token.INT, Value: "27"},
			},
			Ellipsis: 0,
			Rparen:   0,
		}
	case "uuid", "UUID":
		fake = &ast.CallExpr{
			Fun:    ast.NewIdent("models.UUID"),
			Lparen: 0,
			Args: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("uuid"),
						Sel: ast.NewIdent("NewString"),
					},
					Lparen:   0,
					Args:     nil,
					Ellipsis: 0,
					Rparen:   0,
				},
			},
			Ellipsis: 0,
			Rparen:   0,
		}
	case "time.Time":
		fake = &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("faker"),
								Sel: ast.NewIdent("New"),
							},
							Lparen:   0,
							Args:     nil,
							Ellipsis: 0,
							Rparen:   0,
						},
						Sel: ast.NewIdent("Time"),
					},
					Lparen:   0,
					Args:     nil,
					Ellipsis: 0,
					Rparen:   0,
				},
				Sel: ast.NewIdent("Time"),
			},
			Lparen: 0,
			Args: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("time"),
						Sel: ast.NewIdent("Now"),
					},
					Lparen:   0,
					Args:     nil,
					Ellipsis: 0,
					Rparen:   0,
				},
			},
			Ellipsis: 0,
			Rparen:   0,
		}
	default:
		fake = &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("faker"),
				Sel: ast.NewIdent("Todo()"),
			},
			Lparen:   0,
			Args:     nil,
			Ellipsis: 0,
			Rparen:   0,
		}
	}
	if strings.HasPrefix(t, "*") {
		fake = &ast.CallExpr{
			Fun:      ast.NewIdent("utils.Pointer"),
			Lparen:   0,
			Args:     []ast.Expr{fake},
			Ellipsis: 0,
			Rparen:   0,
		}
	}
	return fake
}

func FakeNumberFunc(t string) string {
	switch t {
	case "int", "int64", "int8", "int16", "int32", "float32", "float64":
		return strcase.ToCamel(t)
	case "uint", "uint8", "uint16", "uint32", "uint64":
		return fmt.Sprintf("UInt%s", strings.TrimPrefix(t, "uint"))
	default:
		return "Todo"
	}
}

func SyncInterface(inter *Interface) error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, inter.Path, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var structureExists bool
	var structure *ast.TypeSpec
	_ = structureExists
	_ = structure
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok && t.Name.String() == inter.Name {
			structure = t
			structureExists = true
			return false
		}
		return true
	})
	if structure == nil {
		structure = &ast.TypeSpec{
			Doc:        nil,
			Name:       ast.NewIdent(inter.Name),
			TypeParams: nil,
			Assign:     0,
			Type: &ast.InterfaceType{
				Interface:  0,
				Methods:    &ast.FieldList{},
				Incomplete: false,
			},
			Comment: nil,
		}
	}
	for _, method := range inter.Methods {
		ast.Inspect(structure, func(node ast.Node) bool {
			if st, ok := node.(*ast.InterfaceType); ok && st.Methods != nil {
				for _, meth := range st.Methods.List {
					for _, fieldName := range meth.Names {
						if fieldName.Name == method.Name {
							return false
						}
					}
				}
				fn := &ast.FuncType{
					Func:       0,
					TypeParams: nil,
					Params: &ast.FieldList{
						Opening: 0,
						List:    []*ast.Field{},
					},
					Results: &ast.FieldList{
						Opening: 0,
						List:    nil,
						Closing: 0,
					},
				}
				for _, par := range method.Args {
					fn.Params.List = append(fn.Params.List, &ast.Field{
						Doc:     nil,
						Names:   []*ast.Ident{ast.NewIdent(par.Name)},
						Type:    ast.NewIdent(par.Type),
						Tag:     nil,
						Comment: nil,
					})
				}
				for _, res := range method.Results {
					fn.Results.List = append(fn.Results.List, &ast.Field{
						Doc:     nil,
						Names:   nil,
						Type:    ast.NewIdent(res.Type),
						Tag:     nil,
						Comment: nil,
					})
				}
				st.Methods.List = append(st.Methods.List, &ast.Field{
					Doc:     nil,
					Names:   []*ast.Ident{ast.NewIdent(method.Name)},
					Type:    fn,
					Tag:     nil,
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
	if err := os.WriteFile(inter.Path, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}
