package generator_test

import (
	"os"
	"testing"

	"github.com/riabininkf/go-migrator/pkg/generator"
	"github.com/stretchr/testify/assert"
)

func TestSql_Generate(t *testing.T) {
	sql := generator.NewSQL()
	defer func() {
		assert.NoError(t, os.Remove("test.sql"))
	}()

	err := sql.Generate("test", "./")
	assert.FileExists(t, "test.sql")
	assert.NoError(t, err)
}
