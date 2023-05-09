package mods

import "golang.org/x/exp/slices"

type Layer struct {
	Auth     bool
	Events   bool
	Name     string
	Variable string
	Methods  []*Method
}

func (i *Layer) GetMethod(t MethodType) *Method {
	index := slices.IndexFunc(i.Methods, func(method *Method) bool { return method.Type == t })
	if index >= 0 {
		return i.Methods[index]
	}
	return nil
}

func (i *Layer) GetCreateMethod() *Method {
	return i.GetMethod(MethodTypeCreate)
}
func (i *Layer) GetUpdateMethod() *Method {
	return i.GetMethod(MethodTypeUpdate)
}
func (i *Layer) GetDeleteMethod() *Method {
	return i.GetMethod(MethodTypeDelete)
}
func (i *Layer) GetListMethod() *Method {
	return i.GetMethod(MethodTypeList)
}

func (i *Layer) GetGetMethod() *Method {
	return i.GetMethod(MethodTypeGet)
}
