package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"whispergo/internal/state"
)

func (a *App) newStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show recording status",
		RunE: func(cmd *cobra.Command, args []string) error {
			a.logInfof(cmd.ErrOrStderr(), "status invoked")
			s, err := state.Load()
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Fprintln(cmd.OutOrStdout(), "not recording")
					return nil
				}
				a.logErrorf(cmd.ErrOrStderr(), "state load failed: %v", err)
				return fmt.Errorf("failed to load state: %w", err)
			}

			if !a.recorder.IsRunning(s.PID) {
				fmt.Fprintln(cmd.OutOrStdout(), "not recording")
				return nil
			}

			fmt.Fprintln(cmd.OutOrStdout(), "recording")
			fmt.Fprintf(cmd.OutOrStdout(), "pid: %d\n", s.PID)
			fmt.Fprintf(cmd.OutOrStdout(), "audio: %s\n", s.AudioPath)
			fmt.Fprintf(cmd.OutOrStdout(), "started: %s\n", s.StartedAt.Format("2006-01-02T15:04:05Z07:00"))
			return nil
		},
	}
}
