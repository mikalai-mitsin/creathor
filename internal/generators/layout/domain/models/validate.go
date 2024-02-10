package models

import (
	"bytes"
	mods "github.com/018bf/creathor/internal/domain"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"
	"strings"
)

type Validate struct {
	typeSpec *ast.TypeSpec
	fileName string
	domain   *mods.Domain
}

func NewValidate(typeSpec *ast.TypeSpec, fileName string, domain *mods.Domain) *Validate {
	return &Validate{typeSpec: typeSpec, fileName: fileName, domain: domain}
}
func (m *Validate) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", m.domain.Name, "models", m.fileName)
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var validate *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if fun, ok := node.(*ast.FuncDecl); ok && fun.Name.Name == "Validate" {
			ast.Inspect(fun.Recv, func(node ast.Node) bool {
				if ident, ok := node.(*ast.Ident); ok &&
					ident.String() == m.typeSpec.Name.String() {
					validate = fun
					return false
				}
				return true
			})
		}
		return true
	})
	if validate == nil {
		validate = m.method()
		file.Decls = append(file.Decls, validate)
	}
	m.fill(validate.Body)
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(filename, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (m *Validate) fill(body *ast.BlockStmt) {
	for _, param := range m.checkers() {
		ast.Inspect(body, func(node ast.Node) bool {
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
}

func (m *Validate) checker(name *ast.Ident, typeName ast.Expr) *ast.CallExpr {
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
					Sel: name,
				},
			},
		},
	}
	if _, ok := typeName.(*ast.StarExpr); !ok {
		call.Args = append(call.Args, &ast.SelectorExpr{
			X:   ast.NewIdent("validation"),
			Sel: ast.NewIdent("Required"),
		})
	}
	if ident, ok := typeName.(*ast.Ident); ok && ident.String() == "uuid.UUID" {
		call.Args = append(call.Args, &ast.SelectorExpr{
			X:   ast.NewIdent("is"),
			Sel: ast.NewIdent("UUID"),
		})
	}
	//if !strings.HasPrefix(param.Type, "*") {
	//	call.Args = append(call.Args, &ast.SelectorExpr{
	//		X:   ast.NewIdent("validation"),
	//		Sel: ast.NewIdent("Required"),
	//	})
	//}
	//if strings.ToLower(param.Type) == "uuid" {
	//	call.Args = append(call.Args, &ast.SelectorExpr{
	//		X:   ast.NewIdent("is"),
	//		Sel: ast.NewIdent("UUID"),
	//	})
	//}
	if strings.Contains(strings.ToLower(name.String()), "email") {
		call.Args = append(call.Args, &ast.SelectorExpr{
			X:   ast.NewIdent("is"),
			Sel: ast.NewIdent("EmailFormat"),
		})
	}
	return call
}

func (m *Validate) checkers() []*ast.CallExpr {
	var fields []*ast.CallExpr
	if st, ok := m.typeSpec.Type.(*ast.StructType); ok && st.Fields != nil {
		for _, field := range st.Fields.List {
			for _, name := range field.Names {
				fields = append(fields, m.checker(name, field.Type))
			}
		}
	}
	return fields
}

func (m *Validate) method() *ast.FuncDecl {
	exprs := []ast.Expr{
		ast.NewIdent("m"),
	}
	for _, expr := range m.checkers() {
		exprs = append(exprs, expr)
	}
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("m"),
					},
					Type: &ast.StarExpr{
						X: m.typeSpec.Name,
					},
				},
			},
		},
		Name: ast.NewIdent("Validate"),
		Type: &ast.FuncType{
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
				&ast.AssignStmt{
					Lhs: []ast.Expr{ast.NewIdent("err")},
					Tok: token.DEFINE,
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
					Cond: &ast.BinaryExpr{
						X:  ast.NewIdent("err"),
						Op: token.NEQ,
						Y:  ast.NewIdent("nil"),
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   ast.NewIdent("errs"),
											Sel: ast.NewIdent("FromValidationError"),
										},
										Args: []ast.Expr{
											ast.NewIdent("err"),
										},
									},
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{ast.NewIdent("nil")},
				},
			},
		},
	}
}
