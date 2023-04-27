package usecases

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

type UseCaseInterfaceCrud struct {
	mod *configs.Mod
}

func NewUseCaseInterfaceCrud(mod *configs.Mod) *UseCaseInterfaceCrud {
	return &UseCaseInterfaceCrud{mod: mod}
}

func (i UseCaseInterfaceCrud) file() *ast.File {
	return &ast.File{
		Name: ast.NewIdent("usecases"),
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

func (i UseCaseInterfaceCrud) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "domain", "usecases", i.mod.Filename)
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = i.file()
	}
	var structureExists bool
	var structure *ast.TypeSpec
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok && t.Name.String() == i.mod.Name {
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
							"//%s - domain layer use case interface",
							i.mod.UseCase.Name,
						),
					},
					{
						Text: fmt.Sprintf(
							"//go:generate mockgen -build_flags=-mod=mod -destination mock/%s . %s",
							i.mod.Filename,
							i.mod.UseCase.Name,
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

func (i UseCaseInterfaceCrud) astInterface() *ast.TypeSpec {
	methods := make([]*ast.Field, len(i.mod.UseCase.Methods))
	for i, method := range i.mod.UseCase.Methods {
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
		Name: &ast.Ident{
			Name: i.mod.UseCase.Name,
		},
		Type: &ast.InterfaceType{
			Methods: &ast.FieldList{
				List: methods,
			},
		},
	}
}
