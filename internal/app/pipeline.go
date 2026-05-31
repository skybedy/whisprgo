package app

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"whisprgo/internal/cleanup"
	"whisprgo/internal/output"
	"whisprgo/internal/transcription"
)

func (a *App) transcribeAudio(ctx context.Context, audioPath string) (string, error) {
	if strings.TrimSpace(a.cfg.Provider) != "openai" {
		return "", fmt.Errorf("unsupported transcription provider: %s", a.cfg.Provider)
	}

	model := strings.TrimSpace(a.cfg.Transcription.Model)
	if model == "" {
		return "", fmt.Errorf("transcription.model is empty")
	}

	provider, err := transcription.NewOpenAIProviderFromSecrets(nil)
	if err != nil {
		return "", err
	}

	timeout := time.Duration(a.cfg.Transcription.TimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 120 * time.Second
	}
	reqCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return provider.Transcribe(reqCtx, audioPath, model)
}

func (a *App) maybeCleanupText(ctx context.Context, input string) (string, error) {
	if !a.cfg.Cleanup.Enabled {
		return input, nil
	}
	model := strings.TrimSpace(a.cfg.Cleanup.Model)
	if model == "" {
		return "", fmt.Errorf("cleanup.model is empty")
	}
	prompt := strings.TrimSpace(a.cfg.Cleanup.Prompt)
	if prompt == "" {
		return "", fmt.Errorf("cleanup.prompt is empty")
	}

	cleaner, err := cleanup.NewOpenAICleanerFromSecrets(nil)
	if err != nil {
		return "", err
	}

	timeout := time.Duration(a.cfg.Cleanup.TimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 120 * time.Second
	}
	reqCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return cleaner.Clean(reqCtx, input, model, prompt)
}

func (a *App) maybeOutputText(text string) error {
	switch strings.ToLower(strings.TrimSpace(a.cfg.Output.Mode)) {
	case "stdout", "":
		return nil
	case "clipboard":
		return output.CopyToClipboard(text)
	case "file":
		path := strings.TrimSpace(a.cfg.Output.FilePath)
		if path == "" {
			return fmt.Errorf("output.file_path is empty for output.mode=file")
		}
		return os.WriteFile(path, []byte(text), 0o644)
	default:
		return fmt.Errorf("unsupported output.mode: %s", a.cfg.Output.Mode)
	}
}
