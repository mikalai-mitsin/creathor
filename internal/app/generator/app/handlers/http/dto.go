package http

import (
	"bytes"
	"fmt"
	"github.com/mikalai-mitsin/creathor/internal/pkg/domain"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"
	"path/filepath"
)

type DTOGenerator struct {
	domain *domain.Domain
}

func NewDTOGenerator(domain *domain.Domain) *DTOGenerator {
	return &DTOGenerator{domain: domain}
}

func (g *DTOGenerator) filename() string {
	return filepath.Join(
		"internal",
		"app",
		g.domain.DirName(),
		"handlers",
		"http",
		"dto.go",
	)
}

func (g *DTOGenerator) Sync() error {
	fileset := token.NewFileSet()
	filename := g.filename()
	if err := os.MkdirAll(path.Dir(filename), 0777); err != nil {
		return err
	}
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = g.file()
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(filename, buff.Bytes(), 0777); err != nil {
		return err
	}
	if err := g.syncDTOStruct(); err != nil {
		return err
	}
	if err := g.syncDTOConstructor(); err != nil {
		return err
	}
	if err := g.syncDTOListType(); err != nil {
		return err
	}
	if err := g.syncListDTOConstructor(); err != nil {
		return err
	}
	if err := g.syncFilterDTOStruct(); err != nil {
		return err
	}
	if err := g.syncFilterDTOConstructor(); err != nil {
		return err
	}
	if err := g.syncFilterDTOToEntity(); err != nil {
		return err
	}
	if err := g.syncUpdateDTOStruct(); err != nil {
		return err
	}
	if err := g.syncUpdateDTOConstructor(); err != nil {
		return err
	}
	if err := g.syncUpdateDTOToEntity(); err != nil {
		return err
	}

	if err := g.syncCreateDTOStruct(); err != nil {
		return err
	}
	if err := g.syncCreateDTOConstructor(); err != nil {
		return err
	}
	if err := g.syncCreateDTOToEntity(); err != nil {
		return err
	}
	return nil
}

func (g *DTOGenerator) file() *ast.File {
	return &ast.File{
		Name: ast.NewIdent("handlers"),
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"context"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"fmt"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"time"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/pkg/errs"`, g.domain.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: g.domain.EntitiesImportPath(),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/pkg/pointer"`, g.domain.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/pkg/log"`, g.domain.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/pkg/uuid"`, g.domain.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"github.com/go-chi/chi/v5"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"github.com/go-chi/render"`,
						},
					},
				},
			},
		},
	}
}

// Item DTO

func (g *DTOGenerator) dtoStruct() *ast.TypeSpec {
	structure := &ast.TypeSpec{
		Name: ast.NewIdent(g.domain.GetHTTPItemDTOName()),
		Type: &ast.StructType{
			Fields: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
								Name: "ID",
							},
						},
						Type: &ast.Ident{
							Name: "uuid.UUID",
						},
						Tag: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "`json:\"id\"`",
						},
					},
					{
						Names: []*ast.Ident{
							{
								Name: "UpdatedAt",
							},
						},
						Type: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "time",
							},
							Sel: &ast.Ident{
								Name: "Time",
							},
						},
						Tag: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "`json:\"updated_at\"`",
						},
					},
					{
						Names: []*ast.Ident{
							{
								Name: "CreatedAt",
							},
						},
						Type: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "time",
							},
							Sel: &ast.Ident{
								Name: "Time",
							},
						},
						Tag: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "`json:\"created_at\"`",
						},
					},
				},
			},
		},
	}
	for _, param := range g.domain.GetMainModel().Params {
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
					Type:  ast.NewIdent(param.JsonType()),
					Tag: &ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf("`json:\"%s\"`", param.Tag()),
					},
					Comment: nil,
				})
				return false
			}
			return true
		})
	}
	return structure
}

func (g *DTOGenerator) syncDTOStruct() error {
	fileset := token.NewFileSet()
	filename := g.filename()
	if err := os.MkdirAll(path.Dir(filename), 0777); err != nil {
		return err
	}
	file, err := parser.ParseFile(fileset, g.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var structureExists bool
	var structure *ast.TypeSpec
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok && t.Name.String() == g.domain.GetHTTPItemDTOName() {
			structure = t
			structureExists = true
			return false
		}
		return true
	})
	if structure == nil {
		structure = g.dtoStruct()
	}
	for _, param := range g.domain.GetMainModel().Params {
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
					Type:  ast.NewIdent(param.JsonType()),
					Tag: &ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf("`json:\"%s\"`", param.Tag()),
					},
					Comment: nil,
				})
				return false
			}
			return true
		})
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
	if err := os.WriteFile(g.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (g *DTOGenerator) dtoConstructor() *ast.FuncDecl {
	dto := &ast.CompositeLit{
		Type: ast.NewIdent(g.domain.GetHTTPItemDTOName()),
		Elts: []ast.Expr{},
	}
	for _, param := range g.domain.GetMainModel().Params {
		elt := &ast.KeyValueExpr{
			Key:   ast.NewIdent(param.GetName()),
			Value: nil,
		}
		if param.IsSlice() {
			elt.Value = &ast.CompositeLit{
				Type: ast.NewIdent(param.JsonType()),
			}
		} else {
			if param.JsonType() == param.Type {
				elt.Value = &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "entity",
					},
					Sel: &ast.Ident{
						Name: param.GetName(),
					},
				}
			} else {
				elt.Value = &ast.CallExpr{
					Fun: &ast.Ident{
						Name: param.JsonType(),
					},
					Args: []ast.Expr{
						&ast.SelectorExpr{
							X: &ast.Ident{
								Name: "entity",
							},
							Sel: &ast.Ident{
								Name: param.GetName(),
							},
						},
					},
				}
			}
		}
		dto.Elts = append(dto.Elts, elt)
	}
	constructor := &ast.FuncDecl{
		Name: &ast.Ident{
			Name: g.domain.GetHTTPItemDTOConstructorName(),
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
								Name: "entity",
							},
						},
						Type: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "entities",
							},
							Sel: &ast.Ident{
								Name: g.domain.GetMainModel().Name,
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent(g.domain.GetHTTPItemDTOName()),
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
						&ast.Ident{
							Name: "dto",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						dto,
					},
				},
			},
		},
	}
	for _, param := range g.domain.GetMainModel().Params {
		if !param.IsSlice() {
			continue
		}
		var valueToAppend ast.Expr
		if param.SliceType() == param.PostgresDTOSliceType() {
			valueToAppend = ast.NewIdent("param")
		} else {
			valueToAppend = &ast.CallExpr{
				Fun: &ast.Ident{
					Name: param.PostgresDTOSliceType(),
				},
				Args: []ast.Expr{
					ast.NewIdent("param"),
				},
			}
		}
		rang := &ast.RangeStmt{
			Key: &ast.Ident{
				Name: "_",
			},
			Value: ast.NewIdent("param"),
			Tok:   token.DEFINE,
			X: &ast.SelectorExpr{
				X: &ast.Ident{
					Name: "entity",
				},
				Sel: &ast.Ident{
					Name: param.GetName(),
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.SelectorExpr{
								X: &ast.Ident{
									Name: "dto",
								},
								Sel: &ast.Ident{
									Name: param.GetName(),
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
											Name: "dto",
										},
										Sel: &ast.Ident{
											Name: param.GetName(),
										},
									},
									valueToAppend,
								},
							},
						},
					},
				},
			},
		}
		constructor.Body.List = append(constructor.Body.List, rang)
	}
	constructor.Body.List = append(
		constructor.Body.List,
		&ast.ReturnStmt{
			Results: []ast.Expr{
				ast.NewIdent("dto"),
				ast.NewIdent("nil"),
			},
		},
	)
	return constructor
}

func (g *DTOGenerator) syncDTOConstructor() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, g.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var structureConstructorExists bool
	var structureConstructor *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok &&
			t.Name.String() == g.domain.GetHTTPItemDTOConstructorName() {
			structureConstructorExists = true
			structureConstructor = t
			return false
		}
		return true
	})
	if structureConstructor == nil {
		structureConstructor = g.dtoConstructor()
	}
	for _, param := range g.domain.GetMainModel().Params {
		param := param
		ast.Inspect(structureConstructor, func(node ast.Node) bool {
			if cl, ok := node.(*ast.CompositeLit); ok {
				if i, ok := cl.Type.(*ast.Ident); ok &&
					i.String() == g.domain.GetHTTPItemDTOName() {
					_ = i
					for _, elt := range cl.Elts {
						elt := elt
						if kv, ok := elt.(*ast.KeyValueExpr); ok {
							if key, ok := kv.Key.(*ast.Ident); ok &&
								key.String() == param.GetName() {
								return false
							}
						}
					}
					elt := &ast.KeyValueExpr{
						Key: &ast.Ident{
							Name: param.GetName(),
						},
						Value: nil,
					}
					if param.IsSlice() {
						elt.Value = &ast.CompositeLit{
							Type: ast.NewIdent(param.PostgresDTOType()),
						}
					} else {
						if param.PostgresDTOType() == param.Type {
							elt.Value = &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "entity",
								},
								Sel: &ast.Ident{
									Name: param.GetName(),
								},
							}
						} else {
							elt.Value = &ast.CallExpr{
								Fun: &ast.Ident{
									Name: param.PostgresDTOType(),
								},
								Args: []ast.Expr{
									&ast.SelectorExpr{
										X: &ast.Ident{
											Name: "entity",
										},
										Sel: &ast.Ident{
											Name: param.GetName(),
										},
									},
								},
							}
						}
					}
					cl.Elts = append(cl.Elts, elt)
					return false
				}
			}
			return true
		})
	}
	// TODO: add range sync
	if !structureConstructorExists {
		file.Decls = append(file.Decls, structureConstructor)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(g.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

// List DTO

func (g *DTOGenerator) astDTOListType() *ast.TypeSpec {
	return &ast.TypeSpec{
		Name: &ast.Ident{
			Name: g.domain.GetHTTPListDTOName(),
		},
		Type: &ast.StructType{
			Fields: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("Items"),
						},
						Type: &ast.ArrayType{
							Elt: &ast.Ident{
								Name: g.domain.GetHTTPItemDTOName(),
							},
						},
						Tag: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "`json:\"items\"`",
						},
					},
					{
						Names: []*ast.Ident{
							ast.NewIdent("Count"),
						},
						Type: ast.NewIdent("uint64"),
						Tag: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "`json:\"count\"`",
						},
					},
				},
			},
		},
	}
}

func (g *DTOGenerator) syncDTOListType() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, g.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var structureExists bool
	var dtoListType *ast.TypeSpec
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok &&
			t.Name.String() == g.domain.GetHTTPListDTOName() {
			dtoListType = t
			structureExists = true
			return false
		}
		return true
	})
	if dtoListType == nil {
		dtoListType = g.astDTOListType()
	}
	if !structureExists {
		gd := &ast.GenDecl{
			Tok:   token.TYPE,
			Specs: []ast.Spec{dtoListType},
		}
		file.Decls = append(file.Decls, gd)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(g.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (g *DTOGenerator) listDTOConstructor() *ast.FuncDecl {
	stmts := []ast.Stmt{
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent("response"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CompositeLit{
					Type: ast.NewIdent(g.domain.GetHTTPListDTOName()),
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
										Elt: ast.NewIdent(g.domain.GetHTTPItemDTOName()),
									},
									&ast.CallExpr{
										Fun: &ast.Ident{
											Name: "len",
										},
										Args: []ast.Expr{
											ast.NewIdent(g.domain.GetManyVariableName()),
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
		&ast.RangeStmt{
			Key: &ast.Ident{
				Name: "i",
			},
			Value: ast.NewIdent(g.domain.GetOneVariableName()),
			Tok:   token.DEFINE,
			X:     ast.NewIdent(g.domain.GetManyVariableName()),
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.Ident{
								Name: "dto",
							},
							&ast.Ident{
								Name: "err",
							},
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: ast.NewIdent(g.domain.GetHTTPItemDTOConstructorName()),
								Args: []ast.Expr{
									ast.NewIdent(g.domain.GetOneVariableName()),
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
										&ast.CompositeLit{
											Type: ast.NewIdent(g.domain.GetHTTPListDTOName()),
										},
										ast.NewIdent("err"),
									},
								},
							},
						},
					},
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.IndexExpr{
								X: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "response",
									},
									Sel: &ast.Ident{
										Name: "Items",
									},
								},
								Index: &ast.Ident{
									Name: "i",
								},
							},
						},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							&ast.Ident{
								Name: "dto",
							},
						},
					},
				},
			},
		},
		&ast.ReturnStmt{
			Results: []ast.Expr{
				ast.NewIdent("response"),
				ast.NewIdent("nil"),
			},
		},
	}
	return &ast.FuncDecl{
		Name: &ast.Ident{
			Name: g.domain.GetHTTPListDTOConstructorName(),
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent(g.domain.GetManyVariableName()),
						},
						Type: &ast.ArrayType{
							Elt: &ast.SelectorExpr{
								X:   ast.NewIdent("entities"),
								Sel: ast.NewIdent(g.domain.GetMainModel().Name),
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
						Type: ast.NewIdent(g.domain.GetHTTPListDTOName()),
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: stmts,
		},
	}
}

func (g *DTOGenerator) syncListDTOConstructor() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, g.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok &&
			t.Name.String() == g.domain.GetHTTPFilterDTOConstructorName() {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = g.listDTOConstructor()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(g.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

// Filter DTO
func (g *DTOGenerator) filterDTOStruct() *ast.TypeSpec {
	fields := &ast.FieldList{
		List: []*ast.Field{
			{
				Names: []*ast.Ident{
					{
						Name: "PageSize",
					},
				},
				Type: ast.NewIdent("*uint64"),
				Tag: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "`json:\"page_size\"`",
				},
			},
			{
				Names: []*ast.Ident{
					{
						Name: "PageNumber",
					},
				},
				Type: ast.NewIdent("*uint64"),
				Tag: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "`json:\"page_number\"`",
				},
			},
			{
				Names: []*ast.Ident{
					ast.NewIdent("OrderBy"),
				},
				Type: ast.NewIdent("[]string"),
				Tag: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "`json:\"order_by\"`",
				},
			},
			{
				Names: []*ast.Ident{
					ast.NewIdent("IDs"),
				},
				Type: ast.NewIdent("[]uuid.UUID"),
				Tag: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "`json:\"ids\"`",
				},
			},
		},
	}
	if g.domain.SearchEnabled() {
		fields.List = append(fields.List,
			&ast.Field{
				Names: []*ast.Ident{
					ast.NewIdent("Search"),
				},
				Type: ast.NewIdent("string"),
				Tag: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "`json:\"search\"`",
				},
			})
	}
	structure := &ast.TypeSpec{
		Name: ast.NewIdent(g.domain.GetHTTPFilterDTOName()),
		Type: &ast.StructType{
			Fields: fields,
		},
	}
	return structure
}

func (g *DTOGenerator) syncFilterDTOStruct() error {
	fileset := token.NewFileSet()
	filename := g.filename()
	if err := os.MkdirAll(path.Dir(filename), 0777); err != nil {
		return err
	}
	file, err := parser.ParseFile(fileset, g.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var structureExists bool
	var structure *ast.TypeSpec
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok && t.Name.String() == g.domain.GetHTTPFilterDTOName() {
			structure = t
			structureExists = true
			return false
		}
		return true
	})
	if structure == nil {
		structure = g.filterDTOStruct()
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
	if err := os.WriteFile(g.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (g *DTOGenerator) filterDTOConstructor() *ast.FuncDecl {
	exprs := []ast.Expr{
		&ast.KeyValueExpr{
			Key: &ast.Ident{
				Name: "IDs",
			},
			Value: ast.NewIdent("nil"),
		},
		&ast.KeyValueExpr{
			Key: &ast.Ident{
				Name: "PageSize",
			},
			Value: ast.NewIdent("nil"),
		},
		&ast.KeyValueExpr{
			Key: &ast.Ident{
				Name: "PageNumber",
			},
			Value: ast.NewIdent("nil"),
		},
		&ast.KeyValueExpr{
			Key: &ast.Ident{
				Name: "OrderBy",
			},
			Value: &ast.Ident{
				Name: "nil",
			},
		},
	}
	if g.domain.SearchEnabled() {
		exprs = append(exprs, &ast.KeyValueExpr{
			Key: &ast.Ident{
				Name: "Search",
			},
			Value: &ast.Ident{
				Name: `""`,
			},
		})
	}
	stmts := []ast.Stmt{
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: "filter",
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CompositeLit{
					Type: ast.NewIdent(g.domain.GetHTTPFilterDTOName()),
					Elts: exprs,
				},
			},
		},
		&ast.IfStmt{
			Cond: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "r",
								},
								Sel: &ast.Ident{
									Name: "URL",
								},
							},
							Sel: &ast.Ident{
								Name: "Query",
							},
						},
					},
					Sel: &ast.Ident{
						Name: "Has",
					},
				},
				Args: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: "\"page_size\"",
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.Ident{
								Name: "pageSize",
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
										Name: "strconv",
									},
									Sel: &ast.Ident{
										Name: "Atoi",
									},
								},
								Args: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "r",
														},
														Sel: &ast.Ident{
															Name: "URL",
														},
													},
													Sel: &ast.Ident{
														Name: "Query",
													},
												},
											},
											Sel: &ast.Ident{
												Name: "Get",
											},
										},
										Args: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: "\"page_size\"",
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
										&ast.CompositeLit{
											Type: ast.NewIdent(g.domain.GetHTTPFilterDTOName()),
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.CallExpr{
													Fun: &ast.SelectorExpr{
														X: &ast.CallExpr{
															Fun: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "errs",
																},
																Sel: &ast.Ident{
																	Name: "NewInvalidFormError",
																},
															},
														},
														Sel: &ast.Ident{
															Name: "WithParam",
														},
													},
													Args: []ast.Expr{
														&ast.BasicLit{
															Kind:  token.STRING,
															Value: "\"page_size\"",
														},
														&ast.BasicLit{
															Kind:  token.STRING,
															Value: "\"Invalid page_size.\"",
														},
													},
												},
												Sel: &ast.Ident{
													Name: "WithCause",
												},
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
									X:   ast.NewIdent("pointer"),
									Sel: ast.NewIdent("Pointer"),
								},
								Args: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.Ident{
											Name: "uint64",
										},
										Args: []ast.Expr{
											&ast.Ident{
												Name: "pageSize",
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
			Cond: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "r",
								},
								Sel: &ast.Ident{
									Name: "URL",
								},
							},
							Sel: &ast.Ident{
								Name: "Query",
							},
						},
					},
					Sel: &ast.Ident{
						Name: "Has",
					},
				},
				Args: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: "\"page_number\"",
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.Ident{
								Name: "pageNumber",
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
										Name: "strconv",
									},
									Sel: &ast.Ident{
										Name: "Atoi",
									},
								},
								Args: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "r",
														},
														Sel: &ast.Ident{
															Name: "URL",
														},
													},
													Sel: &ast.Ident{
														Name: "Query",
													},
												},
											},
											Sel: &ast.Ident{
												Name: "Get",
											},
										},
										Args: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: "\"page_number\"",
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
										&ast.CompositeLit{
											Type: ast.NewIdent(g.domain.GetHTTPFilterDTOName()),
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.CallExpr{
													Fun: &ast.SelectorExpr{
														X: &ast.CallExpr{
															Fun: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "errs",
																},
																Sel: &ast.Ident{
																	Name: "NewInvalidFormError",
																},
															},
														},
														Sel: &ast.Ident{
															Name: "WithParam",
														},
													},
													Args: []ast.Expr{
														&ast.BasicLit{
															Kind:  token.STRING,
															Value: "\"page_number\"",
														},
														&ast.BasicLit{
															Kind:  token.STRING,
															Value: "\"Invalid page_number.\"",
														},
													},
												},
												Sel: &ast.Ident{
													Name: "WithCause",
												},
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
									X:   ast.NewIdent("pointer"),
									Sel: ast.NewIdent("Pointer"),
								},
								Args: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.Ident{
											Name: "uint64",
										},
										Args: []ast.Expr{
											&ast.Ident{
												Name: "pageNumber",
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
			Cond: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "r",
								},
								Sel: &ast.Ident{
									Name: "URL",
								},
							},
							Sel: &ast.Ident{
								Name: "Query",
							},
						},
					},
					Sel: &ast.Ident{
						Name: "Has",
					},
				},
				Args: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: "\"order_by\"",
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
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "strings",
									},
									Sel: &ast.Ident{
										Name: "Split",
									},
								},
								Args: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "r",
														},
														Sel: &ast.Ident{
															Name: "URL",
														},
													},
													Sel: &ast.Ident{
														Name: "Query",
													},
												},
											},
											Sel: &ast.Ident{
												Name: "Get",
											},
										},
										Args: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: "\"order_by\"",
											},
										},
									},
									&ast.BasicLit{
										Kind:  token.STRING,
										Value: "\",\"",
									},
								},
							},
						},
					},
				},
			},
		},
		&ast.IfStmt{
			Cond: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "r",
								},
								Sel: &ast.Ident{
									Name: "URL",
								},
							},
							Sel: &ast.Ident{
								Name: "Query",
							},
						},
					},
					Sel: &ast.Ident{
						Name: "Has",
					},
				},
				Args: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: "\"ids\"",
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.Ident{
								Name: "ids",
							},
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "strings",
									},
									Sel: &ast.Ident{
										Name: "Split",
									},
								},
								Args: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "r",
														},
														Sel: &ast.Ident{
															Name: "URL",
														},
													},
													Sel: &ast.Ident{
														Name: "Query",
													},
												},
											},
											Sel: &ast.Ident{
												Name: "Get",
											},
										},
										Args: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: "\"ids\"",
											},
										},
									},
									&ast.BasicLit{
										Kind:  token.STRING,
										Value: "\",\"",
									},
								},
							},
						},
					},
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
									Name: "make",
								},
								Args: []ast.Expr{
									&ast.ArrayType{
										Elt: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "uuid",
											},
											Sel: &ast.Ident{
												Name: "UUID",
											},
										},
									},
									&ast.CallExpr{
										Fun: &ast.Ident{
											Name: "len",
										},
										Args: []ast.Expr{
											&ast.Ident{
												Name: "ids",
											},
										},
									},
									&ast.CallExpr{
										Fun: &ast.Ident{
											Name: "len",
										},
										Args: []ast.Expr{
											&ast.Ident{
												Name: "ids",
											},
										},
									},
								},
							},
						},
					},
					&ast.RangeStmt{
						Key: &ast.Ident{
							Name: "i",
						},
						Value: &ast.Ident{
							Name: "id",
						},
						Tok: token.DEFINE,
						X: &ast.Ident{
							Name: "ids",
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.AssignStmt{
									Lhs: []ast.Expr{
										&ast.IndexExpr{
											X: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "filter",
												},
												Sel: &ast.Ident{
													Name: "IDs",
												},
											},
											Index: &ast.Ident{
												Name: "i",
											},
										},
									},
									Tok: token.ASSIGN,
									Rhs: []ast.Expr{
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "uuid",
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
		},
	}
	if g.domain.SearchEnabled() {
		stmts = append(stmts, &ast.IfStmt{
			Cond: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "r",
								},
								Sel: &ast.Ident{
									Name: "URL",
								},
							},
							Sel: &ast.Ident{
								Name: "Query",
							},
						},
					},
					Sel: &ast.Ident{
						Name: "Has",
					},
				},
				Args: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: "\"search\"",
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
									Name: "Search",
								},
							},
						},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "r",
												},
												Sel: &ast.Ident{
													Name: "URL",
												},
											},
											Sel: &ast.Ident{
												Name: "Query",
											},
										},
									},
									Sel: &ast.Ident{
										Name: "Get",
									},
								},
								Args: []ast.Expr{
									&ast.BasicLit{
										Kind:  token.STRING,
										Value: "\"search\"",
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
			ast.NewIdent("nil"),
		},
	})
	return &ast.FuncDecl{
		Name: &ast.Ident{
			Name: g.domain.GetHTTPFilterDTOConstructorName(),
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
								Name: "r",
							},
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("http"),
								Sel: ast.NewIdent("Request"),
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent(g.domain.GetHTTPFilterDTOName()),
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: stmts,
		},
	}
}

func (g *DTOGenerator) syncFilterDTOConstructor() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, g.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok &&
			t.Name.String() == g.domain.GetHTTPFilterDTOConstructorName() {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = g.filterDTOConstructor()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(g.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (g *DTOGenerator) filterDTOToEntity() *ast.FuncDecl {
	exprs := []ast.Expr{
		&ast.KeyValueExpr{
			Key: &ast.Ident{
				Name: "PageSize",
			},
			Value: &ast.SelectorExpr{
				X: &ast.Ident{
					Name: "dto",
				},
				Sel: &ast.Ident{
					Name: "PageSize",
				},
			},
		},
		&ast.KeyValueExpr{
			Key: &ast.Ident{
				Name: "PageNumber",
			},
			Value: &ast.SelectorExpr{
				X: &ast.Ident{
					Name: "dto",
				},
				Sel: &ast.Ident{
					Name: "PageNumber",
				},
			},
		},
		&ast.KeyValueExpr{
			Key: &ast.Ident{
				Name: "OrderBy",
			},
			Value: &ast.SelectorExpr{
				X: &ast.Ident{
					Name: "dto",
				},
				Sel: &ast.Ident{
					Name: "OrderBy",
				},
			},
		},
		&ast.KeyValueExpr{
			Key: &ast.Ident{
				Name: "IDs",
			},
			Value: &ast.SelectorExpr{
				X: &ast.Ident{
					Name: "dto",
				},
				Sel: &ast.Ident{
					Name: "IDs",
				},
			},
		},
	}
	if g.domain.SearchEnabled() {
		exprs = append(exprs, &ast.KeyValueExpr{
			Key: &ast.Ident{
				Name: "Search",
			},
			Value: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "pointer",
					},
					Sel: &ast.Ident{
						Name: "Pointer",
					},
				},
				Args: []ast.Expr{
					&ast.SelectorExpr{
						X: &ast.Ident{
							Name: "dto",
						},
						Sel: &ast.Ident{
							Name: "Search",
						},
					},
				},
			},
		})
	}
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{
							Name: "dto",
						},
					},
					Type: ast.NewIdent(g.domain.GetHTTPFilterDTOName()),
				},
			},
		},
		Name: &ast.Ident{
			Name: "toEntity",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "entities",
							},
							Sel: &ast.Ident{
								Name: g.domain.GetFilterModel().Name,
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
						&ast.Ident{
							Name: "filter",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CompositeLit{
							Type: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "entities",
								},
								Sel: &ast.Ident{
									Name: "PostFilter",
								},
							},
							Elts: exprs,
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						ast.NewIdent("filter"),
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
}

func (g *DTOGenerator) syncFilterDTOToEntity() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, g.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExists bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == "toEntity" {
			if t.Recv.List[0].Type.(*ast.Ident).String() == g.domain.GetHTTPFilterDTOName() {
				methodExists = true
				method = t
				return false
			}
			return true
		}
		return true
	})
	if method == nil {
		method = g.filterDTOToEntity()
	}
	// TODO: add range sync
	if !methodExists {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(g.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

// Update DTO
func (g *DTOGenerator) updateDTOStruct() *ast.TypeSpec {
	structure := &ast.TypeSpec{
		Name: ast.NewIdent(g.domain.GetHTTPUpdateDTOName()),
		Type: &ast.StructType{
			Fields: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
								Name: "ID",
							},
						},
						Type: &ast.Ident{
							Name: "uuid.UUID",
						},
						Tag: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "`json:\"id\"`",
						},
					},
				},
			},
		},
	}
	for _, param := range g.domain.GetUpdateModel().Params {
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
					Type:  ast.NewIdent(param.JsonType()),
					Tag: &ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf("`json:\"%s\"`", param.Tag()),
					},
					Comment: nil,
				})
				return false
			}
			return true
		})
	}
	return structure
}

func (g *DTOGenerator) syncUpdateDTOStruct() error {
	fileset := token.NewFileSet()
	filename := g.filename()
	if err := os.MkdirAll(path.Dir(filename), 0777); err != nil {
		return err
	}
	file, err := parser.ParseFile(fileset, g.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var structureExists bool
	var structure *ast.TypeSpec
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok && t.Name.String() == g.domain.GetHTTPUpdateDTOName() {
			structure = t
			structureExists = true
			return false
		}
		return true
	})
	if structure == nil {
		structure = g.updateDTOStruct()
	}
	for _, param := range g.domain.GetUpdateModel().Params {
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
					Type:  ast.NewIdent(param.JsonType()),
					Tag: &ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf("`json:\"%s\"`", param.Tag()),
					},
					Comment: nil,
				})
				return false
			}
			return true
		})
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
	if err := os.WriteFile(g.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (g *DTOGenerator) updateDTOConstructor() *ast.FuncDecl {
	stmts := []ast.Stmt{
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: "update",
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CompositeLit{
					Type: ast.NewIdent(g.domain.GetHTTPUpdateDTOName()),
				},
			},
		},
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
							X: &ast.Ident{
								Name: "render",
							},
							Sel: &ast.Ident{
								Name: "DecodeJSON",
							},
						},
						Args: []ast.Expr{
							&ast.SelectorExpr{
								X: &ast.Ident{
									Name: "r",
								},
								Sel: &ast.Ident{
									Name: "Body",
								},
							},
							&ast.Ident{
								Name: "update",
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
							&ast.CompositeLit{
								Type: ast.NewIdent(g.domain.GetHTTPUpdateDTOName()),
							},
							ast.NewIdent("err"),
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
						Name: "ID",
					},
				},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: "uuid",
						},
						Sel: &ast.Ident{
							Name: "UUID",
						},
					},
					Args: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "chi",
								},
								Sel: &ast.Ident{
									Name: "URLParam",
								},
							},
							Args: []ast.Expr{
								&ast.Ident{
									Name: "r",
								},
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: "\"id\"",
								},
							},
						},
					},
				},
			},
		},
	}
	stmts = append(stmts, &ast.ReturnStmt{
		Results: []ast.Expr{
			ast.NewIdent("update"),
			ast.NewIdent("nil"),
		},
	})
	return &ast.FuncDecl{
		Name: &ast.Ident{
			Name: g.domain.GetHTTPUpdateDTOConstructorName(),
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
								Name: "r",
							},
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("http"),
								Sel: ast.NewIdent("Request"),
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent(g.domain.GetHTTPUpdateDTOName()),
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: stmts,
		},
	}
}

func (g *DTOGenerator) syncUpdateDTOConstructor() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, g.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok &&
			t.Name.String() == g.domain.GetHTTPUpdateDTOConstructorName() {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = g.updateDTOConstructor()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(g.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (g *DTOGenerator) updateDTOToEntity() *ast.FuncDecl {
	var exprs []ast.Expr
	for _, param := range g.domain.GetUpdateModel().Params {
		exprs = append(exprs, &ast.KeyValueExpr{
			Key: ast.NewIdent(param.GetName()),
			Value: &ast.SelectorExpr{
				X: &ast.Ident{
					Name: "dto",
				},
				Sel: ast.NewIdent(param.GetName()),
			},
		})
	}
	model := &ast.CompositeLit{
		Type: &ast.SelectorExpr{
			X: ast.NewIdent("entities"),
			Sel: &ast.Ident{
				Name: g.domain.GetUpdateModel().Name,
			},
		},
		Elts: exprs,
	}
	method := &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{
							Name: "dto",
						},
					},
					Type: ast.NewIdent(g.domain.GetHTTPUpdateDTOName()),
				},
			},
		},
		Name: &ast.Ident{
			Name: "toEntity",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "entities",
							},
							Sel: &ast.Ident{
								Name: g.domain.GetUpdateModel().Name,
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
						&ast.Ident{
							Name: "update",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						model,
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						ast.NewIdent("update"),
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
	return method
}

func (g *DTOGenerator) syncUpdateDTOToEntity() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, g.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExists bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == "toEntity" {
			if t.Recv.List[0].Type.(*ast.Ident).String() == g.domain.GetHTTPUpdateDTOName() {
				methodExists = true
				method = t
				return false
			}
			return true
		}
		return true
	})
	if method == nil {
		method = g.updateDTOToEntity()
	}
	// TODO: add range sync
	if !methodExists {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(g.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

// Create DTO
func (g *DTOGenerator) createDTOStruct() *ast.TypeSpec {
	structure := &ast.TypeSpec{
		Name: ast.NewIdent(g.domain.GetHTTPCreateDTOName()),
		Type: &ast.StructType{
			Fields: &ast.FieldList{
				List: []*ast.Field{},
			},
		},
	}
	for _, param := range g.domain.GetCreateModel().Params {
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
					Type:  ast.NewIdent(param.JsonType()),
					Tag: &ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf("`json:\"%s\"`", param.Tag()),
					},
					Comment: nil,
				})
				return false
			}
			return true
		})
	}
	return structure
}

func (g *DTOGenerator) syncCreateDTOStruct() error {
	fileset := token.NewFileSet()
	filename := g.filename()
	if err := os.MkdirAll(path.Dir(filename), 0777); err != nil {
		return err
	}
	file, err := parser.ParseFile(fileset, g.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var structureExists bool
	var structure *ast.TypeSpec
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok && t.Name.String() == g.domain.GetHTTPCreateDTOName() {
			structure = t
			structureExists = true
			return false
		}
		return true
	})
	if structure == nil {
		structure = g.createDTOStruct()
	}
	for _, param := range g.domain.GetCreateModel().Params {
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
					Type:  ast.NewIdent(param.JsonType()),
					Tag: &ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf("`json:\"%s\"`", param.Tag()),
					},
					Comment: nil,
				})
				return false
			}
			return true
		})
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
	if err := os.WriteFile(g.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (g *DTOGenerator) createDTOConstructor() *ast.FuncDecl {
	stmts := []ast.Stmt{
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: "create",
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CompositeLit{
					Type: ast.NewIdent(g.domain.GetHTTPCreateDTOName()),
					Elts: []ast.Expr{},
				},
			},
		},
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
							X: &ast.Ident{
								Name: "render",
							},
							Sel: &ast.Ident{
								Name: "DecodeJSON",
							},
						},
						Args: []ast.Expr{
							&ast.SelectorExpr{
								X: &ast.Ident{
									Name: "r",
								},
								Sel: &ast.Ident{
									Name: "Body",
								},
							},
							&ast.Ident{
								Name: "create",
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
							&ast.CompositeLit{
								Type: ast.NewIdent(g.domain.GetHTTPCreateDTOName()),
							},
							ast.NewIdent("err"),
						},
					},
				},
			},
		},
	}
	stmts = append(stmts, &ast.ReturnStmt{
		Results: []ast.Expr{
			ast.NewIdent("create"),
			ast.NewIdent("nil"),
		},
	})
	return &ast.FuncDecl{
		Name: &ast.Ident{
			Name: g.domain.GetHTTPCreateDTOConstructorName(),
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
								Name: "r",
							},
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("http"),
								Sel: ast.NewIdent("Request"),
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent(g.domain.GetHTTPCreateDTOName()),
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: stmts,
		},
	}
}

func (g *DTOGenerator) syncCreateDTOConstructor() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, g.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok &&
			t.Name.String() == g.domain.GetHTTPCreateDTOConstructorName() {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = g.createDTOConstructor()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(g.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (g *DTOGenerator) createDTOToEntity() *ast.FuncDecl {
	var exprs []ast.Expr
	for _, param := range g.domain.GetCreateModel().Params {
		exprs = append(exprs, &ast.KeyValueExpr{
			Key: ast.NewIdent(param.GetName()),
			Value: &ast.SelectorExpr{
				X: &ast.Ident{
					Name: "dto",
				},
				Sel: ast.NewIdent(param.GetName()),
			},
		})
	}
	method := &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{
							Name: "dto",
						},
					},
					Type: ast.NewIdent(g.domain.GetHTTPCreateDTOName()),
				},
			},
		},
		Name: &ast.Ident{
			Name: "toEntity",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "entities",
							},
							Sel: &ast.Ident{
								Name: g.domain.GetCreateModel().Name,
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
						&ast.Ident{
							Name: "create",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CompositeLit{
							Type: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "entities",
								},
								Sel: ast.NewIdent(g.domain.GetCreateModel().Name),
							},
							Elts: exprs,
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						ast.NewIdent("create"),
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
	return method
}

func (g *DTOGenerator) syncCreateDTOToEntity() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, g.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExists bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == "toEntity" {
			if t.Recv.List[0].Type.(*ast.Ident).String() == g.domain.GetHTTPCreateDTOName() {
				methodExists = true
				method = t
				return false
			}
			return true
		}
		return true
	})
	if method == nil {
		method = g.createDTOToEntity()
	}
	// TODO: add range sync
	if !methodExists {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(g.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}
