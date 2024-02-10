package mods

import "go/ast"

// deprecated
type MethodType uint8

// deprecated
const (
	MethodTypeGet = iota
	MethodTypeList
	MethodTypeCount
	MethodTypeCreate
	MethodTypeUpdate
	MethodTypeDelete
)

type Method struct {
	Name   string
	Type   MethodType // deprecated
	Args   []*ast.Field
	Return []*ast.Field
}
