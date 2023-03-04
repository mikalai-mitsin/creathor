package grpc

import (
	"bytes"
	"fmt"
	"github.com/018bf/creathor/internal/configs"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"
)

type Handler struct {
	model *configs.ModelConfig
}

func NewHandler(model *configs.ModelConfig) *Handler {
	return &Handler{
		model: model,
	}
}

func (h Handler) file() *ast.File {
	return &ast.File{
		Name: &ast.Ident{
			Name: "grpc",
		},
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"context\"",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"fmt\"",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"github.com/018bf/example/internal/domain/interceptors\"",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"github.com/018bf/example/internal/domain/models\"",
						},
					},
					&ast.ImportSpec{
						Name: &ast.Ident{
							Name: h.model.ProtoPackage,
						},
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"github.com/018bf/example/pkg/examplepb/v1\"",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"github.com/018bf/example/pkg/log\"",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"github.com/018bf/example/pkg/utils\"",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"google.golang.org/protobuf/types/known/emptypb\"",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"google.golang.org/grpc\"",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"google.golang.org/grpc/metadata\"",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"google.golang.org/protobuf/types/known/structpb\"",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"google.golang.org/protobuf/types/known/timestamppb\"",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"google.golang.org/protobuf/types/known/wrapperspb\"",
						},
					},
				},
			},
		},
	}
}

func (h Handler) filename() string {
	return path.Join("internal", "interfaces", "grpc", h.model.FileName())
}

func (h Handler) createParams() []ast.Expr {
	var exprs []ast.Expr
	for _, param := range h.model.Params {
		var value ast.Expr
		if param.IsSlice() {
			switch param.Type {
			case param.GRPCType():
				value = &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: "input",
						},
						Sel: &ast.Ident{
							Name: param.GRPCGetter(),
						},
					},
				}
			default:
				value = ast.NewIdent("nil")
			}
		} else {
			value = &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "input",
					},
					Sel: &ast.Ident{
						Name: param.GRPCGetter(),
					},
				},
			}
			switch param.Type {
			case "time.Time":
				value = &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: value,
						Sel: &ast.Ident{
							Name: "AsTime",
						},
					},
				}
			case param.GRPCType():
			default:
				value = &ast.CallExpr{
					Fun: ast.NewIdent(param.Type),
					Args: []ast.Expr{
						value,
					},
				}
			}
		}
		exprs = append(exprs, &ast.KeyValueExpr{
			Key:   ast.NewIdent(param.GetName()),
			Value: value,
		})
	}
	return exprs
}

func (h Handler) astEncodeCreate() *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: &ast.Ident{
			Name: fmt.Sprintf("encode%s", h.model.CreateTypeName()),
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
								Name: "input",
							},
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: h.model.ProtoPackage,
								},
								Sel: &ast.Ident{
									Name: h.model.CreateTypeName(),
								},
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "models",
								},
								Sel: &ast.Ident{
									Name: h.model.CreateTypeName(),
								},
							},
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "create",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X: &ast.CompositeLit{
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "models",
									},
									Sel: &ast.Ident{
										Name: h.model.CreateTypeName(),
									},
								},
								Elts: h.createParams(),
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.Ident{
							Name: "create",
						},
					},
				},
			},
		},
	}
}

func (h Handler) syncEncodeCreate() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, h.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == fmt.Sprintf("encode%s", h.model.CreateTypeName()) {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = h.astEncodeCreate()
	}
	for _, expr := range h.createParams() {
		kv, ok := expr.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		ast.Inspect(method, func(node ast.Node) bool {
			if cl, ok := node.(*ast.CompositeLit); ok {
				for _, elt := range cl.Elts {
					if item, ok := elt.(*ast.KeyValueExpr); ok {
						if item.Key.(*ast.Ident).String() == kv.Key.(*ast.Ident).String() {
							return false
						}
					}
				}
				cl.Elts = append(cl.Elts, kv)
			}
			return true
		})
	}
	rangeStmt := &ast.RangeStmt{
		Key: &ast.Ident{
			Name: "_",
		},
		Value: &ast.Ident{
			Name: "param",
		},
		Tok: token.DEFINE,
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X: &ast.Ident{
					Name: "input",
				},
				Sel: &ast.Ident{
					Name: "GetTags",
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X: &ast.Ident{
								Name: "create",
							},
							Sel: &ast.Ident{
								Name: "Tags",
							},
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.Ident{
								Name: "append",
							},
							Args: []ast.Expr{
								&ast.SelectorExpr{
									X: &ast.Ident{
										Name: "create",
									},
									Sel: &ast.Ident{
										Name: "Tags",
									},
								},
								&ast.CallExpr{
									Fun: &ast.Ident{
										Name: "string",
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "param",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	_ = rangeStmt
	// TODO: add ranges
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(h.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (h Handler) updateStmts() []*ast.IfStmt {
	var stmts []*ast.IfStmt
	for _, param := range h.model.Params {
		var body []ast.Stmt
		if param.Type == "time.Time" {
			body = []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X: &ast.Ident{
								Name: "update",
							},
							Sel: &ast.Ident{
								Name: param.GetName(),
							},
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "utils",
								},
								Sel: &ast.Ident{
									Name: "Pointer",
								},
							},
							Args: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("input"),
												Sel: ast.NewIdent(param.GRPCGetter()),
											},
										},
										Sel: ast.NewIdent("AsTime"),
									},
								},
							},
						},
					},
				},
			}
		} else if param.IsSlice() {
			value := &ast.CallExpr{
				Fun: ast.NewIdent(param.SliceType()),
				Args: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "item",
							},
							Sel: &ast.Ident{
								Name: param.GrpcGetFromListValueAs(),
							},
						},
					},
				},
			}
			body = []ast.Stmt{
				&ast.DeclStmt{
					Decl: &ast.GenDecl{
						Tok: token.VAR,
						Specs: []ast.Spec{
							&ast.ValueSpec{
								Names: []*ast.Ident{
									{
										Name: "params",
									},
								},
								Type: &ast.ArrayType{
									Elt: &ast.Ident{
										Name: param.SliceType(),
									},
								},
							},
						},
					},
				},
				&ast.RangeStmt{
					Key: &ast.Ident{
						Name: "_",
					},
					Value: &ast.Ident{
						Name: "item",
					},
					Tok: token.DEFINE,
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "input",
									},
									Sel: &ast.Ident{
										Name: param.GRPCGetter(),
									},
								},
							},
							Sel: &ast.Ident{
								Name: "GetValues",
							},
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.Ident{
										Name: "params",
									},
								},
								Tok: token.ASSIGN,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.Ident{
											Name: "append",
										},
										Args: []ast.Expr{
											&ast.Ident{
												Name: "params",
											},
											value,
										},
									},
								},
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X: &ast.Ident{
								Name: "update",
							},
							Sel: &ast.Ident{
								Name: param.GetName(),
							},
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "utils",
								},
								Sel: &ast.Ident{
									Name: "Pointer",
								},
							},
							Args: []ast.Expr{
								&ast.Ident{
									Name: "params",
								},
							},
						},
					},
				},
			}
		} else {
			var value ast.Expr
			value = &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "input",
							},
							Sel: &ast.Ident{
								Name: param.GRPCGetter(),
							},
						},
					},
					Sel: &ast.Ident{
						Name: "GetValue",
					},
				},
			}
			if param.Type != param.GRPCType() {
				value = &ast.CallExpr{
					Fun:  ast.NewIdent(param.Type),
					Args: []ast.Expr{value},
				}
			}
			body = []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X: &ast.Ident{
								Name: "update",
							},
							Sel: &ast.Ident{
								Name: param.GetName(),
							},
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "utils",
								},
								Sel: &ast.Ident{
									Name: "Pointer",
								},
							},
							Args: []ast.Expr{
								value,
							},
						},
					},
				},
			}
		}
		stmts = append(stmts, &ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: "input",
						},
						Sel: &ast.Ident{
							Name: param.GRPCGetter(),
						},
					},
				},
				Op: token.NEQ,
				Y: &ast.Ident{
					Name: "nil",
				},
			},
			Body: &ast.BlockStmt{
				List: body,
			},
		})
	}
	return stmts
}

func (h Handler) astEncodeUpdate() *ast.FuncDecl {
	body := []ast.Stmt{
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: "update",
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.UnaryExpr{
					Op: token.AND,
					X: &ast.CompositeLit{
						Type: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "models",
							},
							Sel: &ast.Ident{
								Name: h.model.UpdateTypeName(),
							},
						},
						Elts: []ast.Expr{
							&ast.KeyValueExpr{
								Key: &ast.Ident{
									Name: "ID",
								},
								Value: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "models",
										},
										Sel: &ast.Ident{
											Name: "UUID",
										},
									},
									Args: []ast.Expr{
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "input",
												},
												Sel: &ast.Ident{
													Name: "GetId",
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
	}
	for _, stmt := range h.updateStmts() {
		body = append(body, stmt)
	}
	body = append(body, &ast.ReturnStmt{
		Results: []ast.Expr{
			&ast.Ident{
				Name: "update",
			},
		},
	})
	return &ast.FuncDecl{
		Name: &ast.Ident{
			Name: fmt.Sprintf("encode%s", h.model.UpdateTypeName()),
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
								Name: "input",
							},
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: h.model.ProtoPackage,
								},
								Sel: &ast.Ident{
									Name: h.model.UpdateTypeName(),
								},
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "models",
								},
								Sel: &ast.Ident{
									Name: h.model.UpdateTypeName(),
								},
							},
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: body,
		},
	}
}

func (h Handler) syncEncodeUpdate() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, h.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == fmt.Sprintf("encode%s", h.model.UpdateTypeName()) {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = h.astEncodeUpdate()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(h.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (h Handler) astEncodeFilter() *ast.FuncDecl {
	stmts := []ast.Stmt{
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: "filter",
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.UnaryExpr{
					Op: token.AND,
					X: &ast.CompositeLit{
						Type: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "models",
							},
							Sel: &ast.Ident{
								Name: h.model.FilterTypeName(),
							},
						},
						Elts: []ast.Expr{
							&ast.KeyValueExpr{
								Key: &ast.Ident{
									Name: "IDs",
								},
								Value: &ast.Ident{
									Name: "nil",
								},
							},
							&ast.KeyValueExpr{
								Key: &ast.Ident{
									Name: "PageSize",
								},
								Value: &ast.Ident{
									Name: "nil",
								},
							},
							&ast.KeyValueExpr{
								Key: &ast.Ident{
									Name: "PageNumber",
								},
								Value: &ast.Ident{
									Name: "nil",
								},
							},
							&ast.KeyValueExpr{
								Key: &ast.Ident{
									Name: "OrderBy",
								},
								Value: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "input",
										},
										Sel: &ast.Ident{
											Name: "GetOrderBy",
										},
									},
								},
							},
							&ast.KeyValueExpr{
								Key: &ast.Ident{
									Name: "Search",
								},
								Value: &ast.Ident{
									Name: "nil",
								},
							},
						},
					},
				},
			},
		},
		&ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: "input",
						},
						Sel: &ast.Ident{
							Name: "GetPageSize",
						},
					},
				},
				Op: token.NEQ,
				Y: &ast.Ident{
					Name: "nil",
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.SelectorExpr{
								X: &ast.Ident{
									Name: "filter",
								},
								Sel: &ast.Ident{
									Name: "PageSize",
								},
							},
						},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "utils",
									},
									Sel: &ast.Ident{
										Name: "Pointer",
									},
								},
								Args: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "input",
													},
													Sel: &ast.Ident{
														Name: "GetPageSize",
													},
												},
											},
											Sel: &ast.Ident{
												Name: "GetValue",
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
		&ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: "input",
						},
						Sel: &ast.Ident{
							Name: "GetPageNumber",
						},
					},
				},
				Op: token.NEQ,
				Y: &ast.Ident{
					Name: "nil",
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.SelectorExpr{
								X: &ast.Ident{
									Name: "filter",
								},
								Sel: &ast.Ident{
									Name: "PageNumber",
								},
							},
						},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "utils",
									},
									Sel: &ast.Ident{
										Name: "Pointer",
									},
								},
								Args: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "input",
													},
													Sel: &ast.Ident{
														Name: "GetPageNumber",
													},
												},
											},
											Sel: &ast.Ident{
												Name: "GetValue",
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
		&ast.RangeStmt{
			Key: &ast.Ident{
				Name: "_",
			},
			Value: &ast.Ident{
				Name: "id",
			},
			Tok: token.DEFINE,
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "input",
					},
					Sel: &ast.Ident{
						Name: "GetIds",
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.SelectorExpr{
								X: &ast.Ident{
									Name: "filter",
								},
								Sel: &ast.Ident{
									Name: "IDs",
								},
							},
						},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.Ident{
									Name: "append",
								},
								Args: []ast.Expr{
									&ast.SelectorExpr{
										X: &ast.Ident{
											Name: "filter",
										},
										Sel: &ast.Ident{
											Name: "IDs",
										},
									},
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "models",
											},
											Sel: &ast.Ident{
												Name: "UUID",
											},
										},
										Args: []ast.Expr{
											&ast.Ident{
												Name: "id",
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
	}
	if h.model.SearchEnabled() {
		stmts = append(stmts, &ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: "input",
						},
						Sel: &ast.Ident{
							Name: "GetSearch",
						},
					},
				},
				Op: token.NEQ,
				Y: &ast.Ident{
					Name: "nil",
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.SelectorExpr{
								X: &ast.Ident{
									Name: "filter",
								},
								Sel: &ast.Ident{
									Name: "Search",
								},
							},
						},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "utils",
									},
									Sel: &ast.Ident{
										Name: "Pointer",
									},
								},
								Args: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "input",
													},
													Sel: &ast.Ident{
														Name: "GetSearch",
													},
												},
											},
											Sel: &ast.Ident{
												Name: "GetValue",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		})
	}
	stmts = append(stmts, &ast.ReturnStmt{
		Results: []ast.Expr{
			&ast.Ident{
				Name: "filter",
			},
		},
	})
	return &ast.FuncDecl{
		Name: &ast.Ident{
			Name: fmt.Sprintf("encode%s", h.model.FilterTypeName()),
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
								Name: "input",
							},
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: h.model.ProtoPackage,
								},
								Sel: &ast.Ident{
									Name: h.model.FilterTypeName(),
								},
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "models",
								},
								Sel: &ast.Ident{
									Name: h.model.FilterTypeName(),
								},
							},
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: stmts,
		},
	}
}

func (h Handler) syncEncodeFilter() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, h.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == fmt.Sprintf("encode%s", h.model.FilterTypeName()) {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = h.astEncodeFilter()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(h.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (h Handler) astDecodeModel() *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: &ast.Ident{
			Name: fmt.Sprintf("decode%s", h.model.ModelName()),
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
								Name: h.model.Variable(),
							},
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "models",
								},
								Sel: &ast.Ident{
									Name: h.model.ModelName(),
								},
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: h.model.ProtoPackage,
								},
								Sel: &ast.Ident{
									Name: h.model.ModelName(),
								},
							},
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "response",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X: &ast.CompositeLit{
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: h.model.ProtoPackage,
									},
									Sel: &ast.Ident{
										Name: h.model.ModelName(),
									},
								},
								Elts: h.modelParams(),
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.Ident{
							Name: "response",
						},
					},
				},
			},
		},
	}
}

func (h Handler) modelParams() []ast.Expr {
	exprs := []ast.Expr{
		&ast.KeyValueExpr{
			Key: &ast.Ident{
				Name: "Id",
			},
			Value: &ast.CallExpr{
				Fun: &ast.Ident{
					Name: "string",
				},
				Args: []ast.Expr{
					&ast.SelectorExpr{
						X: &ast.Ident{
							Name: h.model.Variable(),
						},
						Sel: &ast.Ident{
							Name: "ID",
						},
					},
				},
			},
		},
		&ast.KeyValueExpr{
			Key: &ast.Ident{
				Name: "UpdatedAt",
			},
			Value: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "timestamppb",
					},
					Sel: &ast.Ident{
						Name: "New",
					},
				},
				Args: []ast.Expr{
					&ast.SelectorExpr{
						X: &ast.Ident{
							Name: h.model.Variable(),
						},
						Sel: &ast.Ident{
							Name: "UpdatedAt",
						},
					},
				},
			},
		},
		&ast.KeyValueExpr{
			Key: &ast.Ident{
				Name: "CreatedAt",
			},
			Value: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "timestamppb",
					},
					Sel: &ast.Ident{
						Name: "New",
					},
				},
				Args: []ast.Expr{
					&ast.SelectorExpr{
						X: &ast.Ident{
							Name: h.model.Variable(),
						},
						Sel: &ast.Ident{
							Name: "CreatedAt",
						},
					},
				},
			},
		},
	}
	for _, param := range h.model.Params {
		var value ast.Expr
		value = &ast.SelectorExpr{
			X: &ast.Ident{
				Name: h.model.Variable(),
			},
			Sel: &ast.Ident{
				Name: param.GetName(),
			},
		}
		if param.Type != param.GRPCType() {
			value = &ast.CallExpr{
				Fun: &ast.Ident{
					Name: param.GRPCType(),
				},
				Args: []ast.Expr{
					value,
				},
			}
		}
		if param.IsSlice() && param.GRPCType() != param.Type {
			value = &ast.CallExpr{
				Fun: &ast.IndexListExpr{
					X: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: "utils",
						},
						Sel: &ast.Ident{
							Name: "ChangeType",
						},
					},
					Indices: []ast.Expr{
						&ast.Ident{
							Name: param.GRPCSliceType(),
						},
						&ast.Ident{
							Name: param.SliceType(),
						},
					},
				},
				Args: []ast.Expr{
					&ast.SelectorExpr{
						X: &ast.Ident{
							Name: h.model.Variable(),
						},
						Sel: &ast.Ident{
							Name: param.GetName(),
						},
					},
				},
			}
		}
		exprs = append(exprs, &ast.KeyValueExpr{
			Key: &ast.Ident{
				Name: param.GRPCParam(),
			},
			Value: value,
		})
	}
	return exprs
}

func (h Handler) syncDecodeModel() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, h.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == fmt.Sprintf("decode%s", h.model.ModelName()) {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = h.astDecodeModel()
	}
	for _, param := range h.modelParams() {
		pr := param.(*ast.KeyValueExpr)
		prKey := pr.Key.(*ast.Ident)
		ast.Inspect(method, func(node ast.Node) bool {
			if cl, ok := node.(*ast.CompositeLit); ok {
				if clType, ok := cl.Type.(*ast.SelectorExpr); ok {
					if clType.Sel.String() != h.model.ModelName() {
						return true
					}
				}
				for _, elt := range cl.Elts {
					if kv, ok := elt.(*ast.KeyValueExpr); ok {
						if ident, ok := kv.Key.(*ast.Ident); ok {
							if ident.String() == prKey.String() {
								return false
							}
						}
					}
				}
				cl.Elts = append(cl.Elts, pr)
			}
			return true
		})
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(h.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (h Handler) astDecodeList() *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: &ast.Ident{
			Name: fmt.Sprintf("decodeList%s", h.model.ModelName()),
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
								Name: h.model.ListVariable(),
							},
						},
						Type: &ast.ArrayType{
							Elt: &ast.StarExpr{
								X: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "models",
									},
									Sel: &ast.Ident{
										Name: h.model.ModelName(),
									},
								},
							},
						},
					},
					{
						Names: []*ast.Ident{
							{
								Name: "count",
							},
						},
						Type: &ast.Ident{
							Name: "uint64",
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: h.model.ProtoPackage,
								},
								Sel: &ast.Ident{
									Name: fmt.Sprintf("List%s", h.model.ModelName()),
								},
							},
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "response",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X: &ast.CompositeLit{
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: h.model.ProtoPackage,
									},
									Sel: &ast.Ident{
										Name: fmt.Sprintf("List%s", h.model.ModelName()),
									},
								},
								Elts: []ast.Expr{
									&ast.KeyValueExpr{
										Key: &ast.Ident{
											Name: "Items",
										},
										Value: &ast.CallExpr{
											Fun: &ast.Ident{
												Name: "make",
											},
											Args: []ast.Expr{
												&ast.ArrayType{
													Elt: &ast.StarExpr{
														X: &ast.SelectorExpr{
															X: &ast.Ident{
																Name: h.model.ProtoPackage,
															},
															Sel: &ast.Ident{
																Name: h.model.ModelName(),
															},
														},
													},
												},
												&ast.BasicLit{
													Kind:  token.INT,
													Value: "0",
												},
												&ast.CallExpr{
													Fun: &ast.Ident{
														Name: "len",
													},
													Args: []ast.Expr{
														&ast.Ident{
															Name: h.model.ListVariable(),
														},
													},
												},
											},
										},
									},
									&ast.KeyValueExpr{
										Key: &ast.Ident{
											Name: "Count",
										},
										Value: &ast.Ident{
											Name: "count",
										},
									},
								},
							},
						},
					},
				},
				&ast.RangeStmt{
					Key: &ast.Ident{
						Name: "_",
					},
					Value: &ast.Ident{
						Name: h.model.Variable(),
					},
					Tok: token.DEFINE,
					X: &ast.Ident{
						Name: h.model.ListVariable(),
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.SelectorExpr{
										X: &ast.Ident{
											Name: "response",
										},
										Sel: &ast.Ident{
											Name: "Items",
										},
									},
								},
								Tok: token.ASSIGN,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.Ident{
											Name: "append",
										},
										Args: []ast.Expr{
											&ast.SelectorExpr{
												X: &ast.Ident{
													Name: "response",
												},
												Sel: &ast.Ident{
													Name: "Items",
												},
											},
											&ast.CallExpr{
												Fun: &ast.Ident{
													Name: fmt.Sprintf("decode%s", h.model.ModelName()),
												},
												Args: []ast.Expr{
													&ast.Ident{
														Name: h.model.Variable(),
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
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.Ident{
							Name: "response",
						},
					},
				},
			},
		},
	}
}

func (h Handler) syncDecodeList() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, h.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == fmt.Sprintf("decodeList%s", h.model.ModelName()) {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = h.astDecodeList()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(h.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (h Handler) astDecodeUpdate() *ast.FuncDecl {
	stmts := []ast.Stmt{
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: "result",
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.UnaryExpr{
					Op: token.AND,
					X: &ast.CompositeLit{
						Type: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: h.model.ProtoPackage,
							},
							Sel: &ast.Ident{
								Name: h.model.UpdateTypeName(),
							},
						},
						Elts: h.updateParams(),
					},
				},
			},
		},
	}
	for _, param := range h.model.Params {
		if !param.IsSlice() {
			continue
		}
		stmts = append(stmts, &ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "update",
					},
					Sel: &ast.Ident{
						Name: param.GetName(),
					},
				},
				Op: token.NEQ,
				Y: &ast.Ident{
					Name: "nil",
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.Ident{
								Name: "params",
							},
							&ast.Ident{
								Name: "err",
							},
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "structpb",
									},
									Sel: &ast.Ident{
										Name: "NewList",
									},
								},
								Args: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "utils",
											},
											Sel: &ast.Ident{
												Name: "ToAnySlice",
											},
										},
										Args: []ast.Expr{
											&ast.StarExpr{
												X: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "update",
													},
													Sel: &ast.Ident{
														Name: param.GetName(),
													},
												},
											},
										},
									},
								},
							},
						},
					},
					&ast.IfStmt{
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
										&ast.Ident{
											Name: "nil",
										},
									},
								},
							},
						},
					},
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.SelectorExpr{
								X: &ast.Ident{
									Name: "result",
								},
								Sel: &ast.Ident{
									Name: param.GRPCParam(),
								},
							},
						},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							&ast.Ident{
								Name: "params",
							},
						},
					},
				},
			},
		})
	}
	stmts = append(stmts, &ast.ReturnStmt{
		Results: []ast.Expr{
			&ast.Ident{
				Name: "result",
			},
		},
	})
	return &ast.FuncDecl{
		Name: &ast.Ident{
			Name: fmt.Sprintf("decode%s", h.model.UpdateTypeName()),
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
								Name: "update",
							},
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "models",
								},
								Sel: &ast.Ident{
									Name: h.model.UpdateTypeName(),
								},
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: h.model.ProtoPackage,
								},
								Sel: &ast.Ident{
									Name: h.model.UpdateTypeName(),
								},
							},
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: stmts,
		},
	}
}

func (h Handler) updateParams() []ast.Expr {
	exprs := []ast.Expr{
		&ast.KeyValueExpr{
			Key: &ast.Ident{
				Name: "Id",
			},
			Value: &ast.CallExpr{
				Fun: &ast.Ident{
					Name: "string",
				},
				Args: []ast.Expr{
					&ast.SelectorExpr{
						X: &ast.Ident{
							Name: "update",
						},
						Sel: &ast.Ident{
							Name: "ID",
						},
					},
				},
			},
		},
	}
	for _, param := range h.model.Params {
		var value ast.Expr
		if param.IsSlice() {
			value = ast.NewIdent("nil")
		} else {
			var v ast.Expr
			v = &ast.StarExpr{
				X: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "update",
					},
					Sel: &ast.Ident{
						Name: param.GetName(),
					},
				},
			}
			if param.GetGRPCWrapperArgumentType() != param.Type {
				v = &ast.CallExpr{
					Fun: &ast.Ident{
						Name: param.GetGRPCWrapperArgumentType(),
					},
					Args: []ast.Expr{v},
				}
			}
			value = &ast.CallExpr{
				Fun:  ast.NewIdent(param.GetGRPCWrapper()),
				Args: []ast.Expr{v},
			}
		}

		exprs = append(exprs, &ast.KeyValueExpr{
			Key: &ast.Ident{
				Name: param.GRPCParam(),
			},
			Value: value,
		})
	}
	return exprs
}

func (h Handler) syncDecodeUpdate() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, h.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == fmt.Sprintf("decode%s", h.model.UpdateTypeName()) {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = h.astDecodeUpdate()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(h.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (h Handler) astStruct() *ast.TypeSpec {
	return &ast.TypeSpec{
		Name: &ast.Ident{
			Name: h.model.GRPCHandlerTypeName(),
		},
		Type: &ast.StructType{
			Fields: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: h.model.ProtoPackage,
							},
							Sel: &ast.Ident{
								Name: fmt.Sprintf("Unimplemented%sServiceServer", h.model.ModelName()),
							},
						},
					},
					{
						Names: []*ast.Ident{
							{
								Name: h.model.InterceptorVariableName(),
							},
						},
						Type: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "interceptors",
							},
							Sel: &ast.Ident{
								Name: h.model.InterceptorTypeName(),
							},
						},
					},
					{
						Names: []*ast.Ident{
							{
								Name: "logger",
							},
						},
						Type: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "log",
							},
							Sel: &ast.Ident{
								Name: "Logger",
							},
						},
					},
				},
			},
		},
	}
}

func (h Handler) syncStruct() error {
	fileset := token.NewFileSet()
	filename := h.filename()
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = h.file()
	}
	var structureExists bool
	var structure *ast.TypeSpec
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok && t.Name.String() == h.model.GRPCHandlerTypeName() {
			structure = t
			structureExists = true
			return false
		}
		return true
	})
	if structure == nil {
		structure = h.astStruct()
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

func (h Handler) astConstructor() *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: &ast.Ident{
			Name: fmt.Sprintf("New%s", h.model.GRPCHandlerTypeName()),
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
								Name: h.model.InterceptorVariableName(),
							},
						},
						Type: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "interceptors",
							},
							Sel: &ast.Ident{
								Name: h.model.InterceptorTypeName(),
							},
						},
					},
					{
						Names: []*ast.Ident{
							{
								Name: "logger",
							},
						},
						Type: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "log",
							},
							Sel: &ast.Ident{
								Name: "Logger",
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: h.model.ProtoPackage,
							},
							Sel: &ast.Ident{
								Name: fmt.Sprintf("%sServiceServer", h.model.ModelName()),
							},
						},
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
								Type: &ast.Ident{
									Name: fmt.Sprintf("%sServiceServer", h.model.ModelName()),
								},
								Elts: []ast.Expr{
									&ast.KeyValueExpr{
										Key: &ast.Ident{
											Name: h.model.InterceptorVariableName(),
										},
										Value: &ast.Ident{
											Name: h.model.InterceptorVariableName(),
										},
									},
									&ast.KeyValueExpr{
										Key: &ast.Ident{
											Name: "logger",
										},
										Value: &ast.Ident{
											Name: "logger",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (h Handler) syncConstructor() error {
	fileset := token.NewFileSet()
	filename := h.filename()
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == fmt.Sprintf("New%s", h.model.GRPCHandlerTypeName()) {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = h.astConstructor()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
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

func (h Handler) astCreateMethod() *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{
							Name: "s",
						},
					},
					Type: &ast.StarExpr{
						X: &ast.Ident{
							Name: h.model.GRPCHandlerTypeName(),
						},
					},
				},
			},
		},
		Name: &ast.Ident{
			Name: "Create",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
								Name: "ctx",
							},
						},
						Type: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "context",
							},
							Sel: &ast.Ident{
								Name: "Context",
							},
						},
					},
					{
						Names: []*ast.Ident{
							{
								Name: "input",
							},
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: h.model.ProtoPackage,
								},
								Sel: &ast.Ident{
									Name: h.model.CreateTypeName(),
								},
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: h.model.ProtoPackage,
								},
								Sel: &ast.Ident{
									Name: h.model.ModelName(),
								},
							},
						},
					},
					{
						Type: &ast.Ident{
							Name: "error",
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: h.model.Variable(),
						},
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
										Name: "s",
									},
									Sel: &ast.Ident{
										Name: h.model.InterceptorVariableName(),
									},
								},
								Sel: &ast.Ident{
									Name: "Create",
								},
							},
							Args: []ast.Expr{
								&ast.Ident{
									Name: "ctx",
								},
								&ast.CallExpr{
									Fun: &ast.Ident{
										Name: fmt.Sprintf("encode%s", h.model.CreateTypeName()),
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "input",
										},
									},
								},
								&ast.TypeAssertExpr{
									X: &ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "ctx",
											},
											Sel: &ast.Ident{
												Name: "Value",
											},
										},
										Args: []ast.Expr{
											&ast.Ident{
												Name: "UserKey",
											},
										},
									},
									Type: &ast.StarExpr{
										X: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "models",
											},
											Sel: &ast.Ident{
												Name: "User",
											},
										},
									},
								},
							},
						},
					},
				},
				&ast.IfStmt{
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
									&ast.Ident{
										Name: "nil",
									},
									&ast.CallExpr{
										Fun: &ast.Ident{
											Name: "decodeError",
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
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.Ident{
								Name: fmt.Sprintf("decode%s", h.model.ModelName()),
							},
							Args: []ast.Expr{
								&ast.Ident{
									Name: h.model.Variable(),
								},
							},
						},
						&ast.Ident{
							Name: "nil",
						},
					},
				},
			},
		},
	}
}

func (h Handler) syncCreateMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, h.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == "Create" {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = h.astCreateMethod()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(h.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (h Handler) astGetMethod() *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{
							Name: "s",
						},
					},
					Type: &ast.StarExpr{
						X: &ast.Ident{
							Name: h.model.GRPCHandlerTypeName(),
						},
					},
				},
			},
		},
		Name: &ast.Ident{
			Name: "Get",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
								Name: "ctx",
							},
						},
						Type: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "context",
							},
							Sel: &ast.Ident{
								Name: "Context",
							},
						},
					},
					{
						Names: []*ast.Ident{
							{
								Name: "input",
							},
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: h.model.ProtoPackage,
								},
								Sel: &ast.Ident{
									Name: fmt.Sprintf("%sGet", h.model.ModelName()),
								},
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: h.model.ProtoPackage,
								},
								Sel: &ast.Ident{
									Name: h.model.ModelName(),
								},
							},
						},
					},
					{
						Type: &ast.Ident{
							Name: "error",
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: h.model.Variable(),
						},
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
										Name: "s",
									},
									Sel: &ast.Ident{
										Name: h.model.InterceptorVariableName(),
									},
								},
								Sel: &ast.Ident{
									Name: "Get",
								},
							},
							Args: []ast.Expr{
								&ast.Ident{
									Name: "ctx",
								},
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "models",
										},
										Sel: &ast.Ident{
											Name: "UUID",
										},
									},
									Args: []ast.Expr{
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "input",
												},
												Sel: &ast.Ident{
													Name: "GetId",
												},
											},
										},
									},
								},
								&ast.TypeAssertExpr{
									X: &ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "ctx",
											},
											Sel: &ast.Ident{
												Name: "Value",
											},
										},
										Args: []ast.Expr{
											&ast.Ident{
												Name: "UserKey",
											},
										},
									},
									Type: &ast.StarExpr{
										X: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "models",
											},
											Sel: &ast.Ident{
												Name: "User",
											},
										},
									},
								},
							},
						},
					},
				},
				&ast.IfStmt{
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
									&ast.Ident{
										Name: "nil",
									},
									&ast.CallExpr{
										Fun: &ast.Ident{
											Name: "decodeError",
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
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.Ident{
								Name: fmt.Sprintf("decode%s", h.model.ModelName()),
							},
							Args: []ast.Expr{
								&ast.Ident{
									Name: h.model.Variable(),
								},
							},
						},
						&ast.Ident{
							Name: "nil",
						},
					},
				},
			},
		},
	}
}

func (h Handler) syncGetMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, h.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == "Get" {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = h.astGetMethod()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}

	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(h.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (h Handler) astListMethod() *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{
							Name: "s",
						},
					},
					Type: &ast.StarExpr{
						X: &ast.Ident{
							Name: h.model.GRPCHandlerTypeName(),
						},
					},
				},
			},
		},
		Name: &ast.Ident{
			Name: "List",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
								Name: "ctx",
							},
						},
						Type: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "context",
							},
							Sel: &ast.Ident{
								Name: "Context",
							},
						},
					},
					{
						Names: []*ast.Ident{
							{
								Name: "filter",
							},
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: h.model.ProtoPackage,
								},
								Sel: &ast.Ident{
									Name: h.model.FilterTypeName(),
								},
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: h.model.ProtoPackage,
								},
								Sel: &ast.Ident{
									Name: fmt.Sprintf("List%s", h.model.ModelName()),
								},
							},
						},
					},
					{
						Type: &ast.Ident{
							Name: "error",
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: h.model.ListVariable(),
						},
						&ast.Ident{
							Name: "count",
						},
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
										Name: "s",
									},
									Sel: &ast.Ident{
										Name: h.model.InterceptorVariableName(),
									},
								},
								Sel: &ast.Ident{
									Name: "List",
								},
							},
							Args: []ast.Expr{
								&ast.Ident{
									Name: "ctx",
								},
								&ast.CallExpr{
									Fun: &ast.Ident{
										Name: fmt.Sprintf("encode%s", h.model.FilterTypeName()),
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "filter",
										},
									},
								},
								&ast.TypeAssertExpr{
									X: &ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "ctx",
											},
											Sel: &ast.Ident{
												Name: "Value",
											},
										},
										Args: []ast.Expr{
											&ast.Ident{
												Name: "UserKey",
											},
										},
									},
									Type: &ast.StarExpr{
										X: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "models",
											},
											Sel: &ast.Ident{
												Name: "User",
											},
										},
									},
								},
							},
						},
					},
				},
				&ast.IfStmt{
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
									&ast.Ident{
										Name: "nil",
									},
									&ast.CallExpr{
										Fun: &ast.Ident{
											Name: "decodeError",
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
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "header",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "metadata",
								},
								Sel: &ast.Ident{
									Name: "Pairs",
								},
							},
							Args: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: "\"count\"",
								},
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "fmt",
										},
										Sel: &ast.Ident{
											Name: "Sprint",
										},
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "count",
										},
									},
								},
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "_",
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "grpc",
								},
								Sel: &ast.Ident{
									Name: "SendHeader",
								},
							},
							Args: []ast.Expr{
								&ast.Ident{
									Name: "ctx",
								},
								&ast.Ident{
									Name: "header",
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.Ident{
								Name: fmt.Sprintf("decodeList%s", h.model.ModelName()),
							},
							Args: []ast.Expr{
								&ast.Ident{
									Name: h.model.ListVariable(),
								},
								&ast.Ident{
									Name: "count",
								},
							},
						},
						&ast.Ident{
							Name: "nil",
						},
					},
				},
			},
		},
	}
}

func (h Handler) syncListMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, h.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == "List" {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = h.astListMethod()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}

	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(h.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (h Handler) astUpdateMethod() *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{
							Name: "s",
						},
					},
					Type: &ast.StarExpr{
						X: &ast.Ident{
							Name: h.model.GRPCHandlerTypeName(),
						},
					},
				},
			},
		},
		Name: &ast.Ident{
			Name: "Update",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
								Name: "ctx",
							},
						},
						Type: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "context",
							},
							Sel: &ast.Ident{
								Name: "Context",
							},
						},
					},
					{
						Names: []*ast.Ident{
							{
								Name: "input",
							},
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: h.model.ProtoPackage,
								},
								Sel: &ast.Ident{
									Name: h.model.UpdateTypeName(),
								},
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: h.model.ProtoPackage,
								},
								Sel: &ast.Ident{
									Name: h.model.ModelName(),
								},
							},
						},
					},
					{
						Type: &ast.Ident{
							Name: "error",
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: h.model.Variable(),
						},
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
										Name: "s",
									},
									Sel: &ast.Ident{
										Name: h.model.InterceptorVariableName(),
									},
								},
								Sel: &ast.Ident{
									Name: "Update",
								},
							},
							Args: []ast.Expr{
								&ast.Ident{
									Name: "ctx",
								},
								&ast.CallExpr{
									Fun: &ast.Ident{
										Name: fmt.Sprintf("encode%s", h.model.UpdateTypeName()),
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "input",
										},
									},
								},
								&ast.TypeAssertExpr{
									X: &ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "ctx",
											},
											Sel: &ast.Ident{
												Name: "Value",
											},
										},
										Args: []ast.Expr{
											&ast.Ident{
												Name: "UserKey",
											},
										},
									},
									Type: &ast.StarExpr{
										X: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "models",
											},
											Sel: &ast.Ident{
												Name: "User",
											},
										},
									},
								},
							},
						},
					},
				},
				&ast.IfStmt{
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
									&ast.Ident{
										Name: "nil",
									},
									&ast.CallExpr{
										Fun: &ast.Ident{
											Name: "decodeError",
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
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.Ident{
								Name: fmt.Sprintf("decode%s", h.model.ModelName()),
							},
							Args: []ast.Expr{
								&ast.Ident{
									Name: h.model.Variable(),
								},
							},
						},
						&ast.Ident{
							Name: "nil",
						},
					},
				},
			},
		},
	}
}

func (h Handler) syncUpdateMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, h.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == "Update" {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = h.astUpdateMethod()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}

	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(h.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (h Handler) astDeleteMethod() *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{
							Name: "s",
						},
					},
					Type: &ast.StarExpr{
						X: &ast.Ident{
							Name: h.model.GRPCHandlerTypeName(),
						},
					},
				},
			},
		},
		Name: &ast.Ident{
			Name: "Delete",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
								Name: "ctx",
							},
						},
						Type: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "context",
							},
							Sel: &ast.Ident{
								Name: "Context",
							},
						},
					},
					{
						Names: []*ast.Ident{
							{
								Name: "input",
							},
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: h.model.ProtoPackage,
								},
								Sel: &ast.Ident{
									Name: fmt.Sprintf("%sDelete", h.model.ModelName()),
								},
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "emptypb",
								},
								Sel: &ast.Ident{
									Name: "Empty",
								},
							},
						},
					},
					{
						Type: &ast.Ident{
							Name: "error",
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
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "s",
										},
										Sel: &ast.Ident{
											Name: h.model.InterceptorVariableName(),
										},
									},
									Sel: &ast.Ident{
										Name: "Delete",
									},
								},
								Args: []ast.Expr{
									&ast.Ident{
										Name: "ctx",
									},
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "models",
											},
											Sel: &ast.Ident{
												Name: "UUID",
											},
										},
										Args: []ast.Expr{
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "input",
													},
													Sel: &ast.Ident{
														Name: "GetId",
													},
												},
											},
										},
									},
									&ast.TypeAssertExpr{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "ctx",
												},
												Sel: &ast.Ident{
													Name: "Value",
												},
											},
											Args: []ast.Expr{
												&ast.Ident{
													Name: "UserKey",
												},
											},
										},
										Type: &ast.StarExpr{
											X: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "models",
												},
												Sel: &ast.Ident{
													Name: "User",
												},
											},
										},
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
									&ast.Ident{
										Name: "nil",
									},
									&ast.CallExpr{
										Fun: &ast.Ident{
											Name: "decodeError",
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
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X: &ast.CompositeLit{
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "emptypb",
									},
									Sel: &ast.Ident{
										Name: "Empty",
									},
								},
							},
						},
						&ast.Ident{
							Name: "nil",
						},
					},
				},
			},
		},
	}
}

func (h Handler) syncDeleteMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, h.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == "Delete" {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = h.astDeleteMethod()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}

	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(h.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (h Handler) Sync() error {
	if err := h.syncStruct(); err != nil {
		return err
	}
	if err := h.syncConstructor(); err != nil {
		return err
	}
	if err := h.syncCreateMethod(); err != nil {
		return err
	}
	if err := h.syncGetMethod(); err != nil {
		return err
	}
	if err := h.syncListMethod(); err != nil {
		return err
	}
	if err := h.syncUpdateMethod(); err != nil {
		return err
	}
	if err := h.syncDeleteMethod(); err != nil {
		return err
	}

	if err := h.syncEncodeCreate(); err != nil {
		return err
	}
	if err := h.syncEncodeFilter(); err != nil {
		return err
	}
	if err := h.syncEncodeUpdate(); err != nil {
		return err
	}
	if err := h.syncDecodeModel(); err != nil {
		return err
	}
	if err := h.syncDecodeList(); err != nil {
		return err
	}
	if err := h.syncDecodeUpdate(); err != nil {
		return err
	}
	return nil
}
