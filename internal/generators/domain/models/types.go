package models

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"

	"github.com/018bf/creathor/internal/configs"
)

type ModelTypes struct {
	project *configs.Project
}

func NewModelTypes(project *configs.Project) *ModelTypes {
	return &ModelTypes{project: project}
}

func (i ModelTypes) file() *ast.File {
	return &ast.File{
		Name: &ast.Ident{
			Name: "models",
		},
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{
							Name: "UUID",
						},
						Type: &ast.Ident{
							Name: "string",
						},
					},
				},
			},
		},
	}
}

func (i ModelTypes) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "domain", "models", "types.go")
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = i.file()
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
