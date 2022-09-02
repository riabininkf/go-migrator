package cmd

import (
	"errors"
	"fmt"

	"github.com/riabininkf/go-migrator/internal/config"
	"github.com/riabininkf/go-migrator/pkg/migrator"
	"github.com/spf13/cobra"
)

func Up() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "up",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			var cnf *config.Config
			if cnf, err = config.New("GOMIGRATOR", cmd); err != nil {
				return fmt.Errorf("can't init configs: %w", err)
			}

			if err = cnf.BindPFlag("path", cmd.Flags().Lookup("path")); err != nil {
				return fmt.Errorf("can't bind flag \"path\" to config: %w", err)
			}

			var path string
			if path = cnf.GetString("path"); len(path) == 0 {
				return errors.New("path to migrations directory is required")
			}

			var m migrator.Migrator
			if m, err = newMigrator(cmd, cnf); err != nil {
				return fmt.Errorf("can't create migrator: %w", err)
			}

			return m.Up(path)
		},
	}

	cmd.Flags().String("db_dsn", "", "Database dsn")
	cmd.Flags().String("config", "", "Path to config file")
	cmd.Flags().String("path", "", "Path to directory with migrations")

	return cmd
}
