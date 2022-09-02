package generator

//go:generate mockery -case=underscore -output=. -inpkg -name=Generator

type Generator interface {
	Generate(name string, path string) error
}
