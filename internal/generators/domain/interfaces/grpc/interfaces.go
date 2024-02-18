package grpc

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"

	mods "github.com/018bf/creathor/internal/domain"
)

type InterceptorInterfaceCrud struct {
	mod *mods.Domain
}

func NewInterceptorInterfaceCrud(mod *mods.Domain) *InterceptorInterfaceCrud {
	return &InterceptorInterfaceCrud{mod: mod}
}

func (i InterceptorInterfaceCrud) file() *ast.File {
	return &ast.File{
		Name: ast.NewIdent("grpc"),
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
							Value: fmt.Sprintf(`"%s/internal/%s/models"`, i.mod.Module, i.mod.Name),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/pkg/uuid"`, i.mod.Module),
						},
					},
				},
			},
		},
	}
}

func (i InterceptorInterfaceCrud) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", i.mod.Name, "handlers", "grpc", "interfaces.go")
	err := os.MkdirAll(path.Dir(filename), 0777)
	if err != nil {
		return err
	}
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = i.file()
	}
	var structureExists bool
	var structure *ast.TypeSpec
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok && t.Name.String() == i.mod.Interceptor.Name {
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
							"//%s - domain layer interceptor interface",
							i.mod.Interceptor.Name,
						),
					},
					{
						Text: fmt.Sprintf(
							"//go:generate mockgen -build_flags=-mod=mod -destination mock/interfaces.go . %s",
							i.mod.Interceptor.Name,
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

func (i InterceptorInterfaceCrud) astInterface() *ast.TypeSpec {
	methods := make([]*ast.Field, len(i.mod.Interceptor.Methods))
	for i, method := range i.mod.Interceptor.Methods {
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
			Name: i.mod.Interceptor.Name,
		},
		Type: &ast.InterfaceType{
			Methods: &ast.FieldList{
				List: methods,
			},
		},
	}
}