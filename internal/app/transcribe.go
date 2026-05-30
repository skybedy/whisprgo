package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func (a *App) newTranscribeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "transcribe /path/to/audio.wav",
		Short: "Transcribe an existing audio file (placeholder)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			audioPath := args[0]
			if _, err := os.Stat(audioPath); err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("audio file does not exist: %s", audioPath)
				}
				return fmt.Errorf("failed to access audio file: %w", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), "transcription is not implemented yet")
			fmt.Fprintf(cmd.OutOrStdout(), "audio: %s\n", audioPath)
			return nil
		},
	}
}
