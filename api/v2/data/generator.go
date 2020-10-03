package data

var Generators = make([]Generator, 0)

type Generator interface {
	Generate(sourceDir string, destinationDir string) ([]*Resource, error)
	IsManagedDir(dir string) bool
}
