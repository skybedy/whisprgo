package app

import (
	"fmt"

	"github.com/spf13/cobra"

	"whisprgo/internal/state"
)

func (a *App) newCancelCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "cancel",
		Short: "Cancel recording",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !state.Exists() {
				fmt.Fprintln(cmd.OutOrStdout(), "nothing to cancel")
				return nil
			}

			if err := state.Delete(); err != nil {
				return fmt.Errorf("failed to delete state: %w", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), "recording cancelled")
			return nil
		},
	}
}
