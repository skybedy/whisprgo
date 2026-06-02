package app

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/spf13/cobra"

	"whisprgo/internal/control"
)

func (a *App) newToggleCommand() *cobra.Command {
	var noTranscribe bool
	var forceCleanup bool

	cmd := &cobra.Command{
		Use:   "toggle",
		Short: "Start or stop recording",
		RunE: func(cmd *cobra.Command, args []string) error {
			if a.shouldUseServeControl() {
				return a.remoteToggle(cmd.Context(), cmd.OutOrStdout(), noTranscribe, forceCleanup)
			}
			return a.handleToggleLocal(cmd.Context(), cmd.OutOrStdout(), cmd.ErrOrStderr(), noTranscribe, forceCleanup)
		},
	}

	cmd.Flags().BoolVar(&noTranscribe, "no-transcribe", false, "stop recording without automatic transcription")
	cmd.Flags().BoolVar(&forceCleanup, "cleanup", false, "enable cleanup only for this run")
	return cmd
}

func (a *App) remoteToggle(ctx context.Context, outWriter io.Writer, noTranscribe bool, forceCleanup bool) error {
	a.logInfof(nil, "toggle delegated to serve")
	resp, err := control.SendToggle(ctx, noTranscribe, forceCleanup)
	if err != nil {
		if !a.isServeReachable(ctx) {
			return serveNotRunningError()
		}
		return err
	}
	if resp.Message != "" {
		fmt.Fprintln(outWriter, resp.Message)
	}
	if resp.Audio != "" {
		fmt.Fprintf(outWriter, "audio: %s\n", resp.Audio)
	}
	if resp.Text != "" {
		fmt.Fprintln(outWriter, resp.Text)
	}
	return nil
}

func (a *App) isServeReachable(ctx context.Context) bool {
	checkCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return control.IsServeReachable(checkCtx)
}
