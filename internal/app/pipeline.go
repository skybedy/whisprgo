package app

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"whisprgo/internal/cleanup"
	"whisprgo/internal/output"
	"whisprgo/internal/parakeet"
	"whisprgo/internal/secrets"
	"whisprgo/internal/transcription"
)

func (a *App) transcribeAudio(ctx context.Context, audioPath string) (string, error) {
	a.logInfof(nil, "transcribe: start audio=%s", audioPath)
	transcriptionProvider := strings.TrimSpace(a.cfg.Transcription.Provider)
	if transcriptionProvider == "" {
		transcriptionProvider = strings.TrimSpace(a.cfg.Provider)
	}

	model := strings.TrimSpace(a.cfg.Transcription.Model)
	if model == "" {
		return "", fmt.Errorf("transcription.model is empty")
	}
	if transcriptionProvider == "parakeet" {
		return a.transcribeWithParakeet(ctx, audioPath, model)
	}

	a.logInfof(nil, "transcribe: provider=%s model=%s", transcriptionProvider, model)
	var (
		transcriptionAPIKey string
		err                 error
	)
	var src secrets.Source
	transcriptionAPIKey, src, err = secrets.GetForRole("transcription", transcriptionProvider)
	if err == nil {
		a.logInfof(nil, "transcribe: secret_source=%s", src)
	} else {
		a.logErrorf(nil, "transcribe: secret_source_unavailable reason=%v", err)
	}

	var provider transcription.Provider
	switch transcriptionProvider {
	case "openai":
		provider, err = transcription.NewOpenAIProvider(transcriptionAPIKey, nil)
	case "mistral":
		provider, err = transcription.NewMistralProvider(transcriptionAPIKey, nil)
	default:
		return "", fmt.Errorf("unsupported transcription provider: %s", transcriptionProvider)
	}
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
	if a.cfg.Logging.IncludeText {
		a.logInfof(nil, "transcribe: text=%q", text)
	}
	return text, nil
}

func (a *App) transcribeWithParakeet(ctx context.Context, audioPath string, model string) (string, error) {
	mode := strings.ToLower(strings.TrimSpace(a.cfg.Transcription.Parakeet.Mode))
	if mode == "" {
		mode = "external"
	}
	a.logInfof(nil, "transcription provider=parakeet mode=%s model=%s", mode, model)

	switch mode {
	case "external":
		wsURL := strings.TrimSpace(a.cfg.Transcription.Parakeet.SherpaWSURL)
		if wsURL == "" {
			wsURL = strings.TrimSpace(a.cfg.Transcription.SherpaWSURL)
		}
		a.logInfof(nil, "transcription provider=parakeet mode=external sherpa_ws_url=%s", wsURL)
		provider, err := transcription.NewParakeetWSProvider(wsURL, a.cfg.Transcription.Language)
		if err != nil {
			return "", err
		}

		timeout := time.Duration(a.cfg.Transcription.TimeoutSeconds) * time.Second
		if timeout <= 0 {
			timeout = 120 * time.Second
		}
		reqCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		startedAt := time.Now()
		text, err := provider.Transcribe(reqCtx, audioPath, model)
		if err != nil {
			a.logErrorf(nil, "transcribe: failed err=%v", err)
			return "", err
		}
		text = a.retryParakeetIfNeeded(reqCtx, provider, audioPath, model, text)
		a.logInfof(nil, "parakeet transcription completed duration=%s text_chars=%d", time.Since(startedAt), len(text))
		if a.cfg.Logging.IncludeText {
			a.logInfof(nil, "transcribe: text=%q", text)
		}
		return text, nil
	case "managed":
		cfg := a.cfg.Transcription.Parakeet
		a.logInfof(nil, "transcription provider=parakeet mode=managed")
		a.logInfof(nil, "parakeet binary=%s", cfg.Binary)
		a.logInfof(nil, "parakeet model_dir=%s", cfg.ModelDir)
		a.logInfof(nil, "parakeet port=%d", cfg.Port)
		a.logInfof(nil, "parakeet num_threads=%d", cfg.NumThreads)
		a.logInfof(nil, "parakeet starting server")
		server, startupDuration, err := parakeet.StartManagedServer(ctx, cfg, parakeetAppLogger{app: a})
		if err != nil {
			a.logErrorf(nil, "parakeet start failed err=%v", err)
			return "", err
		}
		a.logInfof(nil, "parakeet server ready duration=%s", startupDuration)
		defer func() {
			a.logInfof(nil, "parakeet stopping server")
			if stopErr := server.Stop(); stopErr != nil {
				a.logErrorf(nil, "parakeet stop failed err=%v", stopErr)
			}
		}()

		provider, err := transcription.NewParakeetWSProvider(server.URL(), a.cfg.Transcription.Language)
		if err != nil {
			return "", err
		}
		timeout := time.Duration(cfg.RequestTimeoutSeconds) * time.Second
		if timeout <= 0 {
			timeout = 120 * time.Second
		}
		reqCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		startedAt := time.Now()
		text, err := provider.Transcribe(reqCtx, audioPath, model)
		if err != nil {
			a.logErrorf(nil, "transcribe: failed err=%v", err)
			return "", err
		}
		text = a.retryParakeetIfNeeded(reqCtx, provider, audioPath, model, text)
		a.logInfof(nil, "parakeet transcription completed duration=%s text_chars=%d", time.Since(startedAt), len(text))
		if a.cfg.Logging.IncludeText {
			a.logInfof(nil, "transcribe: text=%q", text)
		}
		return text, nil
	default:
		return "", fmt.Errorf("unsupported parakeet mode: %s", mode)
	}
}

func (a *App) retryParakeetIfNeeded(ctx context.Context, provider transcription.Provider, audioPath string, model string, text string) string {
	if shouldRetryParakeetForCzech(a.cfg.Transcription.Language, text) {
		a.logInfof(nil, "transcribe: parakeet retry triggered for language=%s", a.cfg.Transcription.Language)
		retryText, retryErr := provider.Transcribe(ctx, audioPath, model)
		if retryErr != nil {
			a.logErrorf(nil, "transcribe: parakeet retry failed err=%v", retryErr)
		} else {
			text = retryText
			a.logInfof(nil, "transcribe: parakeet retry completed chars=%d", len(text))
		}
	}
	return text
}

type parakeetAppLogger struct {
	app *App
}

func (l parakeetAppLogger) Infof(format string, args ...any) {
	if l.app != nil {
		l.app.logInfof(nil, format, args...)
	}
}

func (l parakeetAppLogger) Errorf(format string, args ...any) {
	if l.app != nil {
		l.app.logErrorf(nil, format, args...)
	}
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
	cleanupProvider := strings.TrimSpace(a.cfg.Cleanup.Provider)
	if cleanupProvider == "" {
		cleanupProvider = strings.TrimSpace(a.cfg.Provider)
	}
	cleanupAPIKey, _, err := secrets.GetForRole("cleanup", cleanupProvider)
	if err != nil {
		return "", err
	}

	var cleaner cleanup.Cleaner
	switch cleanupProvider {
	case "openai":
		cleaner, err = cleanup.NewOpenAICleaner(cleanupAPIKey, nil)
	case "gemini":
		cleaner, err = cleanup.NewGeminiCleaner(cleanupAPIKey, nil)
	case "deepseek":
		cleaner, err = cleanup.NewDeepSeekCleaner(cleanupAPIKey, nil)
	default:
		return "", fmt.Errorf("unsupported cleanup provider: %s", cleanupProvider)
	}
	if err != nil {
		return "", err
	}

	timeout := time.Duration(a.cfg.Cleanup.TimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 120 * time.Second
	}
	reqCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	a.logInfof(nil, "cleanup: provider=%s model=%s timeout=%s", cleanupProvider, model, timeout)
	out, err := cleaner.Clean(reqCtx, input, model, prompt)
	if err != nil {
		a.logErrorf(nil, "cleanup: failed err=%v", err)
		return "", err
	}
	a.logInfof(nil, "cleanup: done output_chars=%d", len(out))
	if a.cfg.Logging.IncludeText {
		a.logInfof(nil, "cleanup: text=%q", out)
	}
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
