package app

import (
	"fmt"

	"github.com/spf13/cobra"

	"whispergo/internal/buildinfo"
)

func (a *App) newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "version: %s\n", buildinfo.Version)
			fmt.Fprintf(cmd.OutOrStdout(), "commit: %s\n", buildinfo.Commit)
			fmt.Fprintf(cmd.OutOrStdout(), "date: %s\n", buildinfo.Date)
		},
	}
}
