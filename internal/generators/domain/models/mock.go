package models

import (
	"bytes"
	"fmt"
	mods "github.com/018bf/creathor/internal/domain"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"

	"github.com/018bf/creathor/internal/fake"
)

type Mock struct {
	typeSpec *ast.TypeSpec
	fileName string
	domain   *mods.Domain
}

func NewMock(typeSpec *ast.TypeSpec, fileName string, domain *mods.Domain) *Mock {
	return &Mock{typeSpec: typeSpec, fileName: fileName, domain: domain}
}

func (m *Mock) constructorName() string {
	return fmt.Sprintf("New%s", m.typeSpec.Name)
}

func (m *Mock) constructor() *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: ast.NewIdent(m.constructorName()),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("t")},
						Type: &ast.StarExpr{
							Star: 0,
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("testing"),
								Sel: ast.NewIdent("T"),
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
								X:   ast.NewIdent("models"),
								Sel: m.typeSpec.Name,
							},
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("t"),
							Sel: ast.NewIdent("Helper"),
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{m.model()},
				},
			},
		},
	}
}

func (m *Mock) model() *ast.UnaryExpr {
	cl := &ast.CompositeLit{
		Type: &ast.SelectorExpr{
			X:   ast.NewIdent("models"),
			Sel: m.typeSpec.Name,
		},
		Elts: []ast.Expr{},
	}
	for _, kv := range m.values() {
		cl.Elts = append(cl.Elts, kv)
	}
	return &ast.UnaryExpr{
		Op: token.AND,
		X:  cl,
	}
}

func (m *Mock) values() []*ast.KeyValueExpr {
	var kvs []*ast.KeyValueExpr
	if st, ok := m.typeSpec.Type.(*ast.StructType); ok && st.Fields != nil {
		for _, field := range st.Fields.List {
			for _, name := range field.Names {
				switch name.String() {
				case "Email":
					kvs = append(kvs, &ast.KeyValueExpr{Key: name, Value: fake.Email(field.Type)})
				default:
					kvs = append(kvs, &ast.KeyValueExpr{Key: name, Value: fake.Value(field.Type)})
				}
			}
		}
	}
	return kvs
}

func (m *Mock) fill(cl *ast.CompositeLit) {
	v := m.values()
	for _, kv := range v {
		var exists bool
		for _, elt := range cl.Elts {
			if e, ok := elt.(*ast.KeyValueExpr); ok {
				if fmt.Sprint(kv.Key) == fmt.Sprint(e.Key) {
					exists = true
				}
			}
		}
		if !exists {
			cl.Elts = append(cl.Elts, kv)
		}
	}
}

func (m *Mock) file() *ast.File {
	return &ast.File{
		Name: &ast.Ident{
			Name: "mock_models",
		},
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/pkg/utils"`, m.domain.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/pkg/uuid"`, m.domain.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: m.domain.ModelsImportPath(),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"github.com/jaswdr/faker"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"testing"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"time"`,
						},
					},
				},
			},
		},
	}
}

func (m *Mock) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "app", m.domain.DirName(), "models", "mock", m.fileName)
	err := os.MkdirAll(path.Dir(filename), 0777)
	if err != nil {
		return err
	}
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = m.file()
	}
	var mock *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if fun, ok := node.(*ast.FuncDecl); ok && fun.Name.Name == m.constructorName() {
			mock = fun
			return false
		}
		return true
	})
	if mock == nil {
		mock = m.constructor()
		file.Decls = append(file.Decls, mock)
	}
	ast.Inspect(mock.Body, func(node ast.Node) bool {
		if cl, ok := node.(*ast.CompositeLit); ok {
			if sel, ok := cl.Type.(*ast.SelectorExpr); ok {
				if sel.Sel.Name != m.typeSpec.Name.Name {
					return true
				}
			}
			m.fill(cl)
			return false
		}
		return true
	})
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(filename, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}
