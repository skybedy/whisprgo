package app

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"whisprgo/internal/transcription"
)

func (a *App) newTranscribeCommand() *cobra.Command {
	var forceCleanup bool

	cmd := &cobra.Command{
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

			transcribeStart := time.Now()
			text, err := a.transcribeAudio(cmd.Context(), audioPath)
			if err != nil {
				if errors.Is(err, transcription.ErrEmptyTranscript) {
					a.logInfof(cmd.ErrOrStderr(), "transcription: no speech detected, skipping cleanup and output")
					return nil
				}
				a.logErrorf(cmd.ErrOrStderr(), "transcribe failed audio=%s: %v", audioPath, err)
				return err
			}
			transcribeDuration := time.Since(transcribeStart)

			if strings.TrimSpace(text) == "" {
				a.logInfof(cmd.ErrOrStderr(), "transcription returned empty text, skipping cleanup and output")
				return nil
			}

			cleanupStart := time.Now()
			cleaned, err := a.maybeCleanupText(cmd.Context(), text, forceCleanup)
			cleanupDuration := time.Since(cleanupStart)
			cleanupRan := err == nil && (a.cfg.Cleanup.Enabled || forceCleanup)
			if err != nil {
				a.logErrorf(cmd.ErrOrStderr(), "cleanup failed, using raw transcription: %v", err)
			} else {
				text = cleaned
			}
			if cleanupRan {
				a.logInfof(cmd.ErrOrStderr(), "pipeline: transcription=%s cleanup=%s total=%s", transcribeDuration.Round(time.Millisecond), cleanupDuration.Round(time.Millisecond), (transcribeDuration+cleanupDuration).Round(time.Millisecond))
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

	cmd.Flags().BoolVar(&forceCleanup, "cleanup", false, "enable cleanup only for this run")
	return cmd
}
