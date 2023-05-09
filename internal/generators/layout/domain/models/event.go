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

type ModelEvent struct {
	project *configs.Project
}

func NewModelEvent(project *configs.Project) *ModelEvent {
	return &ModelEvent{project: project}
}

func (m ModelEvent) file() *ast.File {
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
							Name: "EventOperation",
						},
						Type: &ast.Ident{
							Name: "string",
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.CONST,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{
							{
								Name: "EventTypeCreated",
							},
						},
						Type: &ast.Ident{
							Name: "EventOperation",
						},
						Values: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: `"created"`,
							},
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							{
								Name: "EventTypeUpdated",
							},
						},
						Type: &ast.Ident{
							Name: "EventOperation",
						},
						Values: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: `"updated"`,
							},
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							{
								Name: "EventTypeDeleted",
							},
						},
						Type: &ast.Ident{
							Name: "EventOperation",
						},
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
						Name: &ast.Ident{
							Name: "Event",
						},
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											{
												Name: "Operation",
											},
										},
										Type: &ast.Ident{
											Name: "EventOperation",
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"operation\"`",
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "Message",
											},
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
	filename := path.Join("internal", "domain", "models", "event.go")
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
