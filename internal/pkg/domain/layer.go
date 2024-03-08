package domain

type Layer struct {
	Auth     bool
	Events   bool
	Name     string
	Variable string
	Methods  []*Method
}
