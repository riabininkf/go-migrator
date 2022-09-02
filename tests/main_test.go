package main_test

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/riabininkf/go-migrator/pkg/registry"
	"github.com/riabininkf/go-migrator/pkg/scanner"
	"github.com/stretchr/testify/assert"
)

const dsn = "postgres://tester:password@localhost:5432/postgres?sslmode=disable"

var db *sql.DB

func TestMain(m *testing.M) {
	var err error
	if db, err = sql.Open("postgres", dsn); err != nil {
		fmt.Println("can't connect to database: " + err.Error())
		os.Exit(1)
	}
	defer func() {
		_ = db.Close()
	}()

	if err = db.Ping(); err != nil {
		fmt.Println("error on ping database: " + err.Error())
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func tableExists(t *testing.T, db *sql.DB, name string) bool {
	query := "select exists (select from information_schema.tables where table_name = $1)"
	var exists bool
	err := db.QueryRow(query, name).Scan(&exists)
	assert.NoError(t, err)

	return exists
}

func dropTable(t *testing.T, db *sql.DB, name string) {
	_, err := db.Exec(fmt.Sprintf("drop table %s", name))
	assert.NoError(t, err)
}

func deleteMigration(t *testing.T, db *sql.DB, name string) {
	qeury := fmt.Sprintf("delete from %s where name = $1", registry.DefaultTableName)
	_, err := db.Exec(qeury, name)
	assert.NoError(t, err)
}

func createMigration(t *testing.T, name, path, up, down string) {
	migration := fmt.Sprintf(
		"%s\n%s\n%s\n%s",
		scanner.DefaultCommentUp,
		up,
		scanner.DefaultCommentDown,
		down,
	)

	assert.NoError(t, ioutil.WriteFile(path+name, []byte(migration), 0600))
}
