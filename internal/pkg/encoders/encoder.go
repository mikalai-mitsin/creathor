package encoders

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

	"github.com/mikalai-mitsin/creathor/internal/pkg/astfile"
	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
)

type ProtoEncoder struct {
	domain configs.EntityConfig
}

func NewProtoEncoder(domain configs.EntityConfig) *ProtoEncoder {
	return &ProtoEncoder{domain: domain}
}

func (h ProtoEncoder) createParams() []ast.Expr {
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

func (h ProtoEncoder) encodeCreate() *ast.FuncDecl {
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
								X:   ast.NewIdent(h.domain.ProtoPackage),
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

func (h ProtoEncoder) syncEncodeCreate(filename string) error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	method, methodExist := astfile.FindFunc(file, h.domain.GetGRPCCreateDTOEncodeName())
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
	if err := os.WriteFile(filename, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (h ProtoEncoder) updateStmts() []*ast.IfStmt {
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

func (h ProtoEncoder) encodeUpdate() *ast.FuncDecl {
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
								X:   ast.NewIdent(h.domain.ProtoPackage),
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

func (h ProtoEncoder) syncEncodeUpdate(filename string) error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	method, methodExist := astfile.FindFunc(file, h.domain.GetGRPCUpdateDTOEncodeName())
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
	if err := os.WriteFile(filename, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (h ProtoEncoder) encodeFilter() *ast.FuncDecl {
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
							Value: &ast.CompositeLit{
								Type: &ast.ArrayType{
									Elt: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "entities",
										},
										Sel: &ast.Ident{
											Name: h.domain.OrderingTypeName(),
										},
									},
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
	stmts = append(stmts, &ast.RangeStmt{
		Key: &ast.Ident{
			Name: "_",
		},
		Value: &ast.Ident{
			Name: "orderBy",
		},
		Tok: token.DEFINE,
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X: &ast.Ident{
					Name: "input",
				},
				Sel: &ast.Ident{
					Name: "GetOrderBy",
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
								Name: "OrderBy",
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
										Name: "OrderBy",
									},
								},
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "entities",
										},
										Sel: &ast.Ident{
											Name: h.domain.OrderingTypeName(),
										},
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "orderBy",
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
								X:   ast.NewIdent(h.domain.ProtoPackage),
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

func (h ProtoEncoder) syncEncodeFilter(filename string) error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	method, methodExist := astfile.FindFunc(file, h.domain.GetGRPCFilterDTOEncodeName())
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
	if err := os.WriteFile(filename, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (h ProtoEncoder) file(pkg string) *ast.File {
	importSpec := []ast.Spec{
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: h.domain.EntitiesImportPath(),
			},
		},
		&ast.ImportSpec{
			Name: ast.NewIdent(h.domain.ProtoPackage),
			Path: &ast.BasicLit{
				Kind: token.STRING,
				Value: fmt.Sprintf(
					`"%s/pkg/%s/v1"`,
					h.domain.Module,
					h.domain.ProtoPackage,
				),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: h.domain.AppConfig.ProjectConfig.PointerImportPath(),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: h.domain.AppConfig.ProjectConfig.UUIDImportPath(),
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
		Name: ast.NewIdent(pkg),
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok:   token.IMPORT,
				Specs: importSpec,
			},
		},
	}
}

func (h ProtoEncoder) Sync() error {
	files := map[string]string{
		path.Join(
			"internal",
			"app",
			h.domain.AppConfig.AppName(),
			"handlers",
			"grpc",
			h.domain.DirName(),
			"dto.go",
		): "handlers",
		//path.Join(
		//	"internal",
		//	"app",
		//	h.domain.AppConfig.AppName(),
		//	"handlers",
		//	"kafka",
		//	h.domain.DirName(),
		//	"dto.go",
		//): "handlers",
	}
	for filename, pkg := range files {
		if err := h.sync(pkg, filename); err != nil {
			return err
		}
	}
	return nil
}

func (h ProtoEncoder) sync(pkg, filename string) error {
	fileset := token.NewFileSet()
	if err := os.MkdirAll(path.Dir(filename), 0777); err != nil {
		return err
	}
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = h.file(pkg)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(filename, buff.Bytes(), 0777); err != nil {
		return err
	}
	if err := h.syncEncodeCreate(filename); err != nil {
		return err
	}
	if err := h.syncEncodeFilter(filename); err != nil {
		return err
	}
	if err := h.syncEncodeUpdate(filename); err != nil {
		return err
	}
	return nil
}
