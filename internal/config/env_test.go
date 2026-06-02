package config

import "testing"

func TestApplyEnvOverrides(t *testing.T) {
	cfg := Default()
	t.Setenv("WHISPRGO_PROVIDER", "openai")
	t.Setenv("WHISPRGO_TRANSCRIPTION_PROVIDER", "openai")
	t.Setenv("WHISPRGO_TRANSCRIPTION_MODEL", "whisper-1")
	t.Setenv("WHISPRGO_TRANSCRIPTION_LANGUAGE", "en")
	t.Setenv("WHISPRGO_TRANSCRIPTION_SHERPA_WS_URL", "ws://127.0.0.1:7000")
	t.Setenv("WHISPRGO_TRANSCRIPTION_PARAKEET_MODE", "managed")
	t.Setenv("WHISPRGO_TRANSCRIPTION_PARAKEET_BINARY", "/tmp/sherpa")
	t.Setenv("WHISPRGO_TRANSCRIPTION_PARAKEET_MODEL_DIR", "/tmp/model")
	t.Setenv("WHISPRGO_TRANSCRIPTION_PARAKEET_HOST", "127.0.0.1")
	t.Setenv("WHISPRGO_TRANSCRIPTION_PARAKEET_PORT", "6010")
	t.Setenv("WHISPRGO_TRANSCRIPTION_PARAKEET_NUM_THREADS", "5")
	t.Setenv("WHISPRGO_TRANSCRIPTION_PARAKEET_STARTUP_TIMEOUT_SECONDS", "11")
	t.Setenv("WHISPRGO_TRANSCRIPTION_PARAKEET_REQUEST_TIMEOUT_SECONDS", "121")
	t.Setenv("WHISPRGO_CLEANUP_ENABLED", "false")
	t.Setenv("WHISPRGO_CLEANUP_PROVIDER", "openai")
	t.Setenv("WHISPRGO_CLEANUP_MODEL", "gpt-5-mini")
	t.Setenv("WHISPRGO_OUTPUT_MODE", "clipboard")
	t.Setenv("WHISPRGO_LOGGING_INCLUDE_TEXT", "true")
	t.Setenv("WHISPRGO_LOG_FILE_PATH", "/tmp/whisprgo.log")

	ApplyEnvOverrides(&cfg)

	if cfg.Transcription.Language != "en" {
		t.Fatalf("expected transcription.language override")
	}
	if cfg.Transcription.Provider != "openai" {
		t.Fatalf("expected transcription.provider override")
	}
	if cfg.Transcription.SherpaWSURL != "ws://127.0.0.1:7000" {
		t.Fatalf("expected transcription.sherpa_ws_url override")
	}
	if cfg.Transcription.Parakeet.Mode != "managed" {
		t.Fatalf("expected transcription.parakeet.mode override")
	}
	if cfg.Transcription.Parakeet.Port != 6010 {
		t.Fatalf("expected transcription.parakeet.port override")
	}
	if cfg.Transcription.Parakeet.NumThreads != 5 {
		t.Fatalf("expected transcription.parakeet.num_threads override")
	}
	if cfg.Cleanup.Enabled != false {
		t.Fatalf("expected cleanup.enabled override")
	}
	if cfg.Cleanup.Provider != "openai" {
		t.Fatalf("expected cleanup.provider override")
	}
	if cfg.Output.Mode != "clipboard" {
		t.Fatalf("expected output.mode override")
	}
	if cfg.Logging.IncludeText != true {
		t.Fatalf("expected logging.include_text override")
	}
	if cfg.Logging.FilePath != "/tmp/whisprgo.log" {
		t.Fatalf("expected logging.file_path override")
	}
}
