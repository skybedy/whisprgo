package app

import (
	"fmt"

	"github.com/spf13/cobra"
	"whisprgo/internal/control"
)

func (a *App) newCancelCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "cancel",
		Short: "Cancel recording",
		RunE: func(cmd *cobra.Command, args []string) error {
			if a.shouldUseServeControl() {
				resp, err := control.SendCancel(cmd.Context())
				if err != nil {
					if !a.isServeReachable(cmd.Context()) {
						return serveNotRunningError()
					}
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), resp.Message)
				return nil
			}
			message, err := a.performCancelLocal(cmd.ErrOrStderr())
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), message)
			return nil
		},
	}
}
