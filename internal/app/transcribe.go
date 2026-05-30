package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func (a *App) newTranscribeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "transcribe /path/to/audio.wav",
		Short: "Transcribe an existing audio file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			audioPath := args[0]
			a.logInfof(cmd.ErrOrStderr(), "transcribe invoked audio=%s", audioPath)
			if _, err := os.Stat(audioPath); err != nil {
				if os.IsNotExist(err) {
					a.logErrorf(cmd.ErrOrStderr(), "audio file does not exist: %s", audioPath)
					return fmt.Errorf("audio file does not exist: %s", audioPath)
				}
				a.logErrorf(cmd.ErrOrStderr(), "failed to access audio file %s: %v", audioPath, err)
				return fmt.Errorf("failed to access audio file: %w", err)
			}

			text, err := a.transcribeAudio(cmd.Context(), audioPath)
			if err != nil {
				a.logErrorf(cmd.ErrOrStderr(), "transcribe failed audio=%s: %v", audioPath, err)
				return err
			}

			text, err = a.maybeCleanupText(cmd.Context(), text)
			if err != nil {
				a.logErrorf(cmd.ErrOrStderr(), "cleanup failed: %v", err)
				return err
			}

			if err := a.maybeOutputText(text); err != nil {
				a.logErrorf(cmd.ErrOrStderr(), "output handling failed: %v", err)
				return err
			}

			a.logInfof(cmd.ErrOrStderr(), "transcribe completed audio=%s", audioPath)
			fmt.Fprintln(cmd.OutOrStdout(), text)
			return nil
		},
	}
}
