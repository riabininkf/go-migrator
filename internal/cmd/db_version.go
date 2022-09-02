package cmd

import (
	"fmt"

	"github.com/riabininkf/go-migrator/internal/config"
	"github.com/riabininkf/go-migrator/pkg/migrator"
	"github.com/spf13/cobra"
)

func DBVersion() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "dbversion",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			var cnf *config.Config
			if cnf, err = config.New("GOMIGRATOR", cmd); err != nil {
				return fmt.Errorf("can't init configs: %w", err)
			}

			var m migrator.Migrator
			if m, err = newMigrator(cmd, cnf); err != nil {
				return fmt.Errorf("can't create migrator: %w", err)
			}

			return m.Version()
		},
	}

	cmd.Flags().String("db_dsn", "", "Database dsn")
	cmd.Flags().String("config", "", "Path to config file")

	return cmd
}
