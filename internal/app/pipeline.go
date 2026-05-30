package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"whisprgo/internal/cleanup"
	"whisprgo/internal/output"
	"whisprgo/internal/transcription"
)

func (a *App) transcribeAudio(ctx context.Context, audioPath string) (string, error) {
	if strings.TrimSpace(a.cfg.Provider.Transcription) != "openai" {
		return "", fmt.Errorf("unsupported transcription provider: %s", a.cfg.Provider.Transcription)
	}

	model := strings.TrimSpace(a.cfg.Provider.TranscriptionModel)
	if model == "" {
		return "", fmt.Errorf("provider.transcription_model is empty")
	}

	provider, err := transcription.NewOpenAIProviderFromEnv(nil)
	if err != nil {
		return "", err
	}

	reqCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	return provider.Transcribe(reqCtx, audioPath, model)
}

func (a *App) maybeCleanupText(ctx context.Context, input string) (string, error) {
	if !a.cfg.Cleanup.Enabled {
		return input, nil
	}
	if strings.TrimSpace(a.cfg.Cleanup.Provider) != "openai" {
		return "", fmt.Errorf("unsupported cleanup provider: %s", a.cfg.Cleanup.Provider)
	}
	model := strings.TrimSpace(a.cfg.Cleanup.Model)
	if model == "" {
		return "", fmt.Errorf("cleanup.model is empty")
	}
	prompt := strings.TrimSpace(a.cfg.Cleanup.Prompt)
	if prompt == "" {
		return "", fmt.Errorf("cleanup.prompt is empty")
	}

	cleaner, err := cleanup.NewOpenAICleanerFromEnv(nil)
	if err != nil {
		return "", err
	}

	reqCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	return cleaner.Clean(reqCtx, input, model, prompt)
}

func (a *App) maybeOutputText(text string) error {
	if a.cfg.Output.Clipboard {
		if err := output.CopyToClipboard(text); err != nil {
			return err
		}
	}

	if a.cfg.Output.Paste {
		if err := output.PasteWithXdotool(); err != nil {
			return err
		}
	}

	return nil
}
