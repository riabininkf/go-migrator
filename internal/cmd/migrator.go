package cmd

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/riabininkf/go-migrator/internal/config"
	"github.com/riabininkf/go-migrator/pkg/generator"
	"github.com/riabininkf/go-migrator/pkg/migrator"
	"github.com/riabininkf/go-migrator/pkg/registry"
	"github.com/riabininkf/go-migrator/pkg/scanner"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func newMigrator(cmd *cobra.Command, cnf *config.Config) (migrator.Migrator, error) {
	var err error

	if err = cnf.BindPFlag("db.dsn", cmd.Flags().Lookup("db_dsn")); err != nil {
		return nil, fmt.Errorf("can't bind flag \"db_dsn\" to config: %w", err)
	}

	var dsn string
	if dsn = cnf.GetString("db.dsn"); len(dsn) == 0 {
		return nil, errors.New("database dsn is required")
	}

	var db *sql.DB
	if db, err = sql.Open("postgres", dsn); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	var logger *zap.Logger
	if logger, err = zap.NewDevelopment(); err != nil {
		return nil, fmt.Errorf("can't create logger: %w", err)
	}

	var r registry.Registry
	if r, err = registry.NewPostgres(registry.DefaultTableName, db); err != nil {
		return nil, fmt.Errorf("can't create registry: %w", err)
	}

	return migrator.New(scanner.NewSQL(db), generator.NewSQL(), r, migrator.WithLogger(logger)), nil
}
