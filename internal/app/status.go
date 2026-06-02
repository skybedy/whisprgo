package app

import (
	"fmt"

	"github.com/spf13/cobra"
	"whisprgo/internal/control"
)

func (a *App) newStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show recording status",
		RunE: func(cmd *cobra.Command, args []string) error {
			if a.shouldUseServeControl() {
				resp, err := control.SendStatus(cmd.Context())
				if err != nil {
					if !a.isServeReachable(cmd.Context()) {
						return serveNotRunningError()
					}
					return err
				}
				if resp.Message != "" {
					fmt.Fprintln(cmd.OutOrStdout(), resp.Message)
					return nil
				}
				fmt.Fprintln(cmd.OutOrStdout(), "recording")
				fmt.Fprintf(cmd.OutOrStdout(), "pid: %d\n", resp.PID)
				fmt.Fprintf(cmd.OutOrStdout(), "audio: %s\n", resp.Audio)
				fmt.Fprintf(cmd.OutOrStdout(), "started: %s\n", resp.Started)
				return nil
			}

			s, err := a.performStatusLocal(cmd.ErrOrStderr())
			if err != nil {
				return err
			}
			if !s.Recording {
				fmt.Fprintln(cmd.OutOrStdout(), "not recording")
				return nil
			}
			fmt.Fprintln(cmd.OutOrStdout(), "recording")
			fmt.Fprintf(cmd.OutOrStdout(), "pid: %d\n", s.PID)
			fmt.Fprintf(cmd.OutOrStdout(), "audio: %s\n", s.Audio)
			fmt.Fprintf(cmd.OutOrStdout(), "started: %s\n", s.Started)
			return nil
		},
	}
}
