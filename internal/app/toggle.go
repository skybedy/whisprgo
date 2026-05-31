package app

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"whisprgo/internal/state"
)

func (a *App) newToggleCommand() *cobra.Command {
	var noTranscribe bool

	cmd := &cobra.Command{
		Use:   "toggle",
		Short: "Start or stop recording",
		RunE: func(cmd *cobra.Command, args []string) error {
			a.logInfof(cmd.ErrOrStderr(), "toggle invoked")
			if !state.Exists() {
				if err := a.recorder.CleanupOrphans(); err != nil {
					a.logErrorf(cmd.ErrOrStderr(), "orphan recorder cleanup failed: %v", err)
				}
				session, err := a.recorder.Start()
				if err != nil {
					a.logErrorf(cmd.ErrOrStderr(), "recording start failed: %v", err)
					a.notify("WhisprGo", "Nepodarilo se spustit nahravani.", "dialog-error", cmd.ErrOrStderr())
					return err
				}

				notificationID := 0
				if a.notifier != nil {
					id, err := a.notifier.NotifyWithOptions("WhisprGo", "Nahravani bezi.", "media-record", 0, 0, true)
					if err != nil {
						a.logErrorf(cmd.ErrOrStderr(), "recording notification failed: %v", err)
					}
					notificationID = id
				}

				s := state.State{
					Recording:      true,
					PID:            session.PID,
					AudioPath:      session.AudioPath,
					StartedAt:      session.StartedAt,
					NotificationID: notificationID,
				}

				if err := state.Save(s); err != nil {
					_ = a.recorder.Stop(session.PID)
					a.logErrorf(cmd.ErrOrStderr(), "state save failed after start: %v", err)
					return fmt.Errorf("failed to save state: %w", err)
				}

				a.logInfof(cmd.ErrOrStderr(), "recording started pid=%d audio=%s", s.PID, s.AudioPath)
				fmt.Fprintln(cmd.OutOrStdout(), "recording started")
				return nil
			}

			s, err := state.Load()
			if err != nil {
				if os.IsNotExist(err) {
					return nil
				}
				a.logErrorf(cmd.ErrOrStderr(), "state load failed: %v", err)
				a.notify("WhisprGo", "Nepodarilo se nacist stav nahravani.", "dialog-error", cmd.ErrOrStderr())
				return fmt.Errorf("failed to load state: %w", err)
			}

			if !s.Recording || s.Processing {
				a.logInfof(cmd.ErrOrStderr(), "toggle ignored while transcription is already processing audio=%s", s.AudioPath)
				a.notify("WhisprGo", "Prepis jeste probiha.", "dialog-information", cmd.ErrOrStderr())
				return nil
			}

			if a.recorder.IsRunning(s.PID) {
				if err := a.recorder.Stop(s.PID); err != nil {
					a.logErrorf(cmd.ErrOrStderr(), "recording stop failed pid=%d: %v", s.PID, err)
					a.notify("WhisprGo", "Nepodarilo se zastavit nahravani.", "dialog-error", cmd.ErrOrStderr())
					return err
				}
			}

			if err := a.recorder.WaitForFile(s.AudioPath, 5*time.Second); err != nil {
				a.logErrorf(cmd.ErrOrStderr(), "recording file not ready audio=%s: %v", s.AudioPath, err)
				a.notify("WhisprGo", "Nahravka nebyla dokoncena vcas.", "dialog-error", cmd.ErrOrStderr())
				return err
			}

			a.logInfof(cmd.ErrOrStderr(), "recording stopped pid=%d audio=%s", s.PID, s.AudioPath)
			if a.notifier != nil {
				_, err := a.notifier.NotifyWithOptions("WhisprGo", "Nahravani zastaveno.", "media-playback-stop", 1800, s.NotificationID, false)
				if err != nil {
					a.logErrorf(cmd.ErrOrStderr(), "stop notification failed: %v", err)
				}
			}
			fmt.Fprintln(cmd.OutOrStdout(), "recording stopped")
			fmt.Fprintf(cmd.OutOrStdout(), "audio: %s\n", s.AudioPath)

			if noTranscribe {
				if err := state.Delete(); err != nil {
					a.logErrorf(cmd.ErrOrStderr(), "state delete failed: %v", err)
					return fmt.Errorf("failed to delete state: %w", err)
				}
				return nil
			}

			s.Recording = false
			s.Processing = true
			if err := state.Save(s); err != nil {
				a.logErrorf(cmd.ErrOrStderr(), "state save failed before transcription: %v", err)
				return fmt.Errorf("failed to save state: %w", err)
			}
			defer func() {
				if err := state.Delete(); err != nil {
					a.logErrorf(cmd.ErrOrStderr(), "state delete failed after transcription: %v", err)
				}
			}()

			text, err := a.transcribeAudio(cmd.Context(), s.AudioPath)
			if err != nil {
				a.logErrorf(cmd.ErrOrStderr(), "auto-transcribe failed audio=%s: %v", s.AudioPath, err)
				a.notify("WhisprGo", "Prepis se nezdaril.", "dialog-error", cmd.ErrOrStderr())
				return fmt.Errorf("failed to transcribe %s: %w", s.AudioPath, err)
			}

			cleaned, err := a.maybeCleanupText(cmd.Context(), text)
			if err != nil {
				a.logErrorf(cmd.ErrOrStderr(), "cleanup failed, using raw transcription: %v", err)
			} else {
				text = cleaned
			}

			if err := a.maybeOutputText(text); err != nil {
				a.logErrorf(cmd.ErrOrStderr(), "output handling failed: %v", err)
				return fmt.Errorf("failed to output transcription: %w", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), text)
			return nil
		},
	}

	cmd.Flags().BoolVar(&noTranscribe, "no-transcribe", false, "stop recording without automatic transcription")
	return cmd
}
