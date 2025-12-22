package astfile

import (
	"go/ast"
	"go/token"
	"slices"
)

func TypeExists(file ast.Node, typeName string) bool {
	found := false
	ast.Inspect(file, func(node ast.Node) bool {
		ts, ok := node.(*ast.TypeSpec)
		if !ok {
			return true
		}
		if ts.Name.Name == typeName {
			found = true
			return false
		}
		return true
	})
	return found
}

func FindType(file ast.Node, typeName string) (*ast.TypeSpec, bool) {
	var ts *ast.TypeSpec
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok && t.Name.Name == typeName {
			ts = t
			return false
		}
		return true
	})
	return ts, ts != nil
}

func SetTypeParam(typeSpec *ast.TypeSpec, name, typeName, tag string) {
	st, ok := typeSpec.Type.(*ast.StructType)
	if !ok || st.Fields == nil {
		return
	}

	for _, field := range st.Fields.List {
		for _, fieldName := range field.Names {
			if fieldName.Name == name {
				return
			}
		}
	}
	field := &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(name)},
		Type:  ast.NewIdent(typeName),
	}
	if tag != "" {
		field.Tag = &ast.BasicLit{
			Kind:  token.STRING,
			Value: tag,
		}
	}

	st.Fields.List = append(st.Fields.List, field)
}

func SetTypeField(typeSpec *ast.TypeSpec, field *ast.Field) {
	st, ok := typeSpec.Type.(*ast.StructType)
	if !ok || st.Fields == nil {
		return
	}
	if slices.ContainsFunc(st.Fields.List, func(stField *ast.Field) bool {
		return slices.ContainsFunc(stField.Names, func(stName *ast.Ident) bool {
			return slices.ContainsFunc(field.Names, func(fName *ast.Ident) bool {
				return stName.Name == fName.Name
			})
		})
	}) {
		return
	}
	st.Fields.List = append(st.Fields.List, field)
}

func FindFunc(file ast.Node, name string) (*ast.FuncDecl, bool) {
	var fn *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.Name == name {
			fn = t
			return false
		}
		return true
	})
	return fn, fn != nil
}

func FindMethod(file ast.Node, receiver, name string) (*ast.FuncDecl, bool) {
	var function *ast.FuncDecl

	ast.Inspect(file, func(node ast.Node) bool {
		fn, ok := node.(*ast.FuncDecl)
		if !ok || fn.Name.Name != name || fn.Recv == nil {
			return true
		}
		for _, field := range fn.Recv.List {
			switch t := field.Type.(type) {
			case *ast.Ident:
				if t.Name == receiver {
					function = fn
					return false
				}
			case *ast.StarExpr:
				if ident, ok := t.X.(*ast.Ident); ok && ident.Name == receiver {
					function = fn
					return false
				}
			}
		}
		return true
	})

	return function, function != nil
}

func ConstExists(file ast.Node, name string) bool {
	found := false
	ast.Inspect(file, func(node ast.Node) bool {
		decl, ok := node.(*ast.GenDecl)
		if !ok || decl.Tok != token.CONST {
			return true
		}
		for _, spec := range decl.Specs {
			if vs, ok := spec.(*ast.ValueSpec); ok {
				for _, ident := range vs.Names {
					if ident.Name == name {
						found = true
						return false
					}
				}
			}
		}
		return true
	})
	return found
}

func VarExists(file ast.Node, name string) bool {
	found := false
	ast.Inspect(file, func(node ast.Node) bool {
		decl, ok := node.(*ast.GenDecl)
		if !ok || decl.Tok != token.VAR {
			return true
		}
		for _, spec := range decl.Specs {
			if vs, ok := spec.(*ast.ValueSpec); ok {
				for _, ident := range vs.Names {
					if ident.Name == name {
						found = true
						return false
					}
				}
			}
		}
		return true
	})
	return found
}

func SetParamValue(cl *ast.CompositeLit, parameter *ast.KeyValueExpr) {
	prKey := parameter.Key.(*ast.Ident)
	for _, elt := range cl.Elts {
		if kv, ok := elt.(*ast.KeyValueExpr); ok {
			if ident, ok := kv.Key.(*ast.Ident); ok && ident.String() == prKey.String() {
				return
			}
		}
	}
	cl.Elts = append(cl.Elts, parameter)
}

func FindStructInstance(file ast.Node, typeName string) (*ast.CompositeLit, bool) {
	var compositeLit *ast.CompositeLit
	ast.Inspect(file, func(node ast.Node) bool {
		if cl, ok := node.(*ast.CompositeLit); ok {
			if clType, ok := cl.Type.(*ast.SelectorExpr); ok {
				if clType.Sel.String() == typeName {
					compositeLit = cl
					return false
				}
			}
		}
		return true
	})
	return compositeLit, compositeLit != nil
}

func AppendToFuncBody(method *ast.FuncDecl, val ast.Stmt) {
	if len(method.Body.List) < 1 {
		method.Body.List = append(method.Body.List, val)
	}
	idx := len(method.Body.List) - 1
	method.Body.List = append(method.Body.List, val)
	copy(method.Body.List[idx+1:], method.Body.List[idx:])
	method.Body.List[idx] = val
}
