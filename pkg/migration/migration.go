package migration

//go:generate mockery -case=underscore -output=. -inpkg -name=Migration

import (
	"database/sql"
	"fmt"
)

type (
	Migration interface {
		Name() string
		Up() error
		Down() error
	}

	migration struct {
		db   *sql.DB
		name string
		up   string
		down string
	}
)

func NewMigration(name string, up string, down string, db *sql.DB) Migration {
	return &migration{
		db:   db,
		name: name,
		up:   up,
		down: down,
	}
}

func (m *migration) Name() string {
	return m.name
}

func (m *migration) Up() error {
	var err error

	var tx *sql.Tx
	if tx, err = m.db.Begin(); err != nil {
		return fmt.Errorf("can't begin transaction: %w", err)
	}

	if _, err := tx.Exec(m.up); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("error on applying migrations: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("can't commit transaction: %w", err)
	}

	return nil
}

func (m *migration) Down() error {
	var err error

	var tx *sql.Tx
	if tx, err = m.db.Begin(); err != nil {
		return fmt.Errorf("can't begin transaction: %w", err)
	}

	if _, err := tx.Exec(m.down); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("error on rolling back migrations: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("can't commit transaction: %w", err)
	}

	return nil
}
