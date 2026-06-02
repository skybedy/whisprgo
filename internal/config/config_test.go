package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPathUsesUserConfigDir(t *testing.T) {
	xdg := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", xdg)

	got := Path()
	want := filepath.Join(xdg, "whisprgo", "config.yaml")
	if got != want {
		t.Fatalf("unexpected path: got %q want %q", got, want)
	}
}

func TestLoadDefaultsWhenConfigDoesNotExist(t *testing.T) {
	xdg := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", xdg)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	def := Default()
	if cfg.Provider != def.Provider {
		t.Fatalf("unexpected provider: got %q want %q", cfg.Provider, def.Provider)
	}
	if cfg.Transcription.Model != def.Transcription.Model {
		t.Fatalf("unexpected transcription.model: got %q want %q", cfg.Transcription.Model, def.Transcription.Model)
	}
}

func TestLoadMergesWithDefaults(t *testing.T) {
	xdg := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", xdg)
	cfgPath := filepath.Join(xdg, "whisprgo", "config.yaml")
	if err := os.MkdirAll(filepath.Dir(cfgPath), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	content := []byte("provider: openai\ntranscription:\n  model: whisper-1\ncleanup:\n  enabled: false\n")
	if err := os.WriteFile(cfgPath, content, 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.Transcription.Model != "whisper-1" {
		t.Fatalf("expected transcription.model override, got %q", cfg.Transcription.Model)
	}
	if cfg.Transcription.Provider != "openai" {
		t.Fatalf("expected transcription.provider fallback to openai, got %q", cfg.Transcription.Provider)
	}
	if cfg.Cleanup.Provider != "openai" {
		t.Fatalf("expected cleanup.provider fallback to openai, got %q", cfg.Cleanup.Provider)
	}
	if cfg.Transcription.Language == "" {
		t.Fatalf("expected default transcription.language to remain set")
	}
}

func TestLoadRespectsSectionProviders(t *testing.T) {
	xdg := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", xdg)
	cfgPath := filepath.Join(xdg, "whisprgo", "config.yaml")
	if err := os.MkdirAll(filepath.Dir(cfgPath), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	content := []byte("provider: openai\ntranscription:\n  provider: mistral\ncleanup:\n  provider: openai\n")
	if err := os.WriteFile(cfgPath, content, 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.Transcription.Provider != "mistral" {
		t.Fatalf("expected transcription.provider override, got %q", cfg.Transcription.Provider)
	}
	if cfg.Cleanup.Provider != "openai" {
		t.Fatalf("expected cleanup.provider explicit value, got %q", cfg.Cleanup.Provider)
	}
}

func TestLoadMergesLegacyAndNestedParakeetConfig(t *testing.T) {
	xdg := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", xdg)
	cfgPath := filepath.Join(xdg, "whisprgo", "config.yaml")
	if err := os.MkdirAll(filepath.Dir(cfgPath), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	content := []byte("transcription:\n  provider: parakeet\n  sherpa_ws_url: ws://127.0.0.1:7000\n  parakeet:\n    mode: managed\n    binary: /tmp/sherpa\n    model_dir: /tmp/model\n")
	if err := os.WriteFile(cfgPath, content, 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.Transcription.Parakeet.Mode != "managed" {
		t.Fatalf("expected parakeet.mode override")
	}
	if cfg.Transcription.Parakeet.SherpaWSURL != "ws://127.0.0.1:7000" {
		t.Fatalf("expected parakeet.sherpa_ws_url merge, got %q", cfg.Transcription.Parakeet.SherpaWSURL)
	}
	if cfg.Transcription.Parakeet.Port != 6010 {
		t.Fatalf("expected default parakeet.port, got %d", cfg.Transcription.Parakeet.Port)
	}
}
