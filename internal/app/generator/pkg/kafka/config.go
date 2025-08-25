package kafka

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"

	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
)

type ConfigGenerator struct {
	project *configs.Project
}

func NewConfigGenerator(project *configs.Project) *ConfigGenerator {
	return &ConfigGenerator{project: project}
}

func (u ConfigGenerator) file() *ast.File {
	return &ast.File{
		Name: &ast.Ident{
			Name: "kafka",
		},
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{
							Name: "Config",
						},
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											{
												Name: "Brokers",
											},
										},
										Type: &ast.ArrayType{
											Elt: &ast.Ident{
												Name: "string",
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
}

func (u ConfigGenerator) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "pkg", "kafka", "config.go")
	if err := os.MkdirAll(path.Dir(filename), 0777); err != nil {
		return err
	}
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = u.file()
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
