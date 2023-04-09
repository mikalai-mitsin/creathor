package models

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"

	"github.com/iancoleman/strcase"
)

type Perm struct {
	modelName string
	fileName  string
}

func NewPerm(modelName string, fileName string) *Perm {
	return &Perm{modelName: modelName, fileName: fileName}
}

func (m *Perm) file() *ast.File {
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
				},
			},
			m.perms(),
		},
	}
}

func (m *Perm) perms() *ast.GenDecl {
	return &ast.GenDecl{
		Tok: token.CONST,
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Names: []*ast.Ident{
					{
						Name: fmt.Sprintf("PermissionID%sList", strcase.ToCamel(m.modelName)),
					},
				},
				Type: &ast.Ident{
					Name: "PermissionID",
				},
				Values: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf(`"%s_list"`, strcase.ToSnake(m.modelName)),
					},
				},
			},
			&ast.ValueSpec{
				Names: []*ast.Ident{
					{
						Name: fmt.Sprintf("PermissionID%sDetail", strcase.ToCamel(m.modelName)),
					},
				},
				Type: &ast.Ident{
					Name: "PermissionID",
				},
				Values: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf(`"%s_detail"`, strcase.ToSnake(m.modelName)),
					},
				},
			},
			&ast.ValueSpec{
				Names: []*ast.Ident{
					{
						Name: fmt.Sprintf("PermissionID%sCreate", strcase.ToCamel(m.modelName)),
					},
				},
				Type: &ast.Ident{
					Name: "PermissionID",
				},
				Values: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf(`"%s_create"`, strcase.ToSnake(m.modelName)),
					},
				},
			},
			&ast.ValueSpec{
				Names: []*ast.Ident{
					{
						Name: fmt.Sprintf("PermissionID%sUpdate", strcase.ToCamel(m.modelName)),
					},
				},
				Type: &ast.Ident{
					Name: "PermissionID",
				},
				Values: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf(`"%s_update"`, strcase.ToSnake(m.modelName)),
					},
				},
			},
			&ast.ValueSpec{
				Names: []*ast.Ident{
					{
						Name: fmt.Sprintf("PermissionID%sDelete", strcase.ToCamel(m.modelName)),
					},
				},
				Type: &ast.Ident{
					Name: "PermissionID",
				},
				Values: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf(`"%s_delete"`, strcase.ToSnake(m.modelName)),
					},
				},
			},
		},
	}
}

func (m *Perm) filename() string {
	return filepath.Join("internal", "domain", "models", m.fileName)
}

func (m *Perm) Sync() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, m.filename(), nil, parser.ParseComments)
	if err != nil {
		file = m.file()
	}
	// TODO: add sync nolint:godoc
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(m.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}
