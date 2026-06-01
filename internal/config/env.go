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
