package postgres

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

	"github.com/mikalai-mitsin/creathor/internal/pkg/astfile"
	"github.com/mikalai-mitsin/creathor/internal/pkg/tmpl"

	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"

	"github.com/iancoleman/strcase"
)

type RepositoryGenerator struct {
	domain *configs.EntityConfig
}

func NewRepositoryGenerator(domain *configs.EntityConfig) *RepositoryGenerator {
	return &RepositoryGenerator{domain: domain}
}

func (r RepositoryGenerator) getDTOName() string {
	return fmt.Sprintf("%sDTO", strcase.ToCamel(r.domain.GetMainModel().Name))
}

func (r RepositoryGenerator) getDTOListName() string {
	return fmt.Sprintf("%sListDTO", strcase.ToCamel(r.domain.GetMainModel().Name))
}

func (r RepositoryGenerator) filename() string {
	return filepath.Join(
		"internal",
		"app",
		r.domain.AppConfig.AppName(),
		"repositories",
		"postgres",
		r.domain.DirName(),
		r.domain.FileName(),
	)
}

func (r RepositoryGenerator) Sync() error {
	err := os.MkdirAll(path.Dir(r.filename()), 0777)
	if err != nil {
		return err
	}
	if err := r.syncStruct(); err != nil {
		return err
	}
	if err := r.syncConstructor(); err != nil {
		return err
	}
	if err := r.syncOrderByMap(); err != nil {
		return err
	}
	if err := r.syncEncodeOrderBy(); err != nil {
		return err
	}
	if err := r.syncDTOStruct(); err != nil {
		return err
	}
	if err := r.syncDTOListType(); err != nil {
		return err
	}
	if err := r.syncDTOListToEntities(); err != nil {
		return err
	}
	if err := r.syncDTOConstructor(); err != nil {
		return err
	}
	if err := r.syncDTOToModel(); err != nil {
		return err
	}
	if err := r.syncCreateMethod(); err != nil {
		return err
	}
	if err := r.syncGetMethod(); err != nil {
		return err
	}
	if err := r.syncListMethod(); err != nil {
		return err
	}
	if err := r.syncCountMethod(); err != nil {
		return err
	}
	if err := r.syncUpdateMethod(); err != nil {
		return err
	}
	if err := r.syncDeleteMethod(); err != nil {
		return err
	}
	if err := r.syncMigrations(); err != nil {
		return err
	}
	return nil
}

func (r RepositoryGenerator) dtoStruct() *ast.TypeSpec {
	structure := &ast.TypeSpec{
		Name: ast.NewIdent(r.getDTOName()),
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
							Value: "`db:\"id,omitempty\"`",
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
							Value: "`db:\"updated_at,omitempty\"`",
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
							Value: "`db:\"created_at,omitempty\"`",
						},
					},
				},
			},
		},
	}
	for _, param := range r.domain.GetMainModel().Params {
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
					Type:  ast.NewIdent(param.PostgresDTOType()),
					Tag: &ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf("`db:\"%s\"`", param.Tag()),
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

func (r RepositoryGenerator) syncDTOStruct() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, r.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	structure, structureExists := astfile.FindType(file, r.getDTOName())
	if structure == nil {
		structure = r.dtoStruct()
	}
	for _, param := range r.domain.GetMainModel().Params {
		astfile.SetTypeParam(
			structure,
			param.GetName(),
			param.PostgresDTOType(),
			fmt.Sprintf("`db:\"%s\"`", param.Tag()),
		)
	}
	if !structureExists {
		file.Decls = append(file.Decls, &ast.GenDecl{
			Tok:   token.TYPE,
			Specs: []ast.Spec{structure},
		})
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(r.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (r RepositoryGenerator) dtoConstructor() *ast.FuncDecl {
	dto := &ast.CompositeLit{
		Type: ast.NewIdent(r.getDTOName()),
		Elts: []ast.Expr{},
	}
	for _, param := range r.domain.GetMainModel().Params {
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
		dto.Elts = append(dto.Elts, elt)
	}
	constructor := &ast.FuncDecl{
		Name: ast.NewIdent(fmt.Sprintf("New%sFromEntity", r.getDTOName())),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("entity"),
						},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("entities"),
							Sel: ast.NewIdent(r.domain.GetMainModel().Name),
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent(r.getDTOName()),
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
	for _, param := range r.domain.GetMainModel().Params {
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
			},
		},
	)
	return constructor
}

func (r RepositoryGenerator) syncDTOConstructor() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, r.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var structureConstructorExists bool
	var structureConstructor *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok &&
			t.Name.String() == fmt.Sprintf("New%sFromEntity", r.getDTOName()) {
			structureConstructorExists = true
			structureConstructor = t
			return false
		}
		return true
	})
	if structureConstructor == nil {
		structureConstructor = r.dtoConstructor()
	}
	for _, param := range r.domain.GetMainModel().Params {
		param := param
		ast.Inspect(structureConstructor, func(node ast.Node) bool {
			if cl, ok := node.(*ast.CompositeLit); ok {
				if i, ok := cl.Type.(*ast.Ident); ok &&
					i.String() == r.getDTOName() {
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
	if err := os.WriteFile(r.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (r RepositoryGenerator) dtoToModel() *ast.FuncDecl {
	model := &ast.CompositeLit{
		Type: &ast.SelectorExpr{
			X:   ast.NewIdent("entities"),
			Sel: ast.NewIdent(r.domain.GetMainModel().Name),
		},
		Elts: []ast.Expr{},
	}
	for _, param := range r.domain.GetMainModel().Params {
		par := &ast.KeyValueExpr{
			Key: ast.NewIdent(param.GetName()),
		}
		if param.IsSlice() {
			par.Value = &ast.CompositeLit{
				Type: &ast.ArrayType{
					Elt: ast.NewIdent(param.SliceType()),
				},
			}
		} else {
			if param.PostgresDTOType() == param.Type {
				par.Value = &ast.SelectorExpr{
					X:   ast.NewIdent("dto"),
					Sel: ast.NewIdent(param.GetName()),
				}
			} else {
				paramType := param.Type
				if paramType == "UUID" {
					paramType = "uuid.UUID"
				}
				par.Value = &ast.CallExpr{
					Fun: ast.NewIdent(paramType),
					Args: []ast.Expr{
						&ast.SelectorExpr{
							X:   ast.NewIdent("dto"),
							Sel: ast.NewIdent(param.GetName()),
						},
					},
				}
			}
		}
		model.Elts = append(model.Elts, par)
	}
	method := &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("dto"),
					},
					Type: ast.NewIdent(r.getDTOName()),
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
							Sel: ast.NewIdent(r.domain.GetMainModel().Name),
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("entity"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						model,
					},
				},
			},
		},
	}
	for _, param := range r.domain.GetMainModel().Params {
		if !param.IsSlice() {
			continue
		}
		var valueToAppend ast.Expr
		if param.SliceType() == param.PostgresDTOSliceType() {
			valueToAppend = ast.NewIdent("param")
		} else {
			valueToAppend = &ast.CallExpr{
				Fun: ast.NewIdent(param.SliceType()),
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
				X:   ast.NewIdent("dto"),
				Sel: ast.NewIdent(param.GetName()),
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.SelectorExpr{
								X:   ast.NewIdent("entity"),
								Sel: ast.NewIdent(param.GetName()),
							},
						},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: ast.NewIdent("append"),
								Args: []ast.Expr{
									&ast.SelectorExpr{
										X:   ast.NewIdent("entity"),
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
		method.Body.List = append(method.Body.List, rang)
	}
	method.Body.List = append(
		method.Body.List,
		&ast.ReturnStmt{
			Results: []ast.Expr{
				ast.NewIdent("entity"),
			},
		},
	)
	return method
}

func (r RepositoryGenerator) syncDTOToModel() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, r.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExists bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == "toEntity" {
			methodExists = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = r.dtoToModel()
	}
	// TODO: add range sync
	if !methodExists {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(r.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (r RepositoryGenerator) astStruct() *ast.TypeSpec {
	structure := &ast.TypeSpec{
		Doc:        nil,
		Name:       ast.NewIdent(r.domain.GetRepositoryTypeName()),
		TypeParams: nil,
		Assign:     0,
		Type: &ast.StructType{
			Struct: 0,
			Fields: &ast.FieldList{
				Opening: 0,
				List: []*ast.Field{
					{
						Doc:     nil,
						Names:   []*ast.Ident{ast.NewIdent("readDB")},
						Type:    ast.NewIdent("database"),
						Tag:     nil,
						Comment: nil,
					},
					{
						Doc:     nil,
						Names:   []*ast.Ident{ast.NewIdent("writeDB")},
						Type:    ast.NewIdent("database"),
						Tag:     nil,
						Comment: nil,
					},
					{
						Doc:     nil,
						Names:   []*ast.Ident{ast.NewIdent("logger")},
						Type:    ast.NewIdent("logger"),
						Tag:     nil,
						Comment: nil,
					},
				},
			},
		},
		Comment: nil,
	}
	return structure
}

func (r RepositoryGenerator) file() *ast.File {
	specs := []ast.Spec{
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
				Value: r.domain.AppConfig.ProjectConfig.ErrsImportPath(),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: r.domain.EntitiesImportPath(),
			},
		},

		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: r.domain.AppConfig.ProjectConfig.PointerImportPath(),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: r.domain.AppConfig.ProjectConfig.UUIDImportPath(),
			},
		},
		&ast.ImportSpec{
			Name: ast.NewIdent("sq"),
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: `"github.com/Masterminds/squirrel"`,
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: `"github.com/lib/pq"`,
			},
		},
	}
	if r.domain.SearchEnabled() {
		specs = append(specs, &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: r.domain.AppConfig.ProjectConfig.PostgresImportPath(),
			},
		})
	}
	return &ast.File{
		Name: ast.NewIdent("repositories"),
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok:   token.IMPORT,
				Specs: specs,
			},
		},
	}
}

func (r RepositoryGenerator) syncStruct() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, r.filename(), nil, parser.ParseComments)
	if err != nil {
		file = r.file()
	}
	structure, structureExists := astfile.FindType(file, r.domain.GetRepositoryTypeName())
	if structure == nil {
		structure = r.astStruct()
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
	if err := os.WriteFile(r.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (r RepositoryGenerator) astConstructor() *ast.FuncDecl {
	constructor := &ast.FuncDecl{
		Doc:  nil,
		Recv: nil,
		Name: ast.NewIdent(fmt.Sprintf("New%s", r.domain.GetRepositoryTypeName())),
		Type: &ast.FuncType{
			Func:       0,
			TypeParams: nil,
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Doc:     nil,
						Names:   []*ast.Ident{ast.NewIdent("readDB")},
						Type:    ast.NewIdent("database"),
						Tag:     nil,
						Comment: nil,
					},
					{
						Doc:     nil,
						Names:   []*ast.Ident{ast.NewIdent("writeDB")},
						Type:    ast.NewIdent("database"),
						Tag:     nil,
						Comment: nil,
					},
					{
						Doc:     nil,
						Names:   []*ast.Ident{ast.NewIdent("logger")},
						Type:    ast.NewIdent("logger"),
						Tag:     nil,
						Comment: nil,
					},
				},
			},
			Results: &ast.FieldList{
				Opening: 0,
				List: []*ast.Field{
					{
						Doc:     nil,
						Names:   nil,
						Type:    ast.NewIdent(fmt.Sprintf("*%s", r.domain.GetRepositoryTypeName())),
						Tag:     nil,
						Comment: nil,
					},
				},
				Closing: 0,
			},
		},
		Body: &ast.BlockStmt{
			Lbrace: 0,
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Return: 0,
					Results: []ast.Expr{
						&ast.UnaryExpr{
							OpPos: 0,
							Op:    token.AND,
							X: &ast.CompositeLit{
								Type: ast.NewIdent(r.domain.GetRepositoryTypeName()),
								Elts: []ast.Expr{
									&ast.KeyValueExpr{
										Key:   ast.NewIdent("readDB"),
										Value: ast.NewIdent("readDB"),
									},
									&ast.KeyValueExpr{
										Key:   ast.NewIdent("writeDB"),
										Value: ast.NewIdent("writeDB"),
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
			Rbrace: 0,
		},
	}
	return constructor
}

func (r RepositoryGenerator) syncConstructor() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, r.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	method, methodExist := astfile.FindFunc(file, r.domain.GetRepositoryConstructorName())
	if method == nil {
		method = r.astConstructor()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(r.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (r RepositoryGenerator) astCreateMethod() *ast.FuncDecl {
	var columns []ast.Expr
	var values []ast.Expr
	for _, param := range r.domain.GetMainModel().Params {
		columns = append(columns, &ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf(`"%s"`, param.Tag()),
		})
		values = append(values, &ast.SelectorExpr{
			X:   ast.NewIdent("dto"),
			Sel: ast.NewIdent(param.GetName()),
		})
	}
	fun := &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("r"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(r.domain.GetRepositoryTypeName()),
					},
				},
			},
		},
		Name: ast.NewIdent("Create"),
		Type: &ast.FuncType{
			Func:       0,
			TypeParams: nil,
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("ctx")},
						Type:  ast.NewIdent("context.Context"),
					},
					{
						Names: []*ast.Ident{ast.NewIdent("tx")},
						Type:  &ast.SelectorExpr{X: ast.NewIdent("dtx"), Sel: ast.NewIdent("TX")},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("entity")},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("entities"),
							Sel: ast.NewIdent(r.domain.GetMainModel().Name),
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				// Setup timeout
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("ctx"),
						ast.NewIdent("cancel"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("context"),
								Sel: ast.NewIdent("WithTimeout"),
							},
							Args: []ast.Expr{
								ast.NewIdent("ctx"),
								&ast.SelectorExpr{
									X:   ast.NewIdent("time"),
									Sel: ast.NewIdent("Second"),
								},
							},
						},
					},
				},
				// Defer cancel
				&ast.DeferStmt{
					Call: &ast.CallExpr{
						Fun: ast.NewIdent("cancel"),
					},
				},
				// Create DTO from model
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("dto"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: ast.NewIdent(
								fmt.Sprintf("New%sDTOFromEntity", r.domain.GetMainModel().Name),
							),
							Args: []ast.Expr{
								ast.NewIdent("entity"),
							},
						},
					},
				},
				// Create sq
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("q"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("sq"),
												Sel: ast.NewIdent("Insert"),
											},
											Args: []ast.Expr{
												&ast.BasicLit{
													Kind: token.STRING,
													Value: fmt.Sprintf(
														`"public.%s"`,
														r.domain.TableName(),
													),
												},
											},
										},
										Sel: ast.NewIdent("Columns"),
									},
									Args: columns,
								},
								Sel: ast.NewIdent("Values"),
							},
							Args: values,
						},
					},
				},
				// Build query from sq
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("query"),
						ast.NewIdent("args"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("q"),
										Sel: ast.NewIdent("PlaceholderFormat"),
									},
									Args: []ast.Expr{
										&ast.SelectorExpr{
											X:   ast.NewIdent("sq"),
											Sel: ast.NewIdent("Dollar"),
										},
									},
								},
								Sel: ast.NewIdent("MustSql"),
							},
						},
					},
				},
				// Run query at DB
				&ast.IfStmt{
					Init: &ast.AssignStmt{
						Lhs: []ast.Expr{
							ast.NewIdent("_"),
							ast.NewIdent("err"),
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   ast.NewIdent("tx"),
											Sel: ast.NewIdent("GetSQLTx"),
										},
									},
									Sel: ast.NewIdent("ExecContext"),
								},
								Args: []ast.Expr{
									ast.NewIdent("ctx"),
									ast.NewIdent("query"),
									ast.NewIdent("args"),
								},
								Ellipsis: 653,
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
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									ast.NewIdent("e"),
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   ast.NewIdent("errs"),
											Sel: ast.NewIdent("FromPostgresError"),
										},
										Args: []ast.Expr{
											ast.NewIdent("err"),
										},
									},
								},
							},
							&ast.ReturnStmt{
								Results: []ast.Expr{
									ast.NewIdent("e"),
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
	return fun
}

func (r RepositoryGenerator) syncCreateMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, r.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	method, methodExist := astfile.FindFunc(file, "Create")
	if method == nil {
		method = r.astCreateMethod()
	}
	for _, param := range r.domain.GetMainModel().Params {
		param := param
		if param.GetName() == "ID" {
			continue
		}
		ast.Inspect(method, func(node ast.Node) bool {
			if call, ok := node.(*ast.CallExpr); ok {
				if fun, ok := call.Fun.(*ast.SelectorExpr); ok && fun.Sel.String() == "Columns" {
					for _, arg := range call.Args {
						arg := arg
						if bl, ok := arg.(*ast.BasicLit); ok &&
							bl.Value == fmt.Sprintf(`"%s"`, param.Tag()) {
							return false
						}
					}
					call.Args = append(call.Args, &ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf(`"%s"`, param.Tag()),
					})
					return false
				}
			}
			return true
		})
		ast.Inspect(method, func(node ast.Node) bool {
			if call, ok := node.(*ast.CallExpr); ok {
				if fun, ok := call.Fun.(*ast.SelectorExpr); ok && fun.Sel.String() == "Values" {
					for _, arg := range call.Args {
						arg := arg
						if bl, ok := arg.(*ast.SelectorExpr); ok &&
							bl.Sel.String() == param.GetName() {
							return false
						}
					}
					call.Args = append(call.Args, &ast.SelectorExpr{
						X:   ast.NewIdent("dto"),
						Sel: ast.NewIdent(param.GetName()),
					})
					return false
				}
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
	if err := os.WriteFile(r.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (r RepositoryGenerator) search() ast.Stmt {
	if !r.domain.SearchEnabled() {
		return &ast.EmptyStmt{}
	}
	var columns []ast.Expr
	for _, param := range r.domain.GetMainModel().Params {
		if param.Search {
			columns = append(columns, &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s"`, param.Tag()),
			})
		}
	}
	stmt := &ast.IfStmt{
		Cond: &ast.BinaryExpr{
			X: &ast.SelectorExpr{
				X:   ast.NewIdent("filter"),
				Sel: ast.NewIdent("Search"),
			},
			Op: token.NEQ,
			Y:  ast.NewIdent("nil"),
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("q"),
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("q"),
								Sel: ast.NewIdent("Where"),
							},
							Args: []ast.Expr{
								&ast.CompositeLit{
									Type: &ast.SelectorExpr{
										X:   ast.NewIdent("postgres"),
										Sel: ast.NewIdent("Search"),
									},
									Elts: []ast.Expr{
										&ast.KeyValueExpr{
											Key: ast.NewIdent("Lang"),
											Value: &ast.BasicLit{
												Kind:  token.STRING,
												Value: `"english"`,
											},
										},
										&ast.KeyValueExpr{
											Key: ast.NewIdent("Query"),
											Value: &ast.StarExpr{
												X: &ast.SelectorExpr{
													X:   ast.NewIdent("filter"),
													Sel: ast.NewIdent("Search"),
												},
											},
										},
										&ast.KeyValueExpr{
											Key: ast.NewIdent("Fields"),
											Value: &ast.CompositeLit{
												Type: &ast.ArrayType{
													Elt: ast.NewIdent("string"),
												},
												Elts: columns,
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
	return stmt
}

func (r RepositoryGenerator) listMethod() *ast.FuncDecl {
	tableName := r.domain.TableName()
	var columns []ast.Expr
	for _, param := range r.domain.GetMainModel().Params {
		columns = append(
			columns,
			&ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s.%s"`, tableName, param.Tag()),
			},
		)
	}
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("r"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(r.domain.GetRepositoryTypeName()),
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
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("entities"),
							Sel: ast.NewIdent(r.domain.GetFilterModel().Name),
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.ArrayType{
							Elt: &ast.SelectorExpr{
								X:   ast.NewIdent("entities"),
								Sel: ast.NewIdent(r.domain.GetMainModel().Name),
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
						ast.NewIdent("ctx"),
						ast.NewIdent("cancel"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("context"),
								Sel: ast.NewIdent("WithTimeout"),
							},
							Args: []ast.Expr{
								ast.NewIdent("ctx"),
								&ast.SelectorExpr{
									X:   ast.NewIdent("time"),
									Sel: ast.NewIdent("Second"),
								},
							},
						},
					},
				},
				&ast.DeferStmt{
					Call: &ast.CallExpr{
						Fun: ast.NewIdent("cancel"),
					},
				},
				&ast.DeclStmt{
					Decl: &ast.GenDecl{
						Tok: token.VAR,
						Specs: []ast.Spec{
							&ast.ValueSpec{
								Names: []*ast.Ident{
									ast.NewIdent("dto"),
								},
								Type: ast.NewIdent(r.getDTOListName()),
							},
						},
					},
				},
				&ast.DeclStmt{
					Decl: &ast.GenDecl{
						Tok: token.CONST,
						Specs: []ast.Spec{
							&ast.ValueSpec{
								Names: []*ast.Ident{
									ast.NewIdent("pageSize"),
								},
								Values: []ast.Expr{
									&ast.CallExpr{
										Fun: ast.NewIdent("uint64"),
										Args: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.INT,
												Value: "10",
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
						X: &ast.SelectorExpr{
							X:   ast.NewIdent("filter"),
							Sel: ast.NewIdent("PageSize"),
						},
						Op: token.EQL,
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
											ast.NewIdent("pageSize"),
										},
									},
								},
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("q"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("sq"),
												Sel: ast.NewIdent("Select"),
											},
											Args: columns,
										},
										Sel: ast.NewIdent("From"),
									},
									Args: []ast.Expr{
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: fmt.Sprintf(`"public.%s"`, tableName),
										},
									},
								},
								Sel: ast.NewIdent("Limit"),
							},
							Args: []ast.Expr{
								ast.NewIdent("pageSize"),
							},
						},
					},
				},
				r.search(),
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X: &ast.BinaryExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("filter"),
								Sel: ast.NewIdent("PageNumber"),
							},
							Op: token.NEQ,
							Y:  ast.NewIdent("nil"),
						},
						Op: token.LAND,
						Y: &ast.BinaryExpr{
							X: &ast.StarExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("filter"),
									Sel: ast.NewIdent("PageNumber"),
								},
							},
							Op: token.GTR,
							Y: &ast.BasicLit{
								Kind:  token.INT,
								Value: "1",
							},
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									ast.NewIdent("q"),
								},
								Tok: token.ASSIGN,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   ast.NewIdent("q"),
											Sel: ast.NewIdent("Offset"),
										},
										Args: []ast.Expr{
											&ast.BinaryExpr{
												X: &ast.ParenExpr{
													X: &ast.BinaryExpr{
														X: &ast.StarExpr{
															X: &ast.SelectorExpr{
																X:   ast.NewIdent("filter"),
																Sel: ast.NewIdent("PageNumber"),
															},
														},
														Op: token.SUB,
														Y: &ast.BasicLit{
															Kind:  token.INT,
															Value: "1",
														},
													},
												},
												Op: token.MUL,
												Y: &ast.StarExpr{
													X: &ast.SelectorExpr{
														X:   ast.NewIdent("filter"),
														Sel: ast.NewIdent("PageSize"),
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
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("q"),
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("q"),
								Sel: ast.NewIdent("Limit"),
							},
							Args: []ast.Expr{
								&ast.StarExpr{
									X: &ast.SelectorExpr{
										X:   ast.NewIdent("filter"),
										Sel: ast.NewIdent("PageSize"),
									},
								},
							},
						},
					},
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X: &ast.CallExpr{
							Fun: ast.NewIdent("len"),
							Args: []ast.Expr{
								&ast.SelectorExpr{
									X:   ast.NewIdent("filter"),
									Sel: ast.NewIdent("OrderBy"),
								},
							},
						},
						Op: token.GTR,
						Y: &ast.BasicLit{
							Kind:  token.INT,
							Value: "0",
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									ast.NewIdent("q"),
								},
								Tok: token.ASSIGN,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   ast.NewIdent("q"),
											Sel: ast.NewIdent("OrderBy"),
										},
										Args: []ast.Expr{
											&ast.CallExpr{
												Fun: &ast.Ident{
													Name: "encodeOrderBy",
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
												},
											},
										},
										Ellipsis: 5337,
									},
								},
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("query"),
						ast.NewIdent("args"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("q"),
										Sel: ast.NewIdent("PlaceholderFormat"),
									},
									Args: []ast.Expr{
										&ast.SelectorExpr{
											X:   ast.NewIdent("sq"),
											Sel: ast.NewIdent("Dollar"),
										},
									},
								},
								Sel: ast.NewIdent("MustSql"),
							},
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
									X: &ast.SelectorExpr{
										X:   ast.NewIdent("r"),
										Sel: ast.NewIdent("readDB"),
									},
									Sel: ast.NewIdent("SelectContext"),
								},
								Args: []ast.Expr{
									ast.NewIdent("ctx"),
									&ast.UnaryExpr{
										Op: token.AND,
										X:  ast.NewIdent("dto"),
									},
									ast.NewIdent("query"),
									ast.NewIdent("args"),
								},
								Ellipsis: 5460,
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
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									ast.NewIdent("e"),
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   ast.NewIdent("errs"),
											Sel: ast.NewIdent("FromPostgresError"),
										},
										Args: []ast.Expr{
											ast.NewIdent("err"),
										},
									},
								},
							},
							&ast.ReturnStmt{
								Results: []ast.Expr{
									ast.NewIdent("nil"),
									ast.NewIdent("e"),
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("dto"),
								Sel: ast.NewIdent("toEntities"),
							},
						},
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
}

func (r RepositoryGenerator) syncListMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, r.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	method, methodExist := astfile.FindFunc(file, "List")
	if method == nil {
		method = r.listMethod()
	}
	for _, param := range r.domain.GetMainModel().Params {
		param := param
		column := fmt.Sprintf(`"%s.%s"`, r.domain.TableName(), param.Tag())
		ast.Inspect(method, func(node ast.Node) bool {
			if call, ok := node.(*ast.CallExpr); ok {
				if fun, ok := call.Fun.(*ast.SelectorExpr); ok && fun.Sel.String() == "Select" {
					for _, arg := range call.Args {
						arg := arg
						if bl, ok := arg.(*ast.BasicLit); ok && bl.Value == column {
							return false
						}
					}
					call.Args = append(call.Args, &ast.BasicLit{
						Kind:  token.STRING,
						Value: column,
					})
					return false
				}
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
	if err := os.WriteFile(r.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (r RepositoryGenerator) astCountMethod() *ast.FuncDecl {
	tableName := r.domain.TableName()
	columns := []ast.Expr{
		&ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf(`"%s.id"`, tableName),
		},
		&ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf(`"%s.updated_at"`, tableName),
		},
		&ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf(`"%s.created_at"`, tableName),
		},
	}
	for _, param := range r.domain.GetMainModel().Params {
		columns = append(
			columns,
			&ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s.%s"`, tableName, param.Tag()),
			},
		)
	}
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("r"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(r.domain.GetRepositoryTypeName()),
					},
				},
			},
		},
		Name: ast.NewIdent("Count"),
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
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("entities"),
							Sel: ast.NewIdent(r.domain.GetFilterModel().Name),
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent("uint64"),
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
						ast.NewIdent("ctx"),
						ast.NewIdent("cancel"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("context"),
								Sel: ast.NewIdent("WithTimeout"),
							},
							Args: []ast.Expr{
								ast.NewIdent("ctx"),
								&ast.SelectorExpr{
									X:   ast.NewIdent("time"),
									Sel: ast.NewIdent("Second"),
								},
							},
						},
					},
				},
				&ast.DeferStmt{
					Call: &ast.CallExpr{
						Fun: ast.NewIdent("cancel"),
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("q"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("sq"),
										Sel: ast.NewIdent("Select"),
									},
									Args: []ast.Expr{
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: `"count(id)"`,
										},
									},
								},
								Sel: ast.NewIdent("From"),
							},
							Args: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: fmt.Sprintf(`"public.%s"`, r.domain.TableName()),
								},
							},
						},
					},
				},
				r.search(),
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("query"),
						ast.NewIdent("args"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("q"),
										Sel: ast.NewIdent("PlaceholderFormat"),
									},
									Args: []ast.Expr{
										&ast.SelectorExpr{
											X:   ast.NewIdent("sq"),
											Sel: ast.NewIdent("Dollar"),
										},
									},
								},
								Sel: ast.NewIdent("MustSql"),
							},
						},
					},
				},
				&ast.DeclStmt{
					Decl: &ast.GenDecl{
						Tok: token.VAR,
						Specs: []ast.Spec{
							&ast.ValueSpec{
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
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "r",
										},
										Sel: &ast.Ident{
											Name: "readDB",
										},
									},
									Sel: &ast.Ident{
										Name: "GetContext",
									},
								},
								Args: []ast.Expr{
									&ast.Ident{
										Name: "ctx",
									},
									&ast.UnaryExpr{
										Op: token.AND,
										X: &ast.Ident{
											Name: "count",
										},
									},
									&ast.Ident{
										Name: "query",
									},
									&ast.Ident{
										Name: "args",
									},
								},
								Ellipsis: 378,
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
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.Ident{
										Name: "e",
									},
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "errs",
											},
											Sel: &ast.Ident{
												Name: "FromPostgresError",
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
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.BasicLit{
										Kind:  token.INT,
										Value: "0",
									},
									&ast.Ident{
										Name: "e",
									},
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						ast.NewIdent("count"),
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
}

func (r RepositoryGenerator) syncCountMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, r.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	method, methodExist := astfile.FindFunc(file, "Count")
	if method == nil {
		method = r.astCountMethod()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(r.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (r RepositoryGenerator) getMethod() *ast.FuncDecl {
	tableName := r.domain.TableName()
	var columns []ast.Expr
	for _, param := range r.domain.GetMainModel().Params {
		columns = append(
			columns,
			&ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s.%s"`, tableName, param.Tag()),
			},
		)
	}
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("r"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(r.domain.GetRepositoryTypeName()),
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
							ast.NewIdent("id"),
						},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("uuid"),
							Sel: ast.NewIdent("UUID"),
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("entities"),
							Sel: ast.NewIdent(r.domain.GetMainModel().Name),
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
						ast.NewIdent("ctx"),
						ast.NewIdent("cancel"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("context"),
								Sel: ast.NewIdent("WithTimeout"),
							},
							Args: []ast.Expr{
								ast.NewIdent("ctx"),
								&ast.SelectorExpr{
									X:   ast.NewIdent("time"),
									Sel: ast.NewIdent("Second"),
								},
							},
						},
					},
				},
				&ast.DeferStmt{
					Call: &ast.CallExpr{
						Fun: ast.NewIdent("cancel"),
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("dto"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X: &ast.CompositeLit{
								Type: ast.NewIdent(r.getDTOName()),
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("q"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.CallExpr{
													Fun: &ast.SelectorExpr{
														X:   ast.NewIdent("sq"),
														Sel: ast.NewIdent("Select"),
													},
													Args: columns,
												},
												Sel: ast.NewIdent("From"),
											},
											Args: []ast.Expr{
												&ast.BasicLit{
													Kind:  token.STRING,
													Value: fmt.Sprintf(`"public.%s"`, tableName),
												},
											},
										},
										Sel: ast.NewIdent("Where"),
									},
									Args: []ast.Expr{
										&ast.CompositeLit{
											Type: &ast.SelectorExpr{
												X:   ast.NewIdent("sq"),
												Sel: ast.NewIdent("Eq"),
											},
											Elts: []ast.Expr{
												&ast.KeyValueExpr{
													Key: &ast.BasicLit{
														Kind:  token.STRING,
														Value: `"id"`,
													},
													Value: ast.NewIdent("id"),
												},
											},
										},
									},
								},
								Sel: ast.NewIdent("Limit"),
							},
							Args: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.INT,
									Value: "1",
								},
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("query"),
						ast.NewIdent("args"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("q"),
										Sel: ast.NewIdent("PlaceholderFormat"),
									},
									Args: []ast.Expr{
										&ast.SelectorExpr{
											X:   ast.NewIdent("sq"),
											Sel: ast.NewIdent("Dollar"),
										},
									},
								},
								Sel: ast.NewIdent("MustSql"),
							},
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
									X: &ast.SelectorExpr{
										X:   ast.NewIdent("r"),
										Sel: ast.NewIdent("readDB"),
									},
									Sel: ast.NewIdent("GetContext"),
								},
								Args: []ast.Expr{
									ast.NewIdent("ctx"),
									ast.NewIdent("dto"),
									ast.NewIdent("query"),
									ast.NewIdent("args"),
								},
								Ellipsis: 4211,
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
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									ast.NewIdent("e"),
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X:   ast.NewIdent("errs"),
													Sel: ast.NewIdent("FromPostgresError"),
												},
												Args: []ast.Expr{
													ast.NewIdent("err"),
												},
											},
											Sel: ast.NewIdent("WithParam"),
										},
										Args: []ast.Expr{
											&ast.BasicLit{
												Kind: token.STRING,
												Value: fmt.Sprintf(
													`"%s_id"`,
													strcase.ToSnake(r.domain.GetMainModel().Name),
												),
											},
											&ast.CallExpr{
												Fun: ast.NewIdent("id.String"),
											},
										},
									},
								},
							},
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.CompositeLit{
										Type: &ast.SelectorExpr{
											X:   ast.NewIdent("entities"),
											Sel: ast.NewIdent(r.domain.GetMainModel().Name),
										},
									},
									ast.NewIdent("e"),
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("dto"),
								Sel: ast.NewIdent("toEntity"),
							},
						},
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
}

func (r RepositoryGenerator) syncGetMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, r.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	method, methodExist := astfile.FindFunc(file, "Get")
	if method == nil {
		method = r.getMethod()
	}
	for _, param := range r.domain.GetMainModel().Params {
		param := param
		column := fmt.Sprintf(`"%s.%s"`, r.domain.TableName(), param.Tag())
		ast.Inspect(method, func(node ast.Node) bool {
			if call, ok := node.(*ast.CallExpr); ok {
				if fun, ok := call.Fun.(*ast.SelectorExpr); ok && fun.Sel.String() == "Select" {
					for _, arg := range call.Args {
						arg := arg
						if bl, ok := arg.(*ast.BasicLit); ok && bl.Value == column {
							return false
						}
					}
					call.Args = append(call.Args, &ast.BasicLit{
						Kind:  token.STRING,
						Value: column,
					})
					return false
				}
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
	if err := os.WriteFile(r.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (r RepositoryGenerator) updateMethod() *ast.FuncDecl {
	tableName := r.domain.TableName()
	updateBlock := &ast.BlockStmt{
		List: []ast.Stmt{},
	}
	for _, param := range r.domain.GetMainModel().Params {
		if param.GetName() == "ID" {
			continue
		}
		updateBlock.List = append(updateBlock.List, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent("q"),
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("q"),
						Sel: ast.NewIdent("Set"),
					},
					Args: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s"`, param.Tag()),
						},
						&ast.SelectorExpr{
							X:   ast.NewIdent("dto"),
							Sel: ast.NewIdent(param.GetName()),
						},
					},
				},
			},
		})
	}
	method := &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("r"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(r.domain.GetRepositoryTypeName()),
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
						Names: []*ast.Ident{ast.NewIdent("tx")},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("dtx"),
							Sel: ast.NewIdent("TX"),
						},
					},
					{
						Names: []*ast.Ident{
							ast.NewIdent("entity"),
						},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("entities"),
							Sel: ast.NewIdent(r.domain.GetMainModel().Name),
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
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
						ast.NewIdent("ctx"),
						ast.NewIdent("cancel"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("context"),
								Sel: ast.NewIdent("WithTimeout"),
							},
							Args: []ast.Expr{
								ast.NewIdent("ctx"),
								&ast.SelectorExpr{
									X:   ast.NewIdent("time"),
									Sel: ast.NewIdent("Second"),
								},
							},
						},
					},
				},
				&ast.DeferStmt{
					Call: &ast.CallExpr{
						Fun: ast.NewIdent("cancel"),
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("dto"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: ast.NewIdent(fmt.Sprintf("New%sFromEntity", r.getDTOName())),
							Args: []ast.Expr{
								ast.NewIdent("entity"),
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("q"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("sq"),
										Sel: ast.NewIdent("Update"),
									},
									Args: []ast.Expr{
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: fmt.Sprintf(`"public.%s"`, tableName),
										},
									},
								},
								Sel: ast.NewIdent("Where"),
							},
							Args: []ast.Expr{
								&ast.CompositeLit{
									Type: &ast.SelectorExpr{
										X:   ast.NewIdent("sq"),
										Sel: ast.NewIdent("Eq"),
									},
									Elts: []ast.Expr{
										&ast.KeyValueExpr{
											Key: &ast.BasicLit{
												Kind:  token.STRING,
												Value: `"id"`,
											},
											Value: &ast.SelectorExpr{
												X:   ast.NewIdent("entity"),
												Sel: ast.NewIdent("ID"),
											},
										},
									},
								},
							},
						},
					},
				},
				updateBlock,
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("query"),
						ast.NewIdent("args"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("q"),
										Sel: ast.NewIdent("PlaceholderFormat"),
									},
									Args: []ast.Expr{
										&ast.SelectorExpr{
											X:   ast.NewIdent("sq"),
											Sel: ast.NewIdent("Dollar"),
										},
									},
								},
								Sel: ast.NewIdent("MustSql"),
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("result"),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("tx"),
										Sel: ast.NewIdent("GetSQLTx"),
									},
								},
								Sel: ast.NewIdent("ExecContext"),
							},
							Args: []ast.Expr{
								ast.NewIdent("ctx"),
								ast.NewIdent("query"),
								ast.NewIdent("args"),
							},
							Ellipsis: 6334,
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
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									ast.NewIdent("e"),
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X:   ast.NewIdent("errs"),
													Sel: ast.NewIdent("FromPostgresError"),
												},
												Args: []ast.Expr{
													ast.NewIdent("err"),
												},
											},
											Sel: ast.NewIdent("WithParam"),
										},
										Args: []ast.Expr{
											&ast.BasicLit{
												Kind: token.STRING,
												Value: fmt.Sprintf(
													`"%s_id"`,
													strcase.ToSnake(r.domain.GetMainModel().Name),
												),
											},
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X:   ast.NewIdent("fmt"),
													Sel: ast.NewIdent("Sprint"),
												},
												Args: []ast.Expr{
													&ast.SelectorExpr{
														X:   ast.NewIdent("entity"),
														Sel: ast.NewIdent("ID"),
													},
												},
											},
										},
									},
								},
							},
							&ast.ReturnStmt{
								Results: []ast.Expr{
									ast.NewIdent("e"),
								},
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("affected"),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("result"),
								Sel: ast.NewIdent("RowsAffected"),
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
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X:   ast.NewIdent("errs"),
													Sel: ast.NewIdent("FromPostgresError"),
												},
												Args: []ast.Expr{
													ast.NewIdent("err"),
												},
											},
											Sel: ast.NewIdent("WithParam"),
										},
										Args: []ast.Expr{
											&ast.BasicLit{
												Kind: token.STRING,
												Value: fmt.Sprintf(
													`"%s_id"`,
													strcase.ToSnake(r.domain.GetMainModel().Name),
												),
											},
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X:   ast.NewIdent("fmt"),
													Sel: ast.NewIdent("Sprint"),
												},
												Args: []ast.Expr{
													&ast.SelectorExpr{
														X:   ast.NewIdent("entity"),
														Sel: ast.NewIdent("ID"),
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
						X:  ast.NewIdent("affected"),
						Op: token.EQL,
						Y: &ast.BasicLit{
							Kind:  token.INT,
							Value: "0",
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									ast.NewIdent("e"),
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X:   ast.NewIdent("errs"),
													Sel: ast.NewIdent("NewEntityNotFoundError"),
												},
											},
											Sel: ast.NewIdent("WithParam"),
										},
										Args: []ast.Expr{
											&ast.BasicLit{
												Kind: token.STRING,
												Value: fmt.Sprintf(
													`"%s_id"`,
													strcase.ToSnake(r.domain.GetMainModel().Name),
												),
											},
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X:   ast.NewIdent("fmt"),
													Sel: ast.NewIdent("Sprint"),
												},
												Args: []ast.Expr{
													&ast.SelectorExpr{
														X:   ast.NewIdent("entity"),
														Sel: ast.NewIdent("ID"),
													},
												},
											},
										},
									},
								},
							},
							&ast.ReturnStmt{
								Results: []ast.Expr{
									ast.NewIdent("e"),
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
	return method
}

func (r RepositoryGenerator) syncUpdateMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, r.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	method, methodExist := astfile.FindFunc(file, "Update")
	if method == nil {
		method = r.updateMethod()
	}
	for _, param := range r.domain.GetMainModel().Params {
		param := param
		if param.GetName() == "ID" {
			continue
		}
		exists := false
		for _, stmt := range method.Body.List {
			if update, ok := stmt.(*ast.BlockStmt); ok {
				for _, updateStmt := range update.List {
					ast.Inspect(updateStmt, func(node ast.Node) bool {
						if call, ok := node.(*ast.CallExpr); ok {
							if callSelector, ok := call.Fun.(*ast.SelectorExpr); ok &&
								callSelector.Sel.String() == "Set" {
								for _, arg := range call.Args {
									if bl, ok := arg.(*ast.BasicLit); ok &&
										bl.Value == fmt.Sprintf(
											`"%s"`,
											param.Tag(),
										) {
										exists = true
										return false
									}
								}
							}
						}
						return true
					})
				}
				if !exists {
					update.List = append(
						update.List,
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("q"),
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("q"),
										Sel: ast.NewIdent("Set"),
									},
									Args: []ast.Expr{
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: fmt.Sprintf(`"%s"`, param.Tag()),
										},
										&ast.SelectorExpr{
											X:   ast.NewIdent("dto"),
											Sel: ast.NewIdent(param.GetName()),
										},
									},
								},
							},
						},
					)
				}
			}
		}
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(r.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (r RepositoryGenerator) astDeleteMethod() *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("r"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(r.domain.GetRepositoryTypeName()),
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
						Names: []*ast.Ident{ast.NewIdent("tx")},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("dtx"),
							Sel: ast.NewIdent("TX"),
						},
					},
					{
						Names: []*ast.Ident{
							ast.NewIdent("id"),
						},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("uuid"),
							Sel: ast.NewIdent("UUID"),
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
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
						ast.NewIdent("ctx"),
						ast.NewIdent("cancel"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("context"),
								Sel: ast.NewIdent("WithTimeout"),
							},
							Args: []ast.Expr{
								ast.NewIdent("ctx"),
								&ast.SelectorExpr{
									X:   ast.NewIdent("time"),
									Sel: ast.NewIdent("Second"),
								},
							},
						},
					},
				},
				&ast.DeferStmt{
					Call: &ast.CallExpr{
						Fun: ast.NewIdent("cancel"),
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("q"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("sq"),
										Sel: ast.NewIdent("Delete"),
									},
									Args: []ast.Expr{
										&ast.BasicLit{
											Kind: token.STRING,
											Value: fmt.Sprintf(
												`"public.%s"`,
												r.domain.TableName(),
											),
										},
									},
								},
								Sel: ast.NewIdent("Where"),
							},
							Args: []ast.Expr{
								&ast.CompositeLit{
									Type: &ast.SelectorExpr{
										X:   ast.NewIdent("sq"),
										Sel: ast.NewIdent("Eq"),
									},
									Elts: []ast.Expr{
										&ast.KeyValueExpr{
											Key: &ast.BasicLit{
												Kind:  token.STRING,
												Value: `"id"`,
											},
											Value: ast.NewIdent("id"),
										},
									},
								},
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("query"),
						ast.NewIdent("args"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("q"),
										Sel: ast.NewIdent("PlaceholderFormat"),
									},
									Args: []ast.Expr{
										&ast.SelectorExpr{
											X:   ast.NewIdent("sq"),
											Sel: ast.NewIdent("Dollar"),
										},
									},
								},
								Sel: ast.NewIdent("MustSql"),
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("result"),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("tx"),
										Sel: ast.NewIdent("GetSQLTx"),
									},
								},
								Sel: ast.NewIdent("ExecContext"),
							},
							Args: []ast.Expr{
								ast.NewIdent("ctx"),
								ast.NewIdent("query"),
								ast.NewIdent("args"),
							},
							Ellipsis: 7041,
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
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									ast.NewIdent("e"),
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X:   ast.NewIdent("errs"),
													Sel: ast.NewIdent("FromPostgresError"),
												},
												Args: []ast.Expr{
													ast.NewIdent("err"),
												},
											},
											Sel: ast.NewIdent("WithParam"),
										},
										Args: []ast.Expr{
											&ast.BasicLit{
												Kind: token.STRING,
												Value: fmt.Sprintf(
													`"%s_id"`,
													strcase.ToSnake(r.domain.GetMainModel().Name),
												),
											},
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X:   ast.NewIdent("fmt"),
													Sel: ast.NewIdent("Sprint"),
												},
												Args: []ast.Expr{
													ast.NewIdent("id"),
												},
											},
										},
									},
								},
							},
							&ast.ReturnStmt{
								Results: []ast.Expr{
									ast.NewIdent("e"),
								},
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("affected"),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("result"),
								Sel: ast.NewIdent("RowsAffected"),
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
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									ast.NewIdent("e"),
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X:   ast.NewIdent("errs"),
													Sel: ast.NewIdent("FromPostgresError"),
												},
												Args: []ast.Expr{
													ast.NewIdent("err"),
												},
											},
											Sel: ast.NewIdent("WithParam"),
										},
										Args: []ast.Expr{
											&ast.BasicLit{
												Kind: token.STRING,
												Value: fmt.Sprintf(
													`"%s_id"`,
													strcase.ToSnake(r.domain.GetMainModel().Name),
												),
											},
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X:   ast.NewIdent("fmt"),
													Sel: ast.NewIdent("Sprint"),
												},
												Args: []ast.Expr{
													ast.NewIdent("id"),
												},
											},
										},
									},
								},
							},
							&ast.ReturnStmt{
								Results: []ast.Expr{
									ast.NewIdent("e"),
								},
							},
						},
					},
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X:  ast.NewIdent("affected"),
						Op: token.EQL,
						Y: &ast.BasicLit{
							Kind:  token.INT,
							Value: "0",
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									ast.NewIdent("e"),
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X:   ast.NewIdent("errs"),
													Sel: ast.NewIdent("NewEntityNotFoundError"),
												},
											},
											Sel: ast.NewIdent("WithParam"),
										},
										Args: []ast.Expr{
											&ast.BasicLit{
												Kind: token.STRING,
												Value: fmt.Sprintf(
													`"%s_id"`,
													strcase.ToSnake(r.domain.GetMainModel().Name),
												),
											},
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X:   ast.NewIdent("fmt"),
													Sel: ast.NewIdent("Sprint"),
												},
												Args: []ast.Expr{
													ast.NewIdent("id"),
												},
											},
										},
									},
								},
							},
							&ast.ReturnStmt{
								Results: []ast.Expr{
									ast.NewIdent("e"),
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
}

func (r RepositoryGenerator) syncDeleteMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, r.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	method, methodExist := astfile.FindFunc(file, "Delete")
	if method == nil {
		method = r.astDeleteMethod()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(r.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (r RepositoryGenerator) astDTOListType() *ast.TypeSpec {
	return &ast.TypeSpec{
		Name: ast.NewIdent(r.getDTOListName()),
		Type: &ast.ArrayType{
			Elt: ast.NewIdent(r.getDTOName()),
		},
	}
}

func (r RepositoryGenerator) syncDTOListType() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, r.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	dtoListType, structureExists := astfile.FindType(file, r.getDTOListName())
	if dtoListType == nil {
		dtoListType = r.astDTOListType()
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
	if err := os.WriteFile(r.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (r RepositoryGenerator) astDTOToEntities() *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("list"),
					},
					Type: ast.NewIdent(r.getDTOListName()),
				},
			},
		},
		Name: ast.NewIdent("toEntities"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.ArrayType{
							Elt: &ast.SelectorExpr{
								X:   ast.NewIdent("entities"),
								Sel: ast.NewIdent(r.domain.GetMainModel().Name),
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
						ast.NewIdent("items"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: ast.NewIdent("make"),
							Args: []ast.Expr{
								&ast.ArrayType{
									Elt: &ast.SelectorExpr{
										X:   ast.NewIdent("entities"),
										Sel: ast.NewIdent(r.domain.GetMainModel().Name),
									},
								},
								&ast.CallExpr{
									Fun: ast.NewIdent("len"),
									Args: []ast.Expr{
										ast.NewIdent("list"),
									},
								},
							},
						},
					},
				},
				&ast.RangeStmt{
					Key: ast.NewIdent("i"),
					Tok: token.DEFINE,
					X:   ast.NewIdent("list"),
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.IndexExpr{
										X:     ast.NewIdent("items"),
										Index: ast.NewIdent("i"),
									},
								},
								Tok: token.ASSIGN,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.IndexExpr{
												X:     ast.NewIdent("list"),
												Index: ast.NewIdent("i"),
											},
											Sel: ast.NewIdent("toEntity"),
										},
									},
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						ast.NewIdent("items"),
					},
				},
			},
		},
	}
}

func (r RepositoryGenerator) astOrderByMap() *ast.GenDecl {
	var values []ast.Expr
	for cnt, column := range r.domain.OrderingMap() {
		values = append(values, &ast.KeyValueExpr{
			Key: &ast.SelectorExpr{
				X: &ast.Ident{
					Name: "entities",
				},
				Sel: &ast.Ident{
					Name: cnt,
				},
			},
			Value: &ast.BasicLit{
				Kind:  token.STRING,
				Value: column,
			},
		})
	}
	return &ast.GenDecl{
		Tok: token.VAR,
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Names: []*ast.Ident{
					{
						Name: "orderByMap",
					},
				},
				Values: []ast.Expr{
					&ast.CompositeLit{
						Type: &ast.MapType{
							Key: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "entities",
								},
								Sel: &ast.Ident{
									Name: r.domain.OrderingTypeName(),
								},
							},
							Value: &ast.Ident{
								Name: "string",
							},
						},
						Elts: values,
					},
				},
			},
		},
	}
}

func (r RepositoryGenerator) syncOrderByMap() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, r.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	varExists := astfile.VarExists(file, "orderByMap")
	if !varExists {
		file.Decls = append(file.Decls, r.astOrderByMap())
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(r.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (r RepositoryGenerator) syncDTOListToEntities() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, r.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	method, methodExist := astfile.FindFunc(file, "toEntities")
	if method == nil {
		method = r.astDTOToEntities()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(r.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (r RepositoryGenerator) astEncodeOrderBy() *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: &ast.Ident{
			Name: "encodeOrderBy",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
								Name: "orderBy",
							},
						},
						Type: &ast.ArrayType{
							Elt: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "entities",
								},
								Sel: &ast.Ident{
									Name: r.domain.OrderingTypeName(),
								},
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.ArrayType{
							Elt: &ast.Ident{
								Name: "string",
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
							Name: "columns",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.Ident{
								Name: "make",
							},
							Args: []ast.Expr{
								&ast.ArrayType{
									Elt: &ast.Ident{
										Name: "string",
									},
								},
								&ast.CallExpr{
									Fun: &ast.Ident{
										Name: "len",
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
				&ast.RangeStmt{
					Key: &ast.Ident{
						Name: "i",
					},
					Value: &ast.Ident{
						Name: "item",
					},
					Tok: token.DEFINE,
					X: &ast.Ident{
						Name: "orderBy",
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.Ident{
										Name: "column",
									},
									&ast.Ident{
										Name: "exists",
									},
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.IndexExpr{
										X: &ast.Ident{
											Name: "orderByMap",
										},
										Index: &ast.Ident{
											Name: "item",
										},
									},
								},
							},
							&ast.IfStmt{
								Cond: &ast.UnaryExpr{
									Op: token.NOT,
									X: &ast.Ident{
										Name: "exists",
									},
								},
								Body: &ast.BlockStmt{
									List: []ast.Stmt{
										&ast.BranchStmt{
											Tok: token.CONTINUE,
										},
									},
								},
							},
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.IndexExpr{
										X: &ast.Ident{
											Name: "columns",
										},
										Index: &ast.Ident{
											Name: "i",
										},
									},
								},
								Tok: token.ASSIGN,
								Rhs: []ast.Expr{
									&ast.Ident{
										Name: "column",
									},
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.Ident{
							Name: "columns",
						},
					},
				},
			},
		},
	}
}

func (r RepositoryGenerator) syncEncodeOrderBy() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, r.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	method, methodExist := astfile.FindFunc(file, "encodeOrderBy")
	if method == nil {
		method = r.astEncodeOrderBy()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(r.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

var destinationPath = "."

func (r RepositoryGenerator) syncMigrations() error {
	pattern := fmt.Sprintf("*_%s.up.sql", r.domain.TableName())
	dir, err := os.ReadDir(path.Join(
		destinationPath,
		"internal",
		"pkg",
		"postgres",
		"migrations",
	))
	if err != nil {
		return err
	}
	for _, file := range dir {
		match, err := filepath.Match(pattern, file.Name())
		if err != nil {
			return err
		}
		if match {
			return nil
		}
	}

	files := []*tmpl.Template{
		{
			SourcePath: "templates/internal/pkg/postgres/migrations/crud.up.sql.tmpl",
			DestinationPath: path.Join(
				destinationPath,
				"internal",
				"pkg",
				"postgres",
				"migrations",
				r.domain.MigrationUpFileName(),
			),
			Name: "migration up",
		},
		{
			SourcePath: "templates/internal/pkg/postgres/migrations/crud.down.sql.tmpl",
			DestinationPath: path.Join(
				destinationPath,
				"internal",
				"pkg",
				"postgres",
				"migrations",
				r.domain.MigrationDownFileName(),
			),
			Name: "migration down",
		},
	}
	for _, file := range files {
		if err := file.RenderToFile(r.domain); err != nil {
			return err
		}
	}
	return nil
}
