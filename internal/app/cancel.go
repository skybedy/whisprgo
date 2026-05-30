package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"whisprgo/internal/state"
)

func (a *App) newCancelCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "cancel",
		Short: "Cancel recording",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := state.Load()
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Fprintln(cmd.OutOrStdout(), "nothing to cancel")
					return nil
				}
				return fmt.Errorf("failed to load state: %w", err)
			}

			if !a.recorder.IsRunning(s.PID) {
				_ = state.Delete()
				fmt.Fprintln(cmd.OutOrStdout(), "nothing to cancel")
				return nil
			}

			if err := a.recorder.Stop(s.PID); err != nil {
				return err
			}

			if err := state.Delete(); err != nil {
				return fmt.Errorf("failed to delete state: %w", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), "recording cancelled")
			return nil
		},
	}
}
