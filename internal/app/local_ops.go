package app

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"whisprgo/internal/state"
)

type toggleResult struct {
	Message string
	Audio   string
	Text    string
}

type statusResult struct {
	Recording bool
	PID       int
	Audio     string
	Started   string
}

func (a *App) handleToggleLocal(ctx context.Context, out io.Writer, errOut io.Writer, noTranscribe bool) error {
	res, err := a.performToggleLocal(ctx, errOut, noTranscribe)
	if err != nil {
		return err
	}
	if res.Message != "" {
		fmt.Fprintln(out, res.Message)
	}
	if res.Audio != "" {
		fmt.Fprintf(out, "audio: %s\n", res.Audio)
	}
	if res.Text != "" {
		fmt.Fprintln(out, res.Text)
	}
	return nil
}

func (a *App) performToggleLocal(ctx context.Context, errOut io.Writer, noTranscribe bool) (toggleResult, error) {
	a.logInfof(errOut, "toggle invoked")
	if !state.Exists() {
		if err := a.recorder.CleanupOrphans(); err != nil {
			a.logErrorf(errOut, "orphan recorder cleanup failed: %v", err)
		}
		session, err := a.recorder.Start()
		if err != nil {
			a.logErrorf(errOut, "recording start failed: %v", err)
			a.notify("WhisprGo", "Nepodarilo se spustit nahravani.", "dialog-error", errOut)
			return toggleResult{}, err
		}

		notificationID := 0
		if a.notifier != nil {
			id, err := a.notifier.NotifyWithOptions("WhisprGo", "Nahravani bezi.", "media-record", 0, 0, true)
			if err != nil {
				a.logErrorf(errOut, "recording notification failed: %v", err)
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
			a.logErrorf(errOut, "state save failed after start: %v", err)
			return toggleResult{}, fmt.Errorf("failed to save state: %w", err)
		}

		a.logInfof(errOut, "recording started pid=%d audio=%s", s.PID, s.AudioPath)
		return toggleResult{Message: "recording started"}, nil
	}

	s, err := state.Load()
	if err != nil {
		if os.IsNotExist(err) {
			return toggleResult{}, nil
		}
		a.logErrorf(errOut, "state load failed: %v", err)
		a.notify("WhisprGo", "Nepodarilo se nacist stav nahravani.", "dialog-error", errOut)
		return toggleResult{}, fmt.Errorf("failed to load state: %w", err)
	}

	if !s.Recording || s.Processing {
		a.logInfof(errOut, "toggle ignored while transcription is already processing audio=%s", s.AudioPath)
		a.notify("WhisprGo", "Prepis jeste probiha.", "dialog-information", errOut)
		return toggleResult{}, nil
	}

	if a.recorder.IsRunning(s.PID) {
		if err := a.recorder.Stop(s.PID); err != nil {
			a.logErrorf(errOut, "recording stop failed pid=%d: %v", s.PID, err)
			a.notify("WhisprGo", "Nepodarilo se zastavit nahravani.", "dialog-error", errOut)
			return toggleResult{}, err
		}
	}

	if err := a.recorder.WaitForFile(s.AudioPath, 5*time.Second); err != nil {
		a.logErrorf(errOut, "recording file not ready audio=%s: %v", s.AudioPath, err)
		a.notify("WhisprGo", "Nahravka nebyla dokoncena vcas.", "dialog-error", errOut)
		return toggleResult{}, err
	}

	a.logInfof(errOut, "recording stopped pid=%d audio=%s", s.PID, s.AudioPath)
	if a.notifier != nil {
		_, err := a.notifier.NotifyWithOptions("WhisprGo", "Nahravani zastaveno.", "media-playback-stop", 1800, s.NotificationID, false)
		if err != nil {
			a.logErrorf(errOut, "stop notification failed: %v", err)
		}
	}

	result := toggleResult{Message: "recording stopped", Audio: s.AudioPath}
	if noTranscribe {
		if err := state.Delete(); err != nil {
			a.logErrorf(errOut, "state delete failed: %v", err)
			return toggleResult{}, fmt.Errorf("failed to delete state: %w", err)
		}
		return result, nil
	}

	s.Recording = false
	s.Processing = true
	if err := state.Save(s); err != nil {
		a.logErrorf(errOut, "state save failed before transcription: %v", err)
		return toggleResult{}, fmt.Errorf("failed to save state: %w", err)
	}
	defer func() {
		if err := state.Delete(); err != nil {
			a.logErrorf(errOut, "state delete failed after transcription: %v", err)
		}
	}()

	text, err := a.transcribeAudio(ctx, s.AudioPath)
	if err != nil {
		a.logErrorf(errOut, "auto-transcribe failed audio=%s: %v", s.AudioPath, err)
		a.notify("WhisprGo", "Prepis se nezdaril.", "dialog-error", errOut)
		return toggleResult{}, fmt.Errorf("failed to transcribe %s: %w", s.AudioPath, err)
	}

	cleaned, err := a.maybeCleanupText(ctx, text)
	if err != nil {
		a.logErrorf(errOut, "cleanup failed, using raw transcription: %v", err)
	} else {
		text = cleaned
	}

	if err := a.maybeOutputText(text); err != nil {
		a.logErrorf(errOut, "output handling failed: %v", err)
		return toggleResult{}, fmt.Errorf("failed to output transcription: %w", err)
	}

	result.Text = text
	return result, nil
}

func (a *App) performStatusLocal(errOut io.Writer) (statusResult, error) {
	a.logInfof(errOut, "status invoked")
	s, err := state.Load()
	if err != nil {
		if os.IsNotExist(err) {
			return statusResult{}, nil
		}
		a.logErrorf(errOut, "state load failed: %v", err)
		return statusResult{}, fmt.Errorf("failed to load state: %w", err)
	}
	if !a.recorder.IsRunning(s.PID) {
		return statusResult{}, nil
	}
	return statusResult{
		Recording: true,
		PID:       s.PID,
		Audio:     s.AudioPath,
		Started:   s.StartedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (a *App) performCancelLocal(errOut io.Writer) (string, error) {
	a.logInfof(errOut, "cancel invoked")
	s, err := state.Load()
	if err != nil {
		if os.IsNotExist(err) {
			return "nothing to cancel", nil
		}
		a.logErrorf(errOut, "state load failed: %v", err)
		return "", fmt.Errorf("failed to load state: %w", err)
	}
	if !a.recorder.IsRunning(s.PID) {
		_ = state.Delete()
		return "nothing to cancel", nil
	}
	if err := a.recorder.Stop(s.PID); err != nil {
		a.logErrorf(errOut, "recording stop failed pid=%d: %v", s.PID, err)
		return "", err
	}
	if err := state.Delete(); err != nil {
		a.logErrorf(errOut, "state delete failed: %v", err)
		return "", fmt.Errorf("failed to delete state: %w", err)
	}
	return "recording cancelled", nil
}
