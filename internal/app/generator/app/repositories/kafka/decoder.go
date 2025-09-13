package kafka

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
							ast.NewIdent(h.domain.GetOneVariableName()),
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
	for _, param := range h.domain.GetMainModel().Params {
		var value ast.Expr
		value = &ast.SelectorExpr{
			X:   ast.NewIdent(h.domain.GetOneVariableName()),
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
						X:   ast.NewIdent(h.domain.GetOneVariableName()),
						Sel: ast.NewIdent(param.GetName()),
					},
				},
			}
		}
		if param.IsID() {
			value = &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent(h.domain.GetOneVariableName()),
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
	method, methodExist := astfile.FindFunc(file, h.domain.GetGRPCMainDecodeName())
	if method == nil {
		method = h.decode()
	}
	mainModelLit, _ := astfile.FindStructInstance(method, h.domain.GetMainModel().Name)
	for _, param := range h.modelParams() {
		astfile.SetParamValue(mainModelLit, param)
	}
	for _, param := range h.domain.GetMainModel().Params {
		param := *param
		if param.IsOptional() {
			var value ast.Expr
			param.Type = strings.TrimPrefix(param.Type, "*")
			value = &ast.StarExpr{
				X: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: h.domain.GetOneVariableName(),
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
							Name: h.domain.GetOneVariableName(),
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

func (h ProtoDecoder) file() *ast.File {
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
		Name: ast.NewIdent("repositories"),
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
		h.domain.AppConfig.AppName(),
		"repositories",
		"kafka",
		h.domain.DirName(),
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
	return nil
}
