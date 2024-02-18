package mods

import "go/ast"

type Method struct {
	Name   string
	Args   []*ast.Field
	Return []*ast.Field
}
