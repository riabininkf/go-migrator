package generator

import (
	"fmt"
	"io/ioutil"
	pathLib "path"

	"github.com/riabininkf/go-migrator/pkg/scanner"
)

type sql struct {
	commentUp   string
	commentDown string
}

func NewSQL() Generator {
	return &sql{
		commentUp:   scanner.DefaultCommentUp,
		commentDown: scanner.DefaultCommentDown,
	}
}

func (s *sql) Generate(name string, path string) error {
	filename := pathLib.Join(path, fmt.Sprintf("%s.sql", name))

	return ioutil.WriteFile(filename, []byte(fmt.Sprintf("%s\n\n%s\n", s.commentUp, s.commentDown)), 0600)
}
