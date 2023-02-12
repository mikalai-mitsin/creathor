package models

type Interface struct {
	Path     string
	Name     string
	Comments []string
	Methods  []*Method
}
