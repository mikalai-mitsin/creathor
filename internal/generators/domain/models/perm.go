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

	"golang.org/x/exp/slices"

	"github.com/iancoleman/strcase"
)

type Perm struct {
	modelName string
	fileName  string
	domain    *mods.Domain
}

func NewPerm(modelName string, fileName string, domain *mods.Domain) *Perm {
	return &Perm{modelName: modelName, fileName: fileName, domain: domain}
}

func (m *Perm) file() *ast.File {
	file := &ast.File{
		Name: ast.NewIdent("models"),
		Decls: []ast.Decl{
			m.perms(),
		},
	}
	return file
}

func (m *Perm) perms() *ast.GenDecl {
	return &ast.GenDecl{
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{
					Text: "// Model permissions.",
				},
			},
		},
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
	return filepath.Join("internal", "app", "user", "models", fmt.Sprintf("permission_%s", m.fileName))
}

func (m *Perm) Sync() error {
	fileset := token.NewFileSet()
	if err := os.MkdirAll(path.Dir(m.filename()), 0777); err != nil {
		return err
	}
	file, err := parser.ParseFile(fileset, m.filename(), nil, parser.ParseComments)
	if err != nil {
		file = m.file()
	}
	m.fill(file)
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(m.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (m *Perm) fill(file *ast.File) {
	var perms *ast.GenDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if genDecl, ok := node.(*ast.GenDecl); ok && genDecl.Doc != nil {
			contains := slices.ContainsFunc(genDecl.Doc.List, func(comment *ast.Comment) bool {
				return comment.Text == "// Model permissions."
			})
			if contains {
				perms = genDecl
				return false
			}
			return true
		}
		return true
	})
	if perms == nil {
		perms = m.perms()
		file.Decls = append(file.Decls, perms)
	}
}
