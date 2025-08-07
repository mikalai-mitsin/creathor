package astfile

import (
	"go/ast"
	"go/token"
)

func TypeExists(file ast.Node, typeName string) bool {
	var loggerExists bool
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok {
			if t.Name.String() == typeName {
				loggerExists = true
			}
			return true
		}
		return true
	})
	return loggerExists
}

func FindType(file ast.Node, typeName string) (*ast.TypeSpec, bool) {
	var structure *ast.TypeSpec
	var structureExists bool
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok && t.Name.String() == typeName {
			structure = t
			structureExists = true
			return false
		}
		return true
	})
	return structure, structureExists
}

func SetTypeParam(typeSpec *ast.TypeSpec, name, typeName, tag string) {
	ast.Inspect(typeSpec, func(node ast.Node) bool {
		if st, ok := node.(*ast.StructType); ok && st.Fields != nil {
			for _, field := range st.Fields.List {
				for _, fieldName := range field.Names {
					if fieldName.Name == name {
						return false
					}
				}
			}
			field := &ast.Field{
				Doc:     nil,
				Names:   []*ast.Ident{ast.NewIdent(name)},
				Type:    ast.NewIdent(typeName),
				Tag:     nil,
				Comment: nil,
			}
			if tag != "" {
				field.Tag = &ast.BasicLit{
					Kind:  token.STRING,
					Value: tag,
				}
			}
			st.Fields.List = append(st.Fields.List, field)
			return false
		}
		return true
	})
}
