package generator

type Generator interface {
	Sync() error
}
