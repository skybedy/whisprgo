package app

import (
	"fmt"

	"github.com/spf13/cobra"

	"whisprgo/internal/config"
)

func (a *App) newConfigCommand() *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Config helpers",
	}

	configCmd.AddCommand(&cobra.Command{
		Use:   "path",
		Short: "Print config path",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), config.DisplayPath())
		},
	})

	return configCmd
}
