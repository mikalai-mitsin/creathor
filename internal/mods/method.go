package mods

import "go/ast"

type MethodType uint8

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
	Type   MethodType
	Args   []*ast.Field
	Return []*ast.Field
}
