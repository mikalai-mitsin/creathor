package grpc

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"
	"strings"

	"github.com/mikalai-mitsin/creathor/internal/pkg/app"
)

type HandlerGenerator struct {
	domain *app.BaseEntity
}

func NewHandlerGenerator(domain *app.BaseEntity) *HandlerGenerator {
	return &HandlerGenerator{
		domain: domain,
	}
}

func (h HandlerGenerator) file() *ast.File {
	importSpec := []ast.Spec{
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: `"context"`,
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: h.domain.EntitiesImportPath(),
			},
		},
		&ast.ImportSpec{
			Name: ast.NewIdent(h.domain.ProtoModule),
			Path: &ast.BasicLit{
				Kind: token.STRING,
				Value: fmt.Sprintf(
					`"%s/pkg/%s/v1"`,
					h.domain.Module,
					h.domain.ProtoModule,
				),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s/internal/pkg/pointer"`, h.domain.Module),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s/internal/pkg/uuid"`, h.domain.Module),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: `"google.golang.org/protobuf/types/known/emptypb"`,
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: `"google.golang.org/protobuf/types/known/timestamppb"`,
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: `"google.golang.org/protobuf/types/known/wrapperspb"`,
			},
		},
	}
	for _, param := range h.domain.GetUpdateModel().Params {
		if param.IsSlice() {
			importSpec = append(importSpec, &ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"google.golang.org/protobuf/types/known/structpb"`,
				},
			})
			break
		}
	}
	return &ast.File{
		Name: ast.NewIdent("handlers"),
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok:   token.IMPORT,
				Specs: importSpec,
			},
		},
	}
}

func (h HandlerGenerator) filename() string {
	return path.Join("internal", "app", h.domain.AppName(), "handlers", "grpc", h.domain.DirName(), h.domain.FileName())
}

func (h HandlerGenerator) createParams() []ast.Expr {
	var exprs []ast.Expr
	for _, param := range h.domain.GetCreateModel().Params {
		var value ast.Expr
		if param.IsSlice() {
			switch param.Type {
			case param.GRPCType():
				value = &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("input"),
						Sel: ast.NewIdent(param.GRPCGetter()),
					},
				}
			default:
				value = ast.NewIdent("nil")
			}
		} else {
			value = &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("input"),
					Sel: ast.NewIdent(param.GRPCGetter()),
				},
			}
			switch param.Type {
			case "time.Time":
				value = &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   value,
						Sel: ast.NewIdent("AsTime"),
					},
				}
			case param.GRPCType():
			default:
				if param.IsID() {
					value = &ast.CallExpr{
						Fun:  ast.NewIdent("uuid.MustParse"),
						Args: []ast.Expr{value},
					}
				} else {
					value = &ast.CallExpr{
						Fun: ast.NewIdent(param.Type),
						Args: []ast.Expr{
							value,
						},
					}
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

func (h HandlerGenerator) encodeCreate() *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: ast.NewIdent(fmt.Sprintf("encode%s", h.domain.GetCreateModel().Name)),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("input"),
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent(h.domain.ProtoModule),
								Sel: ast.NewIdent(h.domain.GetCreateModel().Name),
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("entities"),
							Sel: ast.NewIdent(h.domain.GetCreateModel().Name),
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("create"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CompositeLit{
							Type: &ast.SelectorExpr{
								X:   ast.NewIdent("entities"),
								Sel: ast.NewIdent(h.domain.GetCreateModel().Name),
							},
							Elts: h.createParams(),
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						ast.NewIdent("create"),
					},
				},
			},
		},
	}
}

func (h HandlerGenerator) syncEncodeCreate() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, h.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok &&
			t.Name.String() == fmt.Sprintf("encode%s", h.domain.GetCreateModel().Name) {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = h.encodeCreate()
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
		Key:   ast.NewIdent("_"),
		Value: ast.NewIdent("param"),
		Tok:   token.DEFINE,
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("input"),
				Sel: ast.NewIdent("GetTags"),
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X:   ast.NewIdent("create"),
							Sel: ast.NewIdent("Tags"),
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: ast.NewIdent("append"),
							Args: []ast.Expr{
								&ast.SelectorExpr{
									X:   ast.NewIdent("create"),
									Sel: ast.NewIdent("Tags"),
								},
								&ast.CallExpr{
									Fun: ast.NewIdent("string"),
									Args: []ast.Expr{
										ast.NewIdent("param"),
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

func (h HandlerGenerator) updateStmts() []*ast.IfStmt {
	var stmts []*ast.IfStmt
	for _, param := range h.domain.GetUpdateModel().Params {
		if param.GetName() == "ID" {
			continue
		}
		var body []ast.Stmt
		if param.Type == "*time.Time" || param.Type == "time.Time" {
			body = []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X:   ast.NewIdent("update"),
							Sel: ast.NewIdent(param.GetName()),
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("pointer"),
								Sel: ast.NewIdent("Of"),
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
							X:   ast.NewIdent("item"),
							Sel: ast.NewIdent(param.GrpcGetFromListValueAs()),
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
									ast.NewIdent("params"),
								},
								Type: &ast.ArrayType{
									Elt: ast.NewIdent(param.SliceType()),
								},
							},
						},
					},
				},
				&ast.RangeStmt{
					Key:   ast.NewIdent("_"),
					Value: ast.NewIdent("item"),
					Tok:   token.DEFINE,
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("input"),
									Sel: ast.NewIdent(param.GRPCGetter()),
								},
							},
							Sel: ast.NewIdent("GetValues"),
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									ast.NewIdent("params"),
								},
								Tok: token.ASSIGN,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: ast.NewIdent("append"),
										Args: []ast.Expr{
											ast.NewIdent("params"),
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
							X:   ast.NewIdent("update"),
							Sel: ast.NewIdent(param.GetName()),
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("pointer"),
								Sel: ast.NewIdent("Of"),
							},
							Args: []ast.Expr{
								ast.NewIdent("params"),
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
							X:   ast.NewIdent("input"),
							Sel: ast.NewIdent(param.GRPCGetter()),
						},
					},
					Sel: ast.NewIdent("GetValue"),
				},
			}
			if !param.IsID() && strings.TrimPrefix(param.Type, "*") != param.GRPCType() {
				value = &ast.CallExpr{
					Fun:  ast.NewIdent(strings.TrimPrefix(param.Type, "*")),
					Args: []ast.Expr{value},
				}
			}
			if param.IsID() {
				value = &ast.CallExpr{
					Fun:  ast.NewIdent("uuid.MustParse"),
					Args: []ast.Expr{value},
				}
			}
			body = []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X:   ast.NewIdent("update"),
							Sel: ast.NewIdent(param.GetName()),
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("pointer"),
								Sel: ast.NewIdent("Of"),
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
						X:   ast.NewIdent("input"),
						Sel: ast.NewIdent(param.GRPCGetter()),
					},
				},
				Op: token.NEQ,
				Y:  ast.NewIdent("nil"),
			},
			Body: &ast.BlockStmt{
				List: body,
			},
		})
	}
	return stmts
}

func (h HandlerGenerator) encodeUpdate() *ast.FuncDecl {
	body := []ast.Stmt{
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent("update"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CompositeLit{
					Type: &ast.SelectorExpr{
						X:   ast.NewIdent("entities"),
						Sel: ast.NewIdent(h.domain.GetUpdateModel().Name),
					},
					Elts: []ast.Expr{
						&ast.KeyValueExpr{
							Key: ast.NewIdent("ID"),
							Value: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("uuid"),
									Sel: ast.NewIdent("MustParse"),
								},
								Args: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   ast.NewIdent("input"),
											Sel: ast.NewIdent("GetId"),
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
			ast.NewIdent("update"),
		},
	})
	return &ast.FuncDecl{
		Name: ast.NewIdent(fmt.Sprintf("encode%s", h.domain.GetUpdateModel().Name)),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("input"),
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent(h.domain.ProtoModule),
								Sel: ast.NewIdent(h.domain.GetUpdateModel().Name),
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("entities"),
							Sel: ast.NewIdent(h.domain.GetUpdateModel().Name),
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

func (h HandlerGenerator) syncEncodeUpdate() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, h.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok &&
			t.Name.String() == fmt.Sprintf("encode%s", h.domain.GetUpdateModel().Name) {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = h.encodeUpdate()
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

func (h HandlerGenerator) encodeFilter() *ast.FuncDecl {
	stmts := []ast.Stmt{
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent("filter"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CompositeLit{
					Type: &ast.SelectorExpr{
						X:   ast.NewIdent("entities"),
						Sel: ast.NewIdent(h.domain.GetFilterModel().Name),
					},
					Elts: []ast.Expr{
						&ast.KeyValueExpr{
							Key:   ast.NewIdent("PageSize"),
							Value: ast.NewIdent("nil"),
						},
						&ast.KeyValueExpr{
							Key:   ast.NewIdent("PageNumber"),
							Value: ast.NewIdent("nil"),
						},
						&ast.KeyValueExpr{
							Key: ast.NewIdent("OrderBy"),
							Value: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("input"),
									Sel: ast.NewIdent("GetOrderBy"),
								},
							},
						},
						&ast.KeyValueExpr{
							Key:   ast.NewIdent("Search"),
							Value: ast.NewIdent("nil"),
						},
					},
				},
			},
		},
		&ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("input"),
						Sel: ast.NewIdent("GetPageSize"),
					},
				},
				Op: token.NEQ,
				Y:  ast.NewIdent("nil"),
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.SelectorExpr{
								X:   ast.NewIdent("filter"),
								Sel: ast.NewIdent("PageSize"),
							},
						},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("pointer"),
									Sel: ast.NewIdent("Of"),
								},
								Args: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X:   ast.NewIdent("input"),
													Sel: ast.NewIdent("GetPageSize"),
												},
											},
											Sel: ast.NewIdent("GetValue"),
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
						X:   ast.NewIdent("input"),
						Sel: ast.NewIdent("GetPageNumber"),
					},
				},
				Op: token.NEQ,
				Y:  ast.NewIdent("nil"),
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.SelectorExpr{
								X:   ast.NewIdent("filter"),
								Sel: ast.NewIdent("PageNumber"),
							},
						},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("pointer"),
									Sel: ast.NewIdent("Of"),
								},
								Args: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X:   ast.NewIdent("input"),
													Sel: ast.NewIdent("GetPageNumber"),
												},
											},
											Sel: ast.NewIdent("GetValue"),
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
	if h.domain.SearchEnabled() {
		stmts = append(stmts, &ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("input"),
						Sel: ast.NewIdent("GetSearch"),
					},
				},
				Op: token.NEQ,
				Y:  ast.NewIdent("nil"),
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.SelectorExpr{
								X:   ast.NewIdent("filter"),
								Sel: ast.NewIdent("Search"),
							},
						},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("pointer"),
									Sel: ast.NewIdent("Of"),
								},
								Args: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X:   ast.NewIdent("input"),
													Sel: ast.NewIdent("GetSearch"),
												},
											},
											Sel: ast.NewIdent("GetValue"),
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
			ast.NewIdent("filter"),
		},
	})
	return &ast.FuncDecl{
		Name: ast.NewIdent(fmt.Sprintf("encode%s", h.domain.GetFilterModel().Name)),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("input"),
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent(h.domain.ProtoModule),
								Sel: ast.NewIdent(h.domain.GetFilterModel().Name),
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("entities"),
							Sel: ast.NewIdent(h.domain.GetFilterModel().Name),
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

func (h HandlerGenerator) syncEncodeFilter() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, h.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok &&
			t.Name.String() == fmt.Sprintf("encode%s", h.domain.GetFilterModel().Name) {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = h.encodeFilter()
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

func (h HandlerGenerator) decode() *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: ast.NewIdent(fmt.Sprintf("decode%s", h.domain.GetMainModel().Name)),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("item"),
						},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("entities"),
							Sel: ast.NewIdent(h.domain.GetMainModel().Name),
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent(h.domain.ProtoModule),
								Sel: ast.NewIdent(h.domain.GetMainModel().Name),
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
						ast.NewIdent("response"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X: &ast.CompositeLit{
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent(h.domain.ProtoModule),
									Sel: ast.NewIdent(h.domain.GetMainModel().Name),
								},
								Elts: h.modelParams(),
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						ast.NewIdent("response"),
					},
				},
			},
		},
	}
}

func (h HandlerGenerator) modelParams() []ast.Expr {
	var exprs []ast.Expr
	for _, param := range h.domain.GetMainModel().Params {
		var value ast.Expr
		value = &ast.SelectorExpr{
			X:   ast.NewIdent("item"),
			Sel: ast.NewIdent(param.GetName()),
		}
		if param.Type != param.GRPCType() {
			value = &ast.CallExpr{
				Fun: ast.NewIdent(param.GRPCType()),
				Args: []ast.Expr{
					value,
				},
			}
		}
		if param.IsSlice() && param.GRPCType() != param.Type {
			value = &ast.CallExpr{
				Fun: &ast.IndexListExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent("pointer"),
						Sel: ast.NewIdent("ChangeType"),
					},
					Indices: []ast.Expr{
						ast.NewIdent(param.GRPCSliceType()),
						ast.NewIdent(param.SliceType()),
					},
				},
				Args: []ast.Expr{
					&ast.SelectorExpr{
						X:   ast.NewIdent("item"),
						Sel: ast.NewIdent(param.GetName()),
					},
				},
			}
		}
		if param.IsID() {
			value = &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent("item"),
						Sel: ast.NewIdent(param.GetName()),
					},
					Sel: ast.NewIdent("String"),
				},
			}
		}
		exprs = append(exprs, &ast.KeyValueExpr{
			Key:   ast.NewIdent(param.GRPCParam()),
			Value: value,
		})
	}
	return exprs
}

func (h HandlerGenerator) syncDecodeModel() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, h.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok &&
			t.Name.String() == fmt.Sprintf("decode%s", h.domain.GetMainModel().Name) {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = h.decode()
	}
	for _, param := range h.modelParams() {
		pr := param.(*ast.KeyValueExpr)
		prKey := pr.Key.(*ast.Ident)
		ast.Inspect(method, func(node ast.Node) bool {
			if cl, ok := node.(*ast.CompositeLit); ok {
				if clType, ok := cl.Type.(*ast.SelectorExpr); ok {
					if clType.Sel.String() != h.domain.GetMainModel().Name {
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

func (h HandlerGenerator) decodeList() *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: ast.NewIdent(fmt.Sprintf("decodeList%s", h.domain.GetMainModel().Name)),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("items"),
						},
						Type: &ast.ArrayType{
							Elt: &ast.SelectorExpr{
								X:   ast.NewIdent("entities"),
								Sel: ast.NewIdent(h.domain.GetMainModel().Name),
							},
						},
					},
					{
						Names: []*ast.Ident{
							ast.NewIdent("count"),
						},
						Type: ast.NewIdent("uint64"),
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent(h.domain.ProtoModule),
								Sel: ast.NewIdent(fmt.Sprintf("List%s", h.domain.GetMainModel().Name)),
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
						ast.NewIdent("response"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X: &ast.CompositeLit{
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent(h.domain.ProtoModule),
									Sel: ast.NewIdent(fmt.Sprintf("List%s", h.domain.GetMainModel().Name)),
								},
								Elts: []ast.Expr{
									&ast.KeyValueExpr{
										Key: ast.NewIdent("Items"),
										Value: &ast.CallExpr{
											Fun: ast.NewIdent("make"),
											Args: []ast.Expr{
												&ast.ArrayType{
													Elt: &ast.StarExpr{
														X: &ast.SelectorExpr{
															X:   ast.NewIdent(h.domain.ProtoModule),
															Sel: ast.NewIdent(h.domain.GetMainModel().Name),
														},
													},
												},
												&ast.BasicLit{
													Kind:  token.INT,
													Value: "0",
												},
												&ast.CallExpr{
													Fun: ast.NewIdent("len"),
													Args: []ast.Expr{
														ast.NewIdent("items"),
													},
												},
											},
										},
									},
									&ast.KeyValueExpr{
										Key:   ast.NewIdent("Count"),
										Value: ast.NewIdent("count"),
									},
								},
							},
						},
					},
				},
				&ast.RangeStmt{
					Key:   ast.NewIdent("_"),
					Value: ast.NewIdent("item"),
					Tok:   token.DEFINE,
					X:     ast.NewIdent("items"),
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.SelectorExpr{
										X:   ast.NewIdent("response"),
										Sel: ast.NewIdent("Items"),
									},
								},
								Tok: token.ASSIGN,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: ast.NewIdent("append"),
										Args: []ast.Expr{
											&ast.SelectorExpr{
												X:   ast.NewIdent("response"),
												Sel: ast.NewIdent("Items"),
											},
											&ast.CallExpr{
												Fun: &ast.Ident{
													Name: fmt.Sprintf(
														"decode%s",
														h.domain.GetMainModel().Name,
													),
												},
												Args: []ast.Expr{
													ast.NewIdent("item"),
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
						ast.NewIdent("response"),
					},
				},
			},
		},
	}
}

func (h HandlerGenerator) syncDecodeList() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, h.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok &&
			t.Name.String() == fmt.Sprintf("decodeList%s", h.domain.GetMainModel().Name) {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = h.decodeList()
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

func (h HandlerGenerator) decodeUpdate() *ast.FuncDecl {
	stmts := []ast.Stmt{
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent("result"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.UnaryExpr{
					Op: token.AND,
					X: &ast.CompositeLit{
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent(h.domain.ProtoModule),
							Sel: ast.NewIdent(h.domain.GetUpdateModel().Name),
						},
						Elts: h.decodeUpdateParams(),
					},
				},
			},
		},
	}
	for _, param := range h.domain.GetUpdateModel().Params {
		if !param.IsSlice() {
			continue
		}
		stmts = append(stmts, &ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X: &ast.SelectorExpr{
					X:   ast.NewIdent("update"),
					Sel: ast.NewIdent(param.GetName()),
				},
				Op: token.NEQ,
				Y:  ast.NewIdent("nil"),
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							ast.NewIdent("params"),
							ast.NewIdent("err"),
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("structpb"),
									Sel: ast.NewIdent("NewList"),
								},
								Args: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   ast.NewIdent("pointer"),
											Sel: ast.NewIdent("ToAnySlice"),
										},
										Args: []ast.Expr{
											&ast.StarExpr{
												X: &ast.SelectorExpr{
													X:   ast.NewIdent("update"),
													Sel: ast.NewIdent(param.GetName()),
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
							X:  ast.NewIdent("err"),
							Op: token.NEQ,
							Y:  ast.NewIdent("nil"),
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.ReturnStmt{
									Results: []ast.Expr{
										ast.NewIdent("nil"),
									},
								},
							},
						},
					},
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.SelectorExpr{
								X:   ast.NewIdent("result"),
								Sel: ast.NewIdent(param.GRPCParam()),
							},
						},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							ast.NewIdent("params"),
						},
					},
				},
			},
		})
	}
	stmts = append(stmts, &ast.ReturnStmt{
		Results: []ast.Expr{
			ast.NewIdent("result"),
		},
	})
	return &ast.FuncDecl{
		Name: ast.NewIdent(fmt.Sprintf("decode%s", h.domain.GetUpdateModel().Name)),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("update"),
						},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("entities"),
							Sel: ast.NewIdent(h.domain.GetUpdateModel().Name),
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent(h.domain.ProtoModule),
								Sel: ast.NewIdent(h.domain.GetUpdateModel().Name),
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

func (h HandlerGenerator) decodeUpdateParams() []ast.Expr {
	var exprs []ast.Expr
	for _, param := range h.domain.GetUpdateModel().Params {
		var value ast.Expr
		if param.IsSlice() {
			value = ast.NewIdent("nil")
		} else {
			var v ast.Expr
			v = &ast.SelectorExpr{
				X:   ast.NewIdent("update"),
				Sel: ast.NewIdent(param.GetName()),
			}
			if strings.HasPrefix(param.Type, "*") {
				v = &ast.StarExpr{
					X: v,
				}
			}
			if param.GetGRPCWrapperArgumentType() != strings.TrimPrefix(param.Type, "*") {
				v = &ast.CallExpr{
					Fun:  ast.NewIdent(param.GetGRPCWrapperArgumentType()),
					Args: []ast.Expr{v},
				}
			}
			if param.IsID() {
				v = &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.SelectorExpr{
							X:   ast.NewIdent("update"),
							Sel: ast.NewIdent(param.GetName()),
						},
						Sel: ast.NewIdent("String"),
					},
				}
			}
			value = &ast.CallExpr{
				Fun:  ast.NewIdent(param.GetGRPCWrapper()),
				Args: []ast.Expr{v},
			}
		}
		exprs = append(exprs, &ast.KeyValueExpr{
			Key:   ast.NewIdent(param.GRPCParam()),
			Value: value,
		})
	}
	return exprs
}

func (h HandlerGenerator) syncDecodeUpdate() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, h.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok &&
			t.Name.String() == fmt.Sprintf("decode%s", h.domain.GetUpdateModel().Name) {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = h.decodeUpdate()
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

func (h HandlerGenerator) structure() *ast.TypeSpec {
	return &ast.TypeSpec{
		Name: ast.NewIdent(h.domain.GetGRPCHandlerTypeName()),
		Type: &ast.StructType{
			Fields: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.SelectorExpr{
							X: ast.NewIdent(h.domain.ProtoModule),
							Sel: &ast.Ident{
								Name: fmt.Sprintf(
									"Unimplemented%sServiceServer",
									h.domain.GetMainModel().Name,
								),
							},
						},
					},
					{
						Names: []*ast.Ident{
							ast.NewIdent(h.domain.GetUseCasePrivateVariableName()),
						},
						Type: ast.NewIdent(h.domain.GetUseCaseInterfaceName()),
					},
					{
						Names: []*ast.Ident{
							ast.NewIdent("logger"),
						},
						Type: ast.NewIdent("logger"),
					},
				},
			},
		},
	}
}

func (h HandlerGenerator) syncStruct() error {
	fileset := token.NewFileSet()
	filename := h.filename()
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = h.file()
	}
	var structureExists bool
	var structure *ast.TypeSpec
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok &&
			t.Name.String() == h.domain.GetGRPCHandlerTypeName() {
			structure = t
			structureExists = true
			return false
		}
		return true
	})
	if structure == nil {
		structure = h.structure()
	}
	if !structureExists {
		gd := &ast.GenDecl{
			Tok:   token.TYPE,
			Specs: []ast.Spec{structure},
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

func (h HandlerGenerator) constructor() *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: ast.NewIdent(fmt.Sprintf("New%s", h.domain.GetGRPCHandlerTypeName())),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent(h.domain.GetUseCasePrivateVariableName()),
						},
						Type: ast.NewIdent(h.domain.GetUseCaseInterfaceName()),
					},
					{
						Names: []*ast.Ident{
							ast.NewIdent("logger"),
						},
						Type: ast.NewIdent("logger"),
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: ast.NewIdent(h.domain.GetGRPCHandlerTypeName()),
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
									Name: fmt.Sprintf(
										"%sServiceServer",
										h.domain.GetMainModel().Name,
									),
								},
								Elts: []ast.Expr{
									&ast.KeyValueExpr{
										Key:   ast.NewIdent(h.domain.GetUseCasePrivateVariableName()),
										Value: ast.NewIdent(h.domain.GetUseCasePrivateVariableName()),
									},
									&ast.KeyValueExpr{
										Key:   ast.NewIdent("logger"),
										Value: ast.NewIdent("logger"),
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

func (h HandlerGenerator) syncConstructor() error {
	fileset := token.NewFileSet()
	filename := h.filename()
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok &&
			t.Name.String() == fmt.Sprintf("New%s", h.domain.GetGRPCHandlerTypeName()) {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = h.constructor()
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

func (h HandlerGenerator) create() *ast.FuncDecl {
	args := []ast.Expr{
		ast.NewIdent("ctx"),
		&ast.CallExpr{
			Fun: ast.NewIdent(fmt.Sprintf("encode%s", h.domain.GetCreateModel().Name)),
			Args: []ast.Expr{
				ast.NewIdent("input"),
			},
		},
	}
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("s"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(h.domain.GetGRPCHandlerTypeName()),
					},
				},
			},
		},
		Name: ast.NewIdent("Create"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("ctx"),
						},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("context"),
							Sel: ast.NewIdent("Context"),
						},
					},
					{
						Names: []*ast.Ident{
							ast.NewIdent("input"),
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent(h.domain.ProtoModule),
								Sel: ast.NewIdent(h.domain.GetCreateModel().Name),
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
								X:   ast.NewIdent(h.domain.ProtoModule),
								Sel: ast.NewIdent(h.domain.GetMainModel().Name),
							},
						},
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("item"),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("s"),
									Sel: ast.NewIdent(h.domain.GetUseCasePrivateVariableName()),
								},
								Sel: ast.NewIdent("Create"),
							},
							Args: args,
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
									ast.NewIdent("nil"),
									ast.NewIdent("err"),
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: ast.NewIdent(fmt.Sprintf("decode%s", h.domain.GetMainModel().Name)),
							Args: []ast.Expr{
								ast.NewIdent("item"),
							},
						},
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
}

func (h HandlerGenerator) syncCreateMethod() error {
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
		method = h.create()
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

func (h HandlerGenerator) get() *ast.FuncDecl {
	args := []ast.Expr{
		ast.NewIdent("ctx"),
		&ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("uuid"),
				Sel: ast.NewIdent("MustParse"),
			},
			Args: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("input"),
						Sel: ast.NewIdent("GetId"),
					},
				},
			},
		},
	}
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("s"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(h.domain.GetGRPCHandlerTypeName()),
					},
				},
			},
		},
		Name: ast.NewIdent("Get"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("ctx"),
						},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("context"),
							Sel: ast.NewIdent("Context"),
						},
					},
					{
						Names: []*ast.Ident{
							ast.NewIdent("input"),
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent(h.domain.ProtoModule),
								Sel: ast.NewIdent(fmt.Sprintf("%sGet", h.domain.GetMainModel().Name)),
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
								X:   ast.NewIdent(h.domain.ProtoModule),
								Sel: ast.NewIdent(h.domain.GetMainModel().Name),
							},
						},
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("item"),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("s"),
									Sel: ast.NewIdent(h.domain.GetUseCasePrivateVariableName()),
								},
								Sel: ast.NewIdent("Get"),
							},
							Args: args,
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
									ast.NewIdent("nil"),
									ast.NewIdent("err"),
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: ast.NewIdent(fmt.Sprintf("decode%s", h.domain.GetMainModel().Name)),
							Args: []ast.Expr{
								ast.NewIdent("item"),
							},
						},
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
}

func (h HandlerGenerator) syncGetMethod() error {
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
		method = h.get()
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

func (h HandlerGenerator) list() *ast.FuncDecl {
	args := []ast.Expr{
		ast.NewIdent("ctx"),
		&ast.CallExpr{
			Fun: ast.NewIdent(fmt.Sprintf("encode%s", h.domain.GetFilterModel().Name)),
			Args: []ast.Expr{
				ast.NewIdent("filter"),
			},
		},
	}
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("s"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(h.domain.GetGRPCHandlerTypeName()),
					},
				},
			},
		},
		Name: ast.NewIdent("List"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("ctx"),
						},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("context"),
							Sel: ast.NewIdent("Context"),
						},
					},
					{
						Names: []*ast.Ident{
							ast.NewIdent("filter"),
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent(h.domain.ProtoModule),
								Sel: ast.NewIdent(h.domain.GetFilterModel().Name),
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
								X:   ast.NewIdent(h.domain.ProtoModule),
								Sel: ast.NewIdent(fmt.Sprintf("List%s", h.domain.GetMainModel().Name)),
							},
						},
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("items"),
						ast.NewIdent("count"),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("s"),
									Sel: ast.NewIdent(h.domain.GetUseCasePrivateVariableName()),
								},
								Sel: ast.NewIdent("List"),
							},
							Args: args,
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
									ast.NewIdent("nil"),
									ast.NewIdent("err"),
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: ast.NewIdent(fmt.Sprintf("decodeList%s", h.domain.GetMainModel().Name)),
							Args: []ast.Expr{
								ast.NewIdent("items"),
								ast.NewIdent("count"),
							},
						},
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
}

func (h HandlerGenerator) syncListMethod() error {
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
		method = h.list()
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

func (h HandlerGenerator) update() *ast.FuncDecl {
	args := []ast.Expr{
		ast.NewIdent("ctx"),
		&ast.CallExpr{
			Fun: ast.NewIdent(fmt.Sprintf("encode%s", h.domain.GetUpdateModel().Name)),
			Args: []ast.Expr{
				ast.NewIdent("input"),
			},
		},
	}
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("s"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(h.domain.GetGRPCHandlerTypeName()),
					},
				},
			},
		},
		Name: ast.NewIdent("Update"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("ctx"),
						},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("context"),
							Sel: ast.NewIdent("Context"),
						},
					},
					{
						Names: []*ast.Ident{
							ast.NewIdent("input"),
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent(h.domain.ProtoModule),
								Sel: ast.NewIdent(h.domain.GetUpdateModel().Name),
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
								X:   ast.NewIdent(h.domain.ProtoModule),
								Sel: ast.NewIdent(h.domain.GetMainModel().Name),
							},
						},
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("item"),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("s"),
									Sel: ast.NewIdent(h.domain.GetUseCasePrivateVariableName()),
								},
								Sel: ast.NewIdent("Update"),
							},
							Args: args,
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
									ast.NewIdent("nil"),
									ast.NewIdent("err"),
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: ast.NewIdent(fmt.Sprintf("decode%s", h.domain.GetMainModel().Name)),
							Args: []ast.Expr{
								ast.NewIdent("item"),
							},
						},
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
}

func (h HandlerGenerator) syncUpdateMethod() error {
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
		method = h.update()
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

func (h HandlerGenerator) delete() *ast.FuncDecl {
	args := []ast.Expr{
		ast.NewIdent("ctx"),
		&ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("uuid"),
				Sel: ast.NewIdent("MustParse"),
			},
			Args: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("input"),
						Sel: ast.NewIdent("GetId"),
					},
				},
			},
		},
	}
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("s"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(h.domain.GetGRPCHandlerTypeName()),
					},
				},
			},
		},
		Name: ast.NewIdent("Delete"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("ctx"),
						},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("context"),
							Sel: ast.NewIdent("Context"),
						},
					},
					{
						Names: []*ast.Ident{
							ast.NewIdent("input"),
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent(h.domain.ProtoModule),
								Sel: ast.NewIdent(fmt.Sprintf("%sDelete", h.domain.GetMainModel().Name)),
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
								X:   ast.NewIdent("emptypb"),
								Sel: ast.NewIdent("Empty"),
							},
						},
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.IfStmt{
					Init: &ast.AssignStmt{
						Lhs: []ast.Expr{
							ast.NewIdent("err"),
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.SelectorExpr{
										X:   ast.NewIdent("s"),
										Sel: ast.NewIdent(h.domain.GetUseCasePrivateVariableName()),
									},
									Sel: ast.NewIdent("Delete"),
								},
								Args: args,
							},
						},
					},
					Cond: &ast.BinaryExpr{
						X:  ast.NewIdent("err"),
						Op: token.NEQ,
						Y:  ast.NewIdent("nil"),
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{
									ast.NewIdent("nil"),
									ast.NewIdent("err"),
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
									X:   ast.NewIdent("emptypb"),
									Sel: ast.NewIdent("Empty"),
								},
							},
						},
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
}

func (h HandlerGenerator) syncDeleteMethod() error {
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
		method = h.delete()
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

func (h HandlerGenerator) Sync() error {
	err := os.MkdirAll(path.Dir(h.filename()), 0777)
	if err != nil {
		return err
	}
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
