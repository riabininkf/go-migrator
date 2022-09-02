package scanner

//go:generate mockery -case=underscore -output=. -inpkg -name=Scanner

import (
	"github.com/riabininkf/go-migrator/pkg/migration"
)

type Scanner interface {
	Scan(path string) ([]migration.Migration, error)
}
