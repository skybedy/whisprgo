package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"whisprgo/internal/state"
)

func (a *App) newToggleCommand() *cobra.Command {
	var noTranscribe bool

	cmd := &cobra.Command{
		Use:   "toggle",
		Short: "Start or stop recording",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !state.Exists() {
				session, err := a.recorder.Start()
				if err != nil {
					return err
				}

				s := state.State{Recording: true, PID: session.PID, AudioPath: session.AudioPath, StartedAt: session.StartedAt}

				if err := state.Save(s); err != nil {
					_ = a.recorder.Stop(session.PID)
					return fmt.Errorf("failed to save state: %w", err)
				}

				fmt.Fprintln(cmd.OutOrStdout(), "recording started")
				return nil
			}

			s, err := state.Load()
			if err != nil {
				if os.IsNotExist(err) {
					return nil
				}
				return fmt.Errorf("failed to load state: %w", err)
			}

			if a.recorder.IsRunning(s.PID) {
				if err := a.recorder.Stop(s.PID); err != nil {
					return err
				}
			}

			if err := state.Delete(); err != nil {
				return fmt.Errorf("failed to delete state: %w", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), "recording stopped")
			fmt.Fprintf(cmd.OutOrStdout(), "audio: %s\n", s.AudioPath)

			if noTranscribe {
				return nil
			}

			text, err := a.transcribeAudio(cmd.Context(), s.AudioPath)
			if err != nil {
				return fmt.Errorf("failed to transcribe %s: %w", s.AudioPath, err)
			}

			text, err = a.maybeCleanupText(cmd.Context(), text)
			if err != nil {
				return fmt.Errorf("failed to cleanup transcription: %w", err)
			}

			if err := a.maybeOutputText(text); err != nil {
				return fmt.Errorf("failed to output transcription: %w", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), text)
			return nil
		},
	}

	cmd.Flags().BoolVar(&noTranscribe, "no-transcribe", false, "stop recording without automatic transcription")
	return cmd
}
