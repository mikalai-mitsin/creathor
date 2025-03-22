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

type ModelPermission struct {
	project *configs.Project
}

func NewModelPermission(project *configs.Project) *ModelPermission {
	return &ModelPermission{project: project}
}

func (i ModelPermission) file() *ast.File {
	return &ast.File{
		Name: ast.NewIdent("entities"),
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: ast.NewIdent("PermissionID"),
						Type: ast.NewIdent("string"),
					},
				},
			},
			&ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								ast.NewIdent("p"),
							},
							Type: ast.NewIdent("PermissionID"),
						},
					},
				},
				Name: ast.NewIdent("String"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: ast.NewIdent("string"),
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.CallExpr{
									Fun: ast.NewIdent("string"),
									Args: []ast.Expr{
										ast.NewIdent("p"),
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
						Name: ast.NewIdent("Permission"),
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											ast.NewIdent("ID"),
										},
										Type: ast.NewIdent("PermissionID"),
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`db:\"id,omitempty\" json:\"id\"   form:\"id\"`",
										},
									},
									{
										Names: []*ast.Ident{
											ast.NewIdent("Name"),
										},
										Type: ast.NewIdent("string"),
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
						Name: ast.NewIdent("GroupID"),
						Type: ast.NewIdent("string"),
					},
				},
			},
			&ast.GenDecl{
				Tok: token.CONST,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("GroupIDAdmin"),
						},
						Type: ast.NewIdent("GroupID"),
						Values: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: "\"admin\"",
							},
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("GroupIDUser"),
						},
						Type: ast.NewIdent("GroupID"),
						Values: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: "\"user\"",
							},
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("GroupIDGuest"),
						},
						Type: ast.NewIdent("GroupID"),
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
						Name: ast.NewIdent("Group"),
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											ast.NewIdent("ID"),
										},
										Type: ast.NewIdent("GroupID"),
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`db:\"id,omitempty\" json:\"id\"`",
										},
									},
									{
										Names: []*ast.Ident{
											ast.NewIdent("Name"),
										},
										Type: ast.NewIdent("string"),
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
	filename := path.Join("internal", "app", "user", "entities", "permission.go")
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
