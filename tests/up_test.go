package main_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/riabininkf/go-migrator/pkg/generator"
	"github.com/riabininkf/go-migrator/pkg/migrator"
	"github.com/riabininkf/go-migrator/pkg/registry"
	"github.com/riabininkf/go-migrator/pkg/scanner"
	"github.com/stretchr/testify/assert"
)

func TestUp(t *testing.T) {
	r, err := registry.NewPostgres(registry.DefaultTableName, db)
	assert.NoError(t, err)

	m := migrator.New(scanner.NewSQL(db), generator.NewSQL(), r)

	var (
		tableName     = fmt.Sprintf("test_up_migration_%d", time.Now().Second())
		migrationName = "create_table_" + tableName
		fileName      = migrationName + ".sql"
		path          = "./"
	)

	createMigration(t, fileName, path, "create table "+tableName+" ()", "")
	assert.False(t, tableExists(t, db, tableName))

	defer func() {
		assert.NoError(t, os.Remove(path+fileName))
		dropTable(t, db, tableName)
		deleteMigration(t, db, fileName)
	}()

	assert.NoError(t, m.Up(path))
	assert.True(t, tableExists(t, db, tableName))
}
