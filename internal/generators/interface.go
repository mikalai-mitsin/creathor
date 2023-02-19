package generators

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
)

type Method struct {
	Name    string
	Args    []*Param
	Results []*Param
}

type Interface struct {
	Path     string
	Name     string
	Comments []string
	Methods  []*Method
}

func (i Interface) SyncInterface() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, i.Path, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var structureExists bool
	var structure *ast.TypeSpec
	_ = structureExists
	_ = structure
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok && t.Name.String() == i.Name {
			structure = t
			structureExists = true
			return false
		}
		return true
	})
	if structure == nil {
		structure = &ast.TypeSpec{
			Doc:        nil,
			Name:       ast.NewIdent(i.Name),
			TypeParams: nil,
			Assign:     0,
			Type: &ast.InterfaceType{
				Interface:  0,
				Methods:    &ast.FieldList{},
				Incomplete: false,
			},
			Comment: nil,
		}
	}
	for _, method := range i.Methods {
		ast.Inspect(structure, func(node ast.Node) bool {
			if st, ok := node.(*ast.InterfaceType); ok && st.Methods != nil {
				for _, meth := range st.Methods.List {
					for _, fieldName := range meth.Names {
						if fieldName.Name == method.Name {
							return false
						}
					}
				}
				fn := &ast.FuncType{
					Func:       0,
					TypeParams: nil,
					Params: &ast.FieldList{
						Opening: 0,
						List:    []*ast.Field{},
					},
					Results: &ast.FieldList{
						Opening: 0,
						List:    nil,
						Closing: 0,
					},
				}
				for _, par := range method.Args {
					fn.Params.List = append(fn.Params.List, &ast.Field{
						Doc:     nil,
						Names:   []*ast.Ident{ast.NewIdent(par.Name)},
						Type:    ast.NewIdent(par.Type),
						Tag:     nil,
						Comment: nil,
					})
				}
				for _, res := range method.Results {
					fn.Results.List = append(fn.Results.List, &ast.Field{
						Doc:     nil,
						Names:   nil,
						Type:    ast.NewIdent(res.Type),
						Tag:     nil,
						Comment: nil,
					})
				}
				st.Methods.List = append(st.Methods.List, &ast.Field{
					Doc:     nil,
					Names:   []*ast.Ident{ast.NewIdent(method.Name)},
					Type:    fn,
					Tag:     nil,
					Comment: nil,
				})
				return true
			}
			return true
		})
	}
	if !structureExists {
		gd := &ast.GenDecl{
			Doc:    nil,
			TokPos: 0,
			Tok:    token.TYPE,
			Lparen: 0,
			Specs:  []ast.Spec{structure},
			Rparen: 0,
		}
		file.Decls = append(file.Decls, gd)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(i.Path, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}
