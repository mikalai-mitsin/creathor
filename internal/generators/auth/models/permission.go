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

type ModelPermission struct {
	project *configs.Project
}

func NewModelPermission(project *configs.Project) *ModelPermission {
	return &ModelPermission{project: project}
}

func (i ModelPermission) file() *ast.File {
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
							Name: "PermissionID",
						},
						Type: &ast.Ident{
							Name: "string",
						},
					},
				},
			},
			&ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								{
									Name: "p",
								},
							},
							Type: &ast.Ident{
								Name: "PermissionID",
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "String",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.Ident{
									Name: "string",
								},
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.Ident{
										Name: "string",
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "p",
										},
									},
								},
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
							Name: "Permission",
						},
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											{
												Name: "ID",
											},
										},
										Type: &ast.Ident{
											Name: "PermissionID",
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`db:\"id,omitempty\" json:\"id\"   form:\"id\"`",
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "Name",
											},
										},
										Type: &ast.Ident{
											Name: "string",
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`db:\"name\"         json:\"name\" form:\"name\"`",
										},
									},
								},
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
							Name: "GroupID",
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
								Name: "GroupIDAdmin",
							},
						},
						Type: &ast.Ident{
							Name: "GroupID",
						},
						Values: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: "\"admin\"",
							},
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							{
								Name: "GroupIDUser",
							},
						},
						Type: &ast.Ident{
							Name: "GroupID",
						},
						Values: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: "\"user\"",
							},
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							{
								Name: "GroupIDGuest",
							},
						},
						Type: &ast.Ident{
							Name: "GroupID",
						},
						Values: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: "\"guest\"",
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
							Name: "Group",
						},
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											{
												Name: "ID",
											},
										},
										Type: &ast.Ident{
											Name: "GroupID",
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`db:\"id,omitempty\" json:\"id\"`",
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "Name",
											},
										},
										Type: &ast.Ident{
											Name: "string",
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`db:\"name\"         json:\"name\"`",
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

func (i ModelPermission) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "app", "user", "models", "permission.go")
	if err := os.MkdirAll(path.Dir(filename), 0777); err != nil {
		return err
	}
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
