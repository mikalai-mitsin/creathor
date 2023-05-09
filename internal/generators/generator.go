package generators

type Generator interface {
	Sync() error
}
