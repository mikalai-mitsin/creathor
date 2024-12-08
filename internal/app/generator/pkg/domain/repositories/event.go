package repositories

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"

	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
)

type RepositoryInterfaceEvent struct {
	project *configs.Project
}

func NewRepositoryInterfaceEvent(project *configs.Project) *RepositoryInterfaceEvent {
	return &RepositoryInterfaceEvent{project: project}
}

func (i RepositoryInterfaceEvent) astInterface() *ast.GenDecl {
	return &ast.GenDecl{
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{
					Text: "//EventRepository - domain layer repository interface",
				},
				{
					Text: "//go:generate mockgen -build_flags=-mod=mod -destination mock/interfaces.go . EventRepository",
				},
			},
		},
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: &ast.Ident{
					Name: "EventRepository",
				},
				Type: &ast.InterfaceType{
					Methods: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{
										Name: "Send",
									},
								},
								Type: &ast.FuncType{
									Params: &ast.FieldList{
										List: []*ast.Field{
											{
												Names: []*ast.Ident{
													{
														Name: "ctx",
													},
												},
												Type: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "context",
													},
													Sel: &ast.Ident{
														Name: "Context",
													},
												},
											},
											{
												Names: []*ast.Ident{
													{
														Name: "event",
													},
												},
												Type: &ast.StarExpr{
													X: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "entities",
														},
														Sel: &ast.Ident{
															Name: "Event",
														},
													},
												},
											},
										},
									},
									Results: &ast.FieldList{
										List: []*ast.Field{
											{
												Type: &ast.Ident{
													Name: "error",
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
		},
	}
}

func (i RepositoryInterfaceEvent) file() *ast.File {
	return &ast.File{
		Name: &ast.Ident{
			Name: "repositories",
		},
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"context"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/domain/entities"`, i.project.Module),
						},
					},
				},
			},
			i.astInterface(),
		},
	}
}

func (i RepositoryInterfaceEvent) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "domain", "repositories", "event.go")
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = i.file()
	}
	var exists bool
	ast.Inspect(file, func(node ast.Node) bool {
		if typeSpec, ok := node.(*ast.TypeSpec); ok && typeSpec.Name.String() == "EventRepository" {
			exists = true
			return false
		}
		return true
	})
	if !exists {
		file.Decls = append(file.Decls, i.astInterface())
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
