package scanner

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/riabininkf/go-migrator/pkg/migration"
)

var ErrInvalidMigrationFormat = errors.New("invalid migration format")

const (
	DefaultCommentUp   = `-- migrate up`
	DefaultCommentDown = `-- migrate down`
)

type _sql struct {
	db          *sql.DB
	commentUp   string
	commentDown string
}

func NewSQL(db *sql.DB) Scanner {
	return &_sql{
		db:          db,
		commentUp:   DefaultCommentUp,
		commentDown: DefaultCommentDown,
	}
}

func (s *_sql) Scan(path string) ([]migration.Migration, error) {
	var err error

	var files map[string]string
	if files, err = loadFiles(path); err != nil {
		return nil, fmt.Errorf("can't load files from migrations directory: %w", err)
	}

	migrations := make([]migration.Migration, 0, len(files))
	for name, query := range files {
		query = strings.ReplaceAll(query, "\n", " ")
		parts := strings.Split(query, s.commentUp)
		if len(parts) != 2 {
			return nil, fmt.Errorf("can't read \"%s\" migration: %w", name, ErrInvalidMigrationFormat)
		}

		parts = strings.Split(parts[1], s.commentDown)
		if len(parts) != 2 {
			return nil, fmt.Errorf("can't read \"%s\" migration: %w", name, ErrInvalidMigrationFormat)
		}

		migrations = append(migrations, migration.NewMigration(name, parts[0], parts[1], s.db))
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Name() < migrations[j].Name()
	})

	return migrations, nil
}

func loadFiles(path string) (map[string]string, error) {
	var err error

	migrations := map[string]string{}
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if ext := filepath.Ext(path); ext != ".sql" {
			return nil
		}

		b, fileErr := ioutil.ReadFile(path)
		if fileErr != nil {
			return fmt.Errorf("can't read file %s: %w", path, err)
		}

		migrations[info.Name()] = string(b)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("can't scan migrations directory: %w", err)
	}

	return migrations, nil
}
