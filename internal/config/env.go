package config

import (
	"os"
	"strconv"
	"strings"
)

func ApplyEnvOverrides(cfg *Config) {
	if cfg == nil {
		return
	}

	if v := strings.TrimSpace(os.Getenv("WHISPRGO_PROVIDER")); v != "" {
		cfg.Provider = v
	}
	if v := strings.TrimSpace(os.Getenv("WHISPRGO_TRANSCRIPTION_PROVIDER")); v != "" {
		cfg.Transcription.Provider = v
	}
	if v := strings.TrimSpace(os.Getenv("WHISPRGO_TRANSCRIPTION_MODEL")); v != "" {
		cfg.Transcription.Model = v
	}
	if v := strings.TrimSpace(os.Getenv("WHISPRGO_TRANSCRIPTION_LANGUAGE")); v != "" {
		cfg.Transcription.Language = v
	}
	if v := strings.TrimSpace(os.Getenv("WHISPRGO_TRANSCRIPTION_SHERPA_WS_URL")); v != "" {
		cfg.Transcription.SherpaWSURL = v
		cfg.Transcription.Parakeet.SherpaWSURL = v
	}
	if v := strings.TrimSpace(os.Getenv("WHISPRGO_TRANSCRIPTION_PARAKEET_MODE")); v != "" {
		cfg.Transcription.Parakeet.Mode = v
	}
	if v := strings.TrimSpace(os.Getenv("WHISPRGO_TRANSCRIPTION_PARAKEET_SHERPA_WS_URL")); v != "" {
		cfg.Transcription.Parakeet.SherpaWSURL = v
	}
	if v := strings.TrimSpace(os.Getenv("WHISPRGO_TRANSCRIPTION_PARAKEET_BINARY")); v != "" {
		cfg.Transcription.Parakeet.Binary = v
	}
	if v := strings.TrimSpace(os.Getenv("WHISPRGO_TRANSCRIPTION_PARAKEET_MODEL_DIR")); v != "" {
		cfg.Transcription.Parakeet.ModelDir = v
	}
	if v := strings.TrimSpace(os.Getenv("WHISPRGO_TRANSCRIPTION_PARAKEET_HOST")); v != "" {
		cfg.Transcription.Parakeet.Host = v
	}
	if v := strings.TrimSpace(os.Getenv("WHISPRGO_TRANSCRIPTION_PARAKEET_PORT")); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Transcription.Parakeet.Port = n
		}
	}
	if v := strings.TrimSpace(os.Getenv("WHISPRGO_TRANSCRIPTION_PARAKEET_NUM_THREADS")); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Transcription.Parakeet.NumThreads = n
		}
	}
	if v := strings.TrimSpace(os.Getenv("WHISPRGO_TRANSCRIPTION_PARAKEET_STARTUP_TIMEOUT_SECONDS")); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Transcription.Parakeet.StartupTimeoutSeconds = n
		}
	}
	if v := strings.TrimSpace(os.Getenv("WHISPRGO_TRANSCRIPTION_PARAKEET_REQUEST_TIMEOUT_SECONDS")); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Transcription.Parakeet.RequestTimeoutSeconds = n
		}
	}
	if v := strings.TrimSpace(os.Getenv("WHISPRGO_CLEANUP_MODEL")); v != "" {
		cfg.Cleanup.Model = v
	}
	if v := strings.TrimSpace(os.Getenv("WHISPRGO_CLEANUP_PROVIDER")); v != "" {
		cfg.Cleanup.Provider = v
	}
	if v := strings.TrimSpace(os.Getenv("WHISPRGO_OUTPUT_MODE")); v != "" {
		cfg.Output.Mode = v
	}
	if v := strings.TrimSpace(os.Getenv("WHISPRGO_CLEANUP_ENABLED")); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			cfg.Cleanup.Enabled = b
		}
	}
	if v := strings.TrimSpace(os.Getenv("WHISPRGO_LOGGING_INCLUDE_TEXT")); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			cfg.Logging.IncludeText = b
		}
	}
	if v := strings.TrimSpace(os.Getenv("WHISPRGO_LOG_FILE_PATH")); v != "" {
		cfg.Logging.FilePath = v
	}
}
