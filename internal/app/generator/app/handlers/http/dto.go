package http

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"
	"path/filepath"

	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
)

type DTOGenerator struct {
	domain *configs.EntityConfig
}

func NewDTOGenerator(domain *configs.EntityConfig) *DTOGenerator {
	return &DTOGenerator{domain: domain}
}

func (g *DTOGenerator) filename() string {
	return filepath.Join(
		"internal",
		"app",
		g.domain.AppConfig.AppName(),
		"handlers",
		"http",
		g.domain.DirName(),
		fmt.Sprintf("%s_dto.go", g.domain.SnakeName()),
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
							Value: `"time"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"net/http"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"strings"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"strconv"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: g.domain.AppConfig.ProjectConfig.ErrsImportPath(),
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
							Value: g.domain.AppConfig.ProjectConfig.PointerImportPath(),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: g.domain.AppConfig.ProjectConfig.UUIDImportPath(),
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
							ast.NewIdent("ID"),
						},
						Type: ast.NewIdent("uuid.UUID"),
						Tag: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "`json:\"id\"`",
						},
					},
					{
						Names: []*ast.Ident{
							ast.NewIdent("UpdatedAt"),
						},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("time"),
							Sel: ast.NewIdent("Time"),
						},
						Tag: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "`json:\"updated_at\"`",
						},
					},
					{
						Names: []*ast.Ident{
							ast.NewIdent("CreatedAt"),
						},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("time"),
							Sel: ast.NewIdent("Time"),
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
					X:   ast.NewIdent("entity"),
					Sel: ast.NewIdent(param.GetName()),
				}
			} else {
				elt.Value = &ast.CallExpr{
					Fun: ast.NewIdent(param.JsonType()),
					Args: []ast.Expr{
						&ast.SelectorExpr{
							X:   ast.NewIdent("entity"),
							Sel: ast.NewIdent(param.GetName()),
						},
					},
				}
			}
		}
		dto.Elts = append(dto.Elts, elt)
	}
	constructor := &ast.FuncDecl{
		Name: ast.NewIdent(g.domain.GetHTTPItemDTOConstructorName()),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("entity"),
						},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("entities"),
							Sel: ast.NewIdent(g.domain.GetMainModel().Name),
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
						ast.NewIdent("dto"),
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
				Fun: ast.NewIdent(param.PostgresDTOSliceType()),
				Args: []ast.Expr{
					ast.NewIdent("param"),
				},
			}
		}
		rang := &ast.RangeStmt{
			Key:   ast.NewIdent("_"),
			Value: ast.NewIdent("param"),
			Tok:   token.DEFINE,
			X: &ast.SelectorExpr{
				X:   ast.NewIdent("entity"),
				Sel: ast.NewIdent(param.GetName()),
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.SelectorExpr{
								X:   ast.NewIdent("dto"),
								Sel: ast.NewIdent(param.GetName()),
							},
						},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: ast.NewIdent("append"),
								Args: []ast.Expr{
									&ast.SelectorExpr{
										X:   ast.NewIdent("dto"),
										Sel: ast.NewIdent(param.GetName()),
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
						Key:   ast.NewIdent(param.GetName()),
						Value: nil,
					}
					if param.IsSlice() {
						elt.Value = &ast.CompositeLit{
							Type: ast.NewIdent(param.PostgresDTOType()),
						}
					} else {
						if param.PostgresDTOType() == param.Type {
							elt.Value = &ast.SelectorExpr{
								X:   ast.NewIdent("entity"),
								Sel: ast.NewIdent(param.GetName()),
							}
						} else {
							elt.Value = &ast.CallExpr{
								Fun: ast.NewIdent(param.PostgresDTOType()),
								Args: []ast.Expr{
									&ast.SelectorExpr{
										X:   ast.NewIdent("entity"),
										Sel: ast.NewIdent(param.GetName()),
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
		Name: ast.NewIdent(g.domain.GetHTTPListDTOName()),
		Type: &ast.StructType{
			Fields: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("Items"),
						},
						Type: &ast.ArrayType{
							Elt: ast.NewIdent(g.domain.GetHTTPItemDTOName()),
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
							Key: ast.NewIdent("Items"),
							Value: &ast.CallExpr{
								Fun: ast.NewIdent("make"),
								Args: []ast.Expr{
									&ast.ArrayType{
										Elt: ast.NewIdent(g.domain.GetHTTPItemDTOName()),
									},
									&ast.CallExpr{
										Fun: ast.NewIdent("len"),
										Args: []ast.Expr{
											ast.NewIdent(g.domain.GetManyVariableName()),
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
		&ast.RangeStmt{
			Key:   ast.NewIdent("i"),
			Value: ast.NewIdent(g.domain.GetOneVariableName()),
			Tok:   token.DEFINE,
			X:     ast.NewIdent(g.domain.GetManyVariableName()),
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							ast.NewIdent("dto"),
							ast.NewIdent("err"),
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
							X:  ast.NewIdent("err"),
							Op: token.NEQ,
							Y:  ast.NewIdent("nil"),
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
									X:   ast.NewIdent("response"),
									Sel: ast.NewIdent("Items"),
								},
								Index: ast.NewIdent("i"),
							},
						},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							ast.NewIdent("dto"),
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
		Name: ast.NewIdent(g.domain.GetHTTPListDTOConstructorName()),
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
					ast.NewIdent("PageSize"),
				},
				Type: ast.NewIdent("*uint64"),
				Tag: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "`json:\"page_size\"`",
				},
			},
			{
				Names: []*ast.Ident{
					ast.NewIdent("PageNumber"),
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
			Key:   ast.NewIdent("PageSize"),
			Value: ast.NewIdent("nil"),
		},
		&ast.KeyValueExpr{
			Key:   ast.NewIdent("PageNumber"),
			Value: ast.NewIdent("nil"),
		},
		&ast.KeyValueExpr{
			Key:   ast.NewIdent("OrderBy"),
			Value: ast.NewIdent("nil"),
		},
	}
	if g.domain.SearchEnabled() {
		exprs = append(exprs, &ast.KeyValueExpr{
			Key:   ast.NewIdent("Search"),
			Value: ast.NewIdent(`""`),
		})
	}
	stmts := []ast.Stmt{
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent("filter"),
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
								X:   ast.NewIdent("r"),
								Sel: ast.NewIdent("URL"),
							},
							Sel: ast.NewIdent("Query"),
						},
					},
					Sel: ast.NewIdent("Has"),
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
							ast.NewIdent("pageSize"),
							ast.NewIdent("err"),
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("strconv"),
									Sel: ast.NewIdent("Atoi"),
								},
								Args: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.SelectorExpr{
														X:   ast.NewIdent("r"),
														Sel: ast.NewIdent("URL"),
													},
													Sel: ast.NewIdent("Query"),
												},
											},
											Sel: ast.NewIdent("Get"),
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
							X:  ast.NewIdent("err"),
							Op: token.NEQ,
							Y:  ast.NewIdent("nil"),
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
																X:   ast.NewIdent("errs"),
																Sel: ast.NewIdent("NewInvalidFormError"),
															},
														},
														Sel: ast.NewIdent("WithParam"),
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
												Sel: ast.NewIdent("WithCause"),
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
										Fun: ast.NewIdent("uint64"),
										Args: []ast.Expr{
											ast.NewIdent("pageSize"),
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
								X:   ast.NewIdent("r"),
								Sel: ast.NewIdent("URL"),
							},
							Sel: ast.NewIdent("Query"),
						},
					},
					Sel: ast.NewIdent("Has"),
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
							ast.NewIdent("pageNumber"),
							ast.NewIdent("err"),
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("strconv"),
									Sel: ast.NewIdent("Atoi"),
								},
								Args: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.SelectorExpr{
														X:   ast.NewIdent("r"),
														Sel: ast.NewIdent("URL"),
													},
													Sel: ast.NewIdent("Query"),
												},
											},
											Sel: ast.NewIdent("Get"),
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
							X:  ast.NewIdent("err"),
							Op: token.NEQ,
							Y:  ast.NewIdent("nil"),
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
																X:   ast.NewIdent("errs"),
																Sel: ast.NewIdent("NewInvalidFormError"),
															},
														},
														Sel: ast.NewIdent("WithParam"),
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
												Sel: ast.NewIdent("WithCause"),
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
										Fun: ast.NewIdent("uint64"),
										Args: []ast.Expr{
											ast.NewIdent("pageNumber"),
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
								X:   ast.NewIdent("r"),
								Sel: ast.NewIdent("URL"),
							},
							Sel: ast.NewIdent("Query"),
						},
					},
					Sel: ast.NewIdent("Has"),
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
								X:   ast.NewIdent("filter"),
								Sel: ast.NewIdent("OrderBy"),
							},
						},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("strings"),
									Sel: ast.NewIdent("Split"),
								},
								Args: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.SelectorExpr{
														X:   ast.NewIdent("r"),
														Sel: ast.NewIdent("URL"),
													},
													Sel: ast.NewIdent("Query"),
												},
											},
											Sel: ast.NewIdent("Get"),
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
	}
	if g.domain.SearchEnabled() {
		stmts = append(stmts, &ast.IfStmt{
			Cond: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("r"),
								Sel: ast.NewIdent("URL"),
							},
							Sel: ast.NewIdent("Query"),
						},
					},
					Sel: ast.NewIdent("Has"),
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
								X:   ast.NewIdent("filter"),
								Sel: ast.NewIdent("Search"),
							},
						},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.SelectorExpr{
												X:   ast.NewIdent("r"),
												Sel: ast.NewIdent("URL"),
											},
											Sel: ast.NewIdent("Query"),
										},
									},
									Sel: ast.NewIdent("Get"),
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
		Name: ast.NewIdent(g.domain.GetHTTPFilterDTOConstructorName()),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("r"),
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
			Key: ast.NewIdent("PageSize"),
			Value: &ast.SelectorExpr{
				X:   ast.NewIdent("dto"),
				Sel: ast.NewIdent("PageSize"),
			},
		},
		&ast.KeyValueExpr{
			Key: ast.NewIdent("PageNumber"),
			Value: &ast.SelectorExpr{
				X:   ast.NewIdent("dto"),
				Sel: ast.NewIdent("PageNumber"),
			},
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
							Name: g.domain.OrderingTypeName(),
						},
					},
				},
			},
		},
	}
	if g.domain.SearchEnabled() {
		exprs = append(exprs, &ast.KeyValueExpr{
			Key: ast.NewIdent("Search"),
			Value: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("pointer"),
					Sel: ast.NewIdent("Of"),
				},
				Args: []ast.Expr{
					&ast.SelectorExpr{
						X:   ast.NewIdent("dto"),
						Sel: ast.NewIdent("Search"),
					},
				},
			},
		})
	}
	body := []ast.Stmt{
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent("filter"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CompositeLit{
					Type: &ast.SelectorExpr{
						X:   ast.NewIdent("entities"),
						Sel: ast.NewIdent(g.domain.GetFilterModel().Name),
					},
					Elts: exprs,
				},
			},
		},
	}
	body = append(body, &ast.RangeStmt{
		Key: &ast.Ident{
			Name: "_",
		},
		Value: &ast.Ident{
			Name: "orderBy",
		},
		Tok: token.DEFINE,
		X: &ast.SelectorExpr{
			X: &ast.Ident{
				Name: "dto",
			},
			Sel: &ast.Ident{
				Name: "OrderBy",
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
											Name: g.domain.OrderingTypeName(),
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
	body = append(body, &ast.ReturnStmt{
		Results: []ast.Expr{
			ast.NewIdent("filter"),
			ast.NewIdent("nil"),
		},
	})
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("dto"),
					},
					Type: ast.NewIdent(g.domain.GetHTTPFilterDTOName()),
				},
			},
		},
		Name: ast.NewIdent("toEntity"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("entities"),
							Sel: ast.NewIdent(g.domain.GetFilterModel().Name),
						},
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: body,
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
							ast.NewIdent("ID"),
						},
						Type: ast.NewIdent("uuid.UUID"),
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
				ast.NewIdent("update"),
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
					ast.NewIdent("err"),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("render"),
							Sel: ast.NewIdent("DecodeJSON"),
						},
						Args: []ast.Expr{
							&ast.SelectorExpr{
								X:   ast.NewIdent("r"),
								Sel: ast.NewIdent("Body"),
							},
							&ast.UnaryExpr{
								Op: token.AND,
								X:  ast.NewIdent("update"),
							},
						},
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
					X:   ast.NewIdent("update"),
					Sel: ast.NewIdent("ID"),
				},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("uuid"),
						Sel: ast.NewIdent("MustParse"),
					},
					Args: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("chi"),
								Sel: ast.NewIdent("URLParam"),
							},
							Args: []ast.Expr{
								ast.NewIdent("r"),
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
		Name: ast.NewIdent(g.domain.GetHTTPUpdateDTOConstructorName()),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("r"),
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
				X:   ast.NewIdent("dto"),
				Sel: ast.NewIdent(param.GetName()),
			},
		})
	}
	model := &ast.CompositeLit{
		Type: &ast.SelectorExpr{
			X:   ast.NewIdent("entities"),
			Sel: ast.NewIdent(g.domain.GetUpdateModel().Name),
		},
		Elts: exprs,
	}
	method := &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("dto"),
					},
					Type: ast.NewIdent(g.domain.GetHTTPUpdateDTOName()),
				},
			},
		},
		Name: ast.NewIdent("toEntity"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("entities"),
							Sel: ast.NewIdent(g.domain.GetUpdateModel().Name),
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
						ast.NewIdent("update"),
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
				ast.NewIdent("create"),
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
					ast.NewIdent("err"),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("render"),
							Sel: ast.NewIdent("DecodeJSON"),
						},
						Args: []ast.Expr{
							&ast.SelectorExpr{
								X:   ast.NewIdent("r"),
								Sel: ast.NewIdent("Body"),
							},
							&ast.UnaryExpr{
								Op: token.AND,
								X:  ast.NewIdent("create"),
							},
						},
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
		Name: ast.NewIdent(g.domain.GetHTTPCreateDTOConstructorName()),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("r"),
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
				X:   ast.NewIdent("dto"),
				Sel: ast.NewIdent(param.GetName()),
			},
		})
	}
	method := &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("dto"),
					},
					Type: ast.NewIdent(g.domain.GetHTTPCreateDTOName()),
				},
			},
		},
		Name: ast.NewIdent("toEntity"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("entities"),
							Sel: ast.NewIdent(g.domain.GetCreateModel().Name),
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
						ast.NewIdent("create"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CompositeLit{
							Type: &ast.SelectorExpr{
								X:   ast.NewIdent("entities"),
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
