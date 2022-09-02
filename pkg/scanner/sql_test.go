package scanner_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/riabininkf/go-migrator/pkg/scanner"
	"github.com/stretchr/testify/assert"
)

func Test_sql_Scan(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	var (
		up   = "select 'migrate up'"
		down = "select 'migrate down'"
	)

	mock.ExpectBegin()
	mock.ExpectExec(fmt.Sprintf(" %s ", up)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec(fmt.Sprintf(" %s ", down)).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	s := scanner.NewSQL(db)

	content := fmt.Sprintf("%s\n%s\n%s\n%s", scanner.DefaultCommentUp, up, scanner.DefaultCommentDown, down)
	err = ioutil.WriteFile("./test.sql", []byte(content), 0600)
	assert.NoError(t, err)
	defer os.Remove("./test.sql")

	migrations, err := s.Scan("./")
	assert.NoError(t, err)

	assert.Len(t, migrations, 1)
	assert.NoError(t, migrations[0].Up())
	assert.NoError(t, migrations[0].Down())
}
