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
	"path/filepath"
)

type Structure struct {
	fileName string
	name     string
	domain   *mods.Domain
	params   []*ast.Field
}

func NewStructure(fileName string, name string, params []*ast.Field, domain *mods.Domain) *Structure {
	return &Structure{
		fileName: fileName,
		name:     name,
		domain:   domain,
		params:   params,
	}
}

func (m *Structure) file() *ast.File {
	return &ast.File{
		Name: ast.NewIdent("models"),
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
						Name: ast.NewIdent("validation"),
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"github.com/go-ozzo/ozzo-validation/v4"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"github.com/go-ozzo/ozzo-validation/v4/is"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/pkg/uuid"`, m.domain.Module),
						},
					},
				},
			},
		},
	}
}

func (m *Structure) filename() string {
	return filepath.Join("internal", m.domain.Name, "models", m.fileName)
}

func (m *Structure) fill(structure *ast.TypeSpec) {
	for _, param := range m.params {
		ast.Inspect(structure, func(node ast.Node) bool {
			if st, ok := node.(*ast.StructType); ok && st.Fields != nil {
				for _, field := range st.Fields.List {
					for _, fieldName := range field.Names {
						for _, name := range param.Names {
							if fieldName.Name == name.String() {
								return false
							}
						}

					}
				}
				st.Fields.List = append(st.Fields.List, param)
				return true
			}
			return true
		})
	}
}

func (m *Structure) spec() *ast.TypeSpec {
	return &ast.TypeSpec{
		Name: ast.NewIdent(m.name),
		Type: &ast.StructType{
			Fields: &ast.FieldList{
				List: m.params,
			},
		},
	}
}

func (m *Structure) Sync() error {
	fileset := token.NewFileSet()
	filename := m.filename()
	err := os.MkdirAll(path.Dir(filename), 0777)
	if err != nil {
		return err
	}
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = m.file()
	}
	var structure *ast.TypeSpec
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok && t.Name.String() == m.name {
			structure = t
			return false
		}
		return true
	})
	if structure == nil {
		structure = m.spec()
		gd := &ast.GenDecl{
			Tok:   token.TYPE,
			Specs: []ast.Spec{structure},
		}
		file.Decls = append(file.Decls, gd)
	}
	m.fill(structure)
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(filename, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}
