package app

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"whisprgo/internal/cleanup"
	"whisprgo/internal/output"
	"whisprgo/internal/secrets"
	"whisprgo/internal/transcription"
)

func (a *App) transcribeAudio(ctx context.Context, audioPath string) (string, error) {
	a.logInfof(nil, "transcribe: start audio=%s", audioPath)
	if strings.TrimSpace(a.cfg.Provider) != "openai" {
		return "", fmt.Errorf("unsupported transcription provider: %s", a.cfg.Provider)
	}

	model := strings.TrimSpace(a.cfg.Transcription.Model)
	if model == "" {
		return "", fmt.Errorf("transcription.model is empty")
	}
	a.logInfof(nil, "transcribe: provider=%s model=%s", a.cfg.Provider, model)
	if src, err := secrets.SourceFor("openai"); err == nil {
		a.logInfof(nil, "transcribe: secret_source=%s", src)
	} else {
		a.logErrorf(nil, "transcribe: secret_source_unavailable reason=%v", err)
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
	a.logInfof(nil, "transcribe: timeout=%s", timeout)
	text, err := provider.Transcribe(reqCtx, audioPath, model)
	if err != nil {
		a.logErrorf(nil, "transcribe: failed err=%v", err)
		return "", err
	}
	a.logInfof(nil, "transcribe: done chars=%d", len(text))
	return text, nil
}

func (a *App) maybeCleanupText(ctx context.Context, input string) (string, error) {
	if !a.cfg.Cleanup.Enabled {
		a.logInfof(nil, "cleanup: disabled")
		return input, nil
	}
	a.logInfof(nil, "cleanup: enabled input_chars=%d", len(input))
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
	a.logInfof(nil, "cleanup: model=%s timeout=%s", model, timeout)
	out, err := cleaner.Clean(reqCtx, input, model, prompt)
	if err != nil {
		a.logErrorf(nil, "cleanup: failed err=%v", err)
		return "", err
	}
	a.logInfof(nil, "cleanup: done output_chars=%d", len(out))
	return out, nil
}

func (a *App) maybeOutputText(text string) error {
	mode := strings.ToLower(strings.TrimSpace(a.cfg.Output.Mode))
	shouldCopy := mode == "clipboard" || mode == "paste" || a.cfg.Output.CopyToClipboard
	shouldPaste := mode == "paste" || a.cfg.Output.PasteToActiveWindow
	a.logInfof(nil, "output: mode=%s copy=%t paste=%t text_chars=%d", mode, shouldCopy, shouldPaste, len(text))

	if shouldCopy {
		backend, err := output.CopyToClipboardDetailed(text)
		if err != nil {
			a.logErrorf(nil, "output: copy failed err=%v", err)
			return err
		}
		a.logInfof(nil, "output: copy ok backend=%s", backend)
	}

	if shouldPaste {
		a.logInfof(nil, "output: paste start")
		pasted, title, err := output.PasteWithXdotool(250*time.Millisecond, nil)
		if err != nil {
			a.logErrorf(nil, "output: paste failed err=%v", err)
			return err
		}
		if !pasted {
			a.logInfof(nil, "paste skipped for active window %q", title)
		} else {
			a.logInfof(nil, "output: paste ok active_window=%q", title)
		}
	}

	switch mode {
	case "stdout", "", "clipboard", "paste":
		return nil
	case "file":
		path := strings.TrimSpace(a.cfg.Output.FilePath)
		if path == "" {
			return fmt.Errorf("output.file_path is empty for output.mode=file")
		}
		if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
			a.logErrorf(nil, "output: file write failed path=%s err=%v", path, err)
			return err
		}
		a.logInfof(nil, "output: file write ok path=%s", path)
		return nil
	default:
		return fmt.Errorf("unsupported output.mode: %s", a.cfg.Output.Mode)
	}
}
