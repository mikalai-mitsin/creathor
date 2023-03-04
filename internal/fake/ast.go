package fake

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"go/ast"
	"go/token"
	"strings"
)

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
				Sel: ast.NewIdent("Sentence"),
			},
			Lparen:   0,
			Args:     []ast.Expr{&ast.BasicLit{Kind: token.INT, Value: "15"}},
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
	case "[]uuid", "[]UUID":
		fake = &ast.CompositeLit{
			Type: &ast.ArrayType{
				Lbrack: 0,
				Len:    nil,
				Elt:    ast.NewIdent("models.UUID"),
			},
			Lbrace: 0,
			Elts: []ast.Expr{
				&ast.CallExpr{
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
				},
				&ast.CallExpr{
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
				},
			},
			Rbrace:     0,
			Incomplete: false,
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
				Sel: ast.NewIdent("Todo"),
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
