package entities

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

type ModelEvent struct {
	project *configs.Project
}

func NewModelEvent(project *configs.Project) *ModelEvent {
	return &ModelEvent{project: project}
}

func (m ModelEvent) file() *ast.File {
	return &ast.File{
		Name: ast.NewIdent("entities"),
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: ast.NewIdent("EventOperation"),
						Type: ast.NewIdent("string"),
					},
				},
			},
			&ast.GenDecl{
				Tok: token.CONST,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("EventTypeCreated"),
						},
						Type: ast.NewIdent("EventOperation"),
						Values: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: `"created"`,
							},
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("EventTypeUpdated"),
						},
						Type: ast.NewIdent("EventOperation"),
						Values: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: `"updated"`,
							},
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("EventTypeDeleted"),
						},
						Type: ast.NewIdent("EventOperation"),
						Values: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: `"deleted"`,
							},
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: ast.NewIdent("Event"),
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											ast.NewIdent("Operation"),
										},
										Type: ast.NewIdent("EventOperation"),
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"operation\"`",
										},
									},
									{
										Names: []*ast.Ident{
											ast.NewIdent("Message"),
										},
										Type: ast.NewIdent("string"),
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"message\"`",
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

func (m ModelEvent) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "domain", "entities", "event.go")
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = m.file()
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
