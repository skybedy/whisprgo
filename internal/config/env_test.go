package config

import "testing"

func TestApplyEnvOverrides(t *testing.T) {
	cfg := Default()
	t.Setenv("WHISPRGO_PROVIDER", "openai")
	t.Setenv("WHISPRGO_TRANSCRIPTION_MODEL", "whisper-1")
	t.Setenv("WHISPRGO_TRANSCRIPTION_LANGUAGE", "en")
	t.Setenv("WHISPRGO_CLEANUP_ENABLED", "false")
	t.Setenv("WHISPRGO_CLEANUP_MODEL", "gpt-5-mini")
	t.Setenv("WHISPRGO_OUTPUT_MODE", "clipboard")

	ApplyEnvOverrides(&cfg)

	if cfg.Transcription.Language != "en" {
		t.Fatalf("expected transcription.language override")
	}
	if cfg.Cleanup.Enabled != false {
		t.Fatalf("expected cleanup.enabled override")
	}
	if cfg.Output.Mode != "clipboard" {
		t.Fatalf("expected output.mode override")
	}
}
