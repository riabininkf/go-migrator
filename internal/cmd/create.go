package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/riabininkf/go-migrator/pkg/generator"
	"github.com/spf13/cobra"
)

func Create() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "create",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			var path string
			if path, err = cmd.Flags().GetString("path"); err != nil {
				return fmt.Errorf("can't get flag \"path\": %w", err)
			}

			var name string
			if name, err = cmd.Flags().GetString("name"); err != nil {
				return fmt.Errorf("can't get flag \"name\": %w", err)
			}

			version := time.Now().UTC().Format("2006_01_02_15_04_05")
			fullName := fmt.Sprintf("%s_%s", version, strings.ReplaceAll(name, " ", "_"))

			return generator.NewSQL().Generate(path, fullName)
		},
	}

	cmd.Flags().String("path", "", "Path to folder with migrations")
	_ = cmd.MarkFlagRequired("path")

	cmd.Flags().String("name", "", "Name of migrations")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}
