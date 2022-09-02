package main_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/riabininkf/go-migrator/pkg/generator"
	"github.com/riabininkf/go-migrator/pkg/migrator"
	"github.com/riabininkf/go-migrator/pkg/registry"
	"github.com/riabininkf/go-migrator/pkg/scanner"
	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	r, err := registry.NewPostgres(registry.DefaultTableName, db)
	assert.NoError(t, err)

	m := migrator.New(scanner.NewSQL(db), generator.NewSQL(), r)

	var (
		name = "test_generating_migration"
		path = "./"
	)

	err = m.Generate(name, path)
	assert.NoError(t, err)

	migrations := make([]string, 0, 1)

	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if ext := filepath.Ext(path); ext != ".sql" {
			return nil
		}

		if strings.Contains(info.Name(), name) {
			migrations = append(migrations, info.Name())
			assert.NoError(t, os.Remove(path))
		}

		return nil
	})

	assert.NoError(t, err)
	assert.Len(t, migrations, 1)
}
