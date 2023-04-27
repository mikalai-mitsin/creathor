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

	"github.com/018bf/creathor/internal/configs"
)

type RepositoryInterfaceCrud struct {
	mod *configs.Mod
}

func NewRepositoryInterfaceCrud(mod *configs.Mod) *RepositoryInterfaceCrud {
	return &RepositoryInterfaceCrud{mod: mod}
}

func (i RepositoryInterfaceCrud) file() *ast.File {
	return &ast.File{
		Name: ast.NewIdent("repositories"),
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
							Value: fmt.Sprintf(`"%s/internal/domain/models"`, i.mod.Module),
						},
					},
				},
			},
		},
	}
}

func (i RepositoryInterfaceCrud) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "domain", "repositories", i.mod.Filename)
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = i.file()
	}
	var structureExists bool
	var structure *ast.TypeSpec
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok && t.Name.String() == i.mod.Repository.Name {
			structure = t
			structureExists = true
			return false
		}
		return true
	})
	if structure == nil {
		structure = i.astInterface()
	}
	if !structureExists {
		gd := &ast.GenDecl{
			Doc: &ast.CommentGroup{
				List: []*ast.Comment{
					{
						Text: fmt.Sprintf(
							"//%s - domain layer repository interface",
							i.mod.Repository.Name,
						),
					},
					{
						Text: fmt.Sprintf(
							"//go:generate mockgen -build_flags=-mod=mod -destination mock/%s . %s",
							i.mod.Filename,
							i.mod.Repository.Name,
						),
					},
				},
			},
			Tok:   token.TYPE,
			Specs: []ast.Spec{structure},
		}
		file.Decls = append(file.Decls, gd)
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

func (i RepositoryInterfaceCrud) astInterface() *ast.TypeSpec {
	methods := make([]*ast.Field, len(i.mod.Repository.Methods))
	for i, method := range i.mod.Repository.Methods {
		methods[i] = &ast.Field{
			Names: []*ast.Ident{
				{
					Name: method.Name,
				},
			},
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: method.Args,
				},
				Results: &ast.FieldList{
					List: method.Return,
				},
			},
		}
	}
	return &ast.TypeSpec{
		Name: ast.NewIdent(i.mod.Repository.Name),
		Type: &ast.InterfaceType{
			Methods: &ast.FieldList{
				List: methods,
			},
		},
	}
}
