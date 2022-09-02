package main

import (
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/riabininkf/go-migrator/internal/cmd"
	"github.com/spf13/cobra"
)

func main() {
	root := cobra.Command{}
	root.AddCommand(
		cmd.Create(),
		cmd.Down(),
		cmd.Up(),
		cmd.Redo(),
		cmd.Status(),
		cmd.DBVersion(),
	)

	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
