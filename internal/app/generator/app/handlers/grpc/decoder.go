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

	"github.com/mikalai-mitsin/creathor/internal/pkg/astfile"
	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
)

type ProtoDecoder struct {
	entityConfig configs.EntityConfig
}

func NewProtoDecoder(entityConfig configs.EntityConfig) *ProtoDecoder {
	return &ProtoDecoder{
		entityConfig: entityConfig,
	}
}

func (h ProtoDecoder) decode() *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: ast.NewIdent(h.entityConfig.GetGRPCMainDecodeName()),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent(h.entityConfig.GetOneVariableName()),
						},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("entities"),
							Sel: ast.NewIdent(h.entityConfig.GetMainModel().Name),
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent(h.entityConfig.ProtoPackage),
								Sel: ast.NewIdent(h.entityConfig.GetMainModel().Name),
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
									X:   ast.NewIdent(h.entityConfig.ProtoPackage),
									Sel: ast.NewIdent(h.entityConfig.GetMainModel().Name),
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

func (h ProtoDecoder) modelParams() []*ast.KeyValueExpr {
	var exprs []*ast.KeyValueExpr
	for _, param := range h.entityConfig.GetMainModel().Params {
		var value ast.Expr
		value = &ast.SelectorExpr{
			X:   ast.NewIdent(h.entityConfig.GetOneVariableName()),
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
						X:   ast.NewIdent(h.entityConfig.GetOneVariableName()),
						Sel: ast.NewIdent(param.GetName()),
					},
				},
			}
		}
		if param.IsID() {
			value = &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent(h.entityConfig.GetOneVariableName()),
						Sel: ast.NewIdent(param.GetName()),
					},
					Sel: ast.NewIdent("String"),
				},
			}
		}
		if param.IsOptional() {
			value = ast.NewIdent("nil")
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
	method, methodExist := astfile.FindFunc(file, h.entityConfig.GetGRPCMainDecodeName())
	if method == nil {
		method = h.decode()
	}
	mainModelLit, _ := astfile.FindStructInstance(method, h.entityConfig.GetMainModel().Name)
	for _, param := range h.modelParams() {
		astfile.SetParamValue(mainModelLit, param)
	}
	for _, param := range h.entityConfig.GetMainModel().Params {
		param := *param
		if param.IsOptional() {
			var value ast.Expr
			param.Type = strings.TrimPrefix(param.Type, "*")
			value = &ast.StarExpr{
				X: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: h.entityConfig.GetOneVariableName(),
					},
					Sel: &ast.Ident{
						Name: param.GetName(),
					},
				},
			}
			if param.Type != param.GRPCType() {
				value = &ast.CallExpr{
					Fun: ast.NewIdent(param.GRPCType()),
					Args: []ast.Expr{
						value,
					},
				}
			}
			astfile.AppendToFuncBody(method, &ast.IfStmt{
				Cond: &ast.BinaryExpr{
					X: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: h.entityConfig.GetOneVariableName(),
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
								&ast.SelectorExpr{
									X: &ast.Ident{
										Name: "response",
									},
									Sel: &ast.Ident{
										Name: param.GRPCParam(),
									},
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								value,
							},
						},
					},
				},
			})
		}
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	var buff bytes.Buffer
	if err := printer.Fprint(&buff, fileset, file); err != nil {
		return err
	}
	return os.WriteFile(filename, buff.Bytes(), 0777)
}

func (h ProtoDecoder) decodeList() *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: ast.NewIdent(h.entityConfig.GetGRPCMainListDecodeName()),
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
								Sel: ast.NewIdent(h.entityConfig.GetMainModel().Name),
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
								X: ast.NewIdent(h.entityConfig.ProtoPackage),
								Sel: ast.NewIdent(
									fmt.Sprintf("List%s", h.entityConfig.GetMainModel().Name),
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
									X: ast.NewIdent(h.entityConfig.ProtoPackage),
									Sel: ast.NewIdent(
										fmt.Sprintf("List%s", h.entityConfig.GetMainModel().Name),
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
																h.entityConfig.ProtoPackage,
															),
															Sel: ast.NewIdent(
																h.entityConfig.GetMainModel().Name,
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
					Value: ast.NewIdent(h.entityConfig.GetOneVariableName()),
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
													Name: h.entityConfig.GetGRPCMainDecodeName(),
												},
												Args: []ast.Expr{
													ast.NewIdent(h.entityConfig.GetOneVariableName()),
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
	method, methodExist := astfile.FindFunc(file, h.entityConfig.GetGRPCMainListDecodeName())
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
							X:   ast.NewIdent(h.entityConfig.ProtoPackage),
							Sel: ast.NewIdent(h.entityConfig.GetUpdateModel().Name),
						},
						Elts: h.decodeUpdateParams(),
					},
				},
			},
		},
	}
	for _, param := range h.entityConfig.GetUpdateModel().Params {
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
		Name: ast.NewIdent(h.entityConfig.GetGRPCUpdateDecodeName()),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("update"),
						},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("entities"),
							Sel: ast.NewIdent(h.entityConfig.GetUpdateModel().Name),
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent(h.entityConfig.ProtoPackage),
								Sel: ast.NewIdent(h.entityConfig.GetUpdateModel().Name),
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
	for _, param := range h.entityConfig.GetUpdateModel().Params {
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
	method, methodExist := astfile.FindFunc(file, h.entityConfig.GetGRPCUpdateDecodeName())
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

func (h ProtoDecoder) file() *ast.File {
	importSpec := []ast.Spec{
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: h.entityConfig.ImportPathEntities(),
			},
		},
		&ast.ImportSpec{
			Name: ast.NewIdent(h.entityConfig.ProtoPackage),
			Path: &ast.BasicLit{
				Kind: token.STRING,
				Value: fmt.Sprintf(
					`"%s/pkg/%s/v1"`,
					h.entityConfig.Module,
					h.entityConfig.ProtoPackage,
				),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: h.entityConfig.AppConfig.ProjectConfig.PointerImportPath(),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: h.entityConfig.AppConfig.ProjectConfig.UUIDImportPath(),
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
	for _, param := range h.entityConfig.GetUpdateModel().Params {
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

func (h ProtoDecoder) Sync() error {
	filename := path.Join(
		"internal",
		"app",
		h.entityConfig.AppConfig.AppName(),
		"handlers",
		"grpc",
		h.entityConfig.DirName(),
		"dto.go",
	)
	fileset := token.NewFileSet()
	if err := os.MkdirAll(path.Dir(filename), 0777); err != nil {
		return err
	}
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = h.file()
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
