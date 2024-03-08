package fake

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/iancoleman/strcase"
)

func baseValue(t ast.Expr) ast.Expr {
	var fake ast.Expr
	var typeName string
	switch u := t.(type) {
	case *ast.SelectorExpr:
		typeName = fmt.Sprintf("%s.%s", u.X, u.Sel)
	case *ast.Ident:
		typeName = u.String()
	default:
		typeName = fmt.Sprint(u)
	}
	switch typeName {
	case "int",
		"int64",
		"int8",
		"int16",
		"int32",
		"float32",
		"float64",
		"uint",
		"uint8",
		"uint16",
		"uint32",
		"uint64":
		fake = &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("faker"),
						Sel: ast.NewIdent("New"),
					},
				},
				Sel: ast.NewIdent(numberFunc(typeName)),
			},
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
						},
						Sel: ast.NewIdent("Lorem"),
					},
				},
				Sel: ast.NewIdent("Sentence"),
			},
			Args: []ast.Expr{&ast.BasicLit{Kind: token.INT, Value: "15"}},
		}
	case "uuid", "UUID", "uuid.UUID":
		fake = &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("uuid"),
				Sel: ast.NewIdent("NewUUID"),
			},
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
						},
						Sel: ast.NewIdent("Time"),
					},
				},
				Sel: ast.NewIdent("Time"),
			},
			Args: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("time"),
						Sel: ast.NewIdent("Now"),
					},
				},
			},
		}
	case "models.GroupID", "GroupID":
		fake = &ast.SelectorExpr{
			X:   ast.NewIdent("models"),
			Sel: ast.NewIdent("GroupIDUser"),
		}
	default:
		fake = &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("faker"),
				Sel: ast.NewIdent("Todo"),
			},
		}
	}
	return fake
}

func Value(t ast.Expr) ast.Expr {
	var fake ast.Expr
	switch value := t.(type) {
	case *ast.ArrayType:
		if fmt.Sprint(value.Elt) == "UUID" {
			value.Elt = ast.NewIdent("uuid.UUID")
		}
		fake = &ast.CompositeLit{
			Type: &ast.ArrayType{
				Elt: value.Elt,
			},
			Elts: []ast.Expr{
				baseValue(value.Elt),
				baseValue(value.Elt),
			},
		}
	case *ast.StarExpr:
		switch x := value.X.(type) {
		case *ast.ArrayType:
			fake = &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("pointer"),
					Sel: ast.NewIdent("Pointer"),
				},
				Args: []ast.Expr{
					&ast.CompositeLit{
						Type: &ast.ArrayType{
							Elt: x.Elt,
						},
						Elts: []ast.Expr{
							baseValue(x.Elt),
							baseValue(x.Elt),
						},
					},
				},
			}
		default:
			fake = &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("pointer"),
					Sel: ast.NewIdent("Pointer"),
				},
				Args: []ast.Expr{baseValue(x)},
			}
		}
	case *ast.Ident:
		fake = baseValue(value)
	case *ast.SelectorExpr:
		fake = baseValue(value)
	default:
		fake = baseValue(ast.NewIdent("TODO"))
	}
	return fake
}

func baseEmail() ast.Expr {
	return &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("faker"),
							Sel: ast.NewIdent("New"),
						},
					},
					Sel: ast.NewIdent("Internet"),
				},
			},
			Sel: ast.NewIdent("Email"),
		},
	}
}

func Email(t ast.Expr) ast.Expr {
	var fake ast.Expr
	switch t.(type) {
	case *ast.StarExpr:
		fake = &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("pointer"),
				Sel: ast.NewIdent("Pointer"),
			},
			Args: []ast.Expr{baseEmail()},
		}
	case *ast.Ident:
		fake = baseEmail()
	default:
		fake = baseValue(ast.NewIdent("TODO"))
	}
	return fake
}

func numberFunc(t string) string {
	switch t {
	case "int", "int64", "int8", "int16", "int32", "float32", "float64":
		return strcase.ToCamel(t)
	case "uint", "uint8", "uint16", "uint32", "uint64":
		return fmt.Sprintf("UInt%s", strings.TrimPrefix(t, "uint"))
	default:
		return "Todo"
	}
}
