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

func FindFunc(file ast.Node, name string) (*ast.FuncDecl, bool) {
	var funcExists bool
	var function *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok &&
			t.Name.String() == name {
			funcExists = true
			function = t
			return false
		}
		return true
	})
	return function, funcExists
}

func FindMethod(file ast.Node, receiver, name string) (*ast.FuncDecl, bool) {
	var funcExists bool
	var function *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == name {
			for _, ident := range t.Recv.List {
				if ti, ok := ident.Type.(*ast.Ident); ok && ti.String() == receiver {
					funcExists = true
					function = t
				}
			}
			return false
		}
		return true
	})
	return function, funcExists
}

func ConstExists(file ast.Node, name string) bool {
	var constExists bool
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.GenDecl); ok && t.Tok == token.CONST {
			for _, spec := range t.Specs {
				if s, ok := spec.(*ast.ValueSpec); ok && s.Names[0].Name == name {
					constExists = true
				}
			}
			return false
		}
		return true
	})
	return constExists
}

func VarExists(file ast.Node, name string) bool {
	var constExists bool
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.GenDecl); ok && t.Tok == token.VAR {
			for _, spec := range t.Specs {
				if s, ok := spec.(*ast.ValueSpec); ok && s.Names[0].Name == name {
					constExists = true
				}
			}
			return false
		}
		return true
	})
	return constExists
}
