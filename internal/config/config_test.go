package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaultsWhenConfigDoesNotExist(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	def := Default()
	if cfg.Provider.Transcription != def.Provider.Transcription {
		t.Fatalf("unexpected provider.transcription: got %q want %q", cfg.Provider.Transcription, def.Provider.Transcription)
	}
	if cfg.Provider.TranscriptionModel != def.Provider.TranscriptionModel {
		t.Fatalf("unexpected provider.transcription_model: got %q want %q", cfg.Provider.TranscriptionModel, def.Provider.TranscriptionModel)
	}
	if cfg.Output.Clipboard != def.Output.Clipboard {
		t.Fatalf("unexpected output.clipboard: got %v want %v", cfg.Output.Clipboard, def.Output.Clipboard)
	}
}

func TestLoadMergesWithDefaults(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	cfgDir := filepath.Join(home, ".config", "whisprgo")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}

	content := []byte("provider:\n  transcription_model: whisper-1\ncleanup:\n  enabled: true\n")
	if err := os.WriteFile(filepath.Join(cfgDir, "config.yaml"), content, 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.Provider.Transcription != "openai" {
		t.Fatalf("expected default provider.transcription to remain openai, got %q", cfg.Provider.Transcription)
	}
	if cfg.Provider.TranscriptionModel != "whisper-1" {
		t.Fatalf("expected provider.transcription_model override, got %q", cfg.Provider.TranscriptionModel)
	}
	if !cfg.Cleanup.Enabled {
		t.Fatalf("expected cleanup.enabled override true")
	}
	if cfg.Output.Clipboard != true {
		t.Fatalf("expected default output.clipboard true")
	}
}
