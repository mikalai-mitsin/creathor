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

type ProtoDecoder struct {
	domain configs.EntityConfig
}

func NewProtoDecoder(domain configs.EntityConfig) *ProtoDecoder {
	return &ProtoDecoder{
		domain: domain,
	}
}

func (h ProtoDecoder) decode() *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: ast.NewIdent(h.domain.GetGRPCMainDecodeName()),
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
								X:   ast.NewIdent(h.domain.ProtoPackage),
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
									X:   ast.NewIdent(h.domain.ProtoPackage),
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

func (h ProtoDecoder) modelParams() []ast.Expr {
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

func (h ProtoDecoder) syncDecodeModel(filename string) error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	method, methodExist := astfile.FindFunc(file, h.domain.GetGRPCMainDecodeName())
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
	if err := os.WriteFile(filename, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (h ProtoDecoder) decodeList() *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: ast.NewIdent(h.domain.GetGRPCMainListDecodeName()),
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
								X: ast.NewIdent(h.domain.ProtoPackage),
								Sel: ast.NewIdent(
									fmt.Sprintf("List%s", h.domain.GetMainModel().Name),
								),
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
									X: ast.NewIdent(h.domain.ProtoPackage),
									Sel: ast.NewIdent(
										fmt.Sprintf("List%s", h.domain.GetMainModel().Name),
									),
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
															X: ast.NewIdent(
																h.domain.ProtoPackage,
															),
															Sel: ast.NewIdent(
																h.domain.GetMainModel().Name,
															),
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

func (h ProtoDecoder) syncDecodeList(filename string) error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	method, methodExist := astfile.FindFunc(file, h.domain.GetGRPCMainListDecodeName())
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
	if err := os.WriteFile(filename, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (h ProtoDecoder) decodeUpdate() *ast.FuncDecl {
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
							X:   ast.NewIdent(h.domain.ProtoPackage),
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
		Name: ast.NewIdent(h.domain.GetGRPCUpdateDecodeName()),
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
								X:   ast.NewIdent(h.domain.ProtoPackage),
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

func (h ProtoDecoder) decodeUpdateParams() []ast.Expr {
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

func (h ProtoDecoder) syncDecodeUpdate(filename string) error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	method, methodExist := astfile.FindFunc(file, h.domain.GetGRPCUpdateDecodeName())
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
	if err := os.WriteFile(filename, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (h ProtoDecoder) file(pkg string) *ast.File {
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

func (h ProtoDecoder) Sync() error {
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
		//	"repositories",
		//	"kafka",
		//	h.domain.DirName(),
		//	"dto.go",
		//): "repositories",
	}
	for filename, pkg := range files {
		if err := h.sync(pkg, filename); err != nil {
			return err
		}
	}
	return nil
}

func (h ProtoDecoder) sync(pkg, filename string) error {
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
	if err := h.syncDecodeModel(filename); err != nil {
		return err
	}
	if err := h.syncDecodeList(filename); err != nil {
		return err
	}
	if err := h.syncDecodeUpdate(filename); err != nil {
		return err
	}
	return nil
}
